//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMatchRosterRememberedIdentities(t *testing.T) {
	o := newMatchRosterOwner(t)
	pl := o.units[0].UpdateDataPlayer().Player
	type row struct {
		Name           string
		Query, Present uint32
		Records        [][]byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-remembered", rows, "b3e7df8ed3b67d56f4b2155435e3e66c9d535f9753cced5fcf7f731a8c7a99b2")
	}()
	for _, name := range []string{"", "a", "abcdefghi", "abcdefghijk"} {
		for _, class := range []byte{0, 1, 2, 127, 128, 255} {
			for _, group := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				matchRosterCall("forget", nil, nil, nil)
				namebuf := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
				clear(namebuf)
				copy(namebuf, name)
				*(*byte)(unsafe.Add(pl.C(), 2251)) = class
				objectXferSetWord(pl.C(), 2068, group)
				matchRosterCall("remember", nil, pl, nil)
				records := legacy.PortTestMatchRosterRemembered()
				if len(records) != 1 {
					t.Fatal("remember count", len(records))
				}
				for _, qname := range []string{name, "different"} {
					for _, qclass := range []byte{class, class + 1} {
						for _, qgroup := range []uint32{group, group + 1} {
							label := fmt.Sprintf("name%q/class%d/group%x/query%q-%d-%x", name, class, group, qname, qclass, qgroup)
							t.Run(label, func(t *testing.T) {
								query := legacy.PortTestMatchRosterQuery(qname, int(qclass), qgroup)
								want := uint32(1)
								// C promotes the signed query char and the stored unsigned byte.
								if qname == name && (int(int8(qclass)) != int(class) || qgroup != group) {
									want = 0
								}
								if query != want {
									t.Fatalf("identity admission %d want %d", query, want)
								}
								present := matchRosterCall("remembered", nil, pl, nil)
								if present != 1 {
									t.Fatal("remembered identity missing")
								}
								rows = append(rows, row{label, query, present, records})
							})
						}
					}
				}
			}
		}
	}
	matchRosterCall("forget", nil, nil, nil)
	if legacy.PortTestMatchRosterRemembered() != nil {
		t.Fatal("list did not empty")
	}
	if matchRosterCall("remembered", nil, pl, nil) != 0 {
		t.Fatal("forgotten identity remained")
	}
}

func TestMatchRosterRememberedList(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name    string
		Removed uint32
		Records [][]byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-remembered-list", rows, "e696e388a5a4198a78cc7de930d86ff97733f1966be2a54494779984472087f3")
	}()
	for _, count := range []int{0, 1, 2, 3, 8} {
		for repeat := 0; repeat < 3; repeat++ {
			matchRosterCall("forget", nil, nil, nil)
			for i := 0; i < count; i++ {
				pl := o.units[i%3].UpdateDataPlayer().Player
				b := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
				clear(b)
				copy(b, fmt.Sprintf("name%d", i%3))
				objectXferSetWord(pl.C(), 2068, uint32(i))
				*(*byte)(unsafe.Add(pl.C(), 2251)) = byte(i % 3)
				matchRosterCall("remember", nil, pl, nil)
			}
			before := legacy.PortTestMatchRosterRemembered()
			if len(before) != count {
				t.Fatal("list count", len(before), count)
			}
			removed := matchRosterCall("forget", nil, nil, nil)
			want := uint32(0)
			if count != 0 {
				want = 1
			}
			if removed != want || len(legacy.PortTestMatchRosterRemembered()) != 0 {
				t.Fatal("list cleanup")
			}
			rows = append(rows, row{fmt.Sprintf("count%d/repeat%d", count, repeat), removed, before})
		}
	}
}
