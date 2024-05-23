package logging

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

type Logit struct {
	logger *zap.SugaredLogger
	sigs   chan os.Signal
}

var LogInstance *Logit
var nonce *sync.Once = &sync.Once{}

func Access() *zap.SugaredLogger {
	if LogInstance == nil {
		nonce.Do(func() {
			logger, _ := zap.NewProduction()
			LogInstance = &Logit{logger: logger.Sugar(), sigs: make(chan os.Signal, 1)}
			LogInstance.logger.Info("Starting up logging")
			signal.Notify(LogInstance.sigs, syscall.SIGINT, syscall.SIGTERM)
			go LogInstance.runCleanup()
		})
	}
	return LogInstance.logger
}

func (logit *Logit) Info(template string, vars ...interface{}) {
	logit.logger.Infof(template, vars)
}

func (logit *Logit) Error(template string, vars ...interface{}) {
	logit.logger.Errorf(template, vars)
}

func (logit *Logit) runCleanup() {
	<-logit.sigs
	logit.logger.Info("Shutting down logging")
	logit.logger.Sync()
}
