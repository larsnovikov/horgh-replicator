package postgresql

import (
	"fmt"
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strconv"
	"strings"
)

const (
	Type   = "postgresql"
	Insert = "INSERT INTO \"%s\"(%s) VALUES(%s);"
	Update = "UPDATE \"%s\" SET %s WHERE %s=%s;"
	Delete = "DELETE FROM \"%s\" WHERE %s=$1"
)

type Model struct {
	table       string
	key         string
	keyPosition int
	fields      map[string]connectors.ConfigField
	params      map[string]interface{}
}

func (model *Model) ParseConfig() {
	helpers.ParseDBConfig()
}

func (model Model) GetFields() map[string]connectors.ConfigField {
	return model.fields
}

func (model *Model) ParseKey(row []interface{}) {
	// TODO в зависимости от типа поля ключа, тут могут быть разные приведения типов
	// а если будет составной ключ вообще будет тяжко
	model.params[model.key] = int(row[model.keyPosition].(int32))
}

func (model *Model) GetConfigStruct() interface{} {
	return &connectors.ConfigSlave{}
}

func (model *Model) GetTable() string {
	return model.table
}

func (model *Model) SetConfig(config interface{}) {
	model.table = config.(*connectors.ConfigSlave).Table

	model.fields = make(map[string]connectors.ConfigField)
	for pos, value := range config.(*connectors.ConfigSlave).Fields {
		if model.key == "" && value.Key == true {
			model.key = value.Name
			model.keyPosition = pos
		}

		model.fields[value.Name] = value
	}
}

func (model *Model) SetParams(params map[string]interface{}) {
	model.params = params
}

func (model *Model) Insert() bool {
	var params []interface{}
	var fieldNames []string
	var fieldValues []string

	i := 0
	for _, value := range model.fields {
		i++
		fieldNames = append(fieldNames, value.Name)
		fieldValues = append(fieldValues, "$"+strconv.Itoa(i))

		params = append(params, model.params[value.Name])
	}

	query := fmt.Sprintf(Insert, model.table, strings.Join(fieldNames, ","), strings.Join(fieldValues, ","))

	return model.Connection().Exec(map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Update() bool {
	var params []interface{}
	var fields []string

	i := 0
	for _, value := range model.fields {
		i++
		fields = append(fields, value.Name+"=$"+strconv.Itoa(i))

		params = append(params, model.params[value.Name])
	}

	// add key to params
	params = append(params, model.params[model.key])

	i++
	query := fmt.Sprintf(Update, model.table, strings.Join(fields, ","), model.key, "$"+strconv.Itoa(i))

	return model.Connection().Exec(map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Delete() bool {
	var params []interface{}
	query := fmt.Sprintf(Delete, model.table, model.key)

	params = append(params, model.params[model.key])

	return model.Connection().Exec(map[string]interface{}{
		"query":  query,
		"params": params,
	})
}

func (model *Model) Connection() helpers.Storage {
	helpers.ConnPool.Slave = GetConnection(helpers.ConnPool.Slave, constants.DBSlave).(helpers.Storage)
	return helpers.ConnPool.Slave
}
