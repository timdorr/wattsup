package cmd

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSetupSignalHandler(t *testing.T) {
	t.Run("panics on second call", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		// a nil channel will block forever, so we need to reset it
		onlyOneSignalHandler = make(chan struct{})
		setupSignalHandler()
		setupSignalHandler()
	})
}

func TestHandleSignals(t *testing.T) {
	t.Run("cancels context on SIGTERM", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan os.Signal, 1)

		go handleSignals(cancel, c)

		c <- syscall.SIGTERM

		select {
		case <-ctx.Done():
			// success
		case <-time.After(1 * time.Second):
			t.Fatal("context was not canceled after 1s")
		}
	})
}
