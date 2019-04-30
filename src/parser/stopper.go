package parser

import (
	"horgh-replicator/src/models/slave"
)

func StopListen() {
	Stop()
	slave.Stop()
}
