package models

import (
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/replication"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/models/slave"
)

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

func Insert(model AbstractModel, header *replication.EventHeader) bool {
	if model.BeforeSave() == true && model.Insert() == true {
		log.Infof(constants.MessageInserted, header.Timestamp, model.TableName(), header.LogPos)
		return true
	}

	return false
}

func Update(model AbstractModel, header *replication.EventHeader) bool {
	if model.BeforeSave() == true && model.Update() == true {
		log.Infof(constants.MessageUpdated, header.Timestamp, model.TableName(), header.LogPos)
		return true
	}

	return false
}

func Delete(model AbstractModel, header *replication.EventHeader) bool {
	if model.Delete() == true {
		log.Infof(constants.MessageDeleted, header.Timestamp, model.TableName(), header.LogPos)
		return true
	}

	return false
}
