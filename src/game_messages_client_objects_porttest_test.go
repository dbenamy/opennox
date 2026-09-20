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

// Exercise the original dispatch with real drawable owners. Expectations below
// describe the wire fields independently of the existing sprite helpers.
func TestGameMessageClientObjects(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	legacy.ClientSetPlayerNetCode(999)
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	type row struct {
		Connected, Kind, Special, Status                    int
		Frame                                               uint32
		Return                                              int
		Pos                                                 image.Point
		Anim, Start, Slave, Previous, Seen, Active, Pending uint32
		Direction                                           byte
		Input                                               []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{47, 48} {
			for special := 0; special < 2; special++ {
				for _, frame := range []uint32{0, 100, 0xffffffff} {
					for status := 0; status < 256; status++ {
						c.resetCase(env, pix, 1, frame)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
						if dr == nil {
							t.Fatal("drawable allocation")
						}
						dr.ObjClass = 0
						if special != 0 {
							dr.ObjClass = object.Class(2)
						}
						dr.ObjSubClass = 0x40000
						dr.NetCode32 = 17
						dr.AnimInd = 8
						dr.AnimStart = 77
						dr.AnimFrameSlave = 91
						dr.Field_78 = 92
						dr.Field_72 = 93
						dr.Field_80 = 94
						dr.Field_122 = 95
						dr.AnimDir = 3
						xy := []image.Point{image.Pt(0, 65535), image.Pt(32767, 32768), image.Pt(321, 654)}[status%3]
						data := []byte{byte(kind), 17, 0, 4, 0, 0, 0, 0, 0}
						binary.LittleEndian.PutUint16(data[5:], uint16(xy.X))
						binary.LittleEndian.PutUint16(data[7:], uint16(xy.Y))
						if kind == 48 {
							data = append(data, byte(status), byte(255-status))
						}
						wantInput := bytes.Clone(data)
						want := row{on, kind, special, status, frame, len(data), image.Pt(300, 400), 8, 77, 91, 92, 93, 94, 95, 3, wantInput}
						if on != 0 {
							want.Pos = xy
							if xy.X >= 5888 || xy.Y >= 5888 {
								want.Pos = image.Pt(50, 50)
							}
							want.Seen = frame
							want.Active = frame
							want.Pending = 0
							if kind == 48 {
								want.Direction = byte(status>>4) & 7
								if want.Direction > 3 {
									want.Direction++
								}
								want.Input[9] &= 15
								if special == 0 {
									want.Slave = uint32(255 - status)
									want.Previous = 91
								}
								want.Anim = uint32(status & 15)
								if want.Anim != 8 {
									want.Start = frame
								}
							}
						}
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
						got := row{on, kind, special, status, frame, n, dr.PosVec, dr.AnimInd, dr.AnimStart, dr.AnimFrameSlave, dr.Field_78, dr.Field_72, dr.Field_80, dr.Field_122, dr.AnimDir, data}
						if got.Return != want.Return || got.Pos != want.Pos || got.Anim != want.Anim || got.Start != want.Start || got.Slave != want.Slave || got.Previous != want.Previous || got.Seen != want.Seen || got.Active != want.Active || got.Pending != want.Pending || got.Direction != want.Direction || !bytes.Equal(got.Input, want.Input) || c.Objs.Count != 1 || c.Objs.ByNetCode(17) != dr {
							t.Fatalf("object update got %+v want %+v count %d", got, want, c.Objs.Count)
						}
						rows = append(rows, got)
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-objects", rows)
}

func TestGameMessageClientObjectCreation(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	legacy.ClientSetPlayerNetCode(999)
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	oldCamera, oldTimeout := legacy.Nox_xxx_cliUpdateCameraPos_435600, legacy.InputSetKeyTimeoutLegacy
	var camera image.Point
	var cameraCalls, timeouts int
	legacy.Nox_xxx_cliUpdateCameraPos_435600 = func(x, y int) { camera = image.Pt(x, y); cameraCalls++ }
	legacy.InputSetKeyTimeoutLegacy = func(key byte) {
		if key != 9 {
			t.Fatalf("timeout key %d", key)
		}
		timeouts++
	}
	t.Cleanup(func() {
		legacy.Nox_xxx_cliUpdateCameraPos_435600 = oldCamera
		legacy.InputSetKeyTimeoutLegacy = oldTimeout
	})
	typ := c.Things.TypeByInd(4)
	oldClass := typ.ObjClass
	t.Cleanup(func() { typ.ObjClass = oldClass })
	type row struct {
		On, Kind, Type, Code, Class, Return, Count, CameraCalls, Timeouts int
		LoadError                                                         bool
		Pos, Camera                                                       image.Point
		Seen, Anim, Slave                                                 uint32
		Input                                                             []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{47, 48} {
			for _, typeID := range []int{0, 4} {
				for _, code := range []uint16{0, 17} {
					for _, class := range []object.Class{0, 0x20400000} {
						for _, pos := range []image.Point{image.Pt(0, 0), image.Pt(5887, 5887), image.Pt(65535, 32768)} {
							c.resetCase(env, pix, 1, 100)
							c.Objs.LoadError = false
							binary.LittleEndian.PutUint32(connected, uint32(on))
							typ.ObjClass = class
							camera = image.Pt(-1, -1)
							cameraCalls = 0
							timeouts = 0
							wireCode := code
							if class != 0 {
								wireCode |= 0x8000
							}
							data := []byte{byte(kind), byte(wireCode), byte(wireCode >> 8), byte(typeID), 0, 0, 0, 0, 0}
							binary.LittleEndian.PutUint16(data[5:], uint16(pos.X))
							binary.LittleEndian.PutUint16(data[7:], uint16(pos.Y))
							if kind == 48 {
								data = append(data, 0xf9, 253)
							}
							wantInput := bytes.Clone(data)
							cameraOnly := kind == 48 && wireCode == 0 && typeID == 0
							create := on != 0 && !cameraOnly && typeID != 0
							fail := on != 0 && !cameraOnly && typeID == 0
							wantCount, wantCamera := 0, 0
							if create {
								wantCount = 1
								if kind == 48 {
									wantInput[9] = 9
								}
							}
							if on != 0 && cameraOnly {
								wantCamera = 1
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							r := row{On: on, Kind: kind, Type: typeID, Code: int(wireCode), Class: int(class), Return: n, Count: c.Objs.Count, CameraCalls: cameraCalls, Timeouts: timeouts, LoadError: c.Objs.LoadError, Camera: camera, Input: data}
							if n != len(data) || r.Count != wantCount || r.LoadError != fail || cameraCalls != wantCamera || timeouts != wantCamera || !bytes.Equal(data, wantInput) {
								t.Fatalf("creation %+v expected count%d error%v camera%d input%v", r, wantCount, fail, wantCamera, wantInput)
							}
							if wantCamera != 0 && camera != pos {
								t.Fatalf("camera %v want %v", camera, pos)
							}
							if create {
								dr := c.Objs.ByNetCode(wireCode)
								if dr == nil {
									t.Fatalf("missing created object %+v", r)
								}
								r.Pos = dr.PosVec
								r.Seen = dr.Field_72
								r.Anim = dr.AnimInd
								r.Slave = dr.AnimFrameSlave
								if r.Pos != pos || r.Seen != 100 || dr.Field_80 != 100 || !dr.ObjFlags.Has(object.FlagActive) {
									t.Fatalf("created state %+v", r)
								}
								if kind == 48 && (r.Anim != 9 || r.Slave != 253 || dr.AnimDir != 8) {
									t.Fatalf("created animation %+v", r)
								}
							}
							rows = append(rows, r)
						}
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-object-creation", rows)
}
