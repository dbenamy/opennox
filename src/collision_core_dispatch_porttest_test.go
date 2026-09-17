//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"reflect"
	"testing"
)

func TestCollisionCoreDispatch(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	ids := collisionCoreIDs(a, b)
	damage, words, restore := legacy.PortTestWorldDamageObserver()
	t.Cleanup(restore)
	type row struct {
		Kind   string
		Flags  [2]uint32
		Normal [2]uint32
		Calls  [][4]uint32
		Damage [6]uint32
		Queues [3]uint32
		Hits   [][5]uint32
	}
	var rows []row
	for _, kind := range []string{"pair", "wall", "damage"} {
		for _, flags := range [][2]uint32{{0, 0}, {8, 4}, {4, 8}, {12, 12}} {
			for _, normal := range []types.Pointf{{3, -4}, {math.Float32frombits(0x80000000), 0}, {0.125, -0.25}} {
				o.resetHits()
				o.resetQueues()
				o.resetCalls()
				*words = [6]uint32{}
				worldGeometryResetObject(a, 1001, 100, 100, false)
				worldGeometryResetObject(b, 1002, 105, 100, false)
				a.ObjFlags = object.Flags(flags[0])
				b.ObjFlags = object.Flags(flags[1])
				a.Collide = o.callback
				b.Collide = o.callback
				a.Damage = damage
				b.Damage = damage
				a.Field115 = 0
				b.Field115 = 0
				a.Field116 = 0
				b.Field116 = 0
				target := b
				sentinel := uint32(0)
				if kind != "pair" {
					target = nil
				}
				if kind == "damage" {
					sentinel = 6
				}
				legacy.PortTestCollisionCoreAddHit(a, target, sentinel, &normal)
				before := o.hits(ids)
				legacy.PortTestCollisionCore("dispatch", nil, nil, nil, 0)
				calls := o.calls(ids)
				var want [][4]uint32
				switch kind {
				case "pair":
					want = [][4]uint32{{1001, 1002, math.Float32bits(normal.X), math.Float32bits(normal.Y)}, {1002, 1001, math.Float32bits(-normal.X), math.Float32bits(-normal.Y)}}
				case "wall":
					want = [][4]uint32{{1001, 0, math.Float32bits(normal.X), math.Float32bits(normal.Y)}}
				case "damage":
					want = make([][4]uint32, 0)
				}
				if !reflect.DeepEqual(calls, want) {
					t.Fatalf("dispatch order/normals %s: %v want %v", kind, calls, want)
				}
				d := *words
				if d[1] != 0 {
					d[1] = collisionCoreID(t, ids, d[1])
				}
				if kind == "damage" {
					if d != ([6]uint32{1, 1001, 0, 0, 2, 12}) {
						t.Fatal("damage sentinel", d)
					}
				} else if d != ([6]uint32{}) {
					t.Fatal("unexpected damage callback", d)
				}
				if !reflect.DeepEqual(before, o.hits(ids)) {
					t.Fatal("dispatch consumed or modified Hit records")
				}
				rows = append(rows, row{kind, flags, [2]uint32{math.Float32bits(normal.X), math.Float32bits(normal.Y)}, calls, d, o.queues(ids), o.hits(ids)})
			}
		}
	}
	o.resetQueues()
	a.Field115 = 0
	b.Field115 = 0
	a.Field116 = 0
	b.Field116 = 0
	a.Damage = nil
	b.Damage = nil
	spellbookCapture(t, "collision-core-dispatch", rows, "8ffdb699a1072cb056baccb6253f6c8666391883ebca9a81ee8784e06fa27d8c")
}
