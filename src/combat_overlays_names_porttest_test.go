//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"testing"
)

// Observe actual glyph drawing while retaining the real font and renderer.
// Measuring a string uses GlyphAdvance, so this records rendered text only.
type combatGlyphFace struct {
	font.Face
	drawn []rune
}

func (f *combatGlyphFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	f.drawn = append(f.drawn, r)
	return f.Face.Glyph(dot, r)
}
func TestCombatOverlayFeedLiteralNames(t *testing.T) {
	type record struct {
		Name   string
		Role   int
		Glyphs string
	}
	var rows []record
	for _, name := range []string{"Alice", "A%%B", "世界"} {
		for role := 0; role < 3; role++ {
			t.Run(fmt.Sprintf("%s/role=%d", name, role), func(t *testing.T) {
				o := newCombatOverlayOwner(t)
				face := &combatGlyphFace{Face: basicfont.Face7x13}
				_, restore := o.c.r.GetFonts().PortTestWindowFont(face, "large")
				t.Cleanup(restore)
				o.players[0].SetName(name)
				var row [6]uint32
				row[role] = 7
				legacy.PortTestCombatFeedRow(&row)
				got := string(face.drawn)
				if got != name {
					t.Fatalf("rendered name glyphs %q want %q", got, name)
				}
				rows = append(rows, record{name, role, got})
			})
		}
	}
	spellbookCapture(t, "combat-overlay-feed-literal-names", rows, "88a778ab89d59a8dea16b4d5353642f3461a837b590f92370a8ba532ffa3a7b0")
}
