package helpers

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/constants"
	"os"
	"reflect"
)

func MakeSlice(input interface{}) []interface{} {
	s := reflect.ValueOf(input)
	if s.Kind() != reflect.Slice {
		log.Fatal(constants.ErrorSliceCreation)
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func MakeTablePosKey(hash string) (pos string, name string) {
	pos = fmt.Sprintf(constants.PositionPosPrefix, hash)
	name = fmt.Sprintf(constants.PositionNamePrefix, hash)

	return pos, name
}

func MakeHash(dbName string, table string) string {
	return dbName + "." + table
}

func ReadConfig(configName string) *os.File {
	fileName := fmt.Sprintf(constants.ConfigPath, configName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Fatalf(constants.ErrorNoModelFile, fileName)
	}

	jsonFile, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	return jsonFile
}
