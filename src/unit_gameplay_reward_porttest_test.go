//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestUnitGameplayRewardOwners(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for mode := 0; mode < 11; mode++ {
		for _, coop := range []bool{false, true} {
			for _, target := range []float32{0, 25.25, math.Nextafter32(25.25, 26), 100, 30000} {
				s := controlsBase(65)
				controlsStats(&s)
				s.Lifecycle.GameFlags &^= 2048
				if coop {
					s.Lifecycle.GameFlags |= 2048
				}
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.ActorWords = map[int]uint32{28: math.Float32bits(25.25)}
				s.Callbacks.Shop.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3684: 10}
				xp := float32(25.25)
				eligible := mode == 0 || mode == 4 || mode == 7 || mode == 8 || mode == 10
				if coop && eligible && target > xp {
					xp = float32(float64(float32(target-xp))*float64(float32(.01)) + 1 + float64(xp))
				}
				a.Controls.UnitReward = &legacy.PortTestUnitRewardSpec{Mode: mode, Target: math.Float32bits(target), WantXP: math.Float32bits(xp)}
				cases = append(cases, s)
			}
		}
	}
	callbackHash(t, "unit-gameplay-reward-owners", controlsRun(t, cases), "7b8d6910a70cfd9473d1f18c84330372264be074d4492de9674175d8b10099c5")
}
