package monitor

import (
	"time"

	"github.com/goburrow/modbus"
)

type ModbusClientImpl struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func NewModbusClient(portFileName string, id int) ModbusClient {
	handler := modbus.NewRTUClientHandler(portFileName)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.Timeout = 3 * time.Second
	handler.SlaveId = byte(id)

	return &ModbusClientImpl{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

func (m *ModbusClientImpl) Connect() error {
	return m.handler.Connect()
}

func (m *ModbusClientImpl) Close() error {
	return m.handler.Close()
}

func (m *ModbusClientImpl) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return m.client.ReadHoldingRegisters(address, quantity)
}
