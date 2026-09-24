//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestLifecycleIdentityLifetime(t *testing.T) {
	typ, freeType := alloc.New(server.ObjectType{})
	defer freeType()
	obj, freeObject := alloc.New(server.Object{})
	defer freeObject()
	old := legacy.Nox_xxx_unitMonsterInit_4F0040
	defer func() { legacy.Nox_xxx_unitMonsterInit_4F0040 = old }()
	calls := 0
	legacy.Nox_xxx_unitMonsterInit_4F0040 = func(got *server.Object) {
		if got != obj {
			t.Fatal("stored init key object identity")
		}
		calls++
	}
	seen := make(map[unsafe.Pointer]string)
	total := 0
	for _, create := range []bool{true, false} {
		names := legacy.PortTestLifecycleRegistryNames(create)
		wantCount := 12
		if create {
			wantCount = 8
		}
		if len(names) != wantCount {
			t.Fatal("lifecycle name count", create, len(names))
		}
		for _, name := range names {
			total++
			key, size := server.PortTestLifecycleRegistry(name, create)
			if key == nil {
				t.Fatal("nil lifecycle key", name)
			}
			if previous, ok := seen[key]; ok {
				if create || name != "ShopkeeperInit" || previous != "MonsterInit" {
					t.Fatal("unexpected lifecycle alias", previous, name)
				}
			} else {
				seen[key] = name
			}
			wantSize := uintptr(0)
			switch name {
			case "ShopkeeperInit":
				wantSize = unsafe.Sizeof(server.ShopkeeperInitData{})
			case "GoldInit":
				wantSize = unsafe.Sizeof(server.GoldInitData{})
			case "SkullInit", "DirectionInit":
				wantSize = 8
			}
			if size != wantSize {
				t.Fatal("lifecycle data size", name, size, wantSize)
			}
			if create {
				typ.Create = key
			} else {
				typ.Init = key
				obj.Init = typ.Init
			}
			collisionRegistryGrow(128)
			got, _ := server.PortTestLifecycleRegistry(name, create)
			if got != key {
				t.Fatal("lifecycle key changed after GC", name)
			}
			if create {
				if typ.Create != key {
					t.Fatal("stored create key changed", name)
				}
			} else {
				if typ.Init != key || obj.Init != key {
					t.Fatal("stored init key changed", name)
				}
				if name == "MonsterInit" || name == "ShopkeeperInit" {
					obj.CallInit()
				}
			}
		}
	}
	runtime.GC()
	if total != 20 || len(seen) != 19 || calls != 2 {
		t.Fatal("lifecycle identity/dispatch counts", total, len(seen), calls)
	}
}

// Gold methods classify by the registered init identity, independent of amount.
func TestLifecycleGoldIdentityConsumers(t *testing.T) {
	u, freeObject := alloc.New(server.Object{})
	defer freeObject()
	data, freeData := alloc.New(server.GoldInitData{})
	defer freeData()
	u.InitData = unsafe.Pointer(data)
	obj := (*Object)(u)
	for _, name := range legacy.PortTestLifecycleRegistryNames(false) {
		u.Init, _ = server.PortTestLifecycleRegistry(name, false)
		data.Amount = 123
		want := 0
		if name == "GoldInit" {
			want = 123
		}
		if got := obj.GetGold(); got != want {
			t.Fatal("gold classifier", name, got, want)
		}
	}
	u.Init, _ = server.PortTestLifecycleRegistry("GoldInit", false)
	for _, amount := range []int{0, 1, -1, 2147483647, -2147483648} {
		obj.SetGold(amount)
		if data.Amount != uint32(amount) || obj.GetGold() != amount {
			t.Fatal("gold set/get", amount, data.Amount, obj.GetGold())
		}
		for _, delta := range []int{1, -2, 2147483647, -2147483648} {
			before := data.Amount
			obj.ChangeGold(delta)
			want := before + uint32(delta)
			if data.Amount != want || obj.GetGold() != int(want) {
				t.Fatal("gold change", before, delta, data.Amount, want)
			}
		}
	}
	var absent *Object
	absent.SetGold(7)
	absent.ChangeGold(-3)
	if absent.GetGold() != 0 {
		t.Fatal("nil gold")
	}
}
