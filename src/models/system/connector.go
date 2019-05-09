package system

import (
	"database/sql"
	"horgh-replicator/src/connectors/mysql/slave"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
)

func Exec(mode string, params helpers.Query) bool {
	switch mode {
	case constants.DBMaster:
		helpers.ConnPool.Master = slave.GetConnection(helpers.ConnPool.Master, constants.DBMaster).(helpers.ConnectionMaster)
		return helpers.ConnPool.Master.Exec(params)
	}

	return false
}

func Get(mode string, params helpers.Query) *sql.Rows {
	switch mode {
	default:
		helpers.ConnPool.Master = slave.GetConnection(helpers.ConnPool.Master, constants.DBMaster).(helpers.ConnectionMaster)
		return helpers.ConnPool.Master.Get(params)
	}
}
