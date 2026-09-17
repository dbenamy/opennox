//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestSpatialTargetingProjectileOwner(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := newCreatureXferObject(t, o.s, "Monster")
	owner := newCreatureXferObject(t, o.s, "Monster")
	sa, sb, so := *a, *b, *owner
	t.Cleanup(func() { *a = sa; *b = sb; *owner = so })
	type row struct {
		Owner                                                bool
		SourceClass, SourceSubclass, TargetClass, OwnerClass uint32
		Teams                                                int
		Enemy                                                bool
		Return                                               int32
	}
	var rows []row
	admitted, rejected := 0, 0
	for _, owned := range []bool{false, true} {
		for _, ac := range []object.Class{1, 8} {
			for _, sub := range []object.SubClass{0, 2} {
				for _, bc := range []object.Class{2, 8} {
					for _, oc := range []object.Class{2, 8} {
						for teams := 0; teams < 3; teams++ {
							*a = sa
							*b = sb
							*owner = so
							a.ObjClass = ac
							a.ObjSubClass = sub
							b.ObjClass = bc
							owner.ObjClass = oc
							a.ObjFlags = 0
							b.ObjFlags = 0
							owner.ObjFlags = 0
							a.Collide = o.callback
							b.Collide = o.callback
							a.ObjOwner = nil
							if owned {
								a.ObjOwner = owner
							}
							a.TeamVal.ID = 0
							b.TeamVal.ID = 0
							owner.TeamVal.ID = 0
							if teams != 0 {
								b.TeamVal.ID = 1
								owner.TeamVal.ID = server.TeamID(teams)
							}
							enemy := o.s.IsEnemyTo(owner, b)
							rv := legacy.PortTestSpatialEligible(a, b, false)
							if rv == 0 {
								rejected++
							} else {
								admitted++
							}
							if !owned && rv != 1 {
								t.Fatal("unowned ordinary projectile rejected")
							}
							if owned && ac == 1 && sub == 0 && bc == 2 && oc == 2 && teams == 1 && rv != 0 {
								t.Fatal("projectile hit owner's team")
							}
							if owned && ac == 1 && sub == 2 && rv != 1 {
								t.Fatal("owner exemption subclass not honored")
							}
							rows = append(rows, row{owned, uint32(ac), uint32(sub), uint32(bc), uint32(oc), teams, enemy, rv})
						}
					}
				}
			}
		}
	}
	if admitted == 0 || rejected == 0 {
		t.Fatal("owner decisions not exercised")
	}
	spellbookCapture(t, "spatial-targeting-projectile-owner", rows, "ebb54f2384ded8e7f3d00a3ab5e981565afb31c4ef0ea3fee0b096444ea08c9d")
}
