//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookKnowledgeRemoval(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	bar := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1049220), 64)
	oldBar := append([]uint32(nil), bar...)
	t.Cleanup(func() { copy(bar, oldBar) })
	t.Cleanup(legacy.PortTestBookQuickbar(unsafe.Pointer(&bar[0])))
	raw := unsafe.Slice(memmap.PtrUint8(0x587000, 132100), 32)
	oldFamily := append([]byte(nil), raw...)
	t.Cleanup(func() { copy(raw, oldFamily) })
	copy(raw, blobdata.PortTestBookGuideFamily())
	*memmap.PtrPtr(0x587000, 132124) = memmap.PtrOff(0x587000, 132100)
	guides := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 740076)), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	for i := 1; i <= 40; i++ {
		guides[7*i+1] = 1
	}
	var rows []struct {
		Book  spellbookResult
		Slots []uint32
	}
	for _, guide := range []bool{false, true} {
		for family := 0; family < 4; family++ {
			for selected := 0; selected < 5; selected++ {
				if guide && family > 1 {
					continue
				}
				o.resetBook(t)
				if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
					t.Fatal("book setup")
				}
				*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = 2
				clear(bar)
				bar[50] = uint32(selected)
				parent := []uint32{0, 0x1000, 0x4000, 0x10000}[family]
				child := []uint32{0, 0x2000, 0x8000, 0x20000}[family]
				configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny) | parent, Valid: true}, {Index: 2, Flags: uint32(things.SpellClassAny) | child, Valid: true}, {Index: 3, Flags: uint32(things.SpellClassAny) | child, Valid: false}, {Index: 4, Flags: uint32(things.SpellClassAny), Valid: true}})
				id := uint32(1)
				levels := unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696)), 137)
				remove := map[uint32]bool{1: true}
				op := "sub_45D320"
				if family != 0 {
					remove[2] = true
				}
				if guide {
					op = "sub_45D400"
					levels = unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4244)), 41)
					*o.words["dword_5d4594_1046868"] = 1
					*o.words["dword_5d4594_1046872"] = 1
					remove = map[uint32]bool{1: true}
					if family == 1 {
						id = 24
						remove = map[uint32]bool{24: true, 7: true, 8: true, 25: true, 26: true}
					}
				}
				for i := 1; i < len(levels); i++ {
					levels[i] = 7
				}
				keys := []uint32{1, 2, 3, 4, 1}
				if guide {
					keys = []uint32{1, 24, 7, 8, 25, 26, 36, 40}
				}
				for slot := 0; slot < 25; slot++ {
					entry := keys[slot%len(keys)]
					if guide {
						entry += 74
					}
					bar[2*slot] = entry
					bar[2*slot+1] = 0xa1000000 + uint32(slot)
				}
				ret := o.bookCall(op, id)
				for i := 1; i < len(levels); i++ {
					want := uint32(7)
					if remove[uint32(i)] {
						want = 0
					}
					if levels[i] != want {
						t.Fatalf("guide%v family%d level%d: %d want%d", guide, family, i, levels[i], want)
					}
				}
				for slot := 0; slot < 25; slot++ {
					key := keys[slot%len(keys)]
					want := key
					if guide {
						want += 74
					}
					if remove[key] {
						want = 0
					}
					if bar[2*slot] != want || bar[2*slot+1] != 0xa1000000+uint32(slot) {
						t.Fatal("removal must clear every matching slot and preserve flags")
					}
				}
				if bar[50] != uint32(selected) {
					t.Fatal("selected row retained")
				}
				rows = append(rows, struct {
					Book  spellbookResult
					Slots []uint32
				}{o.bookSnapshot(fmt.Sprintf("guide%v-family%d-row%d", guide, family, selected), ret), append([]uint32(nil), bar...)})
			}
		}
	}
	spellbookCapture(t, "removal", rows, "")
}
