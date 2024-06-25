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
// 	var avroManager management.AvroManager = management.InitAvroManager("sqlite3", "./schemaRegistry.db")

// 	avroManager.CreateSchema("ace", "v1", `{
// 		"type": "record",
// 		"name": "simple",
// 		"namespace": "org.hamba.avro",
// 		"fields" : [
// 			{"name": "a", "type": "long"},
// 			{"name": "b", "type": "string"}
// 		]
// 	}`)

// 	res := avroManager.GetSchema("ace", "v1")

// 	fmt.Println(res.String())
// 	fmt.Println(res.Type())
// }
