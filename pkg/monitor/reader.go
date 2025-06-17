package monitor

import (
	"encoding/binary"

	"github.com/apex/log"
)

func (m *Monitor) read() error {
	for _, reg := range m.registers {
		result, err := m.client.ReadHoldingRegisters(reg.Address, 1)
		if err != nil {
			log.WithError(err).WithField("device", m.deviceName).WithField("register", reg.Name).Error("Failed to read register")
			return err
		}

		value := binary.BigEndian.Uint16(result)

		log.WithField("device", m.deviceName).WithField("register", reg.Name).Infof("Read value: %d", value)
	}

	return nil
}
