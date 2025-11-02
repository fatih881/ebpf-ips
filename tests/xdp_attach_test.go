package tests

import (
	"net"
	"testing"

	localebpf "github.com/fatih881/ebpf-ips/core/ebpf"
	localnl "github.com/fatih881/ebpf-ips/core/netlink"
	"github.com/vishvananda/netlink"
)

func TestAttach(t *testing.T) {
	attrs := netlink.LinkAttrs{
		Name: "xdpdummy0",
	}
	dummyLink := &netlink.Dummy{
		LinkAttrs: attrs,
	}
	err := netlink.LinkAdd(dummyLink)
	if err != nil {
		t.Fatalf("warning: cannot add dummy for xdp : %v", err)
	}
	defer netlink.LinkDel(dummyLink)
	objs := &localebpf.IpsObjects{}
	err = localebpf.LoadIpsObjects(objs, nil)
	if err != nil {
		t.Fatalf("warning: cannot load program to kernel : %v", err)
	}
	attachresult, err := localnl.Attach([]string{"xdpdummy0"}, objs)
	if err != nil {
		t.Fatalf("warning: function cannot attach program : %v", err)
	}
	defer objs.Close()
	iface, err := net.InterfaceByName("xdpdummy0")
	if err != nil {
		t.Fatalf("warning: cannot get interface index for dummy : %v", err)
	}
	linkobj, err := netlink.LinkByName("xdpdummy0")
	if err != nil {
		t.Fatalf("warning: error at query : %v", err)
	}
	queryresult := linkobj.Attrs().Xdp

	if queryresult == nil {
		t.Fatal("warning: cannot find xdp at dummy")
	}
	if attachresult[iface.Index] == nil {

		t.Fatalf("warning: function returned nil as ActiveLinks variable on xdpdummy0 : %v ", err)

	}
	//TODO:IMPLEMENT : CHECKING ATTACHRESULT IS AQUALS TO PROGRAM ID OR NOT.

}
