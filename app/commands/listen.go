package commands

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh2-replicator/app/constants"
)

var cmdListen = cobra.Command{
	Use:   "listen",
	Short: "Listen master",
	Long:  "Listen master",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof(constants.StartMasterListen)
	},
}

func Listen() *cobra.Command {
	return &cmdListen
}
