//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientUpdatesTransfers(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	type result struct {
		Op, Variant, Binding, Failure int
		Return                        uint32
		Calls                         []effectsSpawnCall
		Drawables                     [][]uint32
		Globals                       []uint32
		Logic, Other                  int
	}
	var out []result
	for _, op := range []int{0, 4, 6, 21} {
		for variant := 0; variant < 16; variant++ {
			for binding := 0; binding < 8; binding++ {
				for _, failure := range []int{0, 1, 2} {
					seed := uint32(1 + variant*31)
					c.resetCase(effects, pix, seed, 127)
					env.Reset(0)
					vp := []noxrender.Viewport{
						{Screen: image.Rect(0, 0, 96, 96), World: image.Rect(0, 0, 96, 96), Size: image.Pt(96, 96)},
						{Screen: image.Rect(10, 12, 106, 108), World: image.Rect(200, 300, 296, 396), Size: image.Pt(96, 96)},
						{Screen: image.Rect(-10, -12, 86, 84), World: image.Rect(0, 0, 96, 96), Size: image.Pt(96, 96)},
						{Screen: image.Rect(0, 0, 48, 48), World: image.Rect(300, 400, 348, 448), Size: image.Pt(48, 48)},
					}[variant/4]
					*c.Viewport() = vp
					endpoints := [2]image.Point{vp.World.Min.Add(image.Pt([]int{-10, 1, 47, 100}[variant%4], 25)), vp.World.Min.Add(image.Pt(90, []int{-10, 1, 47, 100}[variant%4]))}
					parent := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(320, 448))
					parent.ZVal = []uint16{0, 20, 32767, 65535}[variant/4]
					data := unsafe.Slice((*byte)(parent.C()), 512)
					missing := binding >= 3 && binding <= 5
					if binding == 0 {
						for i, p := range endpoints {
							binary.LittleEndian.PutUint16(data[437+4*i:], uint16(p.X))
							binary.LittleEndian.PutUint16(data[439+4*i:], uint16(p.Y))
						}
					} else {
						data[432] = 1
						for i, p := range endpoints {
							code := uint32(101 + i)
							if binding == 2 || binding == 7 {
								code |= 0x8000
							}
							if binding == 6 {
								code |= 0xabcd0000
							}
							if binding == 7 {
								code |= 0x12340000
							}
							binary.LittleEndian.PutUint32(data[437+4*i:], code)
							if binding == 5 || binding == 3 && i == 1 || binding == 4 && i == 0 {
								continue
							}
							dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, p)
							dr.NetCode32 = uint32(101 + i)
							if binding == 2 || binding == 7 {
								dr.ObjClass |= 0x20000000
							}
						}
					}
					before := len(c.Calls)
					beforeCount := c.Objs.Count
					beforeOther := c.srv.Rand.Other.Index()
					c.FailEvery = failure
					got := legacy.PortTestClientUpdate(op, c.Viewport(), parent, [4]int32{3, int32(variant % 2)})
					calls := append([]effectsSpawnCall(nil), c.Calls[before:]...)
					gates := 1
					if op == 21 {
						gates = 2
					}
					if len(calls) > gates || missing && len(calls) != 0 {
						t.Fatal("transfer gating/missing endpoint")
					}
					successes := 0
					for _, call := range calls {
						typ := 3
						if op == 4 {
							typ = 1
						}
						if op == 6 {
							typ = 2
						}
						if call.Type != typ {
							t.Fatal("transfer particle type")
						}
						if call.Ref == 0 {
							continue
						}
						successes++
						var child *client.Drawable
						for p, ref := range c.refs {
							if ref == call.Ref {
								child = p
								break
							}
						}
						if child == nil {
							t.Fatal("unowned transfer particle")
						}
						b := unsafe.Slice((*byte)(child.C()), 512)
						target := image.Pt(int(binary.LittleEndian.Uint16(b[432:])), int(binary.LittleEndian.Uint16(b[434:])))
						expected := [2]image.Point{image.Pt(int(uint16(endpoints[0].X)), int(uint16(endpoints[0].Y))), image.Pt(int(uint16(endpoints[1].X)), int(uint16(endpoints[1].Y)))}
						if op == 21 {
							if target != expected[0] && target != expected[1] {
								t.Fatal("bidirectional charm target")
							}
						} else {
							index := 0
							if op == 0 && variant%2 == 0 {
								index = 1
							}
							if target != expected[index] {
								t.Fatal("transfer target coordinates")
							}
						}
						if b[443] < 6 || b[443] > 12 || b[444] < 3 || b[444] > 10 {
							t.Fatal("transfer speed/radius range")
						}
					}
					if c.Objs.Count != beforeCount+successes || len(c.Deleted) != 0 {
						t.Fatal("transfer ownership")
					}
					if c.srv.Rand.Logic.Index() != int(seed) || c.srv.Rand.Other.Index()-beforeOther != gates+2*len(calls)+2*successes {
						t.Fatal("transfer RNG ordering/count")
					}
					want := uint32(1)
					if op == 0 {
						want = 0
					}
					if got != want {
						t.Fatal("transfer return")
					}
					out = append(out, result{op, variant, binding, failure, got, calls, c.snapshotDrawables(t), env.Snapshot(), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				}
			}
		}
	}
	effectsCapture(t, "update-transfers", out, len(out), "45569e43638732ef34eac69611e6385dccf6f4b3e09ee0e075d679ab66358d05")
}
