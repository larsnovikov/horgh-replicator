package master

import (
	mysqlMaster "horgh-replicator/src/connectors/mysql/master"
	postgresqlMaster "horgh-replicator/src/connectors/postgresql/master"
	"os"
)

type AbstractMaster interface {
	Listen()
	BuildSlave(table string)
}

var master AbstractMaster
var built = false

func GetModel() AbstractMaster {
	if built == false {
		switch os.Getenv("MASTER_TYPE") {
		case "postgresql":
			master = &postgresqlMaster.Model{}
		default:
			master = &mysqlMaster.Model{}
		}
		built = true
	}

	return master
}
