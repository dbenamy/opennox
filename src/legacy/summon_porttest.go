//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_1.h"
#include "client__gui__guisumn.h"
extern uint32_t dword_5d4594_1320988;
extern uint32_t dword_5d4594_1320992;
extern uint32_t dword_5d4594_1321024;
extern uint32_t dword_5d4594_1321032;
extern uint32_t dword_5d4594_1321036;
extern uint32_t dword_5d4594_1321040;
extern uint32_t dword_5d4594_1321044;
extern uint32_t dword_5d4594_1321196;
extern uint32_t dword_5d4594_1321204;
extern uint32_t dword_5d4594_1321208;
extern uint32_t nox_xxx_screenWidth_587000_184452;
*/
import "C"
import (
	"image"
	"unsafe"
)

// PortTestSummonInvoke calls the production Go owner directly.
func PortTestSummonInvoke(op string, a [5]uint32) uint32 {
	record := func(v uint32) *summonRecord { return (*summonRecord)(unsafe.Pointer(uintptr(v))) }
	pos := func(v uint32) *[2]int32 { return (*[2]int32)(unsafe.Pointer(uintptr(v))) }
	switch op {
	case "sub_4C1CA0":
		return uint32(summonSetCommand(a[0]))
	case "nox_xxx_guiDrawSummonBox_4C1FE0":
		return uint32(summonDraw(bookWindow(a[0])))
	case "nox_xxx_wndSummonGet_4C2410":
		return summonGet(pos(a[0]))
	case "nox_xxx_guiDrawSummon_4C2440":
		return summonIcon(int(a[0]))
	case "nox_xxx_guiHideSummonWindow_4C2470":
		return summonClose()
	case "sub_4C24A0":
		return 1
	case "nox_xxx_wndSummonBigButtonProc_4C24B0":
		return uint32(summonBigEvent(bookWindow(a[0]), a[1], a[2]))
	case "sub_4C2A00":
		return uint32(summonOutline(image.Pt(int(a[0]), int(a[1])), a[2], a[3], GoWStringP(unsafe.Pointer(uintptr(a[4])))))
	case "nox_client_orderCreature":
		summonOrder(record(a[0]), a[1])
		return 0
	case "nox_xxx_clientOrderCreature_4C2A60":
		return uint32(summonCommandEvent(bookWindow(a[0]), a[1], 0))
	case "nox_xxx_wndSummonProc_4C2B10":
		return uint32(summonBoxEvent(bookWindow(a[0]), a[1], a[2]))
	case "sub_4C2BD0":
		return 0
	case "sub_4C2BE0":
		return 1
	case "sub_4C2BF0":
		return summonClearGrid()
	case "sub_4C2C20":
		return uint32(sub_4C2C20((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), C.uint(a[2])))
	case "sub_4C2C60":
		return uint32(uintptr(summonSlotTooltip(bookWindow(a[0]), image.Pt(int(pos(a[1])[0]), int(pos(a[1])[1])))))
	case "sub_4C2D60":
		return summonAddress(summonFirst())
	case "sub_4C2D90":
		return summonAddress(summonNext(record(a[0])))
	case "sub_4C2DD0":
		return uint32(summonMobile(record(a[0])))
	case "sub_4C2E00":
		return uint32(summonAnyMobile())
	case "nox_xxx_cliSummonCreat_4C2E50":
		return uint32(int8(summonAdd(a[0], a[1], a[2] != 0)))
	case "sub_4C2EF0":
		return uint32(summonClass(int(a[0])))
	case "sub_4C2F20":
		return summonAddress(summonAllocate())
	case "sub_4C2F70":
		return summonLayout()
	case "sub_4C2FD0":
		return uint32(summonPlace(record(a[0])))
	case "sub_4C3030":
		return uint32(summonPaint(pos(a[0]), int32(a[1]), a[2]))
	case "sub_4C30C0":
		return uint32(summonAvailable(pos(a[0]), int32(a[1])))
	case "nox_xxx_cliSummonOnDieOrBanish_4C3140":
		summonRemove(a[0], a[1] != 0)
		return 0
	case "sub_4C31D0":
		return summonAddress(summonFind(a[0]))
	case "sub_4C3210":
		return summonDeactivate(record(a[0]))
	case "nox_xxx_sprite_4C3220":
		return uint32(bool2int(summonFind(*(*uint32)(unsafe.Pointer(uintptr(a[0]) + 128))) != nil))
	case "sub_4C3260":
		return uint32(bool2int(summonFirst() != nil))
	case "nox_xxx_guiSummonCreatureLoad_4C1D80":
		return uint32(summonCreate())
	case "nox_xxx_wndSummonCreateList_4C2560":
		summonMenu(image.Pt(int(pos(a[0])[0]), int(pos(a[0])[1])))
		return 0
	case "sub_4C27F0":
		return uint32(summonDrawMenu(bookWindow(a[0])))
	case "sub_4C2CE0":
		return uint32(summonCommandTooltip())
	}
	panic("unknown summon operation: " + op)
}

func PortTestSummonWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1320988":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320988)),
		"dword_5d4594_1320992":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320992)),
		"dword_5d4594_1321024":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321024)),
		"dword_5d4594_1321032":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321032)),
		"dword_5d4594_1321036":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321036)),
		"dword_5d4594_1321040":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321040)),
		"dword_5d4594_1321044":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321044)),
		"dword_5d4594_1321196":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321196)),
		"dword_5d4594_1321204":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321204)),
		"dword_5d4594_1321208":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321208)),
		"nox_xxx_screenWidth_587000_184452": (*uint32)(unsafe.Pointer(&C.nox_xxx_screenWidth_587000_184452)),
	}
	old := make(map[string]uint32, len(words))
	for name, p := range words {
		old[name] = *p
	}
	return words, func() {
		for name, p := range words {
			*p = old[name]
		}
	}
}
func PortTestSummonCallbacks() map[string]unsafe.Pointer {
	return map[string]unsafe.Pointer{
		"sub_4C1CA0":                            C.sub_4C1CA0,
		"nox_xxx_guiDrawSummonBox_4C1FE0":       nil,
		"nox_xxx_wndSummonGet_4C2410":           nil,
		"nox_xxx_guiDrawSummon_4C2440":          nil,
		"nox_xxx_guiHideSummonWindow_4C2470":    nil,
		"sub_4C24A0":                            nil,
		"nox_xxx_wndSummonBigButtonProc_4C24B0": nil,
		"sub_4C2A00":                            nil,
		"nox_client_orderCreature":              nil,
		"nox_xxx_clientOrderCreature_4C2A60":    nil,
		"nox_xxx_wndSummonProc_4C2B10":          nil,
		"sub_4C2BD0":                            nil,
		"sub_4C2BE0":                            nil,
		"sub_4C2BF0":                            nil,
		"sub_4C2C20":                            C.sub_4C2C20,
		"sub_4C2C60":                            nil,
		"sub_4C2D60":                            nil,
		"sub_4C2D90":                            nil,
		"sub_4C2DD0":                            nil,
		"sub_4C2E00":                            nil,
		"nox_xxx_cliSummonCreat_4C2E50":         nil, // Preserve stable capture IDs after retiring the C export.
		"sub_4C2EF0":                            nil,
		"sub_4C2F20":                            nil,
		"sub_4C2F70":                            nil,
		"sub_4C2FD0":                            nil,
		"sub_4C3030":                            nil,
		"sub_4C30C0":                            nil,
		"nox_xxx_cliSummonOnDieOrBanish_4C3140": nil, // Preserve stable capture IDs after retiring the C export.
		"sub_4C31D0":                            nil,
		"sub_4C3210":                            nil,
		"nox_xxx_sprite_4C3220":                 nil,
		"sub_4C3260":                            C.sub_4C3260,
		"nox_xxx_guiSummonCreatureLoad_4C1D80":  nil,
		"nox_xxx_wndSummonCreateList_4C2560":    nil,
		"sub_4C27F0":                            nil,
		"sub_4C2CE0":                            C.sub_4C2CE0,
	}
}
