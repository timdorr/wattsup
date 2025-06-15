package monitor

import (
	"go.bug.st/serial"
)

func ListDevices() ([]string, error) {
	return serial.GetPortsList()
}
