//go:build porttest

package opennox

import (
	"bytes"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestJournalStorageBoundaries(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	names := []string{"", "x", strings.Repeat("A", 62), strings.Repeat("B", 63), strings.Repeat("C", 64), strings.Repeat("D", 255), "before\x00after", "Café-Ω"}
	for _, name := range names {
		for _, flags := range []uint16{0, 1, 2, 4, 8, 15, 0x8000, 0xffff} {
			o.resetJournal(t)
			r := o.call(t, 0, 1, name, flags)
			rows = append(rows, r)
			n := o.players[1].Journal
			if n == nil {
				t.Fatal("entry allocation")
			}
			want := []byte(strings.SplitN(name, "\x00", 2)[0])
			if len(want) > 63 {
				want = want[:63]
			}
			var buf [64]byte
			copy(buf[:], want)
			if n.EntryBuf != buf || n.Field3 != flags || n.Field4 != 0 || n.Next != nil || n.Prev != nil {
				t.Fatal("initial journal node fields")
			}
			original := n
			r = o.call(t, 5, 1, string(want), flags^0xffff)
			rows = append(rows, r)
			if o.players[1].Journal != original || n.Field3 != flags^0xffff {
				t.Fatal("update must preserve node and replace all flags")
			}
			r = o.call(t, 2, 1, string(want), 0)
			rows = append(rows, r)
			if r.Return != 1 || o.players[1].Journal != nil {
				t.Fatal("remove singleton")
			}
			r = o.call(t, 2, 1, string(want), 0)
			rows = append(rows, r)
			if r.Return != 0 {
				t.Fatal("missing removal")
			}
		}
	}
	journalCapture(t, "storage-boundaries", rows, "6c302597d56e95369d2bc7c8e4d80bdefa7c8857cc58227a08d4688689ccb763")
}
func TestJournalLinksAndDuplicates(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	for _, index := range []int{0, 1, 2, 3} {
		o.resetJournal(t)
		for _, name := range []string{"tail", "duplicate", "middle", "duplicate"} {
			rows = append(rows, o.call(t, 0, 7, name, 2))
		}
		head := o.players[7].Journal
		older := head.Next.Next
		rows = append(rows, o.call(t, 5, 7, "duplicate", 8))
		if head.Field3 != 8 || older.Field3 != 2 {
			t.Fatal("update must choose first duplicate")
		}
		names := []string{"duplicate", "middle", "tail", "absent"}
		before := head
		next := head.Next
		r := o.call(t, 2, 7, names[index], 0)
		rows = append(rows, r)
		if index == 3 {
			if r.Return != 0 || o.players[7].Journal != before {
				t.Fatal("missing removal changed list")
			}
		} else if r.Return != 1 {
			t.Fatal("remove live node")
		}
		if index == 0 && o.players[7].Journal != next {
			t.Fatal("remove first duplicate")
		}

	}
	journalCapture(t, "links-duplicates", rows, "2bd079ad9678663567eadf4d348e7fb6931b1b9aacdb9f2bc819071e0fc62978")
}
func TestJournalMaskRemoval(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	for _, mask := range []uint16{0, 1, 2, 3, 4, 8, 14, 15, 0x8000, 0xffff} {
		o.resetJournal(t)
		flags := []uint16{0, 1, 2, 3, 4, 8, 15, 0x8000, 0xffff}
		for _, f := range flags {
			o.call(t, 0, 1, "Short", f)
		}
		r := o.call(t, 8, 1, "", mask)
		rows = append(rows, r)
		if r.Return != 0 {
			t.Fatal("mask removal return")
		}
		var want []uint16
		for i := len(flags) - 1; i >= 0; i-- {
			if flags[i]&mask == 0 {
				want = append(want, flags[i])
			}
		}
		if len(r.Lists[0]) != len(want) {
			t.Fatal("mask removal cardinality")
		}
		for i, n := range r.Lists[0] {
			if n.Flags != want[i] {
				t.Fatal("mask removal ordering")
			}
		}
	}
	journalCapture(t, "mask-removal", rows, "176a53525d14a411c04e9f2bf90e86b394e02ba8fa3aad2235f0aa4d82b9f7d5")
}
func TestJournalSequences(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	type entry struct {
		name  string
		flags uint16
		node  *server.PlayerJournal
	}
	names := []string{"", "Short", "Wrapped", "Multiline", "Unicode", "missing"}
	for _, seed := range []uint32{1, 0x12345678, 0xffffffff} {
		o.resetJournal(t)
		var model []entry
		rng := seed
		for step := 0; step < 384; step++ {
			rng = rng*1664525 + 1013904223
			op := int(rng % 3)
			name := names[(rng>>8)%uint32(len(names))]
			flags := uint16(rng >> 16)
			switch op {
			case 0:
				o.call(t, 0, 7, name, flags)
				model = append([]entry{{name, flags, o.players[7].Journal}}, model...)
			case 1:
				o.call(t, 5, 7, name, flags)
				for i := range model {
					if model[i].name == name {
						model[i].flags = flags
						break
					}
				}
			case 2:
				o.call(t, 2, 7, name, 0)
				for i := range model {
					if model[i].name == name {
						model = append(model[:i], model[i+1:]...)
						break
					}
				}
			}
			n := o.players[7].Journal
			for _, m := range model {
				if n == nil || n != m.node || n.Field3 != m.flags || string(bytes.TrimRight(n.EntryBuf[:], "\x00")) != m.name {
					t.Fatal("journal differs from independent sequence model")
				}
				n = n.Next
			}
			if n != nil {
				t.Fatal("extra journal nodes")
			}
			if step%8 == 7 {
				rows = append(rows, o.snapshot(t, op, 0, false))
			}
		}
	}
	journalCapture(t, "sequences", rows, "41e21ac6be3d1098bffaaae129314df51413c942726741e809ea968f54eedd2b")
}
func TestJournalLifecycle(t *testing.T) {
	o := newJournalOwner(t)
	for cycle := 0; cycle < 1000; cycle++ {
		o.call(t, 0, 1, "Short", 1)
		o.call(t, 0, 1, "Wrapped", 2)
		o.call(t, 0, 1, "Multiline", 8)
		r := o.call(t, 8, 1, "", 0xffff)
		if r.Return != 0 || o.players[1].Journal != nil {
			t.Fatal("journal lifecycle left entries")
		}
	}
	// Exercise the exact empty list boundary through the production C call.
	if legacy.PortTestJournal(8, uintptr(o.units[0].CObj()), 0, 0xffff) != 0 {
		t.Fatal("empty mask cleanup")
	}
}
