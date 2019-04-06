package mysql

import (
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/connectors2"
)

const Type = "mysql"

type Model struct {
	Table string
}

func (model Model) GetConfig() interface{} {
	return connectors2.ConfigSlave{}
}

func (model Model) Insert() bool {
	var params []interface{}
	//for _, value := range model.Fields {
	//	// берем values и формируем строку вида (id, name, status, created) и (?, ?, ?, ?) и массив params
	//}

	// собираем запрос вида
	query := `INSERT INTO ` + model.Table + `(id, name, status, created) VALUES(?, ?, ?, ?);`

	return connectors.Exec(Type, map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model Model) Update() bool {
	return true
}

func (model Model) Delete() bool {
	return true
}
