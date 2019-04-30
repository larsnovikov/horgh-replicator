package system

import (
	"github.com/spf13/cobra"
	"horgh-replicator/src/parser"
	"horgh-replicator/src/tools"
)

var CmdListen = &cobra.Command{
	Use:   "listen",
	Short: "Listen master binlog",
	Long:  `Listen master binlog`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tools.BeforeExit = func() bool {
			parser.StopListen()

			return true
		}
		parser.BinlogListener()
	},
}
