package netlink

import (
	"net"

	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
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

// StartLinkManager must be started before AttachExistingInterfaces function.
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

// AttachExistingInterfaces :
// AttachExistingInterfaces attaches XDP to all available interfaces on the system at startup.
// After this function,we will handle new links event-drivenly with kernel subscribing.
// It attaches the XDP program to all the interfaces but loopback and down interfaces,this logic is implemented in findinterfaces.go
func AttachExistingInterfaces(objs *ebpfExport.IpsObjects, logger *zap.Logger) error {
	interfaces, err := FindInterfaces(logger)
	if err != nil {
		return err
	}
	var successCount, failCount int
	for _, ifacename := range interfaces {
		iface, err := net.InterfaceByName(ifacename)
		if err != nil {
			logger.Debug("interface lookup failed",
				zap.Error(err))
			failCount++
			continue
		}
		kernelindex := iface.Index
		linkObj, err := netlink.LinkByName(ifacename)
		if err != nil {
			logger.Debug("netlink fetch failed",
				zap.String("ifaceName", iface.Name),
				zap.Int("kernelIndex", kernelindex),
				zap.Error(err))
			attachReply(kernelindex, 0, nil, err)
			failCount++
			continue
		}
		attachinfo, err := AttachTypeManager(kernelindex, linkObj, objs, logger)
		if err != nil {
			logger.Warn("attach failed",
				zap.String("ifaceName", iface.Name),
				zap.Int("kernelIndex", kernelindex),
				zap.Error(err))
			attachReply(kernelindex, 0, nil, err)
			failCount++
			continue
		}
		attachReply(kernelindex, attachinfo.Flag, attachinfo.LinkObject, nil)
		successCount++
	}
	logger.Info("Attach Manager loop completed",
		zap.Int("success", successCount),
		zap.Int("fail", failCount),
		zap.Int("total", len(interfaces)))
	return nil
}
