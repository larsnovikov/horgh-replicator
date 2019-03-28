package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/helpers"
	"strconv"
)

type sqlConnection struct {
	base *sql.DB
}

func (conn sqlConnection) Ping() bool {
	if conn.base.Ping() == nil {
		return true
	}

	return false
}

func (conn sqlConnection) Exec(params map[string]interface{}) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params["query"]), helpers.MakeSlice(params["params"])...)
	if err != nil {
		// TODO Надо проверять почему произошла ошибка.
		// Если duplicate on insert - игнорить
		// Поменялась структура - паниковать
		return false
	}

	return true
}

func (conn sqlConnection) Get(params map[string]interface{}) *sql.Rows {
	rows, err := conn.base.Query(fmt.Sprintf("%v", params["query"]), helpers.MakeSlice(params["params"])...)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func GetMysqlConnection(connection Connection, dbName string) interface{} {
	if connection == nil || connection.Ping() == false {
		cred := helpers.GetCredentials(dbName).(helpers.CredentialsDB)
		conn, err := sql.Open("mysql", buildMysqlString(cred))
		if err != nil || conn.Ping() != nil {
			connection = Retry(dbName, cred.Credentials, connection, GetMysqlConnection).(Connection)
		} else {
			connection = sqlConnection{conn}
		}
	}

	return connection
}

func buildMysqlString(cred helpers.CredentialsDB) string {
	return cred.User + ":" + cred.Pass + "@tcp(" + cred.Host + ":" + strconv.Itoa(cred.Port) + ")/" + cred.DBname + "?charset=utf8&parseTime=True&loc=Local"
}
