package netlink

import (
	"fmt"
	"log"

	"github.com/cilium/ebpf/link"
	ebpf_export "github.com/fatih881/ebpf-ips/core/ebpf"
)

// XDPLinks is a list of all XDP link objects.
// ActiveLinks is a map where the key is the interface index (ID) and the value is the XDP link object.
var (
	XDPLinks    = make([]link.Link, 0)
	ActiveLinks = make(map[int]link.Link)
)

// Attachmanager function is the API between single purpose functions and main.go.
// It attaches the XDP program to all the interfaces but loopback and down interfaces,this logic is implemented in findinterfaces.go
// Returns the Active map link for potential future use(monitoring active interfaces,cleanup etc.)
func AttachManager(objs *ebpf_export.IpsObjects) (map[int]link.Link, error) {
	interfaces, err := FindInterfaces()
	if err != nil {
		return nil, fmt.Errorf("warning: cannot find interfaces from host : %w", err)
	}
	ActiveLinks, err = Attach(interfaces, objs)
	if err != nil {
		return nil, fmt.Errorf("warning: cannot attach interfaces to host : %w", err)
	}
	log.Printf("info: successfully attached XDP to %d interfaces ", len(ActiveLinks))
	return ActiveLinks, nil
}
