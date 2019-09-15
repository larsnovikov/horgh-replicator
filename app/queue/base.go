package queue

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/queue/publisher"
	"horgh2-replicator/app/queue/subscriber"
)

type Connection struct {
	Publisher  publisher.Publisher
	Subscriber subscriber.Subscriber
}

func New(config configs.QueueConfig) (Connection, error) {
	return Connection{
		Publisher:  publisher.New(config),
		Subscriber: subscriber.New(config),
	}, nil
}
