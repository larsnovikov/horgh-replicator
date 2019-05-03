package exit

import (
	"github.com/siddontang/go-log/log"
	"os"
)

func HandleBefore() bool {
	FirstStop = false
	for _, method := range BeforeExitPool {
		if method() == false {
			return false
		}
	}

	return true
}

func Fatal(msg string, args ...interface{}) {
	if FirstStop == false || HandleBefore() {
		log.Warnf(msg, args...)
		os.Exit(1)
	}
}
