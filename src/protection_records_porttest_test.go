//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionRecordsABI(t *testing.T) {
	rng := rand.New(rand.NewSource(0x56f590))
	for _, n := range []int{0, 1, 2, 7, 64} {
		for _, key := range []uint32{0, 1, 0x80000000, 0xffffffff} {
			for trial := 0; trial < 100; trial++ {
				var values [][2]uint32
				for i := 0; i < n; i++ {
					values = append(values, [2]uint32{rng.Uint32(), rng.Uint32()})
				}
				id := rng.Uint32()
				if n != 0 && trial%2 == 0 {
					id = values[trial%n][0] ^ key
				}
				if n > 1 && trial%3 == 0 {
					values[n-1][0] = values[0][0]
					id = values[0][0] ^ key
				}
				index, a, b := rng.Intn(n+3)-1, rng.Intn(n+2)-1, rng.Intn(n+2)-1
				if trial%7 == 0 {
					a, b = 0, 0
				} // same record must still increment counter
				if trial%11 == 0 {
					index = -2147483648
				}
				if trial%13 == 0 {
					index = 2147483647
				}
				counter := rng.Uint32()
				if trial%5 == 0 {
					counter = 0xffffffff
				}
				wantLookup, wantAt := -1, -1
				for i, v := range values {
					if v[0]^key == id {
						wantLookup = i
						break
					}
				}
				if index >= 0 && index < n {
					wantAt = index
				}
				want := append([][2]uint32(nil), values...)
				wantCounter := counter
				if a >= 0 && a < n && b >= 0 && b < n {
					want[a], want[b] = want[b], want[a]
					wantCounter++
				}
				got := legacy.PortTestRecords(values, key, int(id), index, a, b, counter)
				if got.Lookup != wantLookup || got.At != wantAt || got.Counter != wantCounter || !got.LinksUnchanged || !reflect.DeepEqual(got.Values, want) {
					t.Fatalf("n=%d key=%08x trial=%d index=%d swap=%d,%d: got %+v; want lookup=%d at=%d counter=%d values=%v", n, key, trial, index, a, b, got, wantLookup, wantAt, wantCounter, want)
				}
			}
		}
	}
}
