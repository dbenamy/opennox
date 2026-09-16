//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
)

func (q *quickbarOwner) ownFullQuickbar(t *testing.T) {
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if strings.HasPrefix(name, "QuickBar") || strings.HasPrefix(name, "SpellbookButton") {
			q.loads = append(q.loads, name)
			index := 0
			for _, b := range []byte(name) {
				index = (index + int(b)) % len(q.images)
			}
			return q.images[index]
		}
		return oldLoad(name)
	}
	for _, p := range blobdata.PortTestQuickbarConstants() {
		dst := unsafe.Slice(memmap.PtrUint8(0x587000, p.Offset), len(p.Data))
		saved := append([]byte(nil), dst...)
		t.Cleanup(func() { copy(dst, saved) })
		copy(dst, p.Data)
	}
	var entries []strman.Entry
	for _, name := range []string{"CloseSpellbookTT", "NoIcon", "OneSpellPerTrap", "OpenSpellBookTT", "OpenSpellbookTT", "RestrictedTrapSpell", "SpellSet", "ToolTipAllSpellSets", "ToolTipCastAtOther", "ToolTipCastOnMe", "ToolTipLayTrap", "ToolTipNextSpellSet", "ToolTipNextTrap", "ToolTipPrevSpellSet", "ToolTipPrevTrap", "ToolTipSummonBomber", "ToolTipTrapConstruct", "TrapError", "TrapSet"} {
		text := name
		if name == "SpellSet" || name == "TrapSet" {
			text += " %d"
		}
		entries = append(entries, strman.Entry{ID: strman.ID("guispell.c:" + name), Vals: []strman.Variant{{Str: text}}})
	}
	configure, restore := q.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
}
func TestQuickbarLifecycle(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var records []quickbarResult
	for _, width := range []uint32{640, 1280} {
		for _, class := range []uint32{0, 1, 2, 3} {
			for _, mode := range []uint32{0, 0x2000, 0x2000 | 4096, uint32(noxflags.GameModeQuest)} {
				q.reset(t)
				*q.words["nox_win_width"] = width
				*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = byte(class)
				if class == 3 {
					*q.words["dword_8531A0_2576"] = 0
				}
				noxflags.SetGame(noxflags.GameFlag(mode))
				q.collect()
				before := len(q.windows)
				ret := q.call("nox_xxx_quickBarCreate_45E190")
				q.collect()
				if class == 3 {
					q.check(t, ret == 0 && len(q.windows) == before, "missing player creates no quickbar windows")
				} else {
					q.check(t, ret == 1 && len(q.windows) > before && q.call("sub_460D40") == 1, "quickbar lifecycle creates real windows")
					q.check(t, *q.quickWords["dword_5d4594_1047548"] == (width-320)/2, "quickbar centered at target width")
				}
				label := fmt.Sprintf("width%d-class%d-mode%d", width, class, mode)
				records = append(records, q.snapshot(label+"-create", ret))
				if class != 3 {
					q.call("sub_460B90", 0)
					records = append(records, q.snapshot(label+"-hide", 0))
					q.call("sub_460B90", 1)
					records = append(records, q.snapshot(label+"-show", 0))
					// The C destructor leaves these two child addresses in storage.
					// Preserve identities while the actual windows are still live.
					for _, v := range []uint32{*q.quickWords["dword_5d4594_1049524"], memmap.Uint32(0x5D4594, 1049528)} {
						q.c.dataRefs[v] = q.normalize(v)
					}
					q.call("sub_460D50")
					q.check(t, q.call("sub_460D40") == 0 && *q.quickWords["dword_5d4594_1049532"] == 0, "destroy clears quickbar/capture owner")
					records = append(records, q.snapshot(label+"-destroy", 0))
				}
			}
		}
	}
	spellbookCapture(t, "quickbar-lifecycle", records, "88a114c61e0fb134ffc3eabff07fe1fbd5fd04a8496c3ccb7b9aa01440f7e523")
}
