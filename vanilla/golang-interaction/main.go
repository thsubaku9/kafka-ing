package main

import (
	"kafkaing/logging"
	"kafkaing/management"
	"time"

	"github.com/segmentio/kafka-go"

	"go.uber.org/zap"
)

func init() {
	logging.Initialize()
}

func printBytes(m kafka.Message) {
	zap.L().Sugar().Debugf("Read : %s\n", string(m.Value))
}

func main() {
	cm := management.GenerateNewCm("mcswirl", 0, "gid1")
	producerChannel := cm.EstablishConnection("tcp", "localhost:9092", printBytes)

	producerChannel <- []byte("mario")
	producerChannel <- []byte("luigi")

	time.Sleep(5 * time.Second)

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

/*
type SimpleRecord struct {
	A int64  `avro:"a"`
	B string `avro:"b"`
}
*/
