//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func creationBase(op int) legacy.PortTestRoamSpec {
	s := callbackBase(33 + op)
	s.Callbacks.Creation = &legacy.PortTestCreationSpec{TypeIndex: 15, Class: 0x1000008, Flags: 4, WeaponType: 15, ArmorType: 15, Durability: 125, Effects: true, InitFill: 0x6a, UseFill: 0x35, AutoFlag: 0x7f, Balance: map[string]float64{"QuestDurabilityMultiplier": 1.5, "DefaultAmmoAmount": 20, "DefaultAmmoAmountQuest": 30, "QuestStaffChargeMultiplier": 1.5, "ToxicCloudLifetime": 2.5}}
	if op == 0 {
		s.Callbacks.Creation.Class = 2
	}
	return s
}
func TestObjectCreationSmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op < 11; op++ {
		specs = append(specs, creationBase(op))
	}
	callbackHash(t, "creation-smoke", legacy.PortTestRoam(specs), creationHashes["creation-smoke"])
}
func TestObjectCreationAutoSpells(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	names := []string{"UrchinShaman", "Wizard", "WizardWhite", "Beholder", "Lich", "LichLord", "Demon", "WizardGreen", "WillOWisp", ""}
	for _, name := range names {
		for _, fps := range []uint32{0, 1, 2, 3, 30, 65535, 65536, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, flags := range []uint32{0, 4096} {
				s := creationBase(0)
				s.Callbacks.Creation.Type = name
				s.Owner.FPS = fps
				s.Lifecycle.GameFlags = flags
				for id := uint32(1); id <= 136; id++ {
					s.Spells.Permissions = append(s.Spells.Permissions, [2]uint32{id, 0x12340000 + id})
				}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "creation-auto", legacy.PortTestRoam(specs), creationHashes["creation-auto"])
}
func TestObjectCreationEquipment(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	values := []float64{0, 1, 1.5, -1, 1.9999999, 65536, 1e20, math.Inf(1), math.Inf(-1), math.NaN()}
	for _, op := range []int{1, 2} {
		for _, name := range []string{"", "OblivionHeart", "OblivionWierdling", "OblivionOrb"} {
			for _, sub := range []uint32{0, 2, 4, 8, 0x80, 0x40000, 0x10000, 0x47f0000} {
				for j, v := range values {
					s := creationBase(op)
					c := s.Callbacks.Creation
					c.Type = name
					c.Subclass = sub
					if sub >= 0x10000 {
						c.Class = 0x1008
					}
					s.Lifecycle.GameFlags = 4096
					c.Durability = []uint32{0, 1, 65535, 0x12345678}[j%4]
					c.Balance["QuestDurabilityMultiplier"] = v
					c.Balance["DefaultAmmoAmountQuest"] = v
					c.Balance["QuestStaffChargeMultiplier"] = v
					if j%3 == 0 {
						c.NilHealth = true
					}
					if j%4 == 0 {
						c.Effects = false
					}
					if name != "" {
						c.WeaponType = uint32(32 + len(specs)%3)
						c.ArmorType = c.WeaponType
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "creation-equipment", legacy.PortTestRoam(specs), creationHashes["creation-equipment"])
}

func TestObjectCreationContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	names := []string{"UrchinShaman", "Wizard", "WizardWhite", "Beholder", "Lich", "LichLord", "Demon", "WizardGreen", "WillOWisp", "OblivionHeart", "OblivionWierdling", "OblivionOrb"}
	// Cold lookup, partial failed lookup, and deliberately different cached IDs.
	for _, name := range names {
		for mode := 0; mode < 4; mode++ {
			op := 0
			if len(name) >= 8 && name[:8] == "Oblivion" {
				op = 1
			}
			s := creationBase(op)
			c := s.Callbacks.Creation
			c.Type = name
			if mode == 1 {
				c.Disabled = map[string]bool{name: true}
			}
			if mode == 2 {
				for i := range c.Cache {
					c.Cache[i] = uint32(100 + i)
				}
			}
			if mode == 3 {
				for i := range c.Cache {
					c.Cache[i] = uint32(23 + i)
				}
				c.TypeIndex = 27
				c.Type = ""
			}
			specs = append(specs, s)
		}
	}
	// Defined equipment branches independently cross mode, class, health, lookup
	// and modifier availability rather than correlating them with float cases.
	for _, op := range []int{1, 2} {
		for _, ind := range []uint16{15, 32, 33, 34} {
			for bits := 0; bits < 32; bits++ {
				s := creationBase(op)
				c := s.Callbacks.Creation
				c.TypeIndex = ind
				c.WeaponType = uint32(ind)
				c.ArmorType = uint32(ind)
				if bits&1 != 0 {
					s.Lifecycle.GameFlags = 4096
				}
				if bits&2 != 0 {
					c.NilHealth = true
				}
				if bits&4 != 0 {
					c.WeaponType = 99
					c.ArmorType = 99
				}
				c.Effects = bits&8 != 0
				if bits&16 != 0 {
					c.Subclass = 0x82
				} else {
					c.Class = 0x1008
					c.Subclass = 0x40000
				}
				specs = append(specs, s)
			}
		}
	}
	for _, op := range []int{3, 4, 5, 6, 7, 8, 10} {
		for _, flags := range []uint32{0, 4, 0x40, 0xffffffff} {
			for _, fill := range []byte{0, 0x5a, 0xff} {
				s := creationBase(op)
				s.Callbacks.Creation.Flags = flags
				s.Callbacks.Creation.InitFill = fill
				s.Callbacks.Creation.UseFill = fill
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "creation-contracts", legacy.PortTestRoam(specs), creationHashes["creation-contracts"])
}
func TestObjectCreationCloud(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, fps := range []uint32{0, 1, 30, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, life := range []float64{0, 1.25, -1.25, math.Nextafter(1, 0), 1e20, math.Inf(1), math.NaN()} {
			for _, warm := range []bool{false, true} {
				s := creationBase(9)
				c := s.Callbacks.Creation
				s.Owner.FPS = fps
				c.Balance["ToxicCloudLifetime"] = life
				if warm {
					c.Cache[12] = 10
				}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "creation-cloud", legacy.PortTestRoam(specs), creationHashes["creation-cloud"])
}

func TestObjectCreationCorpus(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	state := uint32(0x54c0c0)
	next := func() uint32 { state ^= state << 13; state ^= state >> 17; state ^= state << 5; return state }
	for n := 0; n < 2048; n++ {
		op := int(next() % 11)
		s := creationBase(op)
		c := s.Callbacks.Creation
		s.Seed = int(next() % 1000)
		s.Owner.FPS = next()
		s.Lifecycle.GameFlags = []uint32{0, 2048, 4096, 6144}[next()%4]
		c.TypeIndex = uint16(15 + next()%20)
		c.Class = []uint32{2, 8, 0x1000008, 0x1008}[next()%4]
		c.Subclass = []uint32{0, 2, 4, 8, 0x80, 0x82, 0x40000, 0x47f0000}[next()%8]
		c.Flags = next()
		c.WeaponType = uint32(c.TypeIndex)
		c.ArmorType = uint32(c.TypeIndex)
		if next()%3 == 0 {
			c.WeaponType = 99
		}
		if next()%3 == 0 {
			c.ArmorType = 99
		}
		c.Durability = next()
		c.NilHealth = next()%4 == 0
		c.Effects = next()%2 == 0
		c.InitFill = byte(next())
		c.UseFill = byte(next())
		c.AutoFlag = byte(next())
		for _, key := range []string{"QuestDurabilityMultiplier", "DefaultAmmoAmount", "DefaultAmmoAmountQuest", "QuestStaffChargeMultiplier", "ToxicCloudLifetime"} {
			c.Balance[key] = float64(math.Float32frombits(next()))
		}
		if next()%2 == 0 {
			for i := 0; i < 12; i++ {
				c.Cache[i] = uint32(23 + i)
			}
			c.Cache[12] = 10
		}
		for id := uint32(1); id <= 136; id++ {
			s.Spells.Permissions = append(s.Spells.Permissions, [2]uint32{id, next()})
		}
		specs = append(specs, s)
	}
	callbackHash(t, "creation-corpus", legacy.PortTestRoam(specs), creationHashes["creation-corpus"])
}

func TestObjectCreationPrecision(t *testing.T) {
	s := creationBase(1)
	s.Lifecycle.GameFlags = 4096
	s.Callbacks.Creation.Class = 0x1008
	s.Callbacks.Creation.Subclass = 0x40000
	s.Callbacks.Creation.UseFill = 0
	s.Callbacks.Creation.Balance["QuestStaffChargeMultiplier"] = math.MaxFloat32
	// C rounds the balance load to float32, but keeps the doubled multiplier
	// in double precision: 0*double(2*MaxFloat32) is zero, not 0*float32(+Inf).
	callbackHash(t, "creation-precision", legacy.PortTestRoam([]legacy.PortTestRoamSpec{s}), creationHashes["creation-precision"])
}

var creationHashes = map[string]string{
	"creation-smoke":     "369dea251575244f090a0f35caaf64ec5a61c05d549d39e5f87d928d9cafdda1",
	"creation-auto":      "f181718a31e4bd9f548da39120ddd20e09eebc0479d7fea263fdcdf9ccb69e17",
	"creation-equipment": "0173bc94d782ea6c63cff2d7615360fbb721896b7dc930e671a72c84291c209d",
	"creation-contracts": "2a3b17723d0baad95c04351d2130ced551bdfb1b97eab6a586383a95e75be56c",
	"creation-cloud":     "814d05ad41ae7560735566cb115907b4a72e2f3f1fd261af6c9c9d8143e1da25",
	"creation-corpus":    "6bf58d5246c63bef4e22f9be339513b53df6ac01280234d1545235ec310b9ee9",
	"creation-precision": "cdbbad34bf6bf71b987296520e88aa6bc54a5b5b3158532305081cd002eff176",
}
