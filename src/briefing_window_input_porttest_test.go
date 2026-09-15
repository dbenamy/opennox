//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestBriefingWindowInputGating(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for _, ready := range []uint32{0, 1, 2} {
		for _, gate := range []uint32{0, 1, 2} {
			for _, ev := range []uintptr{0, 17, 18, 23} {
				for _, mode := range []uint32{0, 1, 255} {
					o.resetWindow(t)
					o.constructBriefing(t)
					legacy.PortTestBriefingWindow(7, 254, 1, 2, 0)
					*memmap.PtrUint32(0x5D4594, 831248) = ready
					*memmap.PtrUint32(0x5D4594, 527720) = gate
					*o.briefWords["dword_5d4594_831220"] = mode
					legacy.Dialogs.PlayFile("fixture_pending", 100)
					v := legacy.PortTestBriefingWindow(2, uintptr(*o.briefWords["dword_5d4594_831236"]), ev, 0, 0)
					dismissed := ready != 0 && gate != 1 && ev != 17 && ev != 18
					if v != 1 || (memmap.Uint32(0x5D4594, 832488) == 1) != dismissed {
						t.Fatal("input dismissal gate")
					}
					if (legacy.Dialogs.FileToRead() == "") != dismissed {
						t.Fatal("input dialogue cleanup")
					}
					if (*o.briefWords["dword_5d4594_831256"] == 1) != (dismissed && mode == 0) {
						t.Fatal("loss transition")
					}
					if (legacy.Get_nox_gameDisableMapDraw_5d4594_2650672() == 1) != (dismissed && mode == 255) {
						t.Fatal("credits fade transition")
					}
					rows = append(rows, o.windowCapture(t, 2, v))
				}
			}
		}
	}
	briefingWindowCapture(t, "input-gates", rows, "2bf11a87ab6696367a6e0762e9246956d141e4143a9044218a3d0cae8f606c10")
}
