//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestBriefingWindowDrawModeSelection(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for _, width := range []int{640, 800} {
		for mode := uint32(0); mode < 8; mode++ {
			o.resetWindow(t)
			o.resize(width, 480)
			o.constructBriefing(t)
			legacy.PortTestBriefingWindow(7, 254, 1, 2, 0)
			*memmap.PtrUint32(0x5D4594, 832472) = mode
			*o.briefWords["dword_5d4594_831244"] = 0
			o.c.srv.SetFrame(30)
			w := o.briefWindow().ChildByID(1010)
			o.displayText = nil
			o.loads = nil
			v := legacy.PortTestBriefingWindow(4, txptr(w.C()), txptr(w.DrawData().C()), 0, 0)
			if v != 1 {
				t.Fatal("draw return")
			}
			if o.c.r.Data().ClipRect() != o.pix.Rect {
				t.Fatal("draw must restore full-buffer clipping")
			}
			rows = append(rows, o.windowCapture(t, 4, v))
		}
	}
	briefingWindowCapture(t, "draw-modes", rows, "a5d747c5dfdadfa530f70050f2d57a145710980f5676b274ef247d9b61a6527b")
}
func TestBriefingWindowTimedDismissal(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	duration := uint64(nox_client_getBriefDuration())
	type gateCase struct {
		scroll     float32
		mode       uint32
		elapsed    uint64
		ready, ack uint32
		voice      bool
	}
	cases := []gateCase{
		{0, 2, duration - 1, 1, 0, false}, {0, 2, duration, 1, 0, false}, {0, 2, duration + 1, 1, 0, false},
		{0, 2, ^uint64(0), 1, 0, false}, {0, 2, duration + 1, 0, 0, false}, {0, 2, duration + 1, 2, 0, false},
		{0, 2, duration + 1, 1, 0, true}, {1.25, 2, duration + 1, 1, 0, false}, {-0.75, 2, duration + 1, 1, 0, false},
		{-1.25, 2, duration + 1, 1, 0, false}, {0, 1, duration + 1, 1, 0, false}, {0, 1, duration + 1, 1, 1, false},
		{0, 4, duration + 1, 1, 0, false}, {0, 4, duration + 1, 1, 1, false}, {0, 0, duration + 1, 1, 0, false},
	}
	for _, chapterMode := range []uint32{0, 1, 255} {
		for _, g := range cases {
			o.resetWindow(t)
			o.constructBriefing(t)
			legacy.PortTestBriefingWindow(7, 254, 1, 2, 0)
			*o.briefWords["dword_5d4594_831220"] = chapterMode
			*o.briefWords["dword_5d4594_831276"] = math.Float32bits(g.scroll)
			*memmap.PtrUint32(0x5D4594, 831280) = 0
			*memmap.PtrUint32(0x5D4594, 832472) = g.mode
			*memmap.PtrUint32(0x5D4594, 832488) = g.ack
			*o.briefWords["dword_5d4594_831244"] = g.ready
			*o.briefWords["dword_5d4594_832480"] = 1
			o.ticks = 100000
			*memmap.PtrUint64(0x5D4594, 831292) = o.ticks - g.elapsed
			if g.voice {
				legacy.Dialogs.PlayFile("fixture_pending", 100)
			}
			o.c.srv.SetFrame(31)
			w := o.briefWindow().ChildByID(1010)
			o.displayText = nil
			o.loads = nil
			v := legacy.PortTestBriefingWindow(4, txptr(w.C()), txptr(w.DrawData().C()), 0, 0)
			scroll := g.scroll
			if chapterMode == 255 {
				scroll -= 1
			}
			dismissed := int64(scroll) <= 0 && g.ready == 1 && !g.voice && g.elapsed > duration && (g.ack == 1 || g.mode&5 == 0)
			if v != 1 || *o.briefWords["dword_5d4594_831276"] != math.Float32bits(scroll) {
				t.Fatal("scroll precision/return")
			}
			if (*memmap.PtrUint64(0x5D4594, 831292) == 0) != dismissed {
				t.Fatal("timed dismissal gate")
			}
			wantSound := dismissed && g.mode&2 != 0
			if (len(o.sounds) == 1) != wantSound || len(o.sounds) > 1 {
				t.Fatal("special sound count")
			}
			if wantSound && o.sounds[0] != [2]int{582, 100} {
				t.Fatal("special sound identity/volume")
			}
			rows = append(rows, o.windowCapture(t, 4, v))
			if dismissed {
				before := *o.briefWords["dword_5d4594_831276"]
				w.Draw()
				if *o.briefWords["dword_5d4594_831276"] != before {
					t.Fatal("completion must install no-op draw")
				}
			}
		}
	}
	briefingWindowCapture(t, "timed-dismissal", rows, "a2aba758cbbbdf61177a2892ff27a1bf3fea99175155a2da75f9fa46d413dac9")
}
