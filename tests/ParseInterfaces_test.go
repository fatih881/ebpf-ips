package tests

import (
	"reflect"
	"testing"

	"github.com/showbgpsummary/telemetry-ebpf/core"
)

func TestParseInterfaces(t *testing.T) {
	expectedoutput := []string{"wlan0", "eth0", "eth1"}
	input := []byte("wlan0\neth0\neth1\n")
	actualoutput := core.ParseInterfaces(input)
	if !reflect.DeepEqual(actualoutput, expectedoutput) {
		t.Errorf("WARNING : error at findinterfaces.go,ParseInterfaces.Expected value : %v\n Actual value := %v", expectedoutput, actualoutput)
	}
}
