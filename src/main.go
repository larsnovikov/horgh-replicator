package main

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/slave"
	"go-binlog-replication/src/parser"
	"go-binlog-replication/src/tools"
)

func main() {
	helpers.MakeCredentials()
	slave.MakeSlavePool()

	var cmdListen = &cobra.Command{
		Use:   "listen",
		Short: "Listen master binlog",
		Long:  `Listen master binlog`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			parser.BinlogListener()
		},
	}

	var cmdLoad = &cobra.Command{
		Use:   "load",
		Short: "Create queries to master",
		Long:  `Create queries to master`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			tools.Load()
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdListen, cmdLoad)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf(constants.ErrorCobraStarter, err)
	}
}
