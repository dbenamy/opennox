//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientTurnUndead(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "UndeadKiller")
	for _, r := range blobdata.PortTestClientPresentationTables() {
		copy(serverConfigOwnBytes(t, r.Base, r.Offset, len(r.Data)), r.Data)
	}
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1217508, 4)
	type row struct {
		On, Mode, Fail, Return int
		Frame                  uint32
		Pos                    [2]int16
		Cache                  uint32
		Calls                  []effectsSpawnCall
		Particles              [][8]uint32
		RNG                    [2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for _, fail := range []int{0, 1, 2, 43} {
				for _, frame := range []uint32{0, 100, 0xffffffff} {
					for _, pos := range [][2]int16{{0, 0}, {123, -456}, {-32768, 32767}} {
						c.resetCase(env, pix, 31, frame)
						c.FailEvery = fail
						clear(cache)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						typ := c.Things.IndByID("UndeadKiller")
						if mode == 2 {
							typ = 4
						}
						if mode != 0 {
							binary.LittleEndian.PutUint32(cache, uint32(typ))
						}
						wantCache := binary.LittleEndian.Uint32(cache)
						if on != 0 {
							wantCache = uint32(typ)
						}
						data := []byte{160, byte(pos[0]), byte(uint16(pos[0]) >> 8), byte(pos[1]), byte(uint16(pos[1]) >> 8)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(160), data)
						count := 0
						if on != 0 {
							count = 43
						}
						if n != 5 || !bytes.Equal(input, data) || len(c.Calls) != count || binary.LittleEndian.Uint32(cache) != wantCache || c.srv.Rand.Logic.Index() != 31 || c.srv.Rand.Other.Index() != 32 {
							t.Fatal("turn undead gate/cache/input/RNG")
						}
						byRef := make(map[uint32][8]uint32)
						for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
							angle := uint16(dr.Field_127)
							vx := math.Float32bits(*memmap.PtrFloat32(0x587000, 194136+uintptr(angle)*8) * 4)
							vy := math.Float32bits(*memmap.PtrFloat32(0x587000, 194140+uintptr(angle)*8) * 4)
							if dr.Field_117 != vx || dr.Field_118 != vy || dr.Field_119 != 0 || dr.AnimStart != frame || int32(dr.Field_81) != int32(pos[0]) || int32(dr.Field_82) != int32(pos[1]) || dr.PosVec != image.Pt(int(pos[0]), int(pos[1])) || dr.Field_115 == nil || dr.InClientUpdateList != 1 || dr.ObjFlags&0x200000 == 0 {
								t.Fatal("turn undead trajectory/owner")
							}
							byRef[c.refs[dr]] = [8]uint32{uint32(angle), vx, vy, dr.Field_81, dr.Field_82, dr.AnimStart, uint32(dr.ObjFlags), dr.InClientUpdateList}
						}
						var particles [][8]uint32
						for i, call := range c.Calls {
							alive := fail == 0 || (i+1)%fail != 0
							if call.Type != typ || call.Position != image.Pt(int(pos[0]), int(pos[1])) || (call.Ref != 0) != alive {
								t.Fatal("turn undead allocation")
							}
							if alive {
								p, ok := byRef[call.Ref]
								if !ok || p[0] != uint32(i*6) {
									t.Fatal("turn undead angle order")
								}
								particles = append(particles, p)
							}
						}
						if c.Objs.Count != len(particles) {
							t.Fatal("turn undead live count")
						}
						rows = append(rows, row{on, mode, fail, n, frame, pos, wantCache, append([]effectsSpawnCall(nil), c.Calls...), particles, [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-turn-undead", rows)
}
