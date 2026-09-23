//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type objectCreateRegistryRow struct {
	Name              string
	Class, Flags      uint32
	Field5, Field38   uint32
	Health            [2]uint16
	Init, Update, Use []byte
}

func objectCreateRegistryBytes(p unsafe.Pointer, n uintptr) []byte {
	if p == nil || n == 0 {
		return nil
	}
	return append([]byte(nil), unsafe.Slice((*byte)(p), int(n))...)
}

func objectCreateRegistrySnapshot(name string, u *server.Object) objectCreateRegistryRow {
	r := objectCreateRegistryRow{
		Name: name, Class: uint32(u.ObjClass), Flags: uint32(u.ObjFlags),
		Field5: u.Field5, Field38: u.Field38,
		Init:   objectCreateRegistryBytes(u.InitData, 1724),
		Update: objectCreateRegistryBytes(u.UpdateData, 2200),
		Use:    objectCreateRegistryBytes(u.UseData.Ptr, 128),
	}
	if u.HealthData != nil {
		r.Health = [2]uint16{u.HealthData.Cur, u.HealthData.Max}
	}

	return r
}

func TestLifecycleRegistryCreateRealObjects(t *testing.T) {
	s := newObjectXferOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0x201000))
	cases := []struct {
		name  string
		class object.Class
	}{
		{"MonsterCreate", object.ClassMonster},
		{"ArmorCreate", object.ClassArmor},
		{"WeaponCreate", object.ClassWeapon},
		{"ObeliskCreate", object.ClassImmobile},
		{"AnimCreate", object.ClassImmobile},
		{"TriggerCreate", object.ClassTrigger},
		{"MonsterGeneratorCreate", object.ClassMonsterGenerator},
		{"RewardMarkerCreate", object.ClassSimple},
	}
	var rows []objectCreateRegistryRow
	for _, tc := range cases {
		before := [3]uint32{*memmap.PtrUint32(0x5D4594, 2491660), *memmap.PtrUint32(0x5D4594, 2491664), *memmap.PtrUint32(0x5D4594, 2491668)}
		sub := object.SubClass(0)
		if tc.name == "WeaponCreate" {
			sub = 2
		}
		u := s.PortTestCreateWithRegisteredCreate(tc.name, nil, tc.class, sub)
		if u == nil {
			t.Fatal("allocation", tc.name)
		}
		row := objectCreateRegistrySnapshot(tc.name, u)
		switch tc.name {
		case "MonsterCreate":
			m := u.UpdateDataMonster()
			for _, off := range []int{476, 484, 488, 1196, 2176, 2192} {
				if objectXferGetWord(u.UpdateData, off) != 0 {
					t.Fatalf("unexpected live pointer in isolated monster template at %d", off)
				}
			}
			if m.Aggression != 0.5 || m.Aggression2 != 0.5 || m.SightRange != 150 || m.ScriptLookingForEnemy.Func != -1 || m.ScriptDeath.Func != -1 {
				t.Fatalf("monster defaults %+v", m)
			}
		case "ArmorCreate", "WeaponCreate":
			if row.Health != [2]uint16{150, 150} {
				t.Fatalf("positive durability scaling %s: %v", tc.name, row.Health)
			}
			if tc.name == "WeaponCreate" && !bytes.Equal(row.Use[:3], []byte{30, 30, 0}) {
				t.Fatal("weapon ammo initialization")
			}
		case "ObeliskCreate":
			if *(*uint32)(u.UpdateData) != 50 || u.Field38 != ^uint32(0) {
				t.Fatalf("obelisk create state %d/%08x", *(*uint32)(u.UpdateData), u.Field38)
			}
		case "AnimCreate":
			if u.Field5 != 2 {
				t.Fatalf("animation status %x", u.Field5)
			}
		case "TriggerCreate":
			if !bytes.Equal(row.Update[54:60], []byte{90, 90, 90, 10, 10, 10}) {
				t.Fatalf("trigger defaults %v", row.Update[54:60])
			}
		case "MonsterGeneratorCreate":
			w := unsafe.Slice((*uint32)(u.UpdateData), 24)
			for _, i := range []int{13, 15, 17, 19} {
				if w[i] != ^uint32(0) {
					t.Fatalf("generator sentinel[%d]=%08x", i, w[i])
				}
			}
			if w[23] != 2 {
				t.Fatalf("generator mode %d", w[23])
			}
		case "RewardMarkerCreate":
			w := unsafe.Slice((*uint32)(u.InitData), 54)
			if w[0] != 255 || w[53] != 0 {
				t.Fatalf("marker defaults %d/%d", w[0], w[53])
			}
		}
		after := [3]uint32{*memmap.PtrUint32(0x5D4594, 2491660), *memmap.PtrUint32(0x5D4594, 2491664), *memmap.PtrUint32(0x5D4594, 2491668)}
		if before != after {
			t.Fatalf("create globals changed for %s: %v -> %v", tc.name, before, after)
		}
		rows = append(rows, row)
		objectCreateRegistryFree(s, u) // Keep generator update data live through its destructor.
	}
	spellbookCapture(t, "object-create-registry", rows, "f5bd472acce1ea179a129cc4a2f035954afc653729426d3121f1d368f2d88097")
}

func TestLifecycleRegistryCreateRawNilAndLateBinding(t *testing.T) {
	s := newObjectXferOwner(t)
	legacy.PortTestDeathForwardReset()
	nilObj := s.PortTestCreateWithRegisteredCreate("", nil, object.ClassSimple, 0)
	if nilObj == nil {
		t.Fatal("nil create allocation")
	}
	if ptr, count := legacy.PortTestDeathForwardSnapshot(); ptr != nil || count != 0 {
		t.Fatalf("nil callback dispatched %p/%d", ptr, count)
	}
	objectCreateRegistryFree(s, nilObj)

	forward := legacy.PortTestDeathForwardReset()
	rawObj := s.PortTestCreateWithRegisteredCreate("", forward, object.ClassSimple, 0)
	if rawObj == nil {
		t.Fatal("raw create allocation")
	}
	if ptr, count := legacy.PortTestDeathForwardSnapshot(); ptr != rawObj.CObj() || count != 1 {
		t.Fatalf("raw callback %p/%d want %p/1", ptr, count, rawObj.CObj())
	}
	objectCreateRegistryFree(s, rawObj)

	old := legacy.Nox_xxx_monsterCreateFn_54C480
	t.Cleanup(func() { legacy.Nox_xxx_monsterCreateFn_54C480 = old })
	calls := 0
	for index, marker := range []uint32{0x11223344, 0x55667788} {
		legacy.Nox_xxx_monsterCreateFn_54C480 = func(u *server.Object) { calls++; u.Field5 = marker }
		u := s.PortTestCreateWithRegisteredCreate("MonsterCreate", nil, object.ClassMonster, 0)
		if u == nil || u.Field5 != marker || calls != index+1 {
			t.Fatalf("late-bound monster create %08x", marker)
		}
		objectCreateRegistryFree(s, u)
	}
}

func objectCreateRegistryFree(s *server.Server, u *server.Object) {
	// FreeObject consumes UseData, but only clears these remaining pointers.
	owned := []unsafe.Pointer{u.InitData, u.UpdateData, u.CollideData, u.Field189, unsafe.Pointer(u.HealthData)}
	s.Objs.FreeObject(u)
	for _, p := range owned {
		if p != nil {
			alloc.FreePtr(p)
		}
	}
}
