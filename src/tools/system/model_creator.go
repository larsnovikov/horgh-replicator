package system

import (
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh-replicator/src/connectors"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/master"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools/exit"
	"io/ioutil"
	"os"
	"regexp"
)

type FieldDefinition struct {
	Field   string
	Type    string
	Null    []byte
	Key     []byte
	Default []byte
	Extra   string
}

var CmdModelCreator = &cobra.Command{
	Use:   "create-model",
	Short: "Create model.json for master table. Format: [table]",
	Long:  "Create model.json for master table. Format: [table]",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		makeModel(tableName)
	},
}

func makeModel(table string) {
	fields := getStructure(table)
	data := getJson(table, fields)
	fileName := build(table, data)

	log.Infof(constants.MessageConfigCreated, table, fileName)
}

func getStructure(table string) []FieldDefinition {
	query := fmt.Sprintf("DESCRIBE %s", table)
	rows := master.Get(helpers.Query{
		Query:  query,
		Params: []interface{}{},
	})
	var fields []FieldDefinition
	for rows.Next() {
		field := FieldDefinition{}
		err := rows.Scan(&field.Field, &field.Type, &field.Null, &field.Key, &field.Default, &field.Extra)
		if err != nil {
			exit.Fatal(constants.ErrorTableStructure, table, err)
		}

		fields = append(fields, field)
	}

	err := rows.Err()
	if err != nil {
		exit.Fatal(constants.ErrorTableStructure, table, err)
	}

	return fields
}

func getJson(table string, fields []FieldDefinition) string {
	// field names for master config section
	var fieldNames []string
	// fields in config format
	var preparedFields []connectors.ConfigField

	// foreach fields add create field slices
	for _, el := range fields {
		fieldNames = append(fieldNames, el.Field)

		primaryKey := false

		if fmt.Sprintf("%s", el.Key) == "PRI" {
			primaryKey = true
		}
		tmpField := connectors.ConfigField{
			Name: el.Field,
			Key:  primaryKey,
			Mode: getType(el),
		}
		preparedFields = append(preparedFields, tmpField)
	}

	masterConfig := slave.ConfigMaster{
		Table:  table,
		Fields: fieldNames,
	}

	slaveConfig := connectors.ConfigSlave{
		Table:  table,
		Fields: preparedFields,
	}

	config := slave.Config{
		Master: masterConfig,
		Slave:  slaveConfig,
	}

	jsonOut, err := json.MarshalIndent(&config, "", "\t")
	if err != nil {
		exit.Fatal(constants.ErrorBuildModelConfig, err)
	}

	return string(jsonOut)
}

func getType(definition FieldDefinition) string {
	re := regexp.MustCompile(`(\([0-9]+\))`)
	defType := re.ReplaceAllString(definition.Type, "")

	switch defType {
	case "varchar":
		return "string"
	case "text":
		return "string"
	case "tinyint":
		return "bool"
	case "int":
		return "int"
	case "timestamp":
		return "timestamp"
	case "float":
		return "float"
	case "real":
		return "float"
	case "decimal":
		return "float"
	case "time":
		return "time"
	case "date":
		return "date"
	case "datetime":
		return "datetime"
	default:
		exit.Fatal(constants.ErrorFieldTypeConversion, defType)
	}

	return ""
}

func build(table string, data string) string {
	fileName := fmt.Sprintf(constants.ConfigPath, table)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := ioutil.WriteFile(fileName, []byte(data), 0644)
		if err != nil {
			exit.Fatal(constants.ErrorBuildModelConfig, err)
		}
	} else {
		exit.Fatal(constants.ErrorModelFileExists, fileName)
	}

	return fileName
}
