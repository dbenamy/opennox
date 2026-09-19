//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestUnitGameplayMonsterTimers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{28, 29} {
		for _, frame := range []uint32{0, 1, 100, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, deadline := range []uint32{0, 1, frame - 1, frame, frame + 1, 0xffffffff} {
				for _, fps := range []uint32{1, 30, 60} {
					for gate := 0; gate < 4; gate++ {
						s := stateBase(op)
						s.Seed = len(cases) + 1
						sp := s.MonsterState
						sp.Frame, sp.Deadline, sp.FPS = frame, deadline, fps
						sp.NilUnit, sp.AnimData = gate&1 != 0, gate&2 != 0
						cases = append(cases, s)
					}
				}
			}
		}
	}
	stateHash(t, "unit-gameplay-monster-timers", legacy.PortTestRoam(cases), "2ef4ef3fe482bb27703e319d41f6cd41aaeebb1abbbe1c7f028b575d9ab581f4")
}
