package core

import (
	"hash/crc32"
	"log"
)

func Getsum(log_without_checksum string) uint32 {
	MIN_EXPECTED_LINES := 1 // Can be dynamicly learned in future.
	if len(log_without_checksum) < MIN_EXPECTED_LINES {
		return 1
	} // Only calculates checksum if log is valid
	checksumoflog := crc32.MakeTable(0x82f63b78) // checksum algoritms can changed via here.Only changing hexadecimal part is enough.
	/* Castagnoli's polynomial, used in iSCSI.
	 Has better error detection characteristics than IEEE.
	 https://dx.doi.org/10.1109/26.231911
	source: https://pkg.go.dev/hash/crc32#example-MakeTable */
	checksum := crc32.Checksum([]byte(log_without_checksum), checksumoflog)
	log.Printf("Checksum of %s : %08x\n", log_without_checksum, checksum)
	return checksum

}
