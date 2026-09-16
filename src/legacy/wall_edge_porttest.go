//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
)

// Compare the real Go owner against unchanged original-C expectations.
func PortTestWallEdge(img noxrender.ImageHandle, pos image.Point, first, second *[3]uint32, bottom, left, right, crop, flags int) {
	wallEdgeDraw(img, pos, first, second, bottom, left, right, crop)
}
