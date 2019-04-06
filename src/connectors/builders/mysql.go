package builders

import (
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/models"
)

func Insert(slave models.Slave) bool {
	var params []interface{}
	for _, value := range slave.Fields() {
		// берем values и формируем строку вида (id, name, status, created) и (?, ?, ?, ?) и массив params
	}

	// собираем запрос вида
	query := `INSERT INTO ` + slave.TableName() + `(id, name, status, created) VALUES(?, ?, ?, ?);`

	return connectors.Exec(slave.Type(), map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func Update(slave models.Slave) bool {
	return true
}

func Delete(slave models.Slave) bool {
	return true
}
