package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/cilium/ebpf/link"
	"github.com/fatih881/ebpf-ips/core/ebpf"
	localnl "github.com/fatih881/ebpf-ips/core/netlink"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	objs := &ebpf.IpsObjects{}
	err := ebpf.LoadIpsObjects(objs, nil)
	if err != nil {
		log.Fatalf("Error loading objects: %v", err)
		return
	}
	defer func(objs *ebpf.IpsObjects) {
		err := objs.Close()
		if err != nil {
			log.Fatal("Error closing objects", zap.Error(err))
			return
		}
	}(objs)
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't create logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("Can't stop logger: %v", err)
		}
	}(logger)
	Updates := make(chan netlink.LinkUpdate)
	done := make(chan struct{})
	go localnl.Subscribetokernel(Updates, done, logger)
	WriteChan := make(chan localnl.WriteChanMessage)
	ReadChan := make(chan chan map[int]link.Link)
	StopChan := make(chan struct{})
	DeleteChan := make(chan int)
	go localnl.StartLinkManager(WriteChan, ReadChan, StopChan, DeleteChan, logger)
	err = localnl.AttachExistingInterfaces(objs, WriteChan, logger)
	if err != nil {
		log.Fatalf("Attach Existing Interfaces returned err : %v", err)
		return
	}
	go localnl.HandleKernelMessage(WriteChan, Updates, ReadChan, DeleteChan, objs, logger)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	close(done)
}
