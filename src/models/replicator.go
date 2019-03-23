package models

import (
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
)

type replicator struct {
	Key   string `gorm:"column:param_key"`
	Value string `gorm:"column:param_value"`
}

func GetValue(key string) string {
	query := `SELECT * FROM param_values WHERE param_key=? LIMIT 1;`
	params := []interface{}{
		key,
	}

	res := helpers.Get(map[string]interface{}{
		"query":  query,
		"params": params,
	})

	var row replicator
	for res.Next() {
		err := res.Scan(&row.Key, &row.Value)
		if err != nil {
			panic(err.Error())
		}
	}

	return row.Value
}

func SetValue(key string, value string) bool {
	query := `INSERT INTO param_values(param_key, param_value) VALUES(?, ?) ON DUPLICATE KEY UPDATE param_value=?;`
	params := []interface{}{
		key,
		value,
		value,
	}

	res := helpers.Exec(constants.DBReplicator, map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
