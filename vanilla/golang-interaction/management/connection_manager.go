package management

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ConnectionManager struct {
	topic        string
	partition    int
	ctx          context.Context
	conn         *kafka.Conn
	producerChan chan []byte
	consumerChan chan kafka.Message
}

func GenerateNewCm(topic string, partition int) *ConnectionManager {
	return &ConnectionManager{topic: topic, partition: partition, ctx: context.Background()}
}

func (cm *ConnectionManager) EstablishConnection(network, address string, consumerFn func([]byte)) chan<- []byte {
	var err error
	cm.conn, err = kafka.DialLeader(cm.ctx, network, address, cm.topic, cm.partition)
	if err != nil {
		zap.L().Sugar().Error("failed to establish connection %s", err)
		return nil
	}

	cm.producerChan = make(chan []byte)
	cm.consumerChan = make(chan kafka.Message)

	go cm.runBackgroundProducer()

	if consumerFn != nil {
		go cm.runBackgroundConsumer(consumerFn)
	}

	return cm.producerChan
}

func (cm *ConnectionManager) runBackgroundProducer() {
	for {
		byteArr := <-cm.producerChan
		msg := kafka.Message{Topic: cm.topic, Partition: cm.partition, Value: byteArr}
		_, err := cm.conn.WriteMessages(msg)
		if err != nil {
			zap.L().Sugar().Error("failed to send msg %s", err)
		} else {
			zap.L().Sugar().Info("sucessfully sent msg to %s %s", msg.Topic, string(rune(msg.Partition)))
		}
	}

}

func (cm *ConnectionManager) runBackgroundConsumer(consumerFn func([]byte)) {
	b := make([]byte, 10e3)
	for {
		_, err := cm.conn.Read(b)
		if err != nil {
			zap.L().Sugar().Error("failed to consume msg %s", err)
			continue
		}
		consumerFn(b)
	}

}

func (cm *ConnectionManager) ShutdownHook() {
	cm.conn.Close()
	zap.L().Sugar().Info("Kafka Connection Manager shutdowned")
}
