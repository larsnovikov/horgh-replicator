package subscriber

import "github.com/segmentio/kafka-go"

type Subscriber struct {
	conn *kafka.Conn
}

func New(connection *kafka.Conn) Subscriber {
	return Subscriber{
		conn: connection,
	}
}

func (s Subscriber) Handle() {
}
