//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientParticleDrawing(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t, "BlueRainSpark")
	var active [4]int
	var imagePos, imageSize image.Point
	c.r.HookImageDrawXxx = func(pos, size image.Point) { imagePos, imageSize = pos, size }
	defer func() {
		if failure := recover(); failure != nil {
			t.Fatalf("particle draw op%d variant%d allocation%d step%d image%v size%v: %v", active[0], active[1], active[2], active[3], imagePos, imageSize, failure)
		}
	}()

	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	flags := noxflags.GetGame()
	t.Cleanup(func() { noxflags.ResetGame(); noxflags.SetGame(flags) })
	c.callbackRefs = make(map[unsafe.Pointer]uint32)
	for op := 0; op < 12; op++ {
		if fn := legacy.PortTestClientParticleDrawCallback(op); fn != nil {
			c.callbackRefs[fn] = 0xe2000000 + uint32(op)
		}
	}
	type result struct {
		Op, Variant, Failure, Step int
		Return                     int64
		Pixels                     string
		Render, Globals            []uint32
		Drawables                  [][]uint32
		Calls                      []effectsSpawnCall
		Deleted                    []uint32
		Logic, Other               int
	}
	var out []result
	for op := 0; op < 12; op++ {
		for variant := 0; variant < 32; variant++ {
			for _, failure := range []int{0, 1, 2} {
				seed := uint32(1 + 31*variant)
				c.resetCase(effects, pix, seed, 120)
				env.Reset()
				noxflags.ResetGame()
				if variant&16 != 0 {
					noxflags.SetGame(noxflags.GameFlag(0x200000))
				}
				fps := []uint32{30, 60, 90, 120}[variant/8]
				c.srv.SetTickRate(fps)
				vp := []noxrender.Viewport{
					{Screen: image.Rect(0, 0, 96, 96), World: image.Rect(0, 0, 96, 96), Size: image.Pt(96, 96)},
					{Screen: image.Rect(10, 12, 106, 108), World: image.Rect(200, 300, 296, 396), Size: image.Pt(96, 96)},
					{Screen: image.Rect(-10, -12, 86, 84), World: image.Rect(0, 0, 96, 96), Size: image.Pt(96, 96)},
					{Screen: image.Rect(0, 0, 48, 48), World: image.Rect(300, 400, 348, 448), Size: image.Pt(48, 48)},
				}[variant/8]
				*c.Viewport() = vp
				pos := vp.World.Min.Add([]image.Point{image.Pt(0, 0), image.Pt(10, 10), image.Pt(11, 11), image.Pt(48, 48), image.Pt(85, 85), image.Pt(86, 86), image.Pt(95, 40), image.Pt(-1, 48)}[variant%8])
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
				dr.DrawFuncPtr = legacy.PortTestClientParticleDrawCallback(op)
				dr.ZVal = []uint16{0, 1, 8, 65535}[variant/8]
				dr.ZVal2 = uint16(variant % 3)
				dr.Field_8 = uint32(pos.X - []int{0, 1, 20, -20}[variant%4])
				dr.Field_9 = uint32(pos.Y - []int{0, 20, 1, -20}[variant%4])
				words := unsafe.Slice((*uint32)(dr.C()), 128)
				bytes := unsafe.Slice((*byte)(dr.C()), 512)
				if op == 2 || op == 3 {
					words[108], words[109] = uint32(pos.X-20), uint32(pos.Y-10)
					c.Objs.TransparentDecay(dr, []int{-1, 0, 1, 9, 10, 20, 60, 121}[variant%8])
				}
				if op == 5 {
					words[108], words[109] = 0xf81f, 0xffff
					bytes[440] = []byte{0, 1, 4, 5}[variant%4]
					bytes[441] = byte(1 + variant%3)
					// The lifecycle switches from growth at size five and from
					// shrinkage at size zero; do not invent an impossible phase/size pair.
					if bytes[441] == 1 && bytes[440] == 5 {
						bytes[440] = 4
					}
					if bytes[441] == 2 && bytes[440] == 0 {
						bytes[440] = 1
					}

					bytes[442], bytes[443] = byte(variant%3), byte(1+variant%4)
					bytes[444], bytes[445] = byte(1+variant%3), byte(variant%3)
					bytes[446] = byte(int8([]int{-3, -1, 1, 3}[variant/8]))
					if variant&4 != 0 {
						c.Objs.TransparentDecay(dr, []int{-1, 0, 1, 60}[variant%4])
					}
				}
				if op == 11 {
					words[108], words[109] = 0xbdef, 0xffff
					words[110], words[111] = uint32(pos.X), uint32(pos.Y)
					bytes[448] = []byte{0, 32, 64, 96, 128, 160, 192, 224}[variant%8]
					bytes[449] = byte(int8([]int{-3, -2, 2, 3}[variant/8]))
					bytes[450] = []byte{1, 10, 30, 50}[variant/8]
					bytes[451] = byte(1 + variant%4)
				}
				c.FailEvery = failure
				for step := 0; step < 4; step++ {
					c.srv.SetFrame(120 + uint32(step))
					beforeOther := c.srv.Rand.Other.Index()
					beforeCalls := len(c.Calls)
					beforeCount := c.Objs.Count
					beforeDeleted := len(c.Deleted)
					blank := effectsPixelHash(pix)
					inside := true
					if op == 0 || op == 1 {
						screen := dr.PosVec.Add(vp.Screen.Min).Sub(vp.World.Min)
						screen.Y -= int(int16(dr.ZVal)) + int(int16(dr.ZVal2))
						inside = screen.X-10 >= vp.Screen.Min.X && screen.Y-10 >= vp.Screen.Min.Y && screen.X+10 < vp.Screen.Max.X && screen.Y+10 < vp.Screen.Max.Y
					}
					active = [4]int{op, variant, failure, step}
					var got int64
					if dr.DrawFuncPtr != nil {
						got = int64(dr.CallDraw(c.Viewport()))
					} else {
						got = legacy.PortTestClientDrawParticle(op, c.Viewport(), dr, [4]int32{15})
					}
					calls := c.Calls[beforeCalls:]
					successes := 0
					for _, call := range calls {
						if call.Ref != 0 {
							successes++
						}
					}
					wantRNG, attempts := 0, 0
					switch op {
					case 0, 1:
						if inside {
							wantRNG = 2
						} else if effectsPixelHash(pix) != blank {
							t.Fatal("clipped magic changed pixels")
						}
					case 6:
						if !noxflags.HasGame(noxflags.GameFlag(0x200000)) {
							attempts = 2
							wantRNG = 4 + successes
						}
					case 7, 8, 9:
						attempts = 2
						wantRNG = 6 + successes
					}
					if len(calls) != attempts || c.srv.Rand.Other.Index()-beforeOther != wantRNG || c.srv.Rand.Logic.Index() != int(seed) {
						t.Fatalf("op%d v%d f%d step%d emission/RNG", op, variant, failure, step)
					}
					removed := len(c.Deleted) - beforeDeleted
					if c.Objs.Count != beforeCount+successes-removed {
						t.Fatal("draw ownership count")
					}
					if op == 9 && got != 0 {
						id := c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(uint32(got))))]
						if id == 0 {
							t.Fatal("unowned falling-spark return")
						}
						got = int64(id)
					}
					switch op {
					case 5, 11:
						if (got == 0) != (removed == 1) || (got != 0 && got != 1) {
							t.Fatal("draw deletion return")
						}
					case 9:
						if got != int64(calls[len(calls)-1].Ref) {
							t.Fatal("falling spark last result")
						}
					default:
						if got != 1 || removed != 0 {
							t.Fatal("particle draw return")
						}
					}
					render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
					out = append(out, result{op, variant, failure, step, got, effectsPixelHash(pix), render, env.Snapshot(), c.snapshotDrawables(t), append([]effectsSpawnCall(nil), calls...), append([]uint32(nil), c.Deleted...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
					if removed != 0 {
						break
					}
				}
			}
		}
	}
	effectsCapture(t, "particle-drawing", out, len(out), "7390a7366081a3fed2a09630095746f735ad3f505f019875b2905aa2057f6b76")
}
