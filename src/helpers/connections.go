package helpers

import (
	"database/sql"
	"go-binlog-replication/src/constants"
)

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
		connectionPool.slave = GetMysqlConnection(connectionPool.slave, constants.DBSlave).(ConnectionSlave)
		return connectionPool.slave.Exec(params)
	case constants.DBReplicator:
		connectionPool.replicator = GetMysqlConnection(connectionPool.replicator, constants.DBReplicator).(ConnectionReplicator)
		return connectionPool.replicator.Exec(params)
	}

	return false
}

func Get(params map[string]interface{}) *sql.Rows {
	connectionPool.replicator = GetMysqlConnection(connectionPool.replicator, constants.DBReplicator).(ConnectionReplicator)
	return connectionPool.replicator.Get(params)
}
