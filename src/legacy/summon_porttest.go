//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_1.h"
#include "client__gui__guisumn.h"
void nox_client_orderCreature(int creature, int command);
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
import "unsafe"

// PortTestSummonInvoke calls production C directly; no test algorithm is retained.
func PortTestSummonInvoke(op string, a [5]uint32) uint32 {
	switch op {
	case "sub_4C1CA0":
		return uint32(C.sub_4C1CA0(C.int(a[0])))
	case "nox_xxx_guiDrawSummonBox_4C1FE0":
		return uint32(C.nox_xxx_guiDrawSummonBox_4C1FE0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_wndSummonGet_4C2410":
		return uint32(C.nox_xxx_wndSummonGet_4C2410((*C.int2)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_guiDrawSummon_4C2440":
		return uint32(C.nox_xxx_guiDrawSummon_4C2440(C.int(a[0])))
	case "nox_xxx_guiHideSummonWindow_4C2470":
		return uint32(C.nox_xxx_guiHideSummonWindow_4C2470())
	case "sub_4C24A0":
		return uint32(C.sub_4C24A0())
	case "nox_xxx_wndSummonBigButtonProc_4C24B0":
		return uint32(C.nox_xxx_wndSummonBigButtonProc_4C24B0(C.int(a[0]), C.int(a[1]), C.uint(a[2])))
	case "sub_4C2A00":
		return uint32(C.sub_4C2A00(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), (*C.short)(unsafe.Pointer(uintptr(a[4])))))
	case "nox_client_orderCreature":
		C.nox_client_orderCreature(C.int(a[0]), C.int(a[1]))
		return 0
	case "nox_xxx_clientOrderCreature_4C2A60":
		return uint32(C.nox_xxx_clientOrderCreature_4C2A60(C.int(a[0]), C.uint(a[1])))
	case "nox_xxx_wndSummonProc_4C2B10":
		return uint32(C.nox_xxx_wndSummonProc_4C2B10((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), C.uint(a[1]), C.uint(a[2])))
	case "sub_4C2BD0":
		return uint32(C.sub_4C2BD0())
	case "sub_4C2BE0":
		return uint32(C.sub_4C2BE0())
	case "sub_4C2BF0":
		return uint32(uintptr(unsafe.Pointer(C.sub_4C2BF0())))
	case "sub_4C2C20":
		return uint32(C.sub_4C2C20((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), C.uint(a[2])))
	case "sub_4C2C60":
		return uint32(C.sub_4C2C60((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), (*C.int2)(unsafe.Pointer(uintptr(a[1])))))
	case "sub_4C2D60":
		return uint32(uintptr(unsafe.Pointer(C.sub_4C2D60())))
	case "sub_4C2D90":
		return uint32(uintptr(unsafe.Pointer(C.sub_4C2D90(C.int(a[0])))))
	case "sub_4C2DD0":
		return uint32(C.sub_4C2DD0(C.int(a[0])))
	case "sub_4C2E00":
		return uint32(C.sub_4C2E00())
	case "nox_xxx_cliSummonCreat_4C2E50":
		return uint32(C.nox_xxx_cliSummonCreat_4C2E50(C.int(a[0]), C.int(a[1]), C.int(a[2])))
	case "sub_4C2EF0":
		return uint32(C.sub_4C2EF0(C.int(a[0])))
	case "sub_4C2F20":
		return uint32(uintptr(unsafe.Pointer(C.sub_4C2F20())))
	case "sub_4C2F70":
		return uint32(uintptr(unsafe.Pointer(C.sub_4C2F70())))
	case "sub_4C2FD0":
		return uint32(C.sub_4C2FD0(C.int(a[0])))
	case "sub_4C3030":
		return uint32(C.sub_4C3030((*C.int)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), C.int(a[2])))
	case "sub_4C30C0":
		return uint32(C.sub_4C30C0((*C.int)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1])))
	case "nox_xxx_cliSummonOnDieOrBanish_4C3140":
		C.nox_xxx_cliSummonOnDieOrBanish_4C3140(C.int(a[0]), (unsafe.Pointer)(unsafe.Pointer(uintptr(a[1]))))
		return 0
	case "sub_4C31D0":
		return uint32(uintptr(unsafe.Pointer(C.sub_4C31D0(C.int(a[0])))))
	case "sub_4C3210":
		return uint32(C.sub_4C3210(C.int(a[0])))
	case "nox_xxx_sprite_4C3220":
		return uint32(C.nox_xxx_sprite_4C3220((*C.nox_drawable)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_4C3260":
		return uint32(C.sub_4C3260())
	case "nox_xxx_guiSummonCreatureLoad_4C1D80":
		return uint32(C.nox_xxx_guiSummonCreatureLoad_4C1D80())
	case "nox_xxx_wndSummonCreateList_4C2560":
		C.nox_xxx_wndSummonCreateList_4C2560((*C.int2)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_4C27F0":
		return uint32(C.sub_4C27F0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_4C2CE0":
		return uint32(C.sub_4C2CE0())
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
		"nox_xxx_guiDrawSummonBox_4C1FE0":       C.nox_xxx_guiDrawSummonBox_4C1FE0,
		"nox_xxx_wndSummonGet_4C2410":           C.nox_xxx_wndSummonGet_4C2410,
		"nox_xxx_guiDrawSummon_4C2440":          C.nox_xxx_guiDrawSummon_4C2440,
		"nox_xxx_guiHideSummonWindow_4C2470":    C.nox_xxx_guiHideSummonWindow_4C2470,
		"sub_4C24A0":                            C.sub_4C24A0,
		"nox_xxx_wndSummonBigButtonProc_4C24B0": C.nox_xxx_wndSummonBigButtonProc_4C24B0,
		"sub_4C2A00":                            C.sub_4C2A00,
		"nox_client_orderCreature":              C.nox_client_orderCreature,
		"nox_xxx_clientOrderCreature_4C2A60":    C.nox_xxx_clientOrderCreature_4C2A60,
		"nox_xxx_wndSummonProc_4C2B10":          C.nox_xxx_wndSummonProc_4C2B10,
		"sub_4C2BD0":                            C.sub_4C2BD0,
		"sub_4C2BE0":                            C.sub_4C2BE0,
		"sub_4C2BF0":                            C.sub_4C2BF0,
		"sub_4C2C20":                            C.sub_4C2C20,
		"sub_4C2C60":                            C.sub_4C2C60,
		"sub_4C2D60":                            C.sub_4C2D60,
		"sub_4C2D90":                            C.sub_4C2D90,
		"sub_4C2DD0":                            C.sub_4C2DD0,
		"sub_4C2E00":                            C.sub_4C2E00,
		"nox_xxx_cliSummonCreat_4C2E50":         C.nox_xxx_cliSummonCreat_4C2E50,
		"sub_4C2EF0":                            C.sub_4C2EF0,
		"sub_4C2F20":                            C.sub_4C2F20,
		"sub_4C2F70":                            C.sub_4C2F70,
		"sub_4C2FD0":                            C.sub_4C2FD0,
		"sub_4C3030":                            C.sub_4C3030,
		"sub_4C30C0":                            C.sub_4C30C0,
		"nox_xxx_cliSummonOnDieOrBanish_4C3140": C.nox_xxx_cliSummonOnDieOrBanish_4C3140,
		"sub_4C31D0":                            C.sub_4C31D0,
		"sub_4C3210":                            C.sub_4C3210,
		"nox_xxx_sprite_4C3220":                 C.nox_xxx_sprite_4C3220,
		"sub_4C3260":                            C.sub_4C3260,
		"nox_xxx_guiSummonCreatureLoad_4C1D80":  C.nox_xxx_guiSummonCreatureLoad_4C1D80,
		"nox_xxx_wndSummonCreateList_4C2560":    C.nox_xxx_wndSummonCreateList_4C2560,
		"sub_4C27F0":                            C.sub_4C27F0,
		"sub_4C2CE0":                            C.sub_4C2CE0,
	}
}
