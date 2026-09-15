//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

type objectImageDraw struct {
	Position, Size image.Point
	Image          uint32
}
type objectDrawingRender struct {
	*NoxRender
	owner *objectDrawingOwner
}

func (r *objectDrawingRender) DrawImageAt(img *noxrender.Image, pos image.Point) {
	r.NoxRender.DrawImageAt(img, pos)
	// The production wrapper installs its own low-level hook. Observe its real
	// output metadata after drawing, preserving both that hook and the pixel path.
	var ref uint32
	if img != nil {
		ref = r.owner.c.imageRefs[uint32(uintptr(img.C()))]
	}
	r.owner.drawTrace = append(r.owner.drawTrace, objectImageDraw{
		image.Pt(int(memmap.Int32(0x973F18, 92)), int(memmap.Int32(0x973F18, 84))),
		image.Pt(int(memmap.Uint32(0x973F18, 88)), int(memmap.Uint32(0x973F18, 76))), ref,
	})
}
func (c *objectDrawingClient) R2() legacy.Render2 {
	return &objectDrawingRender{c.effectsTestClient.r, c.owner}
}

type objectDrawingClient struct {
	*effectsTestClient
	owner *objectDrawingOwner
}

func (c *objectDrawingClient) Nox_xxx_spriteDelete_45A4B0(dr *client.Drawable) int {
	o := c.owner
	for i, p := range o.unlinked {
		if p == dr {
			o.rawDeleted = append(o.rawDeleted, c.refs[dr])
			o.unlinked = append(o.unlinked[:i], o.unlinked[i+1:]...)
			return c.Client.Nox_xxx_spriteDelete_45A4B0(dr)
		}
	}
	panic("raw deletion of an unowned/unlinked object")
}
func (o *objectDrawingOwner) installOwnership(t *testing.T) {
	o.proxy = &objectDrawingClient{o.c, o}
	legacy.GetClient = func() legacy.Client { return o.proxy }
	t.Cleanup(func() {
		for len(o.unlinked) > 0 {
			o.proxy.Nox_xxx_spriteDelete_45A4B0(o.unlinked[0])
		}
	})
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	locks := map[string]int{"DoorLockSilverSW": 28, "DoorLockSilverSE": 29, "DoorLockGoldSW": 30, "DoorLockGoldSE": 31}
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		i, ok := locks[name]
		if !ok {
			panic("unowned named object image")
		}
		o.namedCalls = append(o.namedCalls, name)
		return o.images[i]
	}
	o.modRefs = make(map[uint32]uint32)
	for i, name := range []string{"MaterialRed", "TeamRed", "MaterialBlue", "TeamBlue"} {
		str, free := alloc.CString(name)
		t.Cleanup(free)
		*(*unsafe.Pointer)(o.mods[i].C()) = unsafe.Pointer(str)
		o.modRefs[uint32(uintptr(o.mods[i].C()))] = 0xeb000000 + uint32(i)
		if i == 1 || i == 3 {
			off := uintptr(177488 + 8*(i/2))
			*memmap.PtrUint32(0x587000, off) = uint32(uintptr(unsafe.Pointer(str)))
			*memmap.PtrUint32(0x587000, off+4) = uint32(i/2 + 1)
		}
	}
}
func (o *objectDrawingOwner) newUnlinked(t *testing.T) *client.Drawable {
	t.Helper()
	dr := o.c.Nox_new_drawable_for_thing(4)
	if dr == nil {
		t.Fatal("unlinked drawable allocation failed")
	}
	o.c.next++
	o.c.refs[dr] = o.c.next
	o.unlinked = append(o.unlinked, dr)
	return dr
}
func (o *objectDrawingOwner) snapshot(t *testing.T) [][]uint32 {
	rows := o.c.snapshotDrawables(t, o.unlinked...)
	for _, row := range rows {
		words := row[1:]
		op := -1
		if words[75] >= 0xea000000 && words[75] < 0xea000012 {
			op = int(words[75] - 0xea000000)
		}
		slots := []int{}
		refs := o.modRefs
		switch op {
		case 0:
			slots = []int{109, 110, 111, 112}
			refs = o.c.imageRefs
		case 7, 8, 9, 10, 16, 17:
			slots = []int{108, 109, 110, 111}
		case 6:
			if words[108] != 0 {
				n := o.c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(words[108])))]
				if n == 0 {
					t.Fatal("unowned summon child pointer")
				}
				words[108] = n
			}
		}
		for _, slot := range slots {
			if words[slot] == 0 {
				continue
			}
			n, ok := refs[words[slot]]
			if !ok {
				t.Fatalf("unowned object data at op%d word%d", op, slot)
			}
			words[slot] = n
		}
	}
	return rows
}
func (o *objectDrawingOwner) globals(t *testing.T) []uint32 {
	out := []uint32{*memmap.PtrUint32(0x5D4594, 1313720), *memmap.PtrUint32(0x5D4594, 1313724), *memmap.PtrUint32(0x5D4594, 1321512)}
	for _, r := range [][2]uintptr{{0x5D4594, 1321516}, {0x852978, 8}} {
		p := *memmap.PtrUint32(r[0], r[1])
		var n uint32
		if p != 0 {
			n = o.c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(p)))]
			if n == 0 {
				t.Fatal("unowned drawable global")
			}
		}
		out = append(out, n)
	}
	out = append(out, o.c.tiles.lightsOutBuf...)
	return append(out, o.env.NamedState()...)
}
