package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func setupSignalHandler() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go handleSignals(cancel, c)

	return ctx
}

func handleSignals(cancel context.CancelFunc, c chan os.Signal) {
	<-c
	log.Info("Shutting down gracefully...")
	cancel()
	<-c
	log.Fatal("Forcing shutdown...")
	os.Exit(1) // second signal. Exit directly.
}
