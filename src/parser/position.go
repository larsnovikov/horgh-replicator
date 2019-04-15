package parser

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/system"
	"strconv"
	"strings"
)

var curPosition mysql.Position

func makeHash(dbName string, table string) string {
	return dbName + "." + table
}

func getMinPosition(position mysql.Position) mysql.Position {
	tmpLogSuffix, err := strconv.Atoi(strings.Replace(position.Name, constants.MasterLogNamePrefix, "", -1))
	if err != nil {
		log.Fatalf(constants.ErrorGetMinPosition, err)
	}
	fmt.Println(tmpLogSuffix)
	// build current position
	if curPosition.Pos == 0 {
		dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname

		// get all saved positions for operated tables and fin with min pos
		// WARNING! if it is first start for table, replicate it from min pos of another tables
		for _, table := range helpers.GetTables() {
			hash := makeHash(dbName, table)
			pos, name := helpers.MakeTablePosKey(hash)

			tablePosition, err := strconv.ParseUint(system.GetValue(pos), 10, 32)
			if err != nil {
				log.Fatalf(constants.ErrorGetMinPosition, err)
			}
			tableLogFile := system.GetValue(name)

			tableLogSuffix, err := strconv.Atoi(strings.Replace(position.Name, constants.MasterLogNamePrefix, "", -1))
			if err != nil {
				log.Fatalf(constants.ErrorGetMinPosition, err)
			}
			// if log file from storage lower than log file in master - set position from storage
			if tableLogSuffix < tmpLogSuffix {
				position.Pos = uint32(tablePosition)
				position.Name = tableLogFile
			} else {
				// if log file from storage is greater or equal log file from master - check position
				if uint32(tablePosition) < curPosition.Pos || curPosition.Pos == 0 {
					position.Pos = uint32(tablePosition)
					position.Name = tableLogFile
				}
			}
		}
		curPosition = position
	}

	return curPosition
}

// set position for table
func SetPosition(table string, pos mysql.Position) {
	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := makeHash(dbName, table)

	posKey, nameKey := helpers.MakeTablePosKey(hash)

	system.SetValue(posKey, fmt.Sprint(pos.Pos))
	system.SetValue(nameKey, pos.Name)
}
