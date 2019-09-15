package publisher

import (
	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	conn *kafka.Conn
}

func New(connection *kafka.Conn) Publisher {
	return Publisher{
		conn: connection,
	}
}

func (p Publisher) Handle(msg []byte) (int, error) {
	return p.conn.WriteMessages(kafka.Message{Value: msg})
}
