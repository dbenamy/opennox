//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestCollisionRegistryWorldUndeadDamage(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := newCreatureXferObject(t, o.s, "Monster")
	data := o.record(t, 4)
	budget := o.record(t, 80)
	objectXferSetWord(data, 0, uint32(uintptr(budget)))
	a.CollideData = data
	t.Cleanup(func() { a.CollideData = nil; o.s.Objs.DeletedList = nil; o.s.ObjSetOwner(nil, a) })
	callback, observed, restore := legacy.PortTestWorldDamageObserver()
	t.Cleanup(restore)
	b.Damage = callback
	var rows []struct {
		Name                           string
		Remaining, Calls, Amount, Kind uint32
		Deleted                        bool
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-undead-damage", rows)
	}()
	for _, hp := range []uint16{0, 1, 100, 65535} {
		for _, available := range []uint32{0, 1, 99, 100, 101, 65535, 0xffffffff, 0x80000000} {
			for _, undead := range []bool{false, true} {
				for _, owned := range []bool{false, true} {
					name := fmt.Sprintf("hp%d/budget%d/undead%t/owned%t", hp, available, undead, owned)
					t.Run(name, func(t *testing.T) {
						o.reset()
						*observed = [6]uint32{}
						a.ObjFlags = 0
						a.DeletedNext = nil
						o.s.Objs.DeletedList = nil
						o.s.ObjSetOwner(nil, a)
						if owned {
							o.s.ObjSetOwner(&o.units[0], a)
						}
						b.ObjSubClass = 0
						if undead {
							b.ObjSubClass = 0x40
						}
						b.HealthData.Cur = hp
						objectXferSetWord(budget, 72, available)
						collisionRegistryWorld(17, a, b, nil)
						amount := uint32(hp)
						if int32(available) <= int32(hp) {
							amount = available
						}
						calls := uint32(0)
						remaining := available
						deleted := false
						if undead {
							if available != 0 {
								calls = 1
							}
							remaining -= amount
							deleted = int32(available) <= int32(hp)
						}
						if observed[0] != calls || objectXferGetWord(budget, 72) != remaining || a.ObjFlags.Has(object.FlagDestroyed) != deleted {
							t.Fatal("undead budget/deletion")
						}
						if calls != 0 {
							owner := a
							if owned {
								owner = &o.units[0]
							}
							if observed[1] != uint32(uintptr(b.CObj())) || observed[2] != uint32(uintptr(owner.CObj())) || observed[3] != uint32(uintptr(a.CObj())) || observed[4] != amount || observed[5] != 6 {
								t.Fatal("damage callback arguments")
							}
						}
						rows = append(rows, struct {
							Name                           string
							Remaining, Calls, Amount, Kind uint32
							Deleted                        bool
						}{name, objectXferGetWord(budget, 72), observed[0], observed[4], observed[5], deleted})
					})
				}
			}
		}
	}
	for _, normal := range []*types.Pointf{nil, new(types.Pointf)} {
		a.ObjFlags = 0
		a.DeletedNext = nil
		o.s.Objs.DeletedList = nil
		collisionRegistryWorld(17, a, nil, normal)
		if a.ObjFlags.Has(object.FlagDestroyed) != (normal == nil) {
			t.Fatal("empty collision deletion")
		}
	}
}
