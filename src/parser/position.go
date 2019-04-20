package parser

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/mysql"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/system"
	"strconv"
	"strings"
)

var curPosition mysql.Position
var PrevPosition map[string]mysql.Position
var SaveLocks map[string]bool

var channel chan func()

func makeHash(dbName string, table string) string {
	return dbName + "." + table
}

func updatePrevPosition(c chan func()) {
	for {
		method := <-c
		method()
	}
}

func GetSavedPos(table string) mysql.Position {
	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := makeHash(dbName, table)
	pos, name := helpers.MakeTablePosKey(hash)

	tablePosition, err := strconv.ParseUint(system.GetValue(pos), 10, 32)
	if err != nil {
		log.Fatalf(constants.ErrorGetMinPosition, err)
	}
	tableLogFile := system.GetValue(name)

	return mysql.Position{
		Name: tableLogFile,
		Pos:  uint32(tablePosition),
	}
}

func getMinPosition(position mysql.Position) mysql.Position {
	tmpLogSuffix, err := strconv.Atoi(strings.Replace(position.Name, helpers.GetMasterLogFilePrefix(), "", -1))
	if err != nil {
		log.Fatalf(constants.ErrorGetMinPosition, err)
	}

	// build current position
	if curPosition.Pos == 0 {
		// get all saved positions for operated tables and fin with min pos
		// WARNING! if it is first start for table, replicate it from min pos of another tables
		for _, table := range helpers.GetTables() {
			savedPos := GetSavedPos(table)
			tablePosition := savedPos.Pos
			tableLogFile := savedPos.Name

			// TODO переписать на GetlowPosition()
			tableLogSuffix := GetLogFileSuffix(position.Name)

			// if log file from storage lower than log file in master - set position from storage
			if tableLogSuffix < tmpLogSuffix {
				position.Pos = uint32(tablePosition)
				position.Name = tableLogFile
			} else {
				// if log file from storage is greater or equal log file from master - check position
				if uint32(tablePosition) < position.Pos || position.Pos == 0 {
					position.Pos = uint32(tablePosition)
					position.Name = tableLogFile
				}
			}
		}
		curPosition = position

		PrevPosition = make(map[string]mysql.Position)
		SaveLocks = make(map[string]bool)

		channel = make(chan func())
		go updatePrevPosition(channel)

		for _, table := range helpers.GetTables() {
			PrevPosition[table] = curPosition
			SaveLocks[table] = true
		}
	}

	return curPosition
}

// set position for table
func SetPosition(table string, pos mysql.Position) {
	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := makeHash(dbName, table)

	posKey, nameKey := helpers.MakeTablePosKey(hash)

	channel <- func() {
		system.SetValue(posKey, fmt.Sprint(pos.Pos))
		system.SetValue(nameKey, pos.Name)
		PrevPosition[table] = pos
	}
}

func GetLowPosition(pos1 mysql.Position, pos2 mysql.Position) mysql.Position {
	position := mysql.Position{}
	pos1Suffix := GetLogFileSuffix(pos1.Name)
	pos2Suffix := GetLogFileSuffix(pos2.Name)

	// if log file from storage lower than log file in master - set position from storage
	if pos1Suffix < pos2Suffix {
		position.Pos = uint32(pos1.Pos)
		position.Name = pos1.Name
	} else {
		// if log file from storage is greater or equal log file from master - check position
		if uint32(pos1.Pos) < pos2.Pos {
			position.Pos = uint32(pos1.Pos)
			position.Name = pos1.Name
		} else {
			position.Pos = uint32(pos2.Pos)
			position.Name = pos2.Name
		}
	}

	return position
}

func GetLogFileSuffix(name string) int {
	suff, err := strconv.Atoi(strings.Replace(name, helpers.GetMasterLogFilePrefix(), "", -1))
	if err != nil {
		log.Fatalf(constants.ErrorGetMinPosition, err)
	}

	return suff
}
