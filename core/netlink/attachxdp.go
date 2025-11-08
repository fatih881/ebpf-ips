package netlink

import (
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	ebpf_export "github.com/fatih881/ebpf-ips/core/ebpf"
)

// The sole purpose of this function is to attach the interface data it receives to XDP.
// It is the orchestrator's responsibility to decide which function to connect to and which not.
// XDP manager must decide which flag to choose before calling this function.It can be used for fallback if it's a hardware,or direclty used with veth's for resource efficiency.
func Attach(ifaceName string, flag link.XDPAttachFlags, objs *ebpf_export.IpsObjects) (link.Link, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Printf("warning: cannot get interface id for %s: %v", ifaceName, err)
	}
	opts := link.XDPOptions{
		Program:   objs.Firewall,
		Interface: iface.Index,
		Flags:     flag,
	}
	return link.AttachXDP(opts)
}
