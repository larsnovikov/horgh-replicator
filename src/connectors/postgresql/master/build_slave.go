package master

import (
	"fmt"
	slave2 "horgh-replicator/src/connectors/postgresql/slave"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/master"
	"horgh-replicator/src/models/slave"
	toolsHelper "horgh-replicator/src/tools/helpers"
)

func buildModel(tableName string) {
	toolsHelper.Table = tableName
	if canHandle() == true {
		dbType := slave2.Type
		var query string
		// create table for log if not exists
		table := fmt.Sprintf(logTableName, tableName)

		// create table for log if not exists
		query = helpers.GetQuery(dbType, "table", table, table, table, table)
		master.Exec(helpers.Query{
			Query:  query,
			Params: []interface{}{},
		})

		// create method if not exists
		query = helpers.GetQuery(dbType, "func", table)
		master.Exec(helpers.Query{
			Query:  query,
			Params: []interface{}{},
		})

		// create trigger if not exists
		query = helpers.GetQuery(dbType, "trigger", table, table)
		master.Exec(helpers.Query{
			Query:  query,
			Params: []interface{}{},
		})

		// truncate log table
		// start pgsql dump and listen
		//
		toolsHelper.ParseStrings = make(chan string, parseStringSize)
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
