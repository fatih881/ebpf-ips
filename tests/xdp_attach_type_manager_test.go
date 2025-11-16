package tests

import (
	"testing"

	localebpf "github.com/fatih881/ebpf-ips/core/ebpf"
	localnl "github.com/fatih881/ebpf-ips/core/netlink"
	"github.com/mazen160/go-random"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

func TestXDPAttachTypeManager(t *testing.T) {
	//The sum of the generated and fixed values ​cannot be greater than 15. (16 bytes normally, but \0 uses 1.)
	// Source:https://elixir.bootlin.com/linux/v5.6/source/include/uapi/linux/if.h#L33
	randominterfacenumber, err := random.String(11)
	if err != nil {
		t.Fatalf("warning: cannot generate random interface number : %v", err)
	}
	dummyLinkName := "xdp-" + randominterfacenumber
	dummyLink := &netlink.Dummy{
		LinkAttrs: netlink.LinkAttrs{Name: dummyLinkName},
	}
	err = netlink.LinkAdd(dummyLink)
	if err != nil {
		t.Fatalf("cannot add dummy for xdp : %v", err)
	}
	t.Cleanup(func() {
		if err := netlink.LinkDel(dummyLink); err != nil {
			t.Logf("failed to delete dummy : %v", err)
		}
	})
	objs := &localebpf.IpsObjects{}
	err = localebpf.LoadIpsObjects(objs, nil)
	if err != nil {
		t.Fatalf("cannot load program to kernel : %v", err)
	}
	t.Cleanup(func() {
		if err := objs.Close(); err != nil {
			t.Logf("failed to remove objects : %v", err)
		}
	})
	kernel, err := netlink.LinkByName(dummyLinkName)
	if err != nil {
		t.Fatalf("cannot lookup kernel index : %v", err)
	}
	kernelindex := kernel.Attrs().Index
	logger := zap.NewNop()
	_, err = localnl.AttachTypeManager(kernelindex, kernel, objs, logger)
	if err != nil {
		t.Fatalf("function cannot attach program : %v", err)
	}
	linkobj, err := netlink.LinkByName(dummyLinkName)
	if err != nil {
		t.Fatalf("error at query : %v", err)
	}
	if linkobj.Attrs().Xdp == nil {
		t.Fatalf("XDP not attached")
	}

}
