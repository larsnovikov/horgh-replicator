package parser

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/slave"
	"go-binlog-replication/src/models/system"
	"math/rand"
	"runtime/debug"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

var curCanal *canal.Canal

func (h *binlogHandler) canOperate(logTableName string) bool {
	for _, tableName := range helpers.GetTables() {
		if tableName == logTableName {
			return true
		}
	}

	return false
}

func (h *binlogHandler) prepareCanal() {
	// build canal if not exists yet
	if curCanal == nil {
		canalTmp, err := getDefaultCanal()
		if err != nil {
			log.Fatal(constants.ErrorMysqlCanal)
		}
		curCanal = canalTmp
	}
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r, " ", string(debug.Stack()))
		}
	}()

	h.prepareCanal()
	if h.canOperate(e.Table.Name) == false {
		return nil
	}

	slave.GetSlaveByName(e.Table.Name).ClearParams()

	var n int
	var k int

	positionSet := func() {
		// TODO err handler
		pos, _ := curCanal.GetMasterPos()
		pos.Pos = e.Header.LogPos
		SetPosition(e.Table.Name, pos)
		return
	}

	switch e.Action {
	case canal.DeleteAction:
		for _, row := range e.Rows {
			slave.GetSlaveByName(e.Table.Name).GetConnector().ParseKey(row)
			slave.GetSlaveByName(e.Table.Name).Delete(e.Header, positionSet)
		}

		return nil
	case canal.UpdateAction:
		n = 1
		k = 2
	case canal.InsertAction:
		n = 0
		k = 1
	}

	for i := n; i < len(e.Rows); i += k {
		h.ParseBinLog(slave.GetSlaveByName(e.Table.Name), e, i)

		if e.Action == canal.UpdateAction {
			slave.GetSlaveByName(e.Table.Name).Update(e.Header, positionSet)
		} else {
			slave.GetSlaveByName(e.Table.Name).Insert(e.Header, positionSet)
		}
	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func BinlogListener() {
	// set position keys

	c, err := getDefaultCanal()
	if err == nil {
		coords := getMinPosition(c.SyncedPosition())
		c.SetEventHandler(&binlogHandler{})
		err = c.RunFrom(coords)
	}
}

func getDefaultCanal() (*canal.Canal, error) {
	// try to connect to check credentials
	system.Exec(constants.DBMaster, map[string]interface{}{
		"query":  "SELECT 1",
		"params": make([]interface{}, 0),
	})

	master := helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB)

	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", master.Host, master.Port)
	cfg.User = master.User
	cfg.Password = master.Pass
	cfg.Flavor = master.Type
	cfg.ServerID = uint32(rand.Intn(9999999999))

	cfg.Dump.ExecutionPath = ""

	return canal.NewCanal(cfg)
}

func OnRotate(roateEvent *replication.RotateEvent) error {
	return nil
}

func OnTableChanged(schema string, table string) error {
	return nil
}

func OnDDL(nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	return nil
}

func showPos(pos mysql.Position, from string) {
	// log.Info(fmt.Sprintf(constants.MessagePosFrom, from, fmt.Sprint(pos.Pos), pos.Name))
}
