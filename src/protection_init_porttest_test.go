//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionInitABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56f1c0))
	frames := []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff}
	for trial := 0; trial < 200; trial++ {
		frame := gen.Uint32()
		if trial < len(frames) {
			frame = frames[trial]
		}
		seed := int(gen.Int31())
		if trial%7 == 0 {
			seed = -1
		}
		swaps, rekeys := gen.Uint32(), gen.Uint32()
		if trial%5 == 0 {
			swaps, rekeys = 0xffffffff, 0xffffffff
		}
		for _, wrapper := range []bool{false, true} {
			var got legacy.PortTestProtectionInitSnapshot
			for retry := 0; retry < 5; retry++ {
				got = legacy.PortTestProtectionInit(frame, seed, swaps, rekeys, wrapper)
				if !got.Retry {
					break
				}
			}
			if got.Retry {
				t.Fatal("clock repeatedly moved during startup fixture")
			}
			rng, other := prand.New(seed), prand.New(seed+1)
			values := make([][2]uint32, 0, 9)
			for i := 0; i < 9; i++ {
				v := [2]uint32{657757279 + uint32(i), 0}
				if i == 8 {
					v[1] = 1
				}
				if i == 0 {
					values = append(values, v)
					continue
				}
				pos := rng.IntClamp(0, len(values)-1)
				values = append(values, [2]uint32{})
				copy(values[pos+1:], values[pos:])
				values[pos] = v
			}
			matched := false
			for _, c := range got.Candidates {
				if got.Key == frame^c.Random && got.FloatState == c.FloatState && got.FloatRange == c.FloatRange {
					matched = true
				}
			}
			sum := ^got.Key
			for _, v := range values {
				sum ^= v[0] ^ v[1]
			}
			result := uint32(657757287)
			if wrapper {
				result = 0
			}
			if !matched || got.Result != result || got.Sum != sum || got.Count != 9 || got.Sequence != 657757288 || got.FirstSlot != 657757279 || got.LastSlot != 657757287 || got.SwapCount != swaps || got.RekeyCount != rekeys || got.LogicIndex != rng.Index() || got.OtherIndex != other.Index() || !got.LinksValid || !reflect.DeepEqual(got.Values, values) {
				t.Fatalf("trial=%d wrapper=%t got=%+v want sum=%08x values=%v", trial, wrapper, got, sum, values)
			}
			// With shipped constants, the carry remains nonzero after the warmup;
			// this prevents accidentally testing an uninitialized all-zero blob.
			if got.FloatState == ([40]byte{}) {
				t.Fatal("floating RNG fixture lost shipped constants")
			}
		}
	}
}
