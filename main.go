package main

import (
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
	"horgh2-replicator/app"
	"horgh2-replicator/app/constants"

	"horgh2-replicator/app/commands"
)

func main() {
	err := app.Make()
	if err != nil {
		panic(errors.Wrap(err, constants.ErrorMakeContainer))
	}
	var rootCmd = &cobra.Command{Use: constants.AppName}
	rootCmd.AddCommand(
		commands.Listen(),
	)
	err = rootCmd.Execute()
	if err != nil {
		panic(errors.Wrap(err, constants.ErrorHandleCommand))
	}
}
