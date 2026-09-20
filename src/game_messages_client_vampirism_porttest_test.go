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
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientVampirism(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "HealOrb")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200848, 4)
	type row struct {
		On, Mode, Fail, Amount int
		Seed                   uint32
		Coords                 [4]uint16
		Cache                  uint32
		Calls                  []effectsSpawnCall
		Drawables              [][]uint32
		RNG                    [2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for _, fail := range []int{0, 1, 2} {
				for _, amount := range []uint16{0, 1, 3, 4, 7, 8, 27, 28, 31, 32, 33, 255, 256, 65535} {
					for _, seed := range []uint32{1, 1023} {
						for _, xy := range [][4]uint16{{0, 0, 0, 0}, {100, 200, 300, 400}, {65535, 65534, 65533, 65532}} {
							c.resetCase(env, pix, seed, 100)
							c.FailEvery = fail
							binary.LittleEndian.PutUint32(connected, uint32(on))
							typ := uint32(c.Things.IndByID("HealOrb"))
							initial := uint32(0)
							if mode == 1 {
								initial = typ
							}
							if mode == 2 {
								initial = 4
								typ = 4
							}
							binary.LittleEndian.PutUint32(cache, initial)
							data := make([]byte, 11)
							data[0] = 162
							for i, v := range xy {
								binary.LittleEndian.PutUint16(data[1+2*i:], v)
							}
							binary.LittleEndian.PutUint16(data[9:], amount)
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(162), data)
							count := int(amount)/4 + 1
							if count > 8 {
								count = 8
							}
							if on == 0 {
								count = 0
							}
							wantCache := initial
							if on != 0 {
								wantCache = typ
							}
							if n != 11 || !bytes.Equal(data, input) || len(c.Calls) != count || binary.LittleEndian.Uint32(cache) != wantCache {
								t.Fatal("vampirism count/cache/input")
							}
							rng := prand.New(int(seed + 1))
							live := 0
							for i, call := range c.Calls {
								speed := rng.Int(6, 12)
								dy := rng.Int(-20, 20)
								dx := rng.Int(-20, 20)
								alive := fail == 0 || (i+1)%fail != 0
								if call.Type != int(typ) || call.Position != image.Pt(int(xy[2])+dx, int(xy[3])+dy) || (call.Ref != 0) != alive {
									t.Fatal("vampirism type/position/allocation")
								}
								if !alive {
									continue
								}
								live++
								phase := rng.Int(3, 10)
								for dr, ref := range c.refs {
									if ref != call.Ref {
										continue
									}
									raw := unsafe.Slice((*byte)(dr.C()), 512)
									if binary.LittleEndian.Uint16(raw[432:]) != xy[0] || binary.LittleEndian.Uint16(raw[434:]) != xy[1] || raw[443] != byte(speed) || raw[444] != byte(phase) || raw[445] != 0 || raw[446] != 0 || dr.ObjFlags&0x400000 == 0 {
										t.Fatal("vampirism destination/speed/period/lifetime")
									}
								}
							}
							if c.Objs.Count != live || c.srv.Rand.Other.Index() != rng.Index() {
								t.Fatal("vampirism live count/random consumption")
							}
							rows = append(rows, row{on, mode, fail, int(amount), seed, xy, binary.LittleEndian.Uint32(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-vampirism", rows)
}
