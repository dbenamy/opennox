//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
)

func TestScoreboardRefreshBoundaries(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	fps := o.c.srv.TickRate()
	if fps < 2 {
		t.Fatal("fixture tick rate", fps)
	}
	cases := []struct {
		last, frame uint32
		refresh     bool
	}{{0, fps - 1, false}, {0, fps, false}, {0, fps + 1, true}, {^uint32(0), 0, false}, {^uint32(0), fps - 1, false}, {^uint32(0), fps, true}, {^uint32(0) - fps + 1, 0, false}, {^uint32(0) - fps + 1, 1, true}, {^uint32(0) - fps + 1, ^uint32(0), true}}
	for _, mode := range []uint32{1, 2} {
		for _, dirty := range []uint32{0, 1} {
			for _, tc := range cases {
				o.resetRank(t)
				o.constructRank(t)
				p := &o.players[0]
				p.Active = 1
				p.Lessons = 1
				p.Field3680 = 0
				*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = 0
				*o.rankWords["dword_5d4594_1090120"] = mode
				if mode == 1 {
					noxflags.SetGame(noxflags.GameModeQuest)
				}
				o.rankWindow().Show()
				o.c.GUI.Draw()
				if got := *memmap.PtrUint32(0x5D4594, 1084196); got != 1 {
					t.Fatal("initial score", got)
				}
				p.Lessons = 2
				*memmap.PtrUint32(0x5D4594, 1090124) = tc.last
				*o.rankWords["dword_587000_145664"] = dirty
				// Quest draw requests another update after each refresh, so retain that actual
				// owner state when examining its continuous refresh path.
				if mode == 1 {
					*o.rankWords["dword_587000_145664"] = 1
				}
				o.c.srv.SetFrame(tc.frame)
				o.c.GUI.Draw()
				refresh := tc.refresh || dirty != 0 || mode == 1
				want, last := uint32(1), tc.last
				if refresh {
					want, last = 2, tc.frame
				}
				if got := *memmap.PtrUint32(0x5D4594, 1084196); got != want {
					t.Fatalf("mode%d dirty%d last%x frame%x score%d want%d", mode, dirty, tc.last, tc.frame, got, want)
				}
				if got := *memmap.PtrUint32(0x5D4594, 1090124); got != last {
					t.Fatalf("refresh timestamp%x want%x", got, last)
				}
				rows = append(rows, o.rankCapture(t, 3, 1))
			}
		}
	}
	scoreboardCapture(t, "refresh-boundaries", rows, "527161229c963f1a1e51cbc67f0b80c989418d9acc3cd7de5e0aff9a1e94a0c0")
}
