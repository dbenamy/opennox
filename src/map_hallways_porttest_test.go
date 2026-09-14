//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func hallwayConnection(dir, x, y, width int) legacy.PortTestPaintSpec {
	s := populationBase()
	s.Globals["occupancyEnabled"] = roomValue(1)
	s.Globals["roomGlobal4"] = roomArg(2)
	s.Records[0].Words[68] = 32
	s.Records[0].Refs[80] = roomArg(6)
	roomSetGeometry(&s.Records[1], 1, 4, 4, 0, 0)
	roomSetGeometry(&s.Records[2], 1, 4, 4, float32(float64(x)*32.526913), float32(float64(y)*32.526913))
	s.Records[1].Refs[56] = roomArg(3)
	s.Records[2].Refs[60] = roomArg(2)
	p := &s.Records[5]
	p.Words[60] = math.Float32bits(130.10765)
	p.Words[64] = math.Float32bits(130.10765)
	p.Words[76] = 1
	p.Refs[148] = roomArg(2)
	mx, my := 2, 0
	switch dir {
	case 1:
		my = 3
	case 2:
		mx, my = 3, 2
	case 3:
		mx, my = 0, 2
	}
	p.Words[80+16*dir] = uint32(mx)
	p.Words[84+16*dir] = uint32(my)
	p.Words[88+16*dir] = uint32(width)
	p.Words[92+16*dir] = 1
	s.Actions = []legacy.PortTestPaintAction{paintAction(31, roomArg(1))}
	return s
}

func TestMapHallwaysBentConnection(t *testing.T) {
	for dir := 0; dir < 4; dir++ {
		for _, offset := range []int{-6, 6} {
			for _, width := range []int{1, 2} {
				x, y := offset, -12
				switch dir {
				case 1:
					y = 12
				case 2:
					x, y = 12, offset
				case 3:
					x, y = -12, offset
				}
				out := populationRun([]legacy.PortTestPaintSpec{hallwayConnection(dir, x, y, width)})
				if len(out) != 1 || !out[0].Intact || !out[0].ControlOK || out[0].Steps[0].Return != 1 {
					t.Fatalf("dir %d offset %d width %d: bent hallway admission", dir, offset, width)
				}
				halls := 0
				for _, r := range out[0].Steps[0].Records {
					if r.Alive && r.Kind == "input" && len(r.Words) == 94 && r.Words[0] >= 2 && r.Words[0] <= 5 {
						halls++
					}
				}
				if halls != 3 {
					t.Fatalf("dir %d offset %d width %d: hallways got %d want 3", dir, offset, width, halls)
				}
			}
		}
	}
}
