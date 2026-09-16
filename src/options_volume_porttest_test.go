//go:build porttest

package opennox

import (
	"fmt"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/timer"
)

func TestOptionsVolumeCheckbox(t *testing.T) {
	for _, menu := range []bool{false, true} {
		for ch := 0; ch < 3; ch++ {
			t.Run(fmt.Sprintf("menu%v/channel%d", menu, ch), func(t *testing.T) {
				o := newOptionsAudioOwner(t, menu)
				// Use the real Go checkbox. Sliding an enabled channel to zero must mute it.
				for n, off := range []int{126996, 122848, 93156} {
					*o.words[off] = 1
					button := o.buttons[n][1]
					button.DrawData().Field0 |= 4
					*o.words[o.buttonOffset(n)] = uint32(uintptr(button.C()))
					*o.timers[n] = timer.Timer{Current: 0x40000000, Target: 0x40000000}
				}
				ret := legacy.PortTestOptionsEvent(menu, o.root, 16393, o.controls[351+ch], 0)
				if ret != 0 {
					t.Fatalf("return %d", ret)
				}
				if *o.words[[]int{126996, 122848, 93156}[ch]] != 0 || o.buttons[ch][1].DrawData().Field0&4 != 0 {
					t.Fatalf("zero-volume slider did not mute real checkbox: enabled=%d checked=%v", *o.words[[]int{126996, 122848, 93156}[ch]], o.buttons[ch][1].DrawData().Field0&4 != 0)
				}
			})
		}
	}
}
