package parser

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models"
	"runtime/debug"
	"strconv"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

var curPosition mysql.Position
var curCanal *canal.Canal

func canOperateTable(tableName string) bool {
	for _, v := range helpers.GetTables() {
		if v == tableName {
			return true
		}
	}

	return false
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r, " ", string(debug.Stack()))
		}
	}()

	if curCanal == nil {
		canalTmp, err := getDefaultCanal()
		if err != nil {
			log.Fatal(constants.ErrorMysqlCanal)
		}
		curCanal = canalTmp
	}

	if curPosition.Pos == 0 {
		curPosition = getMasterPos(curCanal, false)
	} else {
		curPosition = getMasterPos(curCanal, true)
	}

	var n int
	var k int

	if canOperateTable(e.Table.Name) == false {
		return nil
	}

	model := models.GetModel(e.Table.Name)

	switch e.Action {
	case canal.DeleteAction:
		for _, row := range e.Rows {
			model.ParseKey(row)
			if model.Delete() == true {
				setMasterPosFromCanal(curPosition)
				log.Infof(constants.MessageDeleted, e.Table)
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
		key := e.Table.Schema + "." + e.Table.Name
		switch key {
		case model.SchemaName() + "." + model.TableName():
			h.GetBinLogData(model, e, i)

			if e.Action == canal.UpdateAction {
				oldModel := models.GetModel(e.Table.Name)
				h.GetBinLogData(oldModel, e, i-1)
				if model.Update() == true {
					setMasterPosFromCanal(curPosition)
					log.Infof(constants.MessageUpdated, e.Table)
				}
			} else {
				if model.Insert() == true {
					setMasterPosFromCanal(curPosition)
					log.Infof(constants.MessageInserted, e.Table)
				}
			}
		}

	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func BinlogListener() {
	helpers.MakeCredentials()

	c, err := getDefaultCanal()
	if err == nil {
		coords, err := getMasterPosFromCanal(c, false)
		if err == nil {
			c.SetEventHandler(&binlogHandler{})
			c.RunFrom(coords)
		}
	}
}

func getMasterPosFromCanal(c *canal.Canal, force bool) (mysql.Position, error) {
	// try to get coords from storage
	if force == false {
		position, err := strconv.ParseUint(models.GetValue(constants.LastPositionPos), 10, 32)
		if err == nil {
			pos := mysql.Position{
				models.GetValue(constants.LastPositionName),
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

func setMasterPosFromCanal(position mysql.Position) {
	// save position
	models.SetValue(constants.LastPositionPos, fmt.Sprint(position.Pos))
	models.SetValue(constants.LastPositionName, position.Name)

	curPosition = position
}

func getMasterPos(canal *canal.Canal, force bool) mysql.Position {

	coords, err := getMasterPosFromCanal(canal, force)
	if err != nil {
		log.Fatal(constants.ErrorMysqlPosition)
	}

	return coords
}

func getDefaultCanal() (*canal.Canal, error) {
	master := helpers.GetCredentials(constants.DBMaster)

	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", master.Host, master.Port)
	cfg.User = master.User
	cfg.Password = master.Pass
	cfg.Flavor = master.Type

	cfg.Dump.ExecutionPath = ""

	return canal.NewCanal(cfg)
}

func showPos(pos mysql.Position, from string) {
	// log.Info(fmt.Sprintf(constants.MessagePosFrom, from, fmt.Sprint(pos.Pos), pos.Name))
}
