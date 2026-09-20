//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netstr"
	"testing"
	"unsafe"
)

func TestReplayMessageCopy(t *testing.T) {
	o := newReliableReportsOwner(t)
	s := &Server{Server: o.s}
	for _, payload := range [][]byte{{0x25}, {0x25, 0x55, 0xaa, 0xff}, nil} {
		for i := range o.units {
			pl := o.units[i].UpdateDataPlayer().Player
			pl.Frame3596 = 1
			before := bytes.Clone(payload)
			if !s.onPacket(ntype.PlayerInd(pl.PlayerInd), payload) {
				t.Fatalf("replay message %x was not accepted", payload)
			}
			if pl.Frame3596 != s.Frame() {
				t.Fatalf("replay message %x: frame=%d want=%d", payload, pl.Frame3596, s.Frame())
			}
			if !bytes.Equal(payload, before) {
				t.Fatalf("replay changed its input: %x -> %x", before, payload)
			}
		}
	}
}

func TestReplayMessagePayloadCopy(t *testing.T) {
	o := newReliableReportsOwner(t)
	oldStreams := o.s.NetStr
	o.s.NetStr = &netstr.Streams{}
	t.Cleanup(func() { o.s.NetStr = oldStreams })
	s := &Server{Server: o.s}
	payload := []byte{165, 0, 0xff, 0x80, 0x55, 0xaa, 1, 2, 3, 4, 165, 254, 0x12, 0x34, 0x56, 0x78, 0xfe, 0xdc, 0xba, 0x98}
	before := bytes.Clone(payload)
	for i := range o.units {
		pl := o.units[i].UpdateDataPlayer().Player
		saved := pl.NetData16
		raw := unsafe.Slice((*byte)(unsafe.Pointer(&pl.NetData16[0])), 255*8)
		clear(raw)
		pl.Frame3596 = 1
		ok := s.onPacket(pl.PlayerIndex(), payload)
		got := bytes.Clone(raw)
		pl.NetData16 = saved
		want := make([]byte, 255*8)
		copy(want[:8], payload[2:10])
		copy(want[254*8:], payload[12:20])
		if !ok || !bytes.Equal(got, want) || pl.Frame3596 != s.Frame() {
			t.Fatalf("replayed alias payload differs: accepted=%v frame=%d", ok, pl.Frame3596)
		}
		if !bytes.Equal(payload, before) {
			t.Fatal("replay changed input bytes")
		}
	}
}
