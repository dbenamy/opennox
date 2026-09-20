//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientEffectCreation(t *testing.T) {
	names := []string{"FireBoom", "MediumFireBoom", "CounterspellBoom", "ThinFireBoom", "TeleportPoof", "DamagePoof"}
	c, pix, env := newEffectsFullOwner(t, names...)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200872, 24)
	type row struct {
		On, Mode, Fail, Kind, Return int
		Pos                          image.Point
		Cache                        []byte
		Calls                        []effectsSpawnCall
		Drawables                    [][]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for fail := 0; fail < 2; fail++ {
				for which, kind := range []int{133, 134, 135, 136, 137, 139} {
					for _, pos := range []image.Point{image.Pt(-32768, -1), image.Pt(0, 0), image.Pt(1, 32767), image.Pt(32767, -32768)} {
						c.resetCase(env, pix, 1, 100)
						clear(cache)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						c.FailEvery = fail
						id := c.Things.IndByID(names[which])
						if id == 0 {
							t.Fatal("fixture effect type")
						}
						if mode == 1 {
							binary.LittleEndian.PutUint32(cache[which*4:], uint32(id))
						} else if mode == 2 {
							id = 4
							binary.LittleEndian.PutUint32(cache[which*4:], 4)
						}
						wantCache := bytes.Clone(cache)
						if on != 0 {
							binary.LittleEndian.PutUint32(wantCache[which*4:], uint32(id))
						}
						data := []byte{byte(kind), byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8)}
						input := bytes.Clone(data)
						where := pos
						if kind == 137 || kind == 139 {
							where.Y += 2
						}
						var wantCalls []effectsSpawnCall
						if on != 0 {
							ref := uint32(1)
							if fail != 0 {
								ref = 0
							}
							wantCalls = []effectsSpawnCall{{Type: id, Position: where, Ref: ref}}
						}
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
						count := 0
						if on != 0 && fail == 0 {
							count = 1
						}
						if n != 5 || !bytes.Equal(data, input) || !bytes.Equal(cache, wantCache) || !reflect.DeepEqual(c.Calls, wantCalls) || c.Objs.Count != count {
							t.Fatalf("effect on%d mode%d fail%d kind%d pos%v return%d calls%v want%v count%d", on, mode, fail, kind, pos, n, c.Calls, wantCalls, c.Objs.Count)
						}
						if count != 0 {
							dr := c.Objs.List1
							if dr == nil || dr.TypeIDVal != uint32(id) || dr.PosVec != where || dr.ObjFlags&0x400000 == 0 || c.Objs.List4 != dr {
								t.Fatal("created effect position/type/lifetime list")
							}
						}
						rows = append(rows, row{on, mode, fail, kind, n, pos, bytes.Clone(cache), append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t)})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-effect-creation", rows)
}
