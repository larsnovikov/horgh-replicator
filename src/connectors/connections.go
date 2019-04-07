package connectors

import (
	"database/sql"
	"go-binlog-replication/src/connectors2/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
)

type ConnectionPool struct {
	master     helpers.Storage // used only for loader
	slave      helpers.Storage
	replicator helpers.ConnectionReplicator
}

var connectionPool ConnectionPool

func Exec(mode string, params map[string]interface{}) bool {
	switch mode {
	case constants.DBMaster:
		connectionPool.master = mysql.GetConnection(connectionPool.master, constants.DBMaster).(helpers.Storage)
		return connectionPool.master.Exec(params)
	case constants.DBReplicator:
		connectionPool.replicator = mysql.GetConnection(connectionPool.replicator, constants.DBReplicator).(helpers.ConnectionReplicator)
		return connectionPool.replicator.Exec(params)
	// adapters for slave storages
	case "mysql":
		connectionPool.slave = mysql.GetConnection(connectionPool.slave, constants.DBSlave).(helpers.Storage)
		return connectionPool.slave.Exec(params)
	case "clickhouse":
		connectionPool.slave = GetClickhouseConnection(connectionPool.slave, constants.DBSlave).(helpers.Storage)
		return connectionPool.slave.Exec(params)
	case "postgresql":
		connectionPool.slave = GetPostgresqlConnection(connectionPool.slave, constants.DBSlave).(helpers.Storage)
		return connectionPool.slave.Exec(params)
		/*case "rabbitmq":
		connectionPool.slave = GetRabbitmqConnection(connectionPool.slave, constants.DBSlave).(Storage)*/
	}

	return false
}

func Get(params map[string]interface{}) *sql.Rows {
	connectionPool.replicator = mysql.GetConnection(connectionPool.replicator, constants.DBReplicator).(helpers.ConnectionReplicator)
	return connectionPool.replicator.Get(params)
}
