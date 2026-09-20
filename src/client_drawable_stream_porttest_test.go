//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"os"
	"testing"
)

func drawableStateCapture(t *testing.T, name string, rows any, want string) {
	t.Helper()
	data, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_DRAWABLE_STATE_CAPTURE"); p != "" {
		if err := os.WriteFile(p+"-"+name+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %s", name, hash)
	if want != "" && hash != want {
		t.Fatalf("%s: got %s want frozen C %s", name, hash, want)
	}
}

func TestClientDrawableStreams(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	raw := memmap.Slice(0x5D4594, 1198020)[:255*8+8]
	saved := append([]byte(nil), raw...)
	t.Cleanup(func() { copy(raw, saved) })
	oldCamera := legacy.Nox_xxx_cliUpdateCameraPos_435600
	var camera image.Point
	var cameraCalls int
	legacy.Nox_xxx_cliUpdateCameraPos_435600 = func(x, y int) { camera = image.Pt(x, y); cameraCalls++ }
	t.Cleanup(func() { legacy.Nox_xxx_cliUpdateCameraPos_435600 = oldCamera })
	type result struct {
		Stream, Mode, Status, Class, Step          int
		Frame                                      uint32
		Return                                     int32
		Pos                                        [2]int32
		Camera                                     image.Point
		CameraCalls                                int
		Anim, Start, Slave, Previous, Seen, Active uint32
		Direction                                  byte
		Count                                      int
	}
	var rows []result
	classes := []object.Class{0, 0x200000, 0x200004, 0x200002, 0x20400000}
	for _, stream := range []int{1, 2} {
		for mode := 0; mode < 3; mode++ {
			if stream == 1 && mode == 2 {
				continue
			}
			for ci, cl := range classes {
				for status := 0; status < 256; status++ {
					frame := []uint32{0, 100, 0xffffffff}[status%3]
					c.resetCase(env, pix, 1, frame)
					clear(raw)
					camera = image.Pt(-1, -1)
					cameraCalls = 0
					c.srv.NetList.ResetByInd(server.HostPlayerIndex, netlist.Kind0)
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
					if dr == nil {
						t.Fatal("allocation")
					}
					dr.ObjClass = cl
					dr.ObjSubClass = 0x40000
					dr.NetCode32 = 17
					code := uint16(17)
					if cl&0x20400000 != 0 {
						code |= 0x8000
					}
					binary.LittleEndian.PutUint16(raw[8:], code)
					binary.LittleEndian.PutUint16(raw[10:], 4)
					dr.AnimInd = 8
					dr.AnimStart = 77
					dr.AnimFrameSlave = 91
					dr.Field_78 = 92
					dr.Field_72 = 93
					for step := 0; step < 2; step++ {
						beforeAnim, beforeStart, beforeSlave, beforePrev := dr.AnimInd, dr.AnimStart, dr.AnimFrameSlave, dr.Field_78
						packet := []byte{1}
						if mode == 1 {
							packet = []byte{0xff, byte(code), byte(code >> 8), 4, 0}
						}
						wantPos := [2]int32{321, 654}
						start := [2]int32{300, 400}
						if stream == 1 {
							packet = binary.LittleEndian.AppendUint16(packet, 321)
							packet = binary.LittleEndian.AppendUint16(packet, 654)
						} else if mode == 2 {
							packet = append([]byte{0}, packet...)
							packet = binary.LittleEndian.AppendUint16(packet, 321)
							packet = binary.LittleEndian.AppendUint16(packet, 654)
						} else {
							packet = append(packet, 0x80, 0x7f)
							wantPos = [2]int32{172, 527}
						}
						coordBytes := len(packet)
						packet = append(packet, byte(status))
						if status&0x80 != 0 {
							packet = append(packet, byte(255-step))
						}
						if stream == 1 || cl&4 != 0 {
							packet = append(packet, byte(8+step))
						}
						ret, pos := legacy.PortTestDrawableStream(stream, packet, start)
						wantRet := len(packet)
						animated := stream == 1 || cl&0x200000 != 0
						if !animated {
							wantRet = coordBytes
						}
						if ret != int32(wantRet) || pos != wantPos || dr.PosVec != image.Pt(int(wantPos[0]), int(wantPos[1])) || c.Objs.Count != 1 || c.Objs.ByNetCode(code) != dr {
							t.Fatalf("stream%d mode%d class%x status%x step%d return%d/%d pos%v/%v count%d", stream, mode, cl, status, step, ret, wantRet, pos, wantPos, c.Objs.Count)
						}
						wantAnim, wantStart, wantSlave, wantPrev, wantDir := beforeAnim, beforeStart, beforeSlave, beforePrev, byte(0)
						if animated {
							wantDir = byte(status>>4) & 7
							if wantDir > 3 {
								wantDir++
							}
							anim := uint32(status & 15)
							if stream == 1 || cl&4 != 0 {
								anim = uint32(8 + step)
							}
							// Stream 1 sets animation before frame; stream 2 sets frame first.
							frameAnim := beforeAnim
							if stream == 1 {
								frameAnim = anim
							}
							if (stream == 1 || status&0x80 != 0) && !(cl&2 != 0 && frameAnim == 8) {
								wantPrev = beforeSlave
								wantSlave = 0
								if status&0x80 != 0 {
									wantSlave = uint32(255 - step)
								}
							}
							wantAnim = anim
							if beforeAnim != anim {
								wantStart = frame
							}
						}
						if dr.AnimInd != wantAnim || dr.AnimStart != wantStart || dr.AnimFrameSlave != wantSlave || dr.Field_78 != wantPrev || dr.AnimDir != wantDir || dr.Field_72 != frame || dr.ObjFlags&4 == 0 {
							t.Fatalf("animation mismatch stream%d mode%d class%x status%x step%d: anim/start/frame/prev/dir %d/%d/%d/%d/%d want %d/%d/%d/%d/%d", stream, mode, cl, status, step, dr.AnimInd, dr.AnimStart, dr.AnimFrameSlave, dr.Field_78, dr.AnimDir, wantAnim, wantStart, wantSlave, wantPrev, wantDir)
						}
						wantCalls := 0
						if stream == 1 {
							wantCalls = step + 1
							if camera != dr.PosVec {
								t.Fatal("camera position")
							}
						}
						if cameraCalls != wantCalls {
							t.Fatal("camera callback count")
						}
						rows = append(rows, result{stream, mode, status, ci, step, frame, ret, pos, camera, cameraCalls, dr.AnimInd, dr.AnimStart, dr.AnimFrameSlave, dr.Field_78, dr.Field_72, uint32(dr.ObjFlags), dr.AnimDir, c.Objs.Count})
					}
				}
			}
		}
	}
	drawableStateCapture(t, "streams", rows, "aeede2ef66f63bc96a79e320d4db973c7440abf25c8df72eb2a9426fcb84a154")
}
