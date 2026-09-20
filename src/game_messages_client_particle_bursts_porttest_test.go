//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientParticleBursts(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "BlueSpark", "CyanSpark")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200860, 4)
	cyan := serverConfigOwnBytes(t, 0x5D4594, 1200784, 4)
	type row struct {
		On, Kind, Fail, Return int
		Seed, Frame            uint32
		Radius                 float64
		Pos                    image.Point
		Cache                  []byte
		Calls                  []effectsSpawnCall
		Drawables              [][]uint32
		RNG                    [2]int
	}
	var rows []row
	for _, kind := range []byte{150, 163} {
		for on := 0; on < 2; on++ {
			for _, fail := range []int{0, 1, 2} {
				for _, frame := range []uint32{0, 100, 0xfffffff0} {
					for _, rad := range []float64{0, 1, 30.75, 255.999999} {
						if kind == 150 && rad != 0 {
							continue
						}
						for _, seed := range []uint32{1, 1023} {
							for _, pos := range []image.Point{image.Pt(-32768, -1), image.Pt(100, 200), image.Pt(32767, 32767)} {
								c.resetCase(env, pix, seed, frame)
								c.FailEvery = fail
								clear(cache)
								binary.LittleEndian.PutUint32(connected, uint32(on))
								binary.LittleEndian.PutUint32(cyan, uint32(c.Things.IndByID("CyanSpark")))
								restore := c.srv.Server.PortTestClientUpdateBalance(rad)
								data := []byte{kind, byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8)}
								input := bytes.Clone(data)
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
								restore()
								count := 0
								if on != 0 {
									count = 5
									if kind == 163 {
										count = 150
									}
								}
								if n != 5 || !bytes.Equal(input, data) || len(c.Calls) != count {
									t.Fatal("particle burst gate/count/input")
								}
								typ := c.Things.IndByID("BlueSpark")
								if kind == 163 {
									typ = c.Things.IndByID("CyanSpark")
								}
								wantCache := uint32(0)
								if kind == 150 && on != 0 {
									wantCache = uint32(typ)
								}
								if binary.LittleEndian.Uint32(cache) != wantCache {
									t.Fatal("particle burst cache gate")
								}
								rng := prand.New(int(seed + 1))
								live := 0
								radius := int(float32(rad))
								for i, call := range c.Calls {
									where := pos
									if kind == 163 {
										r := radius/4 + rng.Int(0, radius)
										if r > radius {
											r = radius
										}
										a := rng.Int(0, 255)
										where.X += r * int(*memmap.PtrInt32(0x587000, 192088+uintptr(8*a))) / 16
										where.Y += r * int(*memmap.PtrInt32(0x587000, 192092+uintptr(8*a))) / 16
									}
									if call.Type != typ || call.Position != where {
										t.Fatal("particle burst type/position")
									}
									alive := fail == 0 || (i+1)%fail != 0
									if (call.Ref != 0) != alive {
										t.Fatal("particle burst allocation schedule")
									}
									if !alive {
										continue
									}
									live++
									angle, speed, ttl, z, vel := 0, 0, 0, 0, 0
									if kind == 150 {
										angle = rng.Int(0, 255)
										speed = rng.Int(1333, 4000)
										ttl = rng.Int(5, 20)
										z = 20
										vel = rng.Int(-5, 5)
									} else {
										ttl = rng.Int(30, 40)
										vel = rng.Int(4, 10)
									}
									for dr, ref := range c.refs {
										if ref != call.Ref {
											continue
										}
										raw := unsafe.Slice((*byte)(dr.C()), 512)
										if *txword(dr, 432) != uint32(dr.PosVec.X)<<12 || *txword(dr, 436) != uint32(dr.PosVec.Y)<<12 || *txword(dr, 440) != uint32(speed) || *txword(dr, 444) != frame || *txword(dr, 448) != frame+uint32(ttl) || raw[299] != byte(angle) || int(int8(raw[296])) != vel || int(binary.LittleEndian.Uint16(raw[104:])) != z || dr.ObjFlags&0x400000 == 0 {
											t.Fatal("particle burst fields/lifetime")
										}
									}
								}
								if c.Objs.Count != live || c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != prand.New(int(seed)).Index() {
									t.Fatal("particle burst ownership/RNG")
								}
								rows = append(rows, row{on, int(kind), fail, n, seed, frame, rad, pos, bytes.Clone(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-particle-bursts", rows)
}
