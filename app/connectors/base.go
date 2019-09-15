package connectors

import (
	"database/sql"
	"errors"
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

func NewSlave(config configs.ConnectionConfig) (SlaveConnection, error) {
	switch config.Type {
	case constants.TypeMYSQL:
		return slave.New(config), nil
	case constants.TypePGSQL:
		return slave2.New(config), nil
	}

	return nil, errors.New(constants.ErrorSlaveType)
}

func NewMaster(config configs.ConnectionConfig) (MasterConnection, error) {
	switch config.Type {
	case constants.TypeMYSQL:
		return master.New(config), nil
	}

	return nil, errors.New(constants.ErrorMasterType)
}
