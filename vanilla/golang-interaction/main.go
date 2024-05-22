package main

import (
	"kafkaing/logging"
	"time"
)

func test() {
	logging.LogInstance.Access()
}

func main() {
	go test()
	time.Sleep(time.Second * 4)
}
