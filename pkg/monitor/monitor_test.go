package monitor

import (
	"context"
	"errors"
	"testing"

	"github.com/timdorr/wattsup/pkg/config"
	"github.com/timdorr/wattsup/pkg/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Mock ModbusClient
type mockModbusClient struct {
	connectErr           error
	closeErr             error
	readHoldingRegsErr   error
	readHoldingRegsValue []byte
}

func (m *mockModbusClient) Connect() error {
	return m.connectErr
}

func (m *mockModbusClient) Close() error {
	return m.closeErr
}

func (m *mockModbusClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	if m.readHoldingRegsErr != nil {
		return nil, m.readHoldingRegsErr
	}
	return m.readHoldingRegsValue, nil
}

// Mock DBTX
type mockDBTX struct {
	copyFromErr error
	copyFromN   int64
	metrics     []sql.CreateMetricsParams
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return nil
}

func (m *mockDBTX) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	if m.copyFromErr != nil {
		return 0, m.copyFromErr
	}

	// This is a simplified mock. A more complete mock would iterate through rowSrc
	// and store the data. For this test, we'll just check that it's called.
	return m.copyFromN, nil
}

func TestReadAndStore(t *testing.T) {
	registers := []config.Register{
		{Name: "Reg1", Address: 100},
		{Name: "Reg2", Address: 101},
	}

	t.Run("success", func(t *testing.T) {
		client := &mockModbusClient{
			readHoldingRegsValue: []byte{0, 123}, // 123
		}
		db := &mockDBTX{}
		m := &Monitor{
			deviceName: "testDevice",
			deviceID:   1,
			registers:  registers,
			client:     client,
			db:         sql.New(db),
		}

		err := m.readAndStore()
		if err != nil {
			t.Fatalf("readAndStore() error = %v, wantErr nil", err)
		}
	})

	t.Run("modbus read error", func(t *testing.T) {
		client := &mockModbusClient{
			readHoldingRegsErr: errors.New("modbus error"),
		}
		db := &mockDBTX{}
		m := &Monitor{
			deviceName: "testDevice",
			deviceID:   1,
			registers:  registers,
			client:     client,
			db:         sql.New(db),
		}

		err := m.readAndStore()
		if err == nil {
			t.Fatal("readAndStore() error = nil, wantErr not nil")
		}
	})

	t.Run("db copyfrom error", func(t *testing.T) {
		client := &mockModbusClient{
			readHoldingRegsValue: []byte{0, 123},
		}
		db := &mockDBTX{
			copyFromErr: errors.New("db error"),
		}
		m := &Monitor{
			deviceName: "testDevice",
			deviceID:   1,
			registers:  registers,
			client:     client,
			db:         sql.New(db),
		}

		err := m.readAndStore()
		if err == nil {
			t.Fatal("readAndStore() error = nil, wantErr not nil")
		}
	})
}

func TestNewMonitor(t *testing.T) {
	t.Run("nil db pool", func(t *testing.T) {
		client := &mockModbusClient{}
		_, err := NewMonitor("testDevice", 1, []config.Register{}, client, nil)
		if err == nil {
			t.Fatal("NewMonitor() error = nil, wantErr not nil")
		}
	})
}
