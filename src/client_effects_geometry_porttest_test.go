//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientEffectsCurveRaster(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Geometry, Steps, Glow, Thick, Size int
		Setup                              [2]uint32
		Pixels                             string
		Render, Globals                    []uint32
		Logic, Other                       int
	}
	var out []result
	buf, free := alloc.Make([]byte{}, 48)
	defer free()
	geometry := [][4][2]int32{
		{{10, 10}, {85, 85}, {75, 75}, {75, 75}},
		{{10, 10}, {85, 85}, {100, -50}, {-40, 90}},
		{{-20, 48}, {116, 48}, {136, 0}, {136, 0}},
		{{48, 48}, {48, 48}, {0, 0}, {0, 0}},
	}
	for pi, p := range geometry {
		for _, steps := range []int{-1, 0, 1, 2, 3, 8, 32} {
			for glow := 0; glow < 2; glow++ {
				for thick := 0; thick < 2; thick++ {
					for _, size := range []int{0, 1, 127, 128, 255} {
						c.resetCase(env, pix, 1, 0)
						line := legacy.PortTestClientEffects(38, c.Viewport(), nil, [8]int32{0xfc01}, nil)
						setup := legacy.PortTestClientEffects(39, c.Viewport(), nil, [8]int32{int32(glow), 0x8421, 8, int32(size)}, nil)
						if line != 0xfc01 || setup != uint32(int32(int8(size))) {
							t.Fatal("curve setting return conversion")
						}
						for i := range buf {
							buf[i] = 0xa5
						}
						*(*[4][2]int32)(unsafe.Pointer(&buf[8])) = p
						before := append([]byte(nil), buf...)
						blank := effectsPixelHash(pix)
						legacy.PortTestClientEffects(40, c.Viewport(), nil, [8]int32{int32(steps), int32(thick)}, unsafe.Pointer(&buf[8]))
						if !bytes.Equal(before, buf) {
							t.Fatal("curve raster input guards")
						}
						if steps <= 0 && effectsPixelHash(pix) != blank {
							t.Fatal("empty curve drew pixels")
						}
						if pi == 0 && steps > 0 && steps&(steps-1) == 0 && pix.Pix[pix.PixOffset(85, 85)] == 0 {
							t.Fatal("linear curve missed visible endpoint")
						}
						h := 1.0 / float64(steps)
						if *env.Named("dword_587000_180480") != math.Float32bits(float32(h*h)) || *env.Named("dword_587000_180476") != math.Float32bits(float32(h*h*h)) {
							t.Fatal("curve finite-difference step coefficients")
						}
						if c.srv.Rand.Logic.Index() != 1 || c.srv.Rand.Other.Index() != 2 {
							t.Fatal("curve raster consumed RNG")
						}
						render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
						if pi == 1 && steps == 8 && glow == 1 && thick == 0 && size == 127 {
							effectsDiagnostic(t, "curve", pix)
						}
						out = append(out, result{pi, steps, glow, thick, size, [2]uint32{line, setup}, effectsPixelHash(pix), render, env.Snapshot(), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
					}
				}
			}
		}
	}
	effectsCapture(t, "curve-raster", out, len(out), "b1d085491d177adb5b898153c9ada741a1eaa085b749d4625faeee905998c605")
}

func TestClientEffectsPlasmaSegments(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Slot, Index, Failure int
		Frame, Return        uint32
		Globals              []uint32
		Drawables            [][]uint32
		Calls                []effectsSpawnCall
		Logic, Other         int
	}
	var out []result
	buf, free := alloc.Make([]byte{}, 40)
	defer free()
	for slot := 0; slot < 3; slot++ {
		for _, index := range []int{0, 1, 28} {
			for _, frame := range []uint32{0, 1, 4, 5, 127, 128, 255} {
				for _, failure := range []int{0, 2} {
					c.resetCase(env, pix, 31, frame)
					c.FailEvery = failure
					c.Viewport().World = image.Rect(100, 200, 196, 296)
					c.Viewport().Screen = image.Rect(3, 7, 96, 96)
					*env.Named("dword_5d4594_1316412") = uint32(index)
					for i := range buf {
						buf[i] = 0xa5
					}
					*(*[5]int32)(unsafe.Pointer(&buf[8])) = [5]int32{10, 10, 85, 85, int32(slot)}
					before := append([]byte(nil), buf...)
					got := legacy.PortTestClientEffects(16, c.Viewport(), nil, [8]int32{}, unsafe.Pointer(&buf[8]))
					if !bytes.Equal(before, buf) || *env.Named("dword_5d4594_1316412") != uint32(index+1) {
						t.Fatal("plasma segment input/cursor")
					}
					off := uintptr(28 * (index + 30*slot))
					for i, v := range []uint32{10, 10} {
						if *memmap.PtrUint32(0x5D4594, 1313884+off+uintptr(i*4)) != v {
							t.Fatal("plasma segment start")
						}
					}
					for i, v := range []uint32{85, 85} {
						if *memmap.PtrUint32(0x5D4594, 1313912+off+uintptr(i*4)) != v {
							t.Fatal("plasma segment end")
						}
					}
					rng := prand.New(32)
					want := uint32(int32(int8(byte(frame))))
					attempts := 0
					if frame&4 != 0 {
						chance := rng.Int(0, 10)
						want = uint32(chance)
						if chance > 5 {
							attempts = 2
							want = 0
						}
					}
					if got != want || len(c.Calls) != attempts {
						t.Fatal("plasma segment emission/return")
					}
					for _, call := range c.Calls {
						if call.Type != 8 || call.Position != image.Pt(182, 278) {
							t.Fatal("plasma segment world-space emission")
						}
						if call.Ref != 0 {
							for i := 0; i < 5; i++ {
								rng.Int(0, 255)
							}
						}
					}
					if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != 31 {
						t.Fatal("plasma segment RNG consumption")
					}
					out = append(out, result{slot, index, failure, frame, got, env.Snapshot(), c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				}
			}
		}
	}
	effectsCapture(t, "plasma-segments", out, len(out), "376f4e947e2c40c46692f893d4734f133bc1a5fc4eeb49f2876a84418f17ae5f")
}

func TestClientEffectsPlasmaGeometry(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Op, Angle, Length, Step int
		Frame, Return           uint32
		Pixels                  string
		Render, Globals         []uint32
		Drawables               [][]uint32
		Calls                   []effectsSpawnCall
		Logic, Other            int
	}
	var out []result
	for _, op := range []int{35, 15} {
		for _, angle := range []int{0, 64, 127, 128, 255} {
			for _, length := range []int{0, 1, 39, 40, 1120, 2000} {
				for _, frame := range []uint32{0, 4} {
					c.resetCase(env, pix, 31, frame)
					args := [8]int32{int32(angle), 48, 48, int32(48 + length), 48}
					count := length/40 + 1
					if length/40+2 >= 30 {
						count = 28
					}
					for step := 0; step < 6; step++ {
						c.srv.SetFrame(frame + uint32(step))
						got := legacy.PortTestClientEffects(op, c.Viewport(), nil, args, nil)
						want := uint32(0)
						if op == 15 {
							want = uint32(count)
						}
						if got != want || *env.Named("dword_5d4594_1316408") != uint32(count) || *env.Named("dword_5d4594_1316412") != uint32(count) {
							t.Fatal("plasma geometry segment count")
						}
						for slot := 0; slot < 3; slot++ {
							off := uintptr(840 * slot)
							if *memmap.PtrUint32(0x5D4594, 1313884+off) != 48 || *memmap.PtrUint32(0x5D4594, 1313888+off) != 48 {
								t.Fatal("plasma curve origin")
							}
							phase := *memmap.PtrFloat32(0x5D4594, 1313856+uintptr(slot*4))
							if phase < 0 || phase >= 1 {
								t.Fatal("plasma phase wrap")
							}
						}
						if c.srv.Rand.Logic.Index() != 31 {
							t.Fatal("plasma changed logic RNG")
						}
						render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
						if op == 15 && angle == 64 && length == 40 && frame == 4 && step == 3 {
							effectsDiagnostic(t, "plasma", pix)
						}
						out = append(out, result{op, angle, length, step, frame, got, effectsPixelHash(pix), render, env.Snapshot(), c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
					}
				}
			}
		}
	}
	effectsCapture(t, "plasma-geometry", out, len(out), "dc1bed6dadcd83543ad527ec839f985a754dc864cf35f437901679f7181ebf1e")
}

func TestClientEffectsRayDrawing(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	oldFlags := noxflags.GetGame()
	defer func() { noxflags.ResetGame(); noxflags.SetGame(oldFlags) }()
	type result struct {
		Op, Variant, Step        int
		Return                   uint32
		Pixels                   string
		Render, Globals          []uint32
		Drawables                [][]uint32
		Calls                    []effectsSpawnCall
		Deleted                  []uint32
		MouseReads, Logic, Other int
	}
	var out []result
	for _, op := range []int{11, 12, 13, 14, 17} {
		for variant := 0; variant < 32; variant++ {
			frame := []uint32{0, 4, 127, 0xfffffffc}[variant/4%4]
			c.resetCase(env, pix, 31, frame)
			noxflags.ResetGame()
			if variant&8 != 0 {
				noxflags.SetGame(noxflags.GamePause)
			}
			legacy.PortTestClientEffects(45, c.Viewport(), nil, [8]int32{}, nil)
			c.Mouse = image.Pt(96-variant, variant)
			mode := variant % 4
			from, to := image.Pt(20, 40), image.Pt(20+[]int{0, 1, 63, 64, 127, 128, 255, 512}[variant/4%8], 60)
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(14, image.Pt(48, 70))
			data := unsafe.Slice((*byte)(unsafe.Pointer(dr)), 512)
			data[433] = []byte{0, 1, 2, 3, 127, 128, 129, 255}[variant/4%8]
			if op == 14 {
				binary.LittleEndian.PutUint32(data[433:], uint32(variant/4%3))
			}
			missing := mode == 3
			if mode == 0 {
				for i, v := range []int{from.X, from.Y, to.X, to.Y} {
					binary.LittleEndian.PutUint16(data[437+2*i:], uint16(v))
				}
			} else {
				data[432] = 1
				for i, pos := range []image.Point{from, to} {
					code := uint32(101 + i)
					if mode == 2 {
						code |= 0x8000
					}
					binary.LittleEndian.PutUint32(data[437+i*4:], code)
					if mode == 3 && i == 1 {
						continue
					}
					end := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
					end.NetCode32 = code
					if mode == 2 {
						end.ObjClass |= 0x20000000
						// Preserve and expose the existing difference in code normalization:
						// plasma clears the static flag; the four other draw callbacks do not.
						if op == 17 {
							end.NetCode32 = code & 0x7fff
						}
						if variant&16 != 0 {
							end.NetCode32 ^= 0x8000
							missing = true
						}
					}
				}
			}
			for step := 0; step < 2; step++ {
				c.srv.SetFrame(frame + uint32(step))
				beforeCalls := len(c.Calls)
				beforeOther := c.srv.Rand.Other.Index()
				blank := effectsPixelHash(pix)
				expired := op == 14 && mode == 0 && binary.LittleEndian.Uint32(data[433:]) == 1
				got := legacy.PortTestClientEffects(op, c.Viewport(), dr, [8]int32{}, nil)
				want := uint32(1)
				if expired {
					want = 0
				}
				if got != want {
					t.Fatal("ray draw lifetime")
				}
				if (missing || expired) && (effectsPixelHash(pix) != blank || len(c.Calls) != beforeCalls || c.srv.Rand.Other.Index() != beforeOther) {
					t.Fatal("unavailable ray drew/emitted particles")
				}
				if noxflags.HasGame(noxflags.GamePause) && op != 17 && len(c.Calls) != beforeCalls {
					t.Fatal("paused ray emitted particles")
				}
				if c.srv.Rand.Logic.Index() != 31 {
					t.Fatal("ray drawing changed logic RNG")
				}
				render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
				if op == 13 && variant == 8 && step == 0 {
					effectsDiagnostic(t, "energy-bolt", pix)
				}
				out = append(out, result{op, variant, step, got, effectsPixelHash(pix), render, env.Snapshot(), c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), append([]uint32(nil), c.Deleted...), c.MouseReads, c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				if got == 0 {
					break
				}
			}
		}
	}
	effectsCapture(t, "ray-drawing", out, len(out), "ca31d3dd068c703bc733f2993724aafb1b4e43c8130846734ef8be6f0dc5e44b")
}
