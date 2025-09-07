package monitor

import (
	"context"
	"errors"
	"time"

	"github.com/timdorr/wattsup/pkg/config"
	"github.com/timdorr/wattsup/pkg/sql"

	"github.com/apex/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ModbusClient interface {
	Connect() error
	Close() error
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
}

type Monitor struct {
	deviceName string
	deviceID   int
	registers  []config.Register
	client     ModbusClient
	db         *sql.Queries
}

func NewMonitor(deviceName string, id int, registers []config.Register, client ModbusClient, pool *pgxpool.Pool) (*Monitor, error) {
	if pool == nil {
		return nil, errors.New("database pool is nil")
	}
	return &Monitor{
		deviceName: deviceName,
		deviceID:   id,
		registers:  registers,
		client:     client,
		db:         sql.New(pool),
	}, nil
}

func (m *Monitor) Start(ctx context.Context) error {
	log.WithField("name", m.deviceName).Info("Starting monitor...")

	for {
		err := m.client.Connect()
		if err != nil {
			log.WithError(err).WithField("device", m.deviceName).Error("Failed to connect")

			time.Sleep(5 * time.Second)

			continue
		}

		err = m.watch(ctx)
		if err == nil {
			return nil
		}

		time.Sleep(2 * time.Second)
	}
}

func (m *Monitor) watch(ctx context.Context) error {
	interval := 1 * time.Second
	tick := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			log.WithField("device", m.deviceName).Info("Stopping monitor...")
			m.client.Close()
			return nil
		case <-tick.C:
			tick.Stop()

			err := m.readAndStore()
			if err != nil {
				return err
			}

			tick.Reset(interval)
		}
	}
}
