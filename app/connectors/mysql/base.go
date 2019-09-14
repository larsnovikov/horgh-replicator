package mysql

import (
	"database/sql"
)

type Connection struct {
	connect sql.DB
}

func (c Connection) HealthCheck() bool {
	return true
}

func (c Connection) Connect(configDSN string) error {
	conn, err := sql.Open("mysql", configDSN)
	if err != nil {
		return err
	}

	c.connect = *conn
	return nil
}
