//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestJournalReportsAndLocalUpdates(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	for _, slot := range []int{1, 7, 31} {
		for _, op := range []int{1, 3, 6} {
			for _, exists := range []bool{false, true} {
				o.resetJournal(t)
				if exists {
					o.call(t, 0, slot, "Short", 2)
				}
				o.queueReset()
				oldHeight := memmap.Uint32(0x5D4594, 1064848)
				r := o.call(t, op, slot, "Short", 8)
				rows = append(rows, r)
				changed := op == 1 || exists
				wantPackets := 0
				if slot != 31 && changed {
					wantPackets = 1
				}
				if len(r.Packets) != wantPackets {
					t.Fatalf("journal reports slot%d op%d exists%v", slot, op, exists)
				}
				if wantPackets == 1 {
					p := r.Packets[0]
					kind := byte(1)
					if op == 3 {
						kind = 2
					} else if op == 6 {
						kind = 3
					}
					want := make([]byte, 68)
					want[0] = 213
					want[1] = kind
					copy(want[2:66], "Short")
					if kind != 2 {
						binary.LittleEndian.PutUint16(want[66:], 8)
					}
					if p.Recipient != byte(slot) || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 || !bytes.Equal(p.Data, want) {
						t.Fatal("journal report fields")
					}
				}
				if slot == 31 && changed && (op == 1 || op == 3) {
					if r.Height == oldHeight {
						t.Fatal("local add/remove did not rebuild journal height")
					}
				} else if r.Height != oldHeight {
					t.Fatal("unexpected journal height update")
				}
			}
		}
	}
	journalCapture(t, "reports-local", rows, "ba6ad1ef89ebdc22743800d577070f4e4e71ae78f594fc39647b2aada09d4679")
}
func TestJournalAllPlayers(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	for _, op := range []int{4, 7} {
		for _, mask := range []uint32{0, 1, 2, 3, 4, 5, 6, 7} {
			o.resetJournal(t)
			for i, slot := range []int{1, 7, 31} {
				if mask&(1<<i) != 0 {
					o.call(t, 0, slot, "Short", 2)
					o.call(t, 0, slot, "Short", 4)
				}
			}
			o.queueReset()
			r := o.call(t, op, 1, "Short", 8)
			rows = append(rows, r)
			if r.Return != 0 {
				t.Fatal("all-player traversal return")
			}
			wantPackets := 0
			for i, slot := range []int{1, 7, 31} {
				n := o.players[slot].Journal
				if mask&(1<<i) == 0 {
					if n != nil {
						t.Fatal("all-player call created an entry")
					}
					continue
				}
				if slot != 31 {
					wantPackets++
				}
				if op == 4 {
					if n == nil || n.Next != nil || n.Field3 != 2 {
						t.Fatal("all-player remove must affect first match only")
					}
				} else if n == nil || n.Field3 != 8 || n.Next == nil || n.Next.Field3 != 2 {
					t.Fatal("all-player update must affect first match only")
				}
			}
			if len(r.Packets) != wantPackets {
				t.Fatal("all-player notification count")
			}
		}
	}
	journalCapture(t, "all-players", rows, "83293808f7a775f4bef6ba272c0dafce7c99fae405e97f244c3c261a08f4c9cd")
}
func TestJournalNoLocalPlayer(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	o.call(t, 0, 31, "Short", 2)
	legacy.Set_dword_8531A0_2576(nil)
	before := memmap.Uint32(0x5D4594, 1064848)
	legacy.PortTestJournal(9, 0, 0, 0)
	rows = append(rows, o.snapshot(t, 9, 0, true))
	legacy.PortTestJournal(10, 10, 20, 0)
	rows = append(rows, o.snapshot(t, 10, 0, true))
	if memmap.Uint32(0x5D4594, 1064848) != before {
		t.Fatal("nil local player changed cached height")
	}
	journalCapture(t, "no-local-player", rows, "91f3cd68f7c6cfc0b70bc7c46df12136902270bb0c28a01597160d06c98257ab")
}
