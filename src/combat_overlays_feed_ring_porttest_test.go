//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestCombatOverlayFeedRing(t *testing.T) {
	type record struct {
		Step        int
		Read, Write uint32
		Row         [6]uint32
	}
	var rows []record
	o := newCombatOverlayOwner(t)
	for i := 0; i < 205; i++ {
		var p [11]byte
		binary.LittleEndian.PutUint16(p[8:], uint16(i))
		p[10] = 2
		o.c.srv.SetFrame(uint32(1000 + i))
		legacy.PortTestCombatFeedAdd(unsafe.Pointer(&p[0]))
		r := *(*[6]uint32)(memmap.PtrOff(0x5D4594, 1201428+24*uintptr(i%100)))
		read, write := *o.words["feedRead"], *o.words["feedWrite"]
		if r != [6]uint32{0, 0, 0, uint32(i), 2, uint32(1000 + i)} || write != uint32((i+1)%100) || read != uint32(max(0, i-98)%100) {
			t.Fatalf("ring at insertion %d: %v %d/%d", i, r, read, write)
		}
		rows = append(rows, record{i, read, write, r})
	}
	spellbookCapture(t, "combat-overlay-feed-ring", rows, "8f2d85668556cb4f57fa869d5ea924c97c655ec084ab69515f6fafc094e49d5c")
}
func TestCombatOverlayFeedExpiry(t *testing.T) {
	type record struct {
		Name       string
		Read, Rows uint32
		Pixels     string
	}
	var rows []record
	for _, frame := range []uint32{20, 1000, 0xffffffff} {
		for _, age := range []uint32{0, 89, 90, 91, 100} {
			for _, count := range []int{0, 1, 4, 5, 10} {
				name := fmt.Sprintf("frame=%d/age=%d/count=%d", frame, age, count)
				t.Run(name, func(t *testing.T) {
					o := newCombatOverlayOwner(t)
					o.c.srv.SetFrame(frame)
					for i := 0; i < count; i++ {
						p := (*[6]uint32)(memmap.PtrOff(0x5D4594, 1201428+24*uintptr(i)))
						*p = [6]uint32{0, 0, 0, 0, 0, frame - age}
					}
					*o.words["feedWrite"] = uint32(count)
					legacy.PortTestCombatFeedDraw()
					wantRead, wantRows := uint32(0), uint32(min(count, 4))
					if age > 90 {
						wantRead = uint32(count)
						wantRows = 0
					}
					if *o.words["feedRead"] != wantRead || *o.words["feedRows"] != wantRows {
						t.Fatalf("expiry read=%d rows=%d want %d/%d", *o.words["feedRead"], *o.words["feedRows"], wantRead, wantRows)
					}
					rows = append(rows, record{name, *o.words["feedRead"], *o.words["feedRows"], effectsPixelHash(o.pix)})
				})
			}
		}
	}
	spellbookCapture(t, "combat-overlay-feed-expiry", rows, "7f89fa48222dcf56d9423d495eb94df054fd1f33c50afdd82bf473cb91daa9f4")
}
