//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionRemoveABI(t *testing.T) {
	rng := rand.New(rand.NewSource(0x56f510))
	for trial := 0; trial < 1000; trial++ {
		n := []int{0, 1, 2, 3, 8, 64}[trial%6]
		key, sum := rng.Uint32(), rng.Uint32()
		if trial%7 == 0 {
			key = 0
		}
		count := uint16(n)
		if trial%11 == 0 {
			count = 0
		}
		if trial%13 == 0 {
			count = 65535
		}
		var values [][2]uint32
		var live []int
		for i := 0; i < n; i++ {
			values = append(values, [2]uint32{rng.Uint32(), rng.Uint32()})
			live = append(live, i)
		}
		if n > 0 && trial%5 == 0 {
			values[0][0] = key
		} // decoded ID zero
		if n > 0 && trial%17 == 0 {
			values[0][0] = ^key
		} // decoded ID UINT_MAX
		if n > 1 && trial%3 == 0 {
			values[n-1][0] = values[0][0]
		}
		ids := []uint32{rng.Uint32()}
		if n > 0 {
			ids = append(ids, values[n/2][0]^key, values[0][0]^key, values[n-1][0]^key, values[0][0]^key)
		}
		got := legacy.PortTestRemove(values, key, sum, count, ids)
		for step, id := range ids {
			result, handle := 0, id
			for j, i := range live {
				if values[i][0]^key == id {
					result = 1
					handle = 0
					sum ^= values[i][0] ^ values[i][1]
					count--
					live = append(live[:j], live[j+1:]...)
					break
				}
			}
			var wantIndices []int
			var wantValues [][2]uint32
			for _, i := range live {
				wantIndices = append(wantIndices, i)
				wantValues = append(wantValues, values[i])
			}
			s := got[step]
			if s.Result != result || s.Handle != handle || s.Sum != sum || s.Key != key || s.Count != count || s.Sequence != 0xdeadbeef || !s.LinksValid || !reflect.DeepEqual(s.Indices, wantIndices) || !reflect.DeepEqual(s.Values, wantValues) {
				t.Fatalf("trial=%d step=%d got=%+v want result=%d handle=%08x sum=%08x count=%d indices=%v values=%v", trial, step, s, result, handle, sum, count, wantIndices, wantValues)
			}
		}
		last := got[len(ids)]
		if last.Sum != 0 || last.Key != 0 || last.Count != 0 || last.Sequence != 0xdeadbeef || len(last.Indices) != 0 || !last.LinksValid {
			t.Fatalf("cleanup trial=%d: %+v", trial, last)
		}
	}
}
