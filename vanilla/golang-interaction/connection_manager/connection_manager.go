package connectionmanager

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ConnectionManager struct {
	topic     string
	partition int
	ctx       context.Context
	conn      *kafka.Conn
}

func GenerateNewCm(topic string, partition int) *ConnectionManager {
	return &ConnectionManager{topic: topic, partition: partition, ctx: context.Background()}
}

func (cm *ConnectionManager) establishConnection(network, address string) {
	var err error
	cm.conn, err = kafka.DialLeader(cm.ctx, network, address, cm.topic, cm.partition)
	if err != nil {

	}
}
