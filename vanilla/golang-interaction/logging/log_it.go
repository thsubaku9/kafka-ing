package logging

import (
	"kafkaing/utils"
	"sync"

	"go.uber.org/zap"
)

var nonce *sync.Once = &sync.Once{}

func Initialize() {
	nonce.Do(func() {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
		zap.L().Sugar().Info("Starting up logging")

		go utils.RunCleanup(utils.PrepareSigtermChannel(), func() {
			zap.L().Sugar().Info("Shutting down logging")
			zap.L().Sugar().Sync()
		})
	})
}
