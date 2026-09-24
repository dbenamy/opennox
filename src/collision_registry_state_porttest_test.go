//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestCollisionRegistryObjectState(t *testing.T) {
	for _, op := range []int{40, 41, 42} {
		var cases []legacy.PortTestRoamSpec
		for _, target := range []int{0, 4, 101} {
			for _, enabled := range []bool{false, true} {
				for _, frame := range []uint32{0, 100, 0xffffffff} {
					s := objectStateBase(op)
					s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.State.RegisteredCollision = true
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.State.Target = target
					s.Owner.Frame = frame
					if op != 42 {
						p.Resources.Subject = 3
						a.UpdateWords = map[int]uint32{1272: 0xffffffff}
						if enabled {
							a.UpdateWords[1272] = 17
						} else if op == 41 {
							a.ActorRefs = map[int]int{508: 101}
						}
					} else if enabled {
						p.Equipment.ActiveAbilities = 2
					}
					p.Items[1].Health = true
					p.Items[1].HP = 100
					p.Items[1].MaxHP = 100
					p.Items[1].Flags = 0
					p.Items[1].Class = 0x400000
					a.State.HealthRefs = []int{101}
					p.EffectsUse.Balance["BerserkerDamage"] = []float64{12}
					p.EffectsUse.Balance["BerserkerStunDuration"] = []float64{30}
					p.EffectsUse.Balance["BerserkerPainRatio"] = []float64{0.25}
					cases = append(cases, s)
				}
			}
		}
		name := fmt.Sprintf("collision-registry-state-%02d", op)
		collisionRegistryHash(t, name, effectsTimedRun(t, cases), collisionRegistryOwnerHashes[name])
	}
}
