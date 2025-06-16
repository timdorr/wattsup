package monitor

import (
	"context"

	"github.com/apex/log"
	"go.bug.st/serial"
)

type Monitor struct {
	deviceName   string
	portFileName string
	port         serial.Port
}

func NewMonitor(deviceName string, portFileName string) (*Monitor, error) {
	port, err := newPort(portFileName)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		deviceName:   deviceName,
		portFileName: portFileName,
		port:         port,
	}, nil
}

func (m *Monitor) Start(ctx context.Context) error {
	log.WithField("name", m.deviceName).WithField("device", m.portFileName).Info("Starting monitor...")

	buffer := make([]byte, 1024)

	bufferChan := make(chan []byte, 10)
	go m.watch(ctx, bufferChan)

	for {
		n, err := m.port.Read(buffer)
		if err != nil {
			return err
		}
		log.WithField("device", m.deviceName).Infof("Received %d bytes: %s", n, buffer[:n])
	}
}

func (m *Monitor) watch(ctx context.Context, bufferChan <-chan []byte) error {
	for {
		select {
		case <-ctx.Done():
			log.WithField("device", m.deviceName).Info("Stopping monitor...")
			m.port.Close()
			return nil
		case data := <-bufferChan:
			log.WithField("device", m.deviceName).Infof("Received buffered data: %s", data)
		}
	}
}
