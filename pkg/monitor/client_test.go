package monitor

import (
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	handler := newHandler("/dev/ttyUSB0", 1)

	if handler.SlaveId != 1 {
		t.Errorf("Expected SlaveId 1, got %d", handler.SlaveId)
	}

	if handler.BaudRate != 9600 {
		t.Errorf("Expected BaudRate 9600, got %d", handler.BaudRate)
	}

	if handler.DataBits != 8 {
		t.Errorf("Expected DataBits 8, got %d", handler.DataBits)
	}

	if handler.Parity != "N" {
		t.Errorf("Expected Parity 'N', got %s", handler.Parity)
	}

	if handler.Timeout != 3*time.Second {
		t.Errorf("Expected Timeout 3s, got %v", handler.Timeout)
	}
}

func TestNewClient(t *testing.T) {
	handler := newHandler("/dev/ttyUSB0", 1)
	client := newClient(handler)

	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
}
