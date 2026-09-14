//go:build porttest

package opennox

import (
	"image"
	"math"
	"testing"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientParticleLightProperties(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	type result struct {
		Op, Variant, Constants int
		Params                 [4]int32
		Return                 int64
		Light                  []uint32
		Globals                []uint32
	}
	var out []result
	angles := []int32{-2147483648, -360, -1, 0, 1, 45, 90, 179, 180, 359, 360, 720, 32767, 2147483647}
	intensity := []float32{-10, 0, math.SmallestNonzeroFloat32, 0.5, 1, math.Nextafter32(1, 2), math.Nextafter32(31, 0), 31, math.Nextafter32(31, 32), math.Nextafter32(63, 0), 63, math.Nextafter32(63, 64), 100}
	constants := [][2]float64{{0.5, 65536}, {0, 65536}, {1.5, 1024}, {-0.5, 0}}
	for op := 12; op <= 16; op++ {
		for ci, coef := range constants {
			previousRadius := int64(-1)
			variants := len(angles)
			if op >= 15 {
				variants = len(intensity)
			}
			for variant := 0; variant < variants; variant++ {
				c.resetCase(effects, pix, 17, 120)
				env.Reset()
				env.Constants(coef[0], coef[1])
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
				light := unsafe.Slice((*uint32)(unsafe.Add(dr.C(), 136)), 9)
				for i := range light {
					light[i] = 0x13579abc + uint32(i*0x10001)
				}
				before := append([]uint32(nil), light...)
				p := [4]int32{angles[variant%len(angles)], angles[(variant+3)%len(angles)], angles[(variant+7)%len(angles)]}
				if op >= 15 {
					p[0] = int32(math.Float32bits(intensity[variant]))
				}
				got := legacy.PortTestClientDrawParticle(op, c.Viewport(), dr, p)
				allowed := map[int]bool{}
				switch op {
				case 12:
					allowed = map[int]bool{0: true, 4: true, 5: true, 6: true}
					if got != 1 || light[0] != 2 || light[4] != uint32(p[0]) || light[5] != uint32(p[1]) || light[6] != uint32(p[2]) {
						t.Fatal("light color words or return")
					}
				case 13:
					allowed = map[int]bool{7: true, 8: true}
					if uint16(light[7]) != uint16(got) || light[7]>>16 != before[7]>>16 || light[8] != 0 {
						t.Fatal("light direction layout")
					}
				case 14:
					allowed = map[int]bool{7: true}
					if uint16(light[7]>>16) != uint16(got) || uint16(light[7]) != uint16(before[7]) {
						t.Fatal("light penumbra layout")
					}
				case 15, 16:
					allowed = map[int]bool{1: true, 2: true}
					if op == 16 {
						allowed[3] = true
					}
					want := intensity[variant]
					if want > 63 {
						want = 63
					}
					if light[1] != math.Float32bits(want) || light[2] != uint32(got) || got < 0 {
						t.Fatal("light clamp/radius return")
					}
					if want <= 1 && got != 0 {
						t.Fatal("light minimum intensity")
					}
					if want > 1 && got == 0 {
						t.Fatal("positive light radius missing production coefficients")
					}
					if got < previousRadius {
						t.Fatal("light radius not monotonic")
					}
					previousRadius = got
				}
				if (op == 13 || op == 14) && (p[0] == 0 || coef[1] == 0) && got != int64(coef[0]) {
					t.Fatal("light zero-angle/scale conversion")
				}
				for i, v := range before {
					if !allowed[i] && light[i] != v {
						t.Fatalf("op%d changed light word%d outside contract", op, i)
					}
				}
				if c.Objs.Count != 1 || len(c.Calls) != 1 || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 {
					t.Fatal("light setter changed ownership/RNG")
				}
				out = append(out, result{op, variant, ci, p, got, append([]uint32(nil), light...), env.Snapshot()})
			}
		}
	}
	effectsCapture(t, "particle-light", out, len(out), "1ceb08e8d9bd30d4612500c837222766cd96ad3be88ebf2638eb6efa3754192b")
}

func TestClientParticlePalettes(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	type result struct {
		Seed   uint32
		Return int64
		Colors []uint32
	}
	var out []result
	for _, seed := range []uint32{0, 1, 0x12345678, 0xffffffff} {
		c.resetCase(effects, pix, 17, 120)
		env.Reset()
		colors := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1312500), 256)
		for i := range colors {
			colors[i] = seed + uint32(i)
		}
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
		got := legacy.PortTestClientDrawParticle(17, c.Viewport(), dr, [4]int32{})
		for _, sample := range []struct {
			index   int
			r, g, b byte
		}{
			{0, 0, 0, 0}, {63, 85, 153, 255}, {64, 0, 0, 0}, {95, 247, 247, 0},
			{96, 255, 255, 0}, {127, 255, 8, 0}, {128, 0, 0, 0}, {191, 255, 255, 255},
			{192, 0, 50, 50}, {255, 255, 50, 50},
		} {
			if colors[sample.index] != noxcolor.RGB5551Color(sample.r, sample.g, sample.b).Color32() {
				t.Fatalf("production palette sample%d", sample.index)
			}
		}
		if uint32(got) != colors[255] || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 || c.Objs.Count != 1 {
			t.Fatal("palette return/ownership/RNG")
		}
		out = append(out, result{seed, got, append([]uint32(nil), colors...)})
	}
	effectsCapture(t, "particle-palettes", out, len(out), "617ccd2159c45848fb40dc4d6732f87f6731e5118a65987b284612663dac2857")
}

func TestClientParticleColorInitialization(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	type result struct {
		Seed             uint32
		Step             int
		Return           int64
		Globals, Effects []uint32
	}
	var out []result
	names := []string{"dword_5d4594_1313532", "dword_5d4594_1313536", "dword_5d4594_1313540", "dword_5d4594_1313564"}
	for _, seed := range []uint32{0, 1, 0x12345678, 0xffffffff} {
		c.resetCase(effects, pix, 17, 120)
		env.Reset()
		colors := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1313524), 34)
		for i := range colors {
			colors[i] = seed + uint32(i)
		}
		for _, name := range names {
			*effects.Named(name) = seed
		}
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
		for step := 0; step < 3; step++ {
			got := legacy.PortTestClientDrawParticle(18, c.Viewport(), dr, [4]int32{})
			for _, sample := range []struct {
				off     uintptr
				r, g, b byte
			}{
				{1313524, 255, 255, 0}, {1313528, 255, 100, 0}, {1313544, 0, 200, 200},
				{1313548, 50, 255, 255}, {1313552, 255, 0, 255}, {1313556, 255, 200, 255},
				{1313560, 255, 200, 0}, {1313568, 100, 255, 50}, {1313572, 150, 255, 150},
				{1313576, 255, 255, 0}, {1313580, 0, 220, 0}, {1313584, 150, 255, 150},
				{1313588, 200, 200, 200}, {1313592, 255, 255, 255},
				{1313656, 255, 255, 255}, {1313644, 255, 255, 0}, {1313640, 255, 205, 0},
				{1313628, 255, 55, 0}, {1313624, 238, 55, 0}, {1313596, 118, 55, 0},
			} {
				if *memmap.PtrUint32(0x5D4594, sample.off) != noxcolor.RGB5551Color(sample.r, sample.g, sample.b).Color32() {
					t.Fatalf("color initializer sample%d", sample.off)
				}
			}
			for i, rgb := range [][3]byte{{255, 255, 0}, {0, 0, 255}, {0, 200, 255}, {255, 255, 100}} {
				if *effects.Named(names[i]) != noxcolor.RGB5551Color(rgb[0], rgb[1], rgb[2]).Color32() {
					t.Fatal("named color initializer")
				}
			}
			if uint32(got) != *memmap.PtrUint32(0x5D4594, 1313520) || len(c.Calls) != 1 || c.Objs.Count != 1 || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 {
				t.Fatal("color initializer return/ownership/RNG")
			}
			out = append(out, result{seed, step, got, env.Snapshot(), effects.Snapshot()})
		}
	}
	effectsCapture(t, "particle-color-init", out, len(out), "095076e2cab218dab5797bb2f73a02949f6a9ccf6fe2d6b5e50c03acf18021c9")
}
