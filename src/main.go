package main

import (
	"fmt"
	"go-binlog-example/src/parser"
	"time"
)

func main() {
	go parser.BinlogListener()

	time.Sleep(2 * time.Minute)
	fmt.Print("Thx for watching, goodbuy")
}
