package helpers

import (
	"go-binlog-replication/src/constants"
	"log"
	"reflect"
)

func makeSlice(input interface{}) []interface{} {
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
