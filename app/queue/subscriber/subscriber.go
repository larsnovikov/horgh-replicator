package subscriber

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"horgh2-replicator/app/configs"
)

type Subscriber struct {
	reader *kafka.Reader
}

func New(config configs.QueueConfig) Subscriber {
	return Subscriber{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{config.Host},
			Topic:     config.Topic,
			Partition: 0,
			MinBytes:  10e3, // 10KB
			MaxBytes:  10e6, // 10MB
		}),
	}
}

func (s Subscriber) Handle() {
	ctx := context.Background()
	for {
		m, err := s.reader.FetchMessage(ctx)
		if err != nil {
			break
		}

		err = s.reader.CommitMessages(ctx, m)
		fmt.Println(err) // TODO handle error
	}
}
