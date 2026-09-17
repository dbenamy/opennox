//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerPanelsAdvancedTabs(t *testing.T) {
	type row struct {
		Host     bool
		Repeat   int
		Tabs     []uint32
		Children [][]uint
		Loads    []string
		Settings []byte
	}
	var rows []row
	for _, host := range []bool{false, true} {
		t.Run(fmt.Sprint(host), func(t *testing.T) {
			o := newServerOptionsOwner(t)
			loads := o.installSubpanels(t, true)
			serverOptionsRulesOwner(t, o)
			t.Cleanup(legacy.PortTestServerPanelsCallbackTable())
			_, restore := legacy.PortTestServerPanelsObjectWords()
			t.Cleanup(restore)
			counts := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045472), 2)
			old := append([]uint32(nil), counts...)
			t.Cleanup(func() { copy(counts, old) })
			flags := noxflags.GameFlag(0)
			if host {
				flags = 1
			}
			defer noxflags.PortTestGameFlags(flags)()
			for repeat := 0; repeat < 2; repeat++ {
				for i := 24; i < 52; i++ {
					o.settings[i] = byte(i + repeat)
				}
				legacy.PortTestServerPanelsConstruct("advanced", o.options, unsafe.Pointer(&o.settings[0]))
				rootWord := o.optionWords["panel-1316708"]
				childWord := o.optionWords["panel-1316712"]
				tabWord := o.optionWords["panel-1316704"]
				if *rootWord == 0 || *childWord == 0 {
					t.Fatal("advanced owners")
				}
				w := (*gui.Window)(unsafe.Pointer(uintptr(*rootWord)))
				wantTab := uint32(1)
				if host {
					wantTab = 0
				}
				if *tabWord != wantTab {
					t.Fatal("initial tab")
				}
				r := row{Host: host, Repeat: repeat, Tabs: []uint32{*tabWord}}
				for _, id := range []uint{10165, 10166, 10164, 10167, 10164, 10165, 10166} {
					oldChild := *childWord
					child := w.ChildByID(id)
					if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(child.C()), 0))) != 1 {
						t.Fatal("tab return")
					}
					want := uint32(id - 10163)
					if id == 10167 {
						want = 0
					}
					if *tabWord != want || *childWord == 0 || *childWord == oldChild {
						t.Fatal("tab replacement", id)
					}
					sub := (*gui.Window)(unsafe.Pointer(uintptr(*childWord)))
					wantID := uint(1500)
					if id == 10164 {
						wantID = 1100
					}
					if id == 10167 {
						wantID = 10169
					}
					if sub.ID() != wantID {
						t.Fatal("tab resource", sub.ID(), wantID)
					}
					r.Tabs = append(r.Tabs, *tabWord)
					r.Children = append(r.Children, serverOptionsChildren(t, sub))
					for i := 24; i < 52; i++ {
						o.settings[i] = byte(i + int(id) + repeat)
					}
					legacy.PortTestServerPanelsAdvancedUpdate(unsafe.Pointer(&o.settings[0]))
					updatedMasks := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1045488), 20)
					if !bytes.Equal(updatedMasks, o.settings[24:44]) || *memmap.PtrUint32(0x5D4594, 1045452) != binary.LittleEndian.Uint32(o.settings[44:]) || *memmap.PtrUint32(0x5D4594, 1045456) != binary.LittleEndian.Uint32(o.settings[48:]) {
						t.Fatal("advanced live update")
					}

				}
				masks := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045488), 5)
				want := append([]byte(nil), o.settings...)
				for i := range masks {
					masks[i] = uint32(0x10203040 + i + repeat)
					binary.LittleEndian.PutUint32(want[24+i*4:], masks[i])
				}
				*memmap.PtrUint32(0x5D4594, 1045452) = 0x98765432
				*memmap.PtrUint32(0x5D4594, 1045456) = 0x76543210
				binary.LittleEndian.PutUint32(want[44:], 0x98765432)
				binary.LittleEndian.PutUint32(want[48:], 0x76543210)
				if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(10148).C()), 0))) != 1 {
					t.Fatal("close return")
				}
				if *rootWord != 0 || *childWord != 0 || *o.optionWords["root"] == 0 || !bytes.Equal(o.settings, want) {
					t.Fatal("advanced close ownership/settings")
				}
				o.c.GUI.FreeDestroyed()
				r.Loads = append([]string(nil), (*loads)...)
				r.Settings = append([]byte(nil), o.settings...)
				rows = append(rows, r)
			}
		})
	}
	spellbookCapture(t, "server-panels-advanced-tabs", rows, "4c46c79926c3247a9ad9b80c1a2d33437c7861e23f6ce5f5112d4f5c35e3a230")
}

func TestServerPanelsAdvancedServer(t *testing.T) {
	type row struct {
		Mode    uint32
		Initial uint32
		ID      uint
		State   []byte
		Enabled bool
		Text    string
	}
	var rows []row
	for _, mode := range []uint32{0, 1, 2, 3, 4, 0xffffffff} {
		for _, initial := range []uint32{0, 1, 2, 0xffffffff} {
			t.Run(fmt.Sprintf("%x-%x", mode, initial), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				o.installSubpanels(t, true)
				state := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371516), 100)
				binary.LittleEndian.PutUint32(state[58:], initial)
				binary.LittleEndian.PutUint32(state[62:], initial)
				binary.LittleEndian.PutUint32(state[66:], mode)
				binary.LittleEndian.PutUint32(state[70:], 12345)
				notify := memmap.PtrUint32(0x5D4594, 3588)
				oldNotify := *notify
				t.Cleanup(func() { *notify = oldNotify })
				legacy.PortTestServerPanelsConstruct("advserv", o.options, unsafe.Pointer(&o.settings[0]))
				rootWord := o.optionWords["panel-1316972"]
				if *rootWord == 0 {
					t.Fatal("advanced server root")
				}
				w := (*gui.Window)(unsafe.Pointer(uintptr(*rootWord)))
				for _, id := range []uint{2102, 2103} {
					if (w.ChildByID(id).DrawData().Field0&4 != 0) != (initial != 0) {
						t.Fatal("initial advanced checkbox")
					}
				}
				if w.ChildByID(2110).Flags.IsEnabled() != (mode >= 3) {
					t.Fatal("initial advanced edit state")
				}
				for _, id := range []uint{2102, 2103, 2106, 2109, 2107, 2108} {
					before := append([]byte(nil), state...)
					*notify = 0
					if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(id).C()), 0))) != 1 {
						t.Fatal("advanced event result")
					}
					if id == 2102 || id == 2103 {
						off := 58 + int(id-2102)*4
						binary.LittleEndian.PutUint32(before[off:], binary.LittleEndian.Uint32(before[off:])^1)
					} else {
						binary.LittleEndian.PutUint32(before[66:], uint32(id-2106))
						if *notify != 1 || w.ChildByID(2110).Flags.IsEnabled() != (id == 2109) {
							t.Fatal("advanced rate state")
						}
					}
					if !bytes.Equal(state, before) {
						t.Fatal("advanced settings delta")
					}
					rows = append(rows, row{mode, initial, id, append([]byte(nil), state...), w.ChildByID(2110).Flags.IsEnabled(), ""})
				}
				w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(2130).C()), 0))
				if *rootWord != 0 || *o.optionWords["root"] == 0 {
					t.Fatal("advanced server close ownership")
				}
			})
		}
	}
	spellbookCapture(t, "server-panels-advanced-server", rows, "ac5f1c02738aac5fb551acf78d442b22e9f2e19380833753e9e6d9c4d73b3abe")
}
