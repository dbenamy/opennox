//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientUpdateStream(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	aliases := serverConfigOwnBytes(t, 0x5D4594, 1198020, 255*8+8)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCamera := legacy.Nox_xxx_cliUpdateCameraPos_435600
	var camera image.Point
	calls := 0
	legacy.Nox_xxx_cliUpdateCameraPos_435600 = func(x, y int) { camera = image.Pt(x, y); calls++ }
	t.Cleanup(func() { legacy.Nox_xxx_cliUpdateCameraPos_435600 = oldCamera })
	type row struct {
		On, Mode, Explicit, Return int
		Class, Frame               uint32
		Status                     byte
		Camera                     image.Point
		CameraCalls                int
		Input, Aliases             []byte
		Messages                   [][]byte
		Drawables                  [][]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 5; mode++ {
			for explicit := 0; explicit < 2; explicit++ {
				for _, cl := range []uint32{0, 0x200000, 0x200004} {
					for _, frame := range []uint32{0, 100, 0xffffffff} {
						for _, status := range []byte{0x03, 0xf3} {
							c.resetCase(env, pix, 1, frame)
							clear(aliases)
							c.srv.NetList.ResetAll()
							binary.LittleEndian.PutUint32(connected, uint32(on))
							camera = image.Pt(-1, -1)
							calls = 0
							first := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
							second := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(320, 420))
							if first == nil || second == nil {
								t.Fatal("stream fixture allocation")
							}
							first.NetCode32, second.NetCode32 = 17, 18
							first.ObjClass, second.ObjClass = object.Class(cl), object.Class(cl)
							binary.LittleEndian.PutUint16(aliases[8:], 17)
							binary.LittleEndian.PutUint16(aliases[10:], 4)
							binary.LittleEndian.PutUint16(aliases[16:], 18)
							binary.LittleEndian.PutUint16(aliases[18:], 4)
							data := []byte{164, 1}
							if explicit != 0 {
								data = []byte{164, 255, 17, 0, 4, 0}
							}
							data = binary.LittleEndian.AppendUint16(data, 300)
							data = binary.LittleEndian.AppendUint16(data, 400)
							data = append(data, status)
							if status&0x80 != 0 {
								data = append(data, 255)
							}
							data = append(data, 8)
							firstLength := len(data)
							wantSecond := image.Pt(320, 420)
							if mode == 2 {
								data = append(data, 2, 0x80, 0x7f)
								wantSecond = image.Pt(172, 527)
							}
							if mode == 3 || mode == 4 {
								data = append(data, 0, 2)
								x := uint16(321)
								if mode == 4 {
									x = 65535
								}
								data = binary.LittleEndian.AppendUint16(data, x)
								data = binary.LittleEndian.AppendUint16(data, 654)
								if mode == 3 {
									wantSecond = image.Pt(321, 654)
								}
							}
							if mode == 2 || mode == 3 {
								if cl&0x200000 != 0 {
									data = append(data, status)
									if status&0x80 != 0 {
										data = append(data, 254)
									}
									if cl&4 != 0 {
										data = append(data, 6)
									}
								}
							}
							if mode == 1 || mode == 2 || mode == 3 {
								data = append(data, 0, 0, 0)
							}
							if mode == 4 {
								data = append(data, 0xa5, 0x5a)
							}
							input := bytes.Clone(data)
							wantReturn := len(data)
							if mode == 0 {
								wantReturn = 0
							}
							if mode == 4 {
								wantReturn = firstLength + 6
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(31, netmsg.Op(164), data)
							if n != wantReturn || !bytes.Equal(data, input) || calls != 1 || camera != image.Pt(300, 400) || first.PosVec != image.Pt(300, 400) || second.PosVec != wantSecond || c.Objs.Count != 2 {
								t.Fatalf("stream on%d mode%d explicit%d class%x status%x frame%x return%d want%d positions%v/%v camera%v/%d", on, mode, explicit, cl, status, frame, n, wantReturn, first.PosVec, second.PosVec, camera, calls)
							}
							if first.AnimInd != 8 || first.Field_72 != frame {
								t.Fatal("first stream animation/activity")
							}
							if mode == 2 || mode == 3 {
								if second.Field_72 != frame {
									t.Fatal("second stream activity")
								}
							}
							var messages [][]byte
							c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { messages = append(messages, bytes.Clone(b)); return false })
							if explicit == 0 && len(messages) != 0 || explicit != 0 && len(messages) != 1 {
								t.Fatal("stream alias report count")
							}
							if explicit != 0 {
								m := messages[0]
								if len(m) != 10 || m[0] != 165 || binary.LittleEndian.Uint16(m[2:]) != 17 || binary.LittleEndian.Uint16(m[4:]) != 4 || binary.LittleEndian.Uint32(m[6:]) != 0xffffffff {
									t.Fatal("stream alias report fields")
								}
							}
							rows = append(rows, row{on, mode, explicit, n, cl, frame, status, camera, calls, bytes.Clone(data), bytes.Clone(aliases), messages, c.snapshotDrawables(t)})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-update-stream", rows)
}
