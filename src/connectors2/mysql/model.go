package mysql

import (
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/connectors2"
)

const Type = "mysql"

type Model struct {
	table  string
	key    string
	fields map[string]connectors2.ConfigField
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
	//for _, value := range model.Fields {
	//	// берем values и формируем строку вида (id, name, status, created) и (?, ?, ?, ?) и массив params
	//}

	// собираем запрос вида
	query := `INSERT INTO ` + model.table + `(id, name, status, created) VALUES(?, ?, ?, ?);`

	return connectors.Exec(Type, map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Update() bool {
	return true
}

func (model *Model) Delete() bool {
	return true
}
