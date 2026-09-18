//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
	"unsafe"
)

func TestSessionEntryShadowList(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(legacy.PortTestSessionEntryShadowOwner())
	type state struct {
		Flags      uint32
		Next, Prev int
	}
	type row struct {
		Name  string
		Head  int
		Units [3]state
	}
	var rows []row
	ids := map[uint32]int{0: 0}
	for i := range o.units {
		ids[uint32(uintptr(unsafe.Pointer(&o.units[i])))] = i + 1
	}
	snapshot := func(name string) row {
		r := row{Name: name, Head: ids[uint32(uintptr(unsafe.Pointer(legacy.PortTestSessionEntryShadow("head", nil))))]}
		for i := range o.units {
			u := &o.units[i]
			next, ok := ids[u.Field117]
			if !ok {
				t.Fatal("unknown next")
			}
			prev, ok := ids[u.Field118]
			if !ok {
				t.Fatal("unknown prev")
			}
			r.Units[i] = state{uint32(u.ObjFlags), next, prev}
		}
		return r
	}
	for _, flags := range []uint32{0, 0x20, 0x10000, 0x400000, 0x410000, 0xffffffff} {
		restore := legacy.PortTestSessionEntryShadowOwner()
		for i := range o.units {
			o.units[i].ObjFlags = 0
			o.units[i].Field117 = 0
			o.units[i].Field118 = 0
		}
		u := &o.units[0]
		u.ObjFlags = object.Flags(flags)
		if legacy.PortTestSessionEntryShadow("add", u) != u {
			t.Fatal("add identity")
		}
		r := snapshot(fmt.Sprintf("flags%x", flags))
		want := flags
		head := 0
		if flags&0x410000 == 0 {
			want |= 0x10000
			head = 1
		}
		if r.Head != head || r.Units[0] != (state{Flags: want}) {
			t.Fatal("shadow eligibility", r)
		}
		rows = append(rows, r)
		restore()
	}
	for i := range o.units {
		o.units[i].ObjFlags = 0x20
		o.units[i].Field117 = 0
		o.units[i].Field118 = 0
	}
	for _, c := range []struct {
		op      string
		index   int
		head    int
		links   [3][2]int
		enabled [3]bool
	}{
		{"add", 0, 1, [3][2]int{}, [3]bool{true, false, false}},
		{"add", 1, 2, [3][2]int{{0, 2}, {1, 0}, {0, 0}}, [3]bool{true, true, false}},
		{"add", 2, 3, [3][2]int{{0, 2}, {1, 3}, {2, 0}}, [3]bool{true, true, true}},
		{"add", 1, 3, [3][2]int{{0, 2}, {1, 3}, {2, 0}}, [3]bool{true, true, true}},
		{"remove", 1, 3, [3][2]int{{0, 3}, {1, 3}, {1, 0}}, [3]bool{true, false, true}},
		{"remove", 2, 1, [3][2]int{{0, 0}, {1, 3}, {1, 0}}, [3]bool{true, false, false}},
		{"remove", 0, 0, [3][2]int{{0, 0}, {1, 3}, {1, 0}}, [3]bool{}},
		{"remove", 1, 0, [3][2]int{{0, 0}, {1, 3}, {1, 0}}, [3]bool{}},
		{"add", 1, 2, [3][2]int{{0, 0}, {0, 0}, {1, 0}}, [3]bool{false, true, false}},
	} {
		u := &o.units[c.index]
		if legacy.PortTestSessionEntryShadow(c.op, u) != u {
			t.Fatal("list identity")
		}
		r := snapshot(fmt.Sprintf("%s%d", c.op, c.index))
		want := row{Name: r.Name, Head: c.head}
		for i := range want.Units {
			flags := uint32(0x20)
			if c.enabled[i] {
				flags |= 0x10000
			}
			want.Units[i] = state{flags, c.links[i][0], c.links[i][1]}
		}
		if !reflect.DeepEqual(r, want) {
			t.Fatal("shadow linkage", r, want)
		}
		rows = append(rows, r)
	}
	spellbookCapture(t, "session-entry-shadow-list", rows, "633148d06e5121c4678beb76ca6860369f5d4f88fc7c679566d848032d9ff271")
}
