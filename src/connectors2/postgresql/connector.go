package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strconv"
)

type postgresqlConnection struct {
	base *sql.DB
}

func (conn postgresqlConnection) Ping() bool {
	if conn.base.Ping() == nil {
		return true
	}

	return false
}

func (conn postgresqlConnection) Exec(params map[string]interface{}) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params["query"]), helpers.MakeSlice(params["params"])...)
	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "postgresql", err)
		return false
	}

	return true
}

func GetConnection(connection helpers.Storage, storageType string) interface{} {
	if connection == nil || connection.Ping() == false {
		helpers.ParseDBConfig()
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsDB)
		conn, err := sql.Open("postgres", buildDSN(cred))
		if err != nil || conn.Ping() != nil {
			connection = helpers.Retry(storageType, cred.Credentials, connection, GetConnection).(helpers.Storage)
		} else {
			connection = postgresqlConnection{conn}
		}
	}

	return connection
}

func buildDSN(cred helpers.CredentialsDB) string {
	// TODO constant
	return "host=" + cred.Host + " port=" + strconv.Itoa(cred.Port) + " user=" + cred.User + " password=" + cred.Pass + " dbname=" + cred.DBname + " sslmode=disable"
}
