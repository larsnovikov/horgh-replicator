package master

import (
	"database/sql"
	"horgh-replicator/src/connectors/mysql"
	mysqlMaster "horgh-replicator/src/connectors/mysql/master"
	"horgh-replicator/src/connectors/postgresql"
	postgresqlMaster "horgh-replicator/src/connectors/postgresql/master"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"os"
)

type AbstractMaster interface {
	Listen()
	BuildSlave(table string)
}

var master AbstractMaster

func GetModel() AbstractMaster {
	return master
}

func MakeMaster() {
	switch os.Getenv("MASTER_TYPE") {
	case "postgresql":
		master = &postgresqlMaster.Model{}
		helpers.ConnPool.Master = postgresql.GetConnection(helpers.ConnPool.Master, constants.DBMaster).(helpers.ConnectionMaster)
	case "mysql":
		master = &mysqlMaster.Model{}
		helpers.ConnPool.Master = mysql.GetConnection(helpers.ConnPool.Master, constants.DBMaster).(helpers.ConnectionMaster)
	}
}

func Exec(params helpers.Query) bool {
	return helpers.ConnPool.Master.Exec(params)
}

func Get(params helpers.Query) *sql.Rows {
	return helpers.ConnPool.Master.Get(params)
}
