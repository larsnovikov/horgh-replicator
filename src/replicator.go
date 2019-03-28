package main

import (
	"fmt"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/parser"
	"time"
)

func main() {
	helpers.MakeCredentials()
	for _, tableName := range helpers.GetTables() {
		DBName := helpers.GetCredentials("master").(helpers.CredentialsDB).DBname
		tableHash := fmt.Sprintf("%s.%s", DBName, tableName)
		go parser.BinlogListener(tableHash)
	}

	time.Sleep(60 * time.Minute)
}
