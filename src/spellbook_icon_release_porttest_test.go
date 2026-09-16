//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

type spellbookReleaseResult struct {
	Book     spellbookResult
	Cursor   [7]uint32
	Ability  []uint32
	Bar      []uint32
	Messages [][]byte
	Timeout  uint32
}

func TestSpellbookIconRelease(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	cursor, restore := legacy.PortTestBookAbilityCursorWords()
	t.Cleanup(restore)
	offsets := []uintptr{1047556, 1047916, 1047920, 1047924, 1047928}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = memmap.Uint32(0x5D4594, off)
	}
	t.Cleanup(func() {
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
	})
	ability := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1047788), 30)
	oldAbility := append([]uint32(nil), ability...)
	t.Cleanup(func() { copy(ability, oldAbility) })
	oldTimeouts := inputKeyTimeoutsOld
	t.Cleanup(func() { inputKeyTimeoutsOld = oldTimeouts })
	var rows []spellbookReleaseResult
	kinds := []struct {
		name                 string
		class                byte
		view, flags, instant uint32
	}{
		{"ability-aim", 0, 0, 0, 0}, {"ability-instant", 0, 0, 0, 1},
		{"spell-aim", 1, 0, 0, 0}, {"spell-instant", 1, 0, 0x200, 0},
		{"spell-deferred", 1, 0, 0x2200, 0}, {"guide", 2, 1, 0, 0},
	}
	for _, kind := range kinds {
		for state := 0; state < 4; state++ {
			for hit := 0; hit < 5; hit++ {
				for _, event := range []uint32{6, 7, 8, 12} {
					prepareSpellbookAddition(t, o, bar, 640)
					inputKeyTimeoutsOld = make(map[byte]uint32)
					o.c.srv.NetList.ResetByInd(31, netlist.Kind0)
					for _, off := range offsets {
						*memmap.PtrUint32(0x5D4594, off) = 0
					}
					*cursor[0], *cursor[1] = 0, 0
					clear(ability)
					for id := 1; id <= 5; id++ {
						ability[6*(id-1)] = uint32(id)
					}
					ability[1] = kind.instant
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny) | kind.flags, Valid: true}, {Index: 75, Flags: uint32(things.SpellClassAny) | kind.flags, Valid: true}})
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = kind.class
					*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4232)) = 1
					*o.words["dword_5d4594_1046868"] = kind.view
					*memmap.PtrUint32(0x5D4594, 1046960) = 1
					*memmap.PtrUint32(0x5D4594, 1047508) = 1
					switch state {
					case 1:
						*memmap.PtrUint32(0x5D4594, 1047928) = 1
						*memmap.PtrUint32(0x5D4594, 1047556) = 2
					case 2:
						*cursor[0], *cursor[1] = 1, 3
					case 3:
						*memmap.PtrUint32(0x5D4594, 1047916) = 2
					}
					icon := *o.words["dword_5d4594_1046952"]
					w := (*gui.Window)(unsafe.Pointer(uintptr(icon)))
					origin := w.GlobalPos()
					points := []image.Point{origin.Add(image.Pt(10, 10)), origin.Add(w.SizeVal), origin.Add(w.SizeVal).Add(image.Pt(1, 0)), image.Pt(500, 350), image.Pt(240, 439)}
					pos := points[hit]
					packed := uint32(uint16(pos.X)) | uint32(uint16(pos.Y))<<16
					if o.bookCall("nox_xxx_bookWndFn_45CC10", icon, 5, packed) != 1 || o.c.GUI.Captured() != w || o.c.dragndropSpellType != 1 {
						t.Fatal("icon press must establish drag and capture")
					}
					label := fmt.Sprintf("%s-state%d-hit%d-event%d", kind.name, state, hit, event)
					before := o.bookSnapshot(label, 0)
					ret := o.bookCall("nox_xxx_bookWndFn_45CC10", icon, event, packed)
					r := spellbookReleaseResult{Book: o.bookSnapshot(label, ret), Ability: append([]uint32(nil), ability...), Bar: append([]uint32(nil), bar[:50]...), Timeout: inputKeyTimeoutsOld[5]}
					for i, off := range offsets {
						r.Cursor[i] = memmap.Uint32(0x5D4594, off)
					}
					r.Cursor[5], r.Cursor[6] = *cursor[0], *cursor[1]
					o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
					want := uint32(1)
					if event == 12 {
						want = 0
					}
					if ret != want {
						t.Fatalf("%s result%d", label, ret)
					}
					if event == 6 || event == 7 {
						if o.c.dragndropSpellType != 0 {
							t.Fatalf("%s release must clear drag", label)
						}
						if hit >= 2 && (o.c.GUI.Captured() != nil || *o.words["dword_5d4594_1047540"] != 0) {
							t.Fatalf("%s outside release must clear capture/selection", label)
						}
						if kind.name == "ability-instant" && state == 0 && hit < 2 && !reflect.DeepEqual(r.Messages, [][]byte{{122, 1}}) {
							t.Fatalf("%s immediate ability message %v", label, r.Messages)
						}
					} else {
						before.Return = ret
						if !reflect.DeepEqual(before, r.Book) {
							t.Fatalf("%s non-release event changed book state", label)
						}
					}
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "icon-release", rows, "83a1b3082f043fb1117fb69154d7319018efb282db510fdcf0db188cb6fd543e")
}
