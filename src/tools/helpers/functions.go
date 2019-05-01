package helpers

import (
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/models/system"
	"time"
)

const (
	Timeout = 10
)

func GetHeader() (header slave.Header, positionSet func()) {
	t := time.Now()
	header = slave.Header{
		Timestamp: uint32(t.Unix()),
		LogPos:    Position.Pos,
	}

	// dont set position for every row. Set it for all rows once
	positionSet = func() {
		return
	}

	return header, positionSet
}

func Wait(cond func() bool) {
	for {
		// waiting until save.channel is empty
		time.Sleep(Timeout * time.Second)
		if cond() {
			break
		}
	}
}

func SetPosition() {
	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := helpers.MakeHash(dbName, Table)

	system.SetPosition(hash, Position)

	log.Infof(constants.MessagePositionUpdated, Table)
}
