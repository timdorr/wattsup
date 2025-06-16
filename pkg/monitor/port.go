package monitor

import (
	"go.bug.st/serial"
)

func newPort(portFileName string) (serial.Port, error) {
	return serial.Open(portFileName, &serial.Mode{})
}
