package rabbitmq

import (
	"github.com/siddontang/go-log/log"
	"github.com/streadway/amqp"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strconv"
)

type rabbitmqConnection struct {
	base interface{}
}

func (conn rabbitmqConnection) Ping() bool {
	// TODO проверить подключение
	return true
}

func (conn rabbitmqConnection) Exec(params map[string]interface{}) bool {
	ch := params["channel"].(amqp.Channel)

	err := ch.Publish(
		params["exchange"].(string),
		params["routingKey"].(string),
		params["mandatory"].(bool),
		params["immediate"].(bool),
		amqp.Publishing{
			ContentType: params["contentType"].(string),
			Body:        []byte(params["body"].(string)),
		})

	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "rabbitmq", err)
		return false
	}

	return true
}

func GetConnection(connection helpers.Storage, storageType string) interface{} {
	if connection == nil || connection.Ping() == false {
		helpers.ParseAMQPConfig()
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsAMQP)
		conn, err := amqp.Dial(buildRabbitmqString(cred))
		if err != nil {
			connection = helpers.Retry(storageType, cred.Credentials, connection, GetConnection).(helpers.Storage)
		} else {
			connection = rabbitmqConnection{conn}
		}
	}

	return connection
}

func buildDsn(cred helpers.CredentialsAMQP) string {
	// TODO constant
	return "amqp://" + cred.User + ":" + cred.Pass + "@" + cred.Host + ":" + strconv.Itoa(cred.Port) + "/"
}
