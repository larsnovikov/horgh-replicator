package main

import (
	"github.com/spf13/cobra"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools"
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
