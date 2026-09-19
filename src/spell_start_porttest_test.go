//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestSpellStartTeleport(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, caster := range []int{1, 100} {
		for _, recipient := range []int{0, 1, 100, 101} {
			for _, frame := range []uint32{0, 100, 0xffffffff} {
				for _, coop := range []bool{false, true} {
					for tile := 0; tile <= 5; tile++ {
						for _, blocked := range []bool{false, true} {
							for _, level := range []uint32{1, 3, 5} {
								s := sustainedBase(53)
								s.Owner.Frame = frame
								if coop {
									s.Lifecycle.GameFlags |= 2048
								} else {
									s.Lifecycle.GameFlags &^= 2048
								}
								p := s.Callbacks.Shop
								if blocked {
									p.Inventory.WallMode = 1
								} else {
									p.Inventory.WallMode = 0
								}
								sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
								sp.Effects.Tile = new(int)
								*sp.Effects.Tile = tile
								sp.Effects.Sustained.Start = &legacy.PortTestSpellStartSpec{Blocked: blocked}
								sp.Record.Refs[16] = caster
								sp.Record.Refs[48] = recipient
								sp.Record.Words[8] = level
								sp.Record.Words[52] = math.Float32bits(200)
								sp.Record.Words[56] = math.Float32bits(100)
								sp.Record.Words[68] = 0x12345678
								p.EffectsUse.Balance["TeleportDelay"] = []float64{-2.75, -0.5, 0.75, 1.99999999, 30.875}
								cases = append(cases, s)
							}
						}
					}
				}
			}
		}
	}
	sustainedHash(t, "spell-start-teleport", cases)
}

func TestSpellStartTeleportRounding(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, value := range []float64{-30.875, -1.99999999, -0.5, 0, 0.99999999, 1.99999999, 30.875, 16777217} {
		for _, frame := range []uint32{0, 100, 0xffffffff} {
			s := sustainedBase(53)
			s.Owner.Frame = frame
			s.Lifecycle.GameFlags |= 2048
			p := s.Callbacks.Shop
			sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
			sp.Effects.Tile = new(int)
			sp.Effects.Sustained.Start = &legacy.PortTestSpellStartSpec{}
			sp.Record.Refs[16] = 100
			sp.Record.Refs[48] = 0
			sp.Record.Words[52] = math.Float32bits(200)
			sp.Record.Words[56] = math.Float32bits(100)
			p.EffectsUse.Balance["TeleportDelay"] = []float64{value, value, value, value, value}
			cases = append(cases, s)
		}
	}
	sustainedHash(t, "spell-start-teleport-rounding", cases)
}

func TestSpellStartPixies(t *testing.T) {
	spellStartDirections(t)
	var cases []legacy.PortTestRoamSpec
	for _, seed := range []int{0, 777} {
		for _, level := range []uint32{1, 2, 3, 5} {
			for _, existing := range []int{0, 1, 5} {
				for _, radius := range []float32{-4, 10.25, 110} {
					for _, wall := range []int{0, 1} {
						for _, warm := range []bool{false, true} {
							for _, missing := range []bool{false, true} {
								if missing && existing != 0 {
									continue
								}
								for _, fps := range []uint32{1, 30} {
									s := sustainedBase(54)
									s.Seed = seed
									s.Owner.FPS = fps
									if seed == 777 {
										s.Owner.Frame = 0xffffffff
									}
									p := s.Callbacks.Shop
									p.Inventory.WallMode = wall
									a := p.TemporaryUpdates.World.Objectives.Attack
									if missing {
										a.MissingTypes = []string{"Pixie"}
									}
									sp := a.Controls.SpellLifecycle
									sp.Record.Words[4] = 58
									sp.Record.Words[8] = level
									sp.Effects.Sustained.Start = &legacy.PortTestSpellStartSpec{Existing: existing, WarmCache: warm, Radius: radius}
									p.EffectsUse.Balance["PixieCount"] = []float64{0, -1.5, 2.9, 3.99999999, 5}
									cases = append(cases, s)
								}
							}
						}
					}
				}
			}
		}
	}
	sustainedHash(t, "spell-start-pixies", cases)
}

func spellStartDirections(t *testing.T) {
	t.Helper()
	table := memmap.Slice(0x587000, 194136)[:2048]
	old := bytes.Clone(table)
	copy(table, blobdata.PortTestCombatTables()[194136])
	t.Cleanup(func() { copy(table, old) })
}
func TestSpellStartPixieOwners(t *testing.T) {
	spellStartDirections(t)
	var cases []legacy.PortTestRoamSpec
	for _, nilOwner := range []bool{false, true} {
		for _, warm := range []bool{false, true} {
			for _, missing := range []bool{false, true} {
				for _, wall := range []int{0, 1} {
					for seed := 0; seed < 4; seed++ {
						s := sustainedBase(54)
						s.Seed = seed
						s.Owner.Frame = 0xffffffff
						p := s.Callbacks.Shop
						p.Inventory.WallMode = wall
						a := p.TemporaryUpdates.World.Objectives.Attack
						if missing {
							a.MissingTypes = []string{"Pixie"}
						}
						sp := a.Controls.SpellLifecycle
						sp.Record.Words[4] = 58
						sp.Record.Words[8] = 3
						sp.Effects.Sustained.Start = &legacy.PortTestSpellStartSpec{NilOwner: nilOwner, PlayerOwner: !nilOwner, WarmCache: warm, Radius: 110.12345}
						p.EffectsUse.Balance["PixieCount"] = []float64{1, 2, 3, 4, 5}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	sustainedHash(t, "spell-start-pixie-owners", cases)
}

func TestSpellStartPixieLimits(t *testing.T) {
	spellStartDirections(t)
	var cases []legacy.PortTestRoamSpec
	for _, limit := range []float64{0.99999999, 1.99999999, 3.99999999, -1.99999999} {
		for _, count := range []int{0, 1, 5} {
			for seed := 0; seed < 2; seed++ {
				s := sustainedBase(54)
				s.Seed = seed
				p := s.Callbacks.Shop
				sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
				sp.Record.Words[4] = 58
				sp.Effects.Sustained.Start = &legacy.PortTestSpellStartSpec{Existing: count, Radius: 10.125}
				p.EffectsUse.Balance["PixieCount"] = []float64{limit, limit, limit, limit, limit}
				cases = append(cases, s)
			}
		}
	}
	sustainedHash(t, "spell-start-pixie-limits", cases)
}
func TestSpellStartCharmControl(t *testing.T) {
	got := legacy.PortTestSpellCharmControl([]bool{false, true, true, false, false, true, false, true})
	want := []int32{0, 1, 1, 0, 0, 1, 0, 1}
	if len(got) != len(want) {
		t.Fatalf("charm setter returned %d values, want %d", len(got), len(want))
	}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("charm setter %d: got %d want %d", i, v, want[i])
		}
	}
}
