package core

import (
	"os"
	"os/signal"
	"syscall"
)

func SysCallNotify(stopChan chan struct{}) {
	syscallreceived := make(chan os.Signal, 1)
	signal.Notify(syscallreceived, syscall.SIGINT, syscall.SIGTERM)
	// SIGINT is generated via manual control-c from CLI and SIGTERM is default signal of the kill.
	// log.Println("No signals received,waiting for INT and TERM signals..") Commented because we don't need a log of don't receiving a signal
	// MUST RUN VIA BINARY EXECUTE BECAUSE IT WILL DON'T GET THE TERMINATE SIGNALS IF IT'S A SUB PROCESS.
	// e.g IF USED VIA GO RUN,GO RUN WILL BE THE PARENT PROCESS OF THIS FUNCTION.
	systemcall := <-syscallreceived
	if systemcall != nil {
		stopChan <- struct{}{}
	}
}
