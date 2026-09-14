//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestMapGrowthFrontiers(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	s := growthBase()
	s.Globals["growthInitGrid"] = roomValue(0)
	s.Actions = []legacy.PortTestPaintAction{paintAction(2)}
	cases = append(cases, s)
	for _, kind := range []uint32{0, 1} {
		for _, depth := range []uint32{1, 2, 3} {
			for _, branch := range []uint32{0, 100} {
				for seed := 0; seed < 8; seed++ {
					s := growthBase()
					s.Seed = seed
					s.Records[0].Words[0] = kind
					s.Records[0].Words[64] = math.Float32bits(1000)
					s.Records[0].Words[72] = depth
					s.Records[0].Words[24] = branch
					s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(1)), paintAction(2), paintAction(2)}
					cases = append(cases, s)
				}
			}
		}
	}
	growthCapture(t, "frontiers", cases)
}

func TestMapGrowthObstructions(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, gap := range []int{0, 1, 3, 5} {
			for _, rate := range []uint32{0, 50, 100} {
				for _, flags := range []uint32{0, 2} {
					for seed := 0; seed < 4; seed++ {
						s := growthBase()
						s.Seed = seed
						s.Records[0].Words[64] = math.Float32bits(1000)
						s.Records[0].Words[72] = 3
						s.Records[0].Words[48] = rate
						s.Globals["roomGlobal4"] = roomArg(2)
						x, y := int32(1), int32(-5-gap)
						switch dir {
						case 1:
							y = int32(7 + gap)
						case 2:
							x, y = int32(7+gap), 1
						case 3:
							x, y = int32(-5-gap), 1
						}
						roomSetGeometry(&s.Records[1], 1, 7, 7, 0, 0)
						roomSetGeometry(&s.Records[2], 1, 5, 5, float32(float64(x)*32.526913), float32(float64(y)*32.526913))
						s.Records[2].Words[52] = flags
						s.Records[1].Refs[56] = roomArg(3)
						s.Records[2].Refs[60] = roomArg(2)
						s.Actions = []legacy.PortTestPaintAction{paintAction(3, roomArg(2), roomValue(0), roomValue(0), roomValue(0), roomArg(2))}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	growthCapture(t, "obstructions", cases)
}
