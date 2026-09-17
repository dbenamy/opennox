//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"testing"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestTeamUIHUDConstruction(t *testing.T) {
	type row struct {
		Kind          string
		Missing       bool
		First, Repeat int
		Loads         []string
		Children      int
		Hidden        bool
		Destroyed     bool
	}
	var rows []row
	for _, kind := range []string{"ctf", "ball"} {
		for _, missing := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/missing%t", kind, missing), func(t *testing.T) {
				o := newTeamUIOwner(t)
				o.missing = missing
				r := row{Kind: kind, Missing: missing, First: teamUICall(kind+"-construct", 0, 0)}
				r.Repeat = teamUICall(kind+"-construct", 0, 0)
				w := o.window(kind)
				if missing {
					if w != nil || r.First != 0 || r.Repeat != 0 {
						t.Fatal("missing resource accepted")
					}
				} else {
					if w == nil || r.First != 1 || r.Repeat != 1 {
						t.Fatal("HUD construction")
					}
					r.Hidden = w.GetFlags().IsHidden()
					if !r.Hidden {
						t.Fatal("constructed HUD must start hidden")
					}
					for i := 0; i < 16; i++ {
						if c := w.ChildByID(uint(8811 + i)); c != nil {
							r.Children++
							if kind == "ctf" && c.DrawData().Tooltip() != "home" {
								t.Fatal("initial tooltip")
							}
						}
					}
					want := 1
					if kind == "ctf" {
						want = 16
					}
					if r.Children != want {
						t.Fatal("HUD child ownership")
					}
					if len(o.loads) != 1 {
						t.Fatal("repeat construction loaded resource twice")
					}
				}
				r.Loads = append([]string(nil), o.loads...)
				teamUICall(kind+"-destroy", 0, 0)
				r.Destroyed = o.window(kind) == nil
				if !r.Destroyed {
					t.Fatal("HUD pointer survived destruction")
				}
				teamUICall(kind+"-destroy", 0, 0)
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "team-ui-hud-construction", rows, "1d25a57cb16dc78949560147b28bc75c2f6437fe32557d0eb521d3cf97623db6")
}
func TestTeamUIFlagTooltips(t *testing.T) {
	o := newTeamUIOwner(t)
	if teamUICall("ctf-construct", 0, 0) != 1 {
		t.Fatal("construct")
	}
	w := o.window("ctf")
	type row struct {
		ID       int
		Selected bool
		State    int
		Text     string
		Stored   byte
	}
	var rows []row
	for id := 1; id <= 16; id++ {
		for _, selected := range []bool{false, true} {
			for _, state := range []int{0, 1, 2, 3, 127, 128, 255} {
				c := w.ChildByID(uint(8810 + id))
				c.Flags &^= 32
				if selected {
					c.Flags |= 32
				}
				c.DrawData().SetTooltip(nil, "unchanged")
				teamUICall("ctf-tooltip", id, state)
				text := c.DrawData().Tooltip()
				want := "unchanged"
				switch state {
				case 0:
					want = "home"
				case 1:
					want = "their flag carried"
					if selected {
						want = "your flag carried"
					}
				case 2:
					want = "away"
				}
				if text != want {
					t.Fatalf("id%d selected%t state%d: %q want %q", id, selected, state, text, want)
				}
				stored := memmap.Uint8(0x5D4594, 1045611+uintptr(id))
				if stored != byte(state) {
					t.Fatal("flag status width")
				}
				rows = append(rows, row{id, selected, state, text, stored})
			}
		}
	}
	spellbookCapture(t, "team-ui-flag-tooltips", rows, "ec5440d608a7cc370babf5903c21058254f3612dcb62f06df69acea94bc3bcad")
}
func TestTeamUIFlagDrawing(t *testing.T) {
	o := newTeamUIOwner(t)
	if teamUICall("ctf-construct", 0, 0) != 1 {
		t.Fatal("construct")
	}
	root := o.window("ctf")
	root.SetPos(image.Pt(0, 0))
	w := root.ChildByID(8811)
	type row struct {
		State    int
		Selected bool
		Pixels   string
		Return   int
	}
	var rows []row
	for _, state := range []int{0, 1, 2, 3, 255} {
		for _, selected := range []bool{false, true} {
			w.Flags &^= 32
			if selected {
				w.Flags |= 32
			}
			teamUICall("ctf-tooltip", 1, state)
			clear(o.pix.Pix)
			empty := effectsPixelHash(o.pix)
			ret := legacy.PortTestTeamUIDraw("ctf", w)
			pix := effectsPixelHash(o.pix)
			if ret != 1 || pix == empty {
				t.Fatal("flag image did not render")
			}
			rows = append(rows, row{state, selected, pix, ret})
		}
	}
	if rows[0].Pixels == rows[2].Pixels || rows[2].Pixels == rows[4].Pixels || rows[0].Pixels == rows[1].Pixels {
		t.Fatal("flag states/border did not discriminate")
	}
	spellbookCapture(t, "team-ui-flag-drawing", rows, "006d32be8acfb657e8499370b8d40da0c942be9d15ea1b738265d82490c84b0a")
}
