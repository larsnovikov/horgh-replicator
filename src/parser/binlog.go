package parser

import (
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"go-binlog-example/src/helpers"
	"go-binlog-example/src/models"
	"runtime/debug"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

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
			fmt.Print(r, " ", string(debug.Stack()))
		}
	}()
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
			model.Delete()
			fmt.Printf("[%s] is deleted \n", e.Table)
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

				model.Update()

				fmt.Printf("[%s] update row\n", e.Table)
			} else {
				model.Insert()

				fmt.Printf("[%s] insert row\n", e.Table)
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
		coords, err := c.GetMasterPos()
		if err == nil {
			c.SetEventHandler(&binlogHandler{})
			c.RunFrom(coords)
		}
	}
}

func getDefaultCanal() (*canal.Canal, error) {
	master := helpers.GetCredentials("master")

	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", master.Host, master.Port)
	cfg.User = master.User
	cfg.Password = master.Pass
	cfg.Flavor = master.Type

	cfg.Dump.ExecutionPath = ""

	return canal.NewCanal(cfg)
}
