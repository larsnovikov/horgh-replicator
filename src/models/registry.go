package models

import "go-binlog-replication/src/models/slave"

type AbstractModel interface {
	TableName() string
	SchemaName() string
	Insert() bool
	Update() bool
	Delete() bool
	ParseKey([]interface{})
	BeforeSave() bool
}

func GetModel(name string) interface{ AbstractModel } {
	var model func() interface{ AbstractModel }
	switch name {
	case "user":
		model = func() interface{ AbstractModel } {
			return &slave.User{}
		}
	case "post":
		model = func() interface{ AbstractModel } {
			return &slave.Post{}
		}
	}

	output := model()

	return output
}
