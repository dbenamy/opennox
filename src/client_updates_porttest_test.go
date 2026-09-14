//go:build porttest

package opennox

import (
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

var updateTestTypes = []string{"VioletSpark", "DeathBallSpark", "GreenPuff", "GreenSmoke", "WhiteVortexOrb", "Spark", "Puff", "MagicMissileTailLink", "MagicTailLink"}

func newUpdateTestOwner(t *testing.T) (*effectsTestClient, *noximage.Image16, *legacy.PortTestEffectsEnvironment, *legacy.PortTestClientUpdateEnvironment) {
	t.Helper()
	c, pix, effects := newEffectsFullOwner(t, updateTestTypes...)
	c.callbackRefs = map[unsafe.Pointer]uint32{
		legacy.PortTestEffectsCallback(19):         0xe1000013,
		legacy.PortTestEffectsCallback(21):         0xe1000015,
		legacy.PortTestClientUpdateCloudCallback(): 0xe1000100,
	}
	env := legacy.PortTestNewClientUpdateEnvironment()
	t.Cleanup(env.Restore)
	flags := noxflags.GetGame()
	t.Cleanup(func() { noxflags.ResetGame(); noxflags.SetGame(flags) })
	return c, pix, effects, env
}

func TestClientUpdatesCreation(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	type result struct {
		Op, Variant, Failure int
		Frame, Return        uint32
		Params               [4]int32
		Calls                []effectsSpawnCall
		Drawables            [][]uint32
		Globals              []uint32
		Logic, Other         int
	}
	var out []result
	for _, op := range []int{1, 2, 3, 5, 7, 8, 9, 10, 11, 12, 13, 16, 17, 18, 19, 20, 22, 23, 24, 26} {
		for variant := 0; variant < 32; variant++ {
			for _, failure := range []int{0, 1, 2} {
				frame := []uint32{0, 1, 9, 10, 51, 127, 0xfffffffe, 0xffffffff}[variant%8]
				seed := uint32(1 + 31*variant)
				c.resetCase(effects, pix, seed, frame)
				density := []uint32{0, 1, 3, 8}[variant/8]
				env.Reset(density)
				noxflags.ResetGame()
				if variant&8 != 0 {
					noxflags.SetGame(noxflags.GamePause)
				}
				fps := []uint32{30, 60, 90, 120}[variant/8]
				c.srv.SetTickRate(fps)
				radius := []float64{0, 1, 64.75, 255.99}[variant/8]
				restore := c.srv.PortTestClientUpdateBalance(radius)
				// Register cleanup as well so assertions cannot leave a balance overlay.
				t.Cleanup(restore)
				pos := image.Pt(320, 448)
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
				delta := []image.Point{image.Pt(0, 0), image.Pt(6, 0), image.Pt(7, 0), image.Pt(14, 1), image.Pt(15, 0), image.Pt(-21, 35), image.Pt(64, -32), image.Pt(-80, -48)}[variant%8]
				dr.Field_8 = uint32(pos.X - delta.X)
				dr.Field_9 = uint32(pos.Y - delta.Y)
				dr.ZVal = []uint16{0, 1, 32767, 65535}[variant/8]
				dr.ZVal2 = uint16(variant)
				dr.Field_80 = frame - []uint32{0, 9, 10, 11}[variant/8]
				if variant&16 != 0 {
					dr.Field_120 = 1
				}
				data := unsafe.Slice((*byte)(dr.C()), 512)
				anchor := pos.Sub(delta)
				binary.LittleEndian.PutUint32(data[432:], uint32(anchor.X))
				binary.LittleEndian.PutUint32(data[436:], uint32(anchor.Y))
				p := [4]int32{[]int32{-1, 0, 1, 3}[variant/8], int32([]int{0, 1, 35, 75}[variant/8])}
				if op == 5 {
					p[0] = int32(variant/8 + 1)
				}
				before := len(c.Calls)
				beforeOther := c.srv.Rand.Other.Index()
				c.FailEvery = failure
				got := legacy.PortTestClientUpdate(op, c.Viewport(), dr, p)
				calls := append([]effectsSpawnCall(nil), c.Calls[before:]...)
				successes := 0
				for _, call := range calls {
					if call.Ref != 0 {
						successes++
					}
				}
				count := int(p[0])
				if count < 0 {
					count = 0
				}
				dist := delta.X
				if dist < 0 {
					dist = -dist
				}
				dy := delta.Y
				if dy < 0 {
					dy = -dy
				}
				dist += dy
				links := 0
				if delta.X*delta.X+delta.Y*delta.Y > 200 {
					links = 1
				}
				attempts, wantRNG := -1, -1
				switch op {
				case 1:
					attempts = count
					wantRNG = 3*attempts + 2*successes
				case 2:
					attempts = count
					wantRNG = 4 * successes
				case 3:
					wantRNG = 30 + 3*len(calls) + successes
				case 5:
					attempts = dist / 7
					wantRNG = attempts + 4*successes
				case 7:
					attempts = 20
					if frame&1 != 0 && frame-dr.Field_80 < 10 {
						for angle := frame % 51; angle < 256; angle += 51 {
							attempts += 2
						}
					}
					wantRNG = 40
					for i, call := range calls {
						if call.Ref != 0 {
							if i < 20 {
								wantRNG += 2
							} else {
								wantRNG++
							}
						}
					}
				case 8:
					attempts = int(density) + links
					wantRNG = 2 * int(density)
					for i, call := range calls {
						if i < int(density) && call.Ref != 0 {
							wantRNG += 3
						}
					}
				case 9:
					attempts = 4 + links
					wantRNG = 8
					for i, call := range calls {
						if i >= links && call.Ref != 0 {
							wantRNG++
						}
					}
				case 10:
					attempts = 5
					wantRNG = 10 + successes
				case 11:
					attempts = 1
					wantRNG = 2 + 4*successes
				case 12:
					attempts = 1
					wantRNG = 1 + 2*successes
				case 13:
					attempts = 0
					wantRNG = 0
				case 16, 17, 18, 19, 20:
					attempts = 0
					if dr.Field_120 == 0 && !noxflags.HasGame(noxflags.GamePause) {
						attempts = dist / 7
					}
					wantRNG = attempts + 4*successes
				case 22:
					attempts = 3
					wantRNG = 4 * successes
				case 23:
					attempts = 1
					wantRNG = 4 * successes
				case 24, 26:
					attempts = int(frame & 1)
					wantRNG = 3*attempts + 2*successes
				}
				if attempts >= 0 && len(calls) != attempts {
					t.Fatalf("op%d v%d failure%d attempts%d want%d", op, variant, failure, len(calls), attempts)
				}
				if c.srv.Rand.Other.Index()-beforeOther != wantRNG || c.srv.Rand.Logic.Index() != int(seed) {
					t.Fatalf("op%d v%d failure%d RNG", op, variant, failure)
				}
				if c.Objs.Count != 1+successes || len(c.Deleted) != 0 {
					t.Fatal("update construction ownership")
				}
				if failure == 1 && successes != 0 {
					t.Fatal("allocation-disabled update created a drawable")
				}
				if op == 2 && p[0] > 0 && got != 0 {
					id := c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(got)))]
					if id == 0 {
						t.Fatal("unowned last-created result")
					}
					got = id
				}
				wantReturn := uint32(1)
				switch op {
				case 1:
					wantReturn = 0
					if p[0] <= 0 {
						wantReturn = uint32(p[0])
					}
				case 2:
					if p[0] <= 0 {
						wantReturn = uint32(p[0])
					} else {
						wantReturn = calls[len(calls)-1].Ref
					}
				case 5, 9:
					wantReturn = 0
				}
				if got != wantReturn {
					t.Fatalf("op%d returned%d want%d", op, got, wantReturn)
				}
				out = append(out, result{op, variant, failure, frame, got, p, calls, c.snapshotDrawables(t), env.Snapshot(), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
				restore()
			}
		}
	}
	effectsCapture(t, "update-creation", out, len(out), "eeb05a51bc4138c6511ae3c1bd67d8bd8c2dad802d151f04f187dced6cfe023a")
}

// Allocation failure must leave the last trail anchor available for retry.
// Independent spark emission is observable even when a magic tail is absent.
func TestClientUpdatesTailAllocationRetry(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	for _, op := range []int{8, 9} {
		c.resetCase(effects, pix, 17, 120)
		env.Reset(0) // Isolate the missile tail from its configurable sparks.
		c.srv.SetTickRate(60)
		parent := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(320, 448))
		words := unsafe.Slice((*uint32)(parent.C()), 128)
		words[108], words[109] = 280, 400
		parent.Field_8, parent.Field_9 = 280, 400
		c.FailEvery = 2 // Parent is call one; tail fails, some independent sparks succeed.
		before := len(c.Calls)
		legacy.PortTestClientUpdate(op, c.Viewport(), parent, [4]int32{})
		attempts := 1
		if op == 9 {
			attempts = 5
		}
		if len(c.Calls)-before != attempts || c.Calls[before].Ref != 0 {
			t.Fatalf("op%d failed-tail attempts", op)
		}
		if words[108] != 280 || words[109] != 400 || c.Objs.DeadlineList != nil || parent.Deadline != 0 {
			t.Fatalf("op%d failure advanced anchor or registered decay", op)
		}
		if op == 9 && c.Objs.Count != 3 {
			t.Fatal("magic tail failure suppressed independent sparks")
		}
		c.snapshotDrawables(t)
		c.FailEvery = 0
		before = len(c.Calls)
		legacy.PortTestClientUpdate(op, c.Viewport(), parent, [4]int32{})
		if len(c.Calls)-before != attempts || c.Calls[before].Ref == 0 || c.Calls[before].Position != image.Pt(280, 400) {
			t.Fatalf("op%d did not retry old anchor", op)
		}
		tail := c.Objs.DeadlineList
		lifetime := uint32(60)
		if op == 8 {
			lifetime = 20
		}
		if tail == nil || tail == parent || tail.Deadline != 120+lifetime || tail.Field_87 != nil || tail.Field_88 != nil {
			t.Fatalf("op%d retry decay ownership", op)
		}
		tailWords := unsafe.Slice((*uint32)(tail.C()), 128)
		if words[108] != 320 || words[109] != 448 || tailWords[108] != 320 || tailWords[109] != 448 {
			t.Fatalf("op%d successful anchor endpoints", op)
		}
		before = len(c.Calls)
		legacy.PortTestClientUpdate(op, c.Viewport(), parent, [4]int32{})
		if len(c.Calls)-before != attempts-1 || c.Objs.DeadlineList != tail || tail.Field_87 != nil {
			t.Fatalf("op%d unchanged anchor duplicated tail", op)
		}
		c.snapshotDrawables(t)
	}
}
