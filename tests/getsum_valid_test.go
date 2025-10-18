package tests

import (
	"testing"

	main "github.com/showbgpsummary/telemetry-ebpf/core"
)

func TestValidInput(t *testing.T) {
	input := "checksum123"
	expected := uint32(0x8FA66C63)
	result := main.Getsum(input)
	if result != expected {
		t.Fatalf("checksum does not match the expected value, expected %d,Got %d", expected, result)
	}
}
