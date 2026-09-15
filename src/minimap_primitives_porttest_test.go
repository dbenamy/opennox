//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"os"
	"testing"
)

type minimapPrimitiveResult struct {
	Op     int
	Zoom   uint32
	Inputs [5]int
	Return uint32
	Render objectRenderResult
}

func minimapCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MINIMAP_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %s", label, got)
	// Expectations were captured from the original C owners.
	if got != want {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
func TestMinimapZoomBoundaries(t *testing.T) {
	zoom, restore := legacy.PortTestMinimapZoom()
	t.Cleanup(restore)
	cases := [][3]uint32{
		{0, 500, 10}, {1, 500, 11}, {490, 500, 500}, {499, 500, 509}, {500, 500, 510}, {510, 500, 520},
		{1200, 1190, 1210}, {1750, 1740, 1760}, {2000, 1990, 2010}, {3990, 3980, 4000}, {4000, 3990, 4000}, {4001, 3991, 4000},
		{0x7fffffff, 0x7ffffff5, 4000}, {0x80000000, 0x7ffffff6, 4000}, {0xfffffffa, 500, 4}, {0xffffffff, 500, 9},
	}
	var rows [][3]uint32
	for _, c := range cases {
		for op := 0; op < 2; op++ {
			*zoom = c[0]
			ret := legacy.PortTestMinimap(op, 0, 0, 0, 0, 0)
			if *zoom != c[op+1] || ret != *zoom {
				t.Fatalf("op %d zoom %08x -> %08x, want %08x", op, c[0], *zoom, c[op+1])
			}
			rows = append(rows, [3]uint32{c[0], *zoom, ret})
		}
		*zoom = 1234
		ret := legacy.PortTestMinimap(2, int(c[0]), 0, 0, 0, 0)
		if ret != c[0] || *zoom != c[0] {
			t.Fatal("setter changed raw zoom")
		}
		rows = append(rows, [3]uint32{c[0], *zoom, ret})
	}
	minimapCapture(t, "zoom", rows, "9015214205f48b5e23b92e84a804190a307ad83c87768da20aea2b25bd7f9af3")
}
func TestMinimapPrimitiveRendering(t *testing.T) {
	o := newObjectRenderOwner(t)
	zoom, restore := legacy.PortTestMinimapZoom()
	t.Cleanup(restore)
	var rows []minimapPrimitiveResult
	for _, z := range []uint32{500, 1199, 1200, 1201, 1749, 1750, 1751, 1999, 2000, 2001, 4000} {
		for _, pos := range []image.Point{image.Pt(48, 48), image.Pt(0, 0), image.Pt(95, 95), image.Pt(-4, 50)} {
			for direction := 0; direction < 13; direction++ {
				o.resetRender(1, 120)
				*zoom = z
				blank := effectsPixelHash(o.pix)
				size := int(2300 / z)
				ret := legacy.PortTestMinimap(3, pos.X, pos.Y, direction, size, 0)
				if pos == image.Pt(48, 48) {
					changed := effectsPixelHash(o.pix) != blank
					if (direction <= 10 || z > 2000) != changed {
						t.Fatalf("wall visibility direction%d zoom%d changed%v", direction, z, changed)
					}
				}
				rows = append(rows, minimapPrimitiveResult{3, z, [5]int{pos.X, pos.Y, direction, size, 0}, ret, o.renderResult(t, len(rows), 0, int(ret))})
			}
			for _, op := range []int{5, 6, 7, 9, 10} {
				o.resetRender(1, 120)
				*zoom = z
				o.c.r.Data().SetColor2(noxcolor.RGB5551Color(240, 180, 60))
				blank := effectsPixelHash(o.pix)
				ret := legacy.PortTestMinimap(op, pos.X, pos.Y, 0, 0, 0)
				if pos == image.Pt(48, 48) && effectsPixelHash(o.pix) == blank {
					t.Fatalf("glyph%d zoom%d did not draw", op, z)
				}
				rows = append(rows, minimapPrimitiveResult{op, z, [5]int{pos.X, pos.Y}, ret, o.renderResult(t, len(rows), 0, int(ret))})
			}
		}
	}
	minimapCapture(t, "primitives", rows, "9494a590aa021cdba721c48dbaab95e8869d63eac6212a8103f7a58b007883aa")
}

// Translation and clipping contracts compare actual pixels independently of the
// frozen hashes. The outline deliberately has three sides in the original game.
func TestMinimapPrimitiveContracts(t *testing.T) {
	o := newObjectRenderOwner(t)
	zoom, restore := legacy.PortTestMinimapZoom()
	t.Cleanup(restore)
	var rows []minimapPrimitiveResult
	for _, op := range []int{3, 4, 5, 6, 7, 8, 9, 10} {
		for _, z := range []uint32{500, 1200, 1750, 2000, 2300} {
			var original []uint16
			for _, shift := range []int{0, 7} {
				o.resetRender(1, 120)
				*zoom = z
				o.c.r.Data().SetColor2(noxcolor.RGB5551Color(240, 180, 60))
				a, b, c, d, e := 40+shift, 40+shift, 0, 0, 0
				if op == 3 {
					c, d = 2, 7
				}
				if op == 4 {
					c, d, e = 55+shift, 49+shift, int(noxcolor.RGB5551Color(240, 180, 60).Color16())
				}
				if op == 8 {
					c, d = 15, 11
				}
				ret := legacy.PortTestMinimap(op, a, b, c, d, e)
				if shift == 0 {
					original = append([]uint16(nil), o.pix.Pix...)
				} else {
					for y := 15; y < 75; y++ {
						for x := 15; x < 75; x++ {
							if original[y*o.pix.Stride+x] != o.pix.Pix[(y+shift)*o.pix.Stride+x+shift] {
								t.Fatalf("op%d zoom%d translation mismatch at %d,%d", op, z, x, y)
							}
						}
					}
				}
				rows = append(rows, minimapPrimitiveResult{op, z, [5]int{a, b, c, d, e}, ret, o.renderResult(t, len(rows), 0, int(ret))})
			}
			o.resetRender(1, 120)
			*zoom = z
			o.c.r.Data().SetColor2(noxcolor.RGB5551Color(240, 180, 60))
			clip := image.Rect(40, 40, 50, 50)
			o.c.r.Data().SetClip(true)
			o.c.r.Data().SetClipRect(clip)
			o.c.r.Data().SetClipRect2(image.Rect(40, 40, 49, 49))
			before := append([]uint16(nil), o.pix.Pix...)
			a, b, c, d, e := 40, 40, 0, 0, 0
			if op == 3 {
				c, d = 2, 17
			}
			if op == 4 {
				a, b, c, d, e = 30, 44, 60, 44, int(noxcolor.RGB5551Color(240, 180, 60).Color16())
			}
			if op == 8 {
				c, d = 15, 11
			}
			ret := legacy.PortTestMinimap(op, a, b, c, d, e)
			for y := 0; y < 96; y++ {
				for x := 0; x < 96; x++ {
					i := y*o.pix.Stride + x
					if !image.Pt(x, y).In(clip) && before[i] != o.pix.Pix[i] {
						t.Fatalf("op%d changed pixel outside clip: %d,%d", op, x, y)
					}
				}
			}
			rows = append(rows, minimapPrimitiveResult{op, z, [5]int{a, b, c, d, e}, ret, o.renderResult(t, len(rows), 0, int(ret))})
		}
	}
	minimapCapture(t, "primitive-contracts", rows, "f1e50548eda7abb2b9267dee83e3cb94f88ff08c28c0b0473ba33b48b892ddb8")
}
