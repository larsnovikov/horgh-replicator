package helpers

import (
	"database/sql"
	"fmt"
	_ "github.com/jmoiron/sqlx"
	"github.com/siddontang/go-log/log"
	"strconv"
)

type postgresqlConnection struct {
	base *sql.DB
}

func (conn postgresqlConnection) Ping() bool {
	if conn.base.Ping() != nil {
		return true
	}

	return false
}

func (conn postgresqlConnection) Exec(params map[string]interface{}) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params["query"]), makeSlice(params["params"])...)
	if err != nil {
		// TODO Надо проверять почему произошла ошибка.
		// Если duplicate on insert - игнорить
		// Поменялась структура - паниковать
		return false
	}

	return true
}

func (conn postgresqlConnection) Get(params map[string]interface{}) *sql.Rows {
	rows, err := conn.base.Query(fmt.Sprintf("%v", params["query"]), makeSlice(params["params"])...)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func GetPostgresqlConnection(connection Connection, dbName string) interface{} {
	if connection == nil || connection.Ping() == true {
		cred := GetCredentials(dbName)
		conn, err := sql.Open("postgres", buildPostgresqlString(cred))
		if err != nil || conn.Ping() != nil {
			connection = Retry(dbName, cred, connection, GetClickhouseConnection).(Connection)
		} else {
			connection = postgresqlConnection{conn}
		}
	}

	return connection
}

func buildPostgresqlString(cred Credentials) string {
	return "host=" + cred.Host + " port=" + strconv.Itoa(cred.Port) + " user=" + cred.User + " password=" + cred.Pass + " dbname=" + cred.DBname + " sslmode=disable"
}
