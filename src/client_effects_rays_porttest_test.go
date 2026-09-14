//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientEffectsRayDispatch(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Code                         byte
		Count, Seed, Failure, Coords int
		Return                       uint32
		Calls                        []effectsSpawnCall
		Drawables                    [][]uint32
		Cache, Globals               []uint32
		Logic, Other                 int
	}
	var out []result
	packet, free := alloc.Make([]byte{}, 25)
	defer free()
	types := map[byte]string{125: "PlasmaRay", 140: "DynamicLightning", 141: "DynamicEnergyBolt", 142: "DynamicChainLightning", 143: "OrbRay", 144: "OrbRay", 145: "OrbRay"}
	for _, code := range []byte{0, 124, 125, 126, 140, 141, 142, 143, 144, 145, 255} {
		for _, count := range []int{0, 95, 96} {
			for _, seed := range []int{1, 1023} {
				for _, failure := range []int{0, 1, 2} {
					for ci, coords := range [][4]uint16{{48, 70, 64, 90}, {65535, 0, 0, 65535}} {
						c.resetCase(env, pix, uint32(seed), 127)
						cache := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1303540), 203)
						for i := 0; i < count; i++ {
							dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48+i, 70))
							cache[i] = uint32(uintptr(unsafe.Pointer(dr)))
						}
						*memmap.PtrUint32(0x5D4594, 1304308) = uint32(count)
						c.Calls = nil
						c.FailEvery = failure
						for i := range packet {
							packet[i] = 0xa5
						}
						body := packet[8:17]
						body[0] = code
						for i, v := range coords {
							binary.LittleEndian.PutUint16(body[1+i*2:], v)
						}
						before := append([]byte(nil), packet...)
						got := legacy.PortTestClientEffects(8, c.Viewport(), nil, [8]int32{}, unsafe.Pointer(&body[0]))
						if !bytes.Equal(before, packet) {
							t.Fatal("ray decoder changed guarded packet")
						}
						expectedCount := count
						expectedRNG := 0
						if count >= 96 {
							if got != uint32(count) || len(c.Calls) != 0 {
								t.Fatal("ray limit gate")
							}
						} else if name, known := types[code]; !known {
							if got != uint32(int32(code)-125) || len(c.Calls) != 0 {
								t.Fatal("unknown ray code")
							}
						} else {
							if len(c.Calls) == 0 {
								t.Fatal("missing main ray creation")
							}
							main := c.Calls[len(c.Calls)-1]
							midpoint := image.Pt(int(coords[0])+(int(coords[2])-int(coords[0]))/2, int(coords[1])+(int(coords[3])-int(coords[1]))/2)
							if main.Type != c.Things.IndByID(name) || main.Position != midpoint {
								t.Fatal("ray type/midpoint")
							}
							if main.Ref != 0 {
								dr := (*client.Drawable)(unsafe.Pointer(uintptr(got)))
								if c.refs[dr] != main.Ref {
									t.Fatal("ray return pointer")
								}
								state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), 512)
								if state[432] != 0 || !bytes.Equal(state[437:445], body[1:9]) {
									t.Fatal("ray endpoint payload")
								}
								expectedCount++
								if cache[count] != got {
									t.Fatal("ray cache append")
								}
								got = main.Ref
							} else if got != 0 {
								t.Fatal("failed ray allocation return")
							}
							if code == 143 || code == 145 {
								expectedRNG = 2
							}
							if code == 144 {
								expectedRNG = 3
							}
							expectedRNG += 2 * (len(c.Calls) - 1)
							for _, call := range c.Calls[:len(c.Calls)-1] {
								if call.Ref != 0 {
									expectedRNG++
								}
							}
						}
						if *memmap.PtrUint32(0x5D4594, 1304308) != uint32(expectedCount) {
							t.Fatal("ray cache count")
						}
						rng := prand.New(seed + 1)
						for i := 0; i < expectedRNG; i++ {
							rng.Int(0, 255)
						}
						if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != prand.New(seed).Index() {
							t.Fatal("ray RNG consumption")
						}
						normalized := append([]uint32(nil), cache...)
						for i, v := range normalized[:96] {
							if v != 0 {
								ref := c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(v)))]
								if ref == 0 {
									t.Fatal("unknown cached ray")
								}
								normalized[i] = ref
							}
						}
						out = append(out, result{code, count, seed, failure, ci, got, append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), normalized, env.Snapshot(), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
					}
				}
			}
		}
	}
	effectsCapture(t, "ray-dispatch", out, len(out), "3a2a1c704788a2fb1571665b0364e0e3be16116811407124d13ce9ac837a231e")
}

func TestClientEffectsOrbitUpdate(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Variant            int
		Age, Frame, Return uint32
		Drawables          [][]uint32
		Deleted            []uint32
		Logic, Other       int
	}
	var out []result
	for variant := 0; variant < 32; variant++ {
		for _, age := range []uint32{0, 1, 29, 30, 59, 60, 61} {
			frame := []uint32{0, 1, 127, 0xfffffffe}[variant%4]
			c.resetCase(env, pix, 1, frame)
			pos := image.Pt(48, 70)
			target := image.Pt(148, 170)
			if variant%4 == 0 {
				target = pos
			}
			if variant%4 == 1 {
				target = pos.Add(image.Pt(9, 9))
			}
			if variant%4 == 2 {
				target = pos.Add(image.Pt(10, 0))
			}
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
			dr.AnimStart = frame - age
			data := unsafe.Slice((*byte)(unsafe.Pointer(dr)), 512)
			radius := []uint16{0, 1, 31, 255, 1024, 32767, 65535}[variant%7]
			binary.LittleEndian.PutUint16(data[432:], uint16(target.X))
			binary.LittleEndian.PutUint16(data[434:], uint16(target.Y))
			binary.LittleEndian.PutUint16(data[440:], radius)
			data[442] = []byte{0, 1, 63, 64, 127, 128, 255}[variant%7]
			data[443] = byte(variant % 2)
			got := legacy.PortTestClientEffects(37, c.Viewport(), dr, [8]int32{}, nil)
			want := uint32(1)
			if age >= 60 || variant%4 < 2 {
				want = 0
			}
			if got != want || c.Objs.Count != int(want) {
				t.Fatal("orbit age/proximity deletion")
			}
			if want == 1 && age == 0 && variant%7 == 0 && dr.PosVec != target {
				t.Fatal("zero-radius orbit position")
			}
			if c.srv.Rand.Logic.Index() != 1 || c.srv.Rand.Other.Index() != 2 {
				t.Fatal("orbit update consumed RNG")
			}
			out = append(out, result{variant, age, frame, got, c.snapshotDrawables(t), append([]uint32(nil), c.Deleted...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
		}
	}
	effectsCapture(t, "orbit-update", out, len(out), "8cde798a2aa9871c93f27b67336d428a374a9c7c95dad0611e7f00c2083bc56b")
}
