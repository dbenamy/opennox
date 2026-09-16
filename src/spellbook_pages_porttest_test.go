//go:build porttest

package opennox

import (
	"fmt"
	"testing"

	"github.com/opennox/opennox/v1/common/memmap"
)

func TestSpellbookPageTransitions(t *testing.T) {
	o := newSpellbookOwner(t)
	var rows []spellbookResult
	for _, guide := range []uint32{0, 1} {
		for _, count := range []uint32{1, 2, 17, 18, 19, 36, 37} {
			for _, special := range []uint32{0, 1} {
				if special >= count {
					continue
				}
				o.resetBook(t)
				if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
					t.Fatal("book setup")
				}
				*o.words["dword_5d4594_1046872"] = guide
				*o.words["dword_5d4594_1046868"] = guide
				*o.words["dword_5d4594_1047512"] = special
				*memmap.PtrUint32(0x5D4594, 1047508) = count
				pages := (count-special)/18 + 1
				*memmap.PtrUint32(0x5D4594, 1046940) = pages
				label := fmt.Sprintf("guide%d-count%d-special%d", guide, count, special)
				record := func(step string, ret uint32) { rows = append(rows, o.bookSnapshot(label+"-"+step, ret)) }
				finish := func() {
					op := "nox_xxx_bookClickSpell_45B1F0"
					if guide != 0 {
						op = "nox_xxx_bookClickCreature_45B200"
					}
					if ret := o.bookCall(op); ret != guide || *o.words["dword_5d4594_1046868"] != guide {
						t.Fatal("page animation completion")
					}
				}
				for page := uint32(1); page < pages; page++ {
					ret := o.bookCall("nox_xxx_bookWndProc_45B070", 0, 5)
					if ret != 1 || *o.words["dword_5d4594_1046936"] != page || *o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 1 {
						t.Fatal("next contents page")
					}
					record(fmt.Sprintf("contents%d", page), ret)
					finish()
				}
				for entry := uint32(0); entry < count-special; entry++ {
					ret := o.bookCall("nox_xxx_bookWndProc_45B070", 0, 5)
					if ret != 1 || *o.words["dword_5d4594_1046932"] != entry || *o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 0 {
						t.Fatal("next entry")
					}
					record(fmt.Sprintf("entry%d", entry), ret)
					finish()
				}
				n := len(o.sounds)
				if ret := o.bookCall("nox_xxx_bookWndProc_45B070", 0, 5); ret != 1 || len(o.sounds) != n {
					t.Fatal("last entry must not turn")
				}
				record("last-boundary", 1)
				for entry := int(count-special) - 2; entry >= 0; entry-- {
					ret := o.bookCall("nox_xxx_book_45B210", 0, 5)
					if ret != 1 || *o.words["dword_5d4594_1046932"] != uint32(entry) {
						t.Fatal("previous entry")
					}
					record(fmt.Sprintf("back-entry%d", entry), ret)
					finish()
				}
				ret := o.bookCall("nox_xxx_book_45B210", 0, 5)
				if ret != 1 || *o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 1 || *o.words["dword_5d4594_1046936"] != pages-1 {
					t.Fatal("return to contents")
				}
				record("return-contents", ret)
				finish()
				for page := int(pages) - 2; page >= 0; page-- {
					ret = o.bookCall("nox_xxx_book_45B210", 0, 5)
					if ret != 1 || *o.words["dword_5d4594_1046936"] != uint32(page) {
						t.Fatal("previous contents page")
					}
					record(fmt.Sprintf("back-contents%d", page), ret)
					finish()
				}
				n = len(o.sounds)
				if ret = o.bookCall("nox_xxx_book_45B210", 0, 5); ret != 1 || len(o.sounds) != n {
					t.Fatal("first page must not turn")
				}
				record("first-boundary", ret)
				for _, op := range []string{"nox_xxx_book_45B210", "nox_xxx_bookWndProc_45B070"} {
					if ret = o.bookCall(op, 0, 99); ret != 0 {
						t.Fatal("unhandled event")
					}
					*o.words["dword_5d4594_1047520"] = 1
					if ret = o.bookCall(op, 0, 5); ret != 1 || len(o.sounds) != n {
						t.Fatal("addition animation blocks page turns")
					}
					record(op+"-blocked", ret)
					*o.words["dword_5d4594_1047520"] = 0
				}
				ret = o.bookCall("nox_xxx_bookMoveToPage_45B930", count-special-1)
				if *o.words["dword_5d4594_1046932"] != count-special-1 || *o.words["dword_5d4594_1046936"] != 99 || *o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 0 {
					t.Fatal("direct page selection")
				}
				record("direct", ret)
			}
		}
	}
	spellbookCapture(t, "pages", rows, "0ee40a5e4a73ca8bdb7844497213a7d941cfec30396462d2c4e4993e4bcce8d7")
}
