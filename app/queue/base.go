package queue

import (
	"context"
	"github.com/segmentio/kafka-go"
	"horgh2-replicator/app/configs"
)

type Connection struct {
	connect *kafka.Conn
}

func New(config configs.QueueConfig) (Connection, error) {
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.Host, config.Topic, partition)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		connect: conn,
	}, nil
}
