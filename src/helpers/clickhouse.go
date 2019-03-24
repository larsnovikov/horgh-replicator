package helpers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/kshvakov/clickhouse"
	"log"
	"strconv"
)

type clickhouseConnection struct {
	base *sqlx.DB
}

func (conn clickhouseConnection) Ping() bool {
	if conn.base.Ping() != nil {
		return true
	}

	return false
}

func (conn clickhouseConnection) Exec(params map[string]interface{}) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params["query"]), makeSlice(params["params"])...)
	if err != nil {
		log.Fatal(err)
		// TODO Надо проверять почему произошла ошибка.
		// Если duplicate on insert - игнорить
		// Поменялась структура - паниковать
		return false
	}

	return true
}

func GetClickhouseConnection(connection Connection, dbName string) interface{} {
	if connection == nil || connection.Ping() == true {
		cred := GetCredentials(dbName)
		conn, err := sqlx.Open("clickhouse", buildClickhouseString(cred))
		if err != nil {
			log.Fatal(err)
		}
		connection = clickhouseConnection{conn}
	}

	return connection
}

func buildClickhouseString(cred Credentials) string {
	// TODO разобраться с логином и паролем
	return "tcp://" + cred.Host + ":" + strconv.Itoa(cred.Port) + "?debug=true"
}
