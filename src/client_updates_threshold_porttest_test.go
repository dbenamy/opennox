//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

// The original tail condition is squared distance > 200, including signed
// coordinate deltas. Test the exact cutoff independently of captured hashes.
func TestClientUpdatesTailDistanceThreshold(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	for _, op := range []int{8, 9} {
		for _, delta := range []image.Point{image.Pt(0, 0), image.Pt(14, 1), image.Pt(14, 2), image.Pt(-14, -2), image.Pt(14, 3), image.Pt(-14, 3)} {
			c.resetCase(effects, pix, 13, 120)
			env.Reset(0)
			c.srv.SetTickRate(60)
			pos := image.Pt(320, 448)
			anchor := pos.Sub(delta)
			parent := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
			words := unsafe.Slice((*uint32)(parent.C()), 128)
			words[108], words[109] = uint32(anchor.X), uint32(anchor.Y)
			parent.Field_8, parent.Field_9 = uint32(anchor.X), uint32(anchor.Y)
			legacy.PortTestClientUpdate(op, c.Viewport(), parent, [4]int32{})
			wantTail := delta.Y == 3
			tailCalls := 0
			for _, call := range c.Calls {
				if call.Type == 28 || call.Type == 29 {
					tailCalls++
					if call.Position != anchor || call.Ref == 0 {
						t.Fatal("tail origin or ownership")
					}
				}
			}
			if (tailCalls == 1) != wantTail || tailCalls > 1 || (c.Objs.DeadlineList != nil) != wantTail {
				t.Fatalf("op%d delta%v cutoff", op, delta)
			}
			wantAnchor := anchor
			if wantTail {
				wantAnchor = pos
			}
			if words[108] != uint32(wantAnchor.X) || words[109] != uint32(wantAnchor.Y) {
				t.Fatal("tail cutoff advanced anchor")
			}
			c.snapshotDrawables(t)
		}
	}
}
