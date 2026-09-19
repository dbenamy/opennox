//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"image"
	"testing"
	"unsafe"
)

var presentationShieldNames = []string{"ReflectiveShieldNW", "ReflectiveShieldN", "ReflectiveShieldNE", "ReflectiveShieldW", "ReflectiveShieldE", "ReflectiveShieldSW", "ReflectiveShieldS", "ReflectiveShieldSE"}

type presentationShieldClient struct {
	*objectDrawingClient
	fail, calls int
}

func (c *presentationShieldClient) Nox_new_drawable_for_thing(typ int) *client.Drawable {
	i := c.calls
	c.calls++
	if i == c.fail {
		return nil
	}
	dr := c.effectsTestClient.Nox_new_drawable_for_thing(typ)
	if dr != nil {
		c.next++
		c.refs[dr] = c.next
		c.owner.unlinked = append(c.owner.unlinked, dr)
	}
	return dr
}
func newPresentationShields(t *testing.T, fail int) (*objectRenderOwner, *presentationShieldClient) {
	o := newObjectRenderOwner(t, presentationShieldNames...)
	for _, r := range blobdata.PortTestClientPresentationTables() {
		copy(serverConfigOwnBytes(t, r.Base, r.Offset, len(r.Data)), r.Data)
	}
	clear(serverConfigOwnBytes(t, 0x5D4594, 1217468, 40))
	c := &presentationShieldClient{objectDrawingClient: o.proxy, fail: fail}
	old := legacy.GetClient
	legacy.GetClient = func() legacy.Client { return c }
	t.Cleanup(func() { legacy.GetClient = old })
	t.Cleanup(legacy.PortTestPresentationShieldDestroy)
	return o, c
}
func TestClientPresentationShieldOwnership(t *testing.T) {
	type record struct {
		Failure         int
		Success         bool
		Present, Marked [9]bool
		Deleted         []uint32
	}
	var rows []record
	for fail := -1; fail < 8; fail++ {
		t.Run(fmt.Sprint(fail), func(t *testing.T) {
			o, c := newPresentationShields(t, fail)
			*memmap.PtrUint32(0x5D4594, 1217504) = 99
			ok := legacy.PortTestPresentationShieldInit()
			if ok != (fail < 0) || c.calls != 8 {
				t.Fatal("shield creation boundary")
			}
			var present, marked [9]bool
			j := 0
			for i := 0; i < 9; i++ {
				dr := (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1217468+4*uintptr(i)))
				if i == 4 {
					if dr != nil {
						t.Fatal("center shield")
					}
					continue
				}
				present[i] = dr != nil
				if present[i] != (j != fail) {
					t.Fatal("partial shield allocation")
				}
				if dr != nil {
					marked[i] = dr.ObjFlags&0x1000000 != 0
					if marked[i] != (fail < 0 || j < fail) || dr.TypeIDVal != uint32(o.c.Things.IndByID(presentationShieldNames[j])) {
						t.Fatal("shield marking/type")
					}
				}
				j++
			}
			wantCounter := uint32(99)
			if ok {
				wantCounter = 0
			}
			if memmap.Uint32(0x5D4594, 1217504) != wantCounter {
				t.Fatal("shield counter reset")
			}
			legacy.PortTestPresentationShieldDestroy()
			wantDeleted := 8
			if fail >= 0 {
				wantDeleted--
			}
			if len(o.unlinked) != 0 || len(o.rawDeleted) != wantDeleted {
				t.Fatal("shield ownership cleanup")
			}
			for i := 0; i < 10; i++ {
				if memmap.Uint32(0x5D4594, 1217468+4*uintptr(i)) != 0 {
					t.Fatal("shield stale pointer/counter")
				}
			}
			rows = append(rows, record{fail, ok, present, marked, append([]uint32(nil), o.rawDeleted...)})
			legacy.PortTestPresentationShieldDestroy()
		})
	}
	spellbookCapture(t, "client-presentation-shield-ownership", rows, "feddc7a6eb3a9b0d95050a893601332d1892bee425d3fc003398741753822400")
}
func TestClientPresentationShieldDraw(t *testing.T) {
	o, _ := newPresentationShields(t, -1)
	if !legacy.PortTestPresentationShieldInit() {
		t.Fatal("shield init")
	}
	data, free := alloc.Make([]uint32{}, 2)
	t.Cleanup(free)
	data[0], data[1] = 8, uint32(uintptr(o.images[0].C()))
	player := o.drawable(7, image.Pt(48, 48))
	type record struct {
		Direction int
		Position  image.Point
		Z         int16
		Shield    image.Point
		Pixels    string
	}
	var rows []record
	for dir := 0; dir < 9; dir++ {
		if dir == 4 {
			continue
		}
		shield := (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1217468+4*uintptr(dir)))
		shield.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
		shield.DrawData = unsafe.Pointer(&data[0])
		for _, pos := range []image.Point{{48, 48}, {-100, -100}, {100000, -100000}} {
			for _, z := range []int16{-32768, -10, 0, 10, 32767} {
				player.PosVec, player.ZVal, player.AnimDir = pos, uint16(z), byte(dir)
				want := pos.Add(image.Pt(int(memmap.Int32(0x587000, 161776+8*uintptr(dir))), int(z)+int(memmap.Int32(0x587000, 161780+8*uintptr(dir)))))
				clear(o.pix.Pix)
				blank := effectsPixelHash(o.pix)
				legacy.PortTestPresentationShieldDraw(o.c.Viewport(), player)
				got := effectsPixelHash(o.pix)
				if pos == image.Pt(48, 48) && z == 0 && got == blank {
					t.Fatal("shield rendered no pixels")
				}
				if shield.PosVec != want {
					t.Fatal("shield directional placement/height")
				}
				clear(o.pix.Pix)
				ccall.CallIntPtr2(shield.DrawFuncPtr, o.c.Viewport().C(), shield.C())
				if got != effectsPixelHash(o.pix) {
					t.Fatal("shield did not invoke drawable callback")
				}
				rows = append(rows, record{dir, pos, z, want, got})
			}
		}
	}
	spellbookCapture(t, "client-presentation-shield-draw", rows, "625a6c127c1aa5989c04e14c5a6a7f8a7f143ff7ca02c01898abb3edf889d1da")
}
