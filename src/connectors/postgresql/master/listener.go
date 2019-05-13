package master

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/tools/exit"
	"strconv"
	"time"
)

var AllowHandling = true

func Listen() {
	// TODO get from storage
	lastEventId := 1

	tableName := fmt.Sprintf(logTableName, helpers.GetTable())
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE id>$1`, tableName)
	rows := helpers.ConnPool.Master.Get(helpers.Query{
		Query: query,
		Params: []interface{}{
			lastEventId,
		},
	})

	for rows.Next() {
		var row []interface{}
		err := rows.Scan(row)
		if err != nil {
			exit.Fatal(constants.ErrorLogParsing)
		}
		handleRow(row)
	}
}

func handleRow(row interface{}) bool {
	if AllowHandling == false {
		return false
	}

	// TODO handling row
	return true
}

func stop() bool {
	// stop handle
	log.Infof(constants.MessageStopHandlingBinlog)
	AllowHandling = false

	log.Infof(constants.MessageWait, strconv.Itoa(5), "seconds")
	time.Sleep(5 * time.Second)

	return true
}
