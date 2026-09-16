//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestTileCompositionCounterAndNoop(t *testing.T) {
	o := newTileCompositionOwner(t)
	for axis := 0; axis < 2; axis++ {
		for _, draw := range []bool{false, true} {
			o.populate(0)
			o.resetBuffer()
			vp := tileCompositionViewport(460, 460)
			legacy.PortTestTileCompositionFull(vp)
			*o.counter = 17
			before := o.state()
			ox, oy := int(int32(*o.words["originX"])), int(int32(*o.words["originY"]))
			delta := 23
			if draw {
				delta = 150
			}
			if axis == 0 {
				vp = tileCompositionViewport(ox+delta, oy+40)
				legacy.PortTestTileCompositionHorizontal(vp, ox+delta)
			} else {
				vp = tileCompositionViewport(ox+40, oy+delta)
				legacy.PortTestTileCompositionVertical(vp, oy+delta)
			}
			if draw {
				if *o.counter != 0xffffffff {
					t.Fatalf("axis%d redraw counter%d", axis, *o.counter)
				}
			} else if !reflect.DeepEqual(before, o.state()) {
				t.Fatalf("axis%d no-op changed renderer state", axis)
			}
		}
	}
}
