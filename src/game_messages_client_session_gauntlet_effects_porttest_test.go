//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletGreenBolt(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "GreenZap")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200904, 4)
	type row struct {
		On, Mode, Return int
		Frame            uint32
		Coords           [4]uint16
		Value            uint16
		Cache            uint32
		Drawables        [][]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for _, frame := range []uint32{0, 100, 0xffffffff} {
				for _, xy := range [][4]uint16{{0, 0, 0, 0}, {100, 200, 101, 199}, {0, 65535, 32768, 32767}} {
					for _, value := range []uint16{0, 1, 32768, 65535} {
						c.resetCase(env, pix, 1, frame)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						clear(cache)
						id := c.Things.IndByID("GreenZap")
						if mode == 2 {
							id = 4
						}
						if mode != 0 {
							binary.LittleEndian.PutUint32(cache, uint32(id))
						}
						data := []byte{240, 16}
						for _, v := range xy {
							data = binary.LittleEndian.AppendUint16(data, v)
						}
						data = binary.LittleEndian.AppendUint16(data, value)
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
						count := 0
						if on != 0 {
							count = 1
						}
						if n != 12 || !bytes.Equal(input, data) || binary.LittleEndian.Uint32(cache) != uint32(id) || len(c.Calls) != count || c.Objs.Count != count {
							t.Fatal("quest green bolt gate", on, mode, frame, xy, value, n)
						}
						if on != 0 {
							dr := c.Objs.List1
							b := unsafe.Slice((*byte)(unsafe.Add(dr.C(), 432)), 13)
							if c.Calls[0].Position != image.Pt(int(xy[2]), int(xy[3])) || c.Calls[0].Type != id || b[0] != 0 || binary.LittleEndian.Uint32(b[1:]) != uint32(value) || !bytes.Equal(b[5:], data[2:10]) {
								t.Fatal("quest green bolt fields")
							}
						}
						rows = append(rows, row{on, mode, n, frame, xy, value, binary.LittleEndian.Uint32(cache), c.snapshotDrawables(t)})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-green-bolt", rows)
}
func TestGameMessageClientSessionGauntletParticles(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "FireBoom")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200908, 8)
	pointType := serverConfigOwnBytes(t, 0x5D4594, 1200788, 4)
	type state struct {
		Cache     []byte
		Calls     []effectsSpawnCall
		Drawables [][]uint32
		RNG       int
	}
	snapshot := func() state {
		return state{bytes.Clone(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), c.srv.Rand.Other.Index()}
	}
	type row struct {
		On, Kind, Mode, Fail, Amount int
		Pos                          image.Point
		Return                       int
		State                        state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{25, 26} {
			for mode := 0; mode < 3; mode++ {
				for _, fail := range []int{0, 1, 3} {
					for _, amount := range []byte{0, 1, 128, 255} {
						for _, pos := range []image.Point{{0, 0}, {100, 200}, {-32768, 32767}} {
							setup := func() {
								c.resetCase(env, pix, 1, 100)
								c.FailEvery = fail
								binary.LittleEndian.PutUint32(connected, uint32(on))
								binary.LittleEndian.PutUint32(pointType, 4)
								clear(cache)
								if mode != 0 {
									spark, boom := c.Things.IndByID("GreenSpark"), c.Things.IndByID("FireBoom")
									if mode == 2 {
										spark, boom = 4, 4
									}
									binary.LittleEndian.PutUint32(cache, uint32(spark))
									binary.LittleEndian.PutUint32(cache[4:], uint32(boom))
								}
							}
							setup()
							if kind == 25 {
								if binary.LittleEndian.Uint32(cache) == 0 {
									binary.LittleEndian.PutUint32(cache, uint32(c.Things.IndByID("GreenSpark")))
									binary.LittleEndian.PutUint32(cache[4:], uint32(c.Things.IndByID("FireBoom")))
								}
								dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(int(binary.LittleEndian.Uint32(cache[4:])), pos)
								if dr != nil {
									c.Objs.List34Add(dr)
								}
								if on != 0 {
									p := [2]int32{int32(pos.X), int32(pos.Y)}
									legacy.PortTestClientEffects(7, nil, nil, [8]int32{int32(binary.LittleEndian.Uint32(cache)), int32(amount)}, unsafe.Pointer(&p[0]))
								}
							} else if on != 0 {
								legacy.PortTestClientEffects(2, nil, nil, [8]int32{4, 25, 500, 25, int32(pos.X), int32(pos.Y)}, nil)
							}
							want := snapshot()
							setup()
							data := []byte{240, byte(kind), byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8), amount}
							length := 7
							if kind == 26 {
								length = 6
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data[:length])
							got := snapshot()
							if n != length || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
								t.Fatal("quest particles", on, kind, mode, fail, amount, pos, n)
							}
							rows = append(rows, row{on, kind, mode, fail, int(amount), pos, n, got})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-particles", rows)
}

func TestGameMessageClientSessionGauntletStaticMarker(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On                       int
		Code                     uint16
		Initial, Static, Dynamic byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, code := range []uint16{0, 17, 0x8011, 65535} {
			for _, initial := range []byte{0, 1, 255} {
				c.resetCase(env, pix, 1, 100)
				binary.LittleEndian.PutUint32(connected, uint32(on))
				a := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(1, 2))
				b := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(3, 4))
				a.ObjClass = 0x20400000
				b.ObjClass = 0
				a.NetCode32, b.NetCode32 = 17, 17
				pa, pb := (*byte)(unsafe.Add(a.C(), 432)), (*byte)(unsafe.Add(b.C(), 432))
				*pa, *pb = initial, initial
				data := []byte{240, 15, byte(code), byte(code >> 8)}
				input := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
				want := initial
				if on != 0 && code == 17 {
					want = 0
				}
				if n != 4 || !bytes.Equal(input, data) || *pa != want || *pb != initial {
					t.Fatal("quest static marker", on, code, initial)
				}
				rows = append(rows, row{on, code, initial, *pa, *pb})
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-static-marker", rows)
}
