//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestUpdateRegistryObjectivesFlagDeadline(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, fps := range []uint32{2, 30, 60} {
		for _, stamp := range []uint32{0, 1, 100, 0xfffffffa} {
			for _, delta := range []uint32{0, 1, 30*fps - 1, 30 * fps, 30*fps + 1, 0xffffffff} {
				s := updateRegistryObjectiveBase()
				p := s.Callbacks.Shop
				s.Owner.FPS = fps
				s.Owner.Frame = stamp + delta
				p.Items[0].Class = 0x10000000
				p.TemporaryUpdates.World.Objectives.ColorNames = []string{"RedFlag"}
				p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: math.Float32bits(600), 4: math.Float32bits(650), 8: stamp}
				p.Sequence = []legacy.PortTestShopAction{{Op: 809}, {Op: 809}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-flag-deadline", cases)
}

func TestUpdateRegistryObjectivesBallClock(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, tick := range []uint64{0, 1, 19999, 20000, 20001, 0x7fffffff, 0x80000000, 0xffffffff, 0x100000000, 0x100004e20} {
		for _, stamp := range []uint64{0, 1, 20000, 0xffffffff, 0x100000000} {
			for _, speed := range []float32{0, 1, 5} {
				s := updateRegistryObjectiveBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.World.Objectives.Ticks = []uint64{tick}
				w.UpdateWords[0] = map[int]uint32{8: uint32(stamp), 12: uint32(stamp >> 32), 24: math.Float32bits(speed)}
				w.ItemWords[0][80] = math.Float32bits(3)
				w.ItemWords[0][84] = math.Float32bits(4)
				p.Sequence = []legacy.PortTestShopAction{{Op: 810}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-ball-clock", cases)
}

func TestUpdateRegistryObjectivesObeliskIdle(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, fps := range []uint32{2, 3, 30, 60} {
		for _, frame := range []uint32{0, 1, 14, 15, 16, 30, 0xffffffff} {
			for _, energy := range []uint32{0, 1, 7, 8, 49, 50, 51, 0xffffffff} {
				s := updateRegistryObjectiveBase()
				p := s.Callbacks.Shop
				s.Owner.FPS = fps
				s.Owner.Frame = frame
				p.TemporaryUpdates.World.Objectives.Players = 0
				p.TemporaryUpdates.UpdateWords[0][0] = energy
				p.Sequence = []legacy.PortTestShopAction{{Op: 808}, {Op: 808}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-obelisk-idle", cases)
}

func TestUpdateRegistryObjectivesObeliskTransfer(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mode := range []uint32{0, 4096, 8192, 12288} {
		for _, energy := range []uint32{0, 1, 2, 7, 8, 49, 50} {
			for _, mana := range []uint16{0, 49, 50} {
				for _, distance := range []float32{0, 49.999, 50, 50.001} {
					s := updateRegistryObjectiveBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					o := w.World.Objectives
					s.Lifecycle.GameFlags = mode | 1
					p.Resources.Mana = mana
					p.Resources.MaxMana = 50
					o.PlayerWords[0][56] = math.Float32bits(512 + distance)
					w.UpdateWords[0][0] = energy
					p.Sequence = []legacy.PortTestShopAction{{Op: 808}, {Op: 808}}
					cases = append(cases, s)
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-obelisk-transfer", cases)
}

func TestUpdateRegistryObjectivesCrownDispatch(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{804, 811} {
		for _, target := range []int{0, 100, 4} {
			for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
				for _, accept := range []bool{false, true} {
					s := updateRegistryObjectiveBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					w.Target = target
					w.UpdateRefs[0] = map[int]int{4: target}
					w.World.Objectives.PlayerWords[0][16] = flags
					p.Items[1].Flags = flags
					p.Resources.PickupResult = accept
					p.Sequence = []legacy.PortTestShopAction{{Op: op}}
					cases = append(cases, s)
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-crown-dispatch", cases)
}

func TestUpdateRegistryObjectivesObeliskWands(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mode := range []uint32{0, 4096, 8192, 12288} {
		for _, energy := range []uint32{0, 1, 2, 8} {
			for _, charges := range []uint32{0, 19, 20} {
				for _, rate := range []float64{0, 1, 2, 5} {
					s := updateRegistryObjectiveBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					o := w.World.Objectives
					s.Lifecycle.GameFlags = mode | 1
					p.Resources.Mana = 49
					p.Resources.MaxMana = 50
					w.UpdateWords[0][0] = energy
					o.PlayerUpdateRefs[0] = map[int]int{104: 4}
					p.Items[1].Class = 0x1000
					p.Items[1].Subclass = 0x4000000
					o.UseWords = []map[int]uint32{nil, {108: charges | 20<<8, 112: charges * 5}}
					p.EffectsUse.Balance["OblivionStaffRechargeRate"] = []float64{rate}
					p.Sequence = []legacy.PortTestShopAction{{Op: 808}, {Op: 808}}
					cases = append(cases, s)
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-obelisk-wands", cases)
}

func TestUpdateRegistryObjectivesObeliskEligibility(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, team := range []byte{0, 1, 2} {
		for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
			for _, class := range []byte{0, 1, 2, 3} {
				for _, frame := range []uint32{0, 14, 15, 16, 0xffffffff} {
					s := updateRegistryObjectiveBase()
					p := s.Callbacks.Shop
					s.Lifecycle.GameFlags = 4097
					s.Owner.Frame = frame
					p.Resources.PlayerClass = class
					p.Resources.Mana = 0
					p.Resources.MaxMana = 50
					p.Inventory.Teams = [4]byte{1, 2, 2, team}
					p.TemporaryUpdates.World.Objectives.PlayerWords[0][16] = flags
					p.TemporaryUpdates.UpdateWords[0][0] = 8
					p.Sequence = []legacy.PortTestShopAction{{Op: 808}, {Op: 808}}
					cases = append(cases, s)
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-objectives-obelisk-eligibility", cases)
}
