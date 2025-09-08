package monitor

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/timdorr/wattsup/pkg/sql"

	"github.com/apex/log"
	"github.com/jackc/pgx/v5/pgtype"
)

// Import alias for testing
type CreateMetricsParams = sql.CreateMetricsParams

func (m *Monitor) readAndStore() error {
	var metrics []CreateMetricsParams

	for _, reg := range m.registers {
		result, err := m.client.ReadHoldingRegisters(reg.Address, 1)
		if err != nil {
			log.WithError(err).WithField("device", m.deviceName).WithField("register", reg.Name).Error("Failed to read register")
			return err
		}

		value := int16(binary.BigEndian.Uint16(result))

		log.WithField("device", m.deviceName).WithField("register", reg.Name).Infof("Read value: %d", value)

		metrics = append(metrics, CreateMetricsParams{
			Time:            pgtype.Timestamptz{Time: time.Now(), Valid: true},
			DeviceID:        pgtype.Int4{Int32: int32(m.deviceID), Valid: true},
			RegisterAddress: pgtype.Int4{Int32: int32(reg.Address), Valid: true},
			Value:           pgtype.Int4{Int32: int32(value), Valid: true},
		})
	}

	_, err := m.db.CreateMetrics(context.Background(), metrics)
	if err != nil {
		log.WithField("metrics", metrics).WithError(err).Error("Failed to insert metrics")
		return err
	}

	return nil
}
