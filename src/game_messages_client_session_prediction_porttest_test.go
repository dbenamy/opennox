//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionPrediction(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Type, Marked, Variant, Return, Count int
		Frame                                             uint32
		LoadError, Callback                               bool
		Position, After                                   image.Point
		Start, X, Y, VX, VY, Drag, Token, Listed          uint32
		MotionReturn                                      int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, typ := range []int{0, 4} {
				for marked := 0; marked < 2; marked++ {
					for _, frame := range []uint32{0, 0xffffffff} {
						for variant := 0; variant < 256; variant++ {
							// Invalid types have no motion fields to decode. Keep representative bytes;
							// the original factory lost its allocation on type-lookup failure.
							if present == 0 && typ == 0 && variant != 0 && variant != 127 && variant != 255 {
								continue
							}
							c.resetCase(env, pix, 1, frame)
							c.Objs.LoadError = false
							binary.LittleEndian.PutUint32(connected, uint32(on))
							var dr *client.Drawable
							if present != 0 {
								dr = c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(310, 410))
								if dr == nil {
									t.Fatal("allocation")
								}
								dr.NetCode32 = 17
								dr.AnimStart = 31
								dr.Field_81 = 32
								dr.Field_82 = 33
								dr.Field_117 = 34
								dr.Field_118 = 35
								dr.Field_119 = 36
								dr.Field_127 = 0xcafe1234
							}
							code := uint16(17 | marked<<15)
							data := []byte{181, byte(code), byte(code >> 8), byte(typ), 0, 44, 1, 144, 1, byte(variant), byte(255 - variant), byte(13 * variant), byte(variant), byte(255 - variant)}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(181), data)
							enabled := on != 0 && (present != 0 || typ != 0)
							count := present
							if enabled {
								count = 1
							}
							fail := on != 0 && present == 0 && typ == 0
							if n != 14 || c.Objs.Count != count || c.Objs.LoadError != fail || !bytes.Equal(data, input) {
								t.Fatal("prediction creation/input", on, present, typ, marked, variant, n, c.Objs.Count, c.Objs.LoadError)
							}
							r := row{On: on, Present: present, Type: typ, Marked: marked, Variant: variant, Return: n, Count: c.Objs.Count, Frame: frame, LoadError: c.Objs.LoadError}
							if count != 0 {
								dr = c.Objs.ByNetCode(17)
								if dr == nil {
									t.Fatal("prediction clears code marker")
								}
								r.Position = dr.PosVec
								r.Start = dr.AnimStart
								r.X = dr.Field_81
								r.Y = dr.Field_82
								r.VX = dr.Field_117
								r.VY = dr.Field_118
								r.Drag = dr.Field_119
								r.Token = dr.Field_127
								r.Listed = dr.InClientUpdateList
								r.Callback = dr.Field_115 != nil
								if enabled {
									vx, vy, drag := float32(int8(data[12]))/16, float32(int8(data[13]))/16, float32(int8(data[11]))/16
									token := uint32(binary.LittleEndian.Uint16(data[9:]))
									if present != 0 {
										token |= 0xcafe0000
									}
									if r.Position != image.Pt(300, 400) || r.Start != frame || r.X != 300 || r.Y != 400 || r.VX != math.Float32bits(vx) || r.VY != math.Float32bits(vy) || r.Drag != math.Float32bits(drag) || r.Token != token || r.Listed != 1 || !r.Callback || c.Objs.FirstList5() != dr {
										t.Fatal("prediction fields/list", r)
									}
									_, restore := c.Cli().PortTestObjectRenderSight([]image.Point{{-10000, -10000}, {10000, -10000}, {10000, 10000}, {-10000, 10000}})
									r.MotionReturn = int(client.CallDrawableUpdateResult(dr.Field_115, c.Viewport(), dr))
									restore()
									r.After = dr.PosVec
									want := image.Pt(int(300+vx*(1-drag)), int(400+vy*(1-drag)))
									if r.MotionReturn != 1 || c.Objs.Count != 1 || r.After != want {
										t.Fatal("installed prediction callback", r, want)
									}
								} else if r.Position != image.Pt(310, 410) || r.Start != 31 || r.X != 32 || r.Y != 33 || r.VX != 34 || r.VY != 35 || r.Drag != 36 || r.Token != 0xcafe1234 || r.Listed != 0 || r.Callback {
									t.Fatal("disconnected prediction changed object", r)
								}
							}
							rows = append(rows, r)
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-prediction", rows)
}
