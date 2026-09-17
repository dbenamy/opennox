//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func TestWorldMotionSentryList(t *testing.T) {
	o := newCollisionCoreOwner(t)
	words, restore := legacy.PortTestWorldMotionListGlobals()
	units := []*server.Object{newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)}
	ids := collisionCoreIDs(units...)
	saved := make([]server.Object, len(units))
	records := make([]unsafe.Pointer, len(units))
	for i, u := range units {
		saved[i] = *u
		records[i] = collisionCoreGuarded(t, o, 12)
	}
	t.Cleanup(restore)
	t.Cleanup(func() {
		for i, u := range units {
			*u = saved[i]
		}
	})
	type row struct {
		Order     [3]int
		Remove    int
		Destroyed bool
		Label     string
		Return    uint32
		Head      uint32
		Units     [][6]uint32
	}
	var rows []row
	orderings := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for _, order := range orderings {
		for remove := 0; remove < 3; remove++ {
			for _, destroyed := range []bool{false, true} {
				legacy.PortTestWorldMotionList("sentry-clear", nil, 0)
				for i, u := range units {
					*u = saved[i]
					u.ObjFlags = 4
					u.UpdateData = records[i]
					u.InvNextItem = nil
					u.Field125 = nil
					*(*[3]uint32)(records[i]) = [3]uint32{math.Float32bits(1.5), math.Float32bits(2.25), math.Float32bits(0.125)}
				}
				snapshot := func(label string, rv uint32) {
					r := row{Order: order, Remove: remove, Destroyed: destroyed, Label: label, Return: rv, Head: collisionCoreID(t, ids, *words["sentry"])}
					for _, u := range units {
						data := *(*[3]uint32)(u.UpdateData)
						r.Units = append(r.Units, [6]uint32{uint32(u.ObjFlags), collisionCoreID(t, ids, uint32(uintptr(u.InvNextItem.CObj()))), collisionCoreID(t, ids, uint32(uintptr(u.Field125.CObj()))), data[0], data[1], data[2]})
					}
					rows = append(rows, r)
				}
				var expected []uint32
				for _, i := range order {
					rv := legacy.PortTestWorldMotionList("sentry-update", units[i], 0)
					if rv != 0x80000004 {
						t.Fatal("sentry registration return", rv)
					}
					expected = append([]uint32{uint32(1001 + i)}, expected...)
					snapshot("register", rv)
				}
				target := units[remove]
				var rv uint32
				if destroyed {
					target.ObjFlags |= object.Flags(0x20)
					rv = legacy.PortTestWorldMotionList("sentry-update", target, 0)
				} else {
					rv = legacy.PortTestWorldMotionList("sentry-remove", target, 0)
				}
				if rv != uint32(uintptr(target.CObj())) {
					t.Fatal("sentry removal return")
				}
				snapshot("remove", collisionCoreID(t, ids, rv))
				want := make([]uint32, 0, 2)
				for _, id := range expected {
					if id != uint32(1001+remove) {
						want = append(want, id)
					}
				}
				var got []uint32
				var previous *server.Object
				for p := *words["sentry"]; p != 0; {
					if len(got) >= 3 {
						t.Fatal("sentry list cycle")
					}
					u := (*server.Object)(unsafe.Pointer(uintptr(p)))
					got = append(got, collisionCoreID(t, ids, p))
					if u.Field125 != previous {
						t.Fatal("sentry backward link")
					}
					previous = u
					p = uint32(uintptr(u.InvNextItem.CObj()))
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("sentry removal retained wrong members: %v want %v", got, want)
				}
				if uint32(target.ObjFlags)&0x80000000 != 0 {
					t.Fatal("sentry removal retained membership flag")
				}
				legacy.PortTestWorldMotionList("sentry-remove", target, 0)
				snapshot("remove-again", uint32(1001+remove))
				legacy.PortTestWorldMotionList("sentry-clear", nil, 0)
				snapshot("clear", 0)
				if *words["sentry"] != 0 {
					t.Fatal("sentry reset retained head")
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-sentry-list", rows, "")
}
