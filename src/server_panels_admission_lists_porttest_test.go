//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func serverPanelsListNames(w *gui.Window) []string {
	data := (*gui.ScrollListBoxData)(w.WidgetData)
	var out []string
	for _, it := range unsafe.Slice(data.Items, int(data.Field_11_0)) {
		out = append(out, alloc.GoString16(&it.Text[0]))
	}
	return out
}
func serverPanelsSetText(w *gui.Window, text string) {
	w.Func94(gui.AsWindowEvent(16414, uintptr(unsafe.Pointer(alloc.InternCString16(text))), 0))
}
func serverPanelsGetText(w *gui.Window) string {
	p := gui.EventRespInt(w.Func94(gui.AsWindowEvent(16413, 0, 0)))
	return alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(uint32(p)))))
}

func TestServerPanelsAdmissionLists(t *testing.T) {
	type row struct {
		Allowed                   bool
		Names, After              []string
		Entry                     string
		AddEnabled, RemoveEnabled bool
	}
	var rows []row
	for _, allowed := range []bool{false, true} {
		t.Run(fmt.Sprint(allowed), func(t *testing.T) {
			o := newServerOptionsOwner(t)
			o.installSubpanels(t, true)
			defer noxflags.PortTestGameFlags(1)()
			heads := make(map[uintptr][]uint32)
			for _, off := range []uintptr{371364, 371500} {
				h := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, off)), 3)
				old := append([]uint32(nil), h...)
				t.Cleanup(func() { copy(h, old) })
				addr := uint32(uintptr(unsafe.Pointer(&h[0])))
				h[0], h[1], h[2] = addr, addr, addr
				heads[off] = h
			}
			raw := legacy.PortTestServerPanelsConstruct("access", o.options, unsafe.Pointer(&o.settings[0]))
			w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
			if w == nil {
				t.Fatal("admission constructor")
			}
			event := func(id uint) { w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(id).C()), 0)) }
			listID, radio, head := uint(10105), uint(10207), heads[371500]
			if allowed {
				listID, radio, head = 10109, 10206, heads[371364]
			}
			list := w.ChildByID(listID)
			// Both lists own C-allocated nodes; use their real remove event before
			// restoring sentinel storage, including when an assertion aborts the test.
			t.Cleanup(func() {
				if allowed {
					w.ChildByID(10102).DrawData().Field0 |= 4
				} else {
					w.ChildByID(10102).DrawData().Field0 &^= 4
				}
				for n := 0; head[0] != head[2] && n < 16; n++ {
					list.Func94(gui.AsWindowEvent(16403, 0, 0))
					event(10113)
				}
				if head[0] != head[2] {
					t.Error("admission list cleanup incomplete")
				}
			})
			event(radio)
			if list.Flags.IsHidden() {
				t.Fatal("admission radio selection")
			}
			if allowed {
				w.ChildByID(10102).DrawData().Field0 |= 4
			} else {
				w.ChildByID(10102).DrawData().Field0 &^= 4
			}
			for _, name := range []string{"Player One", "Éowyn", "third"} {
				serverPanelsSetText(w.ChildByID(10111), name)
				event(10112)
			}
			want := []string{"*Player One", "*Éowyn", "*third"}
			if allowed {
				want = []string{"Player One", "Éowyn", "third"}
			}
			names := serverPanelsListNames(list)
			if !reflect.DeepEqual(names, want) {
				t.Fatalf("admission names %v want %v", names, want)
			}
			if serverPanelsGetText(w.ChildByID(10111)) != "" {
				t.Fatal("add did not clear entry")
			}
			list.Func94(gui.AsWindowEvent(16403, 1, 0))
			w.Draw()
			if !w.ChildByID(10113).Flags.IsEnabled() {
				t.Fatal("selected remove disabled")
			}
			event(10113)
			after := serverPanelsListNames(list)
			if !reflect.DeepEqual(after, []string{want[0], want[2]}) {
				t.Fatal("admission remove", after)
			}
			serverPanelsSetText(w.ChildByID(10111), "Next")
			w.Draw()
			if !w.ChildByID(10112).Flags.IsEnabled() {
				t.Fatal("nonempty add disabled")
			}
			rows = append(rows, row{allowed, names, after, serverPanelsGetText(w.ChildByID(10111)), w.ChildByID(10112).Flags.IsEnabled(), w.ChildByID(10113).Flags.IsEnabled()})
		})
	}
	spellbookCapture(t, "server-panels-admission-lists", rows, "112ed82e59995cf9db85846de2d3481a167150e4d36e234938046d83a8853dd2")
}

func TestServerPanelsPlayerSelection(t *testing.T) {
	type row struct {
		Selection []int
		Allowed   bool
		Buttons   []bool
		Lookup    []int
		Names     []string
	}
	var rows []row
	o := newServerOptionsOwner(t)
	o.installSubpanels(t, true)
	for i := range o.players {
		o.players[i].Active = 0
	}
	o.players[31].Active = 1
	o.players[31].PlayerInd = 31
	o.players[31].SetName("Local")
	o.players[7].Active = 1
	o.players[7].PlayerInd = 7
	o.players[7].SetName("Remote")
	raw := legacy.PortTestServerPanelsConstruct("access", o.options, unsafe.Pointer(&o.settings[0]))
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
	if w == nil {
		t.Fatal("player selection root")
	}
	list := w.ChildByID(10200)
	list.Func94(gui.AsWindowEvent(16399, 0, 0))
	for _, name := range []string{"Missing", "Local", "Remote", "remote"} {
		legacy.Sub_455920(alloc.InternCString16(name))
	}
	for _, selection := range [][]int{nil, {0}, {1}, {2}, {0, 1}, {0, 1, 2}, {3}} {
		list.Func94(gui.AsWindowEvent(16403, uintptr(uint32(0xffffffff)), 0))
		for _, i := range selection {
			list.Func94(gui.AsWindowEvent(16405, uintptr(i), 0))
		}
		want := false
		for _, i := range selection {
			if i == 2 || i == 3 {
				want = true
			}
		}
		got := legacy.PortTestServerPanelsAccessSelected()
		if got != want {
			t.Fatalf("selected %v got %t want %t", selection, got, want)
		}
		w.Func94(gui.AsWindowEvent(16400, uintptr(list.C()), 0))
		r := row{Selection: selection, Allowed: got, Names: serverPanelsListNames(list)}
		for _, id := range []uint{10191, 10192} {
			r.Buttons = append(r.Buttons, w.ChildByID(id).Flags.IsEnabled())
			if w.ChildByID(id).Flags.IsEnabled() != want {
				t.Fatal("player action state")
			}
		}
		for _, name := range []string{"Local", "Remote", "remote", "REMOTE", "Absent"} {
			r.Lookup = append(r.Lookup, legacy.PortTestServerPanelsAccessLookup(name))
		}
		if !reflect.DeepEqual(r.Lookup, []int{1, 2, 3, -1, -1}) {
			t.Fatal("player lookup", r.Lookup)
		}
		rows = append(rows, r)
	}
	// Actual notification removes the exact displayed name and repairs selection.
	legacy.Sub_455950(alloc.InternCString16("Remote"))
	if !reflect.DeepEqual(serverPanelsListNames(list), []string{"Missing", "Local", "remote"}) {
		t.Fatal("player departure")
	}
	legacy.Sub_455950(alloc.InternCString16("Absent"))
	legacy.PortTestServerPanelsAccessClose(true)
	if *o.optionWords["panel-1045516"] != 0 {
		t.Fatal("access close")
	}
	legacy.Sub_455920(alloc.InternCString16("After close"))
	legacy.Sub_455950(alloc.InternCString16("After close"))
	spellbookCapture(t, "server-panels-player-selection", rows, "9af5caa5948b059754d81b476943ce57e694faef21fc15dec053676a4ef571d0")
}
