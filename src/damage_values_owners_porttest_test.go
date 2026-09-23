//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

// Frozen from the original integer-result projectile damage owner.
var damageValueHashes = map[string]string{
	"damage-value-basic-02":       "eeacfa3de8fe2f201d53411aa1a55edd5d63074cb5437ab4fb9e102d6580ec5e",
	"damage-value-basic-08":       "26efe969ae70d0577d05fa016f628552b29e35b39076d4bbf2cc206b43af7ebc",
	"damage-value-basic-09":       "574121289fb695aed20ffd6abf658fce9ba0fcb81fcc4dc39f72e9b3a9c1e372",
	"damage-value-basic-10":       "209c91e4d1646ff8ec4d5bd8a60147bbedb298cf33712f6dd554b46a02ff8320",
	"damage-value-basic-14":       "f3e47315b0cbb8c40f4cc1c1f8198726d334faca869c497548a55984105f8ace",
	"damage-value-basic-20":       "3ee123dd759f97d9999539f72b4318d882a8438d159237b63eb7e5f1fa3710b2",
	"damage-value-basic-21":       "629ed7080385c9a74309a9ec3f80d4c2ff398512c07a4f4ff9b594db1529a425",
	"damage-value-basic-22":       "d40cacd0ea3cec927058b33a46bea9593edbd48abfd3067d1c541cfef0aae11f",
	"damage-value-basic-23":       "b7d2cf5279331b8c6d0edab33a14fa76bcd1bfce5b994eeec5e9ffcc496391ab",
	"damage-value-basic-24":       "39133581ded8ca8053ff79e63835f01603221b314203c226eb60722577ecc699",
	"damage-value-generator":      "2f7d2020ce2fda68530c5c7db3f4f7b50590d5edac15d86f4a430944a9254d55",
	"damage-value-golem-wrapping": "dbec9d6105f0cf6b9a2f821e06e853f011734c3fba92d90e05fc95af3addc76d",
	"damage-value-positive":       "43e3a2d989735df453c56b103068c4c5d7000a8aa9096ae5a9ab6ce253703a56",
	"damage-value-source-02":      "b29253349968badb08cae4b541129799bc1dfc52d866a1318f78f66a47eb20a3",
	"damage-value-source-14":      "9b0bdb7b170d4fa26f8312e09129f3571bfbf281553d47831ba2b096f00f7e6a",
}

func damageValueBase(op int) legacy.PortTestRoamSpec {
	s := damageBase(op)
	name, ok := damageRegistryNames[op]
	if !ok {
		panic("unmapped damage callback")
	}
	s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Damage.Registry = name
	s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Damage.RegistryValue = true
	return s
}
func damageValueHash(t *testing.T, name string, specs []legacy.PortTestRoamSpec) {
	t.Helper()
	results := effectsTimedRun(t, specs)
	callbackHash(t, name, results, damageValueHashes[name])
}

func TestDamageValueBasic(t *testing.T) {
	for _, op := range []int{2, 8, 9, 10, 14, 20, 21, 22, 23, 24} {
		t.Run(fmt.Sprint(op), func(t *testing.T) {
			var cases []legacy.PortTestRoamSpec
			for _, subject := range []int{1, 2, 3} {
				for _, amount := range []int32{0, 1, 12, 49, 50, 100} {
					for _, kind := range []int32{0, 1, 2, 5, 7, 9, 11, 12, 15, 17} {
						s := damageValueBase(op)
						p := s.Callbacks.Shop
						p.Resources.Subject = subject
						p.Resources.Subclass = 0x10
						d := p.TemporaryUpdates.World.Objectives.Attack.Damage
						d.Amount = amount
						d.Kind = kind
						cases = append(cases, s)
					}
				}
			}
			damageValueHash(t, fmt.Sprintf("damage-value-basic-%02d", op), cases)
		})
	}
}

func TestDamageValueGenerator(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, hp := range []uint16{1, 33, 34, 66, 67, 100} {
		for _, amount := range []int32{0, 1, 20, 100} {
			for _, frame := range []uint32{0, 20, 21, 29, 30, 0xffffffff} {
				s := damageValueBase(25)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				p.Resources.Subject = 2
				p.Resources.HP = hp
				a.Damage.Amount = amount
				s.Owner.Frame = frame
				a.UpdateWords = map[int]uint32{48: 0xffffffff}
				a.ActorWords = map[int]uint32{536: 0}
				cases = append(cases, s)
			}
		}
	}
	damageValueHash(t, "damage-value-generator", cases)
}

func TestDamageValueSource(t *testing.T) {
	for _, op := range []int{2, 14} {
		var cases []legacy.PortTestRoamSpec
		for _, source := range []int{0, 1, 101} {
			for _, weapon := range []int{0, 3, 4} {
				for _, mode := range []uint32{1, 2049, 4097} {
					for _, gameplay := range []uint32{0, 1} {
						s := damageValueBase(op)
						p := s.Callbacks.Shop
						a := p.TemporaryUpdates.World.Objectives.Attack
						s.Lifecycle.GameFlags = mode
						p.Inventory.Gameplay = gameplay
						a.Damage.Source = source
						a.Damage.Weapon = weapon
						cases = append(cases, s)
					}
				}
			}
		}
		damageValueHash(t, fmt.Sprintf("damage-value-source-%02d", op), cases)
	}
}

func TestDamageValuePositive(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{2, 14, 21, 22} {
		s := damageValueBase(op)
		cases = append(cases, s)
	}
	r := effectsTimedRun(t, cases)
	for i, v := range r {
		if v.Callbacks.Shop.Sequence[0].Return != 1 {
			t.Fatalf("case %d: positive damage rejected", i)
		}
		if hp := v.Callbacks.Shop.Sequence[0].ResourceData[2][0] & 65535; hp != 38 {
			t.Fatalf("case %d HP=%d want 38", i, hp)
		}
	}
	callbackHash(t, "damage-value-positive", r, damageValueHashes["damage-value-positive"])
}

func TestDamageValueGolemWrapping(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, amount := range []int32{-1 << 31, -1, 0, 1, 1<<31 - 1} {
		for _, kind := range []int32{0, 9, 17} {
			s := damageValueBase(22)
			p := s.Callbacks.Shop
			p.Resources.Subject = 2
			p.TemporaryUpdates.World.Objectives.Attack.Damage.Amount = amount
			p.TemporaryUpdates.World.Objectives.Attack.Damage.Kind = kind
			cases = append(cases, s)
		}
	}
	damageValueHash(t, "damage-value-golem-wrapping", cases)
}
