//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func TestPortVotesAllocationAndRemoval(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	if p := legacy.PortTestVoteCreate(0, nil); p != nil {
		t.Fatal("nil player accepted")
	}
	var records []unsafe.Pointer
	for i, kind := range []int{0, 1, 2, 3, 17} {
		p := legacy.PortTestVoteCreate(kind, o.units[0].CObj())
		if p == nil {
			t.Fatal("allocation failed")
		}
		records = append(records, p)
		word := func(off int) uint32 { return *(*uint32)(unsafe.Add(p, off)) }
		min := byte(5)
		if kind == 2 || kind == 3 {
			min = 6
		}
		if kind == 17 {
			min = 9
		}
		if word(0) != uint32(kind) || word(24) != 123 || *(*byte)(unsafe.Add(p, 12)) != min {
			t.Fatalf("record %d initialization", kind)
		}
		if *(*unsafe.Pointer)(unsafe.Add(p, 16)) != unsafe.Add(o.units[0].CObj(), 48) {
			t.Fatal("team reference")
		}
		for _, off := range []int{4, 8, 20, 28, 32, 36, 40} {
			if word(off) != 0 {
				t.Fatalf("record field %d not zero", off)
			}
		}
		list := legacy.PortTestVoteList()
		if len(list) != i+1 || list[0] != p {
			t.Fatal("prepend")
		}
	}
	if nodes := o.state().Nodes; len(nodes) != 1 || !bytes.Equal(nodes[0].Data, []byte{238, 6, 1}) {
		t.Fatalf("first-record announcement: %+v", nodes)
	}
	// Remove middle, head, tail, then the remaining records, validating both links.
	for _, i := range []int{2, 4, 0, 3, 1} {
		legacy.PortTestVoteDelete(records[i])
		legacy.PortTestVoteList()
	}
	if len(legacy.PortTestVoteList()) != 0 {
		t.Fatal("list not empty")
	}
	nodes := o.state().Nodes
	if len(nodes) != 2 || !bytes.Equal(nodes[0].Data, []byte{238, 6, 0}) {
		t.Fatalf("last-record announcement: %+v", nodes)
	}
	// Sparse player indices include bit 31; removing one player traverses and
	// deletes multiple records without losing the saved successor.
	for _, u := range o.units {
		p := legacy.PortTestVoteCreate(0, o.units[0].CObj())
		*(*byte)(unsafe.Add(p, 4)) = 1
		*(*uint32)(unsafe.Add(p, 8)) = uint32(1) << u.UpdateDataPlayer().Player.PlayerInd
	}
	for i := range o.units {
		legacy.PortTestVoteRemovePlayer(o.units[i].CObj())
		if n := len(legacy.PortTestVoteList()); n != 2-i {
			t.Fatalf("removal %d: %d records", i, n)
		}
		legacy.PortTestVoteRemovePlayer(o.units[i].CObj())
	}
}

func TestPortVotesNameMessage(t *testing.T) {
	o := newReliableReportsOwner(t)
	for _, name := range []string{"A", "Vote player", "\u03a9"} {
		o.units[0].UpdateDataPlayer().Player.SetName(name)
		for _, withdraw := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%v", name, withdraw), func(t *testing.T) {
				o.s.NetList.ResetAll()
				text := append(utf16.Encode([]rune(name)), 0)
				legacy.PortTestVoteNameMessage(&text[0], withdraw)
				var packets [][]byte
				o.s.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { packets = append(packets, bytes.Clone(b)); return false })
				want := make([]byte, 52)
				want[0] = 238
				if withdraw {
					want[1] = 2
				}
				for i, c := range text {
					want[2+i*2] = byte(c)
					want[3+i*2] = byte(c >> 8)
				}
				if len(packets) != 1 || !bytes.Equal(packets[0], want) {
					t.Fatalf("name message: got %x; want %x", packets, want)
				}
			})
		}
	}
}

func TestPortVotesCastWithdraw(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	old := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(old)
	t.Cleanup(func() { noxflags.UnsetGamePlay(noxflags.GetGamePlay()); noxflags.SetGamePlay(old) })
	for i := range o.units {
		o.units[i].UpdateDataPlayer().Player.SetName(fmt.Sprintf("player%d", i))
	}
	call := func(kind, caster int, name string, withdraw bool) {
		text := append(utf16.Encode([]rune(name)), 0)
		legacy.PortTestVoteCast(kind, o.units[caster].CObj(), &text[0], withdraw)
	}
	for _, bad := range []string{"missing", "player0", "player2"} {
		call(0, 0, bad, false)
		if len(legacy.PortTestVoteList()) != 0 {
			t.Fatalf("accepted invalid target %s", bad)
		}
	}
	call(0, 0, "player1", false)
	call(0, 0, "player1", false)
	call(0, 2, "player1", false)
	list := legacy.PortTestVoteList()
	if len(list) != 1 {
		t.Fatal("duplicate record")
	}
	p := list[0]
	if *(*byte)(unsafe.Add(p, 4)) != 2 || *(*uint32)(unsafe.Add(p, 8)) != 0x80000002 || *(*unsafe.Pointer)(unsafe.Add(p, 28)) != o.units[1].CObj() {
		t.Fatal("voter set or target")
	}
	call(0, 1, "player1", true) // A missing voter does not alter the record.
	call(0, 0, "player1", true)
	call(0, 0, "player1", true)
	if *(*byte)(unsafe.Add(p, 4)) != 1 || *(*uint32)(unsafe.Add(p, 8)) != 0x80000000 {
		t.Fatal("withdrawal")
	}
	call(0, 2, "player1", true)
	if len(legacy.PortTestVoteList()) != 0 {
		t.Fatal("last voter retained")
	}
}

func TestPortVotesThresholds(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	p := legacy.PortTestVoteCreate(0, o.units[0].CObj())
	// Three live units: unanimity among the other two overrides the configured
	// minimum, but one vote does not. The stored minimum also permits zero.
	for _, min := range []byte{0, 1, 2, 5, 255} {
		for _, count := range []byte{0, 1, 2, 3, 254, 255} {
			*(*byte)(unsafe.Add(p, 12)) = min
			*(*byte)(unsafe.Add(p, 4)) = count
			want := 0
			if count >= min || count >= 2 {
				want = 1
			}
			if got := legacy.PortTestVoteThreshold(p); got != want {
				t.Fatalf("min=%d count=%d got=%d want=%d", min, count, got, want)
			}
		}
	}
}
