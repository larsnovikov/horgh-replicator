package clickhouse

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/kshvakov/clickhouse"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"strconv"
)

const DSN = "tcp://%s:%s?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20"

type connect struct {
	base *sqlx.DB
}

func (conn connect) Ping() bool {
	if conn.base.Ping() == nil {
		return true
	}

	return false
}

func (conn connect) Exec(params map[string]interface{}) bool {
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
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsDB)
		conn, err := sqlx.Open("clickhouse", buildDSN(cred))
		if err != nil || conn.Ping() != nil {
			connection = helpers.Retry(storageType, cred.Credentials, connection, GetConnection).(helpers.Storage)
		} else {
			connection = connect{conn}
		}
	}

	return connection
}

func buildDSN(cred helpers.CredentialsDB) string {
	return fmt.Sprintf(DSN, cred.Host, strconv.Itoa(cred.Port), cred.User, cred.Pass, cred.DBname)
}
