//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
)

func PortTestColorLight(op int, vp *noxrender.Viewport, d *client.Drawable) int {
	switch op {
	case 0:
		colorLightColor(d)
	case 1:
		colorLightIntensity(d)
	case 2:
		return colorLightPenumbra(d)
	case 3:
		colorLightTarget(d)
	case 4:
		colorLightRotate(d)
	case 5:
		return colorLightUpdate(vp, d)
	default:
		panic(op)
	}
	return 0
}
