package netlink

import (
	"net"

	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

var (
	attachSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "xdp_attach_success_total",
	})

	attachFailTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "xdp_attach_failed_total",
	})
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
func getCurrentLinks(readchan chan<- chan map[int]link.Link) map[int]link.Link {
	respChan := make(chan map[int]link.Link)
	readchan <- respChan
	return <-respChan
}

// AttachExistingInterfaces attaches XDP program to all existing interfaces on the system at startup.
func AttachExistingInterfaces(objs *ebpfExport.IpsObjects, logger *zap.Logger) error {
	interfaces, err := FindInterfaces(logger)
	if err != nil {
		return err
	}
	for _, ifacename := range interfaces {
		iface, err := net.InterfaceByName(ifacename)
		if err != nil {
			logger.Debug("interface lookup failed",
				zap.Error(err))
			attachFailTotal.Inc()
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
			attachFailTotal.Inc()
			continue
		}
		attachinfo, err := AttachTypeManager(kernelindex, linkObj, objs, logger)
		if err != nil {
			logger.Warn("attach failed",
				zap.String("ifaceName", iface.Name),
				zap.Int("kernelIndex", kernelindex),
				zap.Error(err))
			attachReply(kernelindex, 0, nil, err)
			attachFailTotal.Inc()
			continue
		}
		attachReply(kernelindex, attachinfo.Flag, attachinfo.LinkObject, nil)
		attachSuccessTotal.Inc()
	}
	logger.Info("attach snapshot completed")
	return nil
}
