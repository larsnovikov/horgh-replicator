package slave

import (
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"runtime"
	"strconv"
	"time"
)

var AllowHandling = true

func save(c chan func() bool) {
	for {
		if AllowHandling == true {
			method := <-c
			method()
		} else {
			// todo откат транзакции
			runtime.Goexit()
		}
	}
}

func Stop() bool {
	// stop handle saving
	log.Infof(constants.MessageStopHandlingSave)
	AllowHandling = false

	log.Infof(constants.MessageWait, strconv.Itoa(5), "seconds")
	time.Sleep(5 * time.Second)

	return true
}
