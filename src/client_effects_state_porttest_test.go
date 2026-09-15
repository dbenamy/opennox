//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type effectsSpawnCall struct {
	Type     int
	Position image.Point
	Ref      uint32
}
type effectsTestClient struct {
	callbackRefs map[unsafe.Pointer]uint32
	imageRefs    map[uint32]uint32
	dataRefs     map[uint32]uint32
	Mouse        image.Point
	MouseReads   int
	*Client
	refs      map[*client.Drawable]uint32
	next      uint32
	FailEvery int
	Calls     []effectsSpawnCall
	Deleted   []uint32
}

func (c *effectsTestClient) GetMousePos() image.Point { c.MouseReads++; return c.Mouse }

func (c *effectsTestClient) Nox_xxx_spriteLoadAdd_45A360_drawable(typ int, pos image.Point) *client.Drawable {
	row := effectsSpawnCall{Type: typ, Position: pos}
	var dr *client.Drawable
	if c.FailEvery == 0 || (len(c.Calls)+1)%c.FailEvery != 0 {
		dr = c.Client.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
		if dr != nil {
			c.next++
			c.refs[dr] = c.next
			row.Ref = c.next
		}
	}
	c.Calls = append(c.Calls, row)
	return dr
}
func (c *effectsTestClient) Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr *client.Drawable) {
	c.Deleted = append(c.Deleted, c.refs[dr])
	c.Client.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
}
func newEffectsFullOwner(t *testing.T, extraNames ...string) (*effectsTestClient, *noximage.Image16, *legacy.PortTestEffectsEnvironment) {
	t.Helper()
	base, pix := newEffectsTestOwner(t, extraNames...)
	base.srv.Rand.Logic, base.srv.Rand.Other = prand.New(1), prand.New(2)
	oldHover := legacy.Get_dword_5d4594_1096640()
	oldCursor := legacy.Get_nox_client_spriteUnderCursorXxx_1096644()
	legacy.Set_dword_5d4594_1096640(nil)
	legacy.Set_nox_client_spriteUnderCursorXxx_1096644(nil)
	t.Cleanup(func() {
		legacy.Set_dword_5d4594_1096640(oldHover)
		legacy.Set_nox_client_spriteUnderCursorXxx_1096644(oldCursor)
	})
	oldServer := legacy.GetServer
	legacy.GetServer = func() legacy.Server { return base.srv }
	t.Cleanup(func() { legacy.GetServer = oldServer })
	c := &effectsTestClient{Client: base, refs: make(map[*client.Drawable]uint32)}
	// newEffectsTestOwner already registered restoration of the original client.
	legacy.GetClient = func() legacy.Client { return c }
	env := legacy.PortTestNewEffectsEnvironment()
	t.Cleanup(env.Restore)
	return c, pix, env
}
func effectsCapture(t *testing.T, label string, out any, count int, want string) {
	t.Helper()
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_CLIENT_EFFECTS_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d results %s", label, count, hash)
	if hash != want {
		t.Fatalf("%s hash %s want audited original-C capture %s", label, hash, want)
	}
}
func TestClientEffectsSparkBounce(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	type result struct {
		Z            uint16
		Velocity     int8
		Frame        uint32
		Step         int
		Return       uint32
		State        []uint32
		Logic, Other int
	}
	var out []result
	for _, z := range []uint16{0, 1, 2, 10, 35, 32767, 32768, 65535} {
		for _, velocity := range []int8{-128, -127, -10, -3, -2, -1, 0, 1, 2, 3, 10, 127} {
			for _, frame := range []uint32{0, 1, 2, 0x7fffffff, 0xffffffff} {
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 70))
				if dr == nil {
					t.Fatal("drawable allocation")
				}
				dr.ZVal, dr.VelZ = z, velocity
				state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
				for step := 0; step < 4; step++ {
					f := frame + uint32(step)
					c.srv.SetFrame(f)
					before := append([]byte(nil), state...)
					nextZ := int16(uint16(int(dr.ZVal) + int(dr.VelZ)))
					nextVelocity := dr.VelZ
					ret := nextZ
					if nextZ >= 0 {
						if f&1 != 0 {
							nextVelocity--
						}
					} else {
						nextZ = -nextZ
						ret = int16(26209 * int(dr.VelZ))
						nextVelocity = int8(-9 * int(dr.VelZ) / 10)
						if nextVelocity < 2 {
							nextZ, nextVelocity = 0, 0
						}
					}
					binary.LittleEndian.PutUint16(before[104:], uint16(nextZ))
					before[296] = byte(nextVelocity)
					got := legacy.PortTestClientEffects(44, c.Viewport(), dr, [8]int32{}, nil)
					if got != uint32(int32(ret)) || !bytes.Equal(before, state) || c.srv.Rand.Logic.Index() != 1 || c.srv.Rand.Other.Index() != 2 {
						t.Fatalf("bounce z%d v%d frame%d step%d return%d want%d", z, velocity, frame, step, got, ret)
					}
					words := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(dr)), 128)...)
					out = append(out, result{z, velocity, frame, step, got, words, c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				}
				c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			}
		}
	}
	effectsCapture(t, "spark-bounce", out, len(out), "75cbc70715d55b4cb395b7043e7451dae827c1a96082100260ac34bf18df350a")
}

func (c *effectsTestClient) resetCase(env *legacy.PortTestEffectsEnvironment, pix *noximage.Image16, seed uint32, frame uint32) {
	for c.Objs.List1 != nil {
		c.Client.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(c.Objs.List1)
	}
	c.refs = make(map[*client.Drawable]uint32)
	c.next = 0
	c.FailEvery = 0
	c.Calls = nil
	c.Deleted = nil
	c.Mouse = image.Pt(48, 48)
	c.MouseReads = 0
	c.srv.SetFrame(frame)
	c.srv.Rand.Logic, c.srv.Rand.Other = prand.New(int(seed)), prand.New(int(seed+1))
	env.Reset()
	clear(pix.Pix)
	d := c.r.Data()
	d.Reset()
	d.SetClip(true)
	d.SetClipRect(pix.Rect)
	d.SetClipRect2(image.Rect(0, 0, 95, 95))
	d.SetRect3(pix.Rect)
	c.r.ClearPoints()
}
func (c *effectsTestClient) snapshotDrawables(t *testing.T) [][]uint32 {
	t.Helper()
	ref := func(p uint32) uint32 {
		if p == 0 {
			return 0
		}
		n := c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(p)))]
		if n == 0 {
			t.Fatal("unknown drawable pointer in captured list")
		}
		return n
	}
	var out [][]uint32
	var prev *client.Drawable
	for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
		if len(out) >= 512 || dr.Field_93 != prev {
			t.Fatal("invalid production drawable list")
		}
		words := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(dr)), 128)...)
		for _, i := range []int{83, 84, 87, 88, 90, 91, 92, 93, 94, 95, 97, 98, 100, 101, 102, 103, 104, 105, 106, 107} {
			words[i] = ref(words[i])
		}
		for _, i := range []int{75, 115} {
			if words[i] != 0 {
				marker, ok := c.callbackRefs[unsafe.Pointer(uintptr(words[i]))]
				if !ok {
					t.Fatalf("unmodeled effect callback at word%d", i)
				}
				words[i] = marker
			}
		}
		for i, refs := range map[int]map[uint32]uint32{2: c.imageRefs, 76: c.dataRefs} {
			if words[i] != 0 {
				marker, ok := refs[words[i]]
				if !ok {
					t.Fatalf("unmodeled owned image/data pointer at word%d", i)
				}
				words[i] = marker
			}
		}
		for _, i := range []int{99, 114, 124} {
			if words[i] != 0 {
				t.Fatalf("unmodeled pointer in drawable word%d", i)
			}
		}
		if words[116] != 0 {
			if unsafe.Pointer(uintptr(words[116])) != legacy.PortTestEffectsCallback(37) {
				t.Fatal("unexpected drawable update callback")
			}
			words[116] = 0xe1000025
		}
		out = append(out, append([]uint32{c.refs[dr]}, words...))
		prev = dr
	}
	if len(out) != c.Objs.Count {
		t.Fatal("drawable count differs from list")
	}
	// Check actual spatial ownership, including both directions of bucket links.
	indexed := make(map[*client.Drawable]bool)
	for x, col := range c.Objs.Index2D {
		for y, head := range col {
			var previous *client.Drawable
			for dr := head; dr != nil; dr = dr.Field_100 {
				if indexed[dr] || dr.Field_101 != previous || dr.PosVec.X/128 != x || dr.PosVec.Y/128 != y {
					t.Fatal("invalid drawable spatial index")
				}
				indexed[dr] = true
				previous = dr
			}
		}
	}
	return out
}

func TestClientEffectsParticleCreation(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Op, Variant, Failure int
		Seed, Frame, Return  uint32
		Calls                []effectsSpawnCall
		Drawables            [][]uint32
		Logic, Other         int
	}
	var out []result
	data, free := alloc.Make([]byte{}, 48)
	defer free()
	for _, op := range []int{0, 1, 2, 3, 4, 5, 7} {
		for variant := 0; variant < 8; variant++ {
			for _, failure := range []int{0, 1, 2, 3} {
				seed := uint32(1 + variant*17 + op*131)
				frame := []uint32{0, 1, 127, 0xfffffffc}[variant%4]
				c.resetCase(env, pix, seed, frame)
				c.FailEvery = failure
				for i := range data {
					data[i] = 0xa5
				}
				buf := data[8:40]
				clear(buf)
				args := [8]int32{}
				attempts := 1
				randomPerSuccess := 1
				ptrResult := false
				switch op {
				case 0, 1:
					coords := [][4]uint16{{48, 70, 64, 90}, {0, 0, 0, 0}, {65535, 1, 2, 65534}, {32767, 32767, 32760, 32760}}[variant%4]
					for i, v := range coords {
						binary.LittleEndian.PutUint16(buf[2*i:], v)
					}
					if op == 0 {
						args = [8]int32{4, int32(variant - 4), int32(4 - variant), int32(variant * 43), int32(variant * 63)}
					} else {
						args = [8]int32{4, int32(variant * 10001), int32(variant % 2), int32(variant * 63)}
					}
				case 2:
					count := []int32{-1, 0, 1, 2, 3, 7, 12, 20}[variant]
					args = [8]int32{4, count, 1 + int32(variant*100), int32(variant * 8), 48, 70}
					attempts = int(count)
					if attempts < 0 {
						attempts = 0
					}
					randomPerSuccess = 4
				case 3:
					args = [8]int32{48, 70, int32(variant * 10001), 4}
					attempts = 2
					randomPerSuccess = 5
				case 4:
					*(*[4]int32)(unsafe.Pointer(&buf[0])) = [4]int32{48, 70, int32(variant * 13), int32(variant * 17)}
					args = [8]int32{4, int32(variant * 10001), int32(variant * 43)}
					ptrResult = true
				case 5:
					*(*[4]int32)(unsafe.Pointer(&buf[0])) = [4]int32{48, 70, 48 + int32(variant*13), 70 + int32(variant*17)}
					args[0] = 4
					attempts = -1
					randomPerSuccess = 5
				case 7:
					*(*[2]int32)(unsafe.Pointer(&buf[0])) = [2]int32{48, 70}
					strength := []int32{0, 1, 2, 63, 127, 128, 254, 255}[variant]
					args = [8]int32{4, strength}
					attempts = 180*int(strength)/255 + 10
					randomPerSuccess = 5
				}
				input := append([]byte(nil), data...)
				got := legacy.PortTestClientEffects(op, c.Viewport(), nil, args, unsafe.Pointer(&buf[0]))
				if !bytes.Equal(input, data) {
					t.Fatal("particle constructor changed guarded inputs")
				}
				if attempts >= 0 && len(c.Calls) != attempts {
					t.Fatalf("op%d variant%d attempts%d want%d", op, variant, len(c.Calls), attempts)
				}
				successes := 0
				for i, call := range c.Calls {
					if call.Type != 4 {
						t.Fatal("wrong particle type")
					}
					wantPos := image.Pt(48, 70)
					if op == 0 {
						wantPos = image.Pt(int(binary.LittleEndian.Uint16(buf[4:])), int(binary.LittleEndian.Uint16(buf[6:]))).Add(image.Pt(int(args[1]), int(args[2])))
					}
					if op == 1 {
						wantPos = image.Pt(int(int16(binary.LittleEndian.Uint16(buf[4:]))), int(int16(binary.LittleEndian.Uint16(buf[6:]))))
					}
					if op != 5 && call.Position != wantPos {
						t.Fatal("particle spawn position")
					}
					if op == 5 {
						if call.Position.X < 48 || call.Position.X > 48+variant*13 || call.Position.Y < 70 || call.Position.Y > 70+variant*17 {
							t.Fatal("line particle outside segment bounds")
						}
						if i > 0 && (call.Position.X < c.Calls[i-1].Position.X || call.Position.Y < c.Calls[i-1].Position.Y) {
							t.Fatal("line particles reversed")
						}
					}
					success := failure == 0 || (i+1)%failure != 0
					if (call.Ref != 0) != success {
						t.Fatal("allocation failure handling")
					}
					if success {
						successes++
					}
				}
				if c.Objs.Count != successes {
					t.Fatal("particle count")
				}
				draws := c.snapshotDrawables(t)
				rng := prand.New(int(seed + 1))
				randomCount := successes * randomPerSuccess
				if op == 5 && variant != 0 {
					randomCount += 1 + len(c.Calls)
				}
				for i := 0; i < randomCount; i++ {
					rng.Int(0, 255)
				}
				if c.srv.Rand.Other.Index() != rng.Index() || c.srv.Rand.Logic.Index() != prand.New(int(seed)).Index() {
					t.Fatalf("op%d variant%d RNG consumption other%d want%d", op, variant, c.srv.Rand.Other.Index(), rng.Index())
				}
				if ptrResult {
					if successes == 0 {
						if got != 0 {
							t.Fatal("failed creation returned a drawable")
						}
					} else {
						if c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(got)))] == 0 {
							t.Fatal("creation returned an unknown drawable")
						}
						got = c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(got)))]
					}
				}
				if op == 2 {
					want := args[1]
					if want > 0 {
						want = 0
					}
					if got != uint32(want) {
						t.Fatal("point particle count return")
					}
				}
				if (op == 3 || op == 7) && got != 0 {
					t.Fatal("particle loop return")
				}
				out = append(out, result{op, variant, failure, seed, frame, got, append([]effectsSpawnCall(nil), c.Calls...), draws, c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
			}
		}
	}
	effectsCapture(t, "particle-creation", out, len(out), "3c01f1d25bf86290df7c4fc126ce8b9241fe1e59a17b133458f0ca2e710cdb1e")
}
