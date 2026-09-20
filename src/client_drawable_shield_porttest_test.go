//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

var drawableShieldNames = []string{"SphericalShieldNW", "SphericalShieldN", "SphericalShieldNE", "SphericalShieldW", "SphericalShieldE", "SphericalShieldSW", "SphericalShieldS", "SphericalShieldSE"}

func TestClientDrawableShields(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, drawableShieldNames...)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	type row struct {
		Dir, Mode, Step int
		Return          uint32
		Count           int
		State           [][]uint32
		Globals         []uint32
	}
	var rows []row
	for dir := 0; dir < 9; dir++ {
		for mode := 0; mode < 8; mode++ {
			c.resetCase(env, pix, 1, 100)
			for _, p := range words {
				*p = 0
			}
			code := uint32(17)
			var target *client.Drawable
			if mode != 0 {
				target = c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
				target.NetCode32 = 17
				if mode&1 != 0 {
					code |= 0x8000
					target.ObjClass = object.Class(0x400000)
				}
			}
			if mode >= 4 {
				// An existing shield either matches the owner, is beyond the search box,
				// or uses a different owner key. Spatial search remains production code.
				pos := image.Pt(300, 403)
				if mode == 5 {
					pos.X += 20
				}
				existing := c.Nox_xxx_spriteLoadAdd_45A360_drawable(c.Things.IndByID(drawableShieldNames[0]), pos)
				*(*uint32)(unsafe.Add(existing.C(), 432)) = code
				if mode == 6 {
					*(*uint32)(unsafe.Add(existing.C(), 432)) = 99
				}
			}
			for step := 0; step < 2; step++ {
				before := c.Objs.Count
				raw := legacy.PortTestDrawableEffect(3, nil, [5]uint32{code, uint32(dir)})
				ret := uint32(raw)
				if raw > 1 {
					dr := (*client.Drawable)(unsafe.Pointer(raw))
					ret = c.refs[dr]
					if ret == 0 {
						t.Fatal("unowned shield")
					}
					if target == nil || dr.PosVec != target.PosVec.Add(image.Pt(0, 3)) || *(*uint32)(unsafe.Add(dr.C(), 432)) != code {
						t.Fatal("shield position/owner")
					}
					if c.Objs.Count != before+1 {
						t.Fatal("shield allocation count")
					}
				} else if c.Objs.Count != before {
					t.Fatal("shield early return allocated")
				}
				if mode == 0 && raw != 0 {
					t.Fatal("missing target")
				}
				if dir != 4 && mode != 0 && step == 1 && raw != 1 {
					t.Fatal("duplicate shield not suppressed")
				}
				// Center has type zero and cannot be allocated by the real factory.
				if dir == 4 && mode > 0 && mode < 4 && raw != 0 {
					t.Fatal("center shield")
				}
				globals := make([]uint32, len(words))
				for i, p := range words {
					globals[i] = *p
				}
				rows = append(rows, row{dir, mode, step, ret, c.Objs.Count, c.snapshotDrawables(t), globals})
			}
		}
	}
	drawableStateCapture(t, "shields", rows, "21ef9ab037949226b7d734d33d391908f0d7934c68e49808e14fbe03d25e09a3")
}

func TestClientDrawableShieldScan(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t, drawableShieldNames...)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	legacy.PortTestDrawableEffect(2, nil, [5]uint32{})
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
	for dir := 0; dir < 9; dir++ {
		for _, owner := range []uint32{0, 17, 0x80000000, 0xffffffff} {
			for _, key := range []uint32{0, 17, 0x80000000, 0xffffffff} {
				for _, initial := range []uint32{0, 1, 55} {
					dr.TypeIDVal = *memmap.PtrUint32(0x5D4594, 1313748+4*uintptr(dir))
					*(*uint32)(unsafe.Add(dr.C(), 432)) = owner
					*words[1] = initial
					legacy.PortTestDrawableEffect(4, dr, [5]uint32{key})
					want := initial
					if owner == key {
						want = 1
					}
					if *words[1] != want {
						t.Fatalf("shield scan dir%d owner%x key%x initial%d got%d want%d", dir, owner, key, initial, *words[1], want)
					}
				}
			}
		}
	}
}
