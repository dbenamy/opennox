//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestWorldCollisionsScriptTriggers(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[0]
	data := o.record(t, 128)
	a.UpdateData = data
	a.CollideData = data
	t.Cleanup(func() { a.UpdateData = nil; a.CollideData = nil })
	var calls [][2]uint32
	index := o.s.NoxScriptVM.AsFuncIndex("world-collision-test", func() {
		calls = append(calls, [2]uint32{o.s.NoxScriptVM.Caller().NetCode, o.s.NoxScriptVM.Trigger().NetCode})
	})
	a.NetCode = 2001
	var rows []struct {
		Name                   string
		Calls                  [][2]uint32
		Flags, Armed, Deadline uint32
	}
	defer func() {
		spellbookCapture(t, "world-collisions-scripts", rows, "cf78abb4d2515a4832c6e2dde7916025f58b1c63fc000a1f776b58e6cfdecf71")
	}()
	for _, op := range []int{14, 18} {
		for _, flags := range []int32{0, 1, 2} {
			for _, armed := range []uint32{0, 1} {
				for _, player := range []bool{false, true} {
					for _, delay := range []uint16{0, 1, 65535} {
						for _, registered := range []bool{false, true} {
							name := fmt.Sprintf("op%d/flags%d/armed%d/player%t/delay%d/registered%t", op, flags, armed, player, delay, registered)
							t.Run(name, func(t *testing.T) {
								clear(unsafe.Slice((*byte)(data), 128))
								calls = nil
								o.s.SetFrame(0xfffffffe)
								a.ObjFlags = 0
								b.ObjClass = object.ClassSimple
								if player {
									b.ObjClass = object.ClassPlayer
								}
								offset := uintptr(0)
								if op == 18 {
									offset = 72
								}
								cb := (*server.ScriptCallback)(unsafe.Add(data, offset))
								cb.Func = int32(index)
								cb.Flags = uint32(flags)
								*(*uint16)(unsafe.Add(data, 20)) = delay
								objectXferSetWord(data, 24, armed)
								if registered {
									name := "TrapDoorCollide"
									if op == 18 {
										name = "MonsterGeneratorCollide"
									}
									a.Collide, _ = server.PortTestWorldCollisionRegistry(name)
									a.CallCollide(int(uintptr(b.CObj())), 0)
								} else {
									legacy.PortTestWorldCollision(op, a, b, nil)
								}
								admitted := op == 14 && armed == 0 || op == 18 && player
								fired := admitted && flags&1 == 0
								if (len(calls) == 1) != fired || len(calls) > 1 {
									t.Fatalf("calls %v", calls)
								}
								if fired && calls[0] != [2]uint32{b.NetCode, a.NetCode} {
									t.Fatal("script caller/trigger")
								}
								wantArmed := armed
								wantDeadline := uint32(0)
								if op == 14 && admitted {
									wantArmed = 1
									if delay != 0 {
										wantDeadline = uint32(0xfffffffe) + uint32(delay)
									}
								}
								if objectXferGetWord(data, 24) != wantArmed || objectXferGetWord(data, 16) != wantDeadline {
									t.Fatal("trigger arm/deadline")
								}
								rows = append(rows, struct {
									Name                   string
									Calls                  [][2]uint32
									Flags, Armed, Deadline uint32
								}{name, append([][2]uint32(nil), calls...), uint32(cb.Flags), objectXferGetWord(data, 24), objectXferGetWord(data, 16)})
							})
						}
					}
				}
			}
		}
	}
}

func TestWorldCollisionsTrapAbility(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	a.CollideData = o.record(t, 28)
	t.Cleanup(func() { a.CollideData = nil })
	calls := 0
	index := o.s.NoxScriptVM.AsFuncIndex("world-trap-ability", func() { calls++ })
	var rows [][3]uint32
	for _, ability := range []server.Ability{1, 2, 3, 4, 5} {
		for _, active := range []uint32{0, 1} {
			o.s.Abils.Reset()
			o.s.Abils.GetFor(b).ExecList = &server.ExecAbilityClass{Abil: ability, Active: active}
			clear(unsafe.Slice((*byte)(a.CollideData), 28))
			(*server.ScriptCallback)(a.CollideData).Func = int32(index)
			calls = 0
			legacy.PortTestWorldCollision(14, a, b, nil)
			want := 1
			if ability == 4 {
				want = 0
			}
			if calls != want || objectXferGetWord(a.CollideData, 24) != uint32(want) {
				t.Fatal("trap ability admission")
			}
			rows = append(rows, [3]uint32{uint32(ability), active, uint32(calls)})
		}
	}
	spellbookCapture(t, "world-collisions-trap-ability", rows, "0834d1620b726d21d23aae65dc0f2fed38d568b4c1f59d0dcbc73c858d6b6a5f")
}
