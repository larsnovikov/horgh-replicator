package tools

import (
	"os"
	"os/signal"
	"syscall"
)

var BeforeExit = func() {}

func MakeHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTSTP, syscall.SIGQUIT)
	go handle(c)
}

func handle(c chan os.Signal) {
	for {
		<-c
		BeforeExit()
		os.Exit(1)
	}
}
