//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestWorldWallsPlayerVisibility(t *testing.T) {
	o := newMinimapOwner(t)
	type result struct {
		Role, Team    int
		Local, Vision bool
		Return        int
		Render        objectRenderResult
	}
	var rows []result
	for role := 0; role < 3; role++ {
		for team := 0; team < 3; team++ {
			for _, localPresent := range []bool{false, true} {
				for _, vision := range []bool{false, true} {
					o.resetMinimap(t)
					local := o.player(t, 0, image.Pt(230, 230))
					target := local
					if role == 1 {
						target = o.player(t, 1, image.Pt(245, 235))
					}
					if role == 2 {
						target = o.drawable(20, image.Pt(245, 235))
					}
					if team > 0 {
						local.TeamVal.ID = 1
						target.TeamVal.ID = 1
						if team == 2 && role != 0 {
							target.TeamVal.ID = 2
						}
					}
					if vision {
						local.Buffs = 1 << 21
					}
					if !localPresent {
						*memmap.PtrPtr(0x852978, 8) = nil
					}
					ret, _ := legacy.PortTestWorldWalls(3, nil, target, nil, image.Point{})
					want := role == 0 || team == 1 || (localPresent && vision)
					if (ret != 0) != want {
						t.Fatalf("role%d team%d local%v vision%v -> %d", role, team, localPresent, vision, ret)
					}
					rows = append(rows, result{role, team, localPresent, vision, ret, o.renderResult(t, role, 0, ret)})
				}
			}
		}
	}
	worldWallsCapture(t, "player-visibility", rows, "fdfd24d199a7ad96fced54d2b512796c56fab6166aac3f378a2425599c2996cd")
}
