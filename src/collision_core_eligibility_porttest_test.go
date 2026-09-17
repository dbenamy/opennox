//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestCollisionCoreEligibility(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	savedA, savedB := *a, *b
	ad, bd := collisionCoreGuarded(t, o, 2200), collisionCoreGuarded(t, o, 2200)
	t.Cleanup(func() { *a, *b = savedA, savedB })
	reset := func() {
		*a, *b = savedA, savedB
		a.UpdateData = ad
		b.UpdateData = bd
		clear(unsafe.Slice((*byte)(ad), 2200))
		clear(unsafe.Slice((*byte)(bd), 2200))
		a.NetCode = 1001
		b.NetCode = 1002
		a.Collide = o.callback
		b.Collide = o.callback
		a.ObjClass = object.ClassSimple
		b.ObjClass = object.ClassSimple
		a.ObjFlags = 0
		b.ObjFlags = 0
		a.TeamVal = server.ObjectTeam{}
		b.TeamVal = server.ObjectTeam{}
	}
	type row struct {
		Label               string
		Class, Flags, Types [2]uint32
		Result              int32
		Cache               [3]uint32
	}
	var rows []row
	call := func(label string, target *server.Object) int32 {
		rv := legacy.PortTestCollisionCore("eligible", a, target, nil, 0)
		rows = append(rows, row{label, [2]uint32{uint32(a.ObjClass), uint32(target.ObjClass)}, [2]uint32{uint32(a.ObjFlags), uint32(target.ObjFlags)}, [2]uint32{uint32(a.TypeInd), uint32(target.TypeInd)}, rv, [3]uint32{*o.words["trigger"], *o.words["powder"], *o.words["hand"]}})
		return rv
	}
	classes := []uint32{8, 2, 4, 0x80, 0x2000, 0x4000, 0x8000, 0x400000}
	flags := []uint32{0, 1, 4, 8, 0x10, 0x20, 0x40, 0x60, 0x400, 0x2000, 0x4000, 0x24000}
	for _, ac := range classes {
		for _, bc := range classes {
			for _, af := range flags {
				for _, bf := range flags {
					reset()
					a.ObjClass = object.Class(ac)
					b.ObjClass = object.Class(bc)
					a.ObjFlags = object.Flags(af)
					b.ObjFlags = object.Flags(bf)
					call("matrix", b)
				}
			}
		}
	}
	reset()
	for _, p := range o.words {
		*p = 0
	}
	if call("cold", b) != 1 {
		t.Fatal("ordinary pair rejected")
	}
	if *o.words["trigger"] == 0 || *o.words["powder"] == 0 || *o.words["hand"] == 0 {
		t.Fatal("lazy lookup not initialized")
	}
	if call("warm", b) != 1 {
		t.Fatal("warm ordinary pair rejected")
	}
	if call("self", a) != 0 {
		t.Fatal("self collision accepted")
	}
	for _, missing := range []int{0, 1, 2} {
		reset()
		if missing != 1 {
			a.Collide = nil
		}
		if missing != 0 {
			b.Collide = nil
		}
		if call("missing-callback", b) != 0 {
			t.Fatal("missing callback accepted")
		}
	}
	for _, swap := range []bool{false, true} {
		reset()
		door, hand := a, b
		if swap {
			door, hand = b, a
		}
		door.Collide = legacy.PortTestCollisionCorePentagram()
		hand.TypeInd = uint16(*o.words["hand"])
		if call("pentagram-hand", b) != 0 {
			t.Fatal("telekinesis/pentagram pair accepted")
		}
		reset()
		trigger, other := a, b
		if swap {
			trigger, other = b, a
		}
		trigger.TypeInd = uint16(*o.words["trigger"])
		trigger.ObjFlags = 0x11
		other.ObjFlags = 0x24000
		if call("trigger-exception", b) != 1 {
			t.Fatal("trigger exception rejected")
		}
		reset()
		missile, powder := a, b
		if swap {
			missile, powder = b, a
		}
		missile.ObjClass = 0x2000
		missile.ObjFlags = 8
		powder.ObjFlags = 8
		powder.TypeInd = uint16(*o.words["powder"])
		if call("powder-exception", b) != 1 {
			t.Fatal("powder exception rejected")
		}
	}
	for _, same := range []bool{false, true} {
		reset()
		a.ObjFlags = 0x400
		a.TeamVal.ID = 1
		b.TeamVal.ID = 2
		if same {
			b.TeamVal.ID = 1
		}
		rv := call("team", b)
		if (rv == 0) != same {
			t.Fatal("team suppression", same, rv)
		}
	}
	// Actual monster retained-target storage, including first/last and absent IDs.
	for _, swap := range []bool{false, true} {
		for _, count := range []int{0, 1, 8} {
			for _, slot := range []int{-1, 0, 7} {
				reset()
				monster, target := a, b
				if swap {
					monster, target = b, a
				}
				monster.ObjClass = object.ClassMonster
				monster.ObjFlags = 0x4000
				target.ObjClass = 0x4000
				data := unsafe.Slice((*byte)(monster.UpdateData), 2200)
				data[2172] = byte(count)
				for i := 0; i < 8; i++ {
					binary.LittleEndian.PutUint32(data[2140+4*i:], uint32(2000+i))
				}
				if slot >= 0 {
					binary.LittleEndian.PutUint32(data[2140+4*slot:], target.NetCode)
				}
				want := int32(0)
				if slot >= 0 && slot < count {
					want = 1
				}
				if rv := legacy.PortTestCollisionCore("retained", monster, target, nil, 0); rv != want {
					t.Fatal("retained target", swap, count, slot, rv, want)
				}
				if rv := call("retained", b); rv != want {
					t.Fatal("eligibility retained target", swap, count, slot, rv, want)
				}
			}
		}
	}
	spellbookCapture(t, "collision-core-eligibility", rows, "9aab85749164616e7ab1dec1e1b1c44644687bae8a8722aa9d8165fa48849280")
}
