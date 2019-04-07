package models

import (
	"encoding/json"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/replication"
	"go-binlog-replication/src/connectors2"
	"go-binlog-replication/src/connectors2/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"io/ioutil"
)

type AbstractConnector interface {
	Insert() bool
	Update() bool
	Delete() bool
	GetConfigStruct() interface{}
	SetConfig(interface{})
	SetParams(map[string]interface{})
	ParseKey([]interface{})
	GetFields() map[string]connectors2.ConfigField
	GetTable() string
	Connection() helpers.Storage
}

type Slave struct {
	connector AbstractConnector
	config    Config
	key       string
	table     string
}

type Config struct {
	Master ConfigMaster `json:"master"`
	Slave  interface{}  `json:"slave"`
}

type ConfigMaster struct {
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}

var slave Slave

// make model, read config by modelName, set var model
func MakeSlave(modelName string) Slave {
	slave = Slave{}
	// TODO в зависимости от типа слейва подключаем разные коннекторы
	slave.connector = &mysql.Model{}

	// добавляем к базовому конфигу конфиг коннектора
	slave.config.Slave = slave.connector.GetConfigStruct()

	// make config
	file := helpers.ReadConfig(modelName)
	byteValue, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(byteValue, &slave.config)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}

	// set model params from config
	slave.connector.SetConfig(slave.config.Slave)

	return slave
}

func (slave Slave) GetBeforeSaveMethods() map[string]func(interface{}, []interface{}) interface{} {
	// TODO fix violent pornography
	functions := map[string]func(interface{}, []interface{}) interface{}{
		"SetValue": func(value interface{}, params []interface{}) interface{} {
			return helpers.SetValue(value, params)
		},
	}

	return functions
}

func (slave Slave) GetConfig() Config {
	return slave.config
}

func (slave Slave) GetConnector() AbstractConnector {
	return slave.connector
}

func (slave Slave) ClearParams() {
	slave.connector.SetParams(map[string]interface{}{})
}

func (slave Slave) TableName() string {
	return slave.GetConnector().GetTable()
}

func (slave Slave) BeforeSave() bool {
	return true
}

func (slave Slave) Insert(header *replication.EventHeader) bool {
	if slave.BeforeSave() == true && slave.connector.Insert() == true {
		log.Infof(constants.MessageInserted, header.Timestamp, slave.TableName(), header.LogPos)
		return true
	}

	slave.logError("insert")

	return false
}

func (slave Slave) Update(header *replication.EventHeader) bool {
	if slave.BeforeSave() == true && slave.connector.Update() == true {
		log.Infof(constants.MessageUpdated, header.Timestamp, slave.TableName(), header.LogPos)
		return true
	}

	slave.logError("update")

	return false
}

func (slave Slave) Delete(header *replication.EventHeader) bool {
	if slave.connector.Delete() == true {
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
