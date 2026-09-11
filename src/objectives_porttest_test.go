//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

var objectiveHashes = map[string]string{
	"objectives-ball-clock":          "948c1d5fb7db76a909fe0417a19a60f8508f37dbb3473adc4daab3b4dc4406e9",
	"objectives-ball-pickup":         "a7c48f6cac89cdc32d264e9d7f80f14c4cc25aa5aef3cee2bdaf560a09d9fa08",
	"objectives-buff-flags":          "a993fa36f2422000a64fe78678bf51ac694901890e37a2471a5c9097a2bc8f57",
	"objectives-crown-dispatch":      "4002891ba2b000d8237d6b5553300039c0569d371565205b60a73836cb3706c6",
	"objectives-ctf-score":           "7440311f389651da4e1e8cdb963104cc381a4f089c52c8ab939d3dc95f586729",
	"objectives-flag-deadline":       "d8e7fa15206e0024d62816b9211ea554d499e2c5abd80b5c04975b8860ade29e",
	"objectives-flag-pickup":         "040aa0d0e738401e3ff9cb4c970a4cfd9c83af0fd46dc9e3c3128c3f391d0162",
	"objectives-home-score":          "b5d9d1e7524578a21f1a2c0664205ed67262a94bdd6d539faa51665c2785e40f",
	"objectives-identity":            "6dc7437e5dc6965893db213c5bb13632c6ec98dbdd052223eee6d9c99209552c",
	"objectives-line-of-sight":       "5c48b96122719b04b2b4ea4f721d1a5631e61198ad4fa57a2052494ef54ffe41",
	"objectives-obelisk-eligibility": "6f668f39f06c6e39342723ad6fce5dfbd2a909a551af776a395ec2185f48bd25",
	"objectives-obelisk-idle":        "2189ebcd8ab9d6333727a341574fde247d25714e0ef06d3e036c30adbc83881a",
	"objectives-obelisk-transfer":    "51f1f1b20311169cf2347d7593bd23b5e2e883503a67120f240fedc23679d0cb",
	"objectives-obelisk-wands":       "3688abbddfa3a3e703cf16e0533464c08bbf117773392ee58c56d5cc5605e23e",
	"objectives-pickup-buffs":        "27da0f7f858bdddec605801b1a769469e7c7f5f4f2b26f19a060cee7b595c378",
	"objectives-positive":            "719eeeaea002ca25bb327e9412b6abb8ce5eeb9b4baaf9aaab0718a039ec7d81",
	"objectives-possession":          "5d2a11d30b1cfe1a4cd8e78e61e113d224d20aadb9032c4147548ee47482449b",
	"objectives-remember-owner":      "171d947dddb758f86a8ae1a273ee78002eeba55475dd1da8c7c7495d280d09a8",
	"objectives-reset":               "3904acad0823e6ad39f88669fc63b6542c1d0e69c66266bbd84cfe79b93b4b22",
}

func objectiveBase() legacy.PortTestRoamSpec {
	s := worldBase()
	p := s.Callbacks.Shop
	w := p.TemporaryUpdates
	s.Lifecycle.GameFlags = 1
	p.EffectsUse.Target = 1
	w.Target = 100
	w.World.Objectives = &legacy.PortTestObjectivesSpec{Players: 1, ObjectList: []int{100, 101, 102}, ManaMultipliers: [3]float32{10, 20, 15}, PlayerWords: make([]map[int]uint32, 3), PlayerUpdateWords: make([]map[int]uint32, 3), PlayerDataWords: make([]map[int]uint32, 3), PlayerRefs: make([]map[int]int, 3), PlayerUpdateRefs: make([]map[int]int, 3), PlayerDataRefs: make([]map[int]int, 3)}
	for i := 0; i < 3; i++ {
		w.World.Objectives.PlayerWords[i] = map[int]uint32{56: math.Float32bits(512), 60: math.Float32bits(512), 176: math.Float32bits(5)}
	}
	return s
}
func objectiveHash(t *testing.T, name string, s []legacy.PortTestRoamSpec) {
	t.Helper()
	callbackHash(t, name, effectsTimedRun(t, s), objectiveHashes[name])
}
func TestObjectivesIdentity(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{806, 807} {
		for _, name := range []string{"", "RedFlag", "BlueFlag", "GoldFlag", "redflag", "unknown"} {
			for _, class := range []uint32{1, 0x10000000, 0x10000001} {
				s := objectiveBase()
				p := s.Callbacks.Shop
				p.Items[0].Class = class
				p.TemporaryUpdates.World.Objectives.ColorNames = []string{name}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
				cases = append(cases, s)
			}
		}
	}
	objectiveHash(t, "objectives-identity", cases)
}
func TestObjectivesFlagDeadline(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, fps := range []uint32{2, 30, 60} {
		for _, stamp := range []uint32{0, 1, 100, 0xfffffffa} {
			for _, delta := range []uint32{0, 1, 30*fps - 1, 30 * fps, 30*fps + 1, 0xffffffff} {
				s := objectiveBase()
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
	objectiveHash(t, "objectives-flag-deadline", cases)
}
func TestObjectivesBallClock(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, tick := range []uint64{0, 1, 19999, 20000, 20001, 0x7fffffff, 0x80000000, 0xffffffff, 0x100000000, 0x100004e20} {
		for _, stamp := range []uint64{0, 1, 20000, 0xffffffff, 0x100000000} {
			for _, speed := range []float32{0, 1, 5} {
				s := objectiveBase()
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
	objectiveHash(t, "objectives-ball-clock", cases)
}
func TestObjectivesRememberOwner(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, target := range []int{0, 100, 101, 4} {
		for _, team := range []byte{0, 1, 2} {
			for _, frame := range []uint32{0, 100, 0xffffffff} {
				s := objectiveBase()
				p := s.Callbacks.Shop
				s.Owner.Frame = frame
				p.Inventory.Teams[0] = team
				p.TemporaryUpdates.Target = target
				p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: 0, 4: 77, 16: 123}
				p.Sequence = []legacy.PortTestShopAction{{Op: 802}}
				cases = append(cases, s)
			}
		}
	}
	objectiveHash(t, "objectives-remember-owner", cases)
}
func TestObjectivesObeliskIdle(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, fps := range []uint32{2, 3, 30, 60} {
		for _, frame := range []uint32{0, 1, 14, 15, 16, 30, 0xffffffff} {
			for _, energy := range []uint32{0, 1, 7, 8, 49, 50, 51, 0xffffffff} {
				s := objectiveBase()
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
	objectiveHash(t, "objectives-obelisk-idle", cases)
}
func TestObjectivesPossession(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{810, 811} {
		for _, owner := range []int{0, 100} {
			for _, remembered := range []int{0, 100, 101} {
				for _, flags := range []uint32{0, 0x20, 0x8000} {
					for _, delta := range []uint32{0, 9, 10, 11, 0xffffffff} {
						s := objectiveBase()
						p := s.Callbacks.Shop
						w := p.TemporaryUpdates
						o := w.World.Objectives
						s.Owner.Frame = 100 + delta
						o.Ticks = []uint64{100, 101}
						w.UpdateWords[0] = map[int]uint32{8: 100, 16: 100, 20: 10, 24: math.Float32bits(2)}
						w.UpdateRefs[0] = map[int]int{0: remembered}
						o.PlayerWords[0][16] = flags
						if owner != 0 {
							w.ItemRefs[0] = map[int]int{508: owner}
							o.PlayerRefs[0] = map[int]int{516: 3}
						}
						p.Sequence = []legacy.PortTestShopAction{{Op: op}}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	objectiveHash(t, "objectives-possession", cases)
}
func TestObjectivesObeliskTransfer(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mode := range []uint32{0, 4096, 8192, 12288} {
		for _, energy := range []uint32{0, 1, 2, 7, 8, 49, 50} {
			for _, mana := range []uint16{0, 49, 50} {
				for _, distance := range []float32{0, 49.999, 50, 50.001} {
					s := objectiveBase()
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
	objectiveHash(t, "objectives-obelisk-transfer", cases)
}
func TestObjectivesReset(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, n := range []int{0, 1, 2} {
		for _, old := range []uint32{0, 1} {
			for _, missing := range []bool{false, true} {
				for _, tick := range []uint64{0, 12345, 0xffffffff, 0x100000001} {
					s := objectiveBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					o := w.World.Objectives
					o.Ticks = []uint64{tick}
					w.World.ItemNames = []string{"gameball", "gameballstart", "gameballstart"}
					if missing {
						o.MissingTypes = []string{"gameball"}
						w.World.ItemNames[0] = ""
					}
					for i := 0; i < n; i++ {
						o.ObjectList = append(o.ObjectList, 4+i)
						w.ItemWords[1+i][56] = math.Float32bits(600 + float32(i)*100)
						w.ItemWords[1+i][60] = math.Float32bits(650)
					}
					o.PlayerDataRefs[0] = map[int]int{3628: 3}
					p.EffectsUse.Balance["FlagballPossDuration"] = []float64{10}
					p.EffectsUse.Balance["FlagballResetVel"] = []float64{2}
					p.Sequence = []legacy.PortTestShopAction{{Op: 800, Value: old}}
					cases = append(cases, s)
				}
			}
		}
	}
	objectiveHash(t, "objectives-reset", cases)
}
func TestObjectivesBallPickup(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, target := range []int{0, 100, 4} {
		for _, team := range []byte{0, 1, 2} {
			for _, remember := range []int{0, 100} {
				for _, owned := range []bool{false, true} {
					for _, frame := range []uint32{44, 45, 46, 100} {
						s := objectiveBase()
						p := s.Callbacks.Shop
						w := p.TemporaryUpdates
						o := w.World.Objectives
						s.Owner.Frame = frame
						o.Players = 3
						p.Inventory.TeamMembers = true
						p.Inventory.Teams = [4]byte{team, team, 2, 0}
						w.Target = target
						w.UpdateRefs[0] = map[int]int{0: remember}
						if owned {
							w.ItemRefs[0] = map[int]int{508: 100}
							o.PlayerRefs[0] = map[int]int{516: 3}
						}
						p.Sequence = []legacy.PortTestShopAction{{Op: 803}}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	objectiveHash(t, "objectives-ball-pickup", cases)
}
func TestObjectivesCrownDispatch(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{804, 811} {
		for _, target := range []int{0, 100, 4} {
			for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
				for _, accept := range []bool{false, true} {
					s := objectiveBase()
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
	objectiveHash(t, "objectives-crown-dispatch", cases)
}
func TestObjectivesFlagPickup(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{801, 812} {
		for _, mode := range []uint32{0, 32, 64, 96} {
			for _, teams := range [][4]byte{{1, 2, 2, 1}, {1, 2, 2, 2}, {0, 0, 0, 2}, {1, 1, 1, 2}} {
				for _, away := range []bool{false, true} {
					s := objectiveBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					o := w.World.Objectives
					s.Lifecycle.GameFlags = mode | 1
					p.Inventory.TeamMembers = true
					p.Inventory.Teams = teams
					o.Players = 3
					p.Items[0].Class = 0x10000000
					o.ColorNames = []string{"RedFlag"}
					home := float32(512)
					if away {
						home = 600
					}
					w.UpdateWords[0] = map[int]uint32{0: math.Float32bits(home), 4: math.Float32bits(512), 8: 50}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}}
					cases = append(cases, s)
				}
			}
		}
	}
	objectiveHash(t, "objectives-flag-pickup", cases)
}
func TestObjectivesHomeScore(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{805, 814, 812} {
		for _, team := range []byte{1, 2} {
			for _, scorer := range []int{0, 100, 101} {
				for _, dead := range []bool{false, true} {
					for _, held := range []bool{false, true} {
						s := objectiveBase()
						p := s.Callbacks.Shop
						w := p.TemporaryUpdates
						o := w.World.Objectives
						s.Lifecycle.GameFlags = 65
						o.Players = 3
						p.Inventory.TeamMembers = true
						p.Inventory.Teams = [4]byte{1, 2, 2, team}
						w.World.ItemNames = []string{"objectiveflag", "gameball", "gameballstart"}
						o.ObjectList = append(o.ObjectList, 5)
						w.Target = 4
						w.UpdateRefs[1] = map[int]int{0: scorer}
						w.ItemWords[2][56] = math.Float32bits(650)
						w.ItemWords[2][60] = math.Float32bits(675)
						w.ItemWords[1][80] = math.Float32bits(3)
						w.ItemWords[1][84] = math.Float32bits(4)
						if dead && scorer >= 100 {
							o.PlayerWords[scorer-100][16] = 0x20
						}
						if held {
							w.ItemRefs[1] = map[int]int{508: 100}
							o.PlayerRefs[0] = map[int]int{516: 4}
							if op != 805 {
								w.Target = 100
							}
						}
						p.Sequence = []legacy.PortTestShopAction{{Op: op}}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	objectiveHash(t, "objectives-home-score", cases)
}
func TestObjectivesCTFScore(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{801, 812} {
		for _, limit := range []uint16{0, 1, 2} {
			for _, color := range []string{"BlueFlag", "GoldFlag"} {
				s := objectiveBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				o := w.World.Objectives
				s.Lifecycle.GameFlags = 33
				p.Inventory.TeamMembers = true
				p.Inventory.Teams = [4]byte{1, 2, 2, 1}
				o.Players = 3
				o.ScoreLimit = limit
				o.ColorNames = []string{"RedFlag", color}
				p.Items[0].Class = 0x10000000
				p.Items[1].Class = 0x10000000
				p.Inventory.Linked = []int{1}
				w.UpdateWords[0] = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
				w.UpdateWords[1] = map[int]uint32{0: math.Float32bits(650), 4: math.Float32bits(675), 8: 50}
				w.ItemWords[1][52] = 2
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				cases = append(cases, s)
			}
		}
	}
	r := effectsTimedRun(t, cases)
	for i, x := range r {
		_, _, pl := objectivePlayerData(x.Callbacks.Shop.Sequence[0], 0)
		if pl[2136/4] != 1 {
			t.Fatalf("CTF case %d score=%d want 1", i, pl[2136/4])
		}
		it := worldObject(t, x, 0, 70001)
		if it[492/4] != 0 || it[56/4] != math.Float32bits(650) || it[60/4] != math.Float32bits(675) {
			t.Fatalf("CTF case %d flag return", i)
		}
	}
	callbackHash(t, "objectives-ctf-score", r, objectiveHashes["objectives-ctf-score"])
}
func TestObjectivesPickupBuffs(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, buffs := range []uint32{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 1 << 17, 1 << 18, 1 << 19, 1 << 20, 1 << 21, 1 << 22, 1 << 23, 1 << 24, 1 << 25, 1 << 26, 1 << 27, 1 << 28, 1 << 29, 1 << 30, 1 << 31, 0xffffffff} {
		s := objectiveBase()
		p := s.Callbacks.Shop
		p.Resources.Buffs = buffs
		p.Sequence = []legacy.PortTestShopAction{{Op: 813}, {Op: 813}}
		cases = append(cases, s)
	}
	objectiveHash(t, "objectives-pickup-buffs", cases)
}

// Decode captured fixture state for independent checks of observable outcomes.
func objectivePlayerData(step legacy.PortTestShopStep, player int) (unit, update, data []uint32) {
	d := step.TemporaryUpdatesData
	off := 3 + int(d[2]) + 7 + 3*16
	off += 7 + int(d[off+6]) + 3*16
	off += 1 + int(d[off])                // recorded transport
	off += 1 + 2*int(d[off]) + 2 + 52 + 4 // ticks, caches, status
	udn := int(unsafe.Sizeof(server.PlayerUpdateData{})) / 4
	pln := int(unsafe.Sizeof(server.Player{})) / 4
	off += player * (193 + udn + pln)
	return d[off : off+193], d[off+193 : off+193+udn], d[off+193+udn : off+193+udn+pln]
}
func objectiveItemUpdate(step legacy.PortTestShopStep, item int) []uint32 {
	d := step.TemporaryUpdatesData
	off := 10 + int(d[2]) + item*16
	return d[off : off+16]
}
func TestObjectivesPositiveEffects(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	s := objectiveBase()
	p := s.Callbacks.Shop
	p.Resources.Mana = 49
	p.Resources.MaxMana = 50
	p.TemporaryUpdates.UpdateWords[0][0] = 8
	p.Sequence = []legacy.PortTestShopAction{{Op: 808}}
	cases = append(cases, s)
	s = objectiveBase()
	p = s.Callbacks.Shop
	p.Resources.Mana = 49
	p.Resources.MaxMana = 50
	p.TemporaryUpdates.UpdateWords[0][0] = 2
	p.TemporaryUpdates.World.Objectives.PlayerUpdateRefs[0] = map[int]int{104: 4}
	p.Items[1].Class = 0x1000
	p.Items[1].Subclass = 0x4000000
	p.TemporaryUpdates.World.Objectives.UseWords = []map[int]uint32{nil, {108: 19 | 20<<8, 112: 95}}
	p.EffectsUse.Balance["OblivionStaffRechargeRate"] = []float64{2}
	p.Sequence = []legacy.PortTestShopAction{{Op: 808}}
	cases = append(cases, s)
	s = objectiveBase()
	p = s.Callbacks.Shop
	p.Items[0].Class = 0x10000000
	p.TemporaryUpdates.World.Objectives.ColorNames = []string{"BlueFlag"}
	p.Sequence = []legacy.PortTestShopAction{{Op: 806}}
	cases = append(cases, s)
	s = objectiveBase()
	p = s.Callbacks.Shop
	p.Inventory.TeamMembers = true
	p.Inventory.Teams = [4]byte{1, 2, 2, 0}
	p.TemporaryUpdates.World.ItemNames = []string{"gameball"}
	p.Sequence = []legacy.PortTestShopAction{{Op: 803}}
	cases = append(cases, s)
	s = objectiveBase()
	p = s.Callbacks.Shop
	s.Owner.Frame = 902
	p.Items[0].Class = 0x10000000
	p.TemporaryUpdates.World.Objectives.ColorNames = []string{"RedFlag"}
	p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: math.Float32bits(600), 4: math.Float32bits(650), 8: 1}
	p.Sequence = []legacy.PortTestShopAction{{Op: 809}}
	cases = append(cases, s)
	r := effectsTimedRun(t, cases)
	for i, want := range []uint32{7, 0} {
		step := r[i].Callbacks.Shop.Sequence[0]
		_, ud, _ := objectivePlayerData(step, 0)
		if got := ud[1] & 0xffff; got != 50 {
			t.Fatalf("mana transfer %d: got %d want 50", i, got)
		}
		if got := objectiveItemUpdate(step, 0)[0]; got != want {
			t.Fatalf("energy transfer %d: got %d want %d", i, got, want)
		}
	}
	if got := r[2].Callbacks.Shop.Sequence[0].Return; got != 2 {
		t.Fatalf("flag identity %d want 2", got)
	}
	ball := worldObject(t, r[3], 0, 70000)
	if ball[508/4] != 54000 || ball[16/4]&0x40 == 0 {
		t.Fatalf("ball pickup owner=%d flags=%x", ball[508/4], ball[16/4])
	}
	flag := worldObject(t, r[4], 0, 70000)
	if flag[56/4] != math.Float32bits(600) || flag[60/4] != math.Float32bits(650) || objectiveItemUpdate(r[4].Callbacks.Shop.Sequence[0], 0)[2] != 0 {
		t.Fatal("flag did not return home and clear deadline")
	}
	callbackHash(t, "objectives-positive", r, objectiveHashes["objectives-positive"])
}
func TestObjectivesObeliskWands(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mode := range []uint32{0, 4096, 8192, 12288} {
		for _, energy := range []uint32{0, 1, 2, 8} {
			for _, charges := range []uint32{0, 19, 20} {
				for _, rate := range []float64{0, 1, 2, 5} {
					s := objectiveBase()
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
	objectiveHash(t, "objectives-obelisk-wands", cases)
}

func TestObjectivesPickupBuffFlags(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var want []uint32
	for _, flags := range []uint32{0, 0x80000, 0x80001} {
		for _, buffs := range []uint32{0, 1, 1 << 7, 1 << 18, 1 << 30, 0xffffffff} {
			s := objectiveBase()
			p := s.Callbacks.Shop
			p.Resources.Buffs = buffs
			remaining := buffs
			for enc := 0; enc < 32; enc++ {
				id := server.EnchantID(enc).Spell()
				if id != 0 {
					p.TemporaryUpdates.World.Objectives.SpellDefinitions = append(p.TemporaryUpdates.World.Objectives.SpellDefinitions, server.PortTestSpellClassDef{Index: uint32(id), Valid: true, Flags: flags})
					if flags&0x80000 != 0 {
						remaining &^= 1 << enc
					}
				}
			}
			p.Sequence = []legacy.PortTestShopAction{{Op: 813}, {Op: 813}}
			cases = append(cases, s)
			want = append(want, remaining)
		}
	}
	r := effectsTimedRun(t, cases)
	for i, x := range r {
		u, _, _ := objectivePlayerData(x.Callbacks.Shop.Sequence[0], 0)
		if u[340/4] != want[i] {
			t.Fatalf("buff case %d got %x want %x", i, u[340/4], want[i])
		}
	}
	callbackHash(t, "objectives-buff-flags", r, objectiveHashes["objectives-buff-flags"])
}
func TestObjectivesObeliskEligibility(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, team := range []byte{0, 1, 2} {
		for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
			for _, class := range []byte{0, 1, 2, 3} {
				for _, frame := range []uint32{0, 14, 15, 16, 0xffffffff} {
					s := objectiveBase()
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
	objectiveHash(t, "objectives-obelisk-eligibility", cases)
}
func TestObjectivesLineOfSight(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{808, 810, 811} {
		for _, wall := range []int{0, 1, 2} {
			s := objectiveBase()
			p := s.Callbacks.Shop
			w := p.TemporaryUpdates
			o := w.World.Objectives
			p.Inventory.WallMode = wall
			p.Resources.Mana = 49
			p.Resources.MaxMana = 50
			w.ItemWords[0][56] = math.Float32bits(120)
			w.ItemWords[0][60] = math.Float32bits(100)
			o.PlayerWords[0][56] = math.Float32bits(160)
			o.PlayerWords[0][60] = math.Float32bits(100)
			w.UpdateWords[0][0] = 8
			if op != 808 {
				w.UpdateWords[0] = map[int]uint32{8: 100, 20: 10}
				w.ItemRefs[0] = map[int]int{508: 100}
				o.PlayerRefs[0] = map[int]int{516: 3}
				o.PlayerWords[0][56] = math.Float32bits(100)
				o.Ticks = []uint64{100}
			}
			p.Sequence = []legacy.PortTestShopAction{{Op: op}}
			cases = append(cases, s)
		}
	}
	r := effectsTimedRun(t, cases)
	_, open, _ := objectivePlayerData(r[0].Callbacks.Shop.Sequence[0], 0)
	_, blocked, _ := objectivePlayerData(r[1].Callbacks.Shop.Sequence[0], 0)
	if open[1]&0xffff != 50 || blocked[1]&0xffff != 49 {
		t.Fatalf("obelisk LOS open mana=%d blocked=%d", open[1]&0xffff, blocked[1]&0xffff)
	}
	callbackHash(t, "objectives-line-of-sight", r, objectiveHashes["objectives-line-of-sight"])
}
