//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestWorldGridWallAttachments(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(c.srv.PortTestMinimapWalls())
	polys := newMapPolygonsOwner(t)
	polygon := polys.construct(t, [][2]float32{{100, 100}, {200, 100}, {200, 200}, {100, 200}})
	*(*byte)(unsafe.Add(polygon, 130)) = 29
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(150, 151))
	if dr == nil {
		t.Fatal("drawable allocation")
	}
	type result struct {
		Drawable               bool
		X, Y, Position, Return int
		Data                   []byte
	}
	var rows []result
	for _, drawable := range []bool{false, true} {
		for position, pos := range []image.Point{{150, 151}, {205, 151}, {300, 300}} {
			for _, p := range []image.Point{{-1, 0}, {0, -1}, {256, 0}, {0, 256}, {0, 0}, {17, 23}, {255, 255}} {
				dr.PosVec = pos
				w := c.srv.Walls.GetWallAtGridRaw(p)
				if w != nil {
					raw := unsafe.Slice((*byte)(w.C()), 36)
					raw[4] = 0x80
					binary.LittleEndian.PutUint16(raw[8:], 0x7e55)
				}
				ret := legacy.PortTestWorldWallAttach(drawable, dr.C(), p.X, p.Y)
				r := result{Drawable: drawable, X: p.X, Y: p.Y, Position: position}
				w = c.srv.Walls.GetWallAtGridRaw(p)
				if p.X < 0 || p.Y < 0 || p.X > 255 || p.Y > 255 {
					if ret != nil || w != nil {
						t.Fatal("out of grid attachment")
					}
				} else {
					if w == nil {
						t.Fatalf("missing wall drawable%v pos%v", drawable, p)
					}
					raw := unsafe.Slice((*byte)(w.C()), 36)
					if raw[4]&0x10 == 0 {
						t.Fatal("attachment flag")
					}
					off := 28
					if drawable {
						off = 32
					}
					if *(*unsafe.Pointer)(unsafe.Add(w.C(), off)) != dr.C() {
						t.Fatal("attachment owner")
					}
					if drawable {
						want := byte(29)
						if position == 2 {
							want = 1
						}
						if raw[8] != want {
							t.Fatalf("polygon byte pos%v got%d want%d", pos, raw[8], want)
						}
						if position == 2 {
							if ret != nil {
								t.Fatal("unexpected polygon")
							}
						} else if ret != polygon {
							t.Fatal("polygon identity")
						}
						r.Return = polys.polygonID(ret)
					} else {
						if ret != w.C() {
							t.Fatal("door wall return")
						}
						r.Return = 1
					}
					r.Data = append([]byte(nil), raw[:16]...)
				}
				rows = append(rows, r)
			}
		}
	}
	drawableStateCapture(t, "world-attachments", rows, "1a2bcd41ef8bf0a054af8affa09c50988c48c4b53cbdb02793742baf5fcd0fc8")
}
