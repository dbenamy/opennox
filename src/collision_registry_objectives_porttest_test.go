//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestCollisionRegistryObjectivesBallPickup(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, target := range []int{0, 100, 4} {
		for _, team := range []byte{0, 1, 2} {
			for _, remember := range []int{0, 100} {
				for _, owned := range []bool{false, true} {
					for _, frame := range []uint32{44, 45, 46, 100} {
						s := collisionRegistryObjectiveBase()
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
	name := "collision-registry-objectives-ball-pickup"
	collisionRegistryHash(t, name, effectsTimedRun(t, cases), collisionRegistryOwnerHashes[name])
}
func TestCollisionRegistryObjectivesCrownDispatch(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{804, 811} {
		for _, target := range []int{0, 100, 4} {
			for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
				for _, accept := range []bool{false, true} {
					s := collisionRegistryObjectiveBase()
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
	name := "collision-registry-objectives-crown-dispatch"
	collisionRegistryHash(t, name, effectsTimedRun(t, cases), collisionRegistryOwnerHashes[name])
}
func TestCollisionRegistryObjectivesFlagPickup(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{801, 812} {
		for _, mode := range []uint32{0, 32, 64, 96} {
			for _, teams := range [][4]byte{{1, 2, 2, 1}, {1, 2, 2, 2}, {0, 0, 0, 2}, {1, 1, 1, 2}} {
				for _, away := range []bool{false, true} {
					s := collisionRegistryObjectiveBase()
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
	name := "collision-registry-objectives-flag-pickup"
	collisionRegistryHash(t, name, effectsTimedRun(t, cases), collisionRegistryOwnerHashes[name])
}
func TestCollisionRegistryObjectivesHomeScore(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{805, 814, 812} {
		for _, team := range []byte{1, 2} {
			for _, scorer := range []int{0, 100, 101} {
				for _, dead := range []bool{false, true} {
					for _, held := range []bool{false, true} {
						s := collisionRegistryObjectiveBase()
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
	name := "collision-registry-objectives-home-score"
	collisionRegistryHash(t, name, effectsTimedRun(t, cases), collisionRegistryOwnerHashes[name])
}
func TestCollisionRegistryObjectivesPositiveEffects(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	s := collisionRegistryObjectiveBase()
	p := s.Callbacks.Shop
	p.Resources.Mana = 49
	p.Resources.MaxMana = 50
	p.TemporaryUpdates.UpdateWords[0][0] = 8
	p.Sequence = []legacy.PortTestShopAction{{Op: 808}}
	cases = append(cases, s)
	s = collisionRegistryObjectiveBase()
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
	s = collisionRegistryObjectiveBase()
	p = s.Callbacks.Shop
	p.Items[0].Class = 0x10000000
	p.TemporaryUpdates.World.Objectives.ColorNames = []string{"BlueFlag"}
	p.Sequence = []legacy.PortTestShopAction{{Op: 806}}
	cases = append(cases, s)
	s = collisionRegistryObjectiveBase()
	p = s.Callbacks.Shop
	p.Inventory.TeamMembers = true
	p.Inventory.Teams = [4]byte{1, 2, 2, 0}
	p.TemporaryUpdates.World.ItemNames = []string{"gameball"}
	p.Sequence = []legacy.PortTestShopAction{{Op: 803}}
	cases = append(cases, s)
	s = collisionRegistryObjectiveBase()
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
	collisionRegistryHash(t, "collision-registry-objectives-positive", r, collisionRegistryOwnerHashes["collision-registry-objectives-positive"])
}

func collisionRegistryObjectiveBase() legacy.PortTestRoamSpec {
	s := objectiveBase()
	s.Callbacks.Shop.TemporaryUpdates.World.Objectives.RegisteredCollision = true
	return s
}
