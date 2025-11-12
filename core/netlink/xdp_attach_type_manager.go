package netlink

import (
	"fmt"
	"log"

	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
)

func AttachTypeManager(iface string, linkObj netlink.Link, objs *ebpfExport.IpsObjects) (NewLink, error) {
	var err error
	var Result NewLink
	switch linkObj.(type) {
	case *netlink.Device:
		Result, err = AttachPhysicalDevice(iface, linkObj, objs, WriteChan)
		if err != nil {
			log.Printf("warning: cannot attach physical device: %v", err)
			return Result, err
		}
	case *netlink.Veth, *netlink.Bridge, *netlink.Vlan, *netlink.Dummy:
		Result, err = AttachVirtualDevice(iface, linkObj, objs, WriteChan)
		if err != nil {
			log.Printf("warning: cannot attach virtual device: %v", err)
			return Result, err
		}
	default:
		log.Printf("warning: cannot get device type for %s", linkObj.Attrs().Name)
		return Result, fmt.Errorf("cannot get device type for %s : %v", linkObj.Attrs().Name, err)
	}
	return Result, nil

}

// AttachPhysicalDevice and AttachVirtualDevice both returns and sends the same values.
// return have the purpose of making checks event-drivenly on main attach manager,and WriteChan is for logging purposes.
func AttachPhysicalDevice(iface string, linkObj netlink.Link, objs *ebpfExport.IpsObjects, WriteChan chan NewLink) (NewLink, error) {

	result, err := Attach(iface, link.XDPOffloadMode, objs)
	if err == nil {
		attachReply := NewLink{ // AttachReply will be replaced a wrapper because it's hard to make a refactor on attachReply right now.
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPOffloadMode,
			LinkObject: result,
		}
		WriteChan <- attachReply
		return attachReply, nil
	} else {
		result, err = Attach(iface, link.XDPDriverMode, objs)
	}
	if err == nil {
		attachReply := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPDriverMode,
			LinkObject: result,
		}
		WriteChan <- attachReply
		return attachReply, nil
	} else {
		result, err = Attach(iface, link.XDPGenericMode, objs)
	}
	if err == nil {
		attachReply := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPGenericMode,
			LinkObject: result,
		}
		WriteChan <- attachReply
		return attachReply, nil
	} else {
		log.Printf("warning: cannot load XDP to %s : %v", iface, err)
	}
	attachReply := NewLink{
		LinkIndex:  linkObj.Attrs().Index,
		Flag:       0,
		LinkObject: result,
	}
	return attachReply, fmt.Errorf("cannot load XDP to %s : %v", iface, err)
}
func AttachVirtualDevice(iface string, linkObj netlink.Link, objs *ebpfExport.IpsObjects, WriteChan chan NewLink) (NewLink, error) {
	result, err := Attach(iface, link.XDPGenericMode, objs)
	if err == nil {
		log.Printf("info: generic mode is attached to %s", iface)
		attachReply := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPGenericMode,
			LinkObject: result,
		}
		WriteChan <- attachReply
		return attachReply, nil
	} else {
		log.Printf("warning: cannot load XDP to %s : %v", iface, err)

	}
	attachReply := NewLink{
		LinkIndex:  linkObj.Attrs().Index,
		Flag:       0,
		LinkObject: result,
	}
	return attachReply, fmt.Errorf("cannot load XDP to %s : %v", iface, err)
}
