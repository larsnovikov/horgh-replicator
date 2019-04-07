package mysql

import (
	"fmt"
	"go-binlog-replication/src/connectors2"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strings"
)

const (
	Type   = "mysql"
	Insert = "INSERT INTO %s(%s) VALUES(%s);"
	Update = "UPDATE %s SET %s WHERE %s=?;"
	Delete = "DELETE FROM %s WHERE %s=?"
)

type Model struct {
	table       string
	key         string
	keyPosition int
	fields      map[string]connectors2.ConfigField
	params      map[string]interface{}
}

func (model Model) GetFields() map[string]connectors2.ConfigField {
	return model.fields
}

func (model *Model) ParseKey(row []interface{}) {
	// TODO в зависимости от типа поля ключа, тут могут быть разные приведения типов
	// а если будет составной ключ вообще будет тяжко
	model.params[model.key] = int(row[model.keyPosition].(int32))
}

func (model *Model) GetConfigStruct() interface{} {
	return &connectors2.ConfigSlave{}
}

func (model *Model) GetTable() string {
	return model.table
}

func (model *Model) SetConfig(config interface{}) {
	model.table = config.(*connectors2.ConfigSlave).Table

	model.fields = make(map[string]connectors2.ConfigField)
	for pos, value := range config.(*connectors2.ConfigSlave).Fields {
		if model.key == "" && value.Key == true {
			model.key = value.Name
			model.keyPosition = pos
		}

		model.fields[value.Name] = value
	}
}

func (model *Model) SetParams(params map[string]interface{}) {
	model.params = params
}

func (model *Model) Insert() bool {
	var params []interface{}
	var fieldNames []string
	var fieldValues []string

	for _, value := range model.fields {
		fieldNames = append(fieldNames, value.Name)
		fieldValues = append(fieldValues, "?")

		params = append(params, model.params[value.Name])
	}

	query := fmt.Sprintf(Insert, model.table, strings.Join(fieldNames, ","), strings.Join(fieldValues, ","))

	return model.Connection().Exec(map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Update() bool {
	var params []interface{}
	var fields []string

	for _, value := range model.fields {
		fields = append(fields, value.Name+"=?")

		params = append(params, model.params[value.Name])
	}

	// add key to params
	params = append(params, model.params[model.key])

	query := fmt.Sprintf(Update, model.table, strings.Join(fields, ","), model.key)

	return model.Connection().Exec(map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Delete() bool {
	var params []interface{}
	query := fmt.Sprintf(Delete, model.table, model.key)

	params = append(params, model.params[model.key])

	return model.Connection().Exec(map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Connection() helpers.Storage {
	helpers.ConnPool.Slave = GetConnection(helpers.ConnPool.Slave, constants.DBSlave).(helpers.Storage)
	return helpers.ConnPool.Slave
}
