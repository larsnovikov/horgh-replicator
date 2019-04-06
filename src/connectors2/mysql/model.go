package mysql

import (
	"fmt"
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/connectors2"
	"strings"
)

const (
	Type   = "mysql"
	Insert = "INSERT INTO %s(%s) VALUES(%s);"
	Update = "UPDATE %s SET %s WHERE %s=?;"
	Delete = "DELETE FROM %s WHERE %s=?"
)

type Model struct {
	table  string
	key    string
	fields map[string]connectors2.ConfigField
	params map[string]interface{}
}

func (model *Model) GetConfigStruct() interface{} {
	return &connectors2.ConfigSlave{}
}

func (model *Model) SetConfig(config interface{}) {
	model.table = config.(*connectors2.ConfigSlave).Table

	model.fields = make(map[string]connectors2.ConfigField)
	for _, value := range config.(*connectors2.ConfigSlave).Fields {
		if model.key == "" && value.Key == true {
			model.key = value.Name
		}

		model.fields[value.Name] = value
	}
}

func (model *Model) Insert() bool {
	var params []interface{}
	var fieldNames []string
	var fieldValues []string

	for _, value := range model.fields {
		// берем values и формируем строку вида (id, name, status, created) и (?, ?, ?, ?) и массив params
		fieldNames = append(fieldNames, value.Name)
		fieldValues = append(fieldValues, "?")

		params = append(params, model.params[value.Name])
	}

	query := fmt.Sprintf(Insert, model.table, strings.Join(fieldNames, ","), strings.Join(fieldValues, ","))

	return connectors.Exec(Type, map[string]interface{}{
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

	query := fmt.Sprintf(Update, model.table, strings.Join(fields, ","), model.key)

	return connectors.Exec(Type, map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Delete() bool {
	var params []interface{}
	query := fmt.Sprintf(Delete, model.table, model.key)

	params = append(params, model.params[model.key])

	return connectors.Exec(Type, map[string]interface{}{
		"query":  query,
		"params": params,
	})
}
