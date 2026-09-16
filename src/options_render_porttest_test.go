//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestOptionsOverlayRendering(t *testing.T) {
	type record struct {
		Pos, Size    [2]int
		Clip, Return int
		Pixels       string
	}
	var rows []record
	o := newOptionsOwner(t)
	o.root.SetPos(image.Pt(3, 4))
	w := o.c.GUI.NewWindowRaw(o.root, 8, 0, 0, 1, 1, nil)
	for ci, clip := range []image.Rectangle{o.pix.Rect, image.Rect(5, 7, 70, 75)} {
		for _, pos := range []image.Point{image.Pt(-8, -9), {}, image.Pt(20, 30), image.Pt(90, 90)} {
			for _, size := range []image.Point{{}, {X: 1, Y: 1}, {X: 10, Y: 15}, {X: 120, Y: 120}} {
				for i := range o.pix.Pix {
					o.pix.Pix[i] = uint16(i*251) ^ 0xabcd
				}
				want := append([]uint16(nil), o.pix.Pix...)
				w.SetPos(pos)
				w.SizeVal = size
				o.c.Render().Data().SetClip(true)
				o.c.Render().Data().SetClipRect(clip)
				rect := image.Rectangle{Min: pos.Add(image.Pt(3, 4)), Max: pos.Add(image.Pt(3, 4)).Add(size)}.Intersect(clip)
				for y := rect.Min.Y; y < rect.Max.Y; y++ {
					for x := rect.Min.X; x < rect.Max.X; x++ {
						i := o.pix.PixOffset(x, y)
						want[i] = (want[i] & 0xfbde) >> 1
					}
				}
				ret := legacy.PortTestOptionsDraw(w)
				for i, p := range o.pix.Pix {
					if p != want[i] {
						t.Fatalf("clip%d pos%v size%v pixel%d got=%x want=%x", ci, pos, size, i, p, want[i])
					}
				}
				if ret != 1 {
					t.Fatalf("return %d", ret)
				}
				raw := make([]byte, 2*len(want))
				for i, p := range o.pix.Pix {
					binary.LittleEndian.PutUint16(raw[2*i:], p)
				}
				rows = append(rows, record{[2]int{pos.X, pos.Y}, [2]int{size.X, size.Y}, ci, ret, fmt.Sprintf("%x", sha256.Sum256(raw))})
			}
		}
	}
	spellbookCapture(t, "options-overlay", rows, "e651c3589a15b14f98f7693dfb215172382429bbb0a02953c921c7fe5391bc46")
}
