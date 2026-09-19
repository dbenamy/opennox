//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestUnitGameplayOrders(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, player := range []int32{1, 7, 31} {
		for _, order := range []int32{-2147483648, -257, -1, 0, 1, 2, 3, 4, 255, 256, 257, 2147483647} {
			s := controlsBase(62)
			o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
			o.Players = 3
			o.Attack.Controls.X, o.Attack.Controls.Y = player, order
			o.Attack.Controls.Reports = &legacy.PortTestGameplayReportsSpec{}
			cases = append(cases, s)
		}
	}
	callbackHash(t, "unit-gameplay-orders", controlsRun(t, cases), "3ee874f06881dc5f757e73dbe215dbe86ae619b899ffa3767d625df4d138429a")
}
func TestUnitGameplayDialogue(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, finish := range []bool{false, true} {
		for gate := 0; gate < 10; gate++ {
			if finish && gate != 0 && gate != 1 && gate != 2 && gate != 9 {
				continue
			}
			for _, response := range []byte{0, 1, 127, 128, 255} {
				for _, kind := range []byte{0, 1, 2, 255} {
					for freeze := 0; freeze < 4; freeze++ {
						s := controlsBase(63)
						c := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls
						c.UnitDialogue = &legacy.PortTestUnitDialogueSpec{Finish: finish, Gate: gate, Response: response, Kind: kind, Frozen: freeze&1 != 0, FreezeLock: freeze&2 != 0}
						c.Reports = &legacy.PortTestGameplayReportsSpec{}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	callbackHash(t, "unit-gameplay-dialogue", controlsRun(t, cases), "153f50d5cb39844f175edff598bf4cd413b77876e80e44e6b5590174cf6dcc51")
}
