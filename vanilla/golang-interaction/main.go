package main

import (
	"kafkaing/logging"
	"kafkaing/management"
	"time"

	"github.com/hamba/avro"
	"github.com/segmentio/kafka-go"

	"go.uber.org/zap"
)

func init() {
	logging.Initialize()
}

func printBytes(m kafka.Message) {
	zap.L().Sugar().Debugf("Read : %s\n", string(m.Value))
}

type Sandwich struct {
	Bread     string         `avro:"bread"`
	Sauces    []string       `avro:"sauces"`
	Stuffings map[string]int `avro:"stuffings"`
	Style     string         `avro:"style"`
	Toasted   bool           `avro:"toasted"`
	IsOpen    bool           `avro:"isOpen"`
	Count     int            `avro:"count"`
}

func createNewSandwich(bread, style string, sauces []string, stuffings map[string]int, toasted, isOpen bool, count int) Sandwich {
	return Sandwich{
		Bread:     bread,
		Sauces:    sauces,
		Stuffings: stuffings,
		Style:     style,
		Toasted:   toasted,
		IsOpen:    isOpen,
		Count:     count,
	}
}

const SandwichAvroModel string = `{
	"type": "record",
	"name": "Sandwich",
	"namespace": "kafaking",
	"fields": [
	  {
		"name": "bread",
		"type": "string"
	  },
	  {
		"name": "sauces",
		"type": { "type": "array", "items": "string" }
	  },
	  {
		"name": "stuffings",
		"type": { "type": "map", "values": "int" }
	  },
	  {
		"name": "style",
		"type": "string"
	  },
	  {
		"name": "toasted",
		"type": "boolean"
	  },
	  {
		"name": "isOpen",
		"type": "boolean"
	  },
	  {
		"name": "count",
		"type": "int"
	  }
	]
  }`

func main() {
	var avroManager management.AvroManager = management.InitAvroManager("sqlite3", "./schemaRegistry.db")

	avroManager.UpsertSchema("sandwich", "v1", SandwichAvroModel)
	sandwichSchema := avroManager.GetSchema("sandwich", "v1")

	cm := management.GenerateNewCm("sandwichShop", 0, "chef1")
	producerChannel := cm.EstablishConnection("tcp", "localhost:9092", func(kmsg kafka.Message) {
		var res Sandwich = Sandwich{}
		avro.Unmarshal(sandwichSchema, kmsg.Value, &res)
		zap.L().Sugar().Infof("Received %v", res)
	})

	{
		s1 := createNewSandwich("Wheat", "soft", []string{"honey mustard", "ice spice"}, map[string]int{"cheese": 1, "tomato": 1, "cucumber": 1}, false, false, 1)
		res, err := avro.Marshal(sandwichSchema, s1)

		if err != nil {
			zap.L().Sugar().Error(err)
		} else {
			producerChannel <- kafka.Message{Key: []byte(s1.Bread), Value: res}
		}
	}

	{
		s1 := createNewSandwich("Rye", "garlic", []string{"honey mustard", "harissa", "bbq"}, map[string]int{"cheddar": 1, "capsicum": 1, "cucumber": 1}, true, false, 2)
		res, err := avro.Marshal(sandwichSchema, s1)

		if err != nil {
			zap.L().Sugar().Error(err)
		} else {
			producerChannel <- kafka.Message{Key: []byte(s1.Bread), Value: res}
		}

	}

	time.Sleep(10 * time.Second)

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
