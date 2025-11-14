package netlink

import (
	"fmt"
	"log"

	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
)

func AttachTypeManager(kernelindex int, linkObj netlink.Link, objs *ebpfExport.IpsObjects) (NewLink, error) {
	var err error
	var Result NewLink
	switch linkObj.(type) {
	case *netlink.Device:
		Result, err = AttachPhysicalDevice(kernelindex, linkObj, objs)
		if err != nil {
			log.Printf("warning: cannot attach physical device: %v", err)
			return Result, err
		}
	case *netlink.Veth, *netlink.Bridge, *netlink.Vlan, *netlink.Dummy:
		Result, err = AttachVirtualDevice(kernelindex, linkObj, objs)
		if err != nil {
			log.Printf("warning: cannot attach virtual device: %v", err)
			return Result, err
		}
	default:
		return Result, fmt.Errorf("cannot get device type for %s", linkObj.Attrs().Name)
	}
	return Result, nil

}
func AttachPhysicalDevice(kernelindex int, linkObj netlink.Link, objs *ebpfExport.IpsObjects) (NewLink, error) {

	result, err := Attach(kernelindex, link.XDPOffloadMode, objs)
	if err == nil {
		replyinfo := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPOffloadMode,
			LinkObject: result,
		}
		return replyinfo, nil
	} else {
		result, err = Attach(kernelindex, link.XDPDriverMode, objs)
	}
	if err == nil {
		replyinfo := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPDriverMode,
			LinkObject: result,
		}
		return replyinfo, nil
	} else {
		result, err = Attach(kernelindex, link.XDPGenericMode, objs)
	}
	if err == nil {
		replyinfo := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPGenericMode,
			LinkObject: result,
		}
		return replyinfo, nil
	} else {
		log.Printf("warning: cannot load XDP to %v : %v", kernelindex, err)
	}
	return NewLink{}, fmt.Errorf("cannot load XDP to %v : %v", kernelindex, err)
}
func AttachVirtualDevice(kernelindex int, linkObj netlink.Link, objs *ebpfExport.IpsObjects) (NewLink, error) {
	result, err := Attach(kernelindex, link.XDPGenericMode, objs)
	if err == nil {
		log.Printf("info: generic mode is attached to %v", kernelindex)
		replyinfo := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPGenericMode,
			LinkObject: result,
		}
		return replyinfo, nil
	} else {
		log.Printf("warning: cannot load XDP to %v : %v", kernelindex, err)

	}
	return NewLink{}, fmt.Errorf("cannot load XDP to %v : %v", kernelindex, err)
}
