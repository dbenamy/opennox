//go:build porttest

package opennox

import (
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionFloatABI(t *testing.T) {
	if cw := legacy.PortTestProtectionFloatCW(); cw&0x0f00 != 0x0200 {
		t.Fatalf("unexpected x87 precision/rounding: %04x", cw)
	}
	gen := rand.New(rand.NewSource(0x56fa40))
	ids := []uint32{0, 657757278, 657757279, 657757280, 0x7fffffff, 0x80000000, 0xffffffff}
	deltas := []uint32{0, 0x80000000, 1, 0x80000001, 0x007fffff, 0x807fffff, 0x3f000000, 0xbf000000, 0x3fc00000, 0xbfc00000, 0x4f7fffff, 0x4f800000, 0xcf800000, 0x5effffff, 0x5f000000, 0xdeffffff, 0xdf000000, 0xdf000001, 0x7f800000, 0xff800000, 0x7fc00000, 0x7f800001, 0xffc00001}
	for e := -70; e <= -15; e++ {
		v := float32(-math.Ldexp(1, e))
		deltas = append(deltas, math.Float32bits(v), math.Float32bits(math.Nextafter32(v, float32(math.Inf(-1)))), math.Float32bits(math.Nextafter32(v, 0)))
	}
	oldValues := []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff}
	lengths := []int{0, 1, 2, 3, 4, 5, 8, 31, 64}
	for trial := 0; trial < 1800; trial++ {
		id := ids[trial%len(ids)]
		if trial%3 == 0 {
			id = uint32(657757279 + gen.Int31n(0x7fffffff-657757279))
		}
		initial := make([][2]uint32, lengths[trial%len(lengths)])
		for i := range initial {
			initial[i] = [2]uint32{gen.Uint32(), gen.Uint32()}
			if trial < 1400 {
				initial[i][1] = oldValues[(trial/len(deltas))%len(oldValues)]
			}
		}
		// Include deliberate misses, matches at all positions, and duplicates.
		if len(initial) != 0 && trial%4 != 0 {
			initial[trial%len(initial)][0] = id
			if trial%5 == 0 {
				initial[len(initial)-1][0] = id
			}
		}
		if trial < len(deltas)*len(oldValues) {
			id = 657757279
			initial = [][2]uint32{{id, oldValues[trial/len(deltas)]}}
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
		if trial < 1400 {
			value = deltas[trial%len(deltas)]
		}
		for mode := 0; mode < 2; mode++ {
			got := legacy.PortTestProtectionFloat(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, mode, id, value)
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
					old := uint32(0)
					if mode == 1 {
						old = values[found][1]
					}
					values[found][1] = protectionFloatOracle(old, math.Float32frombits(value))
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
			if got.Result != wantResult || got.Key != wantKey || got.Sum != wantSum || got.Sequence != sequence || got.Count != uint16(len(values)) || got.SwapCount != wantSwaps || got.RekeyCount != wantRekeys || got.LogicIndex != rng.Index() || got.OtherIndex != other.Index() || got.FloatState != wantFloat || got.FloatRange != wantRange || !got.NodesAndLinksStable || !reflect.DeepEqual(got.Values, values) {
				t.Fatalf("trial=%d mode=%d id=%08x value=%08x got=%+v want result=%08x key=%08x sum=%08x values=%v", trial, mode, id, value, got, wantResult, wantKey, wantSum, values)
			}
		}
	}
}

// Independent arbitrary-precision model of the observed Go-hosted x87 arithmetic.
func protectionFloatOracle(old uint32, v float32) uint32 {
	x := float64(v)
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0
	}
	a := new(big.Float).SetPrec(53).SetMode(big.ToNearestEven).SetUint64(uint64(old))
	a.Add(a, new(big.Float).SetPrec(53).SetFloat64(x))
	n, _ := a.Int(nil)
	if !n.IsInt64() {
		return 0
	}
	return uint32(n.Int64())
}
