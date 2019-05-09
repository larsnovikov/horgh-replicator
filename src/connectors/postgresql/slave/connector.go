package slave

import (
	"database/sql"
	"fmt"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/tools/exit"
	"strconv"
)

const DSN = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

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
		log.Warnf(constants.ErrorExecQuery, "postgresql", err)
		return false
	}

	return true
}

func GetConnection(connection helpers.Storage, storageType string) interface{} {
	if connection == nil || connection.Ping() == false {
		cred := helpers.GetCredentials(storageType).(helpers.CredentialsDB)
		conn, err := sql.Open("postgres", buildDSN(cred))
		if err != nil || conn.Ping() != nil {
			exit.Fatal(constants.ErrorDBConnect, storageType)
		} else {
			connection = connect{conn}
		}
	}

	return connection
}

func buildDSN(cred helpers.CredentialsDB) string {
	return fmt.Sprintf(DSN, cred.Host, strconv.Itoa(cred.Port), cred.User, cred.Pass, cred.DBname)
}
