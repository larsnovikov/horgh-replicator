package main

import (
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/slave"
	"go-binlog-replication/src/parser"
	"time"
)

func main() {
	helpers.MakeCredentials()
	slave.MakeSlavePool()
	parser.BinlogListener()

	time.Sleep(60 * time.Minute)
}
