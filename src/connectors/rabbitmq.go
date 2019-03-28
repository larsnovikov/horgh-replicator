package connectors

import (
	"github.com/streadway/amqp"
	"go-binlog-replication/src/helpers"
)

type rabbitmqConnection struct {
	base interface{}
}

func (conn rabbitmqConnection) Ping() bool {
	// TODO проверить подключение
	return true

	return false
}

func (conn rabbitmqConnection) Exec(params map[string]interface{}) bool {
	// TODO разобрать параметры и послать в RabbitMq
	return true
}

func GetRabbitmqConnection(connection Storage, storageType string) interface{} {
	if connection == nil || connection.Ping() == false {
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsAMQP)
		conn, err := amqp.Dial(buildRabbitmqString(cred))
		if err != nil {
			connection = Retry(storageType, cred.Credentials, connection, GetRabbitmqConnection).(Storage)
		} else {
			connection = rabbitmqConnection{conn}
		}
	}

	return connection
}

func buildRabbitmqString(cred helpers.CredentialsAMQP) string {
	return "amqp://" + cred.User + ":" + cred.Pass + "@" + cred.Host + ":" + cred.Port + "/"
}
