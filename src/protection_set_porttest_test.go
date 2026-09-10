//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionSetABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56f780))
	ids := []uint32{0, 657757278, 657757279, 657757280, 0x7fffffff, 0x80000000, 0xffffffff}
	lengths := []int{0, 1, 2, 3, 4, 5, 8, 31, 64}
	for trial := 0; trial < 500; trial++ {
		id := ids[trial%len(ids)]
		if trial%3 == 0 {
			id = uint32(657757279 + gen.Int31n(0x7fffffff-657757279))
		}
		initial := make([][2]uint32, lengths[trial%len(lengths)])
		for i := range initial {
			initial[i] = [2]uint32{gen.Uint32(), gen.Uint32()}
		}
		// Include deliberate misses, matches at all positions, and duplicates.
		if len(initial) != 0 && trial%4 != 0 {
			initial[trial%len(initial)][0] = id
			if trial%5 == 0 {
				initial[len(initial)-1][0] = id
			}
		}
		key, sum, sequence := gen.Uint32(), gen.Uint32(), gen.Uint32()
		if trial%5 == 0 {
			key = 0
		}
		swapCount, rekeyCount := gen.Uint32(), gen.Uint32()
		if trial%7 == 0 {
			swapCount, rekeyCount = 0xffffffff, 0xffffffff
		}
		frame, floatSeed := gen.Uint32(), gen.Uint32()
		if trial%3 == 0 {
			frame = uint32(trial % 2)
		}
		seed := int(gen.Int31())
		value := gen.Uint32()
		if trial%4 == 0 {
			value = []uint32{0, 0xff, 0x100, 0xffff, 0x10000, 0x80000000, 0xffffffff}[trial%7]
		}
		for mode := 0; mode < 5; mode++ {
			got := legacy.PortTestProtectionSet(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, mode, id, value)
			values := append([][2]uint32(nil), initial...)
			rng, other := prand.New(seed), prand.New(seed+1)
			wantKey, wantSum, wantResult := key, sum, id
			wantSwaps, wantRekeys := swapCount, rekeyCount
			wantFloat, wantRange := got.BeforeFloatState, got.BeforeFloatRange
			if int32(id) >= 657757279 {
				wantResult = 0
				found := -1
				for i, v := range values {
					if v[0] == id {
						found = i
						break
					}
				}
				if found >= 0 {
					v := value
					if mode == 2 {
						v = uint32(uint8(v))
					}
					if mode == 3 {
						v = uint32(uint16(v))
					}
					values[found][1] = v
					wantKey = frame ^ got.ExpectedRandom
					wantResult = wantKey
					for i := 0; i < len(values)/4; i++ {
						a := rng.IntClamp(0, len(values)/2)
						b := rng.IntClamp(len(values)/2+1, len(values)-1)
						values[a], values[b] = values[b], values[a]
						wantSwaps++
					}
					wantRekeys++
					wantSum = ^wantKey
					for _, v := range values {
						wantSum ^= v[0] ^ v[1]
					}
					wantFloat, wantRange = got.ExpectedFloatState, got.ExpectedFloatRange
				}
			}
			if (mode != 4 && got.Result != wantResult) || got.Key != wantKey || got.Sum != wantSum || got.Sequence != sequence || got.Count != uint16(len(values)) || got.SwapCount != wantSwaps || got.RekeyCount != wantRekeys || got.LogicIndex != rng.Index() || got.OtherIndex != other.Index() || got.FloatState != wantFloat || got.FloatRange != wantRange || !got.NodesAndLinksStable || !reflect.DeepEqual(got.Values, values) {
				t.Fatalf("trial=%d mode=%d id=%08x value=%08x got=%+v want result=%08x key=%08x sum=%08x values=%v", trial, mode, id, value, got, wantResult, wantKey, wantSum, values)
			}
		}
	}
}
