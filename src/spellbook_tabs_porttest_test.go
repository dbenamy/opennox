//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSpellbookTabSelection(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	guides := unsafe.Slice(memmap.PtrUint32(0x5D4594, 740076), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	guides[8] = 1
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
	var rows []spellbookResult
	for initial := uint32(0); initial < 2; initial++ {
		for target := uint32(0); target < 2; target++ {
			for knowledge := 0; knowledge < 4; knowledge++ {
				for detail := uint32(0); detail < 2; detail++ {
					for page := uint32(0); page < 2; page++ {
						o.resetBook(t)
						if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
							t.Fatal("book setup")
						}
						*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = 1
						if knowledge&1 != 0 {
							*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3700)) = 1
						}
						if knowledge&2 != 0 {
							*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4248)) = 1
						}
						*o.words["dword_5d4594_1046868"] = initial
						*o.words["dword_5d4594_1046872"] = initial
						*o.words["nox_xxx_aNox_cfg_0_587000_132132"] = 1 - detail
						*o.words["dword_5d4594_1046936"] = page
						o.bookCall("nox_xxx_guiSpellSortList_45ADF0", 1)
						o.collect()
						var tab uint32
						for _, w := range o.windows {
							if *(*uint32)(w.C()) == 1310+10*target {
								tab = uint32(uintptr(w.C()))
								break
							}
						}
						if tab == 0 {
							t.Fatal("actual tab window missing")
						}
						ret := o.bookCall("nox_xxx_bookChildWndProcMB_45B360", tab, 5)
						want := target
						unavailable := initial != target && knowledge&(1<<target) == 0
						if unavailable {
							want = initial
						}
						if ret != 1 || *o.words["dword_5d4594_1046872"] != want {
							t.Fatal("tab selection must retain previous view if target is empty")
						}
						if !unavailable && (*o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 1 || *o.words["dword_5d4594_1046936"] != 0) {
							t.Fatal("tab returns to first contents page")
						}
						label := fmt.Sprintf("initial%d-target%d-known%d-detail%d-page%d", initial, target, knowledge, detail, page)
						rows = append(rows, o.bookSnapshot(label, ret))
						for _, event := range []uint32{0, 4, 6, 7, 8, 99} {
							ret = o.bookCall("nox_xxx_bookChildWndProcMB_45B360", tab, event)
							wantRet := uint32(0)
							if event == 6 || event == 7 {
								wantRet = 1
							}
							if ret != wantRet {
								t.Fatal("tab event handling")
							}
						}
						*o.words["dword_5d4594_1047520"] = 1
						before := *o.words["dword_5d4594_1046872"]
						if o.bookCall("nox_xxx_bookChildWndProcMB_45B360", tab, 5) != 1 || *o.words["dword_5d4594_1046872"] != before {
							t.Fatal("addition blocks tab switching")
						}
						rows = append(rows, o.bookSnapshot(label+"-blocked", 1))
					}
				}
			}
		}
	}
	spellbookCapture(t, "tabs", rows, "")
}
