//go:build porttest

package opennox

import (
	"fmt"
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func TestScriptBindingsMissingObjectPredicates(t *testing.T) {
	o := newWorldCollisionOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	for _, fi := range []asm.Builtin{asm.BuiltinUnknownb8, asm.BuiltinUnknownb9} {
		for _, root := range []bool{false, true} {
			for _, id := range []int32{0, -1, -2, 0x7fffffff} {
				t.Run(fmt.Sprintf("builtin%d/root%v/id%d", fi, root, id), func(t *testing.T) {
					if noxServer.noxScript.ScriptToObject(int(id)) != nil {
						t.Fatal("fixture unexpectedly resolved missing object")
					}
					const sentinel = 0x13579bdf
					o.s.NoxScriptVM.PushU32(sentinel)
					o.s.NoxScriptVM.PushU32(uint32(id))
					if root {
						if err := noxServer.noxScript.callBuiltinNative(fi); err != nil {
							t.Fatal(err)
						}
					} else if result, ok := legacy.CallScriptBuiltin(fi); !ok || result != 0 {
						t.Fatal("builtin dispatch failed", result, ok)
					}
					got := o.s.NoxScriptVM.PopU32()
					if o.s.NoxScriptVM.PopU32() != sentinel {
						t.Fatal("predicate changed surrounding stack")
					}
					if got != 0 {
						t.Error("missing-object predicate returned nonzero")
					}
				})
			}
		}
	}
}

func TestScriptBindingsObjectPredicates(t *testing.T) {
	o := newWorldCollisionOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	legacy.PortTestMonsterCache("reset", nil, 0)
	u := &o.units[0]
	saved := *u
	t.Cleanup(func() { *u = saved })
	u.ScriptIDVal = 100123
	u.ObjFlags = 0
	legacy.PortTestMonsterCache("prepare", u, 0)
	for _, fi := range []asm.Builtin{asm.BuiltinUnknownb8, asm.BuiltinUnknownb9} {
		for _, root := range []bool{false, true} {
			// Exhaust both tested bits and their neighbours, including high bits.
			for bits := uint32(0); bits < 1024; bits++ {
				for _, high := range []uint32{0, 0xfffffc00} {
					u.ObjSubClass = object.SubClass(bits | high)
					if noxServer.noxScript.ScriptToObject(u.ScriptIDVal) != u {
						t.Fatal("fixture did not resolve actual object")
					}
					const sentinel = 0x2468ace0
					o.s.NoxScriptVM.PushU32(sentinel)
					o.s.NoxScriptVM.PushU32(uint32(u.ScriptIDVal))
					if root {
						if err := noxServer.noxScript.callBuiltinNative(fi); err != nil {
							t.Fatal(err)
						}
					} else if result, ok := legacy.CallScriptBuiltin(fi); !ok || result != 0 {
						t.Fatal("builtin dispatch failed", result, ok)
					}
					mask := uint32(0x100)
					if fi == asm.BuiltinUnknownb9 {
						mask = 0x80
					}
					want := uint32(0)
					if bits&mask != 0 {
						want = 1
					}
					if got := o.s.NoxScriptVM.PopU32(); got != want {
						t.Fatalf("builtin %d root %v subclass %08x: got %d want %d", fi, root, bits|high, got, want)
					}
					if o.s.NoxScriptVM.PopU32() != sentinel || uint32(u.ObjSubClass) != bits|high {
						t.Fatal("predicate mutated surrounding stack or object")
					}
				}
			}
		}
	}
}

func TestScriptBindingsRoamByte(t *testing.T) {
	o := newWorldCollisionOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	legacy.PortTestMonsterCache("reset", nil, 0)
	u := &o.units[0]
	saved := *u
	ud, free := alloc.New(server.MonsterUpdateData{})
	t.Cleanup(free)
	t.Cleanup(func() { *u = saved })
	u.UpdateData = unsafe.Pointer(ud)
	u.ScriptIDVal = 100124
	legacy.PortTestMonsterCache("prepare", u, 0)
	for _, class := range []object.Class{object.ClassMonster, object.ClassSimple} {
		for _, flags := range []object.Flags{0, 0x20, 0x8000} {
			for _, root := range []bool{false, true} {
				for value := uint32(0); value < 512; value++ {
					u.ObjClass, u.ObjFlags = class, flags
					ud.Field333 = 0xaabbccdd
					const sentinel = 0x2468ace0
					o.s.NoxScriptVM.PushU32(sentinel)
					o.s.NoxScriptVM.PushU32(uint32(u.ScriptIDVal))
					o.s.NoxScriptVM.PushU32(value | 0x87650000)
					if root {
						if err := noxServer.noxScript.callBuiltinNative(asm.BuiltinSetRoamFlag); err != nil {
							t.Fatal(err)
						}
					} else if result, ok := legacy.CallScriptBuiltin(asm.BuiltinSetRoamFlag); !ok || result != 0 {
						t.Fatal("builtin dispatch failed", result, ok)
					}
					want := uint32(0xaabbccdd)
					if class == object.ClassMonster && flags&0x20 == 0 {
						want = 0xaabbcc00 | uint32(byte(value))
					}
					if ud.Field333 != want || o.s.NoxScriptVM.PopU32() != sentinel {
						t.Fatalf("roam byte or stack mismatch class %x flags %x root %v value %x", class, flags, root, value)
					}
				}
			}
		}
	}
}
