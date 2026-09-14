//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestClientEffectsScreenParticles(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Capacity, Kind, Axis, Direction int
		Return                          uint32
		Particles                       [][]uint32
		Logic, Other                    int
	}
	var out []result
	for _, capacity := range []int{0, 1, 7, 128} {
		for kind := -1; kind <= 5; kind++ {
			for _, axis := range []int{0, 2} {
				for _, direction := range []int{0, 1} {
					c.resetCase(env, pix, 31, 0xfffffffc)
					snapshot, restore := legacy.PortTestEffectsScreenParticles(capacity)
					args := [8]int32{int32(kind), 10, 20, 30, 40, int32(axis), int32(direction)}
					got := legacy.PortTestClientEffects(6, c.Viewport(), nil, args, nil)
					particles := snapshot()
					restore()
					expected := 0
					if kind >= 0 && kind <= 4 {
						expected = capacity
						if expected > 100 {
							expected = 100
						}
					}
					if got != 0 || len(particles) != expected {
						t.Fatalf("screen kind%d capacity%d return%d count%d want%d", kind, capacity, got, len(particles), expected)
					}
					for _, p := range particles {
						if p[1] != uint32(kind) || p[8] != 1 || p[9] != 65536 || p[10] < 2 || p[10] > 5 {
							t.Fatal("screen particle kind/gravity/size")
						}
						colors := [][2]int{{3, 8}, {6, 0}, {7, 9}, {4, 2}, {1, 5}}[kind]
						if p[2] != uint32(colors[0]+1)*0x421 || p[3] != uint32(colors[1]+1)*0x421 {
							t.Fatal("screen particle colors")
						}
						x, y := int32(p[6])>>16, int32(p[7])>>16
						if axis == 2 {
							if x < 10 || x > 40 || y != 20 {
								t.Fatal("horizontal emission bounds")
							}
						} else {
							if x != 10 || y < 20 || y > 60 {
								t.Fatal("vertical emission bounds")
							}
						}
						vx, vy := int32(p[4])>>16, int32(p[5])>>16
						if vy < -40 || vy > -20 {
							t.Fatal("vertical particle velocity")
						}
						if direction == 1 {
							if vx < -20 || vx > 0 {
								t.Fatal("left particle velocity")
							}
						} else {
							if vx < 0 || vx > 20 {
								t.Fatal("right particle velocity")
							}
						}
					}
					rng := prand.New(32)
					for i := 0; i < 400; i++ {
						rng.Int(0, 255)
					}
					if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != 31 {
						t.Fatal("screen particle RNG consumption")
					}
					out = append(out, result{capacity, kind, axis, direction, got, particles, c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				}
			}
		}
	}
	effectsCapture(t, "screen-particles", out, len(out), "5c8be634483a7bf8d902f3526e78c74be5dcc01bbaa3a2364d1dc03200b27e46")
}

func TestClientEffectsParticleCallback(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	callback, hits, restore := legacy.PortTestEffectsParticleEnvironment()
	defer restore()
	type result struct {
		X, Y               uint32
		Callback           bool
		Return             uint32
		Drawable, Particle []uint32
		Calls              int
	}
	var out []result
	buf, free := alloc.Make([]byte{}, 144)
	defer free()
	for _, x := range []uint32{0, 1, 0xffff, 0x10000, 0x7fffffff, 0x80000000, 0xffff0000, 0xffffffff} {
		for _, y := range []uint32{0, 0x10000, 0x80000000, 0xffffffff} {
			for _, enabled := range []bool{false, true} {
				for i := range buf {
					buf[i] = 0xa5
				}
				words := unsafe.Slice((*uint32)(unsafe.Pointer(&buf[8])), 32)
				for i := range words {
					words[i] = uint32(i+7) * 0x13579
				}
				words[20], words[21], words[31] = x, y, 0
				if enabled {
					words[31] = uint32(uintptr(callback))
				}
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 70))
				state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
				binary.LittleEndian.PutUint32(state[432:], uint32(uintptr(unsafe.Pointer(&buf[8]))))
				before := append([]byte(nil), state...)
				expected := append([]byte(nil), buf...)
				if enabled {
					binary.LittleEndian.PutUint32(expected[8:], words[0]+1)
					binary.LittleEndian.PutUint32(expected[8+80:], x+0x10000)
					binary.LittleEndian.PutUint32(expected[8+84:], y^0x10000)
				}
				binary.LittleEndian.PutUint32(before[12:], x>>16)
				binary.LittleEndian.PutUint32(before[16:], y>>16)
				count := *hits
				got := legacy.PortTestClientEffects(32, nil, dr, [8]int32{}, nil)
				wantCalls := 0
				if enabled {
					wantCalls = 1
				}
				if got != 1 || *hits-count != wantCalls || !bytes.Equal(expected, buf) || !bytes.Equal(before, state) {
					t.Fatal("particle callback ordering/state/guards")
				}
				dw := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(dr)), 128)...)
				dw[108] = 1
				pw := append([]uint32(nil), words...)
				pw[31] = uint32(wantCalls)
				out = append(out, result{x, y, enabled, got, dw, pw, *hits - count})
				c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			}
		}
	}
	effectsCapture(t, "particle-callback", out, len(out), "60d9dcdc23eac4872cdeacdbb0b9bfc74e7c539ec1ad7315843984e5da4cca5b")
}
