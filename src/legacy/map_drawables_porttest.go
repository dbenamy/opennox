//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3.h"
extern uint32_t dword_5d4594_527660;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"unsafe"
)

// Calls the real installed map reader; no replacement stream or drawable factory.
func PortTestMapDrawableBase(typ int, outer, inner int16, old bool, consumed *uint32) *client.Drawable {
	r := mapDrawableReader{consumed}
	if old {
		return mapDrawableOld(typ, inner, outer, r)
	}
	return mapDrawableBase(typ, outer, r)
}
func PortTestMapDrawableRecord(op, typ int) int {
	switch op {
	case 0:
		return mapDrawableDispatch(uint16(typ))
	case 1, 2, 3, 4, 5:
		return mapDrawableTyped(op, typ)
	case 6:
		return mapDrawableSection()
	}
	panic(op)
}

func PortTestMapDrawableTeamWord() *uint32 { return (*uint32)(unsafe.Pointer(&C.dword_5d4594_527660)) }
