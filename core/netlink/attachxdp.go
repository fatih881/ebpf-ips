package netlink

import (
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	ebpf_export "github.com/fatih881/ebpf-ips/core/ebpf"
)

// The sole purpose of this function is to attach the interface data it receives to XDP.
// It is the orchestrator's responsibility to decide which function to connect to and which not.
func Attach(interfaces []string, objs *ebpf_export.IpsObjects) (map[int]link.Link, error) {
	for _, ifaceName := range interfaces {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			log.Printf("warning: cannot get interfacename for %s: %v", ifaceName, err)
			continue
		}
		XDPLink, err := link.AttachXDP(link.XDPOptions{
			Program:   objs.Firewall,
			Interface: iface.Index,
			Flags:     link.XDPGenericMode,
		})
		if err != nil {
			log.Printf("warning: cannot attach XDP : %v", err)
			continue
		}
		XDPLinks = append(XDPLinks, XDPLink)
		log.Printf("info: attached XDP to %s", ifaceName)
		ActiveLinks[iface.Index] = XDPLink

	}
	return ActiveLinks, nil
}
