package slave

import (
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"time"
)

type User struct {
	Id      int       `gorm:"column:id"`
	Name    string    `gorm:"column:name"`
	Status  string    `gorm:"column:status"`
	Created time.Time `gorm:"column:created"`
}

func (user *User) BeforeSave() bool {
	user.Status = "***"

	return true
}

func (user *User) ParseKey(row []interface{}) {
	user.Id = int(row[0].(int32))
}

func (User) TableName() string {
	return "user"
}

func (User) SchemaName() string {
	return helpers.GetCredentials(constants.DBMaster).DBname
}

func (User) getType() string {
	return helpers.GetCredentials(constants.DBSlave).Type
}

func (user *User) Insert() bool {
	query := `INSERT INTO ` + user.SchemaName() + `.` + user.TableName() + `(id, name, status, created) VALUES(?, ?, ?, ?);`
	params := []interface{}{
		user.Id,
		user.Name,
		user.Status,
		user.Created,
	}

	res := helpers.Exec(user.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (user *User) Update() bool {
	query := `ALTER TABLE ` + user.SchemaName() + `.` + user.TableName() + ` UPDATE name=?, status=?, created=? WHERE id=?;`
	params := []interface{}{
		user.Name,
		user.Status,
		user.Created,
		user.Id,
	}

	res := helpers.Exec(user.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (user *User) Delete() bool {
	query := `ALTER TABLE ` + user.SchemaName() + `.` + user.TableName() + ` DELETE WHERE id=?`
	params := []interface{}{
		user.Id,
	}

	res := helpers.Exec(user.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
