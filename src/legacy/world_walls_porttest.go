//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/server"
	"image"
)

// Invoke the real Go owners against unchanged frozen C expectations.
func PortTestWorldWalls(op int, vp *noxrender.Viewport, dr *client.Drawable, wl *server.Wall, p image.Point) (int, image.Point) {
	var out image.Point
	switch op {
	case 0:
		return p.Y, vp.ToScreenPos(p)
	case 1:
		return p.Y, vp.ToWorldPos(p)
	case 2:
		worldWallDraw(vp, wl)
		return 0, out
	case 3:
		return bool2int(worldWallPlayerVisible(dr)), out
	case 4:
		return bool2int(worldWallStaticPass(dr)), out
	case 5:
		return bool2int(worldWallDynamicPass(dr)), out
	case 6:
		return bool2int(worldWallSpecialPass(dr)), out
	case 7:
		return bool2int(worldWallInactivePass(dr)), out
	case 8:
		return bool2int(worldWallImageInterval(p.X, p.Y)), out
	}
	panic("unknown world-wall operation")
}
