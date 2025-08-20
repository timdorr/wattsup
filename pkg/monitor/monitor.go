package monitor

import (
	"context"
	"os"
	"time"

	"github.com/timdorr/wattsup/pkg/config"

	"github.com/apex/log"
	"github.com/goburrow/modbus"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Monitor struct {
	deviceName   string
	deviceID     int
	portFileName string
	registers    []config.Register
	handler      *modbus.RTUClientHandler
	client       modbus.Client
	db           *pgxpool.Pool
}

func NewMonitor(deviceName, portFileName string, id int, registers []config.Register, database string) *Monitor {
	handler := newHandler(portFileName, id)

	pool, err := pgxpool.New(context.Background(), database)
	if err != nil {
		log.WithError(err).Error("Failed to create database pool")
		os.Exit(1)
	}

	return &Monitor{
		deviceName:   deviceName,
		deviceID:     id,
		portFileName: portFileName,
		handler:      handler,
		client:       newClient(handler),
		registers:    registers,
		db:           pool,
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	log.WithField("name", m.deviceName).WithField("device", m.portFileName).Info("Starting monitor...")

	for {
		err := m.handler.Connect()
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
			m.handler.Close()
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
