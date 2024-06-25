package logging

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

var nonce *sync.Once = &sync.Once{}

func Initialize() {
	nonce.Do(func() {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
		zap.L().Sugar().Info("Starting up logging")

		exitChannel := make(chan os.Signal, 1)
		signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
		go runCleanup(exitChannel)
	})
}

func runCleanup(exitChannel <-chan os.Signal) {
	<-exitChannel
	zap.L().Sugar().Info("Shutting down logging")
	zap.L().Sugar().Sync()
}
