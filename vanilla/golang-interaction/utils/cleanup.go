package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func PrepareSigtermChannel() <-chan os.Signal {
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
	return exitChannel
}

func RunCleanup(exitChannel <-chan os.Signal, d func()) {
	<-exitChannel
	d()
}
