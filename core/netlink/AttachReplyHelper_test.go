package netlink

import (
	"errors"
	"testing"

	"github.com/cilium/ebpf/link"
)

func TestAttachReply(t *testing.T) {
	TestWriteChan := make(chan WriteChanMessage, 1)
	defer close(TestWriteChan)
	t.Run("success", func(t *testing.T) {
		attachReply(1, link.XDPOffloadMode, nil, TestWriteChan, nil)
		result := <-TestWriteChan
		if result.Err != nil {
			t.Errorf("Error should be nil, but got %v", result.Err)
		}
		if result.NewLink.LinkIndex != 1 {
			t.Errorf("Index should be 1, but got %d", result.NewLink.LinkIndex)
		}
		if result.NewLink.Flag != link.XDPOffloadMode {
			t.Errorf("Flag should be %d, but got %d", link.XDPOffloadMode, result.NewLink.Flag)
		}
	})
	t.Run("fail", func(t *testing.T) {
		attachReply(1, 0, nil, TestWriteChan, errors.New("fail"))
		result := <-TestWriteChan
		if result.Err == nil {
			t.Errorf("Error should not be nil, but got %v", result.Err)
		}
		if result.NewLink.LinkIndex != 1 {
			t.Errorf("Index should be 1, but got %d", result.NewLink.LinkIndex)
		}
		if result.NewLink.Flag != 0 {
			t.Errorf("Flag should be %d, but got %d", 0, result.NewLink.Flag)
		}

	})
}
