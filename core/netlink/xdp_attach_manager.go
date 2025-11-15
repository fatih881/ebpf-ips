package netlink

import (
	"fmt"
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
)

// WriteChan ReadChan StopChan must be called and attached before StartLinkManager.
var WriteChan chan struct {
	NewLink NewLink
	Err     error
}
var ReadChan chan chan map[int]link.Link
var StopChan chan struct{}

type NewLink struct {
	Flag       link.XDPAttachFlags
	LinkIndex  int
	LinkObject link.Link
}

// StartLinkManager must be started before AttachManager function.
func StartLinkManager(writeChan chan NewLink, readChan chan chan map[int]link.Link, stopChan chan struct{}) {
	var (
		// activeLinks is a map where the key is the interface index (ID) and the value is the XDP link object.
		activeLinks = make(map[int]link.Link)
	)
	for {
		select {
		case newLink := <-writeChan:
			activeLinks[newLink.LinkIndex] = newLink.LinkObject
		case req := <-readChan:
			copyactiveLinks := make(map[int]link.Link, len(activeLinks))
			for k, v := range activeLinks {
				copyactiveLinks[k] = v
			}
			req <- copyactiveLinks

		case <-stopChan:
			return
		}
	}
}

// AttachManager :
// AttachManager function is the API between single purpose functions and main.go.
// It attaches the XDP program to all the interfaces but loopback and down interfaces,this logic is implemented in findinterfaces.go
func AttachManager(objs *ebpfExport.IpsObjects) error {
	interfaces, err := FindInterfaces()
	if err != nil {
		return fmt.Errorf("fatal: cannot find interfaces from host : %v", err)
	}
	for _, ifacename := range interfaces {
		iface, err := net.InterfaceByName(ifacename)
		if err != nil {
			log.Printf("Error getting interface %s from host : %v", ifacename, err)
			continue
		}
		kernelindex := iface.Index
		linkObj, err := netlink.LinkByName(ifacename)
		if err != nil {
			log.Printf("warning: cannot fetch attr from %s : %v", ifacename, err)
			attachReply(kernelindex, 0, nil, err)
			continue
		}
		attachinfo, err := AttachTypeManager(kernelindex, linkObj, objs)
		if err != nil {
			log.Printf("warning: cannot attach type %s to %s : %v", ifacename, linkObj.Attrs().Name, err)
			attachReply(kernelindex, 0, nil, err)
			continue
		}
		attachReply(kernelindex, attachinfo.Flag, attachinfo.LinkObject, nil)
	}
	return nil
}
