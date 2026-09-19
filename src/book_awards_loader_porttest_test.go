//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
)

func TestBookAwardsGuideLoader(t *testing.T) {
	o := newSpellbookOwner(t)
	names := make([]string, 41)
	lookup := bookAwardWords(t, 0x587000, 70500, 41)
	table := bookAwardWords(t, 0x5D4594, 740076, 41*7)
	var entries []strman.Entry
	for i := range names {
		names[i] = fmt.Sprintf("Guide%02d", i)
		lookup[i] = bookAwardString(t, names[i])
		entries = append(entries, strman.Entry{ID: strman.ID("creature:" + names[i]), Vals: []strman.Variant{{Str: fmt.Sprintf("Name Ω %d", i)}}}, strman.Entry{ID: strman.ID("creature_desc:" + names[i]), Vals: []strman.Variant{{Str: fmt.Sprintf("Description é %d", i)}}})
	}
	configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	cagePrefix := bookAwardWords(t, 0x587000, 71248, 2)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&cagePrefix[0])), 8), "Creature")
	bookPrefix := bookAwardWords(t, 0x587000, 71264, 2)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&bookPrefix[0])), 8), "GuideBoo")
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	var rows []map[string]any
	for _, missing := range []int{-1, 1, 20, 40} {
		for _, missingImages := range []bool{false, true} {
			restoreTypes := o.c.Cli().PortTestBookGuideTypes(names, missing)
			t.Cleanup(restoreTypes)
			for i := range table {
				table[i] = 0xcccccccc
			}
			before := append([]uint32(nil), table...)
			var calls []string
			legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
				calls = append(calls, name)
				if missingImages {
					return nil
				}
				if len(calls)%2 == 1 {
					return o.images[0]
				}
				return o.images[1]
			}
			ret := bookAwardCall("nox_xxx_loadGuides_427070")
			end := 41
			wantRet := uint32(1)
			if missing >= 1 {
				end = missing
				wantRet = 0
			}
			if ret != wantRet {
				t.Fatal("loader result", missing, ret, wantRet)
			}
			if len(calls) != (end-1)*2 {
				t.Fatal("loader requests", missing, len(calls))
			}
			var state []map[string]any
			for i := 1; i < end; i++ {
				row := table[7*i : 7*i+7]
				name := legacy.GoWStringP(unsafe.Pointer(uintptr(row[0])))
				desc := legacy.GoWStringP(unsafe.Pointer(uintptr(row[2])))
				typ := uint32(1000 + i)
				if i == 7 {
					typ = 0
				}
				size := uint32(4)
				if i%4&1 != 0 {
					size = 1
				} else if i%4&2 != 0 {
					size = 2
				}
				cage, book := uint32(0), uint32(0)
				if !missingImages {
					cage = uint32(uintptr(o.images[0].C()))
					book = uint32(uintptr(o.images[1].C()))
				}
				if name != fmt.Sprintf("Name Ω %d", i) || desc != fmt.Sprintf("Description é %d", i) || row[1] != typ || row[3] != cage || row[4] != book || row[5] != 0 || row[6] != (0xcccccc00|size) {
					t.Fatal("loader row", missing, missingImages, i, name, desc, row)
				}
				if calls[2*(i-1)] != "CreatureCage"+names[i] || calls[2*(i-1)+1] != "GuideBook"+names[i] {
					t.Fatal("image names", calls)
				}
				state = append(state, map[string]any{"id": i, "name": name, "desc": desc, "type": row[1], "sizeAndPadding": row[6], "images": !missingImages})
			}
			if !reflect.DeepEqual(table[:7], before[:7]) || !reflect.DeepEqual(table[7*end:], before[7*end:]) {
				t.Fatal("loader touched unvisited rows")
			}
			rows = append(rows, map[string]any{"missing": missing, "missingImages": missingImages, "return": ret, "calls": calls, "rows": state})
		}
	}
	spellbookCapture(t, "book-awards-guide-loader", rows, "a81fad93f689996661ea885cef1f3200d8e8f1f90d7416cd905de18f072e4716")
}
