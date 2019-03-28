package slave

import (
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"time"
)

type Post struct {
	Id      int       `gorm:"column:id"`
	Title   string    `gorm:"column:title"`
	Text    string    `gorm:"column:text"`
	Created time.Time `gorm:"column:created"`
}

func (post *Post) BeforeSave() bool {
	return true
}

func (post *Post) ParseKey(row []interface{}) {
	post.Id = int(row[0].(int32))
}

func (Post) TableName() string {
	return "post"
}

func (Post) SchemaName() string {
	return helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB).DBname
}

func (Post) getType() string {
	return helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).Type
}

func (post *Post) Insert() bool {
	query := `INSERT INTO ` + post.TableName() + `(id, title, text, created) VALUES($1, $2, $3, $4);`
	params := []interface{}{
		post.Id,
		post.Title,
		post.Text,
		post.Created,
	}

	res := connectors.Exec(post.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (post *Post) Update() bool {
	query := `UPDATE ` + post.TableName() + ` SET title=$1, text=$2, created=$3 WHERE id=$4;`
	params := []interface{}{
		post.Title,
		post.Text,
		post.Created,
		post.Id,
	}

	res := connectors.Exec(post.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (post *Post) Delete() bool {
	query := `DELETE FROM ` + post.TableName() + ` WHERE id=$1`
	params := []interface{}{
		post.Id,
	}

	res := connectors.Exec(post.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
