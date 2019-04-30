package system

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools/exit"
	"horgh-replicator/src/tools/helpers"
)

var CmdDestroyTable = &cobra.Command{
	Use:   "destroy-slave",
	Short: "Destroy slave table from master. Format: [table]",
	Long:  "Destroy slave table from master. Format: [table]",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		beforeExit := func() bool {
			log.Infof(constants.MessageStopTableDestroy)
			return false
		}
		exit.BeforeExitPool = append(exit.BeforeExitPool, beforeExit)

		tableName := args[0]
		destroyModel(tableName)
	},
}

func destroyModel(tableName string) {
	helpers.Table = tableName
	header, positionSet := helpers.GetHeader()

	// delete all from table
	slave.GetSlaveByName(helpers.Table).DeleteAll(&header, positionSet)

	// delete position in db
	helpers.SetPosition()

	helpers.Wait(func() bool {
		return slave.GetSlaveByName(helpers.Table).GetChannelLen() == 0
	})

	log.Infof(constants.MessageTableDestroyed, helpers.Table)
}
