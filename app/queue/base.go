package queue

import (
	"context"
	"github.com/segmentio/kafka-go"
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/queue/publisher"
	"horgh2-replicator/app/queue/subscriber"
)

type Connection struct {
	Connect    *kafka.Conn
	Publisher  publisher.Publisher
	Subscriber subscriber.Subscriber
}

func New(config configs.QueueConfig) (Connection, error) {
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.Host, config.Topic, partition)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		Connect:    conn,
		Publisher:  publisher.New(conn),
		Subscriber: subscriber.New(conn),
	}, nil
}
