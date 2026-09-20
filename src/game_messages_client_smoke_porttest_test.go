//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSmokeBlast(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "Smoke", "Puff")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200896, 8)
	type row struct {
		On, Mode, Fail, Return int
		Seed                   uint32
		Pos                    image.Point
		Cache                  []byte
		Calls                  []effectsSpawnCall
		Drawables              [][]uint32
		RNG                    [2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 4; mode++ {
			for _, fail := range []int{0, 1, 2, 3} {
				for _, seed := range []uint32{0, 1, 1023} {
					for _, pos := range []image.Point{image.Pt(-32768, -1), image.Pt(0, 0), image.Pt(100, 200), image.Pt(32767, -32768)} {
						c.resetCase(env, pix, seed, 100)
						c.FailEvery = fail
						clear(cache)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						smoke, puff := c.Things.IndByID("Smoke"), c.Things.IndByID("Puff")
						if smoke == 0 || puff == 0 {
							t.Fatal("smoke fixture types")
						}
						if mode != 0 {
							if mode == 2 {
								smoke, puff = 4, 4
							}
							if mode == 3 {
								puff = 0
							}
							binary.LittleEndian.PutUint32(cache, uint32(puff))
							binary.LittleEndian.PutUint32(cache[4:], uint32(smoke))
						}
						wantCache := bytes.Clone(cache)
						if on != 0 && mode == 0 {
							binary.LittleEndian.PutUint32(wantCache, uint32(puff))
							binary.LittleEndian.PutUint32(wantCache[4:], uint32(smoke))
						}
						data := []byte{138, byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(138), data)
						wantCalls := 0
						if on != 0 {
							wantCalls = 7
						}
						if n != 5 || !bytes.Equal(input, data) || !bytes.Equal(cache, wantCache) || len(c.Calls) != wantCalls {
							t.Fatalf("smoke on%d mode%d fail%d seed%d pos%v return%d calls%d", on, mode, fail, seed, pos, n, len(c.Calls))
						}
						live := 0
						for i, call := range c.Calls {
							typ := puff
							if i == 0 {
								typ = smoke
							}
							if call.Type != typ {
								t.Fatal("smoke/puff type order")
							}
							dx, dy := call.Position.X-pos.X, call.Position.Y-pos.Y
							if i == 0 && (dx != 0 || dy != 0) || i != 0 && (dx < -15 || dx > 15 || dy < -15 || dy > 15) {
								t.Fatal("smoke signed position/offset")
							}
							shouldLive := typ != 0 && (fail == 0 || (i+1)%fail != 0)
							if (call.Ref != 0) != shouldLive {
								t.Fatal("smoke allocation schedule")
							}
							if shouldLive {
								live++
							}
						}
						if c.Objs.Count != live {
							t.Fatal("smoke object count")
						}
						for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
							ref := c.refs[dr]
							first := len(c.Calls) > 0 && c.Calls[0].Ref == ref
							z := *(*uint16)(unsafe.Add(dr.C(), 104))
							if first && z != 20 || !first && (z < 5 || z > 25) || dr.ObjFlags&0x400000 == 0 {
								t.Fatal("smoke height/lifetime")
							}
						}
						rows = append(rows, row{on, mode, fail, n, seed, pos, bytes.Clone(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-smoke", rows)
}
