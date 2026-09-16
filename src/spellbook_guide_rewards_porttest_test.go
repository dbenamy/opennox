//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
)

func TestSpellbookGuideKnowledge(t *testing.T) {
	o := newSpellbookOwner(t)
	raw := unsafe.Slice(memmap.PtrUint8(0x587000, 132100), 32)
	old := append([]byte(nil), raw...)
	t.Cleanup(func() { copy(raw, old) })
	copy(raw, blobdata.PortTestBookGuideFamily())
	*memmap.PtrPtr(0x587000, 132124) = memmap.PtrOff(0x587000, 132100)
	guides := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 740076)), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	for i := 1; i <= 40; i++ {
		guides[7*i+1] = 1
	}
	var rows []spellbookResult
	for _, id := range []uint32{1, 7, 8, 24, 25, 26, 36, 37, 40} {
		for _, present := range []bool{false, true} {
			for _, known := range []uint32{0, 1, 7} {
				o.resetBook(t)
				if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
					t.Fatal("book setup")
				}
				levels := unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4244)), 41)
				for i := range levels {
					levels[i] = known
				}
				if !present {
					*o.words["dword_8531A0_2576"] = 0
				}
				ret := o.bookCall("nox_xxx_netGuideRewardCli_45D140", id, 0)
				for i := range levels {
					want := known
					family := id == 24 && (i == 7 || i == 8 || i == 25 || i == 26)
					if present && (i == int(id) || family) {
						want = 1
					}
					if levels[i] != want {
						t.Fatalf("guide%d slot%d got%d want%d", id, i, levels[i], want)
					}
				}
				rows = append(rows, o.bookSnapshot(fmt.Sprintf("guide%d-player%v-known%d", id, present, known), ret))
			}
		}
	}
	spellbookCapture(t, "guide-knowledge", rows, "")
}
