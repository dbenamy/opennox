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

func TestGameMessageClientSessionPlayerCamera(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldServer, oldTimes, oldCode, oldCamera := noxServer, inputKeyTimeoutsOld, legacy.ClientPlayerNetCode(), legacy.Nox_xxx_cliUpdateCameraPos_435600
	noxServer = c.srv
	t.Cleanup(func() {
		noxServer = oldServer
		inputKeyTimeoutsOld = oldTimes
		legacy.ClientSetPlayerNetCode(oldCode)
		legacy.Nox_xxx_cliUpdateCameraPos_435600 = oldCamera
	})
	var camera image.Point
	var calls int
	legacy.Nox_xxx_cliUpdateCameraPos_435600 = func(x, y int) { camera = image.Pt(x, y); calls++ }
	type row struct {
		On, High, Local, SameFrame, Existing, Return, Calls int
		Frame, Timeout                                      uint32
		Camera                                              image.Point
		Input                                               []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for high := 0; high < 2; high++ {
			for _, local := range []int{17, 0x8011, 999, 0x10011} {
				for same := 0; same < 2; same++ {
					for exists := 0; exists < 2; exists++ {
						for _, frame := range []uint32{0, 100, 0xffffffff} {
							c.resetCase(env, pix, 1, frame)
							c.Objs.LoadError = false
							binary.LittleEndian.PutUint32(connected, uint32(on))
							legacy.ClientSetPlayerNetCode(local)
							timeout := frame - 1
							if same != 0 {
								timeout = frame
							}
							inputKeyTimeoutsOld = map[byte]uint32{8: timeout}
							code := uint16(17) | uint16(high<<15)
							typ := uint16(0)
							if exists != 0 {
								typ = 4
								dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
								if dr == nil {
									t.Fatal("allocation")
								}
								dr.NetCode32 = 17
								dr.ObjClass = 0
								if high != 0 {
									dr.ObjClass = object.Class(0x20400000)
								}
							}
							data := []byte{195, byte(code), byte(code >> 8), byte(typ), 0, 0xff, 0xff, 0, 0x80, 0xf9, 17, 91}
							want := bytes.Clone(data)

							calls = 0
							camera = image.Pt(-1, -1)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(195), data)
							wantCalls := 0
							wantCamera := image.Pt(-1, -1)
							if int(code) == local && same == 0 {
								wantCalls = 1
								wantCamera = image.Pt(65535, 32768)
							}
							if n != 12 || !bytes.Equal(data, want) || calls != wantCalls || camera != wantCamera || len(inputKeyTimeoutsOld) != 1 || inputKeyTimeoutsOld[8] != timeout || c.Objs.LoadError != (exists == 0) || c.Objs.Count != exists {
								t.Fatalf("camera on%d high%d local%x same%d exists%d frame%x return%d calls%d/%d camera%v/%v", on, high, local, same, exists, frame, n, calls, wantCalls, camera, wantCamera)
							}
							rows = append(rows, row{on, high, local, same, exists, n, calls, frame, timeout, camera, data})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-player-camera", rows)
}
