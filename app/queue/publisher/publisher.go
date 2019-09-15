package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
	"horgh2-replicator/app/configs"
)

type Publisher struct {
	writer *kafka.Writer
}

func New(config configs.QueueConfig) Publisher {
	return Publisher{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{config.Host},
			Topic:    config.Topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (p Publisher) Handle(msg []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{Value: msg})
}
