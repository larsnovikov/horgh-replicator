package helpers

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/constants"
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
