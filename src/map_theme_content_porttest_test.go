//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/legacy"
)

func TestMapThemeSpellSets(t *testing.T) {
	var cases []themeInputCase
	var all []string
	for id := spell.ID(1); id.Valid(); id++ {
		all = append(all, strings.TrimPrefix(id.String(), "SPELL_"))
	}
	texts := []string{"", "END ", "unknown END ", "FIREBALL MAGIC_MISSILE END ", "fireball fireball END ", strings.Join(all, " ") + " END ", strings.Repeat("FIREBALL ", 140) + "END ", "FIREBALL", strings.Repeat("x", 60) + " END "}
	for _, count := range []uint32{0, 135, 136, 137, 138} {
		for _, text := range texts {
			s := themeBase()
			s.Records[0].Words[1096] = count
			for off := 548; off < 1096; off += 4 {
				s.Records[0].Words[off] = 0x12345678
			}
			s.Actions = []legacy.PortTestPaintAction{paintAction(10, roomArg(1), roomValue(legacy.PortTestThemeFile))}
			cases = append(cases, themeInputCase{text, s})
		}
	}
	themeCapture(t, "spell-sets", cases)
}

func TestMapThemeEquipmentSets(t *testing.T) {
	var cases []themeInputCase
	attrs := []string{"QUALITY", "MATERIAL", "PRIMARY_ENCHANTMENT", "SECONDARY_ENCHANTMENT"}
	for _, kind := range []string{"WEAPON", "ARMOR"} {
		for _, template := range []bool{false, true} {
			for mask := 0; mask < 16; mask++ {
				for count := 0; count < 4; count++ {
					var b strings.Builder
					if template {
						b.WriteString(kind + " TEMPLATE ")
						for _, attr := range attrs {
							b.WriteString(attr + " one two END ")
						}
						b.WriteString("END ")
					}
					for item := 0; item < count; item++ {
						fmt.Fprintf(&b, "%s tier%d item%d ", kind, item, item)
						for slot, attr := range attrs {
							if mask&(1<<slot) != 0 {
								b.WriteString(attr + " one -two three END ")
							}
						}
						b.WriteString("END ")
					}
					b.WriteString("END ")
					s := themeBase()
					op := 11
					if kind == "ARMOR" {
						op = 14
					}
					s.Actions = []legacy.PortTestPaintAction{paintAction(op, roomArg(1), roomValue(legacy.PortTestThemeFile)), paintAction(32, roomArg(1))}
					cases = append(cases, themeInputCase{b.String(), s})
				}
			}
		}
	}
	out := themeCapture(t, "equipment-sets", cases)
	for i, r := range out {
		step := r.Steps[0]
		if step.Return != 1 {
			t.Fatalf("equipment case %d parsing", i)
		}
		index := 276
		if i >= 128 {
			index = 278
		}
		for _, row := range step.Records {
			if row.ID == step.Slots[1] && row.Words[index] != uint32(i%4) {
				t.Fatalf("equipment case %d list count", i)
			}
		}
		for _, row := range r.Steps[1].Records {
			if row.Kind == "themeAllocation" && row.Alive {
				t.Fatalf("equipment case %d cleanup", i)
			}
		}
	}
}

func TestMapThemeEquipmentBoundaries(t *testing.T) {
	var cases []themeInputCase
	for _, slot := range []string{"QUALITY", "EFFECTIVENESS", "MATERIAL", "PRIMARY_ENCHANTMENT", "SECONDARY_ENCHANTMENT"} {
		for _, count := range []int{0, 1, 255, 256} {
			var b strings.Builder
			b.WriteString(slot + " ")
			for i := 0; i < count; i++ {
				fmt.Fprintf(&b, "entry%d ", i)
			}
			b.WriteString("END END ")
			s := themeBase()
			s.Actions = []legacy.PortTestPaintAction{paintAction(13, roomArg(4), roomValue(legacy.PortTestThemeFile)), paintAction(12, roomArg(4))}
			cases = append(cases, themeInputCase{b.String(), s})
		}
	}
	for _, text := range []string{"", "QUALITY ", "QUALITY one ", "QUALITY one END ", "UNKNOWN one END END ", "QUALITY " + strings.Repeat("x", 59) + " END END "} {
		s := themeBase()
		s.Actions = []legacy.PortTestPaintAction{paintAction(13, roomArg(4), roomValue(legacy.PortTestThemeFile))}
		cases = append(cases, themeInputCase{text, s})
	}
	themeCapture(t, "equipment-boundaries", cases)
}

func TestMapThemeChoices(t *testing.T) {
	kinds := []string{"OBJECT", "AREAMAP", "CLEAR_COLLIDES", "WEAPON", "ARMOR", "SPELL"}
	var cases []themeInputCase
	for _, kind := range kinds {
		for _, weight := range []string{"*", "0", "1", "99", "100", "-1"} {
			for _, count := range []int{0, 1, 2, 31, 32, 33} {
				var b strings.Builder
				b.WriteString(weight + " ")
				for i := 0; i < count; i++ {
					fmt.Fprintf(&b, "%s entry%d ", kind, i)
				}
				b.WriteString("END ")
				s := themeBase()
				s.Actions = []legacy.PortTestPaintAction{paintAction(20, roomValue(legacy.PortTestThemeFile))}
				cases = append(cases, themeInputCase{b.String(), s})
			}
		}
	}
	for _, text := range []string{"* OBJECT a OR * OBJECT b END ", "* OBJECT a OR * OBJECT b OR * OBJECT c END ", "100 OBJECT a OR * OBJECT b END ", "* OBJECT a OR -5 OBJECT b END ", "", "100 ", "100 OBJECT ", "100 UNKNOWN name END ", "100 OBJECT a OR 50 UNKNOWN b END ", "100 OBJECT " + strings.Repeat("x", 59) + " END "} {
		s := themeBase()
		s.Actions = []legacy.PortTestPaintAction{paintAction(20, roomValue(legacy.PortTestThemeFile))}
		cases = append(cases, themeInputCase{text, s})
	}
	out := themeCapture(t, "choices", cases)
	for i := 0; i < len(kinds)*6*6; i++ {
		if out[i].Steps[0].Return == 0 {
			t.Fatalf("choice case %d rejected", i)
		}
	}
}

func TestMapThemeForeach(t *testing.T) {
	var cases []themeInputCase
	for _, name := range []string{"PaintObject", "PaintDoor", "paintobject", "missing"} {
		for _, suffix := range []string{"CONTAINS 100 OBJECT item END ", "CONTAINS * SPELL FIREBALL OR * ARMOR armor END ", "CONTAINS 100 UNKNOWN bad END ", "CONTAINS ", "wrong ", ""} {
			s := themeBase()
			s.Actions = []legacy.PortTestPaintAction{paintAction(21, roomValue(legacy.PortTestThemeFile))}
			cases = append(cases, themeInputCase{name + " " + suffix, s})
		}
	}
	themeCapture(t, "foreach", cases)
}
