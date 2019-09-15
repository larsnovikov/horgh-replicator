package commands

import (
	"fmt"
	"github.com/pingcap/errors"
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh2-replicator/app"
	"horgh2-replicator/app/constants"
)

var cmdListen = cobra.Command{
	Use:   "listen",
	Short: "Listen master",
	Long:  "Listen master",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		container, err := app.New()
		if err != nil {
			panic(errors.Wrap(err, constants.ErrorMakeContainer))
		}
		log.Infof(constants.StartMasterListen)
		fmt.Println(container)
	},
}

func Listen() *cobra.Command {
	return &cmdListen
}
