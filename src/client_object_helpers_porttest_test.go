//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingHelpers(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	type result struct {
		Op, Input int
		Return    uint32
		Render    []uint32
	}
	var out []result
	for _, typ := range []int{4, 5, 0, -1} {
		o.reset(1, 120)
		got := legacy.PortTestObjectDrawHelper(3, nil, nil, typ)
		if typ == 4 {
			want := uint32(uintptr(o.weapon.C()))&0xffffff00 | uint32(o.weapon.Colors12[6].R)
			if got != want {
				t.Fatal("base material helper ABI return")
			}
			got &= 255 // high bytes checked against the owned definition above
		} else if got != 0 {
			t.Fatal("missing material definition return")
		}
		out = append(out, result{3, typ, got, append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)})
	}
	for i := -1; i < len(o.mods); i++ {
		var arg unsafe.Pointer
		if i >= 0 {
			arg = o.mods[i].C()
		}
		got := legacy.PortTestObjectDrawHelper(0, nil, arg, 0)
		want := uint32(0)
		if i == 1 {
			want = 1
		}
		if i == 3 {
			want = 2
		}
		if got != want {
			t.Fatal("modifier team-name lookup")
		}
		out = append(out, result{0, i, got, nil})
	}
	for _, color := range []int{-1, 0, 1, 2, 3, 255, 256} {
		got := legacy.PortTestObjectDrawHelper(7, nil, nil, color)
		if color == 1 || color == 2 {
			if got != uint32(uintptr(unsafe.Pointer(&c.srv.Teams.Arr[color]))) {
				t.Fatal("team color lookup ownership")
			}
			got = uint32(color)
		} else if got != 0 {
			t.Fatal("missing team-color return")
		}
		out = append(out, result{7, color, got, nil})
	}
	for _, class := range []uint32{0, 0x10000000} {
		for i := -1; i < len(o.mods); i++ {
			o.reset(1, 120)
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
			w := unsafe.Slice((*uint32)(dr.C()), 128)
			w[28] = class
			if i >= 0 {
				w[109] = uint32(uintptr(o.mods[i].C()))
			}
			got := legacy.PortTestObjectDrawHelper(1, dr, nil, 0)
			want := uint32(0)
			if class != 0 && (i == 1 || i == 3) {
				want = uint32(i/2 + 1)
			}
			if got != want {
				t.Fatal("drawable modifier team lookup")
			}
			out = append(out, result{1, i, got, nil})
		}
	}
	*memmap.PtrUint32(0x587000, 177488) = 0
	if legacy.PortTestObjectDrawHelper(0, nil, o.mods[1].C(), 0) != 0 {
		t.Fatal("empty team name table")
	}
	effectsCapture(t, "object-drawing-helpers", out, len(out), "6d350ea30a670458baf1396910e5d219e4ee057a9364a89ac5592ac89da29c00")
}
