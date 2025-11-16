package netlink

import (
	"fmt"

	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

func AttachTypeManager(kernelindex int, linkObj netlink.Link, objs *ebpfExport.IpsObjects, logger *zap.Logger) (NewLink, error) {
	var err error
	var Result NewLink
	switch linkObj.(type) {
	case *netlink.Device:
		Result, err = AttachPhysicalDevice(kernelindex, linkObj, objs)
		if err != nil {
			logger.Warn("attach failed",
				zap.Int("kernelIndex", kernelindex),
				zap.String("ifaceName", linkObj.Attrs().Name),
				zap.Error(err))
			return Result, err
		}
	case *netlink.Veth, *netlink.Bridge, *netlink.Vlan, *netlink.Dummy:
		Result, err = AttachVirtualDevice(kernelindex, linkObj, objs, logger)
		if err != nil {
			logger.Warn("attach failed",
				zap.Int("kernelIndex", kernelindex),
				zap.String("ifaceName", linkObj.Attrs().Name),
				zap.Error(err))
			return Result, err
		}
	default:
		logger.Warn("unknown link type",
			zap.Int("kernelIndex", kernelindex),
			zap.String("ifaceName", linkObj.Attrs().Name))
		return Result, fmt.Errorf("unknown link type")
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
	}
	return NewLink{}, err
}
func AttachVirtualDevice(kernelindex int, linkObj netlink.Link, objs *ebpfExport.IpsObjects, logger *zap.Logger) (NewLink, error) {
	masterindex := linkObj.Attrs().MasterIndex
	if masterindex != 0 {
		logger.Debug("iface is attached to master",
			zap.String("ifaceName", linkObj.Attrs().Name),
			zap.Int("masterIndex", masterindex))
		return NewLink{}, fmt.Errorf("iface is attached to a master")
	}
	result, err := Attach(kernelindex, link.XDPGenericMode, objs)
	if err == nil {
		replyinfo := NewLink{
			LinkIndex:  linkObj.Attrs().Index,
			Flag:       link.XDPGenericMode,
			LinkObject: result,
		}
		return replyinfo, nil
	}
	return NewLink{}, err
}
