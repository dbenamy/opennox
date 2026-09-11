//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

var temporaryHashes = map[string]string{
	"temporary-acquisition":    "23f800d848c47a0084539f7296acac45e0c961785db3b310e1c929cbadccd620",
	"temporary-area-callbacks": "819fb7b42eccaa7620356f0b6fc26d893d08ed11aeeb7c08bb0b1712ffeec132",
	"temporary-area-owners":    "e847ce54985f5f63863019f898a82b0233c989c8aa6a9b3986cf8d4ee152b51d",
	"temporary-break":          "fc108ef077e38a5e9add745b4ddebe44205542af66d1acf264649e74615a3301",
	"temporary-frame-wrap":     "e2a94f1b18934a56c1fda975fc23f66aa9d74190c23d2816b67bfa101290bfeb",
	"temporary-homing":         "bfe1328dbc18848552a38ea6b67d4503250b2198a5ed6521c25d86bf09f274ef",
	"temporary-impact":         "138ebbb8f2c08e608444b2768fff2c32d75f6ec239125c670ff90e395aba7742",
	"temporary-lifetime":       "9af722e94b2f49e5fa0c20838d11d6b08dbf477893196d3b7145ab6d7b0c8cc8",
	"temporary-missing-owners": "7eaaf2179a87626cf8070cac2f568d6ea39729ec44ce7b20193a8cd4c9111231",
	"temporary-owner-effects":  "bf8c2c58c3341e38b2614f9d56cad66a31bbbcd28d1df272bb19dd6c92521667",
	"temporary-positive":       "2db9ae4931a3e239144749592c7fa49209fb0a6dc681462f01d6d9e60ea2cdc9",
	"temporary-random-streams": "7a8f5a17662714a284d092bca069c7a41f0ad42114682320a1ee9eb00e084a42",
	"temporary-spark":          "06a5b0b427a1658f4253cfd5547b87b16ddce24abf01d5cfb8c5757e37af840a",
	"temporary-spawners":       "77d9b032a437f721b226ff4d0576f7c0c6444f6e20313139068886eaa44a6b1e",
	"temporary-tick-rates":     "73ed23edd8a225b7c566e9ee3987ba66503e7738f514ae86ce94bd00e994393e",
	"temporary-trails":         "6a8b2595f571ee5668c8c29beb434834e0ee6364013d68d5b3f74d653301e9cf",
}

func temporaryBase() legacy.PortTestRoamSpec {
	s := effectsUseBase()
	p := s.Callbacks.Shop
	s.Owner.Frame = 100
	p.Inventory.Linked = nil
	p.Inventory.Owned = nil
	p.Inventory.Position.X = 512
	p.Inventory.Position.Y = 512
	p.Inventory.Target.X = 512
	p.Inventory.Target.Y = 512
	p.EffectsUse.Projectiles = true
	p.EffectsUse.Balance = map[string][]float64{"TargetedSpellLifetime": {120}, "UnTargetedSpellLifetime": {60}}
	p.TemporaryUpdates = &legacy.PortTestTemporaryUpdatesSpec{Target: 1, ItemWords: make([]map[int]uint32, 3), UpdateWords: make([]map[int]uint32, 3), ItemRefs: make([]map[int]int, 3), UpdateRefs: make([]map[int]int, 3)}
	for i := range p.Items {
		p.Items[i].Class = 1
		p.Items[i].Subclass = 0
		p.Items[i].Flags = 4
		p.TemporaryUpdates.ItemWords[i] = map[int]uint32{128: 90, 136: 90, 176: math.Float32bits(5), 544: math.Float32bits(2.5)}
	}
	return s
}
func TestTemporaryLifetimeStates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53ADC0, legacy.PortTestTemporary53B8F0, legacy.PortTestTemporary53CB60, legacy.PortTestTemporary53CC90, legacy.PortTestTemporary53DB00} {
		for _, frame := range []uint32{0, 1, 2, 3, 9, 10, 11, 29, 30, 31, 100, 0xffffffff} {
			for _, life := range []uint32{0, 1, 10, 0xffffffff} {
				for _, cb := range []bool{false, true} {
					s := temporaryBase()
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
	callbackHash(t, "temporary-lifetime", effectsTimedRun(t, specs), temporaryHashes["temporary-lifetime"])
}
func TestTemporaryBreakStates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53DB30, legacy.PortTestTemporary53DBB0, legacy.PortTestTemporary53DC30} {
		for _, state := range []uint32{0, 1, 2, 3, 4, 6, 8, 10, 12, 0xffffffff} {
			for _, flags := range []uint32{4, 0x8004, 0x8044} {
				for _, frame := range []uint32{0, 99, 100, 101, 0xffffffff} {
					s := temporaryBase()
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
	callbackHash(t, "temporary-break", effectsTimedRun(t, specs), temporaryHashes["temporary-break"])
}
func TestTemporarySparkCreation(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, stage := range []uint32{0, 1, 2, 3, 4, 5, 0xffffffff} {
		for _, life := range []int{0, 1, 20, -1} {
			for _, owner := range []int{0, 1} {
				s := temporaryBase()
				p := s.Callbacks.Shop
				p.TemporaryUpdates.ItemRefs[0] = map[int]int{508: owner}
				p.TemporaryUpdates.ItemWords[0][80] = math.Float32bits(.125)
				p.TemporaryUpdates.ItemWords[0][84] = math.Float32bits(-1.25)
				p.TemporaryUpdates.ItemWords[0][108] = math.Float32bits(2.5)
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTemporary54FD80, Value: stage, Side: life}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "temporary-spark", effectsTimedRun(t, specs), temporaryHashes["temporary-spark"])
}

func TestTemporaryTrails(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, dir := range []uint32{0, 1, 63, 64, 127, 128, 255} {
		for _, speed := range []float32{0, .125, 2.5, 100} {
			s := temporaryBase()
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
	callbackHash(t, "temporary-trails", results, temporaryHashes["temporary-trails"])
}
func TestTemporaryHoming(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53B940, legacy.PortTestTemporary53BB00, legacy.PortTestTemporary53BDA0, legacy.PortTestTemporary53DCC0} {
		for _, mode := range []int{0, 1, 2, 3} {
			for _, dir := range []uint32{0, 1, 63, 127, 248, 255} {
				for _, frame := range []uint32{96, 100, 200} {
					s := temporaryBase()
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
	callbackHash(t, "temporary-homing", effectsTimedRun(t, specs), temporaryHashes["temporary-homing"])
}
func TestTemporaryOwnerEffects(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53D270, legacy.PortTestTemporary53D330, legacy.PortTestTemporary53D400, legacy.PortTestTemporary53D510} {
		for _, frame := range []uint32{90, 91, 94, 100, 700, 10000} {
			for _, state := range []uint32{0, 0x20, 0x8000} {
				for _, wall := range []bool{false, true} {
					s := temporaryBase()
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
	callbackHash(t, "temporary-owner-effects", effectsTimedRun(t, specs), temporaryHashes["temporary-owner-effects"])
}
func TestTemporarySpawners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53C9A0, legacy.PortTestTemporary53D5A0, legacy.PortTestTemporary53D6E0, legacy.PortTestTemporary53DA60} {
		for _, frame := range []uint32{89, 90, 91, 100, 150, 241} {
			for _, stamp := range []uint32{90, 100, 101} {
				for _, wall := range []bool{false, true} {
					s := temporaryBase()
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
	callbackHash(t, "temporary-spawners", effectsTimedRun(t, specs), temporaryHashes["temporary-spawners"])
}
func TestTemporaryAreaCallbacks(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53CC30, legacy.PortTestTemporary53BD10, legacy.PortTestTemporary53D8C0, legacy.PortTestTemporary53D9D0} {
		for _, class := range []uint32{1, 2, 0x2000, 0x100000} {
			for _, dist := range []float32{0, 39, 40, 45, 45.000004, 100, 600} {
				for _, wall := range []bool{false, true} {
					s := temporaryBase()
					p := s.Callbacks.Shop
					tmp := p.TemporaryUpdates
					tmp.Target = 4
					p.Items[1].Class = class
					tmp.ItemWords[1][56] = math.Float32bits(512 + dist)
					tmp.ItemWords[1][60] = math.Float32bits(512)
					tmp.ItemRefs[0] = map[int]int{508: 1}
					tmp.ItemWords[1][12] = 2
					if wall {
						p.Inventory.WallMode = 1
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "temporary-area-callbacks", effectsTimedRun(t, specs), temporaryHashes["temporary-area-callbacks"])
}
func TestTemporaryAreaOwners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53CB90, legacy.PortTestTemporary53CCB0, legacy.PortTestTemporary53D220, legacy.PortTestTemporary53D850, legacy.PortTestTemporary53D960} {
		for _, frame := range []uint32{90, 93, 98, 100, 101, 150, 151} {
			for _, life := range []uint32{0, 1, 2} {
				for _, indexed := range []bool{false, true} {
					s := temporaryBase()
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
	callbackHash(t, "temporary-area-owners", effectsTimedRun(t, specs), temporaryHashes["temporary-area-owners"])
}

func init() {
	legacy.PortTestTemporaryServer = func(core *server.Server) func() {
		old := core.ExtServer
		oldServer := noxServer
		adapter := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(adapter)
		noxServer = adapter
		return func() { core.ExtServer = old; noxServer = oldServer }
	}
}

func TestTemporaryPositiveCallbacks(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53CC30, legacy.PortTestTemporary53BD10, legacy.PortTestTemporary53D8C0, legacy.PortTestTemporary53D9D0, legacy.PortTestTemporary53CB90} {
		s := temporaryBase()
		p := s.Callbacks.Shop
		tmp := p.TemporaryUpdates
		tmp.Target = 4
		p.Items[1].Class = 2
		if op == legacy.PortTestTemporary53CC30 || op == legacy.PortTestTemporary53CB90 {
			p.Items[1].Class = 0x2000
		}
		tmp.ItemWords[1][56] = math.Float32bits(520)
		tmp.ItemWords[1][60] = math.Float32bits(512)
		tmp.ItemRefs[0] = map[int]int{508: 1}
		if op == legacy.PortTestTemporary53CB90 {
			s.Owner.Frame = 98
			tmp.Indexed = []int{4}
		}
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		specs = append(specs, s)
	}
	results := effectsTimedRun(t, specs)
	for i, r := range results {
		d := r.Callbacks.Shop.Sequence[0].TemporaryUpdatesData
		switch i {
		case 0, 4:
			if d[3] != 1 {
				t.Fatalf("water case %d flag %d want 1", i, d[3])
			}
		case 1:
			if d[4] != 70001 {
				t.Fatalf("nearest candidate %d want 70001", d[4])
			}
		case 2, 3:
			if len(r.Callbacks.Damage) != 5 || r.Callbacks.Damage[0] != 70001 {
				t.Fatalf("damage case %d: %v", i, r.Callbacks.Damage)
			}
		}
	}
	callbackHash(t, "temporary-positive", results, temporaryHashes["temporary-positive"])
}
func TestTemporaryFrameWrap(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53B8F0, legacy.PortTestTemporary53CB60, legacy.PortTestTemporary53CC90, legacy.PortTestTemporary53CCB0, legacy.PortTestTemporary53D220, legacy.PortTestTemporary53DB00, legacy.PortTestTemporary53D850, legacy.PortTestTemporary53DA60} {
		for _, start := range []uint32{0, 0xfffffff0, 0x80000000} {
			for _, elapsed := range []uint32{0, 2, 3, 10, 30, 31, 61, 0x80000000, 0xffffffff} {
				s := temporaryBase()
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
	callbackHash(t, "temporary-frame-wrap", effectsTimedRun(t, specs), temporaryHashes["temporary-frame-wrap"])
}
func TestTemporaryMissingAndOwners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary54FD80, legacy.PortTestTemporary53AEC0, legacy.PortTestTemporary53BB00} {
		s := temporaryBase()
		p := s.Callbacks.Shop
		p.TemporaryUpdates.MissingTypes = []string{"Spark"}
		p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1, Side: 20}}
		specs = append(specs, s)
	}
	for _, op := range []int{legacy.PortTestTemporary53D270, legacy.PortTestTemporary53DCC0} {
		for _, inventory := range []int{0, 4} {
			for _, owner := range []int{0, 1, 4} {
				s := temporaryBase()
				p := s.Callbacks.Shop
				p.TemporaryUpdates.ItemRefs[0] = map[int]int{508: owner, 504: inventory}
				p.TemporaryUpdates.ItemWords[1][16] = 0x24
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "temporary-missing-owners", effectsTimedRun(t, specs), temporaryHashes["temporary-missing-owners"])
}

func TestTemporaryRandomStreams(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53AEC0, legacy.PortTestTemporary53BB00, legacy.PortTestTemporary53C9A0, legacy.PortTestTemporary53D5A0, legacy.PortTestTemporary53DA60} {
		for _, seed := range []int{0, 1, 2, 17, 127, 2047, -1} {
			s := temporaryBase()
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
	callbackHash(t, "temporary-random-streams", results, temporaryHashes["temporary-random-streams"])
}
func TestTemporaryAcquisition(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, mode := range []int{0, 1, 2, 3, 4, 5} {
		for _, dist := range []float32{0, 9.875, 9.9, 10, 20, 599, 600, 601} {
			for _, reverse := range []bool{false, true} {
				s := temporaryBase()
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
	callbackHash(t, "temporary-acquisition", results, temporaryHashes["temporary-acquisition"])
}
func TestTemporaryCollisionAndImpact(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{4, 0x80000004} {
		for _, height := range []float32{-1, 0, .000001, 199.99998, 200, 200.00002} {
			s := temporaryBase()
			p := s.Callbacks.Shop
			p.TemporaryUpdates.ItemWords[0][16] = flags
			p.TemporaryUpdates.ItemWords[0][104] = math.Float32bits(height)
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTemporary53D400}, {Op: legacy.PortTestTemporary53D400}}
			specs = append(specs, s)
		}
	}
	for _, ownerFlags := range []uint32{4, 0x24} {
		for _, ret := range []uint32{0, 1, 0xffffffff, 0x81234567} {
			s := temporaryBase()
			p := s.Callbacks.Shop
			p.TemporaryUpdates.CollideReturn = ret
			p.TemporaryUpdates.UpdateRefs[0] = map[int]int{0: 1}
			p.EffectsUse.UnitWords = map[int]uint32{16: ownerFlags}
			s.Owner.Frame = 200
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTemporary53BDA0}}
			specs = append(specs, s)
		}
	}
	results := effectsTimedRun(t, specs)
	for i := 12; i < len(results); i++ {
		if got := results[i].Callbacks.Shop.Sequence[0].Return; got != specs[i].Callbacks.Shop.TemporaryUpdates.CollideReturn {
			t.Fatalf("collision case %d result %x", i, got)
		}
	}
	callbackHash(t, "temporary-impact", results, temporaryHashes["temporary-impact"])
}

func TestTemporaryTickRates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestTemporary53CB60, legacy.PortTestTemporary53CCB0, legacy.PortTestTemporary53D220, legacy.PortTestTemporary53D270, legacy.PortTestTemporary53D330, legacy.PortTestTemporary53D400, legacy.PortTestTemporary53DA60, legacy.PortTestTemporary53DCC0, legacy.PortTestTemporary53BB00, legacy.PortTestTemporary53BDA0} {
		for _, fps := range []uint32{1, 4, 30, 60} {
			for _, elapsed := range []uint32{fps - 1, fps, fps + 1, 300*fps + 1} {
				s := temporaryBase()
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
	callbackHash(t, "temporary-tick-rates", effectsTimedRun(t, specs), temporaryHashes["temporary-tick-rates"])
}
