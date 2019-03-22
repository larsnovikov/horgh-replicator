package models

import (
	"go-binlog-example/src/helpers"
	"time"
)

type User struct {
	Id      int       `gorm:"column:id"`
	Name    string    `gorm:"column:name"`
	Status  string    `gorm:"column:status"`
	Created time.Time `gorm:"column:created"`
}

func (user User) ParseKey(row []interface{}) {
	user.Id = int(row[0].(int32))
}

func (User) TableName() string {
	return "User"
}

func (User) SchemaName() string {
	return helpers.GetCredentials("master").DBname
}

func (user User) Insert() bool {
	query := `INSERT INTO ` + user.TableName() + `(id, name, status) VALUES(?, ?, ?);`
	params := []interface{}{
		user.Id,
		user.Name,
		user.Status,
	}

	res := helpers.Query(map[string]interface{}{
		"query":  query,
		"params": params,
	})
	helpers.CloseMysql(res)

	return res
}

func (user User) Update() bool {
	query := `UPDATE ` + user.TableName() + ` SET name=?, status=? WHERE id=?;`
	params := []interface{}{
		user.Name,
		user.Status,
		user.Id,
	}

	res := helpers.Query(map[string]interface{}{
		"query":  query,
		"params": params,
	})
	helpers.CloseMysql(res)

	return res
}

func (user User) Delete() bool {
	query := `DELETE FROM ` + user.TableName() + ` WHERE id=?`
	params := []interface{}{
		user.Id,
	}

	res := helpers.Query(map[string]interface{}{
		"query":  query,
		"params": params,
	})
	helpers.CloseMysql(res)

	return res
}
