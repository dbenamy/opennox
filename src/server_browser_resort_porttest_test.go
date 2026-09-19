//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestServerBrowserResortSelection(t *testing.T) {
	o := newListboxOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	defer legacy.PortTestServerBrowserCollectionOwner()()
	*words["dword_587000_87408"] = 1
	clear(serverConfigOwnBytes(t, 0x5D4594, 815120, 1))
	clear(serverConfigOwnBytes(t, 0x587000, 91164, 2))
	set, restoreStrings := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "noxworld.c:Arena", Vals: []strman.Variant{{Str: "Arena"}}}, strman.Entry{ID: "noxworld.c:Open", Vals: []strman.Variant{{Str: "Open"}}})
	defer restoreStrings()
	set(0)
	var lists []*gui.Window
	for _, name := range []string{"nox_wol_wnd_gameList_815012", "dword_5d4594_815016", "dword_5d4594_815020", "dword_5d4594_815024", "dword_5d4594_815028", "dword_5d4594_815032"} {
		o.create(t, 8, 180, 100, nil)
		lists = append(lists, o.win)
		*words[name] = uint32(uintptr(o.win.C()))
		o.win = nil
	}
	raw, free := alloc.Make([]byte{}, 172)
	defer free()
	type row struct {
		Sort     uint32
		Order    []uint32
		Text     [6][]string
		Selected []byte
	}
	var rows []row
	for mode := uint32(0); mode < 10; mode++ {
		*words["nox_wol_servers_sorting_166704"] = 0
		for id, name := range []string{"Beta", "alpha", "ALPHA"} {
			clear(raw)
			copy(raw[12:], "127.0.0.1")
			copy(raw[120:135], name)
			binary.LittleEndian.PutUint32(raw[36:], uint32(id+1))
			raw[103] = byte(id)
			raw[104] = 8
			binary.LittleEndian.PutUint32(raw[96:], uint32(30-id*10))
			legacy.PortTestServerBrowserCollectionAdd(unsafe.Pointer(&raw[0]))
		}
		selected := legacy.PortTestServerBrowserCollectionID(2)
		*words["dword_5d4594_814624"] = uint32(uintptr(selected))
		before := append([]byte(nil), unsafe.Slice((*byte)(selected), 169)[12:]...)
		*words["nox_wol_servers_sorting_166704"] = mode
		cleanup := legacy.PortTestServerBrowserCollectionResort()
		r := row{Sort: mode}
		for _, p := range legacy.PortTestServerBrowserCollectionSnapshot() {
			r.Order = append(r.Order, binary.LittleEndian.Uint32(unsafe.Slice((*byte)(p), 169)[36:]))
		}
		expected := [][]uint32{{2, 3, 1}, {1, 2, 3}, {1, 2, 3}, {3, 2, 1}, {1, 2, 3}, {1, 2, 3}, {3, 2, 1}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}}
		// Equal keys reverse their preceding order because insertion puts ties first.
		if !reflect.DeepEqual(r.Order, expected[mode]) {
			t.Fatalf("resort %d order %v want %v", mode, r.Order, expected[mode])
		}
		for i, w := range lists {
			r.Text[i] = serverPanelsListNames(w)
			if len(r.Text[i]) != 3 {
				t.Fatal("resort rows", mode, i, r.Text[i])
			}
		}
		p := unsafe.Pointer(uintptr(*words["dword_5d4594_814624"]))
		r.Selected = append([]byte(nil), unsafe.Slice((*byte)(p), 169)[12:]...)
		if !reflect.DeepEqual(r.Selected, before) {
			t.Fatal("resort invalidated selected server")
		}
		rows = append(rows, r)
		*words["dword_5d4594_814624"] = 0
		cleanup()
		legacy.PortTestServerBrowserCollectionClear()
		for _, w := range lists {
			w.Func94(gui.AsWindowEvent(16399, 0, 0))
		}
	}
	spellbookCapture(t, "server-browser-resort", rows, "66311dd1593b90d9818319d28a55f53385dbe218e7dcf1937884d3a879bd25a9")
}
