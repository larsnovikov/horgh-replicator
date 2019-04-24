package system

import (
	"horgh-replicator/src/constants"
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

	res := Get(constants.DBReplicator, map[string]interface{}{
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
	result := row.Value

	defer func() {
		_ = res.Close()
	}()

	return result
}

func SetValue(key string, value string) bool {
	query := `INSERT INTO param_values(param_key, param_value) VALUES(?, ?) ON DUPLICATE KEY UPDATE param_value=?;`
	params := []interface{}{
		key,
		value,
		value,
	}

	res := Exec(constants.DBReplicator, map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
