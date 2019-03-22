package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
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
		log.Fatal(err)
		return false
	}

	return true
}

func GetMysqlSlaveConnection() interface{ Connection } {
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

func CloseMysql(rows interface{}) {
	testRow := new(sql.Rows)
	if reflect.TypeOf(rows) == reflect.TypeOf(*testRow) {
		pointer := rows.(sql.Rows)
		defer pointer.Close()
	}
}
