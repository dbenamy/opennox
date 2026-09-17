//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerPanelsObjectList(t *testing.T) {
	type row struct {
		Kind            string
		Flags, Mask     uint32
		Count           int
		Names           []string
		Hidden, Checked []bool
		After           []uint32
	}
	var rows []row
	for _, kind := range []string{"weapon", "armor"} {
		for _, flags := range []uint32{0, 1, 0x4001, 0x8001} {
			for _, mask := range []uint32{0, 0x55555555, 0xffffffff} {
				t.Run(fmt.Sprintf("%s-%x-%x", kind, flags, mask), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					o.installSubpanels(t, true)
					words, restore := legacy.PortTestServerPanelsObjectWords()
					t.Cleanup(restore)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					counts := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045472), 2)
					oldCounts := append([]uint32(nil), counts...)
					t.Cleanup(func() { copy(counts, oldCounts) })
					off, n, start := uintptr(33392), 27, 2
					state := memmap.PtrUint32(0x5D4594, 1045452)
					if kind == "armor" {
						off, n, start = 35496, 26, 0
						state = memmap.PtrUint32(0x5D4594, 1045456)
					}
					data := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), 12*(n+1))
					old := append([]byte(nil), data...)
					t.Cleanup(func() { copy(data, old) })
					clear(data)
					for i := 0; i < n; i++ {
						name := alloc.InternCString16(fmt.Sprintf("%s %02d", kind, i))
						binary.LittleEndian.PutUint32(data[i*12:], uint32(uintptr(unsafe.Pointer(name))))
						binary.LittleEndian.PutUint32(data[i*12+8:], 1<<uint(i))
					}
					*state = mask
					raw := legacy.PortTestServerPanelsConstruct(kind, o.options, unsafe.Pointer(&o.settings[0]))
					if raw == 0 {
						t.Fatal("object constructor")
					}
					w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
					if *words["root"] != uint32(raw) || (*words["kind"] == 1) != (kind == "armor") {
						t.Fatal("object owner")
					}
					list := (*gui.ScrollListBoxData)(w.ChildByID(1510).WidgetData)
					if int(list.Field_11_0) != n-start || int(counts[*words["kind"]]) != n-start {
						t.Fatal("object row count")
					}
					r := row{Kind: kind, Flags: flags, Mask: mask, Count: int(list.Field_11_0)}
					for i, it := range unsafe.Slice(list.Items, int(list.Field_11_0)) {
						name := alloc.GoString16(&it.Text[0])
						if name != fmt.Sprintf("%s %02d", kind, i+start) {
							t.Fatal("object title", i, name)
						}
						r.Names = append(r.Names, name)
					}
					// Existing C skips class indices 0,1,4,5 for checkbox display in both modes.
					// Record that behavior independently from row-name-based event mapping.
					index := 0
					for id := uint(1520); id <= 1533; id++ {
						for index == 0 || index == 1 || index == 4 || index == 5 {
							index++
						}
						child := w.ChildByID(id)
						hidden := child.Flags.IsHidden()
						checked := child.DrawData().Field0&4 != 0
						if hidden != (index >= n-start) || !hidden && checked != (mask&(1<<uint(index)) != 0) {
							t.Fatalf("object checkbox %d index %d", id, index)
						}
						r.Hidden = append(r.Hidden, hidden)
						r.Checked = append(r.Checked, checked)
						index++
					}
					if w.ChildByID(1515).Flags.IsEnabled() != (flags&1 != 0 && flags&0xc000 == 0) {
						t.Fatal("host controls")
					}
					for _, id := range []uint{1516, 1515, 1520, 1521, 1533} {
						child := w.ChildByID(id)
						before := *state
						wasChecked := child.DrawData().Field0&4 != 0
						*o.optionWords["dirty"] = 0
						result := gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(unsafe.Pointer(child)), 0)))
						want := uint32(0)
						switch id {
						case 1516:
							want = 0
						case 1515:
							want = 0xffffffff
						default:
							bit := uint32(1) << uint(start+int(id)-1520)
							want = before | bit
							if wasChecked {
								want = before &^ bit
							}
						}
						if result != 0 || *state != want || *o.optionWords["dirty"] != 1 {
							t.Fatalf("object event %d result %d mask %x want %x", id, result, *state, want)
						}
						r.After = append(r.After, *state)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "server-panels-object-list", rows, "c954add5fcd6c7106eb9c2f0b12490135fd8679b6f5633a674b73b52eed5cae8")
}
