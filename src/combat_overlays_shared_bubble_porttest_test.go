//go:build porttest

package opennox

import (
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
	"strings"
	"testing"
)

func TestCombatOverlaySharedBubbleText(t *testing.T) {
	type record struct{ Text, Glyphs, Pixels string }
	var rows []record
	for _, text := range []string{"AB", "BA"} {
		t.Run(text, func(t *testing.T) {
			o := newChatBubbleRenderOwner(t)
			face := &combatGlyphFace{Face: basicfont.Face7x13}
			t.Cleanup(o.c.r.GetFonts().PortTestDefaultFont(face))
			*o.words["white"] = noxcolor.RGB5551Color(255, 255, 255).Color32()
			*memmap.PtrUint32(0x852978, 4) = noxcolor.RGB5551Color(0, 0, 0).Color32()
			o.create(t, 1234, text, byte(len(text)), uint32(320)|uint32(300)<<16, 20)
			legacy.PortTestChatBubbleDraw(o.c.Viewport())
			got := string(face.drawn)
			if got != strings.Repeat(text, 2) {
				t.Fatalf("highlighted bubble glyphs %q want %q", got, strings.Repeat(text, 2))
			}
			rows = append(rows, record{text, got, effectsPixelHash(o.pix)})
		})
	}
	if len(rows) == 2 && rows[0].Pixels == rows[1].Pixels {
		t.Fatal("equal-width text must produce different glyph pixels")
	}
	spellbookCapture(t, "combat-overlay-shared-bubble-text", rows, "2869508bf1024f25a3de80815f64a6dec513a327cffa2a69da0e363ee7b19159")
}
