package parser

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/schema"
	"horgh-replicator/src/connectors"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"plugin"
	"time"
)

type BinlogParser struct{}

func (m *BinlogParser) ParseBinLog(slave slave.Slave, e *canal.RowsEvent, n int) error {
	masterFields := slave.GetConfig().Master.Fields
	slaveFields := slave.GetConnector().GetFields()

	params := make(map[string]interface{})
	var fieldType string
	var value interface{}
	for key, fieldName := range masterFields {
		fieldType = slaveFields[fieldName].Mode
		row := e.Rows[0]
		if len(e.Rows) > 1 {
			row = e.Rows[1]
		}
		// prepare value before save
		value = m.beforeSave(slaveFields[fieldName].BeforeSave, row[key])
		// prepare value type
		m.prepareType(fieldName, fieldType, value, params)
		// set values to storage
		slave.GetConnector().SetParams(params)
	}

	return nil
}

func (m *BinlogParser) beforeSave(beforeSave connectors.ConfigBeforeSave, value interface{}) interface{} {
	if beforeSave.Handler == "" {
		return value
	}

	mod := fmt.Sprintf(constants.PluginPath, beforeSave.Handler)
	plug, err := plugin.Open(mod)
	if err != nil {
		log.Fatalf(constants.ErrorCachePluginError, err)
	}

	symHandler, err := plug.Lookup("Handler")
	if err != nil {
		log.Fatalf(constants.ErrorCachePluginError, err)
	}

	var handler helpers.Handler
	handler, ok := symHandler.(helpers.Handler)
	if !ok {
		log.Fatalf(constants.ErrorCachePluginError, "unexpected type from module symbol")
	}

	return handler.Handle(value, beforeSave.Params)
}

func (m *BinlogParser) prepareType(fieldName string, fieldType string, value interface{}, params map[string]interface{}) {
	switch fieldType {
	case "bool":
		params[fieldName] = value.(bool)
	case "int":
		params[fieldName] = value.(int32)
	case "string":
		params[fieldName] = value.(string)
	case "float":
		params[fieldName] = value.(float64)
	case "timestamp":
		t, _ := time.Parse("2006-01-02 15:04:05", value.(string))
		params[fieldName] = t
	}
}

func (m *BinlogParser) dateTimeHelper(e *canal.RowsEvent, n int, columnName string) time.Time {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_TIMESTAMP {
		panic("Not dateTime type")
	}
	t, _ := time.Parse("2006-01-02 15:04:05", e.Rows[n][columnId].(string))

	return t
}

func (m *BinlogParser) intHelper(e *canal.RowsEvent, n int, columnName string) int64 {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_NUMBER {
		return 0
	}

	switch e.Rows[n][columnId].(type) {
	case int8:
		return int64(e.Rows[n][columnId].(int8))
	case int32:
		return int64(e.Rows[n][columnId].(int32))
	case int64:
		return e.Rows[n][columnId].(int64)
	case int:
		return int64(e.Rows[n][columnId].(int))
	case uint8:
		return int64(e.Rows[n][columnId].(uint8))
	case uint16:
		return int64(e.Rows[n][columnId].(uint16))
	case uint32:
		return int64(e.Rows[n][columnId].(uint32))
	case uint64:
		return int64(e.Rows[n][columnId].(uint64))
	case uint:
		return int64(e.Rows[n][columnId].(uint))
	}
	return 0
}

func (m *BinlogParser) floatHelper(e *canal.RowsEvent, n int, columnName string) float64 {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_FLOAT {
		panic("Not float type")
	}

	switch e.Rows[n][columnId].(type) {
	case float32:
		return float64(e.Rows[n][columnId].(float32))
	case float64:
		return float64(e.Rows[n][columnId].(float64))
	}
	return float64(0)
}

func (m *BinlogParser) boolHelper(e *canal.RowsEvent, n int, columnName string) bool {

	val := m.intHelper(e, n, columnName)
	if val == 1 {
		return true
	}
	return false
}

func (m *BinlogParser) stringHelper(e *canal.RowsEvent, n int, columnName string) string {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type == schema.TYPE_ENUM {

		values := e.Table.Columns[columnId].EnumValues
		if len(values) == 0 {
			return ""
		}
		if e.Rows[n][columnId] == nil {
			//Если в енум лежит нуул ставим пустую строку
			return ""
		}

		return values[e.Rows[n][columnId].(int64)-1]
	}

	value := e.Rows[n][columnId]

	switch value := value.(type) {
	case []byte:
		return string(value)
	case string:
		return value
	}
	return ""
}

func (m *BinlogParser) getBinlogIdByName(e *canal.RowsEvent, name string) int {
	for id, value := range e.Table.Columns {
		if value.Name == name {
			return id
		}
	}
	panic(fmt.Sprintf(constants.ErrorNoColumn, name, e.Table.Schema, e.Table.Name))
}
