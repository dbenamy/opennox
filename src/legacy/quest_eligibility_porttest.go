//go:build porttest

package legacy

/*
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1568308;
static uint32_t questEligibilityInvoke(int op,uint32_t a) {
 switch(op) {
 case 0:return sub_4F24E0(a);
 case 1:return sub_4F2530(a);
 case 2:return sub_4F2570(a);
 case 3:return sub_4F2590(a);
 case 4:return sub_4F2700(a);
 case 5:return sub_4F27A0(a);
 case 6:return sub_4F27E0(a);
 case 7:return sub_4F28C0(a);
 case 8:return sub_4F2960(a);
 case 9:return sub_4F2B20(a);
 case 10:return sub_4F2B60(a);
 case 11:return sub_4F2C30(a);
 case 12:return nox_xxx_spell_4F2E70(a);
 case 13:return sub_4F2EF0(a);
 default:return 0xDEADBEEF;
 }
}
*/
import "C"

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestEligibilityObject struct {
	Type                   string
	Class, Subclass, Flags uint32
	Mods                   [4]int
	Book                   string
	BookID                 byte
	Inventory              []int
}
type PortTestEligibilityModifier struct {
	Name                         string
	WeaponMask, ArmorMask, Flags uint32
}
type PortTestEligibilityModRow struct {
	Ref                         int
	Name                        *string
	ArmorExclude, WeaponExclude uint32
}
type PortTestEligibilityItemRow struct {
	Name *string
	Type string
}
type PortTestQuestEligibilitySpec struct {
	Objects                                         []PortTestEligibilityObject
	Modifiers                                       []PortTestEligibilityModifier // Appended after the four replenishment modifiers.
	SpecialFlags                                    [4]uint32
	SpellTable, BeastTable                          [][2]uint32
	BeastGroups                                     [][]uint32
	ItemTable                                       []PortTestEligibilityItemRow
	MaterialWeapon, MaterialArmor, Quality, Effects []PortTestEligibilityModRow
	WeaponBit, ArmorBit                             uint32
	Warm                                            bool
}
type PortTestQuestEligibilityResult struct {
	Globals   []uint32
	Objects   [][]uint32
	Modifiers [][]uint32
	Tables    [][]uint32
}
type portTestQuestEligibility struct {
	objects  []*server.Object
	mods     []*server.ModifierEff
	tables   [][]byte
	restores []func()
}

var eligibilityNames = []string{"Diamond", "Emerald", "Ruby", "SulphorousFlareWand", "StreetSneakers", "StreetShirt", "StreetPants", "RedPotion", "BluePotion", "CurePoisonPotion", "HastePotion", "InvisibilityPotion", "ShieldPotion", "VampirismPotion", "FireProtectPotion", "ShockProtectPotion", "PoisonProtectPotion", "InvulnerabilityPotion", "InfravisionPotion", "InfinitePainWand", "EligibilityWeapon", "EligibilityArmor", "EligibilityBook", "EligibilityPlayer"}

func (p *portTestShopPools) eligibilityEnsure() {
	if p.reportEligibility != nil || p.reports == nil || p.reports.Eligibility == nil {
		return
	}
	sp := p.reports.Eligibility
	st := &portTestQuestEligibility{}
	p.reportEligibility = st
	core := p.proxy.core
	st.restores = append(st.restores, core.PortTestEligibilityTypes(eligibilityNames, sp.ArmorBit, sp.WeaponBit))
	special := []string{"Replenishment1", "Replenishment2", "Replenishment3", "Replenishment4"}
	var modNames []*byte
	defs := make([]PortTestEligibilityModifier, 4)
	for i, n := range special {
		defs[i] = PortTestEligibilityModifier{Name: n, Flags: sp.SpecialFlags[i]}
	}
	defs = append(defs, sp.Modifiers...)
	for i, d := range defs {
		ptr := p.objectiveRegion(144)
		clear(unsafe.Slice((*byte)(ptr), 144))
		m := (*server.ModifierEff)(ptr)
		st.mods = append(st.mods, m)
		modNames = append(modNames, (*byte)(p.objectiveString(d.Name)))
		*equipmentWord(ptr, 28) = d.WeaponMask
		*equipmentWord(ptr, 32) = d.ArmorMask
		*equipmentWord(ptr, 36) = d.Flags
		p.identify(ptr, 950000+uint32(i))
	}
	st.restores = append(st.restores, core.PortTestControlsModifiers(st.mods, modNames))
	modRef := func(ref int) unsafe.Pointer {
		if ref == 0 {
			return nil
		}
		if ref < 1 || ref > len(st.mods) {
			panic("eligibility modifier reference")
		}
		return st.mods[ref-1].C()
	}
	for i, d := range sp.Objects {
		ptr := p.objectiveRegion(772)
		clear(unsafe.Slice((*byte)(ptr), 772))
		u := (*server.Object)(ptr)
		*(*uint16)(unsafe.Add(ptr, 4)) = uint16(core.Types.IndByID(d.Type))
		*equipmentWord(ptr, 8) = d.Class
		*equipmentWord(ptr, 12) = d.Subclass
		*equipmentWord(ptr, 16) = d.Flags
		init := p.objectiveRegion(16)
		for j, ref := range d.Mods {
			*(*unsafe.Pointer)(unsafe.Add(init, j*4)) = modRef(ref)
		}
		*(*unsafe.Pointer)(unsafe.Add(ptr, 692)) = init
		book := p.objectiveString(d.Book)
		if d.Subclass&1 != 0 || d.Subclass&2 == 0 {
			*(*byte)(book) = d.BookID
		}
		*(*unsafe.Pointer)(unsafe.Add(ptr, 736)) = book
		p.identify(ptr, 960000+uint32(i))
		st.objects = append(st.objects, u)
	}
	objRef := func(ref int) *server.Object {
		if ref == 0 {
			return nil
		}
		if ref < 1 || ref > len(st.objects) {
			panic("eligibility object reference")
		}
		return st.objects[ref-1]
	}
	for i, d := range sp.Objects {
		var first *server.Object
		for j := len(d.Inventory) - 1; j >= 0; j-- {
			u := objRef(d.Inventory[j])
			u.InvNextItem = first
			first = u
		}
		st.objects[i].InvFirstItem = first
	}
	// Save each complete bounded table before changing its rows or sentinel.
	region := func(base, off uintptr, size int) []byte {
		b := unsafe.Slice((*byte)(memmap.PtrOff(base, off)), size)
		old := append([]byte(nil), b...)
		clear(b)
		st.tables = append(st.tables, b)
		st.restores = append(st.restores, func() { copy(b, old) })
		return b
	}
	guides := region(0x587000, 70500, 41*4)
	for i := 0; i < 41; i++ {
		*(*unsafe.Pointer)(unsafe.Pointer(&guides[i*4])) = p.objectiveString(fmt.Sprintf("EligibilityGuide%02d", i))
	}
	for i, rows := range [][][2]uint32{sp.SpellTable, sp.BeastTable} {
		off, capacity := uintptr(207108), 36
		if i == 1 {
			off, capacity = 207796, 16
		}
		if len(rows) >= capacity {
			panic("eligibility scalar table capacity")
		}
		b := region(0x587000, off, capacity*12)
		for j, row := range rows {
			w := unsafe.Slice((*uint32)(unsafe.Pointer(&b[j*12])), 3)
			w[0], w[1], w[2] = row[0], row[1], 0x13572468
		}
	}
	if len(sp.BeastGroups) > 2 {
		panic("eligibility beast group capacity")
	}
	groups := region(0x587000, 207032, 12)
	for i, ids := range sp.BeastGroups {
		ptr := p.objectiveRegion(4 * (len(ids) + 1))
		w := unsafe.Slice((*uint32)(ptr), len(ids)+1)
		copy(w, ids)
		w[len(ids)] = 0
		*(*unsafe.Pointer)(unsafe.Pointer(&groups[i*4])) = ptr
		p.identify(ptr, 970000+uint32(i))
	}
	if len(sp.ItemTable) > 12 {
		panic("eligibility item table capacity")
	}
	itemTable := region(0x587000, 208176, 13*20)
	for i, row := range sp.ItemTable {
		off := i * 20
		*(*uint32)(unsafe.Pointer(&itemTable[off])) = uint32(i + 1)
		if row.Name != nil {
			*(*unsafe.Pointer)(unsafe.Pointer(&itemTable[off+4])) = p.objectiveString(*row.Name)
		}
		*(*uint32)(unsafe.Pointer(&itemTable[off+8])) = uint32(core.Types.IndByID(row.Type))
	}
	tables := [][]PortTestEligibilityModRow{sp.MaterialWeapon, sp.MaterialArmor, sp.Quality, sp.Effects}
	for i, off := range []uintptr{210704, 210848, 210992, 209336} {
		if len(tables[i]) > 4 {
			panic("eligibility modifier table capacity")
		}
		b := region(0x587000, off, 5*24)
		for j, row := range tables[i] {
			w := unsafe.Slice((*uint32)(unsafe.Pointer(&b[j*24])), 6)
			w[0] = uint32(j + 1)
			w[1] = uint32(uintptr(modRef(row.Ref)))
			if row.Name != nil {
				w[2] = uint32(uintptr(p.objectiveString(*row.Name)))
			}
			w[4], w[5] = row.ArmorExclude, row.WeaponExclude
		}
	}
	oldSpecial := C.dword_5d4594_1568308
	C.dword_5d4594_1568308 = 0
	st.restores = append(st.restores, func() { C.dword_5d4594_1568308 = oldSpecial })
	globals := region(0x5d4594, 1568312, 96)
	if sp.Warm {
		C.dword_5d4594_1568308 = C.uint32_t(uintptr(modRef(1)))
		for i := 0; i < 3; i++ {
			*(*uint32)(unsafe.Pointer(&globals[i*4])) = uint32(uintptr(modRef(i + 2)))
		}
		*(*uint32)(unsafe.Pointer(&globals[12])) = uint32(uintptr(modRef(1)))
		// Type caches are stored in the legacy order, including pants before shirt.
		names := []string{"Diamond", "Emerald", "Ruby", "SulphorousFlareWand", "StreetSneakers", "StreetPants", "StreetShirt"}
		names = append(names, eligibilityNames[7:20]...)
		for i, n := range names {
			*(*uint32)(unsafe.Pointer(&globals[16+i*4])) = uint32(core.Types.IndByID(n))
		}
	}
}
func (p *portTestShopPools) eligibilityRestore() {
	st := p.reportEligibility
	if st == nil {
		return
	}
	for i := len(st.restores) - 1; i >= 0; i-- {
		st.restores[i]()
	}
	p.reportEligibility = nil
}
func (p *portTestShopPools) eligibilityArg(a PortTestGameplayReportArg) uint32 {
	st := p.reportEligibility
	if st == nil || a.Ref < 0 || a.Ref > len(st.objects) || a.Offset != 0 {
		panic("eligibility object argument")
	}
	if a.Ref == 0 {
		return 0
	}
	return uint32(uintptr(st.objects[a.Ref-1].CObj()))
}
func questEligibilityInvoke(op int, args [5]uint32) uint32 {
	if op < 0 || op > 13 {
		panic("eligibility operation")
	}
	return uint32(C.questEligibilityInvoke(C.int(op), C.uint32_t(args[0])))
}
func (p *portTestShopPools) eligibilitySnapshot() *PortTestQuestEligibilityResult {
	st := p.reportEligibility
	if st == nil {
		return nil
	}
	norm := func(ptr unsafe.Pointer, n int) []uint32 {
		out := make([]uint32, n)
		for i, v := range unsafe.Slice((*uint32)(ptr), n) {
			out[i] = p.normalize(v)
		}
		return out
	}
	out := &PortTestQuestEligibilityResult{Globals: []uint32{p.normalize(uint32(C.dword_5d4594_1568308))}}
	out.Globals = append(out.Globals, norm(memmap.PtrOff(0x5d4594, 1568312), 24)...)
	for _, u := range st.objects {
		out.Objects = append(out.Objects, []uint32{uint32(u.TypeInd), uint32(u.ObjClass), uint32(u.ObjSubClass), uint32(u.ObjFlags)})
		out.Objects = append(out.Objects, norm(*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 692)), 4))
	}
	for _, m := range st.mods {
		out.Modifiers = append(out.Modifiers, norm(unsafe.Add(m.C(), 28), 3))
	}
	for _, b := range st.tables {
		out.Tables = append(out.Tables, norm(unsafe.Pointer(&b[0]), len(b)/4))
	}
	return out
}
