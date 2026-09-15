//go:build porttest

package legacy

/*
#include "client__draw__animdraw.h"
#include "client__draw__canidraw.h"
#include "client__draw__staticdraw.h"
#include "client__draw__boulderdraw.h"
#include "client__draw__slavedraw.h"
extern uint32_t dword_5d4594_1321520;
extern int nox_win_width;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestSpriteAnimationCallback(op int) unsafe.Pointer {
	switch op {
	case 0:
		return C.nox_thing_animate_draw
	case 1:
		return C.nox_thing_cond_animate_draw
	case 2:
		return C.nox_thing_static_draw
	case 3:
		return C.nox_thing_static_random_draw
	case 4:
		return C.nox_thing_boulder_draw
	case 5:
		return C.nox_thing_slave_draw
	case 6:
		return C.nox_thing_animate_state_draw
	}
	panic("unknown sprite animation callback")
}

// Own only globals used by the non-player, fixed-lighting drawObject path.
func PortTestSpriteAnimationEnvironment() func() {
	ghost, width := C.dword_5d4594_1321520, C.nox_win_width
	C.dword_5d4594_1321520, C.nox_win_width = 0x7fffffff, 96
	var restores []func()
	for _, region := range [][3]uintptr{{0x587000, 80808, 4}, {0x587000, 185472, 12}, {0x5D4594, 1321512, 8}, {0x973F18, 76, 4}, {0x973F18, 88, 4}} {
		buf := unsafe.Slice((*byte)(memmap.PtrOff(region[0], region[1])), int(region[2]))
		old := append([]byte(nil), buf...)
		clear(buf)
		restores = append(restores, func() { copy(buf, old) })
	}
	for _, off := range []uintptr{185472, 185476, 185480} {
		*memmap.PtrUint32(0x587000, off) = 255
	}
	return func() {
		for i := len(restores) - 1; i >= 0; i-- {
			restores[i]()
		}
		C.dword_5d4594_1321520, C.nox_win_width = ghost, width
	}
}
