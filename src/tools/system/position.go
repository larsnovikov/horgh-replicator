package system

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/spf13/cobra"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/system"
	"strconv"
)

var CmdSetPosition = &cobra.Command{
	Use:   "set-position",
	Short: "Set position for slave table. Format: [table, name, position]",
	Long:  "Set position for slave table. Format: [table, name, position]",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		name := args[1]
		pos, _ := strconv.Atoi(args[2])

		setPosition(tableName, name, pos)
	},
}

func setPosition(table string, name string, pos int) {
	position := mysql.Position{
		Name: name,
		Pos:  uint32(pos),
	}

	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := helpers.MakeHash(dbName, table)

	posKey, nameKey := helpers.MakeTablePosKey(hash)

	system.SetValue(posKey, fmt.Sprint(position.Pos))
	system.SetValue(nameKey, position.Name)
}
