package netlink

import "github.com/cilium/ebpf/link"

var (
	XDPLinks    = make([]link.Link, 0)
	ActiveLinks = make(map[int]link.Link)
)

func xdpmanager() {}
