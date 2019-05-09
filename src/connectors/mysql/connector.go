package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/tools/exit"
	"strconv"
)

const DSN = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

type connect struct {
	base *sql.DB
}

func (conn connect) Ping() bool {
	if conn.base.Ping() == nil {
		return true
	}

	return false
}

func (conn connect) Exec(params helpers.Query) bool {
	_, err := conn.base.Exec(fmt.Sprintf("%v", params.Query), helpers.MakeSlice(params.Params)...)
	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "mysql", err)
		return false
	}

	return true
}

func (conn connect) Get(params helpers.Query) *sql.Rows {
	rows, err := conn.base.Query(fmt.Sprintf("%v", params.Query), helpers.MakeSlice(params.Params)...)
	if err != nil {
		exit.Fatal(err.Error())
	}

	return rows
}

func GetConnection(connection helpers.Storage, storageType string) interface{} {
	if connection == nil || connection.Ping() == false {
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsDB)
		conn, err := sql.Open("mysql", buildDSN(cred))
		if err != nil || conn.Ping() != nil {
			exit.Fatal(constants.ErrorDBConnect, storageType)
		} else {
			connection = connect{conn}
		}
	}

	return connection
}

func buildDSN(cred helpers.CredentialsDB) string {
	return fmt.Sprintf(DSN, cred.User, cred.Pass, cred.Host, strconv.Itoa(cred.Port), cred.DBname)
}
