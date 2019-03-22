package helpers

import (
	_ "github.com/go-sql-driver/mysql"
)

type Connection interface {
	Ping() bool
	Exec(params map[string]interface{}) bool
}

var connection Connection

func Query(params map[string]interface{}) bool {
	slave := GetCredentials("slave").Type

	switch slave {
	case "mysql":
		return GetMysqlSlaveConnection().Exec(params)
	}

	return false
}
