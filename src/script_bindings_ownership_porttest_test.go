//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestScriptBindingsOwnership(t *testing.T) {
	o := newWorldCollisionOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	legacy.PortTestMonsterCache("reset", nil, 0)
	target, other, host := &o.units[0], &o.units[1], &o.units[2]
	saved := append([]server.Object(nil), o.units...)
	pl := o.s.Players.ByInd(31)
	savedPlayer := *pl
	hp, free := alloc.New(server.HealthData{Cur: 73, Max: 109})
	t.Cleanup(free)
	t.Cleanup(func() {
		o.s.Players.Nox_xxx_netUnmarkMinimapObj_417300(31, target, 1)
		copy(o.units, saved)
		*pl = savedPlayer
	})
	type row struct {
		Builtin             asm.Builtin
		Root, Missing, Host bool
		Owner               int
		Before, After       uint32
		ResultOwner         int
		Marker              bool
		Queue               legacy.PortTestReliableReportState
	}
	var rows []row
	for _, fi := range []asm.Builtin{asm.BuiltinMakeFriendly, asm.BuiltinMakeEnemy, asm.BuiltinBecomePet, asm.BuiltinBecomeEnemy} {
		for _, root := range []bool{false, true} {
			for _, missing := range []bool{false, true} {
				for _, present := range []bool{false, true} {
					for _, owner := range []int{0, 1, 2} {
						for _, bits := range []uint32{0, 0x80, 0x100, 0x180, 0xff030180} {
							o.s.Players.Nox_xxx_netUnmarkMinimapObj_417300(31, target, 1)
							copy(o.units, saved)
							*pl = savedPlayer
							o.reset()
							target.ScriptIDVal = 100123
							target.ObjFlags = 0
							target.ObjClass = object.ClassSimple
							target.ObjSubClass = object.SubClass(bits)
							target.HealthData = hp
							other.ObjClass = object.ClassSimple
							target.ObjOwner = nil
							target.Field128 = nil
							host.Field129 = nil
							other.Field129 = nil
							if owner == 1 {
								target.ObjOwner = other
								other.Field129 = target
							}
							if owner == 2 {
								target.ObjOwner = host
								host.Field129 = target
							}
							pl.PlayerUnit = nil
							if present {
								pl.PlayerUnit = host
							}
							legacy.PortTestMonsterCache("reset", nil, 0)
							legacy.PortTestMonsterCache("prepare", target, 0)
							id := uint32(target.ScriptIDVal)
							if missing {
								id = 999
							}
							const sentinel = 0x2468ace0
							o.s.NoxScriptVM.PushU32(sentinel)
							o.s.NoxScriptVM.PushU32(id)
							if root {
								if err := noxServer.noxScript.callBuiltinNative(fi); err != nil {
									t.Fatal(err)
								}
							} else if r, ok := legacy.CallScriptBuiltin(fi); r != 0 || !ok {
								t.Fatal("ownership builtin dispatch", r, ok)
							}
							if o.s.NoxScriptVM.PopU32() != sentinel {
								t.Fatal("ownership builtin changed surrounding stack")
							}
							wantBits, wantOwner, wantMarker := bits, owner, false
							if !missing {
								switch fi {
								case asm.BuiltinMakeFriendly:
									wantBits |= 0x100
									if present {
										wantOwner = 2
									}
								case asm.BuiltinMakeEnemy:
									wantBits &^= 0x100
									wantOwner = 0
								case asm.BuiltinBecomePet:
									if present {
										wantBits |= 0x80
										wantOwner = 2
										wantMarker = true
									}
								case asm.BuiltinBecomeEnemy:
									if present {
										wantBits &^= 0x80
										wantOwner = 0
									}
								}
							}
							gotOwner := 0
							if target.ObjOwner == other {
								gotOwner = 1
							} else if target.ObjOwner == host {
								gotOwner = 2
							} else if target.ObjOwner != nil {
								t.Fatal("unowned parent pointer")
							}
							marked := pl.Field4580 != nil
							if uint32(target.ObjSubClass) != wantBits || gotOwner != wantOwner || marked != wantMarker {
								t.Fatalf("ownership builtin%d root%v missing%v host%v owner%d bits%x: got%x/%d/%v want%x/%d/%v", fi, root, missing, present, owner, bits, target.ObjSubClass, gotOwner, marked, wantBits, wantOwner, wantMarker)
							}
							if (host.Field129 == target) != (wantOwner == 2) || (other.Field129 == target) != (wantOwner == 1) || target.Field128 != nil {
								t.Fatal("ownership list links")
							}
							if marked && (pl.Field4580.Field4 != target || pl.Field4580.Field0 != 1 || pl.Field4580.Field8 != pl.Field4580 || pl.Field4580.Field12 != pl.Field4580) {
								t.Fatal("minimap ownership")
							}
							rows = append(rows, row{fi, root, missing, present, owner, bits, uint32(target.ObjSubClass), gotOwner, marked, o.state()})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-ownership", rows, "b8c5d288106d8e981936055e18e3962855b7e4100fc31b156081aa8889f424f4")
}
