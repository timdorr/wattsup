package monitor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/timdorr/wattsup/pkg/config"
)

func TestNewMonitorWithDeps(t *testing.T) {
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{}
	mockDB := &mockDatabaseQuerier{}

	registers := []config.Register{
		{Name: "Test Reg", Address: 100},
	}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, mockClient, mockDB)

	if monitor.deviceName != "test-device" {
		t.Errorf("Expected device name 'test-device', got %s", monitor.deviceName)
	}

	if monitor.deviceID != 1 {
		t.Errorf("Expected device ID 1, got %d", monitor.deviceID)
	}

	if len(monitor.registers) != 1 {
		t.Errorf("Expected 1 register, got %d", len(monitor.registers))
	}
}

func TestMonitor_ReadAndStore_Success(t *testing.T) {
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{
		readFunc: func(address, quantity uint16) ([]byte, error) {
			if address == 100 {
				return []byte{0x01, 0x2C}, nil // Value 300
			}
			return []byte{0x00, 0x64}, nil // Value 100
		},
	}
	mockDB := &mockDatabaseQuerier{}

	registers := []config.Register{
		{Name: "Test Reg 1", Address: 100},
		{Name: "Test Reg 2", Address: 200},
	}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, mockClient, mockDB)

	err := monitor.readAndStore()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mockDB.callCount != 1 {
		t.Errorf("Expected 1 database call, got %d", mockDB.callCount)
	}

	if len(mockDB.lastParams) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(mockDB.lastParams))
	}

	// Check first metric (address 100, value 300)
	if mockDB.lastParams[0].RegisterAddress.Int32 != 100 {
		t.Errorf("Expected register address 100, got %d", mockDB.lastParams[0].RegisterAddress.Int32)
	}
	if mockDB.lastParams[0].Value.Int32 != 300 {
		t.Errorf("Expected value 300, got %d", mockDB.lastParams[0].Value.Int32)
	}

	// Check second metric (address 200, value 100)
	if mockDB.lastParams[1].RegisterAddress.Int32 != 200 {
		t.Errorf("Expected register address 200, got %d", mockDB.lastParams[1].RegisterAddress.Int32)
	}
	if mockDB.lastParams[1].Value.Int32 != 100 {
		t.Errorf("Expected value 100, got %d", mockDB.lastParams[1].Value.Int32)
	}
}

func TestMonitor_ReadAndStore_ModbusError(t *testing.T) {
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{
		readFunc: func(address, quantity uint16) ([]byte, error) {
			return nil, errors.New("modbus read error")
		},
	}
	mockDB := &mockDatabaseQuerier{}

	registers := []config.Register{
		{Name: "Test Reg", Address: 100},
	}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, mockClient, mockDB)

	err := monitor.readAndStore()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mockDB.callCount != 0 {
		t.Errorf("Expected 0 database calls, got %d", mockDB.callCount)
	}
}

func TestMonitor_ReadAndStore_DatabaseError(t *testing.T) {
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{}
	mockDB := &mockDatabaseQuerier{
		createMetricsFunc: func(ctx context.Context, arg []CreateMetricsParams) (int64, error) {
			return 0, errors.New("database error")
		},
	}

	registers := []config.Register{
		{Name: "Test Reg", Address: 100},
	}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, mockClient, mockDB)

	err := monitor.readAndStore()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mockDB.callCount != 1 {
		t.Errorf("Expected 1 database call, got %d", mockDB.callCount)
	}
}

func TestMonitor_Start_ConnectError(t *testing.T) {
	mockHandler := &mockModbusHandler{
		connectFunc: func() error {
			return errors.New("connection failed")
		},
	}
	mockClient := &mockModbusClient{}
	mockDB := &mockDatabaseQuerier{}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, []config.Register{}, mockHandler, mockClient, mockDB)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := monitor.Start(ctx)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestMonitor_Watch_ContextCancellation(t *testing.T) {
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{}
	mockDB := &mockDatabaseQuerier{}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, []config.Register{}, mockHandler, mockClient, mockDB)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := monitor.watch(ctx)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if mockHandler.connected {
		t.Error("Expected handler to be closed")
	}
}

func TestMonitor_Watch_ReadError(t *testing.T) {
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{
		readFunc: func(address, quantity uint16) ([]byte, error) {
			return nil, errors.New("read error")
		},
	}
	mockDB := &mockDatabaseQuerier{}

	registers := []config.Register{
		{Name: "Test Reg", Address: 100},
	}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, mockClient, mockDB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Allow some time for the tick to occur
	err := monitor.watch(ctx)
	
	// The watch function will return nil when context is cancelled, even if there were read errors
	// during execution, because context cancellation takes precedence
	if err != nil {
		t.Errorf("Expected nil error from context cancellation, got %v", err)
	}
}

func TestMonitor_ReadAndStore_ErrorBeforeContextCancel(t *testing.T) {
	// This test specifically checks that readAndStore returns an error
	// when the modbus read fails, separate from the watch context
	mockHandler := &mockModbusHandler{}
	mockClient := &mockModbusClient{
		readFunc: func(address, quantity uint16) ([]byte, error) {
			return nil, errors.New("read error")
		},
	}
	mockDB := &mockDatabaseQuerier{}

	registers := []config.Register{
		{Name: "Test Reg", Address: 100},
	}

	monitor := NewMonitorWithDeps("test-device", "/dev/test", 1, registers, mockHandler, mockClient, mockDB)

	// Test readAndStore directly
	err := monitor.readAndStore()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
