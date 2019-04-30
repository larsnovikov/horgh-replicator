package system

import (
	"github.com/spf13/cobra"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/parser"
	"horgh-replicator/src/tools/exit"
)

var CmdListen = &cobra.Command{
	Use:   "listen",
	Short: "Listen master binlog",
	Long:  `Listen master binlog`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		exit.BeforeExitPool = append(exit.BeforeExitPool, parser.Stop)
		exit.BeforeExitPool = append(exit.BeforeExitPool, slave.Stop)
		parser.BinlogListener()
	},
}
