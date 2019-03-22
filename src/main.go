package main

import (
	"fmt"
	"go-binlog-replication/src/parser"
	"time"
)

func main() {
	go parser.BinlogListener()

	time.Sleep(10 * time.Minute)
	fmt.Print("Thx for watching, goodbuy")
}
