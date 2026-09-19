//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
	"unsafe"
)

func unitNameStrings(t *testing.T, s *server.Server, entries []strman.Entry) {
	configure, restore := s.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
}
func unitNameEntry(id, value string) strman.Entry {
	return strman.Entry{ID: strman.ID(id), Vals: []strman.Variant{{Str: value}}}
}
func TestUnitGameplayNPCNames(t *testing.T) {
	s := newObjectXferOwner(t)
	u := newObjectXferSimple(t, s)
	typeName := s.Types.ByInd(int(u.TypeInd)).ID()
	inputs := [][2]string{{"", ""}, {"_", ""}, {"a_b", "ab"}, {"first:second:a_b", "ab"}, {"prefix:", ""}, {"A__B_", "AB"}, {"Ren_é", "René"}, {"a\x00ignored", "a"}, {"NPC:Port_Trade_Vendor", "PortTradeVendor"}}
	var entries []strman.Entry
	seen := map[string]bool{}
	for _, row := range append(append([][2]string{}, inputs...), [2]string{typeName, typeName}) {
		key := "NPC:" + row[1]
		if !seen[strings.ToLower(key)] {
			entries = append(entries, unitNameEntry(key, "LOCALIZED:"+strings.ToLower(row[1])))
			seen[strings.ToLower(key)] = true
		}
	}
	unitNameStrings(t, s, entries)
	buf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1563460)), 512)
	old := bytes.Clone(buf)
	t.Cleanup(func() { copy(buf, old) })
	var rows [][2]string
	for i := 0; i <= len(inputs); i++ {
		want := typeName
		u.IDPtr = nil
		free := func() {}
		if i < len(inputs) {
			p, f := alloc.CString(inputs[i][0])
			free = f
			u.IDPtr = unsafe.Pointer(p)
			want = inputs[i][1]
		}
		for j := range buf {
			buf[j] = 0xa5
		}
		got := legacy.PortTestUnitNPCName(u)
		key := alloc.GoString(&buf[0])
		u.IDPtr = nil
		free()
		if got != "LOCALIZED:"+strings.ToLower(want) || key != "NPC:"+want {
			t.Fatalf("NPC case%d: %q/%q want %q", i, got, key, want)
		}
		for _, b := range buf[len(key)+1:] {
			if b != 0xa5 {
				t.Fatal("NPC scratch tail changed")
			}
		}
		rows = append(rows, [2]string{key, got})
	}
	spellbookCapture(t, "unit-gameplay-npc-names", rows, "6070544e573ec745577b9d3e7d211644d93abccaec4c93a470f0799455edff8c")
}
func TestUnitGameplayItemNames(t *testing.T) {
	s := newObjectXferOwner(t)
	u := newObjectXferSimple(t, s)
	typeName := s.Types.ByInd(int(u.TypeInd)).ID()
	unitNameStrings(t, s, []strman.Entry{unitNameEntry("objutil.c:NoInfo", "NO_INFO:%S"), unitNameEntry("objutil.c:NoDescription", "NO_DESCRIPTION")})
	buf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565660)), 2076)
	old := bytes.Clone(buf)
	t.Cleanup(func() { copy(buf, old) })
	oldMods := s.Modif
	t.Cleanup(func() { s.Modif = oldMods })
	def, free := alloc.New(server.Modifier{})
	t.Cleanup(free)
	def.TypeInd = uint32(u.TypeInd)
	mods, freeMods := alloc.New([4]*server.ModifierEff{})
	t.Cleanup(freeMods)
	u.InitData = unsafe.Pointer(mods)
	t.Cleanup(func() { u.InitData = nil })
	var effects [4]*server.ModifierEff
	for i := range effects {
		var f func()
		effects[i], f = alloc.New(server.ModifierEff{})
		t.Cleanup(f)
	}
	text := func(s string) *uint16 { p, f := alloc.CString16(s); t.Cleanup(f); return p }
	words := [4]string{"Ancient", "银", "of Ice", "the Swift"}
	var descriptions [4]*uint16
	for i := range descriptions {
		descriptions[i] = text(words[i])
	}
	empty, wrong, base := text(""), text("WRONG_FIELD"), text("Blade")
	var rows []string
	for _, class := range []uint32{0, 0x1000, 0x1000000, 0x2000000, 0x10000000, 0x12000000} {
		u.ObjClass = object.Class(class)
		for mode := 0; mode < 4; mode++ {
			s.Modif.Dword_5d4594_251600, s.Modif.Dword_5d4594_251608 = def, def
			def.Desc8 = base
			if mode == 1 {
				def.Desc8 = nil
			}
			if mode == 2 {
				def.Desc8 = empty
			}
			if mode == 3 {
				s.Modif.Dword_5d4594_251600, s.Modif.Dword_5d4594_251608 = nil, nil
			}
			for mask := 0; mask < 16; mask++ {
				for nameMode := 0; nameMode < 3; nameMode++ {
					for prefix := 0; prefix < 2; prefix++ {
						clear(buf)
						initial := ""
						if prefix != 0 {
							initial = "~"
							*memmap.PtrUint16(0x5D4594, 1567732) = '~'
						}
						want := initial
						for i := range effects {
							mods[i] = nil
							if mask&(1<<i) != 0 {
								mods[i] = effects[i]
							}
							primary := descriptions[i]
							if nameMode == 1 {
								primary = nil
							}
							if nameMode == 2 {
								primary = empty
							}
							*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(effects[i]), 8)) = unsafe.Pointer(primary)
							*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(effects[i]), 12)) = unsafe.Pointer(wrong)
							if i == 3 {
								*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(effects[i]), 8)) = unsafe.Pointer(wrong)
								*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(effects[i]), 12)) = unsafe.Pointer(primary)
							}
						}
						for i := 0; i < 2; i++ {
							if mask&(1<<i) != 0 && nameMode != 1 {
								if nameMode == 0 {
									want += words[i]
								}
								want += " "
							}
						}
						if mode == 0 {
							want += "Blade"
						}
						for i := 2; i < 4; i++ {
							if mask&(1<<i) != 0 && nameMode != 1 {
								want += " "
								if nameMode == 0 {
									want += words[i]
								}
							}
						}
						if mode == 3 {
							want = "NO_INFO:" + typeName
						}
						if class == 0 {
							want = "NO_DESCRIPTION"
						}
						got, ptr := legacy.PortTestUnitItemName(u)
						if got != want || ptr != unsafe.Pointer(&buf[0]) {
							t.Fatal(fmt.Sprintf("item class%x mode%d mask%x name%d prefix%d: %q want %q", class, mode, mask, nameMode, prefix, got, want))
						}
						rows = append(rows, got)
					}
				}
			}
		}
	}
	spellbookCapture(t, "unit-gameplay-item-names", rows, "92563d27148f547bbaa35e00d4599f6264503b3e8811015a443182b2880f3ab7")
	t.Logf("%d item name cases", len(rows))
}
