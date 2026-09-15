//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestMinimapObjectiveVisibility(t *testing.T) {
	o := newMinimapOwner(t)
	var rows []minimapResult
	for _, zoom := range []uint32{500, 1750, 2300} {
		for _, level := range []int{1, 2} {
			for scene := 0; scene < 12; scene++ {
				o.resetMinimap(t)
				*o.words["zoom"] = zoom
				local := o.player(t, 0, image.Pt(230, 230))
				remote := o.player(t, 1, image.Pt(245, 235))
				local.TeamVal.ID = 1
				remote.TeamVal.ID = 2
				dr := o.drawable(20, image.Pt(220, 225))
				o.c.Objs.MinimapAdd(dr, 1)
				switch scene {
				case 0: // Ordinary tracked object.
				case 1:
					dr.ObjClass = 0x400000
					dr.ObjSubClass = 8
				case 2:
					dr.TypeIDVal = uint32(o.c.Things.IndByID("Crown"))
				case 3:
					dr.TypeIDVal = uint32(o.c.Things.IndByID("Crown"))
					dr.TeamVal.ID = 2
				case 4:
					dr.TypeIDVal = uint32(o.c.Things.IndByID("GameBall"))
				case 5:
					dr.TypeIDVal = uint32(o.c.Things.IndByID("GameBall"))
					dr.TeamVal.ID = 1
				case 6:
					o.c.Objs.MinimapAdd(local, 1)
				case 7:
					o.c.Objs.MinimapAdd(remote, 1)
				case 8:
					noxflags.SetGame(32)
					o.c.Objs.MinimapAdd(remote, 1)
					o.players[1].WeaponEquip = 1
				case 9:
					remote.Buffs = 1 << 30
				case 10:
					remote.TeamVal.ID = 1
				case 11:
					o.players[0].Field3680 = 1
				}
				ret := legacy.PortTestMinimap(13, int(uintptr(unsafe.Pointer(local))), level, 0, 0, 0)
				rows = append(rows, o.capture(t, 13, ret))
			}
		}
	}
	minimapCapture(t, "objective-visibility", rows, "0c41cf00a5a44d1291ae59022238f47ef539dfa964bc961942e3ed56b6be0644")
}
