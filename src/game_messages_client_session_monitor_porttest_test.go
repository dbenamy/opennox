//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionMonitor(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	detach := func() {
		for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
			c.Objs.RemoveHealthBar(dr, 255)
		}
	}
	t.Cleanup(detach)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	allies := combatAllyStorage(t)
	type drawable struct {
		Code, Type, Class, Flags, Frame uint32
		Mask                            byte
		Position                        image.Point
	}
	type state struct {
		Drawables []drawable
		Minimap   []uint32
		Allies    []byte
	}
	type row struct {
		On, Kind, Exists, Marked, Type, Initial, Action, Mask int
		State                                                 state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{210, 219, 220} {
			for exists := 0; exists < 2; exists++ {
				for marked := 0; marked < 2; marked++ {
					for _, typ := range []int{0, 4} {
						for _, initial := range []byte{0, 1, 3} {
							actions := []int{0, 1, 2}
							masks := []int{0, 1, 2, 128, 255}
							if kind != 210 {
								actions = []int{1}
								masks = []int{1}
							}
							for _, action := range actions {
								for _, mask := range masks {
									code := uint16(17 | marked<<15)
									run := func(dispatch bool) state {
										detach()
										c.resetCase(env, pix, 1, 100)
										binary.LittleEndian.PutUint32(connected, uint32(on))
										legacy.PortTestCombatAllyClear()
										if initial == 1 {
											legacy.PortTestCombatAllyAdd(17, 123, 456)
										}
										legacy.PortTestCombatAllyAdd(0x8011, 234, 567)
										legacy.PortTestCombatAllyAdd(999, 345, 678)
										if initial == 3 {
											for id := 100; id < 130; id++ {
												legacy.PortTestCombatAllyAdd(uint32(id), uint16(id), uint16(id+1))
											}
										}
										neighbor := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(250, 250))
										neighbor.NetCode32 = 999
										neighbor.ObjClass = 0
										c.Objs.MinimapAdd(neighbor, 2)
										var dr *client.Drawable
										if exists != 0 {
											dr = c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
											dr.NetCode32 = 17
											dr.ObjClass = 0
											if marked != 0 {
												dr.ObjClass = object.Class(0x20400000)
											}
											if initial != 0 {
												c.Objs.MinimapAdd(dr, initial)
											}
										}
										if dispatch {
											data := binary.LittleEndian.AppendUint16([]byte{byte(kind)}, code)
											if kind != 220 {
												data = binary.LittleEndian.AppendUint16(data, uint16(typ))
											}
											if kind == 210 {
												data = append(data, byte(action), byte(mask))
											}
											input := bytes.Clone(data)
											n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
											if n != len(data) || !bytes.Equal(input, data) {
												t.Fatal("monitor input/length")
											}
										} else if on != 0 {
											switch kind {
											case 210:
												if action == 1 {
													if dr == nil {
														dr = c.Nox_xxx_spriteCreate_48E970(typ, 17, 0, 0)
													}
													if dr != nil {
														c.Objs.MinimapAdd(dr, byte(mask))
													}
												} else if dr != nil {
													c.Objs.RemoveHealthBar(dr, byte(mask))
												}
											case 219:
												if dr == nil {
													dr = c.Nox_xxx_spriteCreate_48E970(typ, 17, 0, 0)
												}
												if dr != nil {
													c.Objs.MinimapAdd(dr, 1)
												}
												legacy.PortTestCombatAllyAdd(17, 0, 0)
											case 220:
												legacy.PortTestCombatAllyRemove(uint32(code))
												if dr != nil {
													c.Objs.RemoveHealthBar(dr, 1)
												}
											}
										}
										r := state{Allies: bytes.Clone(allies)}
										for p := c.Objs.List1; p != nil; p = p.NextPtr {
											r.Drawables = append(r.Drawables, drawable{p.NetCode32, p.TypeIDVal, uint32(p.ObjClass), uint32(p.ObjFlags), p.Field_80, p.Field_71_0, p.PosVec})
										}
										for p := c.Objs.FirstMinimapList(); p != nil; p = p.Nox_xxx_cliNextMinimapObj_459EC0(p) {
											if len(r.Minimap) > 3 {
												t.Fatal("minimap cycle")
											}
											r.Minimap = append(r.Minimap, p.NetCode32)
										}
										if neighbor.Field_71_0 != 2 || neighbor.NetCode32 != 999 || neighbor.PosVec != image.Pt(250, 250) {
											t.Fatal("neighbor changed")
										}
										return r
									}
									want := run(false)
									got := run(true)
									if !reflect.DeepEqual(got, want) {
										t.Fatal("monitor state", on, kind, exists, marked, typ, initial, action, mask, got, want)
									}
									rows = append(rows, row{on, kind, exists, marked, typ, int(initial), action, mask, got})
								}
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-monitor", rows)
}
