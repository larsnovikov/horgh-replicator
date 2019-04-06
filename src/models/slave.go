package models

import (
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/replication"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"io/ioutil"
)

type AbstractConnector interface {
	Insert(slave Slave) bool
	Update(slave Slave) bool
	Delete(slave Slave) bool
}

type Slave struct {
	connector AbstractConnector
	config    Config
	fields    map[string]ConfigField
	key       string
	table     string
	schema    string
}

type Config struct {
	Master ConfigMaster `json:"master"`
	Slave  ConfigSlave  `json:"slave"`
}

type ConfigMaster struct {
	Table string `json:"table"`
}

type ConfigSlave struct {
	Fields []ConfigField `json:"fields"`
}

type ConfigBeforeSave struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type ConfigField struct {
	Name       string           `json:"name"`
	Key        bool             `json:"key"`
	Mode       string           `json:"mode"`
	BeforeSave ConfigBeforeSave `json:"beforeSave"`
}

var slave Slave

// make model, read config by modelName, set var model
func MakeSlave(modelName string) {
	slave = Slave{}
	// make config
	file := helpers.ReadConfig(modelName)
	byteValue, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(byteValue, &slave.config)
	defer func() {
		_ = file.Close()
	}()
	log.Fatal(err)

	slave.fields = make(map[string]ConfigField)
	for _, val := range slave.config.Slave.Fields {
		// set key
		if val.Key == true && slave.key == "" {
			slave.key = val.Name
		}

		// set fields
		slave.fields[val.Name] = val
	}

	// set table
	slave.table = slave.config.Master.Table

	// set schema TODO по идее в конфиге модели должен быть указан название БД но это не точно!
	slave.schema = helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB).DBname
	fmt.Println(slave.fields)

	// TODO don'forget to set connector
	log.Fatal("debug stop")
}

func (slave Slave) Type() string {
	return "111"
}

func (slave Slave) Fields() map[string]ConfigField {
	return slave.fields
}

func (slave Slave) TableName() string {
	return slave.table
}

func (slave Slave) SchemaName() string {
	return slave.schema
}

func (slave Slave) BeforeSave() bool {
	return true
}

func (slave Slave) Insert(header *replication.EventHeader) bool {
	if slave.BeforeSave() == true && slave.connector.Insert(slave) == true {
		log.Infof(constants.MessageInserted, header.Timestamp, slave.TableName(), header.LogPos)
		return true
	}

	slave.logError("insert")

	return false
}

func (slave Slave) Update(header *replication.EventHeader) bool {
	if slave.BeforeSave() == true && slave.connector.Update(slave) == true {
		log.Infof(constants.MessageUpdated, header.Timestamp, slave.TableName(), header.LogPos)
		return true
	}

	slave.logError("update")

	return false
}

func (slave Slave) Delete(header *replication.EventHeader) bool {
	if slave.connector.Delete(slave) == true {
		log.Infof(constants.MessageDeleted, header.Timestamp, slave.TableName(), header.LogPos)
		return true
	}

	slave.logError("delete")

	return false
}

func (slave Slave) logError(operationType string) {
	out, _ := json.Marshal(slave)

	log.Warnf(constants.ErrorSave, operationType, slave.TableName(), string(out))
}
