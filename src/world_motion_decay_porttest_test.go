//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"slices"
	"sort"
	"testing"
	"unsafe"
)

func TestWorldMotionDecayQueue(t *testing.T) {
	o := newCollisionCoreOwner(t)
	units := []*server.Object{newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)}
	ids := collisionCoreIDs(units...)
	words, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	saved := make([]server.Object, len(units))
	for i, u := range units {
		saved[i] = *u
	}
	t.Cleanup(func() {
		for i, u := range units {
			*u = saved[i]
		}
	})
	type row struct {
		Frame        uint32
		Label        string
		Object       int
		Delay        int32
		Return, Head uint32
		Units        [][3]uint32
	}
	var rows []row
	for _, frame := range []uint32{0, 123, 0x7ffffff0, 0xfffffff0} {
		o.s.SetFrame(frame)
		legacy.PortTestWorldMotionList("decay-clear", nil, 0)
		for i, u := range units {
			*u = saved[i]
			u.ObjFlags = 4
			u.Field117 = 0
			u.Field34 = 0xabcdef01
		}
		var want []int
		expiry := map[int]uint32{}
		remove := func(id int) {
			for j, x := range want {
				if x == id {
					want = append(want[:j], want[j+1:]...)
					break
				}
			}
		}
		run := func(op string, id int, delay int32) {
			var u *server.Object
			if id >= 0 {
				u = units[id]
			}
			if op == "decay-set" && u.ObjFlags&0x10000 == 0 {
				remove(id)
				expiry[id] = frame + uint32(delay)
				want = append(want, id)
				sort.SliceStable(want, func(i, j int) bool { return expiry[want[i]] < expiry[want[j]] })
			}
			if op == "decay-remove" {
				remove(id)
			}
			if op == "decay-clear" {
				want = nil
			}
			rv := legacy.PortTestWorldMotionList(op, u, delay)
			if mapped, ok := ids[unsafe.Pointer(uintptr(rv))]; ok {
				rv = mapped
			}
			var got []int
			for p := *words["decay"]; p != 0; {
				if len(got) >= len(units) {
					t.Fatal("decay cycle")
				}
				v := (*server.Object)(unsafe.Pointer(uintptr(p)))
				n := int(collisionCoreID(t, ids, p)) - 1001
				got = append(got, n)
				if v.ObjFlags&0x400000 == 0 || v.Field34 != expiry[n] {
					t.Fatal("decay membership/deadline")
				}
				p = v.Field117
			}
			if !slices.Equal(got, want) {
				t.Fatalf("decay deadline order and FIFO ties %s: %v want %v", op, got, want)
			}
			r := row{Frame: frame, Label: op, Object: id, Delay: delay, Return: rv, Head: collisionCoreID(t, ids, *words["decay"])}
			for _, v := range units {
				r.Units = append(r.Units, [3]uint32{uint32(v.ObjFlags), v.Field34, collisionCoreID(t, ids, v.Field117)})
			}
			rows = append(rows, r)
		}
		run("decay-set", 0, 20)
		run("decay-set", 1, 10)
		run("decay-set", 2, 10)
		run("decay-set", 3, 30)
		run("decay-set", 0, 5)
		run("decay-remove", 2, 0)
		run("decay-remove", 0, 0)
		run("decay-remove", 3, 0)
		run("decay-remove", 1, 0)
		run("decay-remove", 1, 0)
		units[2].ObjFlags |= 0x10000
		before := *units[2]
		run("decay-set", 2, 0)
		if units[2].Field34 != before.Field34 || units[2].Field117 != before.Field117 || units[2].ObjFlags != before.ObjFlags {
			t.Fatal("no-decay object changed")
		}
		units[2].ObjFlags &^= 0x10000
		run("decay-set", 2, -1)
		run("decay-set", 0, 0)
		run("decay-set", 1, 0x7fffffff)
		run("decay-clear", -1, 0)
		if *words["decay"] != 0 {
			t.Fatal("decay clear head")
		}
		for _, u := range units {
			if u.ObjFlags&0x400000 != 0 {
				t.Fatal("clear retained decay flag")
			}
		}
	}
	spellbookCapture(t, "world-motion-decay-queue", rows, "dedfa60ccf2d8ee0881e3343309c8d40fb06a4eccf5bdd1ee6dec61c75eeff6f")
}
func TestWorldMotionDecayExpiry(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b, c, holder := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	units := []*server.Object{a, b, c}
	ids := collisionCoreIDs(a, b, c, holder)
	words, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	saved := []server.Object{*a, *b, *c}
	oldDeleted := o.s.Objs.DeletedList
	t.Cleanup(func() {
		o.s.Objs.DeletedList = oldDeleted
		for i, u := range units {
			*u = saved[i]
		}
	})
	type row struct {
		Frame, Now uint32
		Held       bool
		Head       uint32
		Deleted    []uint32
		Units      [][5]uint32
	}
	var rows []row
	for _, frame := range []uint32{123, 0xfffffff8} {
		for _, held := range []bool{false, true} {
			for _, advance := range []uint32{0, 4, 5, 6, 10, 11, 20} {
				legacy.PortTestWorldMotionList("decay-clear", nil, 0)
				o.s.Objs.DeletedList = nil
				o.s.SetFrame(frame)
				for i, u := range units {
					*u = saved[i]
					u.ObjClass = object.ClassSimple
					u.ObjFlags = 4
					u.Field117 = 0
					u.ObjOwner = nil
					u.InvHolder = nil
					u.DeletedNext = nil
					u.DeletedAt = 0
					legacy.PortTestWorldMotionList("decay-set", u, int32(5+5*i))
				}
				if held {
					b.InvHolder = holder
				}
				now := frame + advance
				o.s.SetFrame(now)
				legacy.PortTestWorldMotionList("decay-tick", nil, 0)
				r := row{Frame: frame, Now: now, Held: held, Head: collisionCoreID(t, ids, *words["decay"])}
				for u := o.s.Objs.DeletedList; u != nil; u = u.DeletedNext {
					if len(r.Deleted) >= 3 {
						t.Fatal("deleted cycle")
					}
					r.Deleted = append(r.Deleted, ids[u.CObj()])
				}
				// The original queue uses unsigned absolute deadlines, including at wrap.
				// A held item is removed only once traversal reaches it, even before expiry.
				candidates := []int{0, 1, 2}
				sort.SliceStable(candidates, func(i, j int) bool { return frame+uint32(5+5*candidates[i]) < frame+uint32(5+5*candidates[j]) })
				reached := true
				for _, i := range candidates {
					u := units[i]
					expired := frame+uint32(5+5*i) <= now
					removed := reached && (i == 1 && held || expired)
					if reached && !(i == 1 && held) && !expired {
						reached = false
					}
					deleted := removed && !(i == 1 && held)
					if (u.ObjFlags&0x20 != 0) != deleted || (u.ObjFlags&0x400000 == 0) != removed {
						t.Fatal("expiry/held-item contract", i, now, removed, deleted, uint32(u.ObjFlags))
					}
					if deleted && (u.DeletedAt != now || u.Field5&0x80 == 0) {
						t.Fatal("expiry delete timestamp/reason")
					}
				}
				for _, u := range units {
					r.Units = append(r.Units, [5]uint32{uint32(u.ObjFlags), u.Field5, u.Field34, u.DeletedAt, collisionCoreID(t, ids, u.Field117)})
				}
				rows = append(rows, r)
			}
		}
	}
	spellbookCapture(t, "world-motion-decay-expiry", rows, "8594706fa8d642fc7b04d80f71a8129cf05a56bbf817dbd47c69660f72c757de")
}
