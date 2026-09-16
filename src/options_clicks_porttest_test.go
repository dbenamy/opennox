//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestOptionsBasicClicks(t *testing.T) {
	type record struct {
		Menu                   bool
		Event, ID, Return, Cut int
		Bits                   uint32
		Sounds                 [][2]int
	}
	var rows []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprint(menu), func(t *testing.T) {
			o := newOptionsAudioOwner(t, menu)
			oldPix := noxPixBuffer.img
			noxPixBuffer.img = o.pix
			t.Cleanup(func() { noxPixBuffer.img = oldPix })
			oldDirty := legacy.Get_dword_5d4594_1193188()
			t.Cleanup(func() { legacy.Set_dword_5d4594_1193188(oldDirty) })
			for _, event := range []int{16389, 16391} {
				for _, id := range []int{311, 312, 313, 314, 331, 332, 333, 334, 999} {
					nox_video_cutSize = 80
					*o.words[172880] = 32
					o.sounds = nil
					ret := legacy.PortTestOptionsEvent(menu, o.root, event, o.controls[id], 0)
					wantCut, wantBits := 80, uint32(32)
					wantSound := 920
					if event == 16391 {
						wantSound = 921
						if id >= 311 && id <= 314 {
							wantCut = []int{65, 75, 85, 100}[id-311]
						}
						if menu && id == 331 {
							wantBits = 8
						}
						if menu && id == 332 {
							wantBits = 16
						}
					}
					if ret != 1 || nox_video_cutSize != wantCut || *o.words[172880] != wantBits || !reflect.DeepEqual(o.sounds, [][2]int{{wantSound, 100}}) {
						t.Fatalf("menu=%v event=%d id=%d return=%d cut=%d bits=%d sounds=%v", menu, event, id, ret, nox_video_cutSize, *o.words[172880], o.sounds)
					}
					rows = append(rows, record{menu, event, id, ret, nox_video_cutSize, *o.words[172880], append([][2]int(nil), o.sounds...)})
				}
			}
		})
	}
	spellbookCapture(t, "options-clicks", rows, "06d6693189c0db0a8a7d121af0cd3bc251c4e1c7bd481c71032144ed61c8a000")
}
