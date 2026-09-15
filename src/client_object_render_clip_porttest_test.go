//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientObjectRenderClipSaveRestore(t *testing.T) {
	o := newObjectRenderOwner(t)
	d := o.c.r.Data()
	type result struct {
		Case, Step int
		Return     int32
		Render     objectRenderResult
	}
	var out []result
	id := 0
	for _, flags := range []uint32{0, 1, 0xffffffff} {
		for _, x := range []int{-2147483648, -17, 0, 95, 2147483647} {
			id++
			o.resetRender(uint32(id), 120)
			clip := image.Rect(x, 3, x+10, 80)
			clip2 := image.Rect(x, 7, x+31, 91)
			d.SetClipRect(clip)
			d.SetClipRect2(clip2)
			*(*uint32)(d.C()) = flags
			got := legacy.PortTestObjectRenderClip(false)
			if got != 0 {
				t.Fatal("restore without saved state")
			}
			out = append(out, result{id, 0, got, o.renderResult(t, id, 0, int(got))})
			legacy.PortTestObjectRenderClip(true)
			out = append(out, result{id, 1, 0, o.renderResult(t, id, 1, 0)})
			d.SetClipRect(image.Rect(2, 4, 9, 16))
			d.SetClipRect2(image.Rect(8, 1, 17, 31))
			*(*uint32)(d.C()) = flags ^ 0xa5a5a5a5
			legacy.PortTestObjectRenderClip(true)
			out = append(out, result{id, 2, 0, o.renderResult(t, id, 2, 0)})
			got = legacy.PortTestObjectRenderClip(false)
			if got != int32(clip2.Max.X) || d.ClipRect() != clip || d.ClipRect2() != clip2 || *(*uint32)(d.C()) != flags {
				t.Fatal("restore must use first saved clipping state")
			}
			out = append(out, result{id, 3, got, o.renderResult(t, id, 3, int(got))})
			got = legacy.PortTestObjectRenderClip(false)
			if got != 0 {
				t.Fatal("restored state remained pending")
			}
			out = append(out, result{id, 4, got, o.renderResult(t, id, 4, int(got))})
		}
	}
	effectsCapture(t, "object-render-clip-state", out, len(out), "575ba34078ad52f3bffbd22abda55c928c606cd21a98bea724979d26a71053ec")
}
