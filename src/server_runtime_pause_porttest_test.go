//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerRuntimePauseLifecycle(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, active := range []uint32{0, 1, 2} {
		for _, kind := range []int32{-1, 0, 1, 2, 0x7fffffff} {
			for mode := 0; mode < 64; mode++ {
				s := controlsBase(98)
				controlsStats(&s)
				s.Lifecycle.GameFlags |= 2048
				ticks := []uint64{0, 10000, 1<<32 + 7, ^uint64(0)}[mode%4]
				s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.RuntimePause = &legacy.PortTestRuntimePauseSpec{Active: active, Kind: kind, Paused: mode&1 != 0, Provided: mode&2 != 0, Cached: mode&4 != 0, Missing: mode&8 != 0, StateOK: mode&16 != 0, BookBusy: mode&32 != 0, Ticks: ticks}
				cases = append(cases, s)
			}
		}
	}
	interactionCapture(t, "server-runtime-pause", controlsRun(t, cases))
	t.Logf("%d pause lifecycle cases", len(cases))
}
