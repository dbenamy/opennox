//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"os"
	"testing"
)

var questEligibilityHashes = map[string]string{
	"limit-rounding":      "0f7be38f1663dd68ff18d14fef4ffc96a4b9be0885dcbd202a6d5086812d511e",
	"scalars":             "5b822ec6366f6688b35e79c8a490357810407c5ad821ec4bd186fde7e8b2b091",
	"books":               "fb7d941e428909476e12906ab612800f8072fd6b615abb33f12189dadc264170",
	"modifier-masks":      "15f40e27ff44ebe33e1f7ccecf3271a157a60d0eba9aa8e951e6b62ad645c52f",
	"effect-slots":        "e7bcb41edfc8974a99c16f4140975230bf0f9d1e177eada6b383486f7f20d1f2",
	"composite-modifiers": "a961365d47f828ba9813d38cce5113836806f31777145c94e3dad8092c811cea",
	"item-classes":        "04c5f6c4724ec116fa9258c814b8c7fde9c185b77dba739adc626ec76ba9ab7d",
	"special-equipment":   "ebc5247aed7aab4276ba12ca1c3bdc0a9ae671fd19cf3b6283b3cd0238e6138e",
	"inventory-limits":    "e0ce7e6bd8b6e0f7c1a7f93aa8afeea702015708af306119d174fdebd5b920fd",
	"destroyed-inventory": "fb026c1ab0056899cc2d2c125781aaf461b99618623699e40f3d0bfabd164d93",
}

func eligibilityBase(op int, args ...legacy.PortTestGameplayReportArg) legacy.PortTestRoamSpec {
	s := gameplayReportsBase(300+op, args...)
	textSpec(&s).Eligibility = &legacy.PortTestQuestEligibilitySpec{
		Objects:      []legacy.PortTestEligibilityObject{{Type: "EligibilityWeapon", Class: 0x1000000}},
		Modifiers:    []legacy.PortTestEligibilityModifier{{Name: "EligibilityMod", WeaponMask: 0xffffffff, ArmorMask: 0xffffffff, Flags: 3}},
		SpecialFlags: [4]uint32{3, 3, 3, 3}, WeaponBit: 4, ArmorBit: 4,
	}
	return s
}
func eligibilitySpec(s *legacy.PortTestRoamSpec) *legacy.PortTestQuestEligibilitySpec {
	return textSpec(s).Eligibility
}
func eligibilityObject(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "eligibility-object", Ref: ref}
}
func eligibilityCapture(t *testing.T, label string, out []legacy.PortTestRoamResult) {
	t.Helper()
	for i, r := range out {
		if !r.Intact || !r.Callbacks.Intact || !r.Spells.Intact || !r.Combat.Intact || !r.MonsterState.Intact || !r.Callbacks.Shop.Intact {
			t.Fatalf("eligibility %s case%d guards", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_QUEST_ELIGIBILITY_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d cases %s", label, len(out), hash)
	want, ok := questEligibilityHashes[label]
	if !ok {
		t.Fatalf("missing audited C capture: %s", label)
	}
	if hash != want {
		t.Fatalf("%s hash %s want%s", label, hash, want)
	}
}
func eligibilityCheck(t *testing.T, label string, cases []legacy.PortTestRoamSpec, wants []uint32) {
	t.Helper()
	out := controlsRun(t, cases)
	for i, r := range out {
		if got := r.Callbacks.Shop.Sequence[0].Return; got != wants[i] {
			t.Fatalf("%s case%d return%d want%d", label, i, got, wants[i])
		}
	}
	eligibilityCapture(t, label, out)
}
func eligibilityBool(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}
func TestQuestEligibilityScalars(t *testing.T) {
	tables := [][][2]uint32{
		nil,
		{{1, 1}, {5, 2}, {9, 1}, {27, 1}, {34, 1}, {41, 1}, {75, 1}, {140, 1}, {0xffffffff, 1}},
		{{5, 0}, {5, 7}, {41, 0}, {41, 0}, {140, 1}},
		{{1, 1}, {0, 1}, {5, 1}, {140, 1}},
	}
	values := []uint32{0x7fffffff, 0x80000000, 0xffffffff, 255, 256}
	for i := uint32(0); i <= 141; i++ {
		values = append(values, i)
	}
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, op := range []int{0, 1, 2, 12, 13} {
		for _, table := range tables {
			for _, value := range values {
				s := eligibilityBase(op, reportValue(value))
				sp := eligibilitySpec(&s)
				sp.SpellTable = table
				sp.BeastTable = table
				sp.BeastGroups = [][]uint32{{99, 4, 7, 0, 8}, {123, 40}}
				allowed := false
				for _, row := range table {
					if row[0] == 0 {
						break
					}
					if row[0] == value && row[1] != 0 {
						allowed = true
						break
					}
				}
				switch op {
				case 0:
					allowed = allowed && value != 0 && value != 34 && value != 27 && value != 9 && value != 41
				case 1:
					allowed = allowed && value != 0
				case 2:
					allowed = int32(value) > 0 && int32(value) < 6
				case 12:
					allowed = allowed || (value >= 46 && value <= 49) || (value >= 122 && value <= 125) || (value >= 75 && value <= 114)
				case 13:
					// Each group's leading marker is skipped. The terminating zero itself can
					// match zero, unlike the ordinary table's sentinel.
					for _, group := range sp.BeastGroups {
						if len(group) == 0 || group[0] == 0 {
							continue
						}
						ids := append(append([]uint32(nil), group[1:]...), 0)
						for _, id := range ids {
							if id == value {
								allowed = true
							}
							if id == 0 {
								break
							}
						}
					}
				}
				cases = append(cases, s)
				wants = append(wants, eligibilityBool(allowed))
			}
		}
	}
	eligibilityCheck(t, "scalars", cases, wants)
}
func TestQuestEligibilityBooks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, op := range []int{3, 4} {
		for _, sub := range []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 0x80000000} {
			for _, id := range []byte{0, 1, 4, 5, 6, 40, 41, 255} {
				s := eligibilityBase(op, eligibilityObject(1))
				sp := eligibilitySpec(&s)
				u := &sp.Objects[0]
				u.Class = 0x100
				u.Subclass = sub
				u.BookID = id
				u.Book = fmt.Sprintf("EligibilityGuide%02d", id)
				sp.SpellTable = [][2]uint32{{1, 1}, {4, 0}, {5, 7}, {255, 1}}
				sp.BeastTable = [][2]uint32{{0, 1}, {1, 1}, {40, 1}}
				// A zero first row terminates the beast table; spell rows are independent.
				want := false
				if sub&1 != 0 {
					want = id == 1 || id == 5 || id == 255
				} else if sub&2 != 0 {
					want = false
				} else if sub&4 != 0 {
					want = id > 0 && id < 6
				}
				cases = append(cases, s)
				wants = append(wants, eligibilityBool(want))
			}
		}
	}
	for _, name := range []string{"EligibilityGuide00", "EligibilityGuide01", "eligibilityguide01", "EligibilityGuide40", "EligibilityGuide41", "", "missing"} {
		s := eligibilityBase(4, eligibilityObject(1))
		sp := eligibilitySpec(&s)
		sp.Objects[0].Class = 0x100
		sp.Objects[0].Subclass = 2
		sp.Objects[0].Book = name
		sp.BeastTable = [][2]uint32{{1, 1}, {40, 1}}
		cases = append(cases, s)
		wants = append(wants, eligibilityBool(name == "EligibilityGuide01" || name == "EligibilityGuide40"))
	}
	eligibilityCheck(t, "books", cases, wants)
}
func TestQuestEligibilityModifierMasks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, op := range []int{6, 7, 8} {
		for _, weapon := range []bool{false, true} {
			for _, bit := range []uint32{0, 1, 4, 0x80000000, 0xffffffff} {
				for _, mask := range []uint32{0, 4, 0xffffffff} {
					for _, exclude := range []uint32{0, 4, 0x80000000, 0xffffffff} {
						for mode := 0; mode < 7; mode++ {
							s := eligibilityBase(op, eligibilityObject(1))
							sp := eligibilitySpec(&s)
							sp.ArmorBit = bit
							sp.WeaponBit = bit
							sp.Modifiers[0].ArmorMask = mask
							sp.Modifiers[0].WeaponMask = mask
							u := &sp.Objects[0]
							u.Type = "EligibilityArmor"
							u.Class = 0x2000000
							if weapon {
								u.Type = "EligibilityWeapon"
								u.Class = 0x1000000
							}
							slot := 0
							if op == 7 {
								slot = 1
							}
							if op == 8 {
								slot = 2
							}
							u.Mods[slot] = 5
							row := legacy.PortTestEligibilityModRow{Ref: 5, Name: lookupString("row"), ArmorExclude: exclude, WeaponExclude: exclude}
							rows := []legacy.PortTestEligibilityModRow{row}
							switch mode {
							case 0:
								rows = nil
							case 2:
								rows[0].Name = nil
							case 3:
								rows = append([]legacy.PortTestEligibilityModRow{{Ref: 4, Name: lookupString("other")}}, row)
							case 4:
								rows = append(rows, legacy.PortTestEligibilityModRow{Ref: 5, Name: lookupString("duplicate")})
							case 5:
								rows = append([]legacy.PortTestEligibilityModRow{{Ref: 4}}, row)
							case 6:
								row.Name = nil
								rows = append([]legacy.PortTestEligibilityModRow{{Ref: 4, Name: lookupString("other")}}, row)
							}
							switch op {
							case 6:
								sp.MaterialWeapon = rows
								sp.MaterialArmor = rows
							case 7:
								sp.Quality = rows
							case 8:
								sp.Effects = rows
							}
							cases = append(cases, s)
							wants = append(wants, eligibilityBool((mode == 1 || mode == 3 || mode == 4) && bit&mask != 0 && bit&exclude == 0))
						}
					}
				}
			}
		}
	}
	eligibilityCheck(t, "modifier-masks", cases, wants)
}
func TestQuestEligibilityEffectSlots(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, slot := range []int{2, 3} {
		for _, flags := range []uint32{0, 1, 2, 3, 4, 0x100, 0xffffffff} {
			for _, ref := range []int{1, 2, 3, 4, 5} {
				for _, marker := range []string{"", "#", "#item", "item"} {
					for _, matching := range []bool{false, true} {
						for _, warm := range []bool{false, true} {
							s := eligibilityBase(8, eligibilityObject(1))
							sp := eligibilitySpec(&s)
							sp.Warm = warm
							sp.Objects[0].Mods[slot] = ref
							if ref <= 4 {
								sp.SpecialFlags[ref-1] = flags
							} else {
								sp.Modifiers[0].Flags = flags
							}
							typ := "EligibilityArmor"
							if matching {
								typ = "EligibilityWeapon"
							}
							sp.ItemTable = []legacy.PortTestEligibilityItemRow{{Name: lookupString(marker), Type: typ}}
							sp.Effects = []legacy.PortTestEligibilityModRow{{Ref: 5, Name: lookupString("generic")}}
							allowed := flags&uint32(1<<(slot-2)) != 0
							if ref <= 4 {
								allowed = allowed && matching && len(marker) > 0 && marker[0] == '#'
							}
							cases = append(cases, s)
							wants = append(wants, eligibilityBool(allowed))
						}
					}
				}
			}
		}
	}
	eligibilityCheck(t, "effect-slots", cases, wants)
}
func TestQuestEligibilityCompositeModifiers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, op := range []int{3, 5, 9} {
		for _, weapon := range []bool{false, true} {
			for present := 0; present < 16; present++ {
				for missing := -1; missing < 4; missing++ {
					s := eligibilityBase(op, eligibilityObject(1))
					sp := eligibilitySpec(&s)
					u := &sp.Objects[0]
					u.Type = "EligibilityArmor"
					u.Class = 0x2000000
					if weapon {
						u.Type = "EligibilityWeapon"
						u.Class = 0x1000000
					}
					sp.ItemTable = []legacy.PortTestEligibilityItemRow{{Name: lookupString("item"), Type: u.Type}}
					row := []legacy.PortTestEligibilityModRow{{Ref: 5, Name: lookupString("row")}}
					sp.MaterialWeapon = row
					sp.MaterialArmor = row
					sp.Quality = row
					sp.Effects = row
					for slot := 0; slot < 4; slot++ {
						if present&(1<<slot) != 0 {
							u.Mods[slot] = 5
						}
					}
					if missing == 0 {
						sp.MaterialWeapon = nil
						sp.MaterialArmor = nil
					}
					if missing == 1 {
						sp.Quality = nil
					}
					if missing >= 2 {
						sp.Effects = nil
					}
					allowed := true
					if missing >= 0 {
						bits := 1 << missing
						if missing >= 2 {
							bits = 12
						}
						allowed = present&bits == 0
					}
					cases = append(cases, s)
					wants = append(wants, eligibilityBool(allowed))
				}
			}
		}
	}
	eligibilityCheck(t, "composite-modifiers", cases, wants)
}
func TestQuestEligibilityItemClasses(t *testing.T) {
	names := []string{"Diamond", "Emerald", "Ruby", "SulphorousFlareWand", "StreetSneakers", "StreetShirt", "StreetPants", "RedPotion", "BluePotion", "EligibilityWeapon", "EligibilityArmor", "EligibilityBook"}
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, warm := range []bool{false, true} {
		for _, name := range names {
			for _, class := range []uint32{0, 0x40, 0x10, 0x100, 0x1000000, 0x2000000, 0x50, 0x1000010} {
				for _, sub := range []uint32{0, 1, 4, 0x10, 0x20000} {
					for _, listed := range []bool{false, true} {
						s := eligibilityBase(3, eligibilityObject(1))
						sp := eligibilitySpec(&s)
						sp.Warm = warm
						u := &sp.Objects[0]
						u.Type = name
						u.Class = class
						u.Subclass = sub
						u.BookID = 1
						sp.SpellTable = [][2]uint32{{1, 1}}
						if listed {
							sp.ItemTable = []legacy.PortTestEligibilityItemRow{{Name: lookupString("item"), Type: name}}
						}
						allowed := listed
						switch {
						case class&0x40 != 0:
							allowed = false
						case class&0x10 != 0:
							allowed = sub&0x1ff78 != 0
						case class&0x100 != 0:
							allowed = sub&1 != 0 || sub&4 != 0
						case name == "Diamond" || name == "Emerald" || name == "Ruby":
							allowed = true
						case name == "SulphorousFlareWand" || name == "StreetSneakers" || name == "StreetShirt" || name == "StreetPants":
							allowed = true
						}
						cases = append(cases, s)
						wants = append(wants, eligibilityBool(allowed))
					}
				}
			}
		}
	}
	for _, name := range []string{"Diamond", "Emerald", "Ruby"} {
		for _, class := range []uint32{0x1000000, 0x2000000} {
			s := eligibilityBase(3, eligibilityObject(1))
			sp := eligibilitySpec(&s)
			sp.Objects[0].Type = name
			sp.Objects[0].Class = class
			sp.Objects[0].Mods = [4]int{5, 5, 5, 5}
			cases = append(cases, s)
			wants = append(wants, 1)
		}
	}
	eligibilityCheck(t, "item-classes", cases, wants)
}
func TestQuestEligibilitySpecialEquipment(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, op := range []int{3, 10} {
		for _, warm := range []bool{false, true} {
			for _, bit := range []uint32{0, 4, 0x10000, 0x10004} {
				for _, mods := range [][4]int{{}, {0, 0, 1, 0}, {5, 0, 1, 0}, {0, 5, 1, 0}, {0, 0, 2, 0}, {0, 0, 1, 5}} {
					s := eligibilityBase(op, eligibilityObject(1))
					sp := eligibilitySpec(&s)
					sp.Warm = warm
					sp.WeaponBit = bit
					sp.Objects[0].Type = "SulphorousFlareWand"
					sp.Objects[0].Mods = mods
					allowed := bit&0x10000 == 0 || mods == [4]int{0, 0, 1, 0}
					cases = append(cases, s)
					wants = append(wants, eligibilityBool(allowed))
				}
			}
		}
	}
	for _, op := range []int{3, 10} {
		for _, name := range []string{"StreetSneakers", "StreetPants", "StreetShirt"} {
			for _, bit := range []uint32{0, 1, 4, 0x400, 0x405, 0x10000} {
				for _, modName := range []string{"UserColor1", "usercolor1", "USERCOLOR1", "UserColo", "UserCol", "UserColX", "Other", ""} {
					for _, slot := range []int{0, 1, 2, 3} {
						s := eligibilityBase(op, eligibilityObject(1))
						sp := eligibilitySpec(&s)
						u := &sp.Objects[0]
						u.Type = name
						u.Class = 0x2000000
						u.Mods[slot] = 5
						sp.ArmorBit = bit
						sp.Modifiers[0].Name = modName
						valid := modName == "UserColor1" || modName == "usercolor1" || modName == "USERCOLOR1" || modName == "UserColo"
						cases = append(cases, s)
						wants = append(wants, eligibilityBool(bit&0x405 == 0 || valid))
					}
				}
			}
		}
	}
	eligibilityCheck(t, "special-equipment", cases, wants)
}
func TestQuestEligibilityInventoryLimits(t *testing.T) {
	potions := []string{"RedPotion", "BluePotion", "CurePoisonPotion", "HastePotion", "InvisibilityPotion", "ShieldPotion", "VampirismPotion", "FireProtectPotion", "ShockProtectPotion", "PoisonProtectPotion", "InvulnerabilityPotion", "InfravisionPotion"}
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, warm := range []bool{false, true} {
		for _, name := range potions {
			for _, count := range []int{0, 1, 9, 10, 11} {
				for _, nested := range []bool{false, true} {
					s := eligibilityBase(11, eligibilityObject(1))
					sp := eligibilitySpec(&s)
					sp.Warm = warm
					sp.Objects[0] = legacy.PortTestEligibilityObject{Type: "EligibilityPlayer", Class: 4}
					s.Callbacks.Shop.EffectsUse.Balance["ForceOfNatureStaffLimit"] = []float64{2}
					for i := 0; i < count; i++ {
						sp.Objects = append(sp.Objects, legacy.PortTestEligibilityObject{Type: name})
						sp.Objects[0].Inventory = append(sp.Objects[0].Inventory, i+2)
					}
					if nested {
						sp.Objects = append(sp.Objects, legacy.PortTestEligibilityObject{Type: "EligibilityBook", Inventory: append([]int(nil), sp.Objects[0].Inventory...)})
						sp.Objects[0].Inventory = []int{len(sp.Objects)}
					}
					cases = append(cases, s)
					wants = append(wants, eligibilityBool(nested || count <= 9))
				}
			}
		}
	}
	limits := []struct {
		bits    uint32
		integer int32
	}{{0, 0}, {0x3f800000, 1}, {0x3fffffff, 1}, {0x40000000, 2}, {0x40200000, 2}, {0xbf800000, -1}, {0x80000000, 0}, {0x7f800000, -2147483648}, {0xff800000, -2147483648}, {0x7fc00001, -2147483648}, {0x4f000000, -2147483648}, {0x4effffff, 2147483520}}
	for _, warm := range []bool{false, true} {
		for _, limit := range limits {
			for _, count := range []int{0, 1, 2, 3} {
				s := eligibilityBase(11, eligibilityObject(1))
				sp := eligibilitySpec(&s)
				sp.Warm = warm
				sp.Objects[0] = legacy.PortTestEligibilityObject{Type: "EligibilityPlayer", Class: 4}
				s.Callbacks.Shop.EffectsUse.Balance["ForceOfNatureStaffLimit"] = []float64{float64(math.Float32frombits(limit.bits))}
				for i := 0; i < count; i++ {
					sp.Objects = append(sp.Objects, legacy.PortTestEligibilityObject{Type: "InfinitePainWand"})
					sp.Objects[0].Inventory = append(sp.Objects[0].Inventory, i+2)
				}
				cases = append(cases, s)
				wants = append(wants, eligibilityBool(int32(count) <= limit.integer))
			}
		}
	}
	for _, ref := range []int{0, 1} {
		for _, warm := range []bool{false, true} {
			s := eligibilityBase(11, eligibilityObject(ref))
			sp := eligibilitySpec(&s)
			sp.Warm = warm
			s.Callbacks.Shop.EffectsUse.Balance["ForceOfNatureStaffLimit"] = []float64{-1}
			cases = append(cases, s)
			wants = append(wants, 1)
		}
	}
	eligibilityCheck(t, "inventory-limits", cases, wants)
}
func TestQuestEligibilityDestroyedInventory(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, name := range []string{"RedPotion", "BluePotion", "CurePoisonPotion", "HastePotion", "InvisibilityPotion", "ShieldPotion", "VampirismPotion", "FireProtectPotion", "ShockProtectPotion", "PoisonProtectPotion", "InvulnerabilityPotion", "InfravisionPotion", "InfinitePainWand"} {
		for _, flags := range []uint32{0, 0x10, 0x20, 0xffffffff} {
			s := eligibilityBase(11, eligibilityObject(1))
			sp := eligibilitySpec(&s)
			sp.Objects[0] = legacy.PortTestEligibilityObject{Type: "EligibilityPlayer", Class: 4}
			s.Callbacks.Shop.EffectsUse.Balance["ForceOfNatureStaffLimit"] = []float64{9}
			for i := 0; i < 10; i++ {
				o := legacy.PortTestEligibilityObject{Type: name}
				if i == 9 {
					o.Flags = flags
				}
				sp.Objects = append(sp.Objects, o)
				sp.Objects[0].Inventory = append(sp.Objects[0].Inventory, i+2)
			}
			cases = append(cases, s)
			wants = append(wants, eligibilityBool(flags&0x20 != 0))
		}
	}
	eligibilityCheck(t, "destroyed-inventory", cases, wants)
}
func TestQuestEligibilityLimitRounding(t *testing.T) {
	// C stores the balance's double in a float before converting it to an int.
	// These values straddle float32 rounding boundaries, not just integer edges.
	limits := []struct {
		value   float64
		integer int32
	}{{1.99999997, 2}, {1.99999990, 1}, {-0.99999999, -1}, {-0.99999995, 0}}
	var cases []legacy.PortTestRoamSpec
	var wants []uint32
	for _, limit := range limits {
		for count := 0; count <= 3; count++ {
			s := eligibilityBase(11, eligibilityObject(1))
			sp := eligibilitySpec(&s)
			sp.Objects[0] = legacy.PortTestEligibilityObject{Type: "EligibilityPlayer", Class: 4}
			s.Callbacks.Shop.EffectsUse.Balance["ForceOfNatureStaffLimit"] = []float64{limit.value}
			for i := 0; i < count; i++ {
				sp.Objects = append(sp.Objects, legacy.PortTestEligibilityObject{Type: "InfinitePainWand"})
				sp.Objects[0].Inventory = append(sp.Objects[0].Inventory, i+2)
			}
			cases = append(cases, s)
			wants = append(wants, eligibilityBool(int32(count) <= limit.integer))
		}
	}
	eligibilityCheck(t, "limit-rounding", cases, wants)
}
