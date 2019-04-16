package clickhouse

import (
	"fmt"
	"go-binlog-replication/src/connectors"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"strings"
)

const (
	Type   = "clickhouse"
	Insert = `INSERT INTO %s.%s(%s) VALUES(%s);`
	Update = `ALTER TABLE %s.%s UPDATE %s WHERE %s=?;`
	Delete = `ALTER TABLE %s.%s DELETE WHERE %s=?;`
)

type Model struct {
	table       string
	schema      string
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
	model.schema = helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname

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

func (model *Model) GetInsert() map[string]interface{} {
	var params []interface{}
	var fieldNames []string
	var fieldValues []string

	for _, value := range model.fields {
		fieldNames = append(fieldNames, "`"+value.Name+"`")
		fieldValues = append(fieldValues, "?")

		params = append(params, model.params[value.Name])
	}

	query := fmt.Sprintf(Insert, model.schema, model.table, strings.Join(fieldNames, ","), strings.Join(fieldValues, ","))

	return map[string]interface{}{
		"query":  query,
		"params": params,
	}
}

func (model *Model) GetUpdate() map[string]interface{} {
	var params []interface{}
	var fields []string

	for _, value := range model.fields {
		if value.Name != model.key {
			fields = append(fields, "`"+value.Name+"`"+"=?")
			params = append(params, model.params[value.Name])
		}
	}

	// add key to params
	params = append(params, model.params[model.key])

	query := fmt.Sprintf(Update, model.schema, model.table, strings.Join(fields, ", "), model.key)

	return map[string]interface{}{
		"query":  query,
		"params": params,
	}
}

func (model *Model) GetDelete() map[string]interface{} {
	var params []interface{}
	query := fmt.Sprintf(Delete, model.schema, model.table, model.key)

	params = append(params, model.params[model.key])

	return map[string]interface{}{
		"query":  query,
		"params": params,
	}
}

func (model *Model) Exec(params map[string]interface{}) bool {
	return model.Connection().Exec(params)
}

func (model *Model) Connection() helpers.Storage {
	helpers.ConnPool.Slave = GetConnection(helpers.ConnPool.Slave, constants.DBSlave).(helpers.Storage)
	return helpers.ConnPool.Slave
}
