//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

func TestMapThemeDecorProperties(t *testing.T) {
	var cases []themeInputCase
	for _, op := range []int{23, 24, 25, 26, 27, 28} {
		var texts []string
		switch op {
		case 23:
			texts = []string{"NONE ", "none ", "START ", "EARLY ", "MIDDLE ", "LATER ", "LATEST ", "END ", "START+END ", "start+early+middle+later+latest+end ", "+START++END+ ", "+++ ", "START+bad ", "bad+END ", "NONE+START ", ""}
		case 24:
			texts = []string{"-1 ", "0 ", "1 ", "255 ", "256 ", "257 ", "1000 ", "12tail ", "unknown ", ""}
		case 25:
			texts = []string{"COMMON ", "UNCOMMON ", "RARE ", "VERY_RARE ", "HARDLY_EVER ", "common ", "unknown ", "100 ", ""}
		case 26:
			texts = []string{"* * ", "0 1 ", "-1 10 ", "10 2 ", "3 * ", "* 5 ", "bad bad ", "1 ", ""}
		case 27, 28:
			texts = []string{"door ", "* ", "IF ", strings.Repeat("d", 59) + " ", "a\x00b ", "αβ ", "part", ""}
		}
		for _, text := range texts {
			s := themeBase()
			for off := 0; off < 224; off += 4 {
				s.Records[4].Words[off] = 0xa5a5a5a5
			}
			s.Actions = []legacy.PortTestPaintAction{paintAction(op, roomArg(5), roomValue(legacy.PortTestThemeFile))}
			cases = append(cases, themeInputCase{text, s})
		}
	}
	themeCapture(t, "decor-properties", cases)
}
