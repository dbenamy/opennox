//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapPopulationPrefabFit(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for flags := 0; flags < 32; flags++ {
		for _, size := range []int{3, 4, 8, 12} {
			for _, blocked := range []bool{false, true} {
				s := populationCacheBase()
				roomSetGeometry(&s.Records[1], 1, int32(size), int32(size), 0, 0)
				s.Globals["dword_5d4594_2487672"] = roomArg(8)
				s.Globals["dword_5d4594_2487676"] = roomValue(1)
				s.Records[7].Words[60] = uint32(flags)
				if blocked {
					s.Records = append(s.Records, roomRecord(28))
					s.Records[8].Words[12] = math.Float32bits(10000)
					s.Records[8].Words[16] = math.Float32bits(10000)
					s.Records[1].Refs[368] = roomArg(9)
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(3, roomArg(1), roomArg(2), roomArg(5), roomValue(0))}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "prefab-fit", cases)
	for i, r := range out {
		step := r.Steps[0]
		size := cases[i].Records[1].Words[12]
		blocked := len(cases[i].Records) > 8
		// Truncating map dimensions to grid cells can admit an attempted
		// coordinate that the subsequent exact rectangle check rejects.
		// Center placement avoids that random edge case.
		flags := cases[i].Records[7].Words[60]
		if step.Globals["prefabLoadCount"] != step.Return || step.Return > 1 {
			t.Fatalf("case %d fit loader admission", i)
		}
		if (blocked || size < 8) && step.Return != 0 {
			t.Fatalf("case %d invalid fit accepted", i)
		}
		if !blocked && size >= 8 && flags&16 != 0 && step.Return != 1 {
			t.Fatalf("case %d centered fit rejected", i)
		}
		if step.Return == 0 {
			continue
		}
		rows := map[uint32][]uint32{}
		for _, rec := range step.Records {
			rows[rec.ID] = rec.Words
		}
		room := rows[step.Slots[2]]
		exclusion := rows[room[92]]
		if len(exclusion) != 7 || exclusion[0] != 1 || exclusion[5] != 0 {
			t.Fatalf("case %d placement exclusion", i)
		}
		x, y := math.Float32frombits(exclusion[1]), math.Float32frombits(exclusion[2])
		if x < 0 || y < 0 || math.Float32frombits(exclusion[3]) > math.Float32frombits(room[11])+0.5 || math.Float32frombits(exclusion[4]) > math.Float32frombits(room[12])+0.5 {
			t.Fatalf("case %d placement outside room", i)
		}
		if flags&16 == 0 {
			if flags&1 != 0 && y != 0 {
				t.Fatalf("case %d north edge", i)
			}
			if flags&8 != 0 && flags&4 == 0 && x != 0 {
				t.Fatalf("case %d west edge", i)
			}
		}
	}
}

func TestMapPopulationHallwayWaypoints(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for kind := uint32(2); kind <= 5; kind++ {
		for dir := 0; dir < 4; dir++ {
			for _, oneWay := range []bool{false, true} {
				s := populationBase()
				s.Globals["roomGlobal4"] = roomArg(2)
				s.Records[1].Refs[56] = roomArg(3)
				s.Records[2].Refs[60] = roomArg(2)
				roomSetGeometry(&s.Records[2], int32(kind), 4, 4, 260.2153, 260.2153)
				s.Records[1].Refs[88+32*dir] = roomArg(3)
				s.Records[1].Words[216] = 1 << uint(dir*8)
				s.Records[2].Refs[88+32*((dir+2)%4)] = roomArg(2)
				s.Records[2].Words[216] = 1 << uint(((dir+2)%4)*8)
				if oneWay {
					s.Records[0].Words[60] = 1
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(16, roomArg(1)), paintAction(17, roomArg(1))}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "hallway-waypoints", cases)
	count := 0
	for _, r := range out {
		for _, rec := range r.Steps[1].Records {
			if rec.Kind == "waypoint" {
				count++
			}
		}
	}
	if count == 0 {
		t.Fatal("no hallway waypoints created")
	}
}
