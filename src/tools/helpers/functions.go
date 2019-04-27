package helpers

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/models/system"
	"time"
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

func Wait() {
	time.Sleep(1 * time.Second)
	for ok := true; ok; ok = slave.GetSlaveByName(Table).GetChannelLen() != 0 {
		// waiting until save.channel is empty
		time.Sleep(1 * time.Second)
	}
}

func SetPosition() {
	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := helpers.MakeHash(dbName, Table)

	posKey, nameKey := helpers.MakeTablePosKey(hash)

	system.SetValue(posKey, fmt.Sprint(Position.Pos))
	system.SetValue(nameKey, Position.Name)

	log.Infof(constants.MessagePositionUpdated, Table)
}
