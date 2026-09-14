//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func growthDoorCase(dir, width, offset int) legacy.PortTestPaintSpec {
	s := growthBase()
	s.Records[0].Words[52] = 100
	s.Globals["roomGlobal4"] = roomArg(2)
	roomSetGeometry(&s.Records[1], 1, 10, 10, 0, 0)
	x, y := int32(offset), int32(-3)
	w, h := int32(width), int32(3)
	switch dir {
	case 1:
		y = 10
	case 2:
		x, y, w, h = 10, int32(offset), 3, int32(width)
	case 3:
		x, y, w, h = -3, int32(offset), 3, int32(width)
	}
	roomSetGeometry(&s.Records[2], int32(dir+2), w, h, float32(float64(x)*32.526913), float32(float64(y)*32.526913))
	s.Records[2].Words[52] = 2
	s.Records[1].Refs[56] = roomArg(3)
	s.Records[2].Refs[60] = roomArg(2)
	s.Actions = []legacy.PortTestPaintAction{paintAction(0)}
	return s
}
func TestMapGrowthDoorProbe(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		cases = append(cases, growthDoorCase(dir, 3, 2))
	}
	out := growthRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("door case %d state", i)
		}
		doors, waypoints := 0, 0
		for _, row := range r.Steps[0].Records {
			if row.Kind == "object" && row.Alive {
				doors++
			}
			if row.Kind == "waypoint" && row.Alive {
				waypoints++
			}
		}
		if doors != 1 || waypoints != 2 {
			t.Fatalf("door case %d got %d doors / %d waypoints", i, doors, waypoints)
		}
	}
}

func TestMapGrowthDoors(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, width := range []int{1, 2, 3, 4, 5, 8} {
			for _, offset := range []int{0, 2, 8} {
				for _, room := range []bool{false, true} {
					for _, flags := range []uint32{0, 2} {
						for _, rate := range []uint32{0, 100} {
							for seed := 0; seed < 2; seed++ {
								s := growthDoorCase(dir, width, offset)
								s.Seed = seed
								s.Records[0].Words[52] = rate
								s.Records[2].Words[52] = flags
								if room {
									s.Records[2].Words[0] = 1
								}
								cases = append(cases, s)
							}
						}
					}
				}
			}
		}
	}
	for dir := 0; dir < 4; dir++ {
		for seed := 0; seed < 16; seed++ {
			s := growthDoorCase(dir, 4, 2)
			s.Seed = seed
			s.Records[0].Words[52] = 50
			cases = append(cases, s)
		}
	}
	growthCapture(t, "doors", cases)
}

func TestMapGrowthDoorTranslations(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, shift := range [][2]int32{{0, 0}, {5, 7}, {-8, -6}, {9, -5}} {
		for dir := 0; dir < 4; dir++ {
			s := growthDoorCase(dir, 3, 2)
			gx, gy := int32(s.Records[2].Words[4]), int32(s.Records[2].Words[8])
			w, h := int32(s.Records[2].Words[12]), int32(s.Records[2].Words[16])
			roomSetGeometry(&s.Records[1], 1, 10, 10, float32(float64(shift[0])*32.526913), float32(float64(shift[1])*32.526913))
			roomSetGeometry(&s.Records[2], int32(dir+2), w, h, float32(float64(gx+shift[0])*32.526913), float32(float64(gy+shift[1])*32.526913))
			cases = append(cases, s)
		}
	}
	out := growthCapture(t, "door-translations", cases)
	positions := make([][2]float32, len(out))
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("translated door case %d state", i)
		}
		doors, waypoints := 0, 0
		for _, row := range r.Steps[0].Records {
			if row.Kind == "object" && row.Alive {
				doors++
				positions[i] = [2]float32{math.Float32frombits(row.Words[14]), math.Float32frombits(row.Words[15])}
			}
			if row.Kind == "waypoint" && row.Alive {
				waypoints++
			}
		}
		if doors != 1 || waypoints != 2 {
			t.Fatalf("translated door case %d: got %d doors / %d waypoints, want 1 / 2", i, doors, waypoints)
		}
		shift := [][2]int32{{0, 0}, {5, 7}, {-8, -6}, {9, -5}}[i/4]
		base := positions[i%4]
		want := [2]float32{base[0] + float32(23*(shift[0]+shift[1])), base[1] + float32(23*(shift[1]-shift[0]))}
		if positions[i] != want {
			t.Fatalf("translated door case %d position got %v want %v", i, positions[i], want)
		}
	}
}
