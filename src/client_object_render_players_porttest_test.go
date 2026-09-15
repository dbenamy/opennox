//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientObjectRenderPlayersAndInvisibility(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		Local, Reference, Info, Observer, Team, Invisible, See, Move, Delta int
		Render                                                              objectRenderResult
	}
	var out []result
	id := 0
	for isLocal := 0; isLocal < 2; isLocal++ {
		for reference := 0; reference < 2; reference++ {
			for info := 0; info < 3; info++ {
				for observer := 0; observer < 2; observer++ {
					for team := 0; team < 4; team++ {
						for invisible := 0; invisible < 2; invisible++ {
							for see := 0; see < 2; see++ {
								for _, move := range []int{0, 1, 3, 4, 7, 8, 16} {
									for _, delta := range []int{0, 1, 2, 8} {
										id++
										start := uint32(120)
										if id%2 != 0 {
											start = 0xfffffff0
										}
										o.resetRender(uint32(id), start)
										local := o.drawable(7, image.Pt(48, 48))
										local.ObjClass = 4
										dr := local
										targetInfo := &o.players[0]
										if isLocal == 0 {
											dr = o.drawable(8, image.Pt(48, 48))
											dr.ObjClass = 4
											targetInfo = &o.players[1]
										}
										o.players[0].Field3680 = uint32(observer)
										targetInfo.Active = 1
										if info == 0 {
											targetInfo.Active = 0
										}
										targetInfo.Field3680 = 0
										if info == 2 {
											targetInfo.Field3680 = 1
										}
										if reference != 0 {
											*memmap.PtrPtr(0x852978, 8) = local.C()
										}
										local.Buffs = uint32(see) << 21
										dr.Buffs |= uint32(invisible)
										switch team {
										case 1:
											local.TeamVal.ID = 1
											dr.TeamVal.ID = 1
										case 2:
											local.TeamVal.ID = 1
											dr.TeamVal.ID = 2
										case 3:
											dr.TeamVal.ID = 1
										}
										dr.Field_8, dr.Field_9 = uint32(48-move), uint32(48+move*(id%2))
										dr.Field_5, dr.Field_10 = start, start-uint32(delta)
										dr.AnimInd = uint32(id%2) * 8
										legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[id%32].C()))
										if isLocal == 0 && info == 2 && len(o.drawTrace) != 0 {
											t.Fatal("remote observer became visible")
										}
										if info != 2 && invisible == 0 && len(o.drawTrace) != 1 {
											t.Fatal("ordinary visible player did not draw")
										}
										out = append(out, result{isLocal, reference, info, observer, team, invisible, see, move, delta, o.renderResult(t, id, 0, 0)})
									}
								}
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-render-player-invisibility", out, len(out), "6d938e4401b4a7eebd9e7e80279c08ef8f8c9982359f7cdefa0f5b5109b8fc0d")
}
