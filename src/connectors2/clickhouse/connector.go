package clickhouse

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/kshvakov/clickhouse"
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strconv"
)

type clickhouseConnection struct {
	base *sqlx.DB
}

func (conn clickhouseConnection) Ping() bool {
	if conn.base.Ping() == nil {
		return true
	}

	return false
}

func (conn clickhouseConnection) Exec(params map[string]interface{}) bool {
	tx, _ := conn.base.Begin()
	_, err := tx.Exec(fmt.Sprintf("%v", params["query"]), helpers.MakeSlice(params["params"])...)

	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "clickhouse", err)
		return false
	}

	defer func() {
		err = tx.Commit()
	}()

	return true
}

func GetConnection(connection helpers.Storage, storageType string) interface{} {
	if connection == nil || connection.Ping() == false {
		helpers.ParseDBConfig()
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsDB)
		conn, err := sqlx.Open("clickhouse", buildDSN(cred))
		if err != nil || conn.Ping() != nil {
			connection = helpers.Retry(storageType, cred.Credentials, connection, GetConnection).(helpers.Storage)
		} else {
			connection = clickhouseConnection{conn}
		}
	}

	return connection
}

func buildDSN(cred helpers.CredentialsDB) string {
	// TODO constant and sprintf
	return "tcp://" + cred.Host + ":" + strconv.Itoa(cred.Port) + "?username=" + cred.User + "&password=" + cred.Pass + "&database=" + cred.DBname + "&read_timeout=10&write_timeout=20"
}
