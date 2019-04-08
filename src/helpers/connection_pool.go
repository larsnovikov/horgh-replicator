package helpers

import (
	"database/sql"
)

type ConnectionReplicator interface {
	Storage
	Get(map[string]interface{}) *sql.Rows
}

type ConnectionPool struct {
	Master     Storage // used only for loader
	Slave      Storage
	Replicator ConnectionReplicator
}

var ConnPool ConnectionPool
