package tests

import (
	"net"
	"reflect"
	"testing"

	"github.com/fatih881/ebpf-ips/core"
)

func TestFilterInterfaces(t *testing.T) {
	testCases := []struct {
		name            string
		inputInterfaces []net.Interface
		expectedOutput  []string
	}{
		{
			name: "1 expected",
			inputInterfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp | net.FlagBroadcast},
				{Name: "lo", Flags: net.FlagUp | net.FlagLoopback},
				{Name: "docker0", Flags: net.FlagBroadcast},
			},
			expectedOutput: []string{"eth0"},
		},
		{
			name: "0 expected",
			inputInterfaces: []net.Interface{
				{Name: "lo", Flags: net.FlagUp | net.FlagLoopback},
				{Name: "eth1", Flags: net.FlagMulticast},
			},
			expectedOutput: []string{},
		},
		{
			name:            "empty",
			inputInterfaces: []net.Interface{},
			expectedOutput:  []string{},
		},
		{
			name: "2 expected",
			inputInterfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
				{Name: "wlan0", Flags: net.FlagUp | net.FlagMulticast},
				{Name: "lo", Flags: net.FlagUp | net.FlagLoopback},
			},
			expectedOutput: []string{"eth0", "wlan0"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, err := core.FilterInterfaces(tc.inputInterfaces)
			if err != nil {
				t.Fatalf("warning: filterinterfaces returned error : %v", err)
			}

			if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
				t.Errorf("expected: %v, output: %v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
