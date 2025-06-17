package monitor

import (
	"time"

	"github.com/goburrow/modbus"
)

func newHandler(portFileName string, id int) *modbus.RTUClientHandler {
	handler := modbus.NewRTUClientHandler(portFileName)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.Timeout = 3 * time.Second
	handler.SlaveId = byte(id)

	return handler
}

func newClient(handler *modbus.RTUClientHandler) modbus.Client {
	client := modbus.NewClient(handler)
	return client
}
