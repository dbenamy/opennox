//go:build porttest

package opennox

import (
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestScriptBindingsMonsterMovement(t *testing.T) {
	o := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	if memmap.Float32(0x587000, 194136) != 1 {
		t.Fatal("missing shipped direction table")
	}
	u := &o.units[0]
	saved := *u
	ud, free := alloc.New(server.MonsterUpdateData{})
	t.Cleanup(free)
	wp, free := alloc.New(server.Waypoint{})
	t.Cleanup(free)
	t.Cleanup(func() { *u = saved })
	oldMover := legacy.Get_dword_5d4594_2386836()
	legacy.Set_dword_5d4594_2386836(12345)
	t.Cleanup(func() { legacy.Set_dword_5d4594_2386836(oldMover) })
	oldChanged := o.s.AI.StackChanged
	t.Cleanup(func() { o.s.AI.StackChanged = oldChanged })
	type row struct {
		Op                      string
		Class, Flags            uint32
		Direction, Stack, Links int
		Index                   int8
		Actions                 [24]server.AIStackItem
		Changed                 bool
	}
	var rows []row
	for _, op := range []string{"roam", "home", "move"} {
		for _, class := range []object.Class{object.ClassMonster, object.ClassSimple} {
			for _, flags := range []object.Flags{0, 0x20, 0x8000} {
				for _, direction := range []int{0, 1, 32, 63, 64, 127, 128, 191, 192, 255} {
					for _, stack := range []int{0, 3, 23} {
						for _, links := range []int{0, 1, 32} {
							*u = saved
							u.ObjClass, u.ObjFlags, u.TypeInd = class, flags, 123
							u.UpdateData = unsafe.Pointer(ud)
							u.PosVec = types.Ptf(100.125, -200.5)
							*ud = server.MonsterUpdateData{Direction94: uint32(direction), Pos95: types.Ptf(-31.25, 73.5), Field333: 0xaabbccdd, AIStackInd: int8(stack)}
							for i := 0; i <= stack; i++ {
								ud.AIStack[i] = server.AIStackItem{Action: 1, Args: [4]uintptr{11, 22, 33, 44}}
							}
							*wp = server.Waypoint{Index: 789, PosVec: types.Ptf(17.25, -18.5), PointsCnt: byte(links)}
							before := ud.AIStack
							o.s.AI.StackChanged = false
							switch op {
							case "roam":
								legacy.Nox_xxx_scriptMonsterRoam_512930(u)
							case "home":
								legacy.Nox_server_gotoHome(u)
							case "move":
								legacy.Nox_server_scriptMoveTo_5123C0(u, wp)
							}
							if class != object.ClassMonster || flags&0x8000 != 0 {
								if ud.AIStack != before || int(ud.AIStackInd) != stack || o.s.AI.StackChanged {
									t.Fatal("ineligible command changed stack")
								}
							} else {
								want := []uint32{32, 10}
								arg := uintptr(10)
								if op == "home" {
									want = []uint32{32, 25, 37}
									arg = 37
								}
								if op == "move" {
									want = []uint32{32, 8}
									arg = 8
									if links != 0 {
										want = []uint32{32, 10, 8}
									}
								}
								if int(ud.AIStackInd) != len(want)-1 || !o.s.AI.StackChanged {
									t.Fatalf("%s active stack size/notification", op)
								}
								for i, v := range want {
									if ud.AIStack[i].Action != v {
										t.Fatalf("%s action %d", op, i)
									}
								}
								if ud.AIStack[0].Args[0] != arg {
									t.Fatal("report action argument")
								}
								if op == "roam" && (ud.AIStack[1].Args[0] != 0 || ud.AIStack[1].Args[2] != 0xdd) {
									t.Fatal("roam flags/target")
								}
								if op == "home" {
									x := float32(float64(memmap.Float32(0x587000, 194136+8*uintptr(direction)))*10 + float64(u.PosVec.X))
									y := float32(float64(memmap.Float32(0x587000, 194140+8*uintptr(direction)))*10 + float64(u.PosVec.Y))
									if ud.AIStack[1].Args[0] != uintptr(math.Float32bits(x)) || ud.AIStack[1].Args[1] != uintptr(math.Float32bits(y)) || ud.AIStack[2].Args[0] != uintptr(math.Float32bits(ud.Pos95.X)) || ud.AIStack[2].Args[1] != uintptr(math.Float32bits(ud.Pos95.Y)) || ud.AIStack[2].Args[2] != 0 {
										t.Fatal("home positions")
									}
								}
								if op == "move" {
									last := &ud.AIStack[ud.AIStackInd]
									if last.Args[0] != uintptr(math.Float32bits(wp.PosVec.X)) || last.Args[1] != uintptr(math.Float32bits(wp.PosVec.Y)) || last.Args[2] != 0 {
										t.Fatal("move waypoint position")
									}
									if links != 0 && (ud.AIStack[1].Args[0] != uintptr(unsafe.Pointer(wp)) || ud.AIStack[1].Args[2] != 0xdd) {
										t.Fatal("move waypoint ownership/flags")
									}
								}
							}
							actions := ud.AIStack
							if op == "move" && class == object.ClassMonster && flags&0x8000 == 0 && links != 0 {
								actions[1].Args[0] = 9001
							}
							rows = append(rows, row{op, uint32(class), uint32(flags), direction, stack, links, ud.AIStackInd, actions, o.s.AI.StackChanged})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-monster-movement", rows, "e79e5828721bff8ab0d0c8e8e7aa7d084d5c01b5a32d3b8b8885a97a2499b993")
}
