//go:build porttest

package opennox

import (
	"image"
	"runtime"
	"testing"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientSliderDrawingStates(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newSliderOwner(t)
	var out []sliderResult
	id := 0
	for horizontal := 0; horizontal < 2; horizontal++ {
		for mode := 0; mode < 2; mode++ {
			for parent := 0; parent < 2; parent++ {
				for enabled := 0; enabled < 2; enabled++ {
					for highlight := 0; highlight < 2; highlight++ {
						for mask := 0; mask < 32; mask++ {
							id++
							size := image.Pt(61, 21)
							if horizontal == 0 {
								size.X, size.Y = size.Y, size.X
							}
							o.create(t, id, horizontal, mode, parent, 1, gui.StatusFlags(enabled*8), 0, 100, size, func(d *gui.WindowData) {
								d.Field0 = uint32(highlight * 2)
								colors := []*uint32{&d.BgColorVal, &d.EnColorVal, &d.DisColorVal, &d.SelColorVal, &d.HlColorVal}
								images := []*noxrender.ImageHandle{&d.BgImageHnd, &d.EnImageHnd, &d.DisImageHnd, &d.SelImageHnd, &d.HlImageHnd}
								for i := range colors {
									if mask&(1<<i) != 0 {
										*colors[i] = 0x80000000
										*images[i] = nil
									}
								}
							})
							for step, state := range []uint32{0, 2, 4} {
								o.thumb.DrawData().Field0 = state
								out = append(out, o.snapshot(t, id, step, 0))
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "slider-drawing-states", out, len(out), "399b603de237d18e5542541e041686cc5caaf158cd1ad0eaf3c93954a9fe6868")
}

func TestClientSliderConstructorContract(t *testing.T) {
	o := newSliderOwner(t)
	for _, style := range []gui.StyleFlags{0, 1, 4, 32} {
		*o.input = gui.SliderData{Min: 10, Max: 30, Field2: 0x13572468, Field3: 17}
		before := *o.input
		*o.draw = gui.WindowData{Style: style}
		w := legacy.Nox_gui_newSlider_4B4EE0(o.parent, 8, 10, 12, 70, 20, o.draw, o.input)
		if w != nil || *o.input != before || o.draw.Window != nil || o.parent.Field100Ptr != nil {
			t.Fatalf("invalid style %x mutated constructor state", style)
		}
	}
	o.create(t, 1, 1, 0, 1, 1, 8, 0, 100, image.Pt(70, 20), func(d *gui.WindowData) { d.Style = 0x18 })
	if o.thumb.Size() != image.Pt(10, 20) {
		t.Fatal("both orientation bits must prefer horizontal")
	}
}
