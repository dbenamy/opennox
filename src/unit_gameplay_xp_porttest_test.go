//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestUnitGameplayExperience(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	thresholds := []float32{0, 10, 30, 60, 100, 150, 210, 280, 360, 450, 550, 660}
	for op := 0; op < 3; op++ {
		for _, level := range []byte{0, 1, 5, 9, 10} {
			threshold := thresholds[int(level)+1]
			for _, xp := range []float32{0, math.Nextafter32(threshold, 0), threshold, math.Nextafter32(threshold, 1000), 10000.125} {
				for _, amount := range []float32{0, 0.125, 10.5, 999.75} {
					for gate := 0; gate < 4; gate++ {
						s := controlsBase(61)
						controlsStats(&s)
						p := s.Callbacks.Shop
						a := p.TemporaryUpdates.World.Objectives.Attack
						p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3684: uint32(level)}
						a.ActorWords = map[int]uint32{28: math.Float32bits(xp)}
						next, wantLevel := xp, level
						if op == 1 {
							next = float32(float64(xp) + float64(amount))
						}
						if op == 2 && xp < amount {
							next = float32(float64(float32(amount-xp))*float64(float32(0.01)) + 1 + float64(xp))
						}
						if (op != 2 || xp < amount) && gate != 3 && next >= threshold {
							wantLevel++
						}
						a.Controls.UnitExperience = &legacy.PortTestUnitExperienceSpec{Op: op, Amount: math.Float32bits(amount), WantXP: math.Float32bits(next), WantLevel: wantLevel, Saving: gate&1 != 0, HasSave: gate&2 != 0}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	callbackHash(t, "unit-gameplay-experience", controlsRun(t, cases), "8252317eec3326e2451dbb02d0209ebca4d196a705b088053b79881496181a1a")
	t.Logf("%d experience cases", len(cases))
}

// Literal results distinguish the subtraction spill and the wide addition from
// both an entirely wide formula and prematurely rounded award addition.
func TestUnitGameplayExperienceRounding(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, row := range [][4]float32{
		{74820.25, 9668373, 95936.53125, 170756.78125},
		{88648.09375, 3324198.25, 32356.501953125, 121004.59375},
		{47860.36328125, 4699923, 46521.625, 94381.984375},
		{26312.0390625, 1288875.125, 12626.630859375, 38938.671875},
		{52092.32421875, 7867120, 78151.2734375, 130243.59375},
		{36293.68359375, 9860146, 98239.515625, 134533.203125},
		{79788.1953125, 1360007.25, 12803.189453125, 92591.3828125},
		{37064.390625, 8656686, 86197.21875, 123261.609375},
	} {
		s := controlsBase(61)
		controlsStats(&s)
		a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
		a.ActorWords = map[int]uint32{28: math.Float32bits(row[0])}
		s.Callbacks.Shop.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3684: 1}
		award := math.Float32bits(row[2])
		a.Controls.UnitExperience = &legacy.PortTestUnitExperienceSpec{Op: 2, Amount: math.Float32bits(row[1]), WantAward: &award, WantXP: math.Float32bits(row[3]), WantLevel: 1, Saving: true, HasSave: true}
		cases = append(cases, s)
	}
	callbackHash(t, "unit-gameplay-experience-rounding", controlsRun(t, cases), "2b3dd1835447472c6cb976ca84051e5dc2690e42547f140e86b62d04fa4f87bf")
}

func TestUnitGameplayExperiencePresentation(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for op := 0; op < 3; op++ {
		for _, missing := range []bool{false, true} {
			s := controlsBase(61)
			controlsStats(&s)
			s.Lifecycle.GameFlags |= 2048
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3684: 1}
			current, amount, want := float32(30), float32(0), float32(30)
			if op == 1 {
				current, amount, want = 25, 5, 30
			}
			if op == 2 {
				current, amount = 29, 30
				want = float32(float64(float32(amount-current))*float64(float32(.01)) + 1 + float64(current))
			}
			a.ActorWords = map[int]uint32{28: math.Float32bits(current)}
			a.Controls.UnitExperience = &legacy.PortTestUnitExperienceSpec{Op: op, Amount: math.Float32bits(amount), WantXP: math.Float32bits(want), WantLevel: 2, Presentation: true, MissingPresentation: missing}
			cases = append(cases, s)
		}
	}
	callbackHash(t, "unit-gameplay-xp-presentation", controlsRun(t, cases), "76fa81b0a0518d313092c614f50d6056ca7ad9ad73c72fa710671c6864a42ee6")
}
