//go:build porttest

package opennox

import (
	"fmt"
	"golang.org/x/image/font/basicfont"
	"testing"
)

func TestTeamUIPlayerListConstruction(t *testing.T) {
	type row struct {
		Lang, Height   int
		Resource       string
		Players, Teams int
		Closed         bool
	}
	var records []row
	for lang := 0; lang < 9; lang++ {
		for _, height := range []int{8, 10, 11, 13} {
			t.Run(fmt.Sprintf("%d/%d", lang, height), func(t *testing.T) {
				o := newTeamUIOwner(t)
				o.configureLanguage(lang)
				face := *basicfont.Face7x13
				face.Height = height
				face.Ascent = height
				t.Cleanup(o.c.Render().GetFonts().PortTestDefaultFont(&face))
				_, restore := o.c.Render().GetFonts().PortTestWindowFont(&face, "small", "large")
				t.Cleanup(restore)
				for pass := 0; pass < 2; pass++ {
					w := o.openPlayers(t)
					wantLang := lang
					if height > 10 {
						wantLang = 2
					}
					if got := o.loads[len(o.loads)-1]; got != fmt.Sprintf("team-ui-player-%d.wnd", wantLang) {
						t.Fatalf("resource %q for lang %d height %d", got, lang, height)
					}
					if w.ChildByID(10501) == nil || w.ChildByID(10502) == nil {
						t.Fatal("missing list widgets")
					}
					players, teams := len(teamUIRowNames(w.ChildByID(10501))), len(teamUIRowNames(w.ChildByID(10502)))
					teamUICall("players-destroy", 0, 1)
					if o.window("players") != nil {
						t.Fatal("destroy left root pointer")
					}
					records = append(records, row{lang, height, o.loads[len(o.loads)-1], players, teams, o.window("players") == nil})
				}
			})
		}
	}
	spellbookCapture(t, "team-ui-list-construction", records, "2c01f4954f4a9c1716c6306e07c5d0716380db1a43eca3914c5279a683434f91")
}
