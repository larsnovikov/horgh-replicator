package main

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/slave"
	"go-binlog-replication/src/tools"
)

func main() {
	helpers.MakeCredentials()
	slave.MakeSlavePool()

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(tools.CmdListen, tools.CmdLoad, tools.CmdSetPosition)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf(constants.ErrorCobraStarter, err)
	}
}
