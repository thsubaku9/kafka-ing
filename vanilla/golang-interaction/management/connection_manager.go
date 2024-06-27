package management

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ConnectionManager struct {
	topic        string
	partition    int
	groupId      string
	ctx          context.Context
	conn         *kafka.Conn
	producerChan chan []byte
}

func GenerateNewCm(topic string, partition int, groupId string) *ConnectionManager {
	return &ConnectionManager{topic: topic, partition: partition, groupId: groupId, ctx: context.Background()}
}

func (cm *ConnectionManager) EstablishConnection(network, address string, consumerFn func(kafka.Message)) chan<- []byte {
	var err error
	cm.conn, err = kafka.DialLeader(cm.ctx, network, address, cm.topic, cm.partition)
	if err != nil {
		zap.L().Sugar().Errorf("failed to establish connection %s", err)
		return nil
	}

	cm.producerChan = make(chan []byte)
	go cm.runBackgroundProducer(address)

	if consumerFn != nil {
		go cm.runBackgroundConsumer(address, cm.groupId, consumerFn)
	}

	return cm.producerChan
}

func (cm *ConnectionManager) runBackgroundProducer(address string) {

	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Balancer: &kafka.Hash{},
	}

	for {
		byteArr := <-cm.producerChan
		msg := kafka.Message{Topic: cm.topic, Partition: cm.partition, Value: byteArr}

		err := writer.WriteMessages(cm.ctx, msg)
		if err != nil {
			zap.L().Sugar().Errorf("failed to send msg %s", err)
		} else {
			zap.L().Sugar().Infof("sucessfully sent msg to %s %s", msg.Topic, string(rune(msg.Partition)))
		}
	}

}

func (cm *ConnectionManager) runBackgroundConsumer(address, groupId string, consumerFn func(kafka.Message)) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     cm.topic,
		GroupID:   groupId,
		Partition: cm.partition,
		MaxBytes:  10e5, // 1MB
	})

	for {
		m, err := reader.ReadMessage(cm.ctx)
		if err != nil {
			zap.L().Sugar().Errorf("failed to consume msg %s", err)
			continue
		}

		zap.L().Sugar().Debugf("Message obtained %s %d %d", m.Topic, m.Partition, m.Offset)
		consumerFn(m)
	}
}

func (cm *ConnectionManager) ShutdownHook() {
	cm.conn.Close()
	zap.L().Sugar().Info("Kafka Connection Manager shutdowned")
}
