package core

import (
	"log"
	"os/exec"
	"strings"
)

func FindInterfaces() ([]string, error) {
	cmd := exec.Command("sh", "-c", "ls /sys/class/net | grep -v lo")
	// sys/class/net haves all the interfaces as a folder."grep -v --invert-match  select non-matching lines"(From grep --help output)
	interfacenames, err := cmd.Output()
	if err != nil {
		log.Println("Can't find interfaces", err)
		return nil, err
	}
	return ParseInterfaces(interfacenames), nil

}
func ParseInterfaces(data []byte) []string {
	cleanlines := strings.Split(strings.TrimSpace(string(data)), "\n")
	return cleanlines

}
