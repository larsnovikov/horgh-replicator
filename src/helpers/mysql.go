package helpers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-log/log"
	"strconv"
)

type sqlConnection struct {
	base *sql.DB
}

func (conn sqlConnection) Ping() bool {
	if conn.base.Ping() != nil {
		return true
	}

	return false
}

func (conn sqlConnection) Exec(params map[string]interface{}) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params["query"]), makeSlice(params["params"])...)
	if err != nil {
		// TODO Надо проверять почему произошла ошибка.
		// Если duplicate on insert - игнорить
		// Поменялась структура - паниковать
		return false
	}

	return true
}

func (conn sqlConnection) Get(params map[string]interface{}) *sql.Rows {
	rows, err := conn.base.Query(fmt.Sprintf("%v", params["query"]), makeSlice(params["params"])...)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func GetMysqlConnection(connection ConnectionSlave, dbName string) interface{} {
	if connection == nil || connection.Ping() == true {
		cred := GetCredentials(dbName)
		conn, err := sql.Open("mysql", buildMysqlString(cred))
		if err != nil {
			log.Fatal(err)
		}
		connection = sqlConnection{conn}
	}

	return connection
}

func buildMysqlString(cred Credentials) string {
	return cred.User + ":" + cred.Pass + "@tcp(" + cred.Host + ":" + strconv.Itoa(cred.Port) + ")/" + cred.DBname + "?charset=utf8&parseTime=True&loc=Local"
}
