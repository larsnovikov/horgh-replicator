package slave

import (
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"runtime"
	"strconv"
	"time"
)

var AllowHandling = true

func save(c chan helpers.QueryAction) {
	for {
		msg := <-c
		if AllowHandling == true {
			msg.Method()
		} else {
			msg.StopMethod()
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
