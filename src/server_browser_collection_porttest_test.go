//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerBrowserCollection(t *testing.T) {
	o := newEntryOwner(t)
	var entries []strman.Entry
	for _, name := range []string{"Quest", "CTF", "Highlander", "KotR", "Flagball", "Chat", "Arena"} {
		entries = append(entries, strman.Entry{ID: strman.ID("noxworld.c:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	set, restoreStrings := o.c.srv.Server.PortTestMeterStrings(entries...)
	defer restoreStrings()
	set(0)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	defer legacy.PortTestServerBrowserCollectionOwner()()
	type entry struct {
		ID              uint32
		Name            string
		Players, Status byte
		Ping            int32
		Mode            uint16
		ModeName        string
		Port            uint16
	}
	entriesIn := []entry{
		{1, "zeta", 0, 0, 9999, 0, "Arena", 0},
		{2, "Alpha", 1, 0x10, 0, 0x20, "CTF", 18590},
		{3, "beta", 32, 0x20, 100, 0x1000, "Quest", 32767},
		{4, "ALPHA", 255, 0x30, 100, 0x80, "Chat", 32768},
		{5, "", 9, 0xff, -1, 0x400, "Highlander", 65535},
		{6, "gamma", 3, 0x0f, 2147483647, 0x10, "KotR", 1},
		{7, "alpha", 32, 0x40, -2147483648, 0x40, "Flagball", 18590},
	}
	type row struct {
		Sort     int
		Insert   uint32
		Return   int
		Order    []uint32
		Keys     []int32
		Payloads [][]byte
	}
	var rows []row
	input, free := alloc.Make([]byte{}, 172)
	defer free()
	for mode := 0; mode < 10; mode++ {
		legacy.PortTestServerBrowserCollectionClear()
		*words["nox_wol_servers_sorting_166704"] = uint32(mode)
		var order []entry
		key := func(e entry) int32 {
			switch mode {
			case 2:
				return int32(e.Players)
			case 3:
				return 32 - int32(e.Players)
			case 6:
				return e.Ping
			case 7:
				return 1000 - e.Ping
			case 8:
				return int32(e.Status & 0x30)
			case 9:
				return 48 - int32(e.Status&0x30)
			}
			return -123
		}
		cmp := func(a, b entry) int {
			switch mode {
			case 0:
				return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
			case 1:
				return strings.Compare(strings.ToLower(b.Name), strings.ToLower(a.Name))
			case 4:
				return strings.Compare(a.ModeName, b.ModeName)
			case 5:
				return strings.Compare(b.ModeName, a.ModeName)
			}
			ka, kb := key(a), key(b)
			if ka < kb {
				return -1
			}
			if ka > kb {
				return 1
			}
			return 0
		}
		for _, e := range entriesIn {
			clear(input)
			binary.LittleEndian.PutUint32(input[8:], uint32(0xffffff85))
			copy(input[12:28], "127.0.0.1")
			binary.LittleEndian.PutUint32(input[36:], e.ID)
			binary.LittleEndian.PutUint32(input[96:], uint32(e.Ping))
			input[100] = e.Status
			input[103] = e.Players
			input[104] = 32
			binary.LittleEndian.PutUint16(input[109:], e.Port)
			copy(input[111:120], "arena")
			copy(input[120:135], e.Name)
			binary.LittleEndian.PutUint16(input[163:], e.Mode)
			before := append([]byte(nil), input...)
			got := legacy.PortTestServerBrowserCollectionAdd(unsafe.Pointer(&input[0]))
			at := sort.Search(len(order), func(i int) bool { return cmp(e, order[i]) <= 0 })
			order = append(order, entry{})
			copy(order[at+1:], order[at:])
			order[at] = e
			if got != at || !reflect.DeepEqual(input, before) {
				t.Fatal("insert index/input mutation", mode, e.ID, got, at)
			}
			r := row{Sort: mode, Insert: e.ID, Return: got}
			actual := legacy.PortTestServerBrowserCollectionSnapshot()
			if len(actual) != len(order) {
				t.Fatal("list size", mode, len(actual), len(order))
			}
			for i, p := range actual {
				raw := unsafe.Slice((*byte)(p), 169)
				id := binary.LittleEndian.Uint32(raw[36:])
				sortKey := int32(binary.LittleEndian.Uint32(raw[8:]))
				if id != order[i].ID || sortKey != key(order[i]) {
					t.Fatal("list order/key", mode, i, id, order[i].ID, sortKey, key(order[i]))
				}
				if legacy.PortTestServerBrowserCollectionAt(int32(i)) != p || legacy.PortTestServerBrowserCollectionID(int32(id)) != p {
					t.Fatal("lookup identity", mode, i, id)
				}
				if legacy.PortTestServerBrowserCollectionMissing("127.0.0.1", order[i].Port) != func() int {
					if order[i].Port >= 32768 {
						return 1
					}
					return 0
				}() {
					t.Fatal("address/port lookup", mode, order[i].Port)
				}
				r.Order = append(r.Order, id)
				r.Keys = append(r.Keys, sortKey)
				r.Payloads = append(r.Payloads, append([]byte(nil), raw[12:]...))
			}
			if legacy.PortTestServerBrowserCollectionAt(-1) != nil || legacy.PortTestServerBrowserCollectionAt(int32(len(order))) != nil || legacy.PortTestServerBrowserCollectionID(12345) != nil || legacy.PortTestServerBrowserCollectionMissing("127.0.0.2", 18590) != 1 {
				t.Fatal("missing lookup")
			}
			rows = append(rows, r)
		}
	}
	legacy.PortTestServerBrowserCollectionClear()
	if len(legacy.PortTestServerBrowserCollectionSnapshot()) != 0 || legacy.PortTestServerBrowserCollectionAt(0) != nil {
		t.Fatal("clear")
	}
	spellbookCapture(t, "server-browser-collection", rows, "83c721cd211500d5bc8b586fdc192ccd2ec96052d7fae124f70b73a798f18bb0")
}
