package vertica

import (
	"fmt"
	"go-binlog-replication/src/_vendor-20190326172913/github.com/jmoiron/sqlx"
	"go-binlog-replication/src/_vendor-20190326172913/github.com/siddontang/go-log/log"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strconv"
)

type verticaConnection struct {
	base *sqlx.DB
}

func (conn verticaConnection) Ping() bool {
	if conn.base.Ping() == nil {
		return true
	}

	return false
}

func (conn verticaConnection) Exec(params map[string]interface{}) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params["query"]), helpers.MakeSlice(params["params"])...)

	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "vertica", err)
		return false
	}

	return true
}

func GetConnection(connection helpers.Storage, dbName string) interface{} {
	if connection == nil || connection.Ping() == false {
		cred := helpers.GetCredentials(dbName).(helpers.CredentialsDB)
		conn, err := sqlx.Open("odbc", buildDSN(cred))
		if err != nil || conn.Ping() != nil {
			connection = helpers.Retry(dbName, cred.Credentials, connection, GetConnection).(helpers.Storage)
		} else {
			connection = verticaConnection{conn}
		}
	}

	return connection
}

func buildDSN(cred helpers.CredentialsDB) string {
	driver := "/opt/vertica/lib64/libverticaodbc.so"
	return "Driver=" + driver + ";ServerName=" + cred.Host + ";Database=" + cred.DBname + ";Port=" + strconv.Itoa(cred.Port) + ";uid=" + cred.User + ";pwd=" + cred.Pass + ";"
}
