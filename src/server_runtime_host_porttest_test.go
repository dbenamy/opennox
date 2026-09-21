//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerRuntimeHostOwnership(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for gate := 0; gate < 32; gate++ {
		for count := 0; count <= 2; count++ {
			for pattern := 0; pattern < 8; pattern++ {
				s := controlsBase(99)
				controlsStats(&s)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls
				a.MonsterRefs = []int{4, 5}
				a.RuntimeHost = &legacy.PortTestRuntimeHostSpec{Host: gate&1 != 0, Unit: gate&2 != 0, Marker: gate&4 != 0, Missing: gate&8 != 0, Cached: gate&16 != 0, Count: count, Reverse: pattern&1 != 0, MonsterMask: pattern >> 1, MonitorMask: (pattern >> 1) ^ 2}
				cases = append(cases, s)
			}
		}
	}
	interactionCapture(t, "server-runtime-host", controlsRun(t, cases))
	t.Logf("%d host ownership cases", len(cases))
}
