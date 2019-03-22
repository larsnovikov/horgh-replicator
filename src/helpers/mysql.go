package helpers

import (
	"database/sql"
	"fmt"
	"log"
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

func GetMysqlConnection(connection ConnectionSlave) interface{ ConnectionSlave } {
	if connection == nil || connection.Ping() == true {
		slave := GetCredentials("slave")
		conn, err := sql.Open("mysql", buildMysqlString(slave))
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
