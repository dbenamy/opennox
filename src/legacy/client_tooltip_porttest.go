//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"unsafe"
)

// PortTestTooltip calls the production owner; the fixture owns every input.
func PortTestTooltip(dr *client.Drawable) *uint16 {
	return (*uint16)(unsafe.Pointer(nox_xxx_clientAskInfoMb_4BF050((*nox_drawable)(unsafe.Pointer(dr)))))
}
func PortTestTooltipCursor(text *uint16) {
	uiCursorTooltip(text)
}
