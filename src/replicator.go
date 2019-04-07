package main

import (
	"fmt"
	"go-binlog-replication/src/helpers"
	slave2 "go-binlog-replication/src/models/slave"
	"go-binlog-replication/src/parser"
	"time"
)

func main() {
	helpers.MakeCredentials()
	for _, tableName := range helpers.GetTables() {
		DBName := helpers.GetCredentials("master").(helpers.CredentialsDB).DBname
		tableHash := fmt.Sprintf("%s.%s", DBName, tableName)
		slave := slave2.MakeSlave(tableName)
		go func() {
			parser.BinlogListener(tableHash, slave)
		}()
		time.Sleep(5 * time.Second)
	}

	time.Sleep(60 * time.Minute)
}
