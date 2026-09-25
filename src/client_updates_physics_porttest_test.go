//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"math"
	"testing"
	"unsafe"
)

func TestClientUpdatesHeightCatchup(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	type result struct {
		Op                            int
		Height, Velocity, Restitution uint32
		Previous, Frame               uint32
		Return                        uint32
		Drawables                     [][]uint32
	}
	var out []result
	for _, op := range []int{14, 15} {
		for _, height := range []float32{-1, 0, 0.5, 1, 2, 20, 255.75, 32768} {
			for _, velocity := range []float32{-8, -2, -0.5, 0, 0.5, 2, 8} {
				for _, restitution := range []float32{0, 5, 10, 20} {
					for _, frames := range [][2]uint32{{0, 0}, {1, 2}, {120, 127}, {120, 150}, {0xfffffffe, 0xffffffff}, {0xfffffffe, 1}} {
						c.resetCase(effects, pix, 1, frames[1])
						env.Reset(0)
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(320, 448))
						data := unsafe.Slice((*byte)(dr.C()), 512)
						binary.LittleEndian.PutUint32(data[432:], frames[0])
						binary.LittleEndian.PutUint32(data[436:], math.Float32bits(height))
						binary.LittleEndian.PutUint32(data[440:], math.Float32bits(velocity))
						binary.LittleEndian.PutUint32(data[444:], math.Float32bits(restitution))
						h, v := height, velocity
						for tick := frames[0]; tick < frames[1]; tick++ {
							if op == 14 {
								if h > 0 {
									h = float32(float64(h) + float64(v))
									v = float32(float64(v) - 1)
								}
								if h <= 0 {
									h, v = 0, 0
								}
							} else {
								next := float64(h) + float64(v)
								h = float32(next)
								if next >= 0 {
									v = float32(float64(v) - 0.5)
								} else {
									bounced := -float64(v) * float64(restitution) * 0.1
									h = 0
									v = float32(bounced)
									if bounced < 2 {
										h, v = 0, 0
									}
								}
							}
						}
						got := legacy.PortTestClientUpdate(op, c.Viewport(), dr, [4]int32{})
						if got != 1 || binary.LittleEndian.Uint32(data[432:]) != frames[1] || binary.LittleEndian.Uint32(data[436:]) != math.Float32bits(h) || binary.LittleEndian.Uint32(data[440:]) != math.Float32bits(v) || dr.ZVal != uint16(int64(h)) || dr.VelZ != int8(int64(v)) {
							t.Fatalf("height catchup op%d h%v v%v r%v frames%v", op, height, velocity, restitution, frames)
						}
						if c.Objs.Count != 1 || len(c.Calls) != 1 || c.srv.Rand.Logic.Index() != 1 || c.srv.Rand.Other.Index() != 2 {
							t.Fatal("height update affected unrelated ownership/RNG")
						}
						out = append(out, result{op, math.Float32bits(height), math.Float32bits(velocity), math.Float32bits(restitution), frames[0], frames[1], got, c.snapshotDrawables(t)})
					}
				}
			}
		}
	}
	effectsCapture(t, "update-height", out, len(out), "bc7c837bf2ca7b2865bd1209e0c01105381fecdd07cf8096ef15d3c002e1a03c")
}

func TestClientUpdatesStoredCloudCallback(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	type result struct {
		Start     uint16
		Speed     byte
		Step      int
		Return    int
		Drawables [][]uint32
	}
	var out []result
	for _, z := range []uint16{0, 1, 32767, 32768, 65534, 65535} {
		for _, speed := range []byte{0, 1, 3, 127, 255} {
			c.resetCase(effects, pix, 31, 127)
			env.Reset(0)
			parent := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(320, 448))
			legacy.PortTestClientUpdate(1, c.Viewport(), parent, [4]int32{1, 0})
			dr := c.Objs.List1
			if dr == parent || dr.Field_115 != legacy.PortTestClientUpdateCloudCallback() {
				t.Fatal("cloud creation failed to install production update callback")
			}
			dr.ZVal = z
			*(*byte)(unsafe.Add(dr.C(), 432)) = speed
			before := c.srv.Rand.Other.Index()
			for step := 0; step < 4; step++ {
				got := int(client.CallDrawableUpdateResult(dr.Field_115, c.Viewport(), dr))
				if got != 1 || dr.ZVal != z+uint16((step+1)*int(speed)) || c.srv.Rand.Other.Index() != before {
					t.Fatal("stored cloud callback increment/wrap")
				}
				out = append(out, result{z, speed, step, got, c.snapshotDrawables(t)})
			}
		}
	}
	effectsCapture(t, "update-cloud-callback", out, len(out), "c9fe11b89e517a17837e6815fb975359a31089fec7a9618359fe891d4ebd0704")
}
