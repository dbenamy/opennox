//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSparkExplosion(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "Spark", "MediumFireBoom", "FireBoom")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	caches := serverConfigOwnBytes(t, 0x5D4594, 1200852, 8)
	large := serverConfigOwnBytes(t, 0x5D4594, 1197380, 4)
	type row struct {
		On, Mode, Fail, Amount, Return int
		Pos                            image.Point
		Cache                          [3]uint32
		Calls                          []effectsSpawnCall
		Drawables                      [][]uint32
		RNG                            [2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 4; mode++ {
			for _, fail := range []int{0, 1, 2} {
				for _, amount := range []byte{0, 1, 169, 170, 171, 254, 255} {
					for _, pos := range []image.Point{image.Pt(-32768, -1), image.Pt(0, 0), image.Pt(100, 200), image.Pt(32767, -32768)} {
						c.resetCase(env, pix, 31, 100)
						c.FailEvery = fail
						clear(caches)
						clear(large)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						ids := [3]uint32{uint32(c.Things.IndByID("Spark")), uint32(c.Things.IndByID("MediumFireBoom")), uint32(c.Things.IndByID("FireBoom"))}
						if mode == 2 {
							ids = [3]uint32{4, 4, 4}
						}
						if mode == 3 {
							ids[1], ids[2] = 0, 0
						}
						if mode != 0 {
							binary.LittleEndian.PutUint32(caches, ids[0])
							binary.LittleEndian.PutUint32(caches[4:], ids[1])
							binary.LittleEndian.PutUint32(large, ids[2])
						}
						data := []byte{147, byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8), amount}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(147), data)
						gotCache := [3]uint32{binary.LittleEndian.Uint32(caches), binary.LittleEndian.Uint32(caches[4:]), binary.LittleEndian.Uint32(large)}
						if n != 6 || !bytes.Equal(input, data) || gotCache != ids {
							t.Fatal("spark explosion ungated cache/input")
						}
						count := 0
						if on != 0 {
							count = 180*int(amount)/255 + 11
						}
						if len(c.Calls) != count {
							t.Fatal("spark explosion count")
						}
						live := 0
						for i, call := range c.Calls {
							typ := ids[0]
							if i == count-1 {
								typ = ids[1]
								if amount > 170 {
									typ = ids[2]
								}
							}
							alive := typ != 0 && (fail == 0 || (i+1)%fail != 0)
							if call.Type != int(typ) || call.Position != pos || (call.Ref != 0) != alive {
								t.Fatalf("spark explosion type/threshold/allocation on%d mode%d fail%d amount%d call%d", on, mode, fail, amount, i)
							}
							if alive {
								live++
							}
						}
						if c.Objs.Count != live {
							t.Fatal("spark explosion live count")
						}
						for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
							if dr.ObjFlags&0x400000 == 0 {
								t.Fatal("spark explosion lifetime list")
							}
						}
						rows = append(rows, row{on, mode, fail, int(amount), n, pos, gotCache, append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-spark-explosion", rows)
}
