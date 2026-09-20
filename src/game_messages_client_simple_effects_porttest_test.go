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

func TestGameMessageClientArrowTrap(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "ArrowTrap1Smoke", "ArrowTrap2Smoke")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200864, 8)
	type row struct {
		On, Mode, Fail, Direction, Return int
		Pos                               image.Point
		Cache                             []byte
		Calls                             []effectsSpawnCall
		Drawables                         [][]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for fail := 0; fail < 2; fail++ {
				for _, dir := range []byte{0, 1, 2, 127, 128, 255} {
					for _, pos := range []image.Point{image.Pt(-32768, -1), image.Pt(0, 0), image.Pt(100, 200), image.Pt(32767, -32768)} {
						c.resetCase(env, pix, 1, 100)
						c.FailEvery = fail
						clear(cache)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						ids := [2]int{c.Things.IndByID("ArrowTrap1Smoke"), c.Things.IndByID("ArrowTrap2Smoke")}
						if mode == 2 {
							ids = [2]int{4, 4}
						}
						if mode != 0 {
							for i, id := range ids {
								binary.LittleEndian.PutUint32(cache[i*4:], uint32(id))
							}
						}
						wantCache := bytes.Clone(cache)
						if on != 0 {
							for i, id := range ids {
								binary.LittleEndian.PutUint32(wantCache[i*4:], uint32(id))
							}
						}
						data := []byte{161, byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8), dir}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(161), data)
						count := 0
						if on != 0 {
							count = 1
						}
						if n != 6 || !bytes.Equal(input, data) || !bytes.Equal(cache, wantCache) || len(c.Calls) != count {
							t.Fatal("arrow trap gate/cache/input")
						}
						if on != 0 {
							typ := ids[1]
							where := pos.Add(image.Pt(-3, 0))
							if dir == 1 {
								typ = ids[0]
								where = pos.Add(image.Pt(15, 0))
							}
							call := c.Calls[0]
							if call.Type != typ || call.Position != where || (call.Ref == 0) != (fail != 0) {
								t.Fatal("arrow trap direction/position/allocation")
							}
							if fail == 0 && (c.Objs.Count != 1 || c.Objs.List1.ObjFlags&0x400000 == 0) {
								t.Fatal("arrow trap lifetime")
							}
						}
						rows = append(rows, row{on, mode, fail, int(dir), n, pos, bytes.Clone(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t)})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-arrow-trap", rows)
}

func TestGameMessageClientGreenBolt(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "GreenZap")
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200844, 4)
	type row struct {
		On, Mode, Fail, Return int
		Coords                 [4]uint16
		Value                  uint16
		Cache                  []byte
		Calls                  []effectsSpawnCall
		Drawables              [][]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for _, xy := range [][4]uint16{{0, 0, 0, 0}, {0, 0, 1, 1}, {1, 1, 0, 0}, {100, 200, 101, 199}, {0, 65535, 65535, 0}, {32767, 32768, 32768, 32767}} {
				for _, value := range []uint16{0, 1, 255, 32767, 32768, 65535} {
					for fail := 0; fail < 2; fail++ {
						c.resetCase(env, pix, 1, 100)
						c.FailEvery = fail
						clear(cache)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						id := c.Things.IndByID("GreenZap")
						if mode == 2 {
							id = 4
						}
						if mode != 0 {
							binary.LittleEndian.PutUint32(cache, uint32(id))
						}
						data := []byte{152}
						for _, v := range xy {
							data = binary.LittleEndian.AppendUint16(data, v)
						}
						data = binary.LittleEndian.AppendUint16(data, value)
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(152), data)
						count := 0
						if on != 0 {
							count = 1
						}
						if n != 11 || !bytes.Equal(data, input) || binary.LittleEndian.Uint32(cache) != uint32(id) || len(c.Calls) != count {
							t.Fatal("green bolt ungated cache/gated allocation/input")
						}
						if on != 0 && fail == 0 {
							where := image.Pt(int(xy[0])+(int(xy[2])-int(xy[0]))/2, int(xy[1])+(int(xy[3])-int(xy[1]))/2)
							if c.Calls[0].Position != where || c.Calls[0].Type != id || c.Objs.Count != 1 {
								t.Fatal("green bolt midpoint truncation/type")
							}
							dr := c.Objs.List1
							b := unsafe.Slice((*byte)(unsafe.Add(dr.C(), 432)), 13)
							if b[0] != 0 || binary.LittleEndian.Uint32(b[1:]) != uint32(value) || !bytes.Equal(b[5:13], data[1:9]) {
								t.Fatalf("green bolt packed effect data on%d mode%d xy%v value%x got%x input%x", on, mode, xy, value, b, data)
							}
						}
						rows = append(rows, row{on, mode, fail, n, xy, value, bytes.Clone(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t)})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-green-bolt", rows)
}
