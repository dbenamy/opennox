//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerStateMinimap(t *testing.T) {
	o := newMatchRosterOwner(t)
	// Own real allocated circular lists and release through the normal unmark owner.
	t.Cleanup(func() {
		for i := 0; i < 32; i++ {
			for j := range o.units {
				o.s.Players.Nox_xxx_netUnmarkMinimapObj_417300(ntype.PlayerInd(i), &o.units[j], ^uint32(0))
			}
		}
	})
	type entry struct {
		Player         int
		Objects, Flags []uint32
	}
	type row struct {
		Name        string
		Lists       []entry
		ObjectFlags [3]uint32
	}
	var rows []row
	capture := func(name string) {
		r := row{Name: name}
		for _, pl := range o.s.Players.List() {
			e := entry{Player: int(pl.PlayerInd)}
			head := pl.Field4580
			for m := head; m != nil; {
				if m.Field8 == nil || m.Field12 == nil || m.Field8.Field12 != m || m.Field12.Field8 != m {
					t.Fatal("broken minimap list")
				}
				if legacy.PortTestPlayerStateMinimap("tracks", m.Field4, int(pl.PlayerInd), 0) != 1 {
					t.Fatal("missing tracked object")
				}
				e.Objects = append(e.Objects, m.Field4.NetCode)
				e.Flags = append(e.Flags, m.Field0)
				if len(e.Objects) > 3 {
					t.Fatal("minimap cycle")
				}
				m = m.Field8
				if m == head {
					break
				}
			}
			r.Lists = append(r.Lists, e)
		}
		for i := range o.units {
			r.ObjectFlags[i] = o.units[i].Field5
		}
		rows = append(rows, r)
	}
	for _, i := range []int{-1, 0, 1, 31, 32, 0x7fffffff} {
		for _, obj := range []*server.Object{nil, &o.units[0]} {
			if legacy.PortTestPlayerStateMinimap("tracks", obj, i, 0) != 0 {
				t.Fatal("empty/invalid query", i)
			}
		}
	}
	for i := range o.units {
		legacy.PortTestPlayerStateMinimap("mark", &o.units[i], 0, uint32(1)<<i)
		for _, p := range o.s.Players.List() {
			if legacy.PortTestPlayerStateMinimap("tracks", &o.units[i], int(p.PlayerInd), 0) != 1 {
				t.Fatal("fanout missed player", p.PlayerInd)
			}
		}
		if o.s.Players.ByIndRaw(0).Field4580 != nil {
			t.Fatal("fanout touched inactive player")
		}
		capture(fmt.Sprintf("mark%d", i))
	}
	legacy.PortTestPlayerStateMinimap("mark", &o.units[1], 0, 8)
	for _, p := range o.s.Players.List() {
		n := 0
		for m := p.Field4580; ; m = m.Field8 {
			n++
			if m.Field4 == &o.units[1] && m.Field0 != 10 {
				t.Fatal("merged flags", m.Field0)
			}
			if m.Field8 == p.Field4580 {
				break
			}
		}
		if n != 3 {
			t.Fatal("duplicate tracking entry", n)
		}
	}
	capture("merge")
	legacy.PortTestPlayerStateMinimap("unmark", &o.units[1], 0, 2)
	for _, p := range o.s.Players.List() {
		if legacy.PortTestPlayerStateMinimap("tracks", &o.units[1], int(p.PlayerInd), 0) != 1 {
			t.Fatal("partial unmark removed object")
		}
	}
	capture("partial")
	for _, i := range []int{1, 2, 0} {
		legacy.PortTestPlayerStateMinimap("unmark", &o.units[i], 0, ^uint32(0))
		for _, p := range o.s.Players.List() {
			if legacy.PortTestPlayerStateMinimap("tracks", &o.units[i], int(p.PlayerInd), 0) != 0 {
				t.Fatal("unmark missed player")
			}
		}
		capture(fmt.Sprintf("unmark%d", i))
	}
	for _, p := range o.s.Players.List() {
		if p.Field4580 != nil {
			t.Fatal("list not empty")
		}
	}
	legacy.PortTestPlayerStateMinimap("unmark", nil, 0, ^uint32(0))
	legacy.PortTestPlayerStateMinimap("mark", nil, 0, 1)
	capture("nil")
	spellbookCapture(t, "player-state-minimap", rows, "26c572bbbcda004d4055c7b2ff2dc1fdb055c5ecc0eeadb9784cb3d1941e3470")
}
