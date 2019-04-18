package tools

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/spf13/cobra"
	"horgh-replicator/src/parser"
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

		SetPosition(tableName, name, pos)
	},
}

func SetPosition(table string, name string, pos int) {
	position := mysql.Position{
		Name: name,
		Pos:  uint32(pos),
	}

	parser.SetPosition(table, position)
}
