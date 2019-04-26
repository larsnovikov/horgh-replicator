package system

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/spf13/cobra"
	helpers2 "horgh-replicator/src/tools/helpers"
	"strconv"
)

var CmdSetPosition = &cobra.Command{
	Use:   "set-position",
	Short: "Set position for slave table. Format: [table, binlog_name, binlog_position]",
	Long:  "Set position for slave table. Format: [table, binlog_name, binlog_position]",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		helpers2.Table = args[0]
		name := args[1]
		pos, _ := strconv.Atoi(args[2])

		helpers2.Position = mysql.Position{
			Name: name,
			Pos:  uint32(pos),
		}

		helpers2.SetPosition()
	},
}
