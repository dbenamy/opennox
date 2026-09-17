//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

type mapPolygonScriptBase = legacy.NoxScript

type mapPolygonScriptTrace struct {
	mapPolygonScriptBase
	events []int
	calls  []int
}

func (s *mapPolygonScriptTrace) ScriptCallback(b *server.ScriptCallback, caller, trigger *server.Object, event server.ScriptEventType) unsafe.Pointer {
	s.events = append(s.events, int(event))
	return s.mapPolygonScriptBase.ScriptCallback(b, caller, trigger, event)
}

type mapPolygonServerTrace struct {
	legacy.Server
	script *mapPolygonScriptTrace
}

func (s *mapPolygonServerTrace) NoxScriptC() legacy.NoxScript { return s.script }

func TestMapPolygonsActorTransitions(t *testing.T) {
	world := newWorldCollisionOwner(t)
	o := newMapPolygonsOwner(t)
	for _, x := range []float32{100, 130} {
		o.construct(t, [][2]float32{{x, 200}, {x + 10, 200}, {x + 10, 210}, {x, 210}})
	}
	original := legacy.GetServer
	trace := &mapPolygonScriptTrace{mapPolygonScriptBase: original().NoxScriptC()}
	proxy := &mapPolygonServerTrace{Server: original(), script: trace}
	legacy.GetServer = func() legacy.Server { return proxy }
	t.Cleanup(func() { legacy.GetServer = original })
	var callbacks [2][2]int32
	for i := range callbacks {
		for j := range callbacks[i] {
			label := 10*(i+1) + j
			callbacks[i][j] = int32(world.s.NoxScriptVM.AsFuncIndex(fmt.Sprintf("polygon-%d", label), func() {
				if world.s.NoxScriptVM.Caller() != &world.units[0] || world.s.NoxScriptVM.Trigger() != nil {
					t.Error("polygon script caller/trigger")
				}
				trace.calls = append(trace.calls, label)
			}))
		}
	}
	u := &world.units[0]
	player := u.UpdateDataPlayer().Player
	playerBytes := unsafe.Slice((*byte)(unsafe.Pointer(player)), int(unsafe.Sizeof(*player)))
	savedUpdate := u.UpdateData
	monsterData := world.record(t, 16)
	t.Cleanup(func() { u.UpdateData = savedUpdate })
	type row struct {
		Name          string
		Cache         uint32
		Level         byte
		Visited       [2]uint32
		Events, Calls []int
	}
	var rows []row
	for _, actor := range []string{"player", "monster"} {
		for _, old := range []uint32{0, 1, 2, 0xdeadface} {
			for _, dest := range []int{0, 1, 2} {
				for _, moving := range []bool{false, true} {
					for _, host := range []bool{false, true} {
						if actor == "monster" && host {
							continue
						}
						name := fmt.Sprintf("%s/old%x/dest%d/moving%t/host%t", actor, old, dest, moving, host)
						t.Run(name, func(t *testing.T) {
							trace.events = nil
							trace.calls = nil
							for i := 1; i <= 2; i++ {
								p := o.polygon(i)
								for j := 0; j < 2; j++ {
									cb := (*server.ScriptCallback)(unsafe.Pointer(&p[112+8*j]))
									cb.Flags = 0
									cb.Func = callbacks[i-1][j]
								}
								p[130] = byte(i + 4)
								binary.LittleEndian.PutUint32(p[136:], 0)
							}
							binary.LittleEndian.PutUint32(o.control[8:], 1)
							u.PosVec = types.Pointf{500, 400}
							if dest != 0 {
								u.PosVec = types.Pointf{float32(102 + 30*(dest-1)), 203}
							}
							u.PrevPos = u.PosVec
							if moving {
								u.PrevPos.X--
							}
							cache := (*uint32)(unsafe.Pointer(&playerBytes[3664]))
							*cache = old
							playerBytes[3668] = 77
							player.PlayerInd = 1
							binary.LittleEndian.PutUint32(playerBytes[3660:], uint32(dest))
							u.ObjClass = object.ClassPlayer
							u.UpdateData = savedUpdate
							if host {
								player.PlayerInd = 31
							}
							if actor == "monster" {
								u.ObjClass = object.ClassMonster
								u.UpdateData = monsterData
								cache = (*uint32)(monsterData)
								*cache = old
							}
							legacy.PortTestMapPolygonActor(actor, u.CObj())
							want := old
							var events, calls []int
							if moving || old == 0xdeadface || host {
								if host && dest == 0 {
									want = 0
								} else {
									target := uint32(dest)
									// Existing player lookup retains its old region when the point leaves
									// all polygons. The host's zero client cache takes a separate path.
									if actor == "player" && !host && dest == 0 && old != 0 && old != 0xdeadface {
										target = old
									}
									if target != 0 {
										want = target
										if old != target && old != 0xdeadface {
											if old != 0 {
												ev := 29
												if actor == "monster" {
													ev = 26
												}
												events = append(events, ev)
												calls = append(calls, int(old)*10+1)
											}
											ev := 28
											if actor == "monster" {
												ev = 25
											}
											events = append(events, ev)
											calls = append(calls, int(target)*10)
										}
									} else if old != 0 && old != 0xdeadface {
										want = 0
										ev := 27
										if actor == "monster" {
											ev = 24
										}
										events = append(events, ev)
										calls = append(calls, int(old)*10+1)
									}
								}
							}
							if *cache != want || !reflect.DeepEqual(trace.events, events) || !reflect.DeepEqual(trace.calls, calls) {
								t.Fatalf("cache %x want %x; events %v want %v; calls %v want %v", *cache, want, trace.events, events, trace.calls, calls)
							}
							rows = append(rows, row{name, *cache, playerBytes[3668], [2]uint32{binary.LittleEndian.Uint32(o.polygon(1)[136:]), binary.LittleEndian.Uint32(o.polygon(2)[136:])}, append([]int(nil), trace.events...), append([]int(nil), trace.calls...)})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "map-polygons-actor-transitions", rows, "91eeec1f6c35c55eba8be10ed8d93c2d68e7c1663ca29686c171a6ff70208838")
}
