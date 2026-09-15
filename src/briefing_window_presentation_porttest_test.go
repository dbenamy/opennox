//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/player"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestBriefingWindowChapterSelection(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for class, name := range []string{"Warrior", "Wizard", "Conjurer"} {
		for chapter := 0; chapter < 11; chapter++ {
			for _, begin := range []uintptr{0, 1} {
				for _, seen := range []uint32{0, 1} {
					o.resetWindow(t)
					o.constructBriefing(t)
					p := &o.players[0]
					legacy.Set_dword_8531A0_2576(p)
					p.Info().SetPlayerClass(player.Class(class))
					slot := (*uint32)(unsafe.Add(p.C(), 4408+4*chapter))
					*slot = seen
					v := legacy.PortTestBriefingWindow(7, uintptr(chapter), begin, 0, 0)
					want := uint64(1)
					if begin == 0 && seen != 0 {
						want = 0
					}
					if v != want {
						t.Fatal("chapter voice eligibility")
					}
					kind := "Loss"
					if begin != 0 {
						kind = "Begin"
					}
					voice := fmt.Sprintf("voice_%s_%s_%d", name, kind, chapter+1)
					if alloc.GoString((*byte)(unsafe.Pointer(uintptr(*o.briefWords["dword_5d4594_831240"])))) != voice {
						t.Fatal("chapter voice selection")
					}
					if begin == 0 && *slot != 1 {
						t.Fatal("loss remembered")
					}
					if begin != 0 && *slot != seen {
						t.Fatal("begin must preserve remembered loss")
					}
					if memmap.Uint32(0x5D4594, 832472) != 0 || *o.briefWords["dword_5d4594_831260"] != 1 {
						t.Fatal("chapter modal state")
					}
					rows = append(rows, o.windowCapture(t, 7, v))
				}
			}
		}
	}
	briefingWindowCapture(t, "chapters", rows, "a4df6e50ccc7ed6b1e4ef82e4aa1f846161932d7bc088d80a608660ee63a47da")
}
func TestBriefingWindowQuestModePrecedence(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for mode := uintptr(0); mode < 8; mode++ {
		for _, caption := range []bool{false, true} {
			o.resetWindow(t)
			o.constructBriefing(t)
			legacy.PortTestBriefing(9, 0, 0, 0)
			if caption {
				legacy.PortTestBriefing(10, txptr(unsafe.Pointer(alloc.InternCString16("A custom briefing."))), 0, 0)
			}
			v := legacy.PortTestBriefingWindow(7, 254, 1, mode, 0)
			want := uint32(2)
			if mode&1 != 0 {
				want = 1
			} else if mode&4 != 0 {
				want = 4
			}
			if v != 0 || memmap.Uint32(0x5D4594, 832472) != want || *o.briefWords["dword_5d4594_831240"] != 0 || *o.briefWords["dword_5d4594_831224"] != 0 {
				t.Fatal("quest precedence/voice")
			}
			rows = append(rows, o.windowCapture(t, 7, v))
		}
	}
	briefingWindowCapture(t, "quest-modes", rows, "bac57ea5e5d5d083c086746c0b0d835d4e0943234f01a6bded878cb5ad2c7b16")
}
func TestBriefingWindowCreditsSelection(t *testing.T) {
	o := newBriefingWindowOwner(t)
	var rows []briefingWindowResult
	for class := 0; class < 3; class++ {
		for _, ending := range []byte{0, 1, 2, 4, 7} {
			o.resetWindow(t)
			o.constructBriefing(t)
			dword_587000_311372 = class
			dword_5d4594_2516476 = ending
			want := nox_xxx_GetEndgameDialog()
			v := legacy.PortTestBriefingWindow(7, 255, 1, 0, 0)
			got := alloc.GoString((*byte)(unsafe.Pointer(uintptr(*o.briefWords["dword_5d4594_831240"]))))
			if v != 1 || got != want || *o.briefWords["dword_5d4594_831220"] != 255 || *o.briefWords["dword_5d4594_831276"] != 1140457472 {
				t.Fatal("credits setup")
			}
			rows = append(rows, o.windowCapture(t, 7, v))
		}
	}
	briefingWindowCapture(t, "credits", rows, "ad2b70eae228be79a0e24321cdc8c508b4493e92f0ab0be37b9fcf17cc7ee03c")
}
