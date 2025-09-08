package monitor

import (
	"context"
	"os"
	"time"

	"github.com/timdorr/wattsup/pkg/config"
	"github.com/timdorr/wattsup/pkg/sql"

	"github.com/apex/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Monitor struct {
	deviceName   string
	deviceID     int
	portFileName string
	registers    []config.Register
	handler      ModbusHandler
	client       ModbusClient
	db           DatabaseQuerier
}

func NewMonitor(deviceName, portFileName string, id int, registers []config.Register, database string) *Monitor {
	handler := newHandler(portFileName, id)
	client := newClient(handler)

	pool, err := pgxpool.New(context.Background(), database)
	if err != nil {
		log.WithError(err).Error("Failed to create database pool")
		os.Exit(1)
	}

	return &Monitor{
		deviceName:   deviceName,
		deviceID:     id,
		portFileName: portFileName,
		handler:      &modbusHandlerWrapper{handler: handler},
		client:       &modbusClientWrapper{client: client},
		registers:    registers,
		db:           sql.New(pool),
	}
}

// NewMonitorWithDeps creates a new monitor with injected dependencies (for testing)
func NewMonitorWithDeps(deviceName, portFileName string, id int, registers []config.Register, handler ModbusHandler, client ModbusClient, db DatabaseQuerier) *Monitor {
	return &Monitor{
		deviceName:   deviceName,
		deviceID:     id,
		portFileName: portFileName,
		handler:      handler,
		client:       client,
		registers:    registers,
		db:           db,
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	log.WithField("name", m.deviceName).WithField("device", m.portFileName).Info("Starting monitor...")

	for {
		err := m.handler.Connect()
		if err != nil {
			log.WithError(err).WithField("device", m.deviceName).Error("Failed to connect")

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
				continue
			}
		}

		err = m.watch(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(2 * time.Second):
			// Continue loop
		}
	}
}

func (m *Monitor) watch(ctx context.Context) error {
	interval := 1 * time.Second
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			log.WithField("device", m.deviceName).Info("Stopping monitor...")
			m.handler.Close()
			return nil
		case <-tick.C:
			err := m.readAndStore()
			if err != nil {
				return err
			}
		}
	}
}
