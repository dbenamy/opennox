//go:build porttest

package opennox

import (
	"encoding/binary"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

func TestClientEffectsMovingOrbDistance(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	for _, distance := range []int{0, 30} {
		c.resetCase(env, pix, 1, 0)
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 70))
		data := unsafe.Slice((*byte)(unsafe.Pointer(dr)), 512)
		binary.LittleEndian.PutUint16(data[432:], uint16(48+distance))
		binary.LittleEndian.PutUint16(data[434:], 70)
		data[443], data[444] = 3, 5
		got := legacy.PortTestClientEffects(34, c.Viewport(), dr, [8]int32{}, nil)
		if distance == 0 {
			if got != 0 || c.Objs.Count != 0 {
				t.Fatal("arrived orb was not deleted")
			}
		} else if got != 1 || c.Objs.Count != 1 || dr.PosVec != image.Pt(50, 70) {
			t.Fatal("moving orb requires the initialized production distance table")
		}
	}
}

func TestClientEffectsStoredOrbitCallback(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Frame                         uint32
		Angle, Direction, Age, Return int
		Drawables                     [][]uint32
		Deleted                       []uint32
		Other                         int
	}
	var out []result
	data, free := alloc.Make([]byte{}, 8)
	defer free()
	for _, frame := range []uint32{123, 0xfffffff0} {
		for _, angle := range []int{0, 127, 255} {
			for direction := 0; direction < 2; direction++ {
				c.resetCase(env, pix, 31, frame)
				*(*[4]int16)(unsafe.Pointer(&data[0])) = [4]int16{50, 50, 100, 100}
				legacy.PortTestClientEffects(1, c.Viewport(), nil, [8]int32{4, int32(angle), int32(direction), 0}, unsafe.Pointer(&data[0]))
				dr := c.Objs.List1
				if dr == nil || dr.ClientUpdateFuncPtr != legacy.PortTestEffectsCallback(37) {
					t.Fatal("orbit creation did not install the production callback")
				}
				for _, age := range []int{0, 1, 30, 59, 60} {
					c.srv.SetFrame(frame + uint32(age))
					got := ccall.CallIntPtr2(dr.ClientUpdateFuncPtr, c.Viewport().C(), dr.C())
					if got < 0 || got > 1 || c.srv.Rand.Other.Index() != 33 || c.srv.Rand.Logic.Index() != 31 {
						t.Fatal("stored orbit callback ABI/RNG")
					}
					if age == 60 && got != 0 {
						t.Fatal("stored orbit callback failed to expire")
					}
					out = append(out, result{frame, angle, direction, age, got, c.snapshotDrawables(t), append([]uint32(nil), c.Deleted...), c.srv.Rand.Other.Index()})
					if got == 0 {
						break
					}
				}
			}
		}
	}
	effectsCapture(t, "stored-orbit", out, len(out), "c6a46b642d3aa98891559ec5c4e6eb8ee5de0500745e8db34c08af09eb4328c0")
}

// Optional visual diagnostics accompany, but never replace, raw-word hashes.
func effectsDiagnostic(t *testing.T, name string, pix *noximage.Image16) {
	t.Helper()
	dir := os.Getenv("OPENNOX_CLIENT_EFFECTS_DIAGNOSTICS")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(filepath.Join(dir, name+".png"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	err = png.Encode(f, pix)
	closeErr := f.Close()
	if err != nil {
		t.Fatal(err)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
}

// A ray's binding mode must not change its particle endpoints. In the original
// coordinate branch, the target Y was uninitialized and X received target Y.
func TestClientEffectsChainParticleEndpoints(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	flags := noxflags.GetGame()
	defer func() { noxflags.ResetGame(); noxflags.SetGame(flags) }()
	noxflags.ResetGame()
	for _, endpoints := range [][2]image.Point{
		{image.Pt(20, 40), image.Pt(83, 60)},
		{image.Pt(90, 25), image.Pt(15, 75)},
		{image.Pt(48, 48), image.Pt(48, 48)},
	} {
		var expected []effectsSpawnCall
		var expectedOther int
		for mode := 0; mode < 2; mode++ {
			c.resetCase(env, pix, 31, 4)
			legacy.PortTestClientEffects(45, c.Viewport(), nil, [8]int32{}, nil)
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(14, image.Pt(48, 70))
			data := unsafe.Slice((*byte)(dr.C()), 512)
			if mode == 0 {
				for i, pos := range endpoints {
					binary.LittleEndian.PutUint16(data[437+4*i:], uint16(pos.X))
					binary.LittleEndian.PutUint16(data[439+4*i:], uint16(pos.Y))
				}
			} else {
				data[432] = 1
				for i, pos := range endpoints {
					code := uint32(101 + i)
					binary.LittleEndian.PutUint32(data[437+4*i:], code)
					endpoint := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
					endpoint.NetCode32 = code
				}
			}
			before := len(c.Calls)
			if legacy.PortTestClientEffects(12, c.Viewport(), dr, [8]int32{}, nil) != 1 {
				t.Fatal("chain ray return")
			}
			calls := append([]effectsSpawnCall(nil), c.Calls[before:]...)
			for i := range calls {
				calls[i].Ref = 0
			}
			if mode == 0 {
				expected = calls
				expectedOther = c.srv.Rand.Other.Index()
			} else if !reflect.DeepEqual(calls, expected) || c.srv.Rand.Other.Index() != expectedOther {
				t.Fatal("coordinate and object endpoints emitted different particles")
			}
			if endpoints[0] != endpoints[1] && len(calls) == 0 {
				t.Fatal("nonempty chain ray emitted no particles")
			}
		}
	}
}
