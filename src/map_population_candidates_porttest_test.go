//go:build porttest

package opennox

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapPopulationPrefabCandidateProbe(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for count := 1; count <= 8; count++ {
		s := populationBase()
		for i := 0; i < count; i++ {
			s.Records = append(s.Records, roomRecord(376))
			r := &s.Records[8+i]
			roomSetGeometry(r, 1, 4, 4, float32(i*64), float32(i*32))
			if i+1 < count {
				r.Refs[56] = roomArg(10 + i)
			}
		}
		s.Globals["roomGlobal4"] = roomArg(9)
		s.Actions = []legacy.PortTestPaintAction{paintAction(30, roomArg(6), roomValue(1))}
		cases = append(cases, s)
	}
	out := populationCapture(t, "prefab-candidates", cases)
	for i, res := range out {
		step := res.Steps[0]
		byID := map[uint32][]uint32{}
		for _, r := range step.Records {
			if len(r.Words) == 94 {
				byID[r.ID] = r.Words
			}
		}
		p := step.Return
		seen := map[uint32]bool{}
		count := 0
		for p != 0 {
			if seen[p] || byID[p] == nil {
				t.Fatalf("case %d invalid/cyclic candidate chain %x", i, p)
			}
			seen[p] = true
			want := step.Slots[9+count]
			if p != want {
				t.Fatalf("case %d candidate %d=%x want %x", i, count, p, want)
			}
			p = byID[p][18]
			count++
		}
		want := i + 1
		if want > 6 {
			want = 6
		}
		if count != want {
			t.Fatalf("case %d candidates %d want %d", i, count, want)
		}
	}
}

func TestMapPopulationNearestCandidates(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	var expected [][]int
	for dir := 0; dir < 4; dir++ {
		for count := 0; count <= 13; count++ {
			for seed := 0; seed < 5; seed++ {
				s := populationBase()
				order := rand.New(rand.NewSource(int64(seed))).Perm(count)
				type candidate struct{ slot, distance int }
				var eligible []candidate
				for i, rank := range order {
					s.Records = append(s.Records, roomRecord(376))
					r := &s.Records[8+i]
					x, y := 2, rank*3+4
					if dir == 0 {
						y = -y
					}
					if dir == 2 {
						x, y = y, x
					}
					if dir == 3 {
						x, y = -y, x
					}
					kind := int32(1)
					if i%5 == 4 {
						kind = 2
					}
					roomSetGeometry(r, kind, 4, 4, float32(float64(x-2)*32.526913), float32(float64(y-2)*32.526913))
					if i+1 < count {
						r.Refs[56] = roomArg(10 + i)
					}
					if kind == 1 {
						eligible = append(eligible, candidate{9 + i, x*x + y*y})
					}
				}
				if count > 0 {
					s.Globals["roomGlobal4"] = roomArg(9)
				}
				sort.Slice(eligible, func(i, j int) bool { return eligible[i].distance < eligible[j].distance })
				var want []int
				for _, c := range eligible {
					if len(want) == 6 {
						break
					}
					want = append(want, c.slot)
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(30, roomArg(6), roomValue(int32(dir)))}
				cases = append(cases, s)
				expected = append(expected, want)
			}
		}
	}
	out := populationCapture(t, "nearest-candidates", cases)
	for i, r := range out {
		step := r.Steps[0]
		rows := map[uint32][]uint32{}
		for _, rec := range step.Records {
			rows[rec.ID] = rec.Words
		}
		ptr := step.Return
		for j, slot := range expected[i] {
			if ptr != step.Slots[slot] {
				t.Fatalf("case %d candidate %d=%x want %x", i, j, ptr, step.Slots[slot])
			}
			ptr = rows[ptr][18]
		}
		if ptr != 0 {
			t.Fatalf("case %d extra candidates", i)
		}
	}
}
