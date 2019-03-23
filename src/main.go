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
		DBName := helpers.GetCredentials("master").DBname
		tableHash := fmt.Sprintf("%s.%s", DBName, tableName)
		go parser.BinlogListener(tableHash)
	}

	time.Sleep(10 * time.Minute)
	fmt.Print("Thx for watching, goodbuy")
}
