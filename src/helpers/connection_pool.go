package helpers

import (
	"database/sql"
)

type ConnectionMaster interface {
	Storage
	Get(map[string]interface{}) *sql.Rows
}

type ConnectionPool struct {
	Master ConnectionMaster
	Slave  Storage
}

var ConnPool ConnectionPool
