//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

var resourceHashes = map[string]string{
	"resources-gold":              "131fd007cfdcbdb879fc0e6663b9abc60d3ee0af646e6d87d20b124cd4c1b4d1",
	"resources-mana":              "1f8b1195e113076d24e85e9c80389a3349d586c8f28e9f11462663953ebab329",
	"resources-health":            "d98a1cd17dd15abc2a3f15240b55891ea45064537cfe77b68fcd49a3e6c8cedb",
	"resources-gold-pickup":       "c21a067b6e21dee0f14f5125269be35f96c6c7aab07285c705ccfc8cbb7bbc1f",
	"resources-damage":            "67f5eef3db129529cfdc6d42630448affa376b270549a66d62fd42a6d13743e1",
	"resources-poison-mutation":   "43b26b1fe972eec9a408d369f4aafcebb97bf5d6f7793f857ff5e4c5dbbd2e19",
	"resources-poison-activation": "22d75834608a91ee322498317061ed9a41942fa0f87ac014287670161609b074",
	"resources-guarded-accessors": "7bb25e9577baa705bbd07eb3281c9178fca4ad933a0ca156db5fdae40fb5773c",
	"resources-owned-health-sync": "18ec3ba62d0d8edc8cd2eb4542b7539f6c3019d924fc7f321b354676aafd187a",
	"resources-healing-disabled":  "cd2c3b6bac8c3d494370220e8cc75e432e50caa365c53322da0f07c3ac0c5c99",
	"resources-pickup-fallback":   "588c9a09848a328d6567cd352e4157b9a7cd501f840ed3e53467156693646935",
}

func resourceBase() legacy.PortTestRoamSpec {
	s := shopBase(4)
	p := s.Callbacks.Shop
	p.Resources = &legacy.PortTestResourceSpec{Subject: 1, HP: 50, OldHP: 49, MaxHP: 100, Mana: 50, OldMana: 49, MaxMana: 100}
	p.Gold = [2]uint32{1000, 2000}
	return s
}
func TestResourcesGold(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, gold := range []uint32{0, 1, 100, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, amount := range []uint32{0, 1, 99, 100, 101, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, protected := range []bool{false, true} {
				s := resourceBase()
				p := s.Callbacks.Shop
				p.Gold[0] = gold
				p.Resources.Protected = protected
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceObjectGold}, {Op: legacy.PortTestResourceGetGold}, {Op: legacy.PortTestResourceGoldAdd, Value: amount}, {Op: legacy.PortTestResourceGetGold}, {Op: legacy.PortTestResourceGoldSub, Value: amount}, {Op: legacy.PortTestResourceGoldSet, Value: amount}, {Op: legacy.PortTestResourceObjectGold}}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		p := specs[i].Callbacks.Shop
		q := v.Callbacks.Shop.Sequence
		if q[0].Return != p.Gold[0] || q[1].Return != p.Gold[0] || q[3].Return != p.Gold[0]+p.Sequence[2].Value {
			t.Fatalf("case %d gold getters/addition", i)
		}
	}
	callbackHash(t, "resources-gold", r, resourceHashes["resources-gold"])
}
func TestResourcesMana(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 0; subject < 4; subject++ {
		for _, mana := range [][2]uint16{{0, 0}, {1, 0}, {50, 100}, {100, 100}, {200, 100}, {65535, 65535}} {
			for _, amount := range []uint32{0, 1, 50, 100, 101, 32767, 32768, 65535, 0x80000000, 0xffffffff} {
				for variant := 0; variant < 4; variant++ {
					s := resourceBase()
					p := s.Callbacks.Shop
					sp := p.Resources
					sp.Subject, sp.Mana, sp.MaxMana = subject, mana[0], mana[1]
					sp.God, sp.Protected = variant&1 != 0, variant&2 != 0
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceGetMana}, {Op: legacy.PortTestResourceGetMaxMana}, {Op: legacy.PortTestResourceManaAdd, Value: amount}, {Op: legacy.PortTestResourceManaSub, Value: amount}, {Op: legacy.PortTestResourceSetMaxMana, Value: amount}, {Op: legacy.PortTestResourceManaRefresh}, {Op: legacy.PortTestResourceGetMana}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "resources-mana", legacy.PortTestRoam(specs), resourceHashes["resources-mana"])
}
func TestResourcesHealth(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject < 4; subject++ {
		for _, hp := range [][2]uint16{{0, 0}, {1, 0}, {0, 100}, {50, 100}, {100, 100}, {200, 100}, {65535, 65535}} {
			for _, amount := range []uint32{0, 1, 100, 32768, 65535, 65536, 0x80000000, 0xffffffff} {
				for variant := 0; variant < 4; variant++ {
					s := resourceBase()
					p := s.Callbacks.Shop
					sp := p.Resources
					sp.Subject, sp.HP, sp.MaxHP = subject, hp[0], hp[1]
					sp.NoHealth, sp.Protected = variant&1 != 0, variant&2 != 0
					sp.Owner = subject == 3
					sp.Holder = subject == 2 && variant&2 != 0
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceGetHP}, {Op: legacy.PortTestResourceGetMaxHP}, {Op: legacy.PortTestResourceAdjustHP, Value: amount}, {Op: legacy.PortTestResourceSetHP, Value: amount}, {Op: legacy.PortTestResourceHPHistory}, {Op: legacy.PortTestResourceSetMaxHP, Value: amount}, {Op: legacy.PortTestResourceRestoreHP}, {Op: legacy.PortTestResourceInformOwner}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "resources-health", legacy.PortTestRoam(specs), resourceHashes["resources-health"])
}
func TestResourcesGoldPickup(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, gold := range []uint32{0, 1, 1000, 0xffffffff} {
		for _, amount := range []uint32{0, 1, 100, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, protected := range []bool{false, true} {
				s := resourceBase()
				p := s.Callbacks.Shop
				p.Gold[0] = gold
				p.Resources.GoldItem = amount
				p.Resources.Protected = protected
				p.Items = []legacy.PortTestShopItem{{Type: 26, Class: 8}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceGoldPickup}, {Op: legacy.PortTestResourceGetGold}}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		p := specs[i].Callbacks.Shop
		q := v.Callbacks.Shop.Sequence
		if q[0].Return != 1 || q[1].Return != p.Gold[0]+p.Resources.GoldItem {
			t.Fatalf("case %d gold pickup", i)
		}
	}
	callbackHash(t, "resources-gold-pickup", r, resourceHashes["resources-gold-pickup"])
}

func TestResourcesDamage(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 0; subject < 4; subject++ {
		for _, hp := range [][2]uint16{{0, 0}, {1, 100}, {50, 100}, {65535, 65535}} {
			for _, amount := range []uint32{0, 1, 49, 50, 51, 0x7fffffff, 0x80000000, 0xffffffff} {
				for variant := 0; variant < 8; variant++ {
					s := resourceBase()
					p := s.Callbacks.Shop
					sp := p.Resources
					sp.Subject, sp.HP, sp.MaxHP = subject, hp[0], hp[1]
					sp.God, sp.DieCallback, sp.Protected = variant&1 != 0, variant&2 != 0, variant&4 != 0
					sp.Harpoon = subject == 1
					sp.PlayerClass = byte(variant % 2)
					if variant == 7 {
						sp.Flags = 0x8000
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceDamage, Value: amount}, {Op: legacy.PortTestResourceDamage, Value: amount}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "resources-damage", legacy.PortTestRoam(specs), resourceHashes["resources-damage"])
}
func TestResourcesPoisonMutation(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 0; subject < 4; subject++ {
		for _, poison := range []byte{0, 1, 127, 255} {
			for _, amount := range []uint32{0, 1, 100, 255, 256, 1000, 0xffffffff} {
				for variant := 0; variant < 4; variant++ {
					s := resourceBase()
					p := s.Callbacks.Shop
					sp := p.Resources
					sp.Subject, sp.Poison, sp.PoisonTime, sp.PoisonTimer = subject, poison, 1234, 999
					sp.NoHealth = variant&1 != 0
					sp.Owner = variant&2 != 0
					sp.Status = 1024
					sp.Subclass = []uint32{0, 0x10, 0x80, 0x90}[variant]
					if variant&1 != 0 {
						s.Lifecycle.GameFlags = 2048
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourcePoisonReduce, Value: amount}, {Op: legacy.PortTestResourcePoisonSet, Value: amount}, {Op: legacy.PortTestResourcePoisonRemove}, {Op: legacy.PortTestResourcePoisonRemove}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "resources-poison-mutation", legacy.PortTestRoam(specs), resourceHashes["resources-poison-mutation"])
}
func TestResourcesPoisonActivation(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject < 4; subject++ {
		for variant := 0; variant < 12; variant++ {
			for _, seed := range []int{0, 1, 2, 7, 42, 123} {
				s := resourceBase()
				s.Seed = seed
				p := s.Callbacks.Shop
				sp := p.Resources
				sp.Subject = subject
				sp.Poison = []byte{0, 1, 100, 255}[variant%4]
				sp.PoisonTime = 1234
				amount := []uint32{0, 1, 10, 256, 0xffffffff, 0x7fffffff}[variant%6]
				maximum := []int{0, 1, 100, 255, -1, 0x7fffffff}[variant%6]
				if variant == 0 {
					sp.Flags = 2
				}
				if variant == 1 {
					sp.Buffs = 1 << 23
				}
				if variant == 2 {
					sp.Status = 1
				}
				if variant == 3 {
					sp.Subclass = 0x200
				}
				if variant >= 4 {
					sp.Buffs = 1 << 18
					p.Balance["PoisonSpellProtection"] = []float64{0, .25, .7, .9}[variant%4]
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourcePoison, Value: amount, Item: maximum}, {Op: legacy.PortTestResourcePoison, Value: amount, Item: maximum}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "resources-poison-activation", legacy.PortTestRoam(specs), resourceHashes["resources-poison-activation"])
}
func TestResourcesGuardedAccessors(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, subject := range []int{0, 2, 3} {
		s := resourceBase()
		p := s.Callbacks.Shop
		p.Resources.Subject = subject
		p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceGetHP}, {Op: legacy.PortTestResourceGetMaxHP}, {Op: legacy.PortTestResourceSetMaxHP, Value: 123}, {Op: legacy.PortTestResourceRestoreHP}, {Op: legacy.PortTestResourceHPHistory}, {Op: legacy.PortTestResourceInformOwner}, {Op: legacy.PortTestResourceGoldSet, Value: 123}, {Op: legacy.PortTestResourceObjectGold}}
		specs = append(specs, s)
	}
	callbackHash(t, "resources-guarded-accessors", legacy.PortTestRoam(specs), resourceHashes["resources-guarded-accessors"])
}

func TestResourcesOwnedHealthSync(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, extra := range []uint32{0, 0x400000, 0x20000000} {
		for variant := 0; variant < 8; variant++ {
			for _, amount := range []uint32{0, 1, 65535, 65536} {
				s := resourceBase()
				p := s.Callbacks.Shop
				sp := p.Resources
				sp.Subject = 3
				sp.Subclass = 0x80
				sp.ExtraClass = extra
				sp.SyncSeed = 0xa5a500ff
				sp.Protected = variant&1 != 0
				sp.Owner = variant&2 != 0
				sp.Holder = variant&4 != 0
				sp.OtherHolder = !sp.Holder
				sp.OtherOwner = !sp.Owner
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceSetHP, Value: amount}, {Op: legacy.PortTestResourceSetHP, Value: amount}, {Op: legacy.PortTestResourceInformOwner}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "resources-owned-health-sync", legacy.PortTestRoam(specs), resourceHashes["resources-owned-health-sync"])
}
func TestResourcesHealingDisabled(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 2, 3} {
		for _, amount := range []uint32{0, 1, 100, 0xffffffff} {
			s := resourceBase()
			p := s.Callbacks.Shop
			p.Resources.Subject = subject
			s.Lifecycle.GameFlags = 0x4000000
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceGetHP}, {Op: legacy.PortTestResourceAdjustHP, Value: amount}, {Op: legacy.PortTestResourceGetHP}}
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		q := v.Callbacks.Shop.Sequence
		if q[0].Return != q[2].Return {
			t.Fatalf("case %d healing disabled", i)
		}
	}
	callbackHash(t, "resources-healing-disabled", r, resourceHashes["resources-healing-disabled"])
}
func TestResourcesPickupFallback(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, subject := range []int{2, 3} {
		for _, result := range []bool{false, true} {
			for _, flags := range []uint32{0, 1, 2, 0xffffffff} {
				s := resourceBase()
				p := s.Callbacks.Shop
				p.Resources.Subject = subject
				p.Resources.PickupResult = result
				p.Resources.GoldItem = 100
				p.Items = []legacy.PortTestShopItem{{Type: 26, Class: 8}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestResourceGoldPickup, Value: flags}}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		want := uint32(0)
		if specs[i].Callbacks.Shop.Resources.PickupResult {
			want = 1
		}
		if v.Callbacks.Shop.Sequence[0].Return != want {
			t.Fatalf("case %d fallback result", i)
		}
	}
	callbackHash(t, "resources-pickup-fallback", r, resourceHashes["resources-pickup-fallback"])
}
