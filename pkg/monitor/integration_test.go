package monitor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/timdorr/wattsup/pkg/config"
)

func TestMonitor_Integration_ConnectionRetry(t *testing.T) {
	// Test that verifies connection retry logic without real timeouts
	connectAttempts := 0
	mockHandler := &mockModbusHandler{
		connectFunc: func() error {
			connectAttempts++
			return errors.New("connection failed")
		},
	}

	mockClient := &mockModbusClient{}
	mockDB := &mockDatabaseQuerier{}

	monitor := NewMonitorWithDeps("test-inverter", "/dev/ttyUSB0", 1, []config.Register{}, mockHandler, mockClient, mockDB)

	// Use a very short timeout so it fails quickly
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Start monitor in goroutine
	done := make(chan error, 1)
	go func() {
		done <- monitor.Start(ctx)
	}()

	// Wait for monitor to stop due to timeout
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Monitor did not stop within timeout")
	}

	// Verify that at least one connection attempt was made
	if connectAttempts == 0 {
		t.Error("Expected at least 1 connection attempt")
	}
}

func TestMonitor_HandleVariousErrors(t *testing.T) {
	tests := []struct {
		name           string
		client         ModbusClient
		db             DatabaseQuerier
		expectedError  bool
		expectedDBCall bool
	}{
		{
			name: "ModbusTimeout",
			client: &mockModbusClient{
				readFunc: func(address, quantity uint16) ([]byte, error) {
					return nil, errors.New("timeout")
				},
			},
			db:             &mockDatabaseQuerier{},
			expectedError:  true,
			expectedDBCall: false,
		},
		{
			name:   "DatabaseConnectionLost",
			client: &mockModbusClient{},
			db: &mockDatabaseQuerier{
				createMetricsFunc: func(ctx context.Context, arg []CreateMetricsParams) (int64, error) {
					return 0, errors.New("connection lost")
				},
			},
			expectedError:  true,
			expectedDBCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := &mockModbusHandler{}
			registers := []config.Register{
				{Name: "Test Reg", Address: 100},
			}

			monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, tt.client, tt.db)

			err := monitor.readAndStore()

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if mockDB, ok := tt.db.(*mockDatabaseQuerier); ok {
				dbCalled := mockDB.callCount > 0
				if tt.expectedDBCall != dbCalled {
					t.Errorf("Expected DB call %v, got %v", tt.expectedDBCall, dbCalled)
				}
			}
		})
	}
}
