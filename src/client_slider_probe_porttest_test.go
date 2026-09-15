//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"math"
	"runtime"
	"testing"
)

func TestClientSliderOwner(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if legacy.PortTestProtectionFloatCW()&0x0f00 != 0x0200 {
		t.Fatal("slider fixture requires hosted53-bit nearest arithmetic")
	}
	o := newSliderOwner(t)
	for horizontal := 0; horizontal < 2; horizontal++ {
		size := image.Pt(20, 70)
		if horizontal != 0 {
			size = image.Pt(70, 20)
		}
		o.create(t, horizontal+1, horizontal, 0, 1, 1, 8, 0, 100, size)
		data := (*gui.SliderData)(o.win.WidgetData)
		if o.input.Field2 != math.Float32bits(0.6) || *data != *o.input {
			t.Fatal("constructor scale or copied slider data")
		}
		if o.win.Flags&0x100 == 0 {
			t.Fatal("slider constructor omitted required flag")
		}
		o.win.Func94(&gui.RawEvent{Event: 16394, Arg1: 50})
		if data.Field3 != 50 {
			t.Fatal("programmatic slider value")
		}
		want := image.Pt(0, 30)
		if horizontal != 0 {
			want = image.Pt(30, 0)
		}
		if o.thumb.Offs() != want {
			t.Fatalf("thumb position got%v want%v", o.thumb.Offs(), want)
		}
		before := effectsPixelHash(o.pix)
		o.win.Draw()
		o.thumb.Draw()
		if effectsPixelHash(o.pix) == before {
			t.Fatal("actual slider/child did not draw")
		}
		o.win.Func94(&gui.RawEvent{Event: 16394, Arg1: 101})
		if data.Field3 != 50 {
			t.Fatal("out-of-range programmatic update accepted")
		}
		o.notices = nil
		o.win.Func94(gui.WindowFocus(true))
		if len(o.notices) != 1 || o.notices[0].Code != 16387 || o.notices[0].A != 1 || o.notices[0].B != 77 {
			t.Fatal("focus notification did not reach real owner")
		}
	}
}
