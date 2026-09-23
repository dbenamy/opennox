//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestUpdateRegistryTemporaryLifetimeStates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53ADC0, legacy.PortTestTemporary53B8F0, legacy.PortTestTemporary53CB60, legacy.PortTestTemporary53CC90, legacy.PortTestTemporary53DB00} {
		for _, frame := range []uint32{0, 1, 2, 3, 9, 10, 11, 29, 30, 31, 100, 0xffffffff} {
			for _, life := range []uint32{0, 1, 10, 0xffffffff} {
				for _, cb := range []bool{false, true} {
					s := updateRegistryTemporaryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					p.TemporaryUpdates.ItemWords[0][128] = 0
					p.TemporaryUpdates.ItemWords[0][136] = 10
					p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: life, 4: life, 12: 4}
					p.TemporaryUpdates.DieCallback = cb
					p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-lifetime", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-lifetime"])
}

func TestUpdateRegistryTemporaryBreakStates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53DB30, legacy.PortTestTemporary53DBB0, legacy.PortTestTemporary53DC30} {
		for _, state := range []uint32{0, 1, 2, 3, 4, 6, 8, 10, 12, 0xffffffff} {
			for _, flags := range []uint32{4, 0x8004, 0x8044} {
				for _, frame := range []uint32{0, 99, 100, 101, 0xffffffff} {
					s := updateRegistryTemporaryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					p.TemporaryUpdates.Updatable = 3
					p.TemporaryUpdates.ItemWords[0][20] = state
					p.TemporaryUpdates.ItemWords[0][16] = flags
					p.TemporaryUpdates.ItemWords[0][136] = 100
					p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-break", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-break"])
}

func TestUpdateRegistryTemporaryTrails(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, dir := range []uint32{0, 1, 63, 64, 127, 128, 255} {
		for _, speed := range []float32{0, .125, 2.5, 100} {
			s := updateRegistryTemporaryBase()
			p := s.Callbacks.Shop
			w := p.TemporaryUpdates.ItemWords[0]
			w[124] = dir
			w[544] = math.Float32bits(speed)
			w[72] = math.Float32bits(508.125)
			w[76] = math.Float32bits(509.875)
			w[88] = math.Float32bits(-.25)
			w[92] = math.Float32bits(.5)
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTemporary53AEC0}}
			specs = append(specs, s)
		}
	}
	results := effectsTimedRun(t, specs)
	for i, r := range results {
		if len(r.Lifecycle.Created) != 8 {
			t.Fatalf("case %d created %d want 8", i, len(r.Lifecycle.Created))
		}
	}
	updateRegistryHash(t, "update-registry-temporary-trails", results, updateRegistryHashes["update-registry-temporary-trails"])
}

func TestUpdateRegistryTemporaryHoming(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53B940, legacy.PortTestTemporary53BB00, legacy.PortTestTemporary53BDA0, legacy.PortTestTemporary53DCC0} {
		for _, mode := range []int{0, 1, 2, 3} {
			for _, dir := range []uint32{0, 1, 63, 127, 248, 255} {
				for _, frame := range []uint32{96, 100, 200} {
					s := updateRegistryTemporaryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					tmp := p.TemporaryUpdates
					tmp.ItemWords[0][124] = dir | (dir << 16)
					tmp.ItemWords[0][80] = math.Float32bits(.25)
					tmp.ItemWords[0][84] = math.Float32bits(-.5)
					tmp.ItemRefs[0] = map[int]int{508: 1, 504: 4}
					tmp.UpdateRefs[0] = map[int]int{0: 1, 4: 4, 8: 4, 12: 4}
					tmp.ItemWords[1][56] = math.Float32bits(520)
					tmp.ItemWords[1][60] = math.Float32bits(515)
					if op == legacy.PortTestTemporary53B940 {
						tmp.UpdateRefs[0] = map[int]int{0: 1, 4: 4, 8: 1}
						tmp.UpdateWords[0] = map[int]uint32{12: 41}
					}
					if op == legacy.PortTestTemporary53DCC0 {
						tmp.UpdateRefs[0] = map[int]int{8: 4, 12: 4}
						tmp.UpdateWords[0] = map[int]uint32{16: math.Float32bits(520), 20: math.Float32bits(515)}
					}
					if mode == 1 {
						tmp.ItemWords[1][16] = 0x24
					}
					if mode == 2 {
						tmp.ItemWords[1][16] = 0x8004
					}
					if mode == 3 {
						tmp.UpdateRefs[0][4] = 0
						tmp.UpdateRefs[0][8] = 0
					}
					tmp.CollideReturn = 0x81234567
					p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-homing", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-homing"])
}

func TestUpdateRegistryTemporaryOwnerEffects(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53D270, legacy.PortTestTemporary53D330, legacy.PortTestTemporary53D400, legacy.PortTestTemporary53D510} {
		for _, frame := range []uint32{90, 91, 94, 100, 700, 10000} {
			for _, state := range []uint32{0, 0x20, 0x8000} {
				for _, wall := range []bool{false, true} {
					s := updateRegistryTemporaryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					tmp := p.TemporaryUpdates
					tmp.ItemRefs[0] = map[int]int{508: 1}
					tmp.ItemWords[0][16] = 4 | state
					tmp.ItemWords[0][104] = math.Float32bits(0)
					tmp.ItemWords[0][156] = math.Float32bits(510)
					tmp.ItemWords[0][160] = math.Float32bits(512)
					tmp.ItemWords[0][72] = math.Float32bits(511)
					tmp.ItemWords[0][76] = math.Float32bits(512)
					tmp.ItemWords[0][136] = 110
					p.EffectsUse.UnitWords = map[int]uint32{16: 4 | state}
					p.EffectsUse.PlayerWords = map[int]uint32{2284: 520, 2288: 512}
					if wall {
						p.Inventory.WallMode = 1
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-owner-effects", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-owner-effects"])
}

func TestUpdateRegistryTemporarySpawners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53C9A0, legacy.PortTestTemporary53D5A0, legacy.PortTestTemporary53D6E0, legacy.PortTestTemporary53DA60} {
		for _, frame := range []uint32{89, 90, 91, 100, 150, 241} {
			for _, stamp := range []uint32{90, 100, 101} {
				for _, wall := range []bool{false, true} {
					s := updateRegistryTemporaryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					tmp := p.TemporaryUpdates
					tmp.ItemWords[0][136] = stamp
					tmp.ItemWords[0][104] = 0
					tmp.ItemRefs[0] = map[int]int{508: 1}
					tmp.UpdateWords[0] = map[int]uint32{0: 17}
					if wall {
						p.Inventory.WallMode = 1
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-spawners", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-spawners"])
}

func TestUpdateRegistryTemporaryAreaOwners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53CB90, legacy.PortTestTemporary53CCB0, legacy.PortTestTemporary53D220, legacy.PortTestTemporary53D850, legacy.PortTestTemporary53D960} {
		for _, frame := range []uint32{90, 93, 98, 100, 101, 150, 151} {
			for _, life := range []uint32{0, 1, 2} {
				for _, indexed := range []bool{false, true} {
					s := updateRegistryTemporaryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					tmp := p.TemporaryUpdates
					tmp.UpdateWords[0] = map[int]uint32{0: life}
					tmp.ItemWords[0][136] = 100
					tmp.ItemRefs[0] = map[int]int{508: 1}
					p.Items[1].Class = 0x2000
					tmp.ItemWords[1][56] = math.Float32bits(525)
					tmp.ItemWords[1][60] = math.Float32bits(512)
					if indexed {
						tmp.Indexed = []int{1, 4}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-area-owners", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-area-owners"])
}

func TestUpdateRegistryTemporaryFrameWrap(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53B8F0, legacy.PortTestTemporary53CB60, legacy.PortTestTemporary53CC90, legacy.PortTestTemporary53CCB0, legacy.PortTestTemporary53D220, legacy.PortTestTemporary53DB00, legacy.PortTestTemporary53D850, legacy.PortTestTemporary53DA60} {
		for _, start := range []uint32{0, 0xfffffff0, 0x80000000} {
			for _, elapsed := range []uint32{0, 2, 3, 10, 30, 31, 61, 0x80000000, 0xffffffff} {
				s := updateRegistryTemporaryBase()
				p := s.Callbacks.Shop
				s.Owner.Frame = start + elapsed
				p.TemporaryUpdates.ItemWords[0][128] = start
				p.TemporaryUpdates.ItemWords[0][136] = start + 10
				p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: 10, 4: 10}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				specs = append(specs, s)
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-frame-wrap", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-frame-wrap"])
}

func TestUpdateRegistryTemporaryMissingAndOwners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary54FD80, legacy.PortTestTemporary53AEC0, legacy.PortTestTemporary53BB00} {
		s := updateRegistryTemporaryBase()
		p := s.Callbacks.Shop
		p.TemporaryUpdates.MissingTypes = []string{"Spark"}
		p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1, Side: 20}}
		specs = append(specs, s)
	}
	for _, op := range []int{legacy.PortTestTemporary53D270, legacy.PortTestTemporary53DCC0} {
		for _, inventory := range []int{0, 4} {
			for _, owner := range []int{0, 1, 4} {
				s := updateRegistryTemporaryBase()
				p := s.Callbacks.Shop
				p.TemporaryUpdates.ItemRefs[0] = map[int]int{508: owner, 504: inventory}
				p.TemporaryUpdates.ItemWords[1][16] = 0x24
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				specs = append(specs, s)
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-missing-owners", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-missing-owners"])
}

func TestUpdateRegistryTemporaryRandomStreams(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53AEC0, legacy.PortTestTemporary53BB00, legacy.PortTestTemporary53C9A0, legacy.PortTestTemporary53D5A0, legacy.PortTestTemporary53DA60} {
		for _, seed := range []int{0, 1, 2, 17, 127, 2047, -1} {
			s := updateRegistryTemporaryBase()
			s.Seed = seed
			p := s.Callbacks.Shop
			tmp := p.TemporaryUpdates
			tmp.ItemWords[0][136] = 100
			if op == legacy.PortTestTemporary53DA60 {
				tmp.ItemWords[0][136] = 99
			}
			tmp.ItemWords[0][72] = math.Float32bits(508.125)
			tmp.ItemWords[0][76] = math.Float32bits(509.875)
			tmp.ItemRefs[0] = map[int]int{508: 1}
			p.Sequence = []legacy.PortTestShopAction{{Op: op}}
			specs = append(specs, s)
		}
	}
	results := effectsTimedRun(t, specs)
	for i, r := range results {
		want := 1
		switch specs[i].Callbacks.Shop.Sequence[0].Op {
		case legacy.PortTestTemporary53AEC0:
			want = 8
		case legacy.PortTestTemporary53C9A0:
			want = 4
		}
		if got := len(r.Lifecycle.Created); got != want {
			t.Fatalf("case %d creations %d want %d", i, got, want)
		}
	}
	updateRegistryHash(t, "update-registry-temporary-random-streams", results, updateRegistryHashes["update-registry-temporary-random-streams"])
}

func TestUpdateRegistryTemporaryAcquisition(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, mode := range []int{0, 1, 2, 3, 4, 5} {
		for _, dist := range []float32{0, 9.875, 9.9, 10, 20, 599, 600, 601} {
			for _, reverse := range []bool{false, true} {
				s := updateRegistryTemporaryBase()
				p := s.Callbacks.Shop
				tmp := p.TemporaryUpdates
				tmp.Indexed = []int{4, 5}
				if reverse {
					tmp.Indexed = []int{5, 4}
				}
				tmp.ItemRefs[0] = map[int]int{508: 1}
				tmp.ItemWords[0][136] = 90
				tmp.ItemWords[1][56] = math.Float32bits(512 + dist)
				tmp.ItemWords[1][60] = math.Float32bits(512)
				tmp.ItemWords[1][12] = 2
				tmp.ItemWords[2][56] = math.Float32bits(512 + dist + 1)
				tmp.ItemWords[2][60] = math.Float32bits(512)
				tmp.ItemWords[2][12] = 2
				switch mode {
				case 1:
					tmp.ItemWords[1][16] = 0x24
				case 2:
					tmp.ItemRefs[1] = map[int]int{508: 1}
				case 3:
					tmp.ItemWords[1][12] = 0
				case 4:
					tmp.ItemWords[2][56] = tmp.ItemWords[1][56]
				case 5:
					tmp.UpdateRefs[0] = map[int]int{4: 4}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTemporary53BB00}}
				specs = append(specs, s)
			}
		}
	}
	results := effectsTimedRun(t, specs)
	for i, r := range results {
		mode := (i / 16)
		di := (i / 2) % 8
		if mode == 0 && di == 4 {
			if got := r.Callbacks.Shop.Sequence[0].TemporaryUpdatesData[4]; got != 70001 {
				t.Fatalf("acquisition %d got %d want closest item", i, got)
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-acquisition", results, updateRegistryHashes["update-registry-temporary-acquisition"])
}

func TestUpdateRegistryTemporaryTickRates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53CB60, legacy.PortTestTemporary53CCB0, legacy.PortTestTemporary53D220, legacy.PortTestTemporary53D270, legacy.PortTestTemporary53D330, legacy.PortTestTemporary53D400, legacy.PortTestTemporary53DA60, legacy.PortTestTemporary53DCC0, legacy.PortTestTemporary53BB00, legacy.PortTestTemporary53BDA0} {
		for _, fps := range []uint32{1, 4, 30, 60} {
			for _, elapsed := range []uint32{fps - 1, fps, fps + 1, 300*fps + 1} {
				s := updateRegistryTemporaryBase()
				s.Owner.Frame = elapsed
				s.Owner.FPS = fps
				p := s.Callbacks.Shop
				tmp := p.TemporaryUpdates
				tmp.ItemWords[0][128] = 0
				tmp.ItemWords[0][136] = 0
				tmp.ItemWords[0][104] = math.Float32bits(1)
				tmp.ItemRefs[0] = map[int]int{508: 1, 504: 4}
				tmp.UpdateRefs[0] = map[int]int{0: 1}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				specs = append(specs, s)
			}
		}
	}
	updateRegistryHash(t, "update-registry-temporary-tick-rates", effectsTimedRun(t, specs), updateRegistryHashes["update-registry-temporary-tick-rates"])
}
