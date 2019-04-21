package main

import (
	"github.com/spf13/cobra"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools"
	"horgh-replicator/src/tools/system"
)

func main() {
	helpers.MakeCredentials()
	slave.MakeSlavePool()
	system.MakeHandler()

	var rootCmd = &cobra.Command{Use: "horgh-replicator"}
	rootCmd.AddCommand(tools.CmdListen, tools.CmdLoad, tools.CmdSetPosition, tools.CmdModelCreator)
	_ = rootCmd.Execute()
}
