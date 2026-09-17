//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestMapPolygonsAmbientColor(t *testing.T) {
	render := newObjectRenderOwner(t)
	o := newMapPolygonsOwner(t)
	o.construct(t, [][2]float32{{100, 200}, {110, 200}, {110, 210}, {100, 210}})
	serverConfigOwnBytes(t, 0x852978, 8, 4)
	previous := serverConfigOwnBytes(t, 0x5D4594, 811364, 8)
	ambient := serverConfigOwnBytes(t, 0x587000, 142296, 12)
	for i, v := range []uint32{31, 63, 127} {
		binary.LittleEndian.PutUint32(ambient[4*i:], v)
	}
	dr := render.drawable(7, image.Pt(102, 203))
	player := unsafe.Slice((*byte)(unsafe.Pointer(&render.players[0])), int(unsafe.Sizeof(render.players[0])))
	type row struct {
		Name  string
		Cache uint32
		Level byte
		Color noxrender.RGB
		Phase uint32
	}
	var rows []row
	for _, local := range []bool{false, true} {
		for _, known := range []bool{false, true} {
			for _, populated := range []bool{false, true} {
				for _, moving := range []bool{false, true} {
					for _, inside := range []bool{false, true} {
						for _, cache := range []uint32{0, 1, 0xdeadface} {
							name := fmt.Sprintf("local%t/known%t/populated%t/moving%t/inside%t/cache%x", local, known, populated, moving, inside, cache)
							t.Run(name, func(t *testing.T) {
								*memmap.PtrPtr(0x852978, 8) = nil
								if local {
									*memmap.PtrPtr(0x852978, 8) = unsafe.Pointer(dr)
								}
								dr.NetCode32 = 99
								if known {
									dr.NetCode32 = 7
								}
								*o.words["polygons"] = 1
								if populated {
									*o.words["polygons"] = 2
								}
								point := image.Pt(500, 400)
								if inside {
									point = image.Pt(102, 203)
								}
								render.c.Viewport().World.Max = point
								binary.LittleEndian.PutUint32(previous, uint32(point.X))
								binary.LittleEndian.PutUint32(previous[4:], uint32(point.Y))
								if moving {
									binary.LittleEndian.PutUint32(previous, uint32(point.X-1))
								}
								binary.LittleEndian.PutUint32(player[3660:], cache)
								player[3668] = 77
								binary.LittleEndian.PutUint32(o.control[8:], 1)
								start := noxrender.RGB{9, 19, 29}
								render.c.r.Data().SetLightColor(start)
								legacy.PortTestMapPolygonColor()
								wantCache, wantLevel, wantColor := cache, byte(77), start
								if local && known {
									if !populated {
										wantCache = 0
										wantLevel = 1
										wantColor = noxrender.RGB{31, 63, 127}
									} else if moving || cache == 0xdeadface {
										if !inside {
											wantCache = 0
											wantLevel = 1
											wantColor = noxrender.RGB{31, 63, 127}
										} else if cache != 1 {
											wantCache = 1
											wantLevel = 9
											wantColor = noxrender.RGB{117, 38, 219}
										}
									}
								}
								got := row{name, binary.LittleEndian.Uint32(player[3660:]), player[3668], render.c.r.Data().GetLightColor(), binary.LittleEndian.Uint32(o.control[8:])}
								if got.Cache != wantCache || got.Level != wantLevel || got.Color != wantColor {
									t.Fatalf("got %+v want %x/%d/%+v", got, wantCache, wantLevel, wantColor)
								}
								rows = append(rows, got)
							})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "map-polygons-ambient-color", rows, "3fe22379e6019d9f50b623343623adc1e1412b1c8a94410e31d2510232b6b46e")
}

func TestMapPolygonsSecretAwards(t *testing.T) {
	world := newWorldCollisionOwner(t)
	o := newMapPolygonsOwner(t)
	o.construct(t, [][2]float32{{100, 200}, {110, 200}, {110, 210}, {100, 210}})
	u := &world.units[0]
	pl := u.UpdateDataPlayer().Player
	b := unsafe.Slice((*byte)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl)))
	type row struct {
		Name                                       string
		Cache, Flags, Visited, Count, Total, Dirty uint32
		Sounds                                     []uint32
		Packets                                    [][]byte
	}
	var rows []row
	for _, quest := range []bool{false, true} {
		for _, secret := range []bool{false, true} {
			for _, visited := range []bool{false, true} {
				for _, first := range []bool{false, true} {
					for _, inactive := range []bool{false, true} {
						name := fmt.Sprintf("quest%t/secret%t/visited%t/first%t/inactive%t", quest, secret, visited, first, inactive)
						t.Run(name, func(t *testing.T) {
							world.reset()
							world.s.PortTestCombatAudioReset()
							noxflags.PortTestGameFlags(0)
							if quest {
								noxflags.PortTestGameFlags(4096)
							}
							pl.PlayerInd = 1
							u.ObjFlags = 0
							if inactive {
								u.ObjFlags = 0x20
							}
							u.PosVec = types.Pointf{102, 203}
							u.PrevPos = types.Pointf{101, 203}
							binary.LittleEndian.PutUint32(b[3664:], 0)
							if first {
								binary.LittleEndian.PutUint32(b[3664:], 0xdeadface)
							}
							binary.LittleEndian.PutUint32(b[4672:], 10)
							binary.LittleEndian.PutUint32(b[4676:], 20)
							binary.LittleEndian.PutUint32(b[4692:], 0x8000)
							p := o.polygon(1)
							flags := uint32(0x12340000)
							if secret {
								flags |= 1
							}
							binary.LittleEndian.PutUint32(p[132:], flags)
							mask := uint32(0x80000000)
							if visited {
								mask |= 2
							}
							binary.LittleEndian.PutUint32(p[136:], mask)
							binary.LittleEndian.PutUint32(o.control[8:], 1)
							legacy.PortTestMapPolygonActor("player", u.CObj())
							awarded := quest && secret && !visited && !first
							wantFlags := flags
							if awarded {
								wantFlags &^= 1
							}
							wantMask := mask
							if !first {
								wantMask |= 2
							}
							wantCount, wantTotal, wantDirty := uint32(10), uint32(20), uint32(0x8000)
							if awarded && !inactive {
								wantCount++
								wantTotal++
								wantDirty |= 16
							}
							r := row{Name: name, Cache: binary.LittleEndian.Uint32(b[3664:]), Flags: binary.LittleEndian.Uint32(p[132:]), Visited: binary.LittleEndian.Uint32(p[136:]), Count: binary.LittleEndian.Uint32(b[4672:]), Total: binary.LittleEndian.Uint32(b[4676:]), Dirty: binary.LittleEndian.Uint32(b[4692:]), Packets: visibilityEffectsPackets(world.s)}
							for _, e := range world.s.PortTestCombatAudioSnapshot() {
								if e.Obj != u || e.ID != 904 || e.Kind != 0 || e.Code != 0 || e.ByPos {
									t.Fatal("secret audio event")
								}
								r.Sounds = append(r.Sounds, uint32(e.ID))
							}
							if r.Cache != 1 || r.Flags != wantFlags || r.Visited != wantMask || r.Count != wantCount || r.Total != wantTotal || r.Dirty != wantDirty || len(r.Sounds) != bool2int(awarded) {
								t.Fatalf("secret result %+v awarded %t", r, awarded)
							}

							for index, packet := range r.Packets {
								var want []byte
								if awarded {
									for j := range world.units {
										other := &world.units[j]
										if int(other.UpdateDataPlayer().Player.PlayerInd) != index {
											continue
										}
										if other == u {
											want = append([]byte{169, 15, 0}, []byte("GeneralPrint:SecretFound\x00")...)
										} else {
											want = []byte{169, 20, 0, 0, 0, 0}
											binary.LittleEndian.PutUint32(want[2:], u.NetCode)
										}
									}
								}
								if !bytes.Equal(packet, want) {
									t.Fatalf("player %d notification: %x want %x", index, packet, want)
								}
							}
							rows = append(rows, r)
							// A repeated stationary check must never award a second time.
							u.PrevPos = u.PosVec
							legacy.PortTestMapPolygonActor("player", u.CObj())
							if binary.LittleEndian.Uint32(b[4672:]) != wantCount || len(world.s.PortTestCombatAudioSnapshot()) != len(r.Sounds) {
								t.Fatal("repeated award")
							}
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "map-polygons-secret-awards", rows, "e9a9d51234dcf62c138e4ac2fc52c109236a72d3b02c20cd533e507c05bd9466")
}

func TestMapPolygonsActorCacheInitialization(t *testing.T) {
	world := newWorldCollisionOwner(t)
	pl := world.units[0].UpdateDataPlayer().Player
	p := unsafe.Pointer(pl)
	b := unsafe.Slice((*byte)(p), int(unsafe.Sizeof(*pl)))
	before := append([]byte(nil), b...)
	if legacy.PortTestMapPolygonActorInit(nil) != nil || legacy.PortTestMapPolygonActorInit(p) != p {
		t.Fatal("initialization return")
	}
	for i, v := range b {
		want := before[i]
		if i >= 3660 && i < 3668 {
			want = []byte{0xce, 0xfa, 0xad, 0xde}[(i-3660)%4]
		}
		if v != want {
			t.Fatalf("initialization byte %d: %x want %x", i, v, want)
		}
	}
}
