package main

import (
	"fmt"
	connectionmanager "kafkaing/connection_manager"
	"kafkaing/logging"
	"time"
)

func init() {
	logging.Access()
}

func printBytes(b []byte) {
	fmt.Printf("Read : %s\n", string(b))
}

func main() {
	cm := connectionmanager.GenerateNewCm("mcswirl", 0)
	producerChannel := cm.EstablishConnection("tcp", "localhost:9092", printBytes)

	logging.LogInstance.Info("eh")
	producerChannel <- []byte("mario")
	producerChannel <- []byte("luigi")

	time.Sleep(3 * time.Second)

	cm.ShutdownHook()
}
