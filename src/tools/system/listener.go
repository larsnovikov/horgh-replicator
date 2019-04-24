package system

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/parser"
	"horgh-replicator/src/tools"
	"strconv"
	"time"
)

var CmdListen = &cobra.Command{
	Use:   "listen",
	Short: "Listen master binlog",
	Long:  `Listen master binlog`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tools.BeforeExit = func() {
			// stop handle binlog
			log.Infof(constants.MessageStopHandlingBinlog)
			parser.AllowHandling = false

			// stop handle saving
			log.Infof(constants.MessageStopHandlingSave)
			slave.AllowHandling = false

			log.Infof(constants.MessageWait, strconv.Itoa(10), "seconds")
			time.Sleep(10 * time.Second)
		}
		parser.BinlogListener()
	},
}
