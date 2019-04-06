package models

import (
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/replication"
	"go-binlog-replication/src/connectors2/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"io/ioutil"
)

type AbstractConnector interface {
	Insert() bool
	Update() bool
	Delete() bool
	GetConfig() interface{}
}

type Slave struct {
	connector AbstractConnector
	config    Config
	key       string
	table     string
	schema    string
}

type Config struct {
	Master ConfigMaster `json:"master"`
	Slave  interface{}  `json:"slave"`
}

type ConfigMaster struct {
	Table string `json:"table"`
}

var slave Slave

// make model, read config by modelName, set var model
func MakeSlave(modelName string) {
	slave = Slave{}
	// TODO в зависимости от типа слейва подключаем разные коннекторы
	slave.connector = mysql.Model{}

	// добавляем к базовому конфигу конфиг коннектора
	slave.config.Slave = slave.connector.GetConfig()

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

	fmt.Println(slave.config.Slave)

	// set table
	slave.table = slave.config.Master.Table

	// set schema TODO по идее в конфиге модели должен быть указан название БД но это не точно!
	slave.schema = helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB).DBname

	log.Fatal("debug stop")
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
