//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestScoreboardTimeAndLessonHeadings(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	oldFlag := o.c.srv.flag3592
	t.Cleanup(func() { o.c.srv.flag3592 = oldFlag })
	o.c.srv.flag3592 = false
	for _, host := range []bool{false, true} {
		for _, enabled := range []uint32{0, 1} {
			for _, permit := range []byte{0, 1} {
				for _, remaining := range []int32{-60001, -1, 0, 999, 1000, 59999, 60000, 125999} {
					o.resetRank(t)
					o.constructRank(t)
					noxflags.SetGame(noxflags.GameModeArena)
					if host {
						noxflags.SetGame(noxflags.GameHost)
					}
					*memmap.PtrUint32(0x587000, 4660) = enabled
					*memmap.PtrUint32(0x5D4594, 3468) = uint32(int64(123456) + int64(remaining))
					for i := uintptr(0); i < 6; i++ {
						*memmap.PtrUint8(0x5D4594, 3500+i) = permit
					}
					ret := legacy.PortTestScoreboard(6, 0, 0, 0)
					w := (*gui.Window)(unsafe.Pointer(uintptr(*o.rankWords["dword_5d4594_1090108"])))
					shown := enabled != 0 && (!host || permit != 0)
					if w.GetFlags().Has(gui.StatusHidden) == shown {
						t.Fatalf("timer visibility host%v enabled%d permit%d", host, enabled, permit)
					}
					if shown {
						want := fmt.Sprintf("Time %d:%02d", remaining/60000, (remaining%60000)/1000)
						got := alloc.GoString16((*gui.StaticTextData)(w.WidgetData).Text)
						if got != want {
							t.Fatalf("remaining%d got%q want%q", remaining, got, want)
						}
					}
					rows = append(rows, o.rankCapture(t, 6, ret))
				}
			}
		}
	}
	scoreboardCapture(t, "time-headings", rows, "6ccc8b3b815323bbcf3877487eba21ac5378b060f9c96e55e9f7c165a32e82fa")
	rows = nil
	for _, host := range []bool{false, true} {
		for _, mode := range []noxflags.GameFlag{noxflags.GameModeArena, noxflags.GameModeElimination, noxflags.GameModeQuest, noxflags.GameModeChat} {
			for _, limit := range []uint16{0, 1, 255, 32767, 32768, 65535} {
				o.resetRank(t)
				o.constructRank(t)
				noxflags.SetGame(mode)
				if host {
					noxflags.SetGame(noxflags.GameHost)
				}
				for i := uintptr(0); i < 6; i++ {
					*memmap.PtrUint16(0x5D4594, 3488+2*i) = limit
				}
				*memmap.PtrUint16(0x5D4594, 371434) = limit
				ret := legacy.PortTestScoreboard(7, 0, 0, 0)
				w := (*gui.Window)(unsafe.Pointer(uintptr(*o.rankWords["dword_5d4594_1090112"])))
				hidden := mode&(noxflags.GameModeQuest|noxflags.GameModeChat) != 0
				if w.GetFlags().Has(gui.StatusHidden) != hidden {
					t.Fatalf("lesson visibility mode%x", mode)
				}
				if !hidden {
					want := fmt.Sprintf("Limit %d", limit)
					got := alloc.GoString16((*gui.StaticTextData)(w.WidgetData).Text)
					if got != want {
						t.Fatalf("lesson%d got%q want%q", limit, got, want)
					}
				}
				rows = append(rows, o.rankCapture(t, 7, ret))
			}
		}
	}
	scoreboardCapture(t, "lesson-headings", rows, "bd1562fbf4aa7aa8430daa1593d42b33a3579ebb64508ef12e0d7210660931d3")
}
