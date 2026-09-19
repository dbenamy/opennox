//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

func TestUnitGameplayReadUse(t *testing.T) {
	oldInc, oldInfinite := questLevelWarpInc, questLevelWarpInfinite
	questLevelWarpInc, questLevelWarpInfinite = 5, false
	defer func() { questLevelWarpInc, questLevelWarpInfinite = oldInc, oldInfinite }()
	var cases []legacy.PortTestRoamSpec
	for _, warp := range []bool{false, true} {
		for _, frame := range []uint32{0, 1, 1000, 0xffffffff} {
			for _, fps := range []uint32{1, 30, 60} {
				for _, age := range []uint32{0, 3*fps - 1, 3 * fps, 3*fps + 1, 0xffffffff} {
					for gate := 0; gate < 8; gate++ {
						s := controlsBase(66)
						s.Owner.Frame, s.Owner.FPS = frame, fps
						s.Combat.Wall = gate & 1
						s.Callbacks.Shop.Inventory.WallMode = gate & 1
						stage := uint32((gate%3)*10 + 4)
						threshold := (stage/5 + 1) * 5
						if stage >= 20 {
							threshold = stage
						}
						c := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls
						c.UnitRead = &legacy.PortTestUnitReadSpec{Warp: warp, NonPlayer: gate&2 != 0, Blocked: gate&1 != 0, Frame: frame, FPS: fps, Timestamp: frame - age, Stage: stage, Allowed: uint32(gate & 4), Threshold: threshold, Text: "PortRead:Message"}
						c.Reports = &legacy.PortTestGameplayReportsSpec{}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	callbackHash(t, "unit-gameplay-read-use", controlsRun(t, cases), "b388b9fa8ac7af7eaeb4c88451b891cdce3e181bc85573963d3cf8144f57c795")
}

func TestUnitGameplayReadMessageLimits(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, text := range []string{"", "X", strings.Repeat("X", 48), strings.Repeat("X", 49), strings.Repeat("X", 255)} {
		for _, warp := range []bool{false, true} {
			for _, allowed := range []uint32{0, 1} {
				s := controlsBase(66)
				s.Owner.Frame, s.Owner.FPS = 1000, 30
				c := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls
				c.UnitRead = &legacy.PortTestUnitReadSpec{Warp: warp, Frame: 1000, FPS: 30, Text: text, Stage: 4, Allowed: allowed, Threshold: 5}
				c.Reports = &legacy.PortTestGameplayReportsSpec{}
				cases = append(cases, s)
			}
		}
	}
	oldInc, oldInfinite := questLevelWarpInc, questLevelWarpInfinite
	questLevelWarpInc, questLevelWarpInfinite = 5, false
	defer func() { questLevelWarpInc, questLevelWarpInfinite = oldInc, oldInfinite }()
	callbackHash(t, "unit-gameplay-read-message-limits", controlsRun(t, cases), "d0abb11062f5e23a88aaa7469265d10f2c536f338afc3debbe4820ad27aba3df")
}
