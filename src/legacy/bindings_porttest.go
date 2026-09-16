//go:build porttest

package legacy

/*
#include "defs.h"
extern int nox_win_width;
extern int nox_win_height;
extern nox_gui_animation* nox_wnd_xxx_1522608;
#include "GAME3_1.h"
#include "client__gui__guiinput.h"
#include "client__shell__inputcfg__inputcfg.h"
extern uint32_t dword_5d4594_1321224;
extern uint32_t dword_5d4594_1321228;
extern uint32_t dword_5d4594_1321232;
extern nox_window* dword_5d4594_1321236;
extern nox_window* dword_5d4594_1321240;
extern nox_window* dword_5d4594_1321244;
extern nox_window* dword_5d4594_1321248;
extern uint32_t dword_5d4594_1321252;
extern uint32_t dword_5d4594_1522604;
extern uint32_t dword_5d4594_1522612;
extern nox_window* dword_5d4594_1522616;
extern nox_window* dword_5d4594_1522620;
extern nox_window* dword_5d4594_1522624;
extern nox_window* dword_5d4594_1522628;
extern uint32_t dword_5d4594_1522632;
*/
import "C"
import "unsafe"

// PortTestBindingWords owns live globals independently of the backing memory blob.
func PortTestBindingWords() (map[int]*uint32, func()) {
	w := map[int]*uint32{
		1321224: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321224)),
		1321228: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321228)),
		1321232: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321232)),
		1321236: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321236)),
		1321240: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321240)),
		1321244: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321244)),
		1321248: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321248)),
		1321252: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321252)),
		1522604: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522604)),
		1522612: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522612)),
		1522616: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522616)),
		1522620: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522620)),
		1522624: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522624)),
		1522628: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522628)),
		1522632: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1522632)),
	}
	old := make(map[int]uint32)
	for n, p := range w {
		old[n] = *p
		*p = 0
	}
	return w, func() {
		for n, p := range w {
			*p = old[n]
		}
	}
}

// PortTestBindingAssign invokes the existing C keyboard/mouse assignment owners.
func PortTestBindingAssign(menu, mouse bool, key uint32) int {
	if menu {
		if mouse {
			return int(C.sub_4CC3C0(C.uint(key)))
		}
		return int(C.sub_4CC280(C.uint(key)))
	}
	if mouse {
		return int(C.sub_4C4100(C.uint(key)))
	}
	return int(C.sub_4C3FC0(C.uint(key)))
}
func PortTestBindingModal(menu bool, win uint32, event int, key uint32, state int) int {
	if menu {
		return int(C.sub_4CC170(C.int(win), C.int(event), (*C.char)(unsafe.Pointer(uintptr(key))), C.int(state)))
	}
	return int(C.sub_4C3EB0(C.int(win), C.int(event), C.uint(key), C.int(state)))
}

func PortTestBindingApply(menu bool) {
	if menu {
		C.sub_4CBD30()
	} else {
		C.sub_4C3620()
	}
}

func PortTestBindingDimensions() [2]*int32 {
	return [2]*int32{(*int32)(unsafe.Pointer(&C.nox_win_width)), (*int32)(unsafe.Pointer(&C.nox_win_height))}
}
func PortTestBindingConstruct(menu bool) int {
	if menu {
		return int(C.sub_4CB880())
	}
	return int(C.sub_4C3760())
}
func PortTestBindingDestroy() int { return int(C.sub_4C4220()) }
func PortTestBindingInvoke(op string, a [4]uint32) uint32 {
	switch op {
	case "sub_4C3500":
		return uint32(C.sub_4C3500())
	case "sub_4C35B0":
		return uint32(C.sub_4C35B0(C.int(a[0])))
	case "sub_4C4260":
		C.sub_4C4260()
		return 0
	case "sub_4C4280":
		return uint32(C.sub_4C4280())
	case "sub_4CBB70":
		return uint32(C.sub_4CBB70())
	case "sub_4CBBB0":
		return uint32(C.sub_4CBBB0())
	case "sub_4C3A60":
		return uint32(C.sub_4C3A60((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), C.uint(a[1]), C.uint(a[2]), C.int(a[3])))
	case "sub_4CC140":
		return uint32(C.sub_4CC140((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), C.uint(a[1]), C.uint(a[2]), C.int(a[3])))
	case "sub_4C3CD0":
		return uint32(C.sub_4C3CD0(C.int(a[0]), C.uint(a[1]), C.int(a[2]), C.int(a[3])))
	case "sub_4CBF60":
		return uint32(C.sub_4CBF60(C.int(a[0]), C.uint(a[1]), C.int(a[2]), C.int(a[3])))
	default:
		panic(op)
	}
}
func PortTestBindingAnimationWord() *uint32 { return (*uint32)(unsafe.Pointer(&C.nox_wnd_xxx_1522608)) }
