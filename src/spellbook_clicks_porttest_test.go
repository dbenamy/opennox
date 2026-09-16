//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
	"unsafe"
)

func TestSpellbookListClicks(t *testing.T) {
	o := newSpellbookOwner(t)
	var rows []spellbookResult
	for class := 0; class < 3; class++ {
		for _, guide := range []uint32{0, 1} {
			for page := 0; page < 3; page++ {
				for _, summon := range []uint32{0, 1} {
					for _, point := range [][2]int{{145, 18}, {145, 19}, {146, 19}, {145, 31}, {145, 32}, {146, 32}, {145, 135}, {146, 135}, {145, 136}, {146, 136}, {0, 19}, {284, 19}} {
						o.resetBook(t)
						if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
							t.Fatal("book setup")
						}
						*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
						*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4232)) = summon
						*o.words["dword_5d4594_1046868"] = guide
						*o.words["dword_5d4594_1046872"] = guide
						*o.words["dword_5d4594_1046936"] = uint32(page)
						*memmap.PtrUint32(0x5D4594, 1047508) = 37
						*memmap.PtrUint32(0x5D4594, 1046940) = 3
						for i := 0; i < 37; i++ {
							*memmap.PtrUint32(0x5D4594, 1046960+4*uintptr(i)) = uint32(i%5 + 1)
						}
						x, y := point[0]+5, point[1]+157
						packed := uint32(uint16(x)) | uint32(uint16(y))<<16
						ret := o.bookCall("nox_xxx_bookListWndProc_45B5F0", *o.words["nox_win_unk1"], 5, packed)
						index := -1
						if point[1] >= 19 && point[1] < 136 {
							index = page*18 + (point[1]-19)/13
							if point[0] > 145 {
								index += 9
							}
							if index >= 37 {
								index = -1
							}
						}
						immediate := index >= 0 && class == 2 && guide == 1 && summon == 0
						want := uint32(index)
						if immediate {
							want = ^uint32(0)
						}
						if ret != 1 || *o.words["nox_xxx_aNox_cfg_0_587000_132136"] != want {
							t.Fatalf("selection class%d guide%d page%d point%v: got%d want%d", class, guide, page, point, *o.words["nox_xxx_aNox_cfg_0_587000_132136"], want)
						}
						label := fmt.Sprintf("class%d-guide%d-page%d-summon%d-x%d-y%d", class, guide, page, summon, point[0], point[1])
						rows = append(rows, o.bookSnapshot(label+"-down", ret))
						ret = o.bookCall("nox_xxx_bookListWndProc_45B5F0", *o.words["nox_win_unk1"], 6, packed)
						if ret != 1 || o.c.dragndropSpellType != 0 || o.c.GUI.Captured() != nil {
							t.Fatal("click release clears drag and capture")
						}
						if index >= 0 && (*o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 0 || *o.words["dword_5d4594_1046932"] != uint32(index)) {
							t.Fatal("click opens selected detail page")
						}
						rows = append(rows, o.bookSnapshot(label+"-up", ret))
					}
				}
			}
		}
	}
	spellbookCapture(t, "clicks", rows, "fec3eae1cf267d9ad253bf4819a05bf98e414c12ad4f2507f9b7d27ba607ec6c")
}
