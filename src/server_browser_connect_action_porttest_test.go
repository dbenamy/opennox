//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"net/netip"
	"testing"
	"unsafe"
)

func TestServerBrowserPasswordAction(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	clocks, restoreClocks := legacy.PortTestServerBrowserClocks()
	defer restoreClocks()
	*words["dword_5d4594_815056"] = 1
	port := serverConfigOwnBytes(t, 0x5D4594, 814604, 4)
	binary.LittleEndian.PutUint32(port, 18590)
	sender := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 10, 10, nil)
	sender.SetID(4001)
	record, free := alloc.Make([]byte{}, 172)
	defer free()
	copy(record[12:28], "127.0.0.1")
	binary.LittleEndian.PutUint16(record[109:], 18590)
	*words["dword_5d4594_814624"] = uint32(uintptr(unsafe.Pointer(&record[0])))
	getText, send, ticks, buttons := legacy.Sub_449E60, legacy.SendXXX_5550D0, legacy.PlatformTicks, legacy.Sub_449EA0
	defer func() {
		legacy.Sub_449E60 = getText
		legacy.SendXXX_5550D0 = send
		legacy.PlatformTicks = ticks
		legacy.Sub_449EA0 = buttons
	}()
	for _, text := range [][]uint16{{}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}, {0xd83d, 0xde00, 0xff, 0x8000, 'a', 'b', 'c', 'd', 'e'}} {
		p, release := alloc.Make([]uint16{}, len(text)+1)
		copy(p, text)
		legacy.Sub_449E60 = func(v int8) int {
			if v != 4 {
				t.Fatal("text field", v)
			}
			return int(uintptr(unsafe.Pointer(&p[0])))
		}
		for _, now := range []uint64{1, 0xffffffff, 1 << 32} {
			*words["dword_5d4594_814548"] = 6
			legacy.PlatformTicks = func() uint64 { return now }
			sent, disabled := 0, 0
			legacy.Sub_449EA0 = func(flags gui.DialogFlags) {
				if flags != 0 {
					t.Fatal("buttons", flags)
				}
				disabled++
			}
			legacy.SendXXX_5550D0 = func(addr netip.AddrPort, data []byte) (int, error) {
				if addr != netip.MustParseAddrPort("127.0.0.1:18590") || len(data) != 22 {
					t.Fatal("destination/size", addr, len(data))
				}
				want := make([]byte, 18)
				for i, v := range text {
					if i == 8 {
						break
					}
					binary.LittleEndian.PutUint16(want[i*2:], v)
				}
				// Original C leaves the four header bytes uninitialized; the sender writes
				// its first three bytes. The Go port also initializes the padding byte.
				if !bytes.Equal(data[4:], want) {
					t.Fatal("password units", data[4:], want)
				}
				sent++
				return 22, nil
			}
			if legacy.PortTestServerBrowserEvent(o.parent.C(), 16391, uint32(uintptr(sender.C())), 0) != 0 || sent != 1 || disabled != 1 || *words["dword_5d4594_814548"] != 3 || *clocks["qword_5d4594_814956"] != uint64(uint32(now)+20000) {
				t.Fatal("password action", now, sent, disabled, *words["dword_5d4594_814548"], *clocks["qword_5d4594_814956"])
			}
		}
		release()
	}
}
