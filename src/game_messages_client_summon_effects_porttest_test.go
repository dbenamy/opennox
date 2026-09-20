//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSummonEffects(t *testing.T) {
	base, pix, env := newEffectsFullOwner(t, "SummonEffect")
	c := &drawableSummonOwner{effectsTestClient: base}
	old := legacy.GetClient
	legacy.GetClient = func() legacy.Client { return c }
	t.Cleanup(func() { legacy.GetClient = old })
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Failure, Direction int
		Owner                  uint16
		Frame                  uint32
		Pos                    image.Point
		Start, Stopped         [][]uint32
		Calls                  []effectsSpawnCall
		Globals                []uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for failure := 0; failure < 3; failure++ {
			for _, owner := range []uint16{0, 7, 0x8007, 0xffff} {
				for _, dir := range []byte{0, 1, 127, 128, 255} {
					for _, frame := range []uint32{0, 100, 0xffffffff} {
						for _, pos := range []image.Point{image.Pt(300, 400), image.Pt(65535, 65534)} {
							c.resetCase(env, pix, 17, frame)
							for _, p := range words {
								*p = 0
							}
							binary.LittleEndian.PutUint32(connected, uint32(on))
							c.FailChild = failure == 2
							if failure == 1 {
								c.FailEvery = 1
							}
							value := owner ^ 0xa5a5
							data := []byte{126}
							for _, v := range []uint16{uint16(pos.X), uint16(pos.Y), owner, 4} {
								data = binary.LittleEndian.AppendUint16(data, v)
							}
							data = append(data, dir)
							data = binary.LittleEndian.AppendUint16(data, value)
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(126), data)
							count := 0
							if on != 0 && failure != 1 {
								count = 1
								if failure == 0 {
									count = 2
								}
							}
							if n != 12 || !bytes.Equal(data, input) || c.Objs.Count != count {
								t.Fatal("summon start gate/input/allocation")
							}
							parent := c.Objs.List1
							var child *client.Drawable
							var start [][]uint32
							if count == 2 {
								child = *(**client.Drawable)(unsafe.Add(parent.C(), 432))
								if child == nil || child.PosVec != pos || child.AnimInd != 8 || parent.AnimStart != frame || *txword(parent, 436) != (uint32(owner)<<16|uint32(value)) {
									t.Fatal("summon wire field widths/order")
								}
								start = c.snapshotDrawables(t, child)
								start[0][109] = c.refs[child]
							} else {
								start = c.snapshotDrawables(t)
							}
							// A disconnected cancel must preserve an already-created summon too.
							binary.LittleEndian.PutUint32(connected, 0)
							cancel := []byte{127, byte(owner), byte(owner >> 8)}
							cancelInput := bytes.Clone(cancel)
							if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(127), cancel); n != 3 || !bytes.Equal(cancel, cancelInput) || c.Objs.Count != count {
								t.Fatal("disconnected summon cancel")
							}
							var after [][]uint32
							if child != nil {
								after = c.snapshotDrawables(t, child)
								after[0][109] = c.refs[child]
							} else {
								after = c.snapshotDrawables(t)
							}
							if !reflect.DeepEqual(start, after) {
								t.Fatal("disconnected cancel changed summon")
							}
							if count == 2 {
								binary.LittleEndian.PutUint32(connected, 1)
								wrong := []byte{127, byte(owner), byte((owner ^ 0x8000) >> 8)}
								if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(127), wrong); n != 3 || c.Objs.Count != 2 {
									t.Fatal("summon owner must retain high bit")
								}
								if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(127), cancel); n != 3 || c.Objs.Count != 50 {
									t.Fatal("summon matching cancel")
								}
								for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
									if dr == parent || dr == child {
										t.Fatal("cancel retained summon allocation")
									}
								}
							}
							globals := make([]uint32, len(words))
							for i, p := range words {
								globals[i] = *p
							}
							rows = append(rows, row{on, failure, int(dir), owner, frame, pos, start, c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), globals})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-summon-effects", rows)
}
