package helpers

import "database/sql"

type ConnectionSlave interface {
	Ping() bool
	Exec(params map[string]interface{}) bool
}

type ConnectionReplicator interface {
	ConnectionSlave
	Get(map[string]interface{}) *sql.Rows
}

type ConnectionPool struct {
	slave      ConnectionSlave
	replicator ConnectionReplicator
}

var connectionPool ConnectionPool

func Exec(mode string, params map[string]interface{}) bool {
	switch mode {
	case "mysql":
		return GetMysqlConnection(connectionPool.slave).Exec(params)
	case "replicator":
		return GetMysqlConnection(connectionPool.replicator).Exec(params)
	}

	return false
}

func Get(params map[string]interface{}) *sql.Rows {
	return connectionPool.replicator.Get(params)
}
