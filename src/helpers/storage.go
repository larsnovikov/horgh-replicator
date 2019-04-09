package helpers

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/constants"
	"strconv"
	"time"
)

type Storage interface {
	Ping() bool
	Exec(params map[string]interface{}) bool
}

var retryCounter = map[string]int{
	constants.DBReplicator: 0,
	constants.DBSlave:      0,
	constants.DBMaster:     0,
}

func Retry(storageType string, cred Credentials, connection Storage, method func(connection Storage, dbName string) interface{}) interface{} {
	if retryCounter[storageType] > cred.RetryAttempts {
		log.Fatal(fmt.Sprintf(constants.ErrorDBConnect, storageType))
	}

	log.Infof(constants.MessageRetryConnect, storageType, strconv.Itoa(cred.RetryTimeout))

	time.Sleep(time.Duration(cred.RetryTimeout) * time.Second)
	retryCounter[storageType]++

	return method(connection, storageType)
}
