package slave

import (
	"encoding/json"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/connectors"
	"horgh-replicator/src/connectors/clickhouse"
	"horgh-replicator/src/connectors/mysql"
	"horgh-replicator/src/connectors/postgresql"
	"horgh-replicator/src/connectors/vertica"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/tools/exit"
	"io/ioutil"
	"os"
	"strings"
)

type AbstractConnector interface {
	GetInsert() map[string]interface{}
	GetUpdate() map[string]interface{}
	GetDelete(all bool) map[string]interface{}
	Exec(map[string]interface{}) bool
	GetConfigStruct() interface{}
	SetConfig(interface{})
	SetParams(map[string]interface{})
	ParseKey([]interface{})
	GetFields() map[string]connectors.ConfigField
	GetTable() string
	Connection() helpers.Storage
	ParseConfig()
}

type Slave struct {
	connector AbstractConnector
	config    Config
	key       string
	table     string
	channel   chan func() bool
}

type Config struct {
	Master ConfigMaster `json:"master"`
	Slave  interface{}  `json:"slave"`
}

type ConfigMaster struct {
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
}

type Header struct {
	Timestamp uint32
	LogPos    uint32
}

var slavePool map[string]Slave

func getModel() AbstractConnector {

	switch os.Getenv("SLAVE_TYPE") {
	case "mysql":
		return &mysql.Model{}
	case "clickhouse":
		return &clickhouse.Model{}
	case "postgresql":
		return &postgresql.Model{}
	case "vertica":
		return &vertica.Model{}
	}

	return &mysql.Model{}
}

func GetSlaveByName(name string) Slave {
	if slave, ok := slavePool[name]; ok {
		return slave
	}

	exit.Fatal(constants.ErrorUndefinedSlave)

	return Slave{}
}

func MakeSlavePool() {
	slavePool = make(map[string]Slave)
	for _, tableName := range helpers.GetTables() {
		table := strings.TrimSpace(tableName)
		makeSlave(table)
	}
}

// make model, read config by modelName, set var model
func makeSlave(modelName string) {
	slave := Slave{}

	slave.connector = getModel()

	// parse .env config
	slave.GetConnector().ParseConfig()

	// add connector config to base config
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

	// make channel
	slave.channel = make(chan func() bool, helpers.GetChannelSize())
	go save(slave.channel)

	slavePool[modelName] = slave
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

func (slave Slave) GetChannelLen() int {
	return len(slave.channel)
}

func (slave Slave) Insert(header *Header, positionSet func()) {
	if slave.BeforeSave() == true {
		params := slave.connector.GetInsert()

		slave.channel <- func() bool {
			if slave.connector.Exec(params) {
				log.Infof(constants.MessageInserted, header.Timestamp, slave.TableName(), header.LogPos)
				positionSet()
				return true
			}

			slave.logError("insert")

			return false
		}
	}
}

func (slave Slave) Update(header *Header, positionSet func()) {
	if slave.BeforeSave() == true {
		params := slave.connector.GetUpdate()

		slave.channel <- func() bool {
			// fmt.Println(params["params"])
			if slave.connector.Exec(params) {
				log.Infof(constants.MessageUpdated, header.Timestamp, slave.TableName(), header.LogPos)
				positionSet()
				return true
			}

			slave.logError("update")

			return false
		}
	}
}

func (slave Slave) Delete(header *Header, positionSet func()) {
	params := slave.connector.GetDelete(false)

	slave.channel <- func() bool {
		if slave.connector.Exec(params) {
			log.Infof(constants.MessageDeleted, header.Timestamp, slave.TableName(), header.LogPos)
			positionSet()
			return true
		}

		slave.logError("delete")

		return false
	}
}

func (slave Slave) DeleteAll(header *Header, positionSet func()) {
	params := slave.connector.GetDelete(true)

	slave.channel <- func() bool {
		if slave.connector.Exec(params) {
			log.Infof(constants.MessageDeletedAll, header.Timestamp, slave.TableName())
			positionSet()
			return true
		}

		slave.logError("delete")

		return false
	}
}

func (slave Slave) logError(operationType string) {
	exit.Fatal(constants.ErrorSave, operationType, slave.TableName())
}
