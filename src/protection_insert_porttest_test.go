//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionInsertABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56f2f0))
	for trial := 0; trial < 500; trial++ {
		seed := int(gen.Int31())
		if trial%7 == 0 {
			seed = -1
		}
		if trial%11 == 0 {
			seed = 0
		}
		n := []int{0, 1, 2, 3, 32, 100, 256}[trial%7]
		key, sum := gen.Uint32(), gen.Uint32()
		if trial%5 == 0 {
			key = 0
		}
		var values [][2]uint32
		for i := 0; i < n; i++ {
			values = append(values, [2]uint32{uint32(i) ^ uint32(seed)*0x9e3779b9, gen.Uint32()})
		}
		got := legacy.PortTestInsert(nil, values, key, sum, seed)
		if len(got) != n {
			t.Fatal("missing insertion snapshots")
		}
		rng := prand.New(seed)
		other := prand.New(seed + 1)
		var order []int
		for i, v := range values {
			if len(order) == 0 {
				order = append(order, i)
			} else {
				j := rng.IntClamp(0, len(order)-1)
				order = append(order, 0)
				copy(order[j+1:], order[j:])
				order[j] = i
			}
			sum ^= v[0] ^ v[1]
			var want [][2]uint32
			for _, j := range order {
				want = append(want, values[j])
			}
			s := got[i]
			if s.Result != 1 || s.Key != key || s.Sum != sum || s.Count != uint16(i+1) || s.LogicIndex != rng.Index() || s.OtherIndex != other.Index() || !s.LinksValid || !reflect.DeepEqual(s.Values, want) {
				t.Fatalf("trial=%d insertion=%d: got=%+v want sum=%08x random=%d values=%v", trial, i, s, sum, rng.Index(), want)
			}
		}
	}
}

func TestProtectionInsertCountBoundary(t *testing.T) {
	for _, n := range []int{32768, 65535} {
		key := uint32(0x87654321)
		sum := ^key
		initial := make([][2]uint32, n)
		for i := range initial {
			initial[i] = [2]uint32{uint32(i + 1), ^uint32(i + 1)}
			sum ^= initial[i][0] ^ initial[i][1]
		}
		value := [2]uint32{0xffffffff, 0x7f800001}
		rng := prand.New(0)
		index := rng.IntClamp(0, n-1)
		want := append([][2]uint32(nil), initial[:index]...)
		want = append(want, value)
		want = append(want, initial[index:]...)
		got := legacy.PortTestInsert(initial, [][2]uint32{value}, key, sum, 0)[0]
		if got.Result != 1 || got.Count != uint16(n+1) || got.Sum != sum^value[0]^value[1] || got.Key != key || got.LogicIndex != rng.Index() || got.OtherIndex != prand.New(1).Index() || !got.LinksValid || !reflect.DeepEqual(got.Values, want) {
			t.Fatalf("count boundary %d failed", n)
		}
	}
}
