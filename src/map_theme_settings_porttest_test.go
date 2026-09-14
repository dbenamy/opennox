//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

func TestMapThemeAlgorithms(t *testing.T) {
	keys := []string{"adjacentPortalRate", "debug", "emptyRoomRate", "hallBranchRate", "hallLengthVariance", "hallLimit", "hallRoomRate", "hallWidthVariance", "irregularRoomRate", "mapSize", "mergeRate", "midHallLength", "midHallWidth", "midRoomSize", "roomVariance", "seed", "skeleton", "useDoors"}
	values := []string{"-1", "0", "1", "3", "100", "+12tail", "0x10", "true", "FALSE", "NONE", "HALL_RING", "1.25", "1e3", "nan", "inf", "-0", "2147483647", "2147483648", "-2147483649", "1.25tail"}
	var cases []themeInputCase
	for _, key := range keys {
		for _, upper := range []bool{false, true} {
			for _, value := range values {
				word := key
				if upper {
					word = strings.ToUpper(word)
				}
				s := themeBase()
				for off := 0; off < 84; off += 4 {
					s.Records[0].Words[off] = 0x12345678
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(9, roomArg(1), roomValue(legacy.PortTestThemeFile))}
				cases = append(cases, themeInputCase{word + " " + value + " END ", s})
			}
		}
	}
	for _, text := range []string{"", "END ", "midHallLength ", "midHallLength 12", "unknown 5 END ", "uſeDoors true END ", "midHallLength 1 midHallLength 2 END ", "IF 0% midHallLength 1 ELSE midHallLength 2 ENDIF END "} {
		s := themeBase()
		s.Actions = []legacy.PortTestPaintAction{paintAction(9, roomArg(1), roomValue(legacy.PortTestThemeFile))}
		cases = append(cases, themeInputCase{text, s})
	}
	out := themeCapture(t, "algorithms", cases)
	for i := 0; i < len(keys)*2*len(values); i++ {
		if out[i].Steps[0].Return != 1 {
			t.Fatalf("algorithm case %d rejected", i)
		}
	}
}

func TestMapThemeFrequencyValidation(t *testing.T) {
	var cases []themeInputCase
	for _, flags := range []uint32{0, 1, 2, 4, 8, 16, 32, 63, 64, 128, 255, 0xffffff00} {
		for _, frequency := range []int32{0, 1, 7, -1} {
			for count := 0; count <= 3; count++ {
				s := themeBase()
				for off := 8; off < 32; off += 4 {
					s.Records[5].Words[off] = 0x12345678
				}
				previous := 0
				for i := 0; i < count; i++ {
					r := roomRecord(224)
					r.Words[64] = flags
					r.Words[72] = uint32(frequency)
					if previous != 0 {
						r.Refs[220] = roomArg(previous)
					}
					s.Records = append(s.Records, r)
					previous = len(s.Records)
				}
				if previous != 0 {
					s.Records[5].Refs[0] = roomArg(previous)
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(29, roomArg(6))}
				cases = append(cases, themeInputCase{"", s})
			}
		}
	}
	out := themeCapture(t, "frequency-validation", cases)
	for i, r := range out {
		if i%4 == 0 && r.Steps[0].Return != 0 {
			t.Fatal("empty decoration set accepted")
		}
	}
}
