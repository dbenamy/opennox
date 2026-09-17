//go:build porttest

package legacy

/*
#include "GAME2.h"
#include "GAME3_1.h"
#include "client__gui__servopts__objlst.h"
#include "client__gui__servopts__spelllst.h"
#include "client__gui__servopts__general.h"
#include "client__gui__servopts__access.h"
extern uint32_t dword_5d4594_1045464;
extern uint32_t dword_5d4594_1045468;
extern uint32_t dword_5d4594_1045460;
*/
import "C"
import "unsafe"
import "github.com/opennox/opennox/v1/common/memmap"
import "github.com/opennox/opennox/v1/legacy/common/alloc"
import "github.com/opennox/opennox/v1/client/gui"

// Exercise original production helpers with caller-owned guarded storage.
func PortTestServerPanelsByte(p *byte, mask byte, on int) bool {
	return unsafe.Pointer(C.sub_453620((*C.uint8_t)(unsafe.Pointer(p)), C.char(mask), C.int(on))) == unsafe.Pointer(p)
}
func PortTestServerPanelsWord(p *uint32, mask uint32, on int) bool {
	return unsafe.Pointer(C.sub_453640((*C.uint32_t)(unsafe.Pointer(p)), C.int(mask), C.int(on))) == unsafe.Pointer(p)
}
func PortTestServerPanelsBit(p *uint32, index, on int) (uint32, bool) {
	r := uint32(C.sub_453FA0(C.int(uintptr(unsafe.Pointer(p))), C.int(index), C.int(on)))
	return r, C.sub_454000(C.int(uintptr(unsafe.Pointer(p))), C.int(index)) != 0
}
func PortTestServerPanelsMaskQuery(mask uint32, armor bool) bool {
	if armor {
		return C.sub_453690(C.int(mask)) != 0
	}
	return C.sub_453660(C.int(mask)) != 0
}
func PortTestServerPanelsClass(index byte, armor bool) bool {
	old := C.dword_5d4594_1045460
	C.dword_5d4594_1045460 = C.uint32_t(bool2int(armor))
	defer func() { C.dword_5d4594_1045460 = old }()
	return C.sub_453080(C.char(index)) != 0
}

func PortTestServerPanelsWeaponStore(p *uint32) bool {
	return unsafe.Pointer(C.sub_4535E0((*C.int)(unsafe.Pointer(p)))) == unsafe.Pointer(p)
}
func PortTestServerPanelsWeaponPointer() *uint32 {
	return (*uint32)(unsafe.Pointer(C.sub_453600()))
}
func PortTestServerPanelsArmorStore(v uint32) uint32 { return uint32(C.sub_4535F0(C.int(v))) }
func PortTestServerPanelsArmorLoad() uint32          { return uint32(C.sub_453610()) }
func PortTestServerPanelsSpellStore(p *uint32)       { C.sub_453F70(unsafe.Pointer(p)) }
func PortTestServerPanelsSpellPointer() *uint32 {
	return (*uint32)(unsafe.Pointer(C.sub_453F90()))
}

// The snapshot helper's last return is either a scalar or a pointer into its
// output. Identify owned pointers before recording; never capture an address.
func PortTestServerPanelsSpellSnapshot(p *uint32) (uint32, int) {
	r := uint32(C.sub_454040((*C.uint32_t)(unsafe.Pointer(p))))
	for i := 0; i < 5; i++ {
		if uintptr(r) == uintptr(unsafe.Pointer(p))+uintptr(4*i) {
			return 0, i
		}
	}
	return r, -1
}
func PortTestServerPanelsSpellApply(p *uint32) int {
	return int(C.sub_4540E0(C.int(uintptr(unsafe.Pointer(p)))))
}

func PortTestServerPanelsObjectWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"kind": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045460)),
		"list": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045464)),
		"root": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045468)),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestServerPanelsConstruct(kind string, parent *gui.Window, settings unsafe.Pointer) int {
	p := C.int(uintptr(unsafe.Pointer(parent)))
	switch kind {
	case "weapon":
		return int(C.nox_xxx_guiObjlistLoad_4530C0(p, 0x1000000))
	case "armor":
		return int(C.nox_xxx_guiObjlistLoad_4530C0(p, 0x2000000))
	case "spell":
		return int(C.nox_xxx_guiSpelllistLoad_453850(p))
	case "access":
		return int(C.nox_xxx_guiServerAccessLoad_4541D0(p))
	case "general":
		return int(C.nox_xxx_gui_4AD320(p))
	case "advanced":
		return int(C.nox_xxx_loadAdvancedWnd_4BDC10((*C.int)(settings)))
	case "advserv":
		return int(C.sub_4BDFD0())
	default:
		panic(kind)
	}
}
func PortTestServerPanelsRefresh(kind string) {
	switch kind {
	case "object":
		C.sub_453750()
	case "spell":
		C.sub_454120()
	case "access":
		C.sub_454740()
	case "general":
		C.sub_4AD840()
	default:
		panic(kind)
	}
}
func PortTestServerPanelsWeaponSnapshot(p *uint32) int {
	return int(C.sub_4536B0((*C.uint32_t)(unsafe.Pointer(p))))
}
func PortTestServerPanelsArmorSnapshot() int { return int(C.sub_453710()) }

func PortTestServerPanelsAccessLookup(name string) int {
	return int(C.sub_4559B0((*C.wchar2_t)(unsafe.Pointer(alloc.InternCString16(name)))))
}
func PortTestServerPanelsAccessSelected() bool                   { return C.sub_455770() != 0 }
func PortTestServerPanelsAccessRefresh()                         { C.sub_455800() }
func PortTestServerPanelsAccessClose(destroy bool)               { C.sub_4557D0(C.int(bool2int(destroy))) }
func PortTestServerPanelsAdvancedUpdate(settings unsafe.Pointer) { C.sub_4BDF70((*C.int)(settings)) }
func PortTestServerPanelsCallbackTable() func() {
	table := unsafe.Slice((*unsafe.Pointer)(memmap.PtrOff(0x587000, 180016)), 4)
	old := append([]unsafe.Pointer(nil), table...)
	table[0] = nil
	table[1] = unsafe.Pointer(C.sub_454120)
	table[2] = unsafe.Pointer(C.sub_453750)
	table[3] = unsafe.Pointer(C.sub_453750)
	return func() { copy(table, old) }
}
