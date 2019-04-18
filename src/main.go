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

	var rootCmd = &cobra.Command{Use: "horgh-replicator"}
	rootCmd.AddCommand(tools.CmdListen, tools.CmdLoad, tools.CmdSetPosition, tools.CmdModelCreator)
	_ = rootCmd.Execute()
	//if err != nil {
	//	log.Fatalf(constants.ErrorCobraStarter, err)
	//}
}
