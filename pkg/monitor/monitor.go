package monitor

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/apex/log"
	"github.com/goburrow/modbus"
)

type Monitor struct {
	deviceName   string
	portFileName string
	handler      *modbus.RTUClientHandler
	client       modbus.Client
}

func NewMonitor(deviceName string, portFileName string, slaveId int) *Monitor {
	handler := newHandler(portFileName, slaveId)

	return &Monitor{
		deviceName:   deviceName,
		portFileName: portFileName,
		handler:      handler,
		client:       newClient(handler),
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
			res, err := m.client.ReadHoldingRegisters(178, 1)
			if err != nil {
				log.WithError(err).WithField("device", m.deviceName).Error("Failed to read holding registers")
				return err
			}

			log.WithField("device", m.deviceName).Infof("Received data: %d", binary.BigEndian.Uint16(res))
		}
	}
}
