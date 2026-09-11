//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"strings"
	"testing"
)

var damageHashes = map[string]string{
	"damage-admission-02": "c7dc4a0611c39e210254813f2d80bf826f736d6ce7d7af57ba28a9711d385ca4",
	"damage-admission-14": "cc9452b9b8cc6234fbaccb3336f3dc2a1fdc88df7428bc149fc92f98af51fcf6",
	"damage-armor":        "049325fd1311d09a48256759f43879739655b0ceccd76ec4a74e27c6804e74fe",
	"damage-ball":         "26d0c3f8a3904c027549935a256a6d79de5fedff17a3c6ba3e4bd65ff1efd823",
	"damage-basic-02":     "4bff2e78ce703bca8e2bb3b365da49f0b51c1b4a8a618a55cf41efdc0a5ed4ca",
	"damage-basic-09":     "092eb30f839233e318249f40abc99647dab3ae2230b7b61960e95a1f7a05ebb7",
	"damage-basic-10":     "092eb30f839233e318249f40abc99647dab3ae2230b7b61960e95a1f7a05ebb7",
	"damage-basic-14":     "dcada449448bbc41782db1456e0a076fc377dfb01b831d0d2918ce99b31cc390",
	"damage-basic-20":     "4bff2e78ce703bca8e2bb3b365da49f0b51c1b4a8a618a55cf41efdc0a5ed4ca",
	"damage-basic-21":     "4bff2e78ce703bca8e2bb3b365da49f0b51c1b4a8a618a55cf41efdc0a5ed4ca",
	"damage-basic-22":     "1d28be03178fa224f1c6191fd7746e1fea23031485d9464ed81bd512562a4aca",
	"damage-basic-23":     "f48a2a099d0db02a3596c5bb94c806704fa6b7742ffeafa365c72a6f68d33ccb",
	"damage-basic-24":     "f5cc1849d66761cb8bd245f36935f9b13ffbdee1b091a9f7bb65294bc456b3f5",
	"damage-blocking":     "d72c0a90e03295d120e8c1bbe97ddf8f5816d4d201a376a339c5b5141ac888c7",
	"damage-buffs-02":     "14c787ed80afc609b9c197e015da2ccc0b9fcc5ff428673da789cb0b6bd33f02",
	"damage-buffs-14":     "0d0efe356aa6c01326d99fef3f2a42018705d3747b16370d9d5324959ce78508",
	"damage-conductivity": "060177bf2937736b4fd3a9536211dce6b2e03b2f137c167bcb861bf0e63d2772",
	"damage-equipment-06": "d24a9ea17755339642c1bda520da771f773af3c8e6751a603d24ed6e47133877",
	"damage-equipment-07": "8861c43d64d08b3b6d23aaad3e6c99595ed00aaed79e5ef63e04634db46d559a",
	"damage-equipment-08": "8861c43d64d08b3b6d23aaad3e6c99595ed00aaed79e5ef63e04634db46d559a",
	"damage-equipment-09": "7be0a098be6d79bdd22114cdbfacdc97fc32f0bb183609594a5e3ef8c22a624b",
	"damage-equipment-10": "fecd5120836d6082fbba0efe3b87fcba24a2bf47f640c63551031f842451fabe",
	"damage-equipment-12": "82d7a8f3db268eb030acfb0f308b51ffeb894b0480b36169a38e114740bc5c32",
	"damage-equipment-16": "20ec2f556d08887edc28dcdef39506df747210d4aedb44a244fab2f41435e8f8",
	"damage-equipment-17": "7fe4956e6e1b2c3654093d7c87d86aabce4da5aa7d809eb9a9080ba62a1a93ab",
	"damage-equipment-18": "8ed0a14f8a9a6991e8676211d34be6ec9ad513a9a39d2f36b4df2e5b90317060",
	"damage-equipment-19": "8ed0a14f8a9a6991e8676211d34be6ec9ad513a9a39d2f36b4df2e5b90317060",
	"damage-fraction":     "9964c112a90ac48a7ae9f8dafe21e740d89e982318344bbab371b3fc12857154",
	"damage-generator":    "ede19cee3cdf5378493d8b76b103a9e26e14f6b2c8fab1e5e440e0b08d139f00",
	"damage-modifiers-02": "37061e516a70c4e9954b35a1b97b22c110ffb1d7289d1d74f6c4ae8c7de19fa4",
	"damage-modifiers-04": "06ccfed4fef3f2d90b82d4c80eebb932e52e9462049af27c9161740bf3c255ab",
	"damage-modifiers-05": "63619f54f4ec0bdd00ecaa6d979715fd9d99bd32ed1513f001414aaf2e175705",
	"damage-modifiers-11": "a791a4c3ed6b8ce0879d7a576a0d655dbd25ff9b31e6da0d69b5c56126839e60",
	"damage-modifiers-13": "a791a4c3ed6b8ce0879d7a576a0d655dbd25ff9b31e6da0d69b5c56126839e60",
	"damage-names":        "211cb3ec6b6902f9a54a947a371041c067493cd35dd7566486e760d147a190f0",
	"damage-nil":          "e0e132fb9df8db778230894dcf81ade50342a3dee7ba1afdaa0881d7a98f78bc",
	"damage-positive":     "97ae2a3328ceb9765ebe811b5eab63609642b6eeaec74b8dbb2a73ce9054db79",
	"damage-reflect":      "751b5d49b283bda180746a401f3235a7c9a1e9106c09df7c3be1b8f2740ae1e5",
	"damage-resistance":   "47e79497f434f4aaa4a922c0e8d5171e86ba1238773ff299b967908ac8e90527",
	"damage-shield":       "a897e3a66c720431ba702a1aeaafbd6a01d1cac9c8fb429798cc8cf1fad28379",
	"damage-source-02":    "303cd589cd08acd42089e4ace3b235c8ed0691bb0cead627a841d584a17b060d",
	"damage-source-14":    "e7ae2b20cc1e7fa8bcf5540c62d325917aaa370df61bbd0e4cff2b999e516674",
}

func damageBase(op int) legacy.PortTestRoamSpec {
	s := attackBase()
	p := s.Callbacks.Shop
	a := p.TemporaryUpdates.World.Objectives.Attack
	a.Actor = 1
	a.Damage = &legacy.PortTestDamageSpec{Amount: 12, FloatBits: math.Float32bits(12.5)}
	a.RecordWords[0] = 12
	p.Inventory.Linked = nil
	p.Inventory.Owned = nil
	p.Equipment.ActiveWeapon = 0
	p.Equipment.WeaponFlags = 0
	p.Resources.DieCallback = true
	p.Inventory.Gameplay = 1
	p.Sequence = []legacy.PortTestShopAction{{Op: 1100 + op}}
	return s
}
func damageHash(t *testing.T, name string, s []legacy.PortTestRoamSpec) {
	t.Helper()
	callbackHash(t, name, effectsTimedRun(t, s), damageHashes[name])
}
func TestDamageNames(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	names := append(legacy.PortTestDamageNames(), "", "unknown", "Fire ")
	for _, name := range names {
		for _, value := range []string{name, strings.ToLower(name), strings.ToUpper(name)} {
			s := damageBase(0)
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Damage.Name = value
			cases = append(cases, s)
		}
	}
	r := effectsTimedRun(t, cases)
	for i, v := range r {
		want := uint32(18)
		if i < 54 && i%3 != 1 {
			want = uint32(i / 3)
		}
		if got := v.Callbacks.Shop.Sequence[0].Return; got != want {
			t.Fatalf("name %q return %d want %d", cases[i].Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Damage.Name, got, want)
		}
	}
	callbackHash(t, "damage-names", r, damageHashes["damage-names"])
}
func TestDamageBasic(t *testing.T) {
	for _, op := range []int{2, 9, 10, 14, 20, 21, 22, 23, 24} {
		t.Run(fmt.Sprint(op), func(t *testing.T) {
			var cases []legacy.PortTestRoamSpec
			for _, subject := range []int{1, 2, 3} {
				for _, amount := range []int32{0, 1, 12, 49, 50, 100} {
					for _, kind := range []int32{0, 1, 2, 5, 7, 9, 11, 15, 17} {
						s := damageBase(op)
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
			damageHash(t, fmt.Sprintf("damage-basic-%02d", op), cases)
		})
	}
}
func TestDamageFraction(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 2, 3} {
		for _, value := range []float32{-12.75, -1, -0.5, 0, 0.125, 0.5, 0.999, 1, 12.75, 65535.5} {
			s := damageBase(15)
			p := s.Callbacks.Shop
			p.Resources.Subject = subject
			p.TemporaryUpdates.World.Objectives.Attack.Damage.FloatBits = math.Float32bits(value)
			p.Sequence = append(p.Sequence, p.Sequence[0], p.Sequence[0])
			cases = append(cases, s)
		}
	}
	damageHash(t, "damage-fraction", cases)
}

func TestDamagePositive(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{2, 14, 21, 22} {
		s := damageBase(op)
		cases = append(cases, s)
	}
	r := effectsTimedRun(t, cases)
	for i, v := range r {
		if hp := v.Callbacks.Shop.Sequence[0].ResourceData[2][0] & 65535; hp != 38 {
			t.Fatalf("case %d HP=%d want 38", i, hp)
		}
	}
	callbackHash(t, "damage-positive", r, damageHashes["damage-positive"])
}
func TestDamageAdmission(t *testing.T) {
	for _, op := range []int{2, 14} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 3} {
			for _, flags := range []uint32{0, 2, 0x20, 0x8000} {
				for _, buff := range []uint32{0, 1 << 23} {
					for _, frame := range []uint32{0, 1, 4, 0xffffffff} {
						s := damageBase(op)
						p := s.Callbacks.Shop
						p.Resources.Subject = subject
						p.Resources.Flags = flags
						p.Resources.Buffs = buff
						s.Owner.Frame = frame
						cases = append(cases, s)
					}
				}
			}
		}
		damageHash(t, fmt.Sprintf("damage-admission-%02d", op), cases)
	}
}
func TestDamageModifiers(t *testing.T) {
	for _, op := range []int{2, 4, 5, 11, 13} {
		var cases []legacy.PortTestRoamSpec
		for _, mask := range []uint32{0, 1, 2, 4, 8, 15} {
			for _, output := range []uint32{0, 1, 37} {
				for _, set := range []bool{false, true} {
					s := damageBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Inventory.Linked = []int{0}
					p.Items[0].Mods = [4]bool{true, true, true, true}
					a.Damage.Weapon = 3
					a.Damage.DefendMask = mask
					a.Damage.PreDamageMask = mask
					a.Damage.SetOutput = set
					a.Damage.Output = output
					if op == 11 || op == 13 {
						a.Actor = 3
						a.Damage.Source = 100
						a.Damage.Output = math.Float32bits(float32(output))
						p.Items[0].Health = true
						p.Items[0].HP = 50
						p.Items[0].MaxHP = 100
						p.TemporaryUpdates.UpdateWords[0][0] = math.Float32bits(0.75)
					}
					p.Sequence = append(p.Sequence, p.Sequence[0])
					cases = append(cases, s)
				}
			}
		}
		damageHash(t, fmt.Sprintf("damage-modifiers-%02d", op), cases)
	}
}
func TestDamageReflect(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, pos := range [][2]float32{{512, 512}, {513, 512}, {512, 513}, {500, 530}, {512.00006, 512.00006}} {
		for _, vel := range [][2]float32{{0, 0}, {3, 4}, {-3, 4}, {1e-20, -1e-20}} {
			for _, dir := range []uint32{0, 127, 128, 255, 256, 32767, 65535} {
				s := damageBase(1)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.Actor = 3
				a.Damage.Source = 100
				p.TemporaryUpdates.ItemWords[0][56] = math.Float32bits(pos[0])
				p.TemporaryUpdates.ItemWords[0][60] = math.Float32bits(pos[1])
				p.TemporaryUpdates.ItemWords[0][80] = math.Float32bits(vel[0])
				p.TemporaryUpdates.ItemWords[0][84] = math.Float32bits(vel[1])
				p.TemporaryUpdates.ItemWords[0][124] = dir
				cases = append(cases, s)
			}
		}
	}
	damageHash(t, "damage-reflect", cases)
}

func TestDamageEquipmentHelpers(t *testing.T) {
	for _, op := range []int{6, 7, 8, 9, 10, 12, 16, 17, 18, 19} {
		var cases []legacy.PortTestRoamSpec
		for _, class := range []uint32{0x1000000, 0x2000000} {
			for _, flags := range []uint32{0, 0x100, 0x120} {
				for _, value := range []float32{0, 0.25, 1, 12.75} {
					s := damageBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Items[0].Class = class
					p.Items[0].Flags = flags
					p.Items[0].Subclass = 0x100
					p.Items[0].Health = true
					p.Items[0].HP = 50
					p.Items[0].MaxHP = 100
					p.Inventory.Linked = []int{0}
					a.Damage.Weapon = 3
					a.Damage.FloatBits = math.Float32bits(value)
					a.Damage.Amount = 0x100
					a.Damage.Kind = 12
					if op == 9 || op == 10 {
						a.Actor = 3
						a.Damage.Amount = 12
					}
					if op == 12 {
						a.Damage.Source = 3
						a.Damage.Amount = 50
						a.Damage.Kind = int32(value)
						a.Damage.PlayerIndex = 1
					}
					// The weapon lookup helper requires a matching equipped item.
					if op == 18 || op == 19 {
						p.Items[0].Flags = 0x100 | (flags & 0x20)
					}
					if a.Actor == 1 {
						a.UpdateWords = map[int]uint32{228: math.Float32bits(1)}
					}
					cases = append(cases, s)
				}
			}
		}
		damageHash(t, fmt.Sprintf("damage-equipment-%02d", op), cases)
	}
}
func TestDamageResistance(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, subclass := range []uint32{0, 0x10, 0x200, 0x400, 0x800, 0xe00} {
		for _, kind := range []int32{0, 1, 5, 7, 9, 12, 17} {
			for _, amount := range []int32{0, 1, 12, 51} {
				s := damageBase(2)
				p := s.Callbacks.Shop
				p.Resources.Subject = 3
				p.Resources.Subclass = subclass
				d := p.TemporaryUpdates.World.Objectives.Attack.Damage
				d.Kind = kind
				d.Amount = amount
				cases = append(cases, s)
			}
		}
	}
	damageHash(t, "damage-resistance", cases)
}
func TestDamageSource(t *testing.T) {
	for _, op := range []int{2, 14} {
		var cases []legacy.PortTestRoamSpec
		for _, source := range []int{0, 1, 101} {
			for _, weapon := range []int{0, 3, 4} {
				for _, mode := range []uint32{1, 2049, 4097} {
					for _, gameplay := range []uint32{0, 1} {
						s := damageBase(op)
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
		damageHash(t, fmt.Sprintf("damage-source-%02d", op), cases)
	}
}

func TestDamageArmor(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 3} {
		for _, coeff := range []float32{0, 0.25, 0.5, 1, 1.25} {
			for kind := int32(0); kind < 18; kind++ {
				for _, god := range []bool{false, true} {
					s := damageBase(14)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Resources.Subject = subject
					p.Resources.Subclass = 0x10
					p.Resources.God = god
					a.Damage.Kind = kind
					a.Damage.Amount = 7
					off := 228
					if subject == 3 {
						off = 2072
					}
					a.UpdateWords = map[int]uint32{off: math.Float32bits(coeff)}
					p.Sequence = append(p.Sequence, p.Sequence[0], p.Sequence[0])
					cases = append(cases, s)
				}
			}
		}
	}
	damageHash(t, "damage-armor", cases)
}
func TestDamageGenerator(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, hp := range []uint16{1, 33, 34, 66, 67, 100} {
		for _, amount := range []int32{0, 1, 20, 100} {
			for _, frame := range []uint32{0, 20, 21, 29, 30, 0xffffffff} {
				s := damageBase(25)
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
	damageHash(t, "damage-generator", cases)
}
func TestDamageBall(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, amount := range []int32{0, 29, 30, 31, 100} {
		for _, owned := range []bool{false, true} {
			s := damageBase(3)
			p := s.Callbacks.Shop
			a := p.TemporaryUpdates.World.Objectives.Attack
			a.Actor = 101
			a.Damage.Source = 100
			a.Damage.Amount = amount
			if owned {
				p.Inventory.Owned = []int{0}
				p.TemporaryUpdates.World.ItemNames = []string{"GameBall"}
			}
			cases = append(cases, s)
		}
	}
	r := effectsTimedRun(t, cases)
	for i, v := range r {
		u := worldObject(t, v, 0, 70000)
		want := uint32(0)
		if i%2 == 1 && i < 4 {
			want = 54000
		}
		if u[508/4] != want {
			t.Fatalf("ball case %d owner %d want %d", i, u[508/4], want)
		}
	}
	callbackHash(t, "damage-ball", r, damageHashes["damage-ball"])
}

func TestDamageBlocking(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, weapon := range []uint32{0, 0x400, 0x8000} {
		for _, shield := range []uint32{0, 0x1000000} {
			for _, state := range []uint32{0, 13, 16, 18} {
				for _, kind := range []int32{0, 9, 11, 16, 17} {
					for _, front := range []bool{false, true} {
						s := damageBase(14)
						p := s.Callbacks.Shop
						o := p.TemporaryUpdates.World.Objectives
						a := o.Attack
						a.Damage.Source = 101
						a.Damage.Weapon = 4
						a.Damage.Kind = kind
						p.Equipment.WeaponFlags = weapon
						p.Equipment.ArmorFlags = shield
						p.Inventory.Linked = []int{0}
						p.Items[0].Subclass = weapon | 2
						p.Items[0].Flags = 0x100
						p.Items[0].Health = true
						p.Items[0].HP = 50
						p.Items[0].MaxHP = 100
						p.EffectsUse.Balance["ItemDamageFromBlockPercentage"] = []float64{0.25}
						a.UpdateWords = map[int]uint32{88: state}
						a.ActorWords = map[int]uint32{124: 0}
						x := float32(480)
						if front {
							x = 540
						}
						p.TemporaryUpdates.ItemWords[1][56] = math.Float32bits(x)
						p.TemporaryUpdates.ItemWords[1][60] = math.Float32bits(512)
						p.TemporaryUpdates.ItemWords[1][72] = math.Float32bits(x)
						p.TemporaryUpdates.ItemWords[1][76] = math.Float32bits(512)
						cases = append(cases, s)
					}
				}
			}
		}
	}
	damageHash(t, "damage-blocking", cases)
}
func TestDamageBuffs(t *testing.T) {
	for _, op := range []int{2, 14} {
		var cases []legacy.PortTestRoamSpec
		for _, buff := range []uint32{0, 1, 1 << 22, 1 << 27} {
			for _, kind := range []int32{0, 5, 9, 12, 15, 16, 17} {
				for _, frame := range []uint32{0, 1, 4, 0xffffffff} {
					s := damageBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Resources.Buffs = buff
					s.Owner.Frame = frame
					a.Damage.Source = 101
					a.Damage.Weapon = 3
					a.Damage.Kind = kind
					p.EffectsUse.Balance["ShockDamage"] = []float64{1, 2, 3, 4, 12}
					p.TemporaryUpdates.ItemWords[0][56] = math.Float32bits(540)
					p.TemporaryUpdates.ItemWords[0][60] = math.Float32bits(512)
					cases = append(cases, s)
				}
			}
		}
		damageHash(t, fmt.Sprintf("damage-buffs-%02d", op), cases)
	}
}

func TestDamageShield(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, amount := range []int32{0, 1, 12, 49, 50, 51, 100} {
		for _, kind := range []int32{0, 5, 15} {
			for _, self := range []bool{false, true} {
				s := damageBase(2)
				p := s.Callbacks.Shop
				p.Resources.Buffs = 1 << 26
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.Damage.Amount = amount
				a.Damage.Kind = kind
				if self {
					a.Damage.Source = 1
				}
				cases = append(cases, s)
			}
		}
	}
	damageHash(t, "damage-shield", cases)
}
func TestDamageConductivity(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{14, 17} {
		for count := 0; count <= 3; count++ {
			for _, material := range []uint32{0, 0x10} {
				for _, flags := range []uint32{0, 0x100} {
					s := damageBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.Damage.Kind = 9
					for i := 0; i < count; i++ {
						p.Inventory.Linked = append(p.Inventory.Linked, i)
						p.Items[i].Class = 0x2000000
						p.Items[i].Subclass = 0x2000000
						p.Items[i].Flags = flags
						p.TemporaryUpdates.ItemWords[i][24] = material
					}
					cases = append(cases, s)
				}
			}
		}
	}
	damageHash(t, "damage-conductivity", cases)
}
func TestDamageNil(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{2, 7, 8, 9, 11, 13, 15, 21, 22, 23, 24} {
		s := damageBase(op)
		s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Actor = 0
		cases = append(cases, s)
	}
	damageHash(t, "damage-nil", cases)
}
