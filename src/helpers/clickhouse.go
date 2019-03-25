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
	tx, _ := conn.base.Begin()
	_, err := tx.Exec(fmt.Sprintf("%v", params["query"]), makeSlice(params["params"])...)
	tx.Commit()
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
			Retry(cred, connection, GetClickhouseConnection)
		}
		connection = clickhouseConnection{conn}
	}

	return connection
}

func buildClickhouseString(cred Credentials) string {
	return "tcp://" + cred.Host + ":" + strconv.Itoa(cred.Port) + "?username=" + cred.User + "&password=" + cred.Pass + "&database=" + cred.DBname + "&read_timeout=10&write_timeout=20"
}
