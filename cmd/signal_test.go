package cmd

import (
	"testing"
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
