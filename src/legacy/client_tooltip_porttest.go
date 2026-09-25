//go:build porttest

package legacy

/*
#include "client__gui__tooltip.h"
#include "GAME2_2.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"unsafe"
)

// PortTestTooltip calls the production C ABI; the fixture owns every input.
func PortTestTooltip(dr *client.Drawable) *uint16 {
	return (*uint16)(unsafe.Pointer(C.nox_xxx_clientAskInfoMb_4BF050((*C.nox_drawable)(unsafe.Pointer(dr)))))
}
func PortTestTooltipCursor(text *uint16) {
	uiCursorTooltip(text)
}
