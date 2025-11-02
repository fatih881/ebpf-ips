package tests

import (
	"os"
	"testing"

	main "github.com/fatih881/ebpf-ips/core"
)

func TestEmptystring(t *testing.T) {
	var testVariable string
	result := main.Getsum(testVariable)
	if result != 1 {
		t.Fatalf("Function can't handle empty strings properly")
		os.Exit(1)
	} else {
		return
	}
}
