package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

//export sub_478030
func sub_478030() C.int { return C.int(uiShopActive()) }

//export sub_478040
func sub_478040() C.int { return C.int(uiShopCancelRequest()) }

//export sub_478080
func sub_478080(code C.int) C.int { return C.int(uiInventoryPointer(uiShopDrawable(uint32(code)).C())) }

//export sub_478850
func sub_478850(point C.int, code C.short, typ, count C.int) {
	uiShopBuyAccept(uint16(code), uint32(typ), uint32(count))
}

//export sub_478E50
func sub_478E50(w, event C.int, packed C.uint) C.int {
	return C.int(uiShopHover(uiInventoryPackedPoint(uintptr(packed))))
}

//export sub_479280
func sub_479280() { uiShopClose() }

//export sub_479300
func sub_479300(typ, code, value C.int, health C.short, mods C.int) C.uint32_t {
	return C.uint32_t(uiShopAdd(uint32(typ), uint32(code), uint32(value), uint16(health), unsafe.Pointer(uintptr(uint32(mods)))))
}

//export sub_479480
func sub_479480(code C.int) C.uint32_t { return C.uint32_t(uiShopRemove(uint32(code))) }

//export sub_479590
func sub_479590() C.int { return C.int(uiShopMode()) }

//export sub_4795A0
func sub_4795A0(mode C.int) { uiShopSetMode(uint32(mode)) }

//export sub_479690
func sub_479690(point C.int, code, typ C.short, count C.int) C.int {
	return C.int(uiShopSellAccept(uint16(code), uint16(typ), uint32(count)))
}

//export sub_479680
func sub_479680() { *uiShopWord(1098616) = 0 }

//export sub_479810
func sub_479810() { *uiShopWord(1098620) = 0 }

//export sub_479820
func sub_479820(point C.int, code C.short) C.int { return C.int(uiShopRepairAccept(uint16(code))) }

//export nox_xxx_cliStartShopDlg_478FD0
func nox_xxx_cliStartShopDlg_478FD0(name *C.wchar2_t, greeting *C.char, typ C.int) C.int {
	return C.int(uiShopStart((*uint16)(unsafe.Pointer(name)), alloc.GoString((*byte)(unsafe.Pointer(greeting))), uint32(typ)))
}

//export sub_479520
func sub_479520(amount C.int) { uiShopGoldWarning(uint32(amount)) }

//export sub_4795E0
func sub_4795E0(code, value C.int) C.int { return C.int(uiShopSellShow(uint32(code), uint32(value))) }

//export sub_479740
func sub_479740(code C.int, value C.uint) { uiShopRepairShow(uint32(code), uint32(value)) }
