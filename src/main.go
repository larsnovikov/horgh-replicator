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
	tools.MakeHandler()

	var rootCmd = &cobra.Command{Use: "horgh-replicator"}
	rootCmd.AddCommand(
		system.CmdListen,
		system.CmdLoad,
		system.CmdSetPosition,
		system.CmdModelCreator,
		system.CmdBuildTable,
		system.CmdDestroyTable,
	)
	_ = rootCmd.Execute()
}
