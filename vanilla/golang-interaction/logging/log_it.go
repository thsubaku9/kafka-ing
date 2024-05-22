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

func (logit *Logit) Access() *zap.SugaredLogger {
	if logit == nil {
		nonce.Do(func() {
			logger, _ := zap.NewProduction()
			logit = &Logit{logger: logger.Sugar(), sigs: make(chan os.Signal, 1)}
			logit.logger.Info("Starting up logging")
			signal.Notify(logit.sigs, syscall.SIGINT, syscall.SIGTERM)
			go logit.runCleanup()
		})
	}
	return logit.logger
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
