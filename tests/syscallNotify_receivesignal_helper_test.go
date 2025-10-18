package tests

import (
	"os"

	"github.com/showbgpsummary/minimalist_nos_mvp/core"
)

// ONLY USED BY syscallNotify_receivesignal_helper_test.go TO HAVE ANOTHER PID.
// SHOULDN'T BE A TEST ON IT'S OWN.
func HelperSignalHandler() {
	stopChan := make(chan struct{})
	go core.SysCallNotify(stopChan)
	<-stopChan
	os.Exit(0)
}
