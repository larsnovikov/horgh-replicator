package system

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/master"
	"horgh-replicator/src/tools/exit"
)

var CmdBuildTable = &cobra.Command{
	Use:   "build-slave",
	Short: "Build slave table from master",
	Long:  "Build slave table from master",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		beforeExit := func() bool {
			log.Infof(constants.MessageStopTableBuild)
			return false
		}
		exit.BeforeExitPool = append(exit.BeforeExitPool, beforeExit)

		master.GetModel().BuildSlave(helpers.GetTable())
	},
}
