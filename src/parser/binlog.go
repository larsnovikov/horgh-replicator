package parser

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/models/system"
	"horgh-replicator/src/tools/exit"
	"runtime"
	"runtime/debug"
	"strconv"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

var AllowHandling = true

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
	if AllowHandling == false {
		runtime.Goexit()
	}

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

	getCalculatedPos := func() mysql.Position {
		logFile := PrevPosition[e.Table.Name].Name
		if e.Header.LogPos < PrevPosition[e.Table.Name].Pos {
			// log file changed
			newLog := strconv.Itoa(GetLogFileSuffix(logFile) + 1)
			for len(newLog) < 6 {
				newLog = "0" + newLog
			}
			logFile = helpers.GetMasterLogFilePrefix() + newLog
			log.Infof(constants.MessageLogFileChanged, e.Table, logFile)
		}
		return mysql.Position{
			Name: logFile,
			Pos:  e.Header.LogPos,
		}
	}

	positionSet := func() {
		SetPosition(e.Table.Name, getCalculatedPos())
		return
	}

	header := slave.Header{
		Timestamp: e.Header.Timestamp,
		LogPos:    e.Header.LogPos,
	}

	switch e.Action {
	case canal.DeleteAction:
		for _, row := range e.Rows {
			slave.GetSlaveByName(e.Table.Name).GetConnector().ParseKey(row)
			if SaveLocks[e.Table.Name] == false || canSave(getCalculatedPos(), e.Table.Name) {
				slave.GetSlaveByName(e.Table.Name).Delete(&header, positionSet)
				SaveLocks[e.Table.Name] = false
			} else {
				log.Infof(constants.MessageIgnoreDelete, header.Timestamp, e.Table.Name, header.LogPos)
				return nil
			}
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
			if SaveLocks[e.Table.Name] == false || canSave(getCalculatedPos(), e.Table.Name) {
				slave.GetSlaveByName(e.Table.Name).Update(&header, positionSet)
				SaveLocks[e.Table.Name] = false
			} else {
				log.Infof(constants.MessageIgnoreUpdate, header.Timestamp, e.Table.Name, header.LogPos)
				return nil
			}
		} else {
			if SaveLocks[e.Table.Name] == false || canSave(getCalculatedPos(), e.Table.Name) {
				slave.GetSlaveByName(e.Table.Name).Insert(&header, positionSet)
				SaveLocks[e.Table.Name] = false
			} else {
				log.Infof(constants.MessageIgnoreInsert, header.Timestamp, e.Table.Name, header.LogPos)
				return nil
			}
		}
	}
	return nil
}

func canSave(pos mysql.Position, table string) bool {
	saved := GetSavedPos(table)
	lowPosition := GetLowPosition(pos, saved)

	if pos.Name == saved.Name && pos.Pos == saved.Pos {
		return false
	}
	if lowPosition.Name == saved.Name && lowPosition.Pos == saved.Pos {
		// low == saved => write
		return true
	} else {
		// low == calculated => no!
		return false
	}
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func BinlogListener() {
	c, err := getDefaultCanal()
	if err == nil {
		position, err := c.GetMasterPos()
		if err != nil {
			exit.Fatal(constants.ErrorParserPosition, err)
		}
		coords := getMinPosition(position)
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
	cfg.ServerID = uint32(helpers.GetSlaveId())

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
