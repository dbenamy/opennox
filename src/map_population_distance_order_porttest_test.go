//go:build porttest

package opennox

import (
	"math"
	"sort"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapPopulationDistanceOrdering(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	var orders [][]int
	for n := 0; n <= 16; n++ {
		for permutation := 0; permutation < 5; permutation++ {
			for _, themes := range []bool{false, true} {
				if themes && n == 0 {
					continue
				}
				s := populationBase()
				order := make([]int, n)
				distances := make([]float32, n)
				for i := 0; i < n; i++ {
					s.Records = append(s.Records, roomRecord(376))
					r := &s.Records[8+i]
					roomSetGeometry(r, 1, 4, 4, 0, 0)
					d := i
					switch permutation {
					case 1:
						d = n - 1 - i
					case 2:
						d = (i * 7) % max(n, 1)
					case 3:
						d = i / 3
					case 4:
						d = (i + 3) % max(n, 1)
					}
					if themes {
						d = i
					}
					distances[i] = float32(d)
					r.Words[356] = math.Float32bits(distances[i])
					if i+1 < n {
						r.Refs[56] = roomArg(10 + i)
						r.Refs[88] = roomArg(10 + i)
						r.Words[216] = 1
					}
					if i > 0 {
						r.Refs[60] = roomArg(8 + i)
					}
					order[i] = i
				}
				if n > 0 {
					s.Globals["roomGlobal4"] = roomArg(9)
				}
				sort.SliceStable(order, func(i, j int) bool { return distances[order[i]] < distances[order[j]] })
				if themes {
					s.Actions = []legacy.PortTestPaintAction{paintAction(22, roomArg(9))}
				} else {
					s.Actions = []legacy.PortTestPaintAction{paintAction(24)}
				}
				cases = append(cases, s)
				orders = append(orders, order)
			}
		}
	}
	out := populationCapture(t, "distance-ordering", cases)
	for i, r := range out {
		step := r.Steps[0]
		rows := map[uint32][]uint32{}
		for _, rec := range step.Records {
			rows[rec.ID] = rec.Words
		}
		p := step.Globals["dword_5d4594_2487584"]
		prev := uint32(0)
		for rank, index := range orders[i] {
			want := step.Slots[9+index]
			if p != want || rows[p][17] != prev {
				t.Fatalf("case %d rank %d distance list", i, rank)
			}
			if cases[i].Actions[0].Op == 22 {
				theme := uint32(1)
				if rank > 0 {
					pct := 100 * rank / (len(orders[i]) - 1)
					switch {
					case pct < 30:
						theme = 2
					case pct < 60:
						theme = 4
					case pct < 90:
						theme = 8
					default:
						theme = 16
					}
				}
				if rank == len(orders[i])-1 {
					theme = 32
				}
				if rows[p][91] != theme {
					t.Fatalf("case %d rank %d theme", i, rank)
				}
			}
			prev = p
			p = rows[p][16]
		}
		if p != 0 {
			t.Fatalf("case %d trailing distance link", i)
		}
	}
}
