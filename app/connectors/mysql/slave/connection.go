package slave

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/connectors/mysql"
)

type Connection struct {
	mysql.Connection
}

func New(config configs.ConnectionConfig) Connection {
	return Connection{}
}
