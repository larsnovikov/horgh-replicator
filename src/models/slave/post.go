package slave

import (
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
	return helpers.GetCredentials(constants.DBMaster).DBname
}

func (Post) getType() string {
	return helpers.GetCredentials(constants.DBSlave).Type
}

func (post *Post) Insert() bool {
	query := `INSERT INTO ` + post.SchemaName() + `.` + post.TableName() + `(id, title, text, created) VALUES(?, ?, ?, ?);`
	params := []interface{}{
		post.Id,
		post.Title,
		post.Text,
		post.Created,
	}

	res := helpers.Exec(post.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (post *Post) Update() bool {
	query := `ALTER TABLE ` + post.SchemaName() + `.` + post.TableName() + ` UPDATE title=?, text=?, created=? WHERE id=?;`
	params := []interface{}{
		post.Title,
		post.Text,
		post.Created,
		post.Id,
	}

	res := helpers.Exec(post.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}

func (post *Post) Delete() bool {
	query := `ALTER TABLE ` + post.SchemaName() + `.` + post.TableName() + ` DELETE WHERE id=?`
	params := []interface{}{
		post.Id,
	}

	res := helpers.Exec(post.getType(), map[string]interface{}{
		"query":  query,
		"params": params,
	})

	return res
}
