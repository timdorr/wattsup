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
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		err := m.client.Connect()
		if err != nil {
			log.WithError(err).WithField("device", m.deviceName).Error("Failed to connect")

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
			}

			continue
		}

		err = m.watch(ctx)
		if err == nil {
			return nil
		}

		log.WithError(err).WithField("device", m.deviceName).Warn("Monitor stopped unexpectedly, reconnecting...")

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(2 * time.Second):
		}
	}
}

func (m *Monitor) watch(ctx context.Context) error {
	defer m.client.Close()

	interval := 1 * time.Second
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			log.WithField("device", m.deviceName).Info("Stopping monitor...")
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
