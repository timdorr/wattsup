package monitor

import (
	"context"
	"github.com/goburrow/modbus"
	"github.com/timdorr/wattsup/pkg/sql"
)

// ModbusClient interface for testing
type ModbusClient interface {
	ReadHoldingRegisters(address, quantity uint16) ([]byte, error)
}

// ModbusHandler interface for testing
type ModbusHandler interface {
	Connect() error
	Close() error
}

// DatabaseQuerier interface for testing
type DatabaseQuerier interface {
	CreateMetrics(ctx context.Context, arg []sql.CreateMetricsParams) (int64, error)
}

// Wrapper for modbus client to implement our interface
type modbusClientWrapper struct {
	client modbus.Client
}

func (w *modbusClientWrapper) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	return w.client.ReadHoldingRegisters(address, quantity)
}

// Wrapper for modbus handler to implement our interface
type modbusHandlerWrapper struct {
	handler *modbus.RTUClientHandler
}

func (w *modbusHandlerWrapper) Connect() error {
	return w.handler.Connect()
}

func (w *modbusHandlerWrapper) Close() error {
	return w.handler.Close()
}
