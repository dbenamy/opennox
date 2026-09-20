//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"reflect"
	"testing"
)

func TestGameMessageClientSessionForget(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	t.Cleanup(noxflags.PortTestGameFlags(0))
	oldSummon := c.dword_5d4594_1046604
	c.dword_5d4594_1046604 = 999
	t.Cleanup(func() { c.dword_5d4594_1046604 = oldSummon })
	type row struct {
		On, Mode, Kind, Local int
		Cutoff, Frame         uint32
		Codes                 []uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 2; mode++ {
			for _, kind := range []int{233, 234} {
				for _, local := range []int{17, 0x8011, 999, 0x10011} {
					for _, cutoff := range []uint32{0, 1, 100, 0x80000000, 0xffffffff} {
						for _, frame := range []uint32{0, 99, 100, 0xffffffff} {
							c.resetCase(env, pix, 1, 100)
							noxflags.ResetGame()
							if mode != 0 {
								noxflags.SetGame(0x2000)
							}
							binary.LittleEndian.PutUint32(connected, uint32(on))
							legacy.ClientSetPlayerNetCode(local)
							active := mode != 0 && (kind == 234 || on != 0 && local == 17)
							var want []uint32
							for i, class := range []uint32{0, 4, 2, 0x400000, 0x20000000, 0} {
								dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300+i, 400))
								if dr == nil {
									t.Fatal("allocation")
								}
								dr.ObjClass = object.Class(class)
								dr.NetCode32 = uint32(100 + i)
								dr.Field_80 = frame
								if i == 5 {
									dr.TypeIDVal = 999
								}
								if !active || class != 0 || i == 5 || frame >= cutoff {
									want = append([]uint32{dr.NetCode32}, want...)
								}
							}
							data := []byte{byte(kind)}
							if kind == 233 {
								data = binary.LittleEndian.AppendUint16(data, 17)
							}
							data = binary.LittleEndian.AppendUint32(data, cutoff)
							if kind == 233 {
								data = append(data, 255, 0)
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							var codes []uint32
							for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
								codes = append(codes, dr.NetCode32)
							}
							if n != len(data) || !bytes.Equal(input, data) || !reflect.DeepEqual(codes, want) {
								t.Fatal("forget gate/full local ID/unsigned age", on, mode, kind, local, cutoff, frame, codes, want)
							}
							rows = append(rows, row{on, mode, kind, local, cutoff, frame, codes})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-forget", rows)
}
