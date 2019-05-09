package master

import (
	"horgh-replicator/src/models/slave"
	toolsHelper "horgh-replicator/src/tools/helpers"
)

const (
	ParseStringSize = 99999999
)

func buildModel(tableName string) {
	toolsHelper.Table = tableName
	if canHandle() == true {
		// create table for log if not exists
		// pgsql create method if not exists
		// pgsql create trigger if not exists
		// truncate log table
		// start pgsql dump and listen
		//
		toolsHelper.ParseStrings = make(chan string, ParseStringSize)
		go parseLine(toolsHelper.ParseStrings)

		readDump()

		toolsHelper.Wait(func() bool {
			return slave.GetSlaveByName(toolsHelper.Table).GetChannelLen() == 0 && len(toolsHelper.ParseStrings) == 0
		})
	}
}

func parseLine(c chan string) {
	//for {
	//line := <-c

	// try to parse like insert
	//if parseInsert(line) == true {
	//	continue
	//}
	//
	//// try to parse like position setter
	//if parsePosition(line) == true {
	//	continue
	//}
	//}
}

func canHandle() bool {
	return true
	//savedPos := GetSavedPos(toolsHelper.Table)
	//if savedPos.Name == "" && savedPos.Pos == 0 {
	//	return true
	//}
	//
	//exit.Fatal(constants.ErrorSlaveBuilt, toolsHelper.Table, toolsHelper.Table)
	//return false
}

func readDump() {

}
