package system

import (
	"database/sql"
	"go-binlog-replication/src/connectors/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
)

//type ConnectionPool struct {
//	master     Storage // used only for loader
//	replicator ConnectionReplicator
//}

func Exec(mode string, params map[string]interface{}) bool {
	switch mode {
	case constants.DBMaster:
		helpers.ConnPool.Master = mysql.GetConnection(helpers.ConnPool.Master, constants.DBMaster).(helpers.Storage)
		return helpers.ConnPool.Master.Exec(params)
	case constants.DBReplicator:
		helpers.ConnPool.Replicator = mysql.GetConnection(helpers.ConnPool.Replicator, constants.DBReplicator).(helpers.ConnectionReplicator)
		return helpers.ConnPool.Replicator.Exec(params)
	}

	return false
}

func Get(params map[string]interface{}) *sql.Rows {
	helpers.ConnPool.Replicator = mysql.GetConnection(helpers.ConnPool.Replicator, constants.DBReplicator).(helpers.ConnectionReplicator)
	return helpers.ConnPool.Replicator.Get(params)
}
