//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestPingAggregates(t *testing.T) {
	type sample struct {
		name   string
		slots  []int
		values map[int][]int
	}
	cases := []sample{
		{"empty", nil, nil},
		{"host-only", []int{31}, map[int][]int{31: {123, 456}}},
		{"single-positive", []int{4}, map[int][]int{4: {8, 8}}},
		{"signed-first-filter", []int{0, 1, 2, 3, 31}, map[int][]int{0: {0, 99}, 1: {-1, 99}, 2: {-2147483648, 99}, 3: {2147483647, 7}, 31: {5, 9}}},
		{"second-read-zero", []int{0, 30}, map[int][]int{0: {7, 0}, 30: {9, 12}}},
		{"unsigned-min", []int{0, 30}, map[int][]int{0: {7, -1}, 30: {9, 12}}},
		{"negative-second", []int{3}, map[int][]int{3: {1, -2147483648}}},
		{"negative-average-truncation", []int{3, 4}, map[int][]int{3: {1, -9}, 4: {1, 2}}},
		{"positive-overflow", []int{3, 4}, map[int][]int{3: {1, 2147483647}, 4: {1, 2147483647}}},
		{"negative-overflow", []int{3, 4}, map[int][]int{3: {1, -2147483648}, 4: {1, -2147483648}}},
		{"sparse-order", []int{30, 1, 17, 31, 0}, map[int][]int{0: {10, 17}, 1: {2, 33}, 17: {9, 50}, 30: {1, 3}, 31: {9, 1}}},
	}
	seed := uint32(0x9182af73)
	next := func() uint32 { seed = seed*1664525 + 1013904223; return seed }
	edge := []int{0, -1, -2147483648, 2147483647, 1, 2, 999, -13}
	for n := 0; n < 256; n++ {
		c := sample{name: fmt.Sprintf("generated-%03d", n), values: make(map[int][]int)}
		for slot := 0; slot < 32; slot++ {
			if next()>>28 < 5 {
				continue
			}
			c.slots = append(c.slots, slot)
			vals := make([]int, 8)
			for i := range vals {
				v := next()
				if v>>28 < 8 {
					vals[i] = edge[(v>>16)%uint32(len(edge))]
				} else {
					vals[i] = int(int32(v))
				}
			}
			c.values[slot] = vals
		}
		cases = append(cases, c)
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Independent invocations cover both ABI and public wrapper before conversion;
			// the final series also checks fresh callback reads across aggregate calls.
			series := [][]legacy.PortTestPingAggregateCall{
				{{Kind: "min"}}, {{Kind: "average"}}, {{Kind: "min", Wrapper: true}}, {{Kind: "average", Wrapper: true}},
				{{Kind: "min"}, {Kind: "average"}, {Kind: "min", Wrapper: true}, {Kind: "average", Wrapper: true}},
			}
			for _, calls := range series {
				got, err := legacy.PortTestPingAggregates(c.slots, c.values, calls)
				if err != nil {
					t.Fatal(err)
				}
				if len(got) != len(calls) {
					t.Fatal("missing results")
				}
				pos := make(map[int]int)
				active := make(map[int]bool)
				for _, slot := range c.slots {
					active[slot] = true
				}
				for j, call := range calls {
					var trace []legacy.PortTestPingTrace
					var included []int
					read := func(slot int) int {
						v := 0
						if pos[slot] < len(c.values[slot]) {
							v = c.values[slot][pos[slot]]
						}
						pos[slot]++
						trace = append(trace, legacy.PortTestPingTrace{Player: slot, Value: v})
						return v
					}
					for slot := 0; slot < 31; slot++ {
						if active[slot] && read(slot) > 0 {
							included = append(included, read(slot))
						}
					}
					want := uint32(0)
					if len(included) > 0 {
						if call.Kind == "min" {
							want = ^uint32(0)
							for _, v := range included {
								if uint32(v) < want {
									want = uint32(v)
								}
							}
						} else {
							var sum int64
							for _, v := range included {
								sum += int64(v)
							}
							want = uint32(int64(int32(sum)) / int64(len(included)))
						}
					}
					if got[j].Value != want || !reflect.DeepEqual(got[j].Trace, trace) || !got[j].PlayersUnchanged {
						t.Fatalf("calls=%+v index=%d got=%+v want=%08x trace=%+v", calls, j, got[j], want, trace)
					}
				}
			}
		})
	}
}
