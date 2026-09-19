//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestScriptBindingsGroupRoam(t *testing.T) {
	o := newWorldCollisionOwner(t)
	data, free := alloc.Make([]server.MonsterUpdateData{}, 3)
	t.Cleanup(free)
	saved := append([]server.Object(nil), o.units...)
	oldList, oldGroups := o.s.Objs.List, o.s.MapGroups
	o.s.MapGroups = server.ServerMapGroups{}
	o.s.MapGroups.Init()
	t.Cleanup(func() { o.s.MapGroups.Free(); o.s.MapGroups = oldGroups; o.s.Objs.List = oldList; copy(o.units, saved) })
	for i := range o.units {
		u := &o.units[i]
		u.UpdateData = unsafe.Pointer(&data[i])
		u.Extent = uint32(101 + i)
		u.ObjClass = object.ClassMonster
		u.ObjNext = nil
		if i+1 < len(o.units) {
			u.ObjNext = &o.units[i+1]
		}
	}
	o.units[1].ObjClass = object.ClassSimple
	o.s.Objs.List = &o.units[0]
	group := func(id uint32, kind byte, ids ...uint32) {
		if o.s.MapGroups.MapLoadAddGroup57C0C0(fmt.Sprint(id), id, kind) != 1 {
			t.Fatal("group allocation")
		}
		for i := len(ids) - 1; i >= 0; i-- {
			if o.s.MapGroups.Sub57C130([]uint32{ids[i], 0}, id) != 1 {
				t.Fatal("group member")
			}
		}
	}
	group(10, 0, 101, 999, 102, 103, 101)
	group(11, 0, 103)
	group(12, 0)
	group(20, 1, 101)
	group(30, 3, 10, 20, 999, 11)
	group(31, 3, 30, 12)
	type row struct {
		Group uint32
		Flags object.Flags
		Root  bool
		Value uint32
		Words [3]uint32
	}
	var rows []row
	for _, id := range []uint32{0, 10, 11, 12, 20, 30, 31, 999} {
		for _, flags := range []object.Flags{0, 0x20, 0x8000} {
			for _, root := range []bool{false, true} {
				for _, value := range []uint32{0, 1, 127, 128, 255, 256, 511, 0x87654321, 0xffffffff} {
					for i := range data {
						data[i] = server.MonsterUpdateData{Field333: 0xaabbccdd}
						o.units[i].ObjFlags = flags
					}
					const sentinel = 0x2468ace0
					o.s.NoxScriptVM.PushU32(sentinel)
					o.s.NoxScriptVM.PushU32(id)
					o.s.NoxScriptVM.PushU32(value)
					if root {
						if err := noxServer.noxScript.callBuiltinNative(asm.BuiltinGroupSetRoamFlag); err != nil {
							t.Fatal(err)
						}
					} else if r, ok := legacy.CallScriptBuiltin(asm.BuiltinGroupSetRoamFlag); r != 0 || !ok {
						t.Fatal("group builtin dispatch", r, ok)
					}
					if o.s.NoxScriptVM.PopU32() != sentinel {
						t.Fatal("group changed surrounding stack")
					}
					var words [3]uint32
					for i := range data {
						want := server.MonsterUpdateData{Field333: 0xaabbccdd}
						selected := id == 10 || id == 30 || id == 31 || id == 11 && i == 2
						if selected && i != 1 && flags&0x20 == 0 {
							want.Field333 = 0xaabbcc00 | uint32(byte(value))
						}
						if data[i] != want {
							t.Fatalf("group %d flags%x root%v value%x member%d changed incorrectly", id, flags, root, value, i)
						}
						words[i] = data[i].Field333
					}
					rows = append(rows, row{id, flags, root, value, words})
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-group-roam", rows, "873acf1bc83385691c40a956f7080978423b455e679dcbda298b537541e2024c")
}
