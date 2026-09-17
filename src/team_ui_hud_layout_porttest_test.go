//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestTeamUIHUDLayout(t *testing.T) {
	type row struct {
		Kind                    string
		Count, Return           int
		Mode, Maximum, Off, End image.Point
		Visible                 uint32
		Children                []bool
		HiddenAfter, ShownAgain bool
	}
	var rows []row
	for _, mode := range []image.Point{image.Pt(640, 480), image.Pt(1024, 768), image.Pt(201, 133)} {
		for _, kind := range []string{"ctf", "ball"} {
			for _, count := range []int{0, 1, 4, 5, 8, 9, 16} {
				if kind == "ball" && count != 0 {
					continue
				}
				t.Run(fmt.Sprintf("%v/%s/count%d", mode, kind, count), func(t *testing.T) {
					o := newTeamUIOwner(t)
					nox_win_width_game, nox_win_height_game = mode.X, mode.Y
					noxVideoMax = image.Pt(800, 600)
					r := row{Kind: kind, Count: count, Mode: mode, Maximum: noxVideoMax, Return: teamUICall(kind+"-open", 0, count)}
					w := o.window(kind)
					if w == nil {
						t.Fatal("open did not construct HUD")
					}
					r.Off, r.End = w.Offs(), w.End()
					width, height := min(mode.X, noxVideoMax.X), min(mode.Y, noxVideoMax.Y)
					want := image.Pt(width-w.Size().X/3-91, height-120)
					if kind == "ctf" {
						want.X = width - w.Size().X - 91
						if count <= 4 {
							want.X = width - w.Size().X/2 - 91
						}
						switch {
						case count == 0:
							want.Y = 0
						case count <= 4:
							want.Y = height - 40*count
						case count <= 8:
							want.Y = height - w.Size().Y/2
						default:
							want.Y = height - w.Size().Y
						}
						r.Visible = memmap.Uint32(0x5D4594, 1045608)
						for i := 0; i < 16; i++ {
							shown := !w.ChildByID(uint(8811 + i)).GetFlags().IsHidden()
							if shown != (i < count) {
								t.Fatal("visible icon count")
							}
							r.Children = append(r.Children, shown)
						}
						if r.Visible != uint32(boolToIntForTeamUI(count > 0)) {
							t.Fatal("CTF visibility owner")
						}
					} else {
						r.Visible = *o.words["ball-visible"]
						if r.Visible != 1 {
							t.Fatal("ball visibility owner")
						}
					}
					if r.Off != want || r.End != r.Off.Add(w.Size()) {
						t.Fatalf("HUD geometry %v/%v want %v", r.Off, r.End, want)
					}
					if r.Return != 0 {
						t.Fatalf("show result %d", r.Return)
					}
					teamUICall(kind+"-show", 0, 0)
					r.HiddenAfter = w.GetFlags().IsHidden()
					if !r.HiddenAfter {
						t.Fatal("hide request")
					}
					teamUICall(kind+"-show", 0, 1)
					r.ShownAgain = !w.GetFlags().IsHidden()
					if r.ShownAgain != (kind == "ball" || count > 0) {
						t.Fatal("reopen visibility guard")
					}
					teamUICall(kind+"-hide", 0, 0)
					teamUICall(kind+"-show", 0, 1)
					if !w.GetFlags().IsHidden() {
						t.Fatal("explicit close did not clear visibility owner")
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "team-ui-hud-layout", rows, "4a5bf58f05c404c33a0cff6ccfeca08cc6754d207280216f7d15871093443644")
}
func boolToIntForTeamUI(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TestTeamUIFlagSelection(t *testing.T) {
	o := newTeamUIOwner(t)
	teamUICall("ctf-construct", 0, 0)
	w := o.window("ctf")
	*memmap.PtrUint8(0x5D4594, 1045628) = 16
	type row struct {
		Selected int
		Bits     uint32
	}
	var rows []row
	for _, id := range []int{1, 16, 4, 8, 1, 0, 17} {
		teamUICall("ctf-select", id, 0)
		var bits uint32
		for i := 1; i <= 16; i++ {
			if w.ChildByID(uint(8810+i)).GetFlags()&32 != 0 {
				bits |= 1 << uint(i-1)
			}
		}
		var want uint32
		if id >= 1 && id <= 16 {
			want = 1 << uint(id-1)
		}
		if bits != want {
			t.Fatalf("selected %d bits %x want %x", id, bits, want)
		}
		rows = append(rows, row{id, bits})
	}
	spellbookCapture(t, "team-ui-flag-selection", rows, "aef8c9ee8e3ccecd7306e06cca6411115895331599a8f44504905bc1b939076f")
}
func TestTeamUIBallDrawing(t *testing.T) {
	o := newTeamUIOwner(t)
	teamUICall("ball-construct", 0, 0)
	w := o.window("ball")
	w.SetPos(image.Pt(16, 16))
	w.DrawData().SetBackgroundImage(o.images[3])
	type row struct {
		Image, Selected, Border bool
		Return                  int
		Pixels                  string
	}
	var rows []row
	border := memmap.PtrUint32(0x5D4594, 1045648)
	original := *border
	for _, hasImage := range []bool{false, true} {
		for _, selected := range []bool{false, true} {
			for _, hasBorder := range []bool{false, true} {
				w.DrawData().SetBackgroundImage(nil)
				if hasImage {
					w.DrawData().SetBackgroundImage(o.images[3])
				}
				w.Flags &^= 32
				if selected {
					w.Flags |= 32
				}
				*border = 0
				if hasBorder {
					*border = original
				}
				clear(o.pix.Pix)
				empty := effectsPixelHash(o.pix)
				ret := legacy.PortTestTeamUIDraw("ball", w)
				pix := effectsPixelHash(o.pix)
				if ret != 1 || (pix != empty) != (hasImage || selected && hasBorder) {
					t.Fatal("ball HUD image/border guard")
				}
				rows = append(rows, row{hasImage, selected, hasBorder, ret, pix})
			}
		}
	}
	spellbookCapture(t, "team-ui-ball-drawing", rows, "34f8bacc3ddbba173f5ab8a22dbd715e5947da2ecb89e14d21cc6b088fc50794")
}
