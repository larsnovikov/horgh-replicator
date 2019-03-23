package models

import (
	"go-binlog-replication/src/helpers"
)

const (
	LastPositionPos  = "last_position_pos"
	LastPositionName = "last_position_name"
)

type replicator struct {
	Id    int    `gorm:"column:id"`
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
		err := res.Scan(&row.Id, &row.Key, &row.Value)
		if err != nil {
			panic(err.Error())
		}
	}

	return row.Value
}

func SetValue(key string, value string) bool {
	query := `UPDATE param_values SET param_value=? WHERE param_key=?;`
	params := []interface{}{
		value,
		key,
	}

	res := helpers.Exec(helpers.GetCredentials("slave").Type, map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
