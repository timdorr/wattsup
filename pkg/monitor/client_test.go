package monitor

import (
	"testing"
)

func TestNewModbusClient(t *testing.T) {
	t.Run("returns non-nil client", func(t *testing.T) {
		client := NewModbusClient("/dev/null", 1)
		if client == nil {
			t.Fatal("NewModbusClient() returned nil")
		}
	})
}
