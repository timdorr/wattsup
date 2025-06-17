package monitor

import (
	"context"
	"time"

	"github.com/timdorr/wattsup/pkg/config"

	"github.com/apex/log"
	"github.com/goburrow/modbus"
)

type Monitor struct {
	deviceName   string
	portFileName string
	registers    []config.Register
	handler      *modbus.RTUClientHandler
	client       modbus.Client
}

func NewMonitor(deviceName, portFileName string, id int, registers []config.Register) *Monitor {
	handler := newHandler(portFileName, id)

	return &Monitor{
		deviceName:   deviceName,
		portFileName: portFileName,
		handler:      handler,
		client:       newClient(handler),
		registers:    registers,
	}
}

func (m *Monitor) Start(ctx context.Context) error {
	log.WithField("name", m.deviceName).WithField("device", m.portFileName).Info("Starting monitor...")

	for {
		err := m.handler.Connect()
		if err != nil {
			log.WithError(err).WithField("device", m.deviceName).Error("Failed to connect")
			return err
		}

		err = m.watch(ctx)
		if err == nil {
			return nil
		}

		time.Sleep(5 * time.Second)
	}
}

func (m *Monitor) watch(ctx context.Context) error {
	tick := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.WithField("device", m.deviceName).Info("Stopping monitor...")
			m.handler.Close()
			return nil
		case <-tick.C:
			err := m.read()
			if err != nil {
				return err
			}
		}
	}
}
