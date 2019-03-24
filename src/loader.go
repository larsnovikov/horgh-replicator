package main

import (
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/helpers"
	"math/rand"
	"strconv"
	"time"
)

const (
	ThreadCount = 5
	LoadTime    = 5
)

func main() {
	log.Info("Start loader")
	helpers.MakeCredentials()

	for i := 0; i < ThreadCount; i++ {
		log.Infof("Create goroutine #%s", strconv.Itoa(i+1))
		go load()
	}

	time.Sleep(LoadTime * time.Minute)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func load() {
	queries := []string{
		"INSERT INTO test.user (`name`, `status`) VALUE ('Jack', 'active');",
		"UPDATE test.user SET `name`='Tommy' ORDER BY RAND() LIMIT 1",
		"DELETE FROM test.user LIMIT 1;",
		"INSERT INTO test.post (`title`, `text`) VALUE ('Title', 'London is the capital of Great Britain');",
		"UPDATE test.post SET title='New title' ORDER BY RAND() LIMIT 1;",
		"DELETE FROM test.post LIMIT 1;",
	}

	rand.Seed(time.Now().UTC().UnixNano())

	var query string
	var result bool

	for {
		query = queries[randInt(0, len(queries))]

		result = helpers.Exec("master", map[string]interface{}{
			"query":  query,
			"params": []interface{}{},
		})

		log.Infof("Result: %s", strconv.FormatBool(result))
	}
}
