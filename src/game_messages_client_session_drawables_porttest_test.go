//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionDrawables(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	particles := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(particles.Restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Static, Special, Kind, Value, Return             int
		Frame                                                         uint32
		Angle                                                         byte
		Class, Flags, AnimFrame, Previous, AnimStart, ActivationFrame uint32
		Light                                                         [4]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for special := 0; special < 2; special++ {
					for _, kind := range []int{178, 179, 180} {
						for _, frame := range []uint32{0, 0xffffffff} {
							for value := 0; value < 256; value++ {
								c.resetCase(env, pix, 1, frame)
								particles.Reset()
								binary.LittleEndian.PutUint32(connected, uint32(on))
								dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
								if dr == nil {
									t.Fatal("allocation")
								}
								dr.ObjClass = 0
								if static != 0 {
									dr.ObjClass = object.Class(0x20400000)
								}
								if special != 0 {
									dr.ObjClass |= object.ClassMonster
								}
								dr.ObjClass |= object.Class(0x80000)
								dr.ObjSubClass = 0x40000
								dr.ObjFlags &^= object.FlagActive
								dr.ObjFlags |= object.Flags(0x80000)
								dr.NetCode32 = 17
								dr.AnimInd = 8
								dr.AnimFrameSlave = 91
								dr.Field_78 = 92
								dr.AnimStart = 31
								dr.Field_72 = 29
								dr.LightIntensity = -3.25
								dr.LightIntensityRad = 0x12345678
								dr.LightIntensityU16 = 0x87654321
								raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
								raw[299] = 0x5a
								before := bytes.Clone(raw)
								want := bytes.Clone(raw)
								if on != 0 && present != 0 {
									switch kind {
									case 178:
										raw[299] = byte(value)
									case 179:
										dr.SetActive()
										dr.SetLightIntensity(float32(16 * value / 10))
										dr.SetFrameMB(8 * value / 50)
										if dr.AnimFrameSlave == 8 {
											dr.AnimFrameSlave = 7
										}
									case 180:
										if value != 0 {
											dr.ObjClass |= object.Class(0x80000)
											dr.SetLightIntensity(41.958)
										} else {
											dr.ObjClass &^= object.Class(0x80000)
											dr.SetLightIntensity(0)
										}
										dr.SetFrameMB(value)
										dr.Field_72 = frame
									}
									copy(want, raw)
									copy(raw, before)
								}
								code := uint16(17)
								if present == 0 {
									code = 18
								}
								if static != 0 {
									code |= 0x8000
								}
								data := []byte{byte(kind), byte(code), byte(code >> 8), byte(value)}
								input := bytes.Clone(data)
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
								if n != 4 || !bytes.Equal(data, input) || !bytes.Equal(raw, want) {
									t.Fatalf("drawable kind%d on%d present%d static%d special%d value%d frame%d return%d", kind, on, present, static, special, value, frame, n)
								}
								rows = append(rows, row{on, present, static, special, kind, value, n, frame, raw[299], uint32(dr.ObjClass), uint32(dr.ObjFlags), dr.AnimFrameSlave, dr.Field_78, dr.AnimStart, dr.Field_72, [4]uint32{dr.LightFlags, math.Float32bits(dr.LightIntensity), dr.LightIntensityRad, dr.LightIntensityU16}})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-drawables", rows)
}
