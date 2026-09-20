//go:build porttest

package legacy

/*
#include "GAME2_3.h"
#include "GAME3.h"
extern uint32_t dword_5d4594_1200776;
extern uint32_t dword_5d4594_1200796;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestDrawableEffect(op int, dr *client.Drawable, a [5]uint32) uintptr {
	switch op {
	case 0:
		xy, free := alloc.Make([]uint16{uint16(a[0]), uint16(a[1])}, 2)
		defer free()
		return uintptr(unsafe.Pointer(C.nox_xxx_netHandleSummonPacket_4B7C40(C.short(a[2]), (*C.ushort)(unsafe.Pointer(&xy[0])), C.ushort(a[3]), C.uchar(a[4]), C.short(a[2]>>16))))
	case 1:
		C.sub_4B7EE0(C.short(a[0]))
		return 0
	case 2:
		return uintptr(drawableShieldLoad())
	case 3:
		return uintptr(unsafe.Pointer(C.nox_xxx_fxShield_4B8090(C.uint(a[0]), C.int(a[1]))))
	case 4:
		key, free := alloc.New(uint32(0))
		*key = a[0]
		defer free()
		drawableShieldScan(dr, *key)
		return 0
	case 5:
		return uintptr(drawableEffectTypes())
	}
	panic("unknown drawable effect")
}

// Named globals and mapped historical addresses are distinct owners.
func PortTestDrawableEffectGlobals() ([]*uint32, func()) {
	words := []*uint32{&drawableSummonSpark, &drawableShieldFound, (*uint32)(unsafe.Pointer(&C.dword_5d4594_1200776)), (*uint32)(unsafe.Pointer(&C.dword_5d4594_1200796))}
	for off := uintptr(1313744); off <= 1313784; off += 4 {
		words = append(words, memmap.PtrUint32(0x5D4594, off))
	}
	for off := uintptr(1200772); off <= 1200828; off += 4 {
		words = append(words, memmap.PtrUint32(0x5D4594, off))
	}
	saved := make([]uint32, len(words))
	for i, p := range words {
		saved[i] = *p
		*p = 0
	}
	return words, func() {
		for i, p := range words {
			*p = saved[i]
		}
	}
}
