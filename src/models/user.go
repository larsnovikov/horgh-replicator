package models

import (
	"go-binlog-replication/src/helpers"
	"time"
)

type User struct {
	Id      int       `gorm:"column:id"`
	Name    string    `gorm:"column:name"`
	Status  string    `gorm:"column:status"`
	Created time.Time `gorm:"column:created"`
}

func (user *User) ParseKey(row []interface{}) {
	user.Id = int(row[0].(int32))
}

func (User) TableName() string {
	return "User"
}

func (User) SchemaName() string {
	return helpers.GetCredentials("master").DBname
}

func getType() string {
	return helpers.GetCredentials("slave").Type
}

func (user *User) Insert() bool {
	query := `INSERT INTO ` + user.TableName() + `(id, name, status, created) VALUES(?, ?, ?, ?);`
	params := []interface{}{
		user.Id,
		user.Name,
		user.Status,
		user.Created,
	}

	res := helpers.Exec(getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (user *User) Update() bool {
	query := `UPDATE ` + user.TableName() + ` SET name=?, status=?, created=? WHERE id=?;`
	params := []interface{}{
		user.Name,
		user.Status,
		user.Created,
		user.Id,
	}

	res := helpers.Exec(getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (user *User) Delete() bool {
	query := `DELETE FROM ` + user.TableName() + ` WHERE id=?`
	params := []interface{}{
		user.Id,
	}

	res := helpers.Exec(getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
