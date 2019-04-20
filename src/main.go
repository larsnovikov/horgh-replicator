package main

import (
	"github.com/siddontang/go-log/log"
	"github.com/spf13/cobra"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/parser"
	"horgh-replicator/src/tools"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func handle(c chan os.Signal) {
	for {
		<-c
		// stop handle binlog
		log.Infof(constants.MessageStopHandlingBinlog)
		parser.AllowHandling = false

		// stop handle saving
		log.Infof(constants.MessageStopHandlingSave)
		slave.AllowHandling = false

		log.Infof(constants.MessageWait, strconv.Itoa(10), "seconds")
		time.Sleep(10 * time.Second)
		// exit
		os.Exit(1)
	}
}

func main() {
	helpers.MakeCredentials()
	slave.MakeSlavePool()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTSTP, syscall.SIGQUIT)
	go handle(c)

	var rootCmd = &cobra.Command{Use: "horgh-replicator"}
	rootCmd.AddCommand(tools.CmdListen, tools.CmdLoad, tools.CmdSetPosition, tools.CmdModelCreator)
	_ = rootCmd.Execute()
}
