//go:build porttest

package legacy

/*
#include "GAME3.h"
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
)

// PortTestEffectsOrbEnvironment restores the actual production cache and color
// words touched by the stationary white-orb probe, including cold type lookup.
func PortTestEffectsOrbEnvironment() func() {
	offsets := []uintptr{1313660, 1313664, 1313668, 1313672, 1313676, 1313680, 1313684, 1313588, 1313592}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		p := memmap.PtrUint32(0x5D4594, off)
		old[i] = *p
		*p = 0
	}
	*memmap.PtrUint32(0x5D4594, 1313588) = 0xffff
	*memmap.PtrUint32(0x5D4594, 1313592) = 0xffff
	return func() {
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
	}
}

func PortTestEffectsOrb(vp *noxrender.Viewport, dr *client.Drawable, move bool) int {
	return int(C.sub_4B6B80((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr)), C.int(bool2int(move))))
}
