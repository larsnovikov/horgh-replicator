package slave

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/connectors/pgsql"
)

type Connection struct {
	pgsql.Connection
}

func New(config configs.ConnectionConfig) Connection {
	return Connection{}
}
