package slave

import (
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/system"
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
	return helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB).DBname
}

func (User) getType() string {
	return helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).Type
}

func (user *User) Insert() bool {
	query := `INSERT INTO "` + user.TableName() + `"(id, name, status, created) VALUES($1, $2, $3, $4);`
	params := []interface{}{
		user.Id,
		user.Name,
		user.Status,
		user.Created,
	}

	res := system.Exec(user.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (user *User) Update() bool {
	query := `UPDATE "` + user.TableName() + `" SET name=$1, status=$2, created=$3 WHERE id=$4;`
	params := []interface{}{
		user.Name,
		user.Status,
		user.Created,
		user.Id,
	}

	res := system.Exec(user.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (user *User) Delete() bool {
	query := `DELETE FROM "` + user.TableName() + `" WHERE id=$1`
	params := []interface{}{
		user.Id,
	}

	res := system.Exec(user.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
