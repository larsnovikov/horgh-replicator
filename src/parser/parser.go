package parser

import (
	"github.com/siddontang/go-mysql/canal"
	"horgh-replicator/src/connectors"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/plugins"
	"strconv"
	"time"
)

type BinlogParser struct{}

func (m *BinlogParser) ParseBinLog(slave slave.Slave, e *canal.RowsEvent, n int) error {
	row := e.Rows[0]
	if len(e.Rows) > 1 {
		row = e.Rows[1]
	}

	return ParseRow(slave, row)
}

func ParseRow(slave slave.Slave, row []interface{}) error {
	masterFields := slave.GetConfig().Master.Fields
	slaveFields := slave.GetConnector().GetFields()

	params := make(map[string]interface{})
	var fieldType string
	var value interface{}
	for key, fieldName := range masterFields {
		fieldType = slaveFields[fieldName].Mode
		// prepare value before save
		value = beforeSave(slaveFields[fieldName].BeforeSave, row[key])
		// prepare value type
		prepareType(fieldName, fieldType, value, params)
		// set values to storage
		slave.GetConnector().SetParams(params)
	}

	return nil
}

func beforeSave(beforeSave connectors.ConfigBeforeSave, value interface{}) interface{} {
	if beforeSave.Handler == "" {
		return value
	}

	return plugins.Handle(beforeSave, value)
}

func prepareType(fieldName string, fieldType string, value interface{}, params map[string]interface{}) {
	switch fieldType {
	case "bool":
		switch value.(type) {
		case int8:
			if value.(int8) == 1 {
				params[fieldName] = true
			} else {
				params[fieldName] = false
			}
		case string:
			if value.(string) == "1" || value.(string) == "true" {
				params[fieldName] = true
			} else {
				params[fieldName] = false
			}
		}
	case "int":
		switch value.(type) {
		case string:
			params[fieldName], _ = strconv.Atoi(value.(string))
		case int8:
			params[fieldName] = int64(value.(int8))
		case int32:
			params[fieldName] = int64(value.(int32))
		case int64:
			params[fieldName] = value.(int64)
		case int:
			params[fieldName] = int64(value.(int))
		case uint8:
			params[fieldName] = int64(value.(uint8))
		case uint16:
			params[fieldName] = int64(value.(uint16))
		case uint32:
			params[fieldName] = int64(value.(uint32))
		case uint64:
			params[fieldName] = int64(value.(uint64))
		case uint:
			params[fieldName] = int64(value.(uint))
		default:
			params[fieldName] = 0
		}
	case "string":
		if value == nil {
			params[fieldName] = ""
		} else {
			switch value := value.(type) {
			case []byte:
				params[fieldName] = string(value)
			case string:
				params[fieldName] = value
			default:
				params[fieldName] = value
			}
		}
	case "float":
		switch value.(type) {
		case string:
			f, _ := strconv.ParseFloat(value.(string), 64)
			params[fieldName] = float64(f)
		case float32:
			params[fieldName] = float64(value.(float32))
		case float64:
			params[fieldName] = float64(value.(float64))
		default:
			params[fieldName] = float64(0)
		}
	case "timestamp":
		t, _ := time.Parse("2006-01-02 15:04:05", value.(string))
		params[fieldName] = t
	}
}
