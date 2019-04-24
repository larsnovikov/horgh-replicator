package helpers

import (
	"database/sql"
)

type ConnectionMaster interface {
	Storage
	Get(map[string]interface{}) *sql.Rows
}

type ConnectionReplicator interface {
	ConnectionMaster
}

type ConnectionPool struct {
	Master     ConnectionMaster
	Slave      Storage
	Replicator ConnectionReplicator
}

var ConnPool ConnectionPool
