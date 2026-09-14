//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapThemeExitsPrefabs(t *testing.T) {
	var cases []themeInputCase
	for _, op := range []int{15, 30, 31} {
		var texts []string
		switch op {
		case 15:
			texts = []string{"", "END ", "OBJECT ", "OBJECT door ", "LINKDATA ", "LINKDATA link END ", "unknown END "}
			for _, dir := range []string{"NORTH", "SOUTH", "EAST", "WEST", "north", "unknown"} {
				for count := 1; count <= 4; count++ {
					texts = append(texts, strings.Repeat("OBJECT door "+dir+" ", count)+"LINKDATA link END ")
				}
			}
		case 30:
			texts = []string{"", "END ", "unknown END ", "MUST_OCCUR END ", "FOREACH PaintObject CONTAINS 100 OBJECT item END END ", "AREAMAP ", "AREAMAP name", "AREAMAP name END "}
			for count := 1; count <= 4; count++ {
				for _, prefix := range []string{"", "MUST_OCCUR "} {
					for _, suffix := range []string{"", "FOREACH PaintObject CONTAINS 100 OBJECT item END ", "FOREACH missing CONTAINS * SPELL FIREBALL END "} {
						texts = append(texts, strings.Repeat(prefix+"AREAMAP name "+suffix, count)+"END ")
					}
				}
			}
		case 31:
			texts = []string{"", "name", "name ", "IF 100% name ENDIF ", strings.Repeat("n", 59) + " "}
		}
		for _, input := range texts {
			s := themeBase()
			s.Actions = []legacy.PortTestPaintAction{paintAction(op, roomArg(1), roomValue(legacy.PortTestThemeFile))}
			if op != 15 {
				s.Actions = append(s.Actions, paintAction(32, roomArg(1)))
			}
			cases = append(cases, themeInputCase{input, s})
		}
	}
	themeCapture(t, "exits-prefabs", cases)
}

func TestMapThemeDecorSets(t *testing.T) {
	var cases []themeInputCase
	bodies := []string{"", "CLEAR_COLLIDES ", "WEAPON tier ARMOR tier SPELL FIREBALL ", "OBJECT name 1 2 ", "AREAMAP name 1 2 ", "OBJECT name 1 2 CONTAINS * OBJECT item END ", "AREAMAP name 1 2 FOREACH PaintObject CONTAINS 100 OBJECT item END ", "CONTAINS 100 OBJECT item END ", "FOREACH PaintObject CONTAINS 100 OBJECT item END ", "SPELL FIREBALL CONTAINS 100 OBJECT item END ", "OBJECT name 1 2 FOREACH PaintObject CONTAINS 100 OBJECT item END ", "unknown ", "OBJECT ", "OBJECT name ", "OBJECT name 1 "}
	for _, kind := range []string{"OBJECT", "AREAMAP"} {
		for _, density := range []string{"0", "-0", "0.1", "1.23456789", "1e20", "nan", "inf"} {
			for _, bounds := range []string{"* *", "0 10", "-1 2", "10 2"} {
				bodies = append(bodies, fmt.Sprintf("%s name DENSITY %s %s ", kind, density, bounds))
			}
		}
	}
	for _, bounds := range []string{"0 1 ", "-1 10 ", "3 2 "} {
		for _, body := range bodies {
			s := themeBase()
			s.Actions = []legacy.PortTestPaintAction{paintAction(19, roomArg(5), roomValue(legacy.PortTestThemeFile)), paintAction(19, roomArg(5), roomValue(legacy.PortTestThemeFile))}
			cases = append(cases, themeInputCase{bounds + body + "END 1 2 CLEAR_COLLIDES END ", s})
		}
	}
	for _, text := range []string{"", "1 ", "1 2 ", "1 2 OBJECT name DENSITY ", "1 2 OBJECT name DENSITY 1 ", "1 2 OBJECT name DENSITY 1 * "} {
		s := themeBase()
		s.Actions = []legacy.PortTestPaintAction{paintAction(19, roomArg(5), roomValue(legacy.PortTestThemeFile))}
		cases = append(cases, themeInputCase{text, s})
	}
	themeCapture(t, "decor-sets", cases)
}

func TestMapThemeDecorations(t *testing.T) {
	var cases []themeInputCase
	bodies := []string{"", "WALL_FLOOR wall floor ", "WALL_FLOOR wall floor WALL_FLOOR wall2 floor2 ", "WALL_FLOOR wall floor EDGING edge floor EDGING_BEVELED edge2 floor2 EDGING_IRREGULAR edge3 floor3 ", "EDGING edge floor ", "OCCUR_CONSTRAINT START+END OCCUR_LIMIT 2 MUST_OCCUR FREQUENCY COMMON ROOM_SIZE_CONSTRAINT 1 * DOOR door DOUBLE_DOOR double ", "DECOR_SET 0 1 OBJECT name 1 2 CONTAINS * OBJECT item END END ", "DECOR_SET 0 1 CLEAR_COLLIDES END DECOR_SET 2 3 SPELL FIREBALL END ", "COPY missing ", "unknown ", "WALL_FLOOR ", "WALL_FLOOR wall "}
	for _, kind := range []string{"ROOM", "HALL", "TEMPLATE", "BACKDROP", "room", "unknown"} {
		for _, body := range bodies {
			s := themeBase()
			s.Actions = []legacy.PortTestPaintAction{paintAction(16, roomArg(1), roomValue(legacy.PortTestThemeFile)), paintAction(32, roomArg(1))}
			cases = append(cases, themeInputCase{kind + " name " + body + "END ", s})
		}
	}
	for _, text := range []string{"", "ROOM ", "ROOM name ", "ROOM name WALL_FLOOR wall floor EDGING edge ", "ROOM name DOOR ", "ROOM name FREQUENCY "} {
		s := themeBase()
		s.Actions = []legacy.PortTestPaintAction{paintAction(16, roomArg(1), roomValue(legacy.PortTestThemeFile))}
		cases = append(cases, themeInputCase{text, s})
	}
	themeCapture(t, "decorations", cases)
}

func TestMapThemeDecorCopies(t *testing.T) {
	var cases []themeInputCase
	for _, source := range []string{"ROOM", "HALL", "TEMPLATE", "BACKDROP"} {
		for _, dest := range []string{"ROOM", "HALL", "TEMPLATE", "BACKDROP"} {
			for _, name := range []string{"base", "BASE", "missing"} {
				for count := 0; count <= 3; count++ {
					body := strings.Repeat("WALL_FLOOR wall floor EDGING edge floor DECOR_SET 0 1 OBJECT item 1 2 END ", count)
					text := source + " base " + body + "END " + dest + " target COPY " + name + " END "
					s := themeBase()
					s.Actions = []legacy.PortTestPaintAction{paintAction(16, roomArg(1), roomValue(legacy.PortTestThemeFile)), paintAction(16, roomArg(1), roomValue(legacy.PortTestThemeFile)), paintAction(32, roomArg(1))}
					cases = append(cases, themeInputCase{text, s})
				}
			}
		}
	}
	themeCapture(t, "decor-copies", cases)
}
