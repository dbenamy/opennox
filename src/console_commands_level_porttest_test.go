//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestConsoleCommandsLevel(t *testing.T) {
	t.Cleanup(legacy.PortTestConsoleContext())
	var cases []legacy.PortTestRoamSpec
	for class := byte(0); class < 3; class++ {
		for _, mode := range []uint32{0, 8192} {
			for _, tc := range []struct {
				text  string
				level int32
			}{{"0", 0}, {"1", 1}, {"9", 9}, {"10", 10}, {"11", 10}, {"127", 10}, {"256", 0}, {"257", 1}, {"invalid", 0}} {
				s := controlsBase(57)
				controlsStats(&s)
				s.Lifecycle.GameFlags |= mode
				p := s.Callbacks.Shop
				p.Resources.PlayerClass = class
				o := p.TemporaryUpdates.World.Objectives
				o.PlayerDataWords[0] = map[int]uint32{3684: 3}
				o.Attack.Controls.Name = &tc.text
				o.Attack.Controls.X = tc.level
				if mode == 8192 {
					o.Attack.Controls.X = 3
				}
				cases = append(cases, s)
			}
		}
	}
	spellbookCapture(t, "console-commands-level", controlsRun(t, cases), "ab4e4e4722faff15825a1cb506dd7978870f18d8870154732a476a01bd5dc064")
}
