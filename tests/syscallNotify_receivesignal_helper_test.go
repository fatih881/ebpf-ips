package tests

import (
	"os"

	"github.com/fatih881/ebpf-ips/core"
)

// ONLY USED BY syscallNotify_receivesignal_test.go TO HAVE ANOTHER PID.
// SHOULDN'T BE A TEST ON IT'S OWN.
func HelperSignalHandler() {
	stopChan := make(chan struct{})
	go core.SysCallNotify(stopChan)
	<-stopChan
	os.Exit(0)
}
