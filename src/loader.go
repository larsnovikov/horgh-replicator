package main

import (
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/helpers"
	"math/rand"
	"strconv"
	"time"
)

const (
	// goroutine count. WARNING if you set more 1, may be concurrency problems
	ThreadCount = 1
	// time to create queries in minutes
	LoadTime = 1
)

var counters map[int]int

func main() {
	log.Info("Start loader")
	helpers.MakeCredentials()
	counters = make(map[int]int)

	for i := 0; i < ThreadCount; i++ {
		log.Infof("Create goroutine #%s", strconv.Itoa(i+1))
		counters[i] = 0
		go load(i)
	}

	time.Sleep(LoadTime * time.Minute)
	showStats()
	log.Info("Stop loader")
}

func showStats() {
	totalQueries := 0

	for i := 0; i < ThreadCount; i++ {
		log.Infof("Goroutine create %s queries per %s minute(s)", strconv.Itoa(counters[i]), strconv.Itoa(LoadTime))
		totalQueries = totalQueries + counters[i]
	}

	queriesPerMinute := totalQueries / LoadTime
	log.Infof("Total queries: %s; Queries per minute: %s", strconv.Itoa(totalQueries), strconv.Itoa(queriesPerMinute))
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func load(id int) {
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

	counter := 0
	for {
		query = queries[randInt(0, len(queries))]

		result = connectors.Exec("master", map[string]interface{}{
			"query":  query,
			"params": []interface{}{},
		})

		if result == true {
			counter++
			counters[id] = counter
		}
	}
}
