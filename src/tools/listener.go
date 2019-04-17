package tools

import (
	"github.com/spf13/cobra"
	"go-binlog-replication/src/parser"
)

var CmdListen = &cobra.Command{
	Use:   "listen",
	Short: "Listen master binlog",
	Long:  `Listen master binlog`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		parser.BinlogListener()
	},
}
