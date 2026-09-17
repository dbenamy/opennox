//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"reflect"
	"testing"
)

func TestCollisionCoreHitQueues(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b, c := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	ids := collisionCoreIDs(a, b, c)
	type row struct {
		Code, Step int
		Hits       [][5]uint32
		Buckets    map[uint32][][2]uint32
	}
	var rows []row
	for _, code := range []uint32{0, 1, 127, 255, 256, 65535} {
		o.resetHits()
		a.NetCode = code
		b.NetCode = 17
		c.NetCode = 273
		type contact struct {
			a, b     *server.Object
			sentinel uint32
		}
		seq := []contact{{a, b, 0}, {b, a, 0}, {a, b, 0}, {a, c, 0}, {a, nil, 0}, {a, nil, 6}, {c, a, 0}, {b, c, 0}}
		var want [][5]uint32
		for i, call := range seq {
			normal := types.Pointf{float32(i) + 0.25, -float32(i) - 0.5}
			aid, bid := ids[call.a.CObj()], call.sentinel
			if call.b != nil {
				bid = ids[call.b.CObj()]
			}
			duplicate := false
			for _, h := range want {
				if h[0] == aid && h[1] == bid || h[0] == bid && h[1] == aid {
					duplicate = true
				}
			}
			if !duplicate {
				hash := call.a.NetCode
				if call.b != nil {
					hash += call.b.NetCode
				}
				want = append([][5]uint32{{aid, bid, math.Float32bits(normal.X), math.Float32bits(normal.Y), hash % 256}}, want...)
			}
			legacy.PortTestCollisionCoreAddHit(call.a, call.b, call.sentinel, &normal)
			got := o.hits(ids)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("queue order/normal/dedup code=%d step=%d: %v want %v", code, i, got, want)
			}
			buckets := legacy.PortTestCollisionCoreBuckets(ids)
			count := 0
			for _, v := range buckets {
				count += len(v)
			}
			if count != len(got) {
				t.Fatal("bucket/list membership count", count, len(got))
			}
			rows = append(rows, row{int(code), i, got, buckets})
		}
		o.resetHits()
		if len(o.hits(ids)) != 0 || len(legacy.PortTestCollisionCoreBuckets(ids)) != 0 {
			t.Fatal("reset retained collision indices")
		}
	}
	spellbookCapture(t, "collision-core-hit-queues", rows, "a11c06406a1cf4b9978e3daa3c13b4a85db129e57659c55887314e4adfd69fd2")
}

func TestCollisionCoreHitCapacity(t *testing.T) {
	o := newCollisionCoreOwner(t)
	units := make([]*server.Object, 46)
	for i := range units {
		units[i] = newObjectXferSimple(t, o.s)
		units[i].NetCode = uint32(i + 1)
	}
	ids := collisionCoreIDs(units...)
	normal := types.Pointf{3, 4}
	for i, a := range units {
		for _, b := range units[i+1:] {
			legacy.PortTestCollisionCoreAddHit(a, b, 0, &normal)
		}
	}
	got := o.hits(ids)
	if len(got) != 1024 {
		t.Fatal("fixed Hit class capacity", len(got))
	}
	type row struct {
		Hits    [][5]uint32
		Buckets map[uint32][][2]uint32
	}
	rows := []row{{got, legacy.PortTestCollisionCoreBuckets(ids)}}
	o.resetHits()
	legacy.PortTestCollisionCoreAddHit(units[0], units[1], 0, &normal)
	if len(o.hits(ids)) != 1 {
		t.Fatal("Hit class not reusable after exhaustion/reset")
	}
	rows = append(rows, row{o.hits(ids), legacy.PortTestCollisionCoreBuckets(ids)})
	spellbookCapture(t, "collision-core-hit-capacity", rows, "294a8d25354209fc8e52523c52e092f6851036dae57f403ebe7cd7399ff39516")
}

func TestCollisionCoreActivationQueue(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b, c := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	units := []*server.Object{a, b, c}
	ids := collisionCoreIDs(units...)
	type row struct {
		Label        string
		Return       uint32
		Queues       [3]uint32
		Links, Flags [3]uint32
	}
	var rows []row
	snapshot := func(label string, rv uint32) {
		r := row{Label: label, Return: rv, Queues: o.queues(ids)}
		for i, u := range units {
			r.Links[i] = collisionCoreID(t, ids, u.Field115)
			r.Flags[i] = u.Field116
		}
		rows = append(rows, r)
	}
	reset := func() {
		o.resetQueues()
		for _, u := range units {
			u.ObjClass = object.ClassSimple
			u.ObjFlags = 4
			u.Collide = o.callback
			u.Update = nil
			u.Field115 = 0
			u.Field116 = 0x40
		}
	}
	reset()
	for _, u := range units {
		if rv := legacy.PortTestCollisionCore("activate", u, nil, nil, 0); rv != 0x41 {
			t.Fatal("activation flags", rv)
		}
		snapshot("append", 0)
	}
	before := len(rows)
	legacy.PortTestCollisionCore("activate", b, nil, nil, 0)
	snapshot("duplicate", 0)
	if rows[before].Queues != rows[before-1].Queues || rows[before].Links != rows[before-1].Links {
		t.Fatal("duplicate activation changed list")
	}
	legacy.PortTestCollisionCore("remove", b, nil, nil, 0)
	snapshot("remove-middle", 0)
	if a.Field115 != uint32(uintptr(c.CObj())) || b.Field115 != 0xffffffff || b.Field116 != 0x40 {
		t.Fatal("middle removal contract")
	}
	legacy.PortTestCollisionCore("remove", c, nil, nil, 0)
	snapshot("remove-tail", 0)
	if o.queues(ids)[2] != 1001 {
		t.Fatal("tail removal contract")
	}
	pop := uint32(legacy.PortTestCollisionCore("pop", nil, nil, nil, 0))
	snapshot("pop-last", collisionCoreID(t, ids, pop))
	if pop != uint32(uintptr(a.CObj())) || o.queues(ids) != ([3]uint32{}) {
		t.Fatal("last pop did not clear queue")
	}
	if legacy.PortTestCollisionCore("next", nil, nil, nil, 0) != 0 {
		t.Fatal("nil traversal")
	}
	reset()
	for _, u := range units {
		legacy.PortTestCollisionCore("activate", u, nil, nil, 0)
	}
	head := uint32(legacy.PortTestCollisionCore("head", nil, nil, nil, 0))
	snapshot("head", collisionCoreID(t, ids, head))
	for _, u := range units {
		next := uint32(legacy.PortTestCollisionCore("next", u, nil, nil, 0))
		snapshot("next", collisionCoreID(t, ids, next))
	}
	legacy.PortTestCollisionCore("remove", a, nil, nil, 0)
	snapshot("remove-head", 0)
	if o.queues(ids)[1] != 1002 {
		t.Fatal("head removal contract")
	}
	// Every remaining element is popped from the production queue, preserving FIFO.
	for _, want := range []uint32{1002, 1003} {
		v := collisionCoreID(t, ids, uint32(legacy.PortTestCollisionCore("pop", nil, nil, nil, 0)))
		if v != want {
			t.Fatal("activation FIFO", v, want)
		}
		snapshot("pop", v)
	}
	for _, class := range []uint32{8, 2, 4, 0x400000, 0x80, 0x2008} {
		for _, flags := range []uint32{0, 4, 8, 12, 64, 68} {
			for _, callback := range []bool{false, true} {
				reset()
				a.ObjClass = object.Class(class)
				a.ObjFlags = object.Flags(flags)
				if !callback {
					a.Collide = nil
				}
				rv := uint32(uint8(legacy.PortTestCollisionCore("activate", a, nil, nil, 0)))
				// Nonqueued returns can be callback low bytes; the observer is aligned to 256.
				snapshot("eligibility", rv)
				if (a.Field116&1 != 0) != (o.queues(ids)[1] != 0) {
					t.Fatal("active bit/queue disagree")
				}
				if legacy.PortTestCollisionCore("active", a, nil, nil, 0) != int32(a.Field116&1) {
					t.Fatal("active membership accessor")
				}
			}
		}
	}
	o.resetQueues()
	for _, u := range units {
		u.Field115 = 0
		u.Field116 = 0
		u.ObjClass = object.ClassSimple
		u.ObjFlags = 0
	}
	spellbookCapture(t, "collision-core-activation-queue", rows, "b32bca07748e41c9bf849644f993665ca91c217de15e602c04363d8c91e03e1b")
}
