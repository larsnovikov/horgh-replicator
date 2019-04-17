package main

import (
	"github.com/spf13/cobra"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/slave"
	"go-binlog-replication/src/tools"
)

func main() {
	helpers.MakeCredentials()
	slave.MakeSlavePool()

	var rootCmd = &cobra.Command{Use: "go-bin-log-replication"}
	rootCmd.AddCommand(tools.CmdListen, tools.CmdLoad, tools.CmdSetPosition)
	_ = rootCmd.Execute()
	//if err != nil {
	//	log.Fatalf(constants.ErrorCobraStarter, err)
	//}
}
