//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestMapPopulationDensity(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, density := range []float32{-1, 0, 0.0625, 0.25, 1} {
		for seed := 0; seed < 8; seed++ {
			s := populationBase()
			s.Seed = seed
			paintString(&s.Records[5], 4, "PaintObject")
			s.Records[5].Words[64] = 1
			s.Records[5].Words[68] = 2
			s.Records[5].Words[72] = 5
			s.Records[5].Words[76] = math.Float32bits(density)
			s.Actions = []legacy.PortTestPaintAction{paintAction(10, roomArg(1), roomArg(2), roomArg(6))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "density", cases)
	for i, r := range out {
		count := 0
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "object" {
				count++
			}
		}
		want := []int{2, 2, 2, 4, 5}[i/8]
		if count != want {
			t.Fatalf("case %d density created %d want %d", i, count, want)
		}
	}
}

func TestMapPopulationPrefabRetries(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 32; seed++ {
		for _, count := range []uint32{0, 1, 3} {
			for _, blocked := range []bool{false, true} {
				s := populationCacheBase()
				s.Seed = seed
				roomSetGeometry(&s.Records[1], 1, 20, 20, 0, 0)
				s.Records[5] = roomRecord(160)
				s.Records[5].Words[0] = 1
				paintString(&s.Records[5], 4, "fixture")
				s.Records[5].Words[68] = count
				s.Records[5].Words[72] = count
				s.Globals["dword_5d4594_2487672"] = roomArg(8)
				s.Globals["dword_5d4594_2487676"] = roomValue(1)
				paintString(&s.Records[7], 0, "fixture")
				if blocked {
					s.Records = append(s.Records, roomRecord(28))
					s.Records[8].Words[12] = math.Float32bits(10000)
					s.Records[8].Words[16] = math.Float32bits(10000)
					s.Records[1].Refs[368] = roomArg(9)
				}
				clear := paintAction(10, roomArg(1), roomArg(2), roomArg(6))
				clear.Writes = []legacy.PortTestMapRoomWrite{{Slot: 6, Words: map[int]uint32{0: 2}}}
				s.Actions = []legacy.PortTestPaintAction{paintAction(10, roomArg(1), roomArg(2), roomArg(6)), clear}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "prefab-retries", cases)
	successes := 0
	for i, r := range out {
		first, last := r.Steps[0], r.Steps[1]
		loads := first.Globals["prefabLoadCount"]
		if loads > cases[i].Records[5].Words[68] {
			t.Fatalf("case %d excess placements", i)
		}
		if len(cases[i].Records) > 8 && loads != 0 {
			t.Fatalf("case %d ignored exclusion", i)
		}
		successes += int(loads)
		generated := map[uint32]bool{}
		for _, rec := range first.Records {
			if rec.Kind == "exclusion" && len(rec.Words) == 7 && rec.Words[0] == 1 {
				generated[rec.ID] = true
			}
		}
		if len(generated) != int(loads) {
			t.Fatalf("case %d generated exclusion count", i)
		}
		for _, rec := range last.Records {
			if generated[rec.ID] {
				if rec.Alive {
					t.Fatalf("case %d generated exclusion survived clear", i)
				}
				delete(generated, rec.ID)
			}
		}
		if len(generated) != 0 {
			t.Fatalf("case %d lost freed exclusion records", i)
		}
	}
	if successes == 0 {
		t.Fatal("no prefab placed before clearing")
	}
}
