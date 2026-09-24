//go:build porttest

package opennox

import (
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestItemIdentityLifetime(t *testing.T) {
	typ, freeType := alloc.New(server.ObjectType{})
	defer freeType()
	obj, freeObject := alloc.New(server.Object{})
	defer freeObject()
	seen := map[unsafe.Pointer]string{}
	families := []struct {
		kind  string
		names []string
	}{
		{"use", []string{"ConsumeUse", "ConsumeConfuseUse", "CastUse", "EnchantUse", "MushroomUse", "PotionUse", "FireWandUse", "ReadUse", "WarpReadUse", "WandUse", "WandCastUse", "SpellRewardUse", "AbilityRewardUse", "FieldGuideUse"}},
		{"drop", []string{"DefaultDrop", "ArmorDrop", "WeaponDrop", "TreasureDrop", "GlyphDrop", "PotionDrop", "TrapDrop", "FoodDrop", "CrownDrop", "AudEventDrop", "AnkhTradableDrop"}},
		{"pickup", []string{"DefaultPickup", "FoodPickup", "UsePickup", "ArmorPickup", "WeaponPickup", "OblivionPickup", "TreasurePickup", "TrapPickup", "PotionPickup", "GoldPickup", "AmmoPickup", "SpellBookPickup", "AbilityBookPickup", "CrownPickup", "AudEventPickup", "AnkhTradablePickup"}},
	}
	sizes := map[string]uintptr{
		"ConsumeUse": unsafe.Sizeof(server.ConsumeUseData{}), "ConsumeConfuseUse": unsafe.Sizeof(server.ConsumeUseData{}), "CastUse": unsafe.Sizeof(server.CastUseData{}), "EnchantUse": unsafe.Sizeof(server.EnchantUseData{}), "PotionUse": unsafe.Sizeof(server.PotionUseData{}), "ReadUse": 260, "WarpReadUse": 260, "WandUse": 116, "WandCastUse": 116, "SpellRewardUse": unsafe.Sizeof(server.SpellRewardUseData{}), "AbilityRewardUse": unsafe.Sizeof(server.AbilityRewardUseData{}), "FieldGuideUse": unsafe.Sizeof(server.FieldGuideUseData{}),
	}
	for _, family := range families {
		for _, name := range family.names {
			key, size := server.PortTestItemIdentity(family.kind, name)
			if key == nil || seen[key] != "" {
				t.Fatal("item identity", name, seen[key])
			}
			seen[key] = name
			if size != sizes[name] {
				t.Fatal("item data size", name, size, sizes[name])
			}
			switch family.kind {
			case "use":
				typ.Use = server.UseFuncPtr{Ptr: key}
				obj.Use = typ.Use
			case "drop":
				typ.Drop = server.DropFuncPtr{Ptr: key}
				obj.Drop = typ.Drop
			case "pickup":
				typ.Pickup = server.PickupFuncPtr{Ptr: key}
				obj.Pickup = typ.Pickup
			}
			collisionRegistryGrow(128)
			got, _ := server.PortTestItemIdentity(family.kind, name)
			if got != key {
				t.Fatal("item key after GC", name)
			}
			switch family.kind {
			case "use":
				if typ.Use.Ptr != key || obj.Use.Ptr != key || obj.Use.Get() == nil {
					t.Fatal("stored use key", name)
				}
			case "drop":
				if typ.Drop.Ptr != key || obj.Drop.Ptr != key || obj.Drop.Get() == nil {
					t.Fatal("stored drop key", name)
				}
			case "pickup":
				if typ.Pickup.Ptr != key || obj.Pickup.Ptr != key || obj.Pickup.Get() == nil {
					t.Fatal("stored pickup key", name)
				}
			}
		}
	}
	if len(seen) != 41 {
		t.Fatal("item identity count", len(seen))
	}
	for name, size := range map[string]uintptr{"AmmoUse": 3, "BowUse": 1} {
		key, n := server.PortTestItemIdentity("use", name)
		if key != nil || n != size {
			t.Fatal("nil use registration", name, key, n)
		}
	}
}

func TestItemMutableIdentityHooks(t *testing.T) {
	u, freeU := alloc.New(server.Object{})
	defer freeU()
	it, freeIt := alloc.New(server.Object{})
	defer freeIt()
	uses := []struct {
		name string
		slot *server.UseFunc
	}{
		{"ConsumeUse", &legacy.Nox_xxx_useConsume_53EE10}, {"ConsumeConfuseUse", &legacy.Nox_xxx_useCiderConfuse_53EF00}, {"CastUse", &legacy.Nox_xxx_useCast_53ED90}, {"EnchantUse", &legacy.Nox_xxx_useEnchant_53ED60}, {"MushroomUse", &legacy.Nox_xxx_useMushroom_53ECE0}, {"PotionUse", &legacy.Nox_xxx_usePotion_53EF70},
	}
	for _, tc := range uses {
		t.Run(tc.name, func(t *testing.T) {
			old := *tc.slot
			defer func() { *tc.slot = old }()
			key, _ := server.PortTestItemIdentity("use", tc.name)
			it.Use = server.UseFuncPtr{Ptr: key}
			for _, result := range []bool{false, true} {
				calls := 0
				*tc.slot = func(a, b *server.Object) bool {
					if a != u || b != it {
						t.Fatal("use arguments")
					}
					calls++
					return result
				}
				collisionRegistryGrow(128)
				if got := it.Use.Get()(u, it); got != result || calls != 1 {
					t.Fatal("late use hook", got, result, calls)
				}
			}
		})
	}
	pickups := []struct {
		name string
		slot *server.PickupFunc
	}{
		{"DefaultPickup", &legacy.Nox_xxx_pickupDefault_4F31E0}, {"PotionPickup", &legacy.Nox_xxx_pickupPotion_4F37D0}, {"AudEventPickup", &legacy.Nox_objectPickupAudEvent_4F3D50},
	}
	for _, tc := range pickups {
		t.Run(tc.name, func(t *testing.T) {
			old := *tc.slot
			defer func() { *tc.slot = old }()
			key, _ := server.PortTestItemIdentity("pickup", tc.name)
			it.Pickup = server.PickupFuncPtr{Ptr: key}
			for _, result := range []bool{false, true} {
				calls := 0
				*tc.slot = func(a, b *server.Object, x, y int) bool {
					if a != u || b != it || x != -2147483648 || y != 2147483647 {
						t.Fatal("pickup arguments")
					}
					calls++
					return result
				}
				collisionRegistryGrow(128)
				if got := it.Pickup.Get()(u, it, -2147483648, 2147483647); got != result || calls != 1 {
					t.Fatal("late pickup hook", got, result, calls)
				}
			}
		})
	}
	t.Run("AudEventDrop", func(t *testing.T) {
		old := legacy.Nox_objectDropAudEvent_4EE2F0
		defer func() { legacy.Nox_objectDropAudEvent_4EE2F0 = old }()
		key, _ := server.PortTestItemIdentity("drop", "AudEventDrop")
		it.Drop = server.DropFuncPtr{Ptr: key}
		pos := types.Ptf(math.Float32frombits(0x80000000), math.Float32frombits(0x7f800000))
		for _, result := range []bool{false, true} {
			calls := 0
			legacy.Nox_objectDropAudEvent_4EE2F0 = func(a, b *server.Object, p types.Pointf) bool {
				if a != u || b != it || math.Float32bits(p.X) != 0x80000000 || math.Float32bits(p.Y) != 0x7f800000 {
					t.Fatal("drop arguments")
				}
				calls++
				return result
			}
			collisionRegistryGrow(128)
			if got := it.Drop.Get()(u, it, pos); got != result || calls != 1 {
				t.Fatal("late drop hook", got, result, calls)
			}
		}
	})
}

func TestItemUseFullResultWords(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	var expected []uint32
	for _, value := range []int32{-2147483648, -1, 0, 1, 256, 2147483647} {
		for _, blocked := range []bool{false, true} {
			s := effectsUseBase()
			p := s.Callbacks.Shop
			result := value
			p.Inventory.UseResult = &result
			p.Inventory.Blocked = blocked
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F8E0}}
			specs = append(specs, s)
			want := uint32(value)
			if blocked {
				want = 1
			}
			expected = append(expected, want)
		}
	}
	rows := legacy.PortTestRoam(specs)
	if len(rows) != len(expected) {
		t.Fatal("use result count", len(rows), len(expected))
	}
	for i, row := range rows {
		if got := row.Callbacks.Shop.Sequence[0].Return; got != expected[i] {
			t.Fatal("full use result", i, got, expected[i])
		}
	}
}

func TestItemUseGetterIdentities(t *testing.T) {
	obj, freeObject := alloc.New(server.Object{})
	defer freeObject()
	for _, tc := range []struct {
		name string
		get  func() unsafe.Pointer
	}{
		{"PotionUse", legacy.Get_nox_xxx_usePotion_53EF70}, {"SpellRewardUse", legacy.Get_nox_xxx_useSpellReward_53F9E0}, {"AbilityRewardUse", legacy.Get_nox_xxx_useAbilityReward_53FAE0}, {"EnchantUse", legacy.Get_nox_xxx_useEnchant_53ED60}, {"CastUse", legacy.Get_nox_xxx_useCast_53ED90}, {"FieldGuideUse", legacy.Get_sub_53F930},
	} {
		key, _ := server.PortTestItemIdentity("use", tc.name)
		if tc.get() != key {
			t.Fatal("item getter", tc.name)
		}
		obj.SetUse(tc.get(), nil)
		collisionRegistryGrow(128)
		if obj.Use.Ptr != key || obj.Use.Get() == nil {
			t.Fatal("stored getter identity", tc.name)
		}
	}
	key, _ := server.PortTestItemIdentity("drop", "AudEventDrop")
	if legacy.Get_nox_objectDropAudEvent_4EE2F0() != key {
		t.Fatal("drop getter")
	}
}
