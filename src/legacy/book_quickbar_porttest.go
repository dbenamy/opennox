//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "client__gui__guispell.h"
extern void* nox_xxx_aClosewoodengat_587000_133480;
*/
import "C"
import "unsafe"

// PortTestBookQuickbar binds the actual quickbar record used by the unchanged
// slot owner. Its layout and storage come from the normal startup blob.
func PortTestBookQuickbar(p unsafe.Pointer) func() {
	old, oldGuard := C.nox_xxx_aClosewoodengat_587000_133480, dword_5d4594_1049496
	C.nox_xxx_aClosewoodengat_587000_133480 = p
	dword_5d4594_1049496 = 0
	return func() { C.nox_xxx_aClosewoodengat_587000_133480 = old; dword_5d4594_1049496 = oldGuard }
}

func PortTestBookQuickbarInit(p unsafe.Pointer, x, y int) uint32 {
	return uint32(quickbarInitWindow((*quickbarRecord)(p), x, y, 5, 0, bookEvent(quickbarSlotEvent), quickbarDrawAbility))
}
func PortTestBookQuickbarCallbacks() map[unsafe.Pointer]uint32 {
	return nil
}

// PortTestBookPauseOwner owns the already-paused reward presentation state.
// Live pause/unpause routines remain unchanged and are exercised on completion.
func PortTestBookPauseOwner() (func(), func()) {
	a, b, c := dword_5d4594_2523804, dword_5d4594_2523780, dword_5d4594_2523776
	return func() { dword_5d4594_2523804 = 1; dword_5d4594_2523780 = 0; dword_5d4594_2523776 = 0 }, func() { dword_5d4594_2523804 = a; dword_5d4594_2523780 = b; dword_5d4594_2523776 = c }
}

func PortTestBookTrapWords() ([]*uint32, func()) {
	words := []*uint32{(*uint32)(unsafe.Pointer(&dword_5d4594_1049500)), (*uint32)(unsafe.Pointer(&dword_5d4594_1049504)), (*uint32)(unsafe.Pointer(&dword_5d4594_1049520)), (*uint32)(unsafe.Pointer(&dword_5d4594_1049536))}
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
