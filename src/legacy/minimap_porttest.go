//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"image"
	"unsafe"
)

// PortTestMinimap invokes the production Go owner against frozen C expectations.
func PortTestMinimap(op int, a, b, c, d, e int) uint32 {
	p := image.Pt(a, b)
	switch op {
	case 0:
		minimapZoomIn()
		return minimapZoom
	case 1:
		minimapZoomOut()
		return minimapZoom
	case 2:
		return uint32(minimapSetZoom(a))
	case 3:
		return uint32(minimapWall(p, byte(c), d))
	case 4:
		return uint32(minimapColoredLine(a, b, c, d, e))
	case 5:
		return uint32(minimapFlag(p))
	case 6:
		return uint32(minimapCrown(p))
	case 7:
		minimapCircle(p)
		return 0
	case 8:
		return uint32(minimapOutline(a, b, c, d))
	case 9:
		minimapPoint(a, b, false)
		return 0
	case 10:
		minimapPoint(a, b, true)
		return 0
	case 11:
		return uint32(minimapLevel((*client.Drawable)(unsafe.Pointer(uintptr(a)))))
	case 12:
		minimapDrawSprite((*client.Drawable)(unsafe.Pointer(uintptr(a))))
		return 0
	case 13:
		return uint32(minimapDraw((*client.Drawable)(unsafe.Pointer(uintptr(a))), b))
	case 14:
		return uint32(minimapDrawAndMessages())
	}
	panic("unknown minimap operation")
}
func PortTestMinimapZoom() (*uint32, func()) {
	old := minimapZoom
	return &minimapZoom, func() { minimapZoom = old }
}
