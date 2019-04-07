package main

import (
	"fmt"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models"
	"go-binlog-replication/src/parser"
	"time"
)

func main() {
	helpers.MakeCredentials()
	for _, tableName := range helpers.GetTables() {
		DBName := helpers.GetCredentials("master").(helpers.CredentialsDB).DBname
		tableHash := fmt.Sprintf("%s.%s", DBName, tableName)
		slave := models.MakeSlave(tableName)
		go func() {
			parser.BinlogListener(tableHash, slave)
		}()
		time.Sleep(5 * time.Second)
	}

	time.Sleep(60 * time.Minute)
}
