package main

import (
	"kafkaing/logging"
	"kafkaing/management"
	"time"

	"go.uber.org/zap"
)

func init() {
	logging.Initialize()
}

func printBytes(b []byte) {
	zap.L().Sugar().Debug("Read : %s\n", string(b))
}

func main() {
	cm := management.GenerateNewCm("mcswirl", 0)
	producerChannel := cm.EstablishConnection("tcp", "localhost:9092", printBytes)

	producerChannel <- []byte("mario")
	producerChannel <- []byte("luigi")

	time.Sleep(3 * time.Second)

	cm.ShutdownHook()
}

// func main() {
// 	management.TestConnection()
// }
