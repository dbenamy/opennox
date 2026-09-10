//go:build porttest

package opennox

import (
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestNetworkExtentABI(t *testing.T) {
	check := func(specs []legacy.PortTestExtentSpec, codes []uint32, noServer bool) {
		t.Helper()
		got := legacy.PortTestNetworkExtent(specs, codes, noServer)
		if len(got) != len(codes) {
			t.Fatal("missing extent snapshots")
		}
		find := func(code uint32) int {
			for i, s := range specs {
				if s.Extent == code && (s.Flags/32)%2 == 0 {
					return i
				}
			}
			return -1
		}
		for i, code := range codes {
			dynamic := code
			if (code/32768)%2 != 0 {
				dynamic = 0
				if n := find(code - 32768); n >= 0 {
					dynamic = specs[n].NetCode
				}
			}
			found := -3
			if !noServer {
				found = find(code)
			}
			if got[i].DynamicResult != dynamic || got[i].FoundIndex != found || !got[i].ObjectsUnchanged {
				t.Fatalf("code=%08x noServer=%t specs=%+v got=%+v want dynamic=%08x found=%d", code, noServer, specs, got[i], dynamic, found)
			}
		}
	}
	// Every unmarked low-word code, including upper-word data, must bypass the server.
	var plain []uint32
	for low := uint32(0); low < 32768; low++ {
		for _, hi := range []uint32{0, 0x10000, 0x7fff0000, 0x80000000, 0xffff0000} {
			plain = append(plain, hi+low)
		}
	}
	check(nil, plain, true)
	fixed := []legacy.PortTestExtentSpec{{Extent: 7, NetCode: 111, Flags: 32}, {Extent: 7, NetCode: 222}, {Extent: 7, NetCode: 333}, {Extent: 0, NetCode: 0xffffffff}, {Extent: 0x80000000, NetCode: 0x7fffffff}, {Extent: 0xffff7fff, NetCode: 0x80000000}, {Extent: 0xffffffff, NetCode: 0}, {Extent: 0x7fff, NetCode: 0xffff}}
	codes := []uint32{0, 1, 7, 0x7fff, 0x8000, 0x8007, 0xffff, 0x10000, 0x17fff, 0x18000, 0x80000000, 0x80008000, 0xffffffff}
	check(nil, codes, false)
	check(fixed, codes, false)
	gen := rand.New(rand.NewSource(0x578b40))
	for trial := 0; trial < 400; trial++ {
		specs := make([]legacy.PortTestExtentSpec, trial%33)
		input := append([]uint32(nil), codes...)
		for i := range specs {
			extent := gen.Uint32()
			if i%2 == 0 {
				extent &^= 0x8000
			}
			if i > 0 && i%3 == 0 {
				extent = specs[i-1].Extent
			}
			flags := gen.Uint32()
			if i%4 == 0 {
				flags = 32
			}
			specs[i] = legacy.PortTestExtentSpec{Extent: extent, NetCode: gen.Uint32(), Flags: flags}
			input = append(input, extent, extent|0x8000)
		}
		for i := 0; i < 16; i++ {
			input = append(input, gen.Uint32())
		}
		check(specs, input, false)
	}
}
