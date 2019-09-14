package connectors

import (
	"database/sql"
	"horgh2-replicator/app/configs"

	"horgh2-replicator/app/connectors/mysql/master"
	"horgh2-replicator/app/connectors/mysql/slave"

	slave2 "horgh2-replicator/app/connectors/pgsql/slave"
	"horgh2-replicator/app/constants"
)

type DBConnection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Connection interface {
	HealthCheck() bool
	Connect(configDSN string) error
}

type MasterConnection interface {
	Connection
	Listen()
}

type SlaveConnection interface {
	Connection
}

func NewSlave(config configs.ConnectionConfig) SlaveConnection {
	switch config.Type {
	case constants.TypeMYSQL:
		return slave.New(config)
	case constants.TypePGSQL:
		return slave2.New(config)
	}

	return nil
}

func NewMaster(config configs.ConnectionConfig) MasterConnection {
	switch config.Type {
	case constants.TypeMYSQL:
		return master.New(config)
	}

	return nil
}
