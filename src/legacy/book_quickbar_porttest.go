//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "client__gui__guispell.h"
extern void* nox_xxx_aClosewoodengat_587000_133480;
extern uint32_t dword_5d4594_1049496;
extern uint32_t dword_5d4594_1049500, dword_5d4594_1049504, dword_5d4594_1049520, dword_5d4594_1049536;
extern uint32_t dword_5d4594_2523804, dword_5d4594_2523780, dword_5d4594_2523776;
*/
import "C"
import "unsafe"

// PortTestBookQuickbar binds the actual quickbar record used by the unchanged
// slot owner. Its layout and storage come from the normal startup blob.
func PortTestBookQuickbar(p unsafe.Pointer) func() {
	old, oldGuard := C.nox_xxx_aClosewoodengat_587000_133480, C.dword_5d4594_1049496
	C.nox_xxx_aClosewoodengat_587000_133480 = p
	C.dword_5d4594_1049496 = 0
	return func() { C.nox_xxx_aClosewoodengat_587000_133480 = old; C.dword_5d4594_1049496 = oldGuard }
}

func PortTestBookQuickbarInit(p unsafe.Pointer, x, y int) uint32 {
	return uint32(C.nox_xxx_quickBarInitWindow_4601F0(C.int(uintptr(p)), C.int(x), C.int(y), 5, 0, C.int(uintptr(unsafe.Pointer(C.nox_xxx_quickBarWnd_45EF50))), C.int(uintptr(unsafe.Pointer(C.nox_xxx_quickBarDrawFn_45FBD0)))))
}
func PortTestBookQuickbarCallbacks() map[unsafe.Pointer]uint32 {
	return map[unsafe.Pointer]uint32{
		unsafe.Pointer(C.nox_xxx_quickbar_45F8D0):       0xee700001,
		unsafe.Pointer(C.nox_xxx_quickBarWnd_45EF50):    0xee700002,
		unsafe.Pointer(C.nox_xxx_quickBarDrawFn_45FBD0): 0xee700003,
	}
}

// PortTestBookPauseOwner owns the already-paused reward presentation state.
// Live pause/unpause routines remain unchanged and are exercised on completion.
func PortTestBookPauseOwner() (func(), func()) {
	a, b, c := C.dword_5d4594_2523804, C.dword_5d4594_2523780, C.dword_5d4594_2523776
	return func() { C.dword_5d4594_2523804 = 1; C.dword_5d4594_2523780 = 0; C.dword_5d4594_2523776 = 0 }, func() { C.dword_5d4594_2523804 = a; C.dword_5d4594_2523780 = b; C.dword_5d4594_2523776 = c }
}

func PortTestBookTrapWords() ([]*uint32, func()) {
	words := []*uint32{(*uint32)(unsafe.Pointer(&C.dword_5d4594_1049500)), (*uint32)(unsafe.Pointer(&C.dword_5d4594_1049504)), (*uint32)(unsafe.Pointer(&C.dword_5d4594_1049520)), (*uint32)(unsafe.Pointer(&C.dword_5d4594_1049536))}
	old := make([]uint32, len(words))
	for i, p := range words {
		old[i] = *p
	}
	return words, func() {
		for i, p := range words {
			*p = old[i]
		}
	}
}
