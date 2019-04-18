package tools

import (
	"github.com/spf13/cobra"
)

var CmdModelCreator = &cobra.Command{
	Use:   "create-model",
	Short: "Create model.json for master table. Format: [table]",
	Long:  "Create model.json for master table. Format: [table]",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		makeModel(tableName)
	},
}

func makeModel(table string) {

}

func getStructure() map[string]interface{} {
	return map[string]interface{}{}
}

func buildJsonFile() bool {
	return true
}
