package tools

import (
	"horgh-replicator/src/constants"
	"horgh-replicator/src/tools/exit"
	"os"
	"os/signal"
	"syscall"
)

func MakeHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTSTP, syscall.SIGQUIT)
	go handle(c)
}

func handle(c chan os.Signal) {
	for {
		<-c
		exit.Fatal(constants.MessageSysCallStop)
	}
}
