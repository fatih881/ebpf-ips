package netlink

import (
	"github.com/cilium/ebpf/link"
	ebpfExport "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

func Subscribetokernel(updates chan netlink.LinkUpdate, done chan struct{}, logger *zap.Logger) {
	err := netlink.LinkSubscribe(updates, done)
	if err != nil {
		logger.Error("Error subscribing to kernel updates", zap.Error(err))
	}
}
func HandleKernelMessage(updates chan netlink.LinkUpdate, readchan chan chan map[int]link.Link, objs *ebpfExport.IpsObjects, logger *zap.Logger) {
	for update := range updates {
		switch update.Header.Type {
		case unix.RTM_NEWLINK:
			if update.Link.Attrs().OperState != netlink.OperUp {
				continue
			}
			activelinks := getCurrentLinks(readchan)
			_, exists := activelinks[int(update.Index)]
			if exists {
				continue
			}
			linkObj, err := netlink.LinkByIndex(int(update.Index))
			if err != nil {
				logger.Error("Error getting link by index", zap.Int("index", int(update.Index)))
				continue
			}
			result, err := AttachTypeManager(int(update.Index), linkObj, objs, logger)
			if err != nil {
				attachFailTotal.Inc()
				continue
			}
			attachReply(result.LinkIndex, result.Flag, result.LinkObject, nil)
			attachSuccessTotal.Inc()
		}
	}
}
