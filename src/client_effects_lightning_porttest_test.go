//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientEffectsPlasmaMissingEndpoints(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	for missing := 0; missing < 3; missing++ {
		c.resetCase(env, pix, 31, 4)
		if missing == 1 {
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(20, 30))
			dr.NetCode32 = 101
		}
		if missing == 2 {
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(60, 70))
			dr.NetCode32 = 202
		}
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(14, image.Pt(48, 70))
		state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), 512)
		state[432] = 1
		binary.LittleEndian.PutUint32(state[437:], 101)
		binary.LittleEndian.PutUint32(state[441:], 202)
		before := append([]byte(nil), state...)
		scratch := env.Snapshot()
		blank := effectsPixelHash(pix)
		calls := len(c.Calls)
		got := legacy.PortTestClientEffects(17, c.Viewport(), dr, [8]int32{}, nil)
		if got != 1 || !bytes.Equal(before, state) || len(c.Calls) != calls || effectsPixelHash(pix) != blank || c.MouseReads != 1 || c.srv.Rand.Logic.Index() != 31 || c.srv.Rand.Other.Index() != 32 {
			t.Fatal("missing plasma endpoint must skip drawing without effect mutations")
		}
		after := env.Snapshot()
		if len(scratch) != len(after) {
			t.Fatal("scratch size")
		}
		for i, v := range scratch {
			if after[i] != v {
				t.Fatal("missing plasma endpoint changed scratch state")
			}
		}
	}
}

func TestClientEffectsLightningSetup(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Missing, Repeat int
		Return          uint32
		Globals         []uint32
	}
	var out []result
	names := []string{"BlueSpark", "YellowSpark", "GreenSpark"}
	for mask := 0; mask < 8; mask++ {
		c.resetCase(env, pix, 1, 0)
		saved := [3]int32{}
		for i, name := range names {
			p := c.Things.TypeByID(name)
			saved[i] = p.Field_1c
			if mask&(1<<i) != 0 {
				p.Field_1c = 0
			}
		}
		for repeat := 0; repeat < 2; repeat++ {
			got := legacy.PortTestClientEffects(45, c.Viewport(), nil, [8]int32{}, nil)
			if got != uint32(c.Things.IndByID("GreenSpark")) {
				t.Fatal("lightning setup return")
			}
			for i, off := range []uintptr{1316520, 1316524, 1316528} {
				if *memmap.PtrUint32(0x5D4594, off) != uint32(c.Things.IndByID(names[i])) {
					t.Fatal("lightning type cache")
				}
			}
			for _, v := range []struct {
				off     uintptr
				r, g, b byte
			}{{1316464, 255, 255, 255}, {1316488, 128, 128, 255}, {1316424, 128, 128, 255}, {1316428, 64, 64, 255}, {1316516, 200, 200, 255}, {1316512, 128, 128, 255}, {1316496, 255, 255, 255}, {1316468, 255, 255, 0}, {1316460, 30, 160, 30}, {1316444, 60, 140, 60}, {1316504, 40, 225, 40}, {1316480, 150, 220, 150}} {
				if *memmap.PtrUint32(0x5D4594, v.off) != noxcolor.RGB5551Color(v.r, v.g, v.b).Color32() {
					t.Fatal("lightning setup color")
				}
			}
			out = append(out, result{mask, repeat, got, env.Snapshot()})
		}
		for i, name := range names {
			c.Things.TypeByID(name).Field_1c = saved[i]
		}
	}
	effectsCapture(t, "lightning-setup", out, len(out), "9e3cc912b3fdef660c1f99c76c04555da2a167a0cff402c0af9de8ed5beed6ce")
}

func TestClientEffectsLightningRecursion(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Depth, Jitter, Geometry, Glow int
		Return                        uint32
		Pixels                        string
		Render, Globals               []uint32
		Logic, Other                  int
	}
	var out []result
	points := [][4]int16{{10, 10, 85, 85}, {85, 10, 10, 85}, {-10, 48, 110, 48}, {48, -10, 48, 110}, {-32768, -32768, 32767, 32767}, {0, 0, 0, 0}}
	for _, depth := range []int{1, 2, 4, 8} {
		for _, jitter := range []int{0, 1, 3, 100} {
			for pi, p := range points {
				for glow := 0; glow < 2; glow++ {
					c.resetCase(env, pix, 31, 0)
					*env.Named("dword_5d4594_1316448") = uint32(depth)
					*env.Named("dword_5d4594_1316476") = uint32(jitter)
					*env.Named("dword_5d4594_1316472") = 0xffff
					*memmap.PtrUint32(0x5D4594, 1316440) = 0x8421
					*memmap.PtrUint32(0x5D4594, 1316508) = uint32(glow)
					*memmap.PtrUint8(0x5D4594, 1316420) = byte(1 + pi)
					a := uint32(uint16(p[0])) | uint32(uint16(p[1]))<<16
					b := uint32(uint16(p[2])) | uint32(uint16(p[3]))<<16
					blank := effectsPixelHash(pix)
					got := legacy.PortTestClientEffects(9, c.Viewport(), nil, [8]int32{int32(a), int32(b)}, nil)
					if got != 0 || *env.Named("dword_5d4594_1316492") != 0 {
						t.Fatal("lightning recursion did not unwind")
					}
					rng := prand.New(32)
					for i := 0; i < 2*((1<<(depth-1))-1); i++ {
						rng.Int(0, 255)
					}
					if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != 31 {
						t.Fatal("lightning recursion RNG tree size")
					}
					hash := effectsPixelHash(pix)
					if depth == 1 && pi == 0 && hash == blank {
						t.Fatal("visible lightning did not render")
					}
					render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
					if depth == 4 && jitter == 3 && pi == 0 && glow == 1 {
						effectsDiagnostic(t, "lightning", pix)
					}
					out = append(out, result{depth, jitter, pi, glow, got, hash, render, env.Snapshot(), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				}
			}
		}
	}
	effectsCapture(t, "lightning-recursion", out, len(out), "3ca18302d14009f453d9149495ea713241dbcf8b5ed614eaa948cefa21d4ab53")
}

func TestClientEffectsLightningPasses(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Length, Mode, Passes int
		Return               uint32
		Pixels               string
		Render, Globals      []uint32
		Logic, Other         int
	}
	var out []result
	buf, free := alloc.Make([]byte{}, 48)
	defer free()
	for _, length := range []int{0, 1, 63, 64, 127, 128, 255, 256, 511, 512, 513} {
		for mode := 0; mode < 4; mode++ {
			for passes := 0; passes < 4; passes++ {
				c.resetCase(env, pix, 17, 0)
				*env.Named("dword_5d4594_1316452") = 0xfc01
				*env.Named("dword_5d4594_1316456") = 0x83e1
				*env.Named("dword_5d4594_1316436") = 0xffff
				*env.Named("dword_5d4594_1316484") = 0x8421
				*memmap.PtrUint8(0x5D4594, 1316420) = 5
				for i := range buf {
					buf[i] = 0xa5
				}
				data := buf[8:40]
				*(*[4]int32)(unsafe.Pointer(&data[0])) = [4]int32{48, 48, int32(48 + length), 48}
				*(*[4]uint16)(unsafe.Pointer(&data[16])) = [4]uint16{0x1234, 0x5678, 0x9abc, 0xdef0}
				before := append([]byte(nil), buf...)
				args := [8]int32{int32(mode), 7, int32(passes & 1), int32(passes >> 1)}
				got := legacy.PortTestClientEffects(10, c.Viewport(), nil, args, unsafe.Pointer(&data[0]))
				if got != uint32(passes>>1) || !bytes.Equal(before, buf) {
					t.Fatal("lightning passes return/input guards")
				}
				depth := 8
				if length < 64 {
					depth = 5
				} else if length < 128 {
					depth = 6
				} else if length < 256 {
					depth = 7
				}
				totalPasses := 2*(passes&1) + (passes >> 1)
				draws := totalPasses * 2 * ((1 << (depth - 2)) - 1)
				rng := prand.New(18)
				for i := 0; i < draws; i++ {
					rng.Int(0, 255)
				}
				if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != 17 {
					t.Fatal("lightning pass RNG count")
				}
				render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
				out = append(out, result{length, mode, passes, got, effectsPixelHash(pix), render, env.Snapshot(), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
			}
		}
	}
	effectsCapture(t, "lightning-passes", out, len(out), "470c70ea76b5fb1f0af0a8f6e74b02c30477107ac8c5c76f2758e2ab29c2aab6")
}
