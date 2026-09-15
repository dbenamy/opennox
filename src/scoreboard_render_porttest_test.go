//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestScoreboardPopulatedRendering(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, width := range []int{640, 800} {
		for mode := 1; mode <= 5; mode++ {
			for _, count := range []int{0, 1, 3, 16, 17, 32} {
				o.resetRank(t)
				o.resize(width, 480)
				o.constructRank(t)
				*o.rankWords["dword_5d4594_1090120"] = uint32(mode)
				if mode == 1 {
					noxflags.SetGame(noxflags.GameModeQuest)
				}
				for i := 0; i < count; i++ {
					p := &o.players[i]
					p.Active = 1
					p.PlayerInd = byte(i)
					p.Lessons = int32(i % 7)
					p.Field2140 = uint32(i % 5)
					p.Field2108 = uint32(i % 7)
					p.Field3680 = []uint32{0, 1, 0x20, 0x21}[i%4]
					p.Field4792 = 1
					*(*uint16)(unsafe.Add(unsafe.Pointer(p), 2148)) = uint16(10 + i*17)
					*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = byte(i % 3)
					*(*byte)(unsafe.Add(unsafe.Pointer(p), 2282)) = []byte{0, 25, 26, 50, 51, 100, 255}[i%7]
					*(*byte)(unsafe.Add(unsafe.Pointer(p), 4816)) = byte(i % 4)
					*(*byte)(unsafe.Add(unsafe.Pointer(p), 4824)) = byte(i % 3)
					*(*byte)(unsafe.Add(unsafe.Pointer(p), 4825)) = byte((i / 3) % 3)
				}
				o.rankWindow().Show()
				o.c.GUI.Draw()
				if got := int(*memmap.PtrUint8(0x5D4594, 1090117)); got != count {
					t.Fatalf("mode%d count%d got%d", mode, count, got)
				}
				rows = append(rows, o.rankCapture(t, 3, 1))
			}
		}
	}
	scoreboardCapture(t, "populated-rendering", rows, "df9f586134a4c0af8cec48088c3469c5d2ef519fc8d6858eeb1bff16f14727ea")
}
func TestScoreboardVisibilityAndModeCallers(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, quest := range []bool{false, true} {
		o.resetRank(t)
		if got := legacy.PortTestScoreboard(30, 0, 0, 0); got != 0 {
			t.Fatal("unconstructed rank enabled")
		}
		legacy.PortTestScoreboard(26, 0, 0, 0)
		o.constructRank(t)
		if quest {
			noxflags.SetGame(noxflags.GameModeQuest)
		}
		want := []uint32{2, 3, 4, 0, 2, 3, 4, 0}
		if quest {
			want = []uint32{1, 2, 3, 4, 0, 1, 2, 3}
		}
		for _, mode := range want {
			sub_4703F0()
			if *o.rankWords["dword_5d4594_1090120"] != mode {
				t.Fatalf("cycle mode got%d want%d", *o.rankWords["dword_5d4594_1090120"], mode)
			}
			visible := legacy.PortTestScoreboard(25, 0, 0, 0)
			enabled := legacy.PortTestScoreboard(30, 0, 0, 0)
			expected := uint32(1)
			if mode == 0 {
				expected = 0
			}
			if visible != expected || enabled != expected {
				t.Fatalf("mode%d visibility/enabled %d/%d", mode, visible, enabled)
			}
			rows = append(rows, o.rankCapture(t, 25, visible))
		}
		o.rankWindow().Hide()
		legacy.PortTestScoreboard(26, 0, 0, 0)
		rows = append(rows, o.rankCapture(t, 26, 0))
		sub_470510()
		rows = append(rows, o.rankCapture(t, 31, 0))
		sub_470550()
		rows = append(rows, o.rankCapture(t, 32, 0))
	}
	scoreboardCapture(t, "visibility-mode-callers", rows, "4a035253a9fb35a2c2b045aa2741cd73d37ffdad6e65b1c4ccecbbc01f8ef8ef")
}
