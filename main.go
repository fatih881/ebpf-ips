package main

import (
	"log"

	"github.com/cilium/ebpf/link"
	"github.com/fatih881/ebpf-ips/core/netlink"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
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
	WriteChan := make(chan netlink.NewLink)
	ReadChan := make(chan chan map[int]link.Link)
	StopChan := make(chan struct{})
	netlink.StartLinkManager(WriteChan, ReadChan, StopChan)
}
