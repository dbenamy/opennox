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

func TestGameMessageClientLightIntensity(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Static, Return int
		Bits                        uint32
		Light                       [4]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for _, value := range []float32{-10, math.Float32frombits(0x80000000), 0, math.SmallestNonzeroFloat32, 0.5, 1, math.Nextafter32(1, 2), 31, math.Nextafter32(31, 32), 63, math.Nextafter32(63, 64), 100} {
					c.resetCase(effects, pix, 1, 100)
					env.Reset()
					binary.LittleEndian.PutUint32(connected, uint32(on))
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
					if dr == nil {
						t.Fatal("allocation")
					}
					dr.ObjClass = 0
					if static != 0 {
						dr.ObjClass = object.Class(0x20400000)
					}
					dr.NetCode32 = 17
					dr.LightIntensity = -3.25
					dr.LightIntensityRad = 0x12345678
					dr.LightIntensityU16 = 0x87654321
					raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
					before := bytes.Clone(raw)
					want := bytes.Clone(raw)
					if on != 0 && present != 0 {
						// The lower-level light primitive is separately qualified against frozen C.
						legacy.PortTestClientDrawParticle(16, c.Viewport(), dr, [4]int32{int32(math.Float32bits(value))})
						copy(want, raw)
						copy(raw, before)
						clamped := value
						if clamped > 63 {
							clamped = 63
						}
						if binary.LittleEndian.Uint32(want[140:]) != math.Float32bits(clamped) || binary.LittleEndian.Uint32(want[148:]) != uint32(int64(float64(clamped)*65536+0.5)) {
							t.Fatal("independent intensity/clamp/fixed-point contract")
						}
					}
					code := uint16(17)
					if present == 0 {
						code = 18
					}
					if static != 0 {
						code |= 0x8000
					}
					data := binary.LittleEndian.AppendUint32([]byte{93, byte(code), byte(code >> 8)}, math.Float32bits(value))
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(93), data)
					if n != 7 || !bytes.Equal(data, input) || !bytes.Equal(raw, want) {
						t.Fatalf("intensity on%d present%d static%d value%g return%d", on, present, static, value, n)
					}
					rows = append(rows, row{on, present, static, n, math.Float32bits(value), [4]uint32{dr.LightFlags, math.Float32bits(dr.LightIntensity), dr.LightIntensityRad, dr.LightIntensityU16}})
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-light-intensity", rows)
}
