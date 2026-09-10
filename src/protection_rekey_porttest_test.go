//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionRekeyABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56f5c0))
	lengths := []int{0, 1, 2, 3, 4, 5, 7, 8, 31, 32, 63, 128, 257}
	frames := []uint32{0, 1, 0xffffffff}
	for trial := 0; trial < 400; trial++ {
		key, sum, sequence := gen.Uint32(), gen.Uint32(), gen.Uint32()
		if trial%7 == 0 {
			key = 0
		}
		swapCount, rekeyCount := gen.Uint32(), gen.Uint32()
		if trial%11 == 0 {
			swapCount, rekeyCount = 0xffffffff, 0xffffffff
		}
		seed := int(gen.Int31())
		if trial%9 == 0 {
			seed = -1
		}
		floatSeed := gen.Uint32()
		initial := make([][2]uint32, lengths[trial%len(lengths)])
		for i := range initial {
			initial[i] = [2]uint32{gen.Uint32(), gen.Uint32()}
		}
		frame := frames[trial%len(frames)]
		for _, wrapper := range []bool{false, true} {
			got := legacy.PortTestRekey(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, wrapper)
			rng, other := prand.New(seed), prand.New(seed+1)
			values := append([][2]uint32(nil), initial...)
			swaps := len(values) >> 2
			for i := 0; i < swaps; i++ {
				a := rng.IntClamp(0, len(values)>>1)
				b := rng.IntClamp((len(values)>>1)+1, len(values)-1)
				values[a], values[b] = values[b], values[a]
			}
			newKey := frame ^ got.ExpectedRandom
			wantSum := ^newKey
			for _, value := range values {
				wantSum ^= value[0] ^ newKey
				wantSum ^= value[1] ^ newKey
			}
			if (!wrapper && got.Result != newKey) || got.Key != newKey || got.Sum != wantSum || got.Sequence != sequence || got.Count != uint16(len(initial)) || got.SwapCount != swapCount+uint32(swaps) || got.RekeyCount != rekeyCount+1 || got.LogicIndex != rng.Index() || got.OtherIndex != other.Index() || got.FloatState != got.ExpectedFloatState || got.FloatRange != got.ExpectedFloatRange || !got.LinksValid || !got.NodesAndLinksStable || !reflect.DeepEqual(got.Values, values) {
				t.Fatalf("trial=%d wrapper=%v got=%+v want key=%08x sum=%08x swaps=%d values=%v", trial, wrapper, got, newKey, wantSum, swaps, values)
			}
		}
	}
}
