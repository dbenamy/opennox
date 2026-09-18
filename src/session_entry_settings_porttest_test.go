//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestSessionEntrySettingsBroadcast(t *testing.T) {
	o := newServerOptionsOwner(t)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1563214, 58)
	clear(cache)
	reset, packets, restore := legacy.PortTestVisibilityReliableOwner()
	defer restore()
	for i := range o.players {
		o.players[i].Active = 0
		o.players[i].PlayerInd = byte(i)
	}
	for _, i := range []int{1, 7, 31} {
		o.players[i].Active = 1
	}
	legacy.PortTestServerOptionsModeName(0)
	o.event(10119, 16385, uintptr(unsafe.Pointer(alloc.InternCString16("Capture flag"))), 0)
	o.options.ChildByID(10122).DrawData().Field0 &^= 4
	type row struct {
		Name, Score, Time string
		Repeat            bool
		Cache             []byte
		Queue             []legacy.PortTestShopPacketResult
	}
	var rows []row
	for _, v := range [][3]string{{"", "", ""}, {"Server name", "17", "9"}, {"Next", "65536", "256"}} {
		o.text(10101, v[0])
		o.text(10134, v[1])
		o.text(10135, v[2])
		want := make([]byte, 58)
		copy(want[9:24], v[0])
		binary.LittleEndian.PutUint16(want[52:], 0x20)
		if v[1] == "17" {
			binary.LittleEndian.PutUint16(want[54:], 17)
			want[56] = 9
		}
		for _, repeat := range []bool{false, true} {
			reset()
			legacy.Sub_4DF020()
			q := packets()
			if !bytes.Equal(cache, want) {
				t.Errorf("settings %v repeat%v cache%x want%x", v, repeat, cache, want)
			}
			wantCount := 2
			if repeat {
				wantCount = 0
			}
			if len(q) != wantCount {
				t.Errorf("settings repeat%v reports%d want%d", repeat, len(q), wantCount)
			}
			recipients := map[byte]bool{}
			for _, p := range q {
				recipients[p.Recipient] = true
				if !bytes.Equal(p.Data, append([]byte{177, 1}, want...)) {
					t.Errorf("settings report%x", p.Data)
				}
			}
			if !repeat && (!recipients[1] || !recipients[7] || recipients[31]) {
				t.Errorf("settings recipients%v", recipients)
			}
			rows = append(rows, row{v[0], v[1], v[2], repeat, bytes.Clone(cache), q})
		}
	}
	spellbookCapture(t, "session-entry-settings-broadcast", rows, "5d516aeb9aa0741dafcba8811a9c5547aa4225c995db9c21385d1b108f77cbcb")
}
