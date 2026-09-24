//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestDamageIdentityLifetime(t *testing.T) {
	typ, freeType := alloc.New(server.ObjectType{})
	defer freeType()
	obj, freeObject := alloc.New(server.Object{})
	defer freeObject()
	seen := map[unsafe.Pointer]bool{}
	names := []string{"DefaultDamage", "SkeletonDamage", "PlayerDamage", "StoneDamage", "MechGolemDamage", "FlammableDamage", "BlackPowderDamage", "ArmorDamage", "WeaponDamage", "BallDamage", "MonsterGeneratorDamage"}
	for _, name := range names {
		key := server.PortTestDamageRegistry(name)
		if key == nil || seen[key] {
			t.Fatal("damage identity", name)
		}
		seen[key] = true
		typ.Damage = key
		obj.Damage = typ.Damage
		collisionRegistryGrow(128)
		if server.PortTestDamageRegistry(name) != key || typ.Damage != key || obj.Damage != key {
			t.Fatal("stored damage identity", name)
		}
		if name == "BallDamage" {
			if got := obj.CallDamageValue(nil, nil, -2147483648, 2147483647); got != 0 {
				t.Fatal("ball damage after GC", got)
			}
		}
	}
	if server.DefaultDamage != server.PortTestDamageRegistry("DefaultDamage") {
		t.Fatal("default damage identity")
	}
	oldDefault, oldPlayer := legacy.Nox_xxx_soundDefaultDamageSound_532E20, legacy.Nox_xxx_soundPlayerDamageSound_5328B0
	defer func() {
		legacy.Nox_xxx_soundDefaultDamageSound_532E20 = oldDefault
		legacy.Nox_xxx_soundPlayerDamageSound_5328B0 = oldPlayer
	}()
	calls := 0
	observe := func(a, b *server.Object) int {
		if a != obj || b != obj {
			t.Fatal("damage sound arguments")
		}
		calls++
		return -1
	}
	legacy.Nox_xxx_soundDefaultDamageSound_532E20 = observe
	legacy.Nox_xxx_soundPlayerDamageSound_5328B0 = observe
	for _, name := range []string{"DefaultDamageSound", "PlayerDamageSound"} {
		key := server.PortTestXferSoundRegistry(name, true)
		if key == nil || seen[key] {
			t.Fatal("damage sound identity", name)
		}
		seen[key] = true
		typ.DamageSound = key
		obj.DamageSound = typ.DamageSound
		collisionRegistryGrow(128)
		if server.PortTestXferSoundRegistry(name, true) != key || obj.DamageSound != key || typ.DamageSound != key {
			t.Fatal("stored damage sound identity", name)
		}
		obj.CallDamageSound(obj)
	}
	if len(seen) != 13 || calls != 2 {
		t.Fatal("damage identity counts", len(seen), calls)
	}
}

func TestDamageNameCStringBoundary(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var expected []uint32
	for i, name := range legacy.PortTestDamageNames() {
		for j, value := range []string{name + "\x00ignored", "\x00" + name, "unknown\x00" + name} {
			s := damageBase(0)
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Damage.Name = value
			cases = append(cases, s)
			want := uint32(18)
			if j == 0 {
				want = uint32(i)
			}
			expected = append(expected, want)
		}
	}
	rows := effectsTimedRun(t, cases)
	if len(rows) != len(expected) {
		t.Fatal("name boundary count", len(rows), len(expected))
	}
	for i, row := range rows {
		if got := row.Callbacks.Shop.Sequence[0].Return; got != expected[i] {
			t.Fatal("C string boundary", i, got, expected[i])
		}
	}
}
