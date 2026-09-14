//go:build porttest

package opennox

import (
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func effectsDistinctColors(env *legacy.PortTestEffectsEnvironment) {
	for off := uintptr(1313528); off <= 1313592; off += 4 {
		*memmap.PtrUint32(0x5D4594, off) = uint32(0x8001 + ((off-1313528)/4+1)*0x421)
	}
	for i, name := range []string{"dword_5d4594_1313532", "dword_5d4594_1313536", "dword_5d4594_1313540", "dword_5d4594_1313564", "nox_color_white_2523948"} {
		*env.Named(name) = uint32(0xfc01 - i*0x823)
	}
}

func TestClientEffectsGlowAndSparks(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Op, Variant, Step int
		Return            uint32
		Pixels            string
		Render, Globals   []uint32
		Drawables         [][]uint32
		Calls             []effectsSpawnCall
		Deleted           []uint32
		Logic, Other      int
	}
	var out []result
	for _, op := range []int{18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 33, 34, 36, 42, 43} {
		for variant := 0; variant < 32; variant++ {
			seed := uint32(1 + op*17 + variant*31)
			frame := []uint32{0, 1, 0xfffffffe, 128}[variant%4]
			c.resetCase(env, pix, seed, frame)
			effectsDistinctColors(env)
			if variant&8 != 0 {
				c.FailEvery = 2
			}
			typ := 4
			if op == 33 || op == 34 || op == 36 {
				typ = 1 + variant%7
			}
			if op == 23 {
				typ = c.Things.IndByID([]string{"RainOrbWhite", "RainOrbBlue"}[variant%2])
			}
			pos := []image.Point{{48, 70}, {9, 32}, {10, 32}, {85, 32}, {86, 32}, {48, 9}, {48, 10}, {48, 85}}[variant%8]
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
			if dr == nil {
				t.Fatal("parent drawable allocation")
			}
			state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
			age := []uint32{0, 1, 9, 10}[variant/4%4]
			begin, end := frame-age, frame-age+10
			dr.ZVal = []uint16{0, 1, 5, 35, 65535}[variant%5]
			dr.ZVal2 = []uint16{0, 3, 65535}[variant%3]
			dr.VelZ = []int8{-10, -1, 0, 1, 5, 127}[variant%6]
			dr.Field_74_4 = byte(variant * 17)
			binary.LittleEndian.PutUint32(state[432:], uint32(pos.X)<<12)
			binary.LittleEndian.PutUint32(state[436:], uint32(pos.Y)<<12)
			binary.LittleEndian.PutUint32(state[440:], uint32(1+variant*97))
			binary.LittleEndian.PutUint32(state[444:], begin)
			binary.LittleEndian.PutUint32(state[448:], end)
			if op == 23 {
				binary.LittleEndian.PutUint32(state[432:], uint32(pos.X-20))
				binary.LittleEndian.PutUint32(state[436:], uint32(pos.Y-30))
				binary.LittleEndian.PutUint16(state[440:], dr.ZVal+3)
				state[442] = byte(3 + variant%10)
			}
			if op == 33 || op == 34 || op == 36 {
				clear(state[432:452])
				target := pos.Add(image.Pt(30, -10))
				if variant%4 == 0 {
					target = pos
				}
				binary.LittleEndian.PutUint16(state[432:], uint16(target.X))
				binary.LittleEndian.PutUint16(state[434:], uint16(target.Y))
				state[443] = byte(1 + variant%12)
				state[444] = []byte{1, 2, 5, 12, 31}[variant%5]
				state[445] = []byte{0, 1, 3}[variant%3]
				state[446] = []byte{0, 1, 2}[variant%3]
			}
			for step := 0; step < 4; step++ {
				f := frame + uint32(step)
				c.srv.SetFrame(f)
				args := [8]int32{0xfc01, 0x8421}
				if op == 36 {
					args[0] = int32(variant % 2)
				}
				beforeOther := c.srv.Rand.Other.Index()
				beforeCalls := len(c.Calls)
				oldHeight := int16(dr.ZVal)
				expectedOther := 0
				if op == 18 || op == 19 || op == 21 {
					remaining := int32(end - f)
					if remaining == int32(end-begin) {
						remaining--
					}
					screen := image.Pt(dr.PosVec.X, dr.PosVec.Y-int(int16(dr.ZVal))-int(int16(dr.ZVal2)))
					if remaining > 0 && screen.X >= 10 && screen.Y >= 10 && screen.X+10 < 96 && screen.Y+10 < 96 {
						expectedOther++
					}
					if op != 18 {
						expectedOther++
					}
				}
				if op == 20 {
					expectedOther = 1
				}
				got := legacy.PortTestClientEffects(op, c.Viewport(), dr, args, nil)
				if got > 1 {
					t.Fatal("draw callback return outside live/deleted contract")
				}
				childSuccess := 0
				for _, call := range c.Calls[beforeCalls:] {
					if call.Ref != 0 {
						childSuccess++
					}
				}
				if op == 22 {
					expectedOther = 5 * childSuccess
				}
				if op == 23 && oldHeight <= 0 {
					expectedOther = 1 + childSuccess
				}
				rng := prand.New(beforeOther)
				for i := 0; i < expectedOther; i++ {
					rng.Int(0, 255)
				}
				if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != prand.New(int(seed)).Index() {
					t.Fatalf("glow op%d variant%d step%d RNG other%d want%d", op, variant, step, c.srv.Rand.Other.Index(), rng.Index())
				}
				render := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)
				out = append(out, result{op, variant, step, got, effectsPixelHash(pix), render, env.Snapshot(), c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), append([]uint32(nil), c.Deleted...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				if got == 0 {
					break
				}
			}
		}
	}
	effectsCapture(t, "glow-sparks", out, len(out), "d0e724ab19068881ccc0f68f43e96ddc4b6ea970c4c27c8f96d4afc7346e66db")
}
