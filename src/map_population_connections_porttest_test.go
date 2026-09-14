//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestMapPopulationPrefabConnections(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, gap := range []int{4, 8, 12} {
			for _, width := range []int{1, 2} {
				for _, candidate := range []bool{false, true} {
					s := populationBase()
					s.Globals["occupancyEnabled"] = roomValue(1)
					s.Globals["roomGlobal4"] = roomArg(2)
					s.Records[0].Words[68] = 32
					s.Records[0].Refs[80] = roomArg(6)
					roomSetGeometry(&s.Records[1], 1, 4, 4, 0, 0)
					x, y := 0, 0
					mx, my := 2, 0
					switch dir {
					case 0:
						y = -4 - gap
					case 1:
						y = 4 + gap
						my = 3
					case 2:
						x = 4 + gap
						mx = 3
						my = 2
					case 3:
						x = -4 - gap
						mx = 0
						my = 2
					}
					roomSetGeometry(&s.Records[2], 1, 4, 4, float32(float64(x)*32.526913), float32(float64(y)*32.526913))
					if candidate {
						s.Records[1].Refs[56] = roomArg(3)
						s.Records[2].Refs[60] = roomArg(2)
					}
					p := &s.Records[5]
					p.Words[60] = math.Float32bits(130.10765)
					p.Words[64] = math.Float32bits(130.10765)
					p.Words[76] = 1
					p.Refs[148] = roomArg(2)
					p.Words[80+16*dir] = uint32(mx)
					p.Words[84+16*dir] = uint32(my)
					p.Words[88+16*dir] = uint32(width)
					p.Words[92+16*dir] = 1
					s.Actions = []legacy.PortTestPaintAction{paintAction(31, roomArg(1))}
					cases = append(cases, s)
				}
			}
		}
	}
	out := populationCapture(t, "prefab-connections", cases)
	for i, r := range out {
		step := r.Steps[0]
		want := uint32(i % 2)
		if step.Return != want {
			t.Fatalf("case %d prefab hallway admission", i)
		}
		halls := 0
		rows := map[uint32][]uint32{}
		for _, rec := range step.Records {
			if !rec.Alive {
				continue
			}
			rows[rec.ID] = rec.Words
			if rec.Kind == "input" && len(rec.Words) == 94 && rec.Words[0] >= 2 && rec.Words[0] <= 5 {
				halls++
			}
		}
		if halls != int(want) {
			t.Fatalf("case %d hallway count %d", i, halls)
		}
		if want != 0 {
			prefab := rows[rows[step.Slots[6]][37]]
			dir := i / 12
			if (prefab[54]>>uint(8*dir))&255 != 1 {
				t.Fatalf("case %d missing prefab connection", i)
			}
		}
	}
}
