package cmd

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSetupSignalHandler_Cancellation(t *testing.T) {
	// Reset the signal handler for testing
	onlyOneSignalHandler = make(chan struct{})

	ctx := setupSignalHandler()

	// Send interrupt signal to trigger cancellation
	go func() {
		time.Sleep(10 * time.Millisecond)
		proc, _ := os.FindProcess(os.Getpid())
		proc.Signal(syscall.SIGTERM)
	}()

	select {
	case <-ctx.Done():
		// Expected behavior
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should have been cancelled")
	}
}

func TestSetupSignalHandler_PanicsOnSecondCall(t *testing.T) {
	// Reset the signal handler for testing
	onlyOneSignalHandler = make(chan struct{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic on second call")
		}
	}()

	setupSignalHandler()
	setupSignalHandler() // Should panic
}
