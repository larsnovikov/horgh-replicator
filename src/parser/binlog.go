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
	"strconv"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
	positionNameKey string
	positionPosKey  string
}

var curPosition mysql.Position
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

	// build current position
	if curPosition.Pos == 0 {
		// first row after start, try to get pos from storage
		curPosition = h.getMasterPos(curCanal, false)
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

	slave.GetByName(e.Table.Name).ClearParams()

	var n int
	var k int

	switch e.Action {
	case canal.DeleteAction:
		for _, row := range e.Rows {
			slave.GetByName(e.Table.Name).GetConnector().ParseKey(row)
			if slave.GetByName(e.Table.Name).Delete(e.Header) {
				h.setMasterPosFromCanal(e)
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
		h.ParseBinLog(slave.GetByName(e.Table.Name), e, i)

		//log.Fatal("debug stop")
		if e.Action == canal.UpdateAction {
			if slave.GetByName(e.Table.Name).Update(e.Header) {
				h.setMasterPosFromCanal(e)
			}
		} else {
			if slave.GetByName(e.Table.Name).Insert(e.Header) {
				h.setMasterPosFromCanal(e)
			}
		}
	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func BinlogListener() {
	// set position keys
	// TODO get min position key
	positionPosKey, positionNameKey := helpers.MakeTablePosKey("test.user")

	c, err := getDefaultCanal()
	if err == nil {
		coords, err := getMasterPosFromCanal(c, positionPosKey, positionNameKey, false)
		if err == nil {
			c.SetEventHandler(&binlogHandler{
				positionNameKey: positionNameKey,
				positionPosKey:  positionPosKey,
			})
			err = c.RunFrom(coords)
		}
	}
}

func getMasterPosFromCanal(c *canal.Canal, positionPosKey string, positionNameKey string, force bool) (mysql.Position, error) {
	// try to get coords from storage
	if force == false {
		position, err := strconv.ParseUint(system.GetValue(positionPosKey), 10, 32)
		if err == nil {
			pos := mysql.Position{
				system.GetValue(positionNameKey),
				uint32(position),
			}

			if pos.Pos != 0 && pos.Name != "" {
				showPos(pos, "Storage")
				return pos, nil
			}
		}
	}

	// get coords from mysql
	pos, err := c.GetMasterPos()
	showPos(pos, "MySQL")

	return pos, err
}

func (h *binlogHandler) setMasterPosFromCanal(event *canal.RowsEvent) {
	curPosition.Pos = event.Header.LogPos
	// save position
	system.SetValue(h.positionPosKey, fmt.Sprint(curPosition.Pos))
	system.SetValue(h.positionNameKey, curPosition.Name)
}

func (h *binlogHandler) getMasterPos(canal *canal.Canal, force bool) mysql.Position {
	coords, err := getMasterPosFromCanal(canal, h.positionPosKey, h.positionNameKey, force)
	if err != nil {
		log.Fatal(constants.ErrorMysqlPosition)
	}

	return coords
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
