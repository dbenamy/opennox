//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestGameMessageClientObjectLifetime(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	local := serverConfigOwnBytes(t, 0x852978, 8, 4)
	oldFade := legacy.Get_nox_client_fadeObjects_80836()
	t.Cleanup(func() { legacy.Set_nox_client_fadeObjects_80836(oldFade) })
	drawData, free := alloc.Make([]uint32{}, 4)
	defer free()
	type row struct {
		On, Present, Static, Local, Fade, Animation, Kind int
		Frame                                             uint32
		Return, Count                                     int
		Active                                            bool
		Deadline, Shadow1, Shadow2, Pending               uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for isLocal := 0; isLocal < 2; isLocal++ {
					for fade := 0; fade < 2; fade++ {
						for animation := 0; animation < 3; animation++ {
							for _, kind := range []int{50, 51} {
								for _, frame := range []uint32{0, 100, 0xfffffff0} {
									// Clear the mapped owner before deleting the previous drawable.
									clear(local)
									c.resetCase(env, pix, 1, frame)
									binary.LittleEndian.PutUint32(connected, uint32(on))
									legacy.Set_nox_client_fadeObjects_80836(fade)
									dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
									if dr == nil {
										t.Fatal("allocation")
									}
									dr.NetCode32 = 17
									dr.ObjClass = 0
									if static != 0 {
										dr.ObjClass = object.Class(0x20400000)
									}
									dr.ObjFlags |= object.FlagActive
									dr.Field_120 = 10
									dr.Field_121 = 11
									dr.Field_122 = 12
									dr.DrawFuncPtr = nil
									dr.DrawData = nil
									if animation != 0 {
										dr.DrawFuncPtr = legacy.Get_nox_thing_animate_draw()
										dr.DrawData = unsafe.Pointer(&drawData[0])
										drawData[3] = uint32(animation - 1)
									}
									if isLocal != 0 {
										binary.LittleEndian.PutUint32(local, uint32(uintptr(unsafe.Pointer(dr))))
									}
									code := uint16(17)
									if present == 0 {
										code = 18
									}
									if static != 0 {
										code |= 0x8000
									}
									data := []byte{byte(kind), byte(code), byte(code >> 8)}
									before := bytes.Clone(data)
									enabled := on != 0 && present != 0
									protected := isLocal != 0 || animation == 2
									remove := enabled && !protected && (kind == 50 || fade == 0)
									wantCount := 1
									if remove && static == 0 {
										wantCount = 0
									}
									n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
									r := row{On: on, Present: present, Static: static, Local: isLocal, Fade: fade, Animation: animation, Kind: kind, Frame: frame, Return: n, Count: c.Objs.Count}
									if n != 3 || !bytes.Equal(data, before) || r.Count != wantCount {
										t.Fatalf("lifetime %+v count want%d", r, wantCount)
									}
									if wantCount != 0 {
										found := c.Objs.ByNetCode(uint16(17) | uint16(static<<15))
										if found != dr {
											t.Fatalf("lifetime lookup %+v", r)
										}
										r.Active = dr.ObjFlags.Has(object.FlagActive)
										r.Deadline = dr.Deadline
										r.Shadow1 = dr.Field_120
										r.Shadow2 = dr.Field_121
										r.Pending = dr.Field_122
										wantDeadline := uint32(0)
										if enabled && kind == 51 && fade != 0 && isLocal == 0 {
											wantDeadline = frame + c.srv.TickRate()
										}
										a, b, d := uint32(10), uint32(11), uint32(12)
										if enabled && kind == 51 {
											a, b, d = 1, 1, 1
										}
										if r.Active == remove || r.Deadline != wantDeadline || r.Shadow1 != a || r.Shadow2 != b || r.Pending != d {
											t.Fatalf("lifetime fields %+v expected active%v deadline%d fields%d/%d/%d", r, !remove, wantDeadline, a, b, d)
										}
									}
									rows = append(rows, r)
								}
							}
						}
					}
				}
			}
		}
	}
	clear(local)
	interactionCapture(t, "game-client-object-lifetime", rows)
}
