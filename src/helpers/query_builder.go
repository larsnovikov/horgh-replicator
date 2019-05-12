package helpers

import (
	"fmt"
	"io/ioutil"
)

const (
	pgPath          = "../files/sql/pg_"
	pgCreateTable   = "create_log_table.sql"
	pgCreateFunc    = "create_log_func.sql"
	pgCreateTrigger = "create_log_trigger.sql"
)

func getFilePath(dbType string, queryType string) string {
	var fileName string
	switch dbType {
	default: // postgree
		fileName = pgPath + "create_log_" + queryType + ".sql"
	}

	return fileName
}

func getQuery(fileName string, params ...interface{}) string {
	// TODO check file exists
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		// TODO error
	}
	data := string(content[:])

	return fmt.Sprintf(data, params)
}

func GetQuery(dbType string, queryType string, params ...interface{}) string {
	fileName := getFilePath(dbType, queryType)

	return getQuery(fileName, params)
}
