package netlink

import (
	"github.com/cilium/ebpf/link"
)

// attachReply wrapper :
// Since we have to use WriteChan for attach attempt no matter it fails or not,
// this wrapper will be used for WriteChan so we can easily change the logic of replying and sending it to channel.
func attachReply(index int, flags link.XDPAttachFlags, LinkObject link.Link, receivederr error) {
	if receivederr == nil {
		newLink := NewLink{
			LinkIndex:  index,
			Flag:       flags,
			LinkObject: LinkObject,
		}
		WriteChan <- struct {
			NewLink NewLink
			Err     error
		}{newLink, nil}
	} else {
		WriteChan <- struct {
			NewLink NewLink
			Err     error
		}{NewLink: NewLink{LinkIndex: index},
			Err: receivederr}
	}
	return
}
