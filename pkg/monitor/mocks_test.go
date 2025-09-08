package monitor

import (
	"context"
)

// Mock implementations for testing
type mockModbusClient struct {
	readFunc func(address, quantity uint16) ([]byte, error)
}

func (m *mockModbusClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	if m.readFunc != nil {
		return m.readFunc(address, quantity)
	}
	return []byte{0x00, 0x64}, nil // Default: value 100
}

type mockModbusHandler struct {
	connectFunc func() error
	closeFunc   func() error
	connected   bool
}

func (m *mockModbusHandler) Connect() error {
	if m.connectFunc != nil {
		err := m.connectFunc()
		if err == nil {
			m.connected = true
		}
		return err
	}
	m.connected = true
	return nil
}

func (m *mockModbusHandler) Close() error {
	if m.closeFunc != nil {
		err := m.closeFunc()
		if err == nil {
			m.connected = false
		}
		return err
	}
	m.connected = false
	return nil
}

type mockDatabaseQuerier struct {
	createMetricsFunc func(ctx context.Context, arg []CreateMetricsParams) (int64, error)
	callCount         int
	lastParams        []CreateMetricsParams
}

func (m *mockDatabaseQuerier) CreateMetrics(ctx context.Context, arg []CreateMetricsParams) (int64, error) {
	m.callCount++
	m.lastParams = arg
	if m.createMetricsFunc != nil {
		return m.createMetricsFunc(ctx, arg)
	}
	return int64(len(arg)), nil
}
