//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestModifierItemUpdateOwner(t *testing.T) {
	actor, freeActor := alloc.New(server.Object{})
	defer freeActor()
	items, freeItems := alloc.Make([]server.Object{}, 3)
	defer freeItems()
	records, freeRecords := alloc.Make([]server.ModifierInitData{}, 3)
	defer freeRecords()
	mods, freeMods := alloc.Make([]server.ModifierEff{}, 4)
	defer freeMods()
	core := new(server.Server)
	actor.InvFirstItem = &items[0]
	for i := range items {
		items[i].InitData = unsafe.Pointer(&records[i])
		if i+1 < len(items) {
			items[i].InvNextItem = &items[i+1]
		}
	}
	// These later list entries verify both class and equipped filtering.
	items[1].ObjClass = object.ClassWeapon
	items[2].ObjClass = object.ClassPlayer
	items[2].ObjFlags = object.FlagEquipped
	for _, class := range []object.Class{0, object.ClassWeapon, object.ClassArmor, object.ClassWand, object.ClassFlag, object.ClassPlayer, object.ClassPlayer | object.ClassWeapon, object.ClassMonster} {
		for _, equipped := range []bool{false, true} {
			items[0].ObjClass = class
			items[0].ObjFlags = 0
			if equipped {
				items[0].ObjFlags = object.FlagEquipped
			}
			for combination := 0; combination < 81; combination++ {
				keys := legacy.PortTestModifierObserverReset(0)
				var last *server.ModifierEff
				calls := 0
				code := combination
				for i := range mods {
					mods[i] = server.ModifierEff{Price20: int32(i + 1)}
					records[0].Modifiers[i] = nil
					state := code % 3
					code /= 3
					if state > 0 {
						records[0].Modifiers[i] = &mods[i]
					}
					if state == 2 {
						mods[i].Update100.Fnc = keys[1]
						calls++
						last = &mods[i]
					}
					records[1].Modifiers[i] = &mods[i]
					records[2].Modifiers[i] = &mods[i]
				}
				allowed := object.ClassWeapon | object.ClassArmor | object.ClassWand | object.ClassFlag
				if !equipped || !class.HasAny(allowed) {
					calls = 0
					last = nil
				}
				beforeItems := [3]server.Object{items[0], items[1], items[2]}
				beforeMods := [4]server.ModifierEff{mods[0], mods[1], mods[2], mods[3]}
				core.ItemsApplyUpdateEffect(actor)
				args, count, kind := legacy.PortTestModifierObserverSnapshot()
				var expected [6]uintptr
				wantKind := 0
				if calls > 0 {
					expected[0] = uintptr(unsafe.Pointer(last))
					expected[1] = uintptr(unsafe.Pointer(actor))
					wantKind = 2
				}
				if count != calls || kind != wantKind || args != expected {
					t.Fatalf("class%x equipped%v combination%d: %d/%d/%v want %d/%d/%v", class, equipped, combination, count, kind, args, calls, wantKind, expected)
				}
				for i := range items {
					if items[i] != beforeItems[i] {
						t.Fatal("item changed", i)
					}
				}
				for i := range mods {
					if mods[i] != beforeMods[i] {
						t.Fatal("modifier changed", i)
					}
				}
			}
		}
	}
}
