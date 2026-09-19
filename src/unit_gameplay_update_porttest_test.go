//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestUnitGameplayUndeadUpdate(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, frame := range []uint32{0, 1, 70, 71, 1000, 0x80000000, 0xffffffff} {
		for _, age := range []uint32{0, 1, 69, 70, 71, 72, 0x7fffffff, 0xffffffff} {
			for _, target := range []bool{false, true} {
				for _, state := range []byte{0, 1, 2, 3, 128, 255} {
					s := controlsBase(64)
					s.Owner.Frame = frame
					s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.UnitUpdate = &legacy.PortTestUnitUpdateSpec{Frame: frame, Spawn: frame - age, Target: target, TargetState: state}
					cases = append(cases, s)
				}
			}
		}
	}
	callbackHash(t, "unit-gameplay-undead-update", controlsRun(t, cases), "f4eb94f69053cec87515303f4af2dced00900641bddd619788698a44f6323ef3")
}
