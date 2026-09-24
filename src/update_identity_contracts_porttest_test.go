//go:build porttest

package opennox

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestUpdateIdentityLifetime(t *testing.T) {
	typ, free := alloc.New(server.ObjectType{})
	defer free()
	obj, free2 := alloc.New(server.Object{})
	defer free2()
	names := legacy.PortTestUpdateRegistryNames()
	if len(names) != 53 {
		t.Fatal("named update count", len(names))
	}
	seen := make(map[unsafe.Pointer]bool)
	for _, name := range names {
		key, size := server.PortTestUnitGameplayRegistration(name, false)
		if key == nil || seen[key] {
			t.Fatal("update identity", name)
		}
		seen[key] = true
		typ.Update = key
		typ.UpdateDataSize = size
		obj.Update = typ.Update
		collisionRegistryGrow(128)
		got, n := server.PortTestUnitGameplayRegistration(name, false)
		if got != key || n != size || typ.Update != key || obj.Update != key || typ.UpdateDataSize != size {
			t.Fatal("stored update identity", name)
		}
	}
	for _, get := range []func() unsafe.Pointer{legacy.Get_nox_xxx_updatePlayerObserver_4E62F0, legacy.Get_nox_xxx___mkgmtime_538280, legacy.Get_nox_xxx_updatePlayerMonsterBot_4FAB20} {
		key := get()
		if key == nil || seen[key] {
			t.Fatal("special update identity")
		}
		seen[key] = true
		obj.Update = key
		collisionRegistryGrow(128)
		if get() != key || obj.Update != key {
			t.Fatal("special update lifetime")
		}
	}
	for _, row := range []struct {
		name string
		get  func() unsafe.Pointer
	}{
		{"PlayerUpdate", legacy.Get_nox_xxx_updatePlayer_4F8100},
		{"PixieUpdate", legacy.Get_nox_xxx_updatePixie_53CD20},
		{"HarpoonUpdate", legacy.Get_nox_xxx_updateHarpoon_54F380},
	} {
		key, _ := server.PortTestUnitGameplayRegistration(row.name, false)
		if row.get() != key {
			t.Fatal("named update getter", row.name)
		}
	}
}

func TestUpdateIdentitySpecialHooks(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	for _, row := range []struct {
		get  func() unsafe.Pointer
		slot *func(*server.Object)
	}{
		{legacy.Get_nox_xxx_updatePlayerObserver_4E62F0, &legacy.Nox_xxx_updatePlayerObserver_4E62F0},
		{legacy.Get_nox_xxx___mkgmtime_538280, &legacy.Nox_xxx___mkgmtime_538280},
	} {
		func() {
			old := *row.slot
			defer func() { *row.slot = old }()
			u.Update = row.get()
			for _, v := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				calls := 0
				*row.slot = func(got *server.Object) {
					if got != u {
						t.Fatal("special update argument")
					}
					calls++
					got.Field34 = v
				}
				collisionRegistryGrow(128)
				u.CallUpdate()
				if calls != 1 || u.Field34 != v {
					t.Fatal("special update hook replacement", calls, u.Field34, v)
				}
			}
		}()
	}
}

func TestUpdateIdentityBotOwner(t *testing.T) {
	for class := byte(0); class < 3; class++ {
		for _, frame := range []uint32{1, 30, 60, 120} {
			run := func(route uint8) []legacy.PortTestRoamResult {
				s := controlsBase(47)
				s.Owner.Frame = frame
				p := s.Callbacks.Shop
				p.Resources.PlayerClass = class
				sp := p.TemporaryUpdates.World.Objectives.Attack.Controls
				sp.Bot = true
				sp.BotUpdateRoute = route
				p.Sequence = []legacy.PortTestShopAction{{Op: 1444}, {Op: 1447}}
				return controlsRun(t, []legacy.PortTestRoamSpec{s})
			}
			direct, stored := run(1), run(2)
			if len(direct) == 0 || !reflect.DeepEqual(direct, stored) {
				t.Fatal("bot direct/stored callback state mismatch", class, frame)
			}
		}
	}
}
