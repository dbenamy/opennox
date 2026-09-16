//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"testing"
	"unsafe"
)

func TestSpellbookVisibility(t *testing.T) {
	o := newSpellbookOwner(t)
	var rows []spellbookResult
	for _, guide := range []uint32{0, 1} {
		for _, reset := range []uint32{0, 1} {
			o.resetBook(t)
			if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
				t.Fatal("book setup")
			}
			*o.words["dword_5d4594_1046872"] = guide
			label := fmt.Sprintf("guide%d-reset%d", guide, reset)
			record := func(step string, ret uint32) { rows = append(rows, o.bookSnapshot(label+"-"+step, ret)) }
			if o.bookCall("sub_45CFC0") != 0 {
				t.Fatal("initial book hidden")
			}
			if o.bookCall("nox_xxx_bookHideMB_45ACA0", reset) != 0 {
				t.Fatal("already hidden")
			}
			for cycle := 0; cycle < 3; cycle++ {
				o.bookCall("nox_xxx_book_45B010", 1)
				if o.bookCall("sub_45CFC0") != 1 || o.bookWindow().GetFlags().IsHidden() {
					t.Fatal("forced show")
				}
				record(fmt.Sprintf("show%d", cycle), 1)
				if !o.bookWindow().Capture(true) {
					t.Fatal("book capture")
				}
				*o.words["nox_xxx_aNox_cfg_0_587000_132132"] = 0
				*o.words["dword_5d4594_1046936"] = 2
				*o.words["dword_5d4594_1047520"] = 1
				n := len(o.sounds)
				if o.bookCall("nox_xxx_bookHideMB_45ACA0", reset) != 0 || len(o.sounds) != n || o.c.GUI.Captured() != o.bookWindow() {
					t.Fatal("active addition must prevent hiding")
				}
				record(fmt.Sprintf("blocked%d", cycle), 0)
				*o.words["dword_5d4594_1047520"] = 0
				ret := o.bookCall("nox_xxx_bookHideMB_45ACA0", reset)
				if ret != 1 || o.bookCall("sub_45CFC0") != 0 || o.c.GUI.Captured() != nil || *o.words["dword_5d4594_1046868"] != guide {
					t.Fatal("hide releases capture")
				}
				if reset != 0 && (*o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 1 || *o.words["dword_5d4594_1046936"] != 0) {
					t.Fatal("hide resets contents when requested")
				}
				record(fmt.Sprintf("hide%d", cycle), ret)
			}
			// Hiding a book must retain another window's capture.
			o.bookCall("nox_xxx_book_45B010", 1)
			other := (*gui.Window)(unsafe.Pointer(uintptr(*o.words["dword_5d4594_1046944"])))
			if other == o.bookWindow() {
				t.Fatal("independent capture fixture")
			}
			if !other.Capture(true) {
				t.Fatal("other capture")
			}
			if o.bookCall("nox_xxx_bookHideMB_45ACA0", 0) != 1 || o.c.GUI.Captured() != other {
				t.Fatal("unrelated capture retained")
			}
			record("other-capture", uint32(uintptr(unsafe.Pointer(other))))
			other.Capture(false)
		}
	}
	spellbookCapture(t, "visibility", rows, "")
}
