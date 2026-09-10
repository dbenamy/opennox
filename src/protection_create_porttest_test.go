//go:build porttest

package opennox

import (
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionCreateABI(t *testing.T) {
	patterns := []uint32{0, 0x80000000, 1, 0x007fffff, 0x00800000, 0x3f800000, 0xbf800000, 0x7f7fffff, 0x7f800000, 0xff800000, 0x7fc00001, 0x7f800001, 0xff800001, 0xffffffff}
	rng := rand.New(rand.NewSource(0x56f280))
	for i := 0; i < 1000; i++ {
		patterns = append(patterns, rng.Uint32())
	}
	for _, value := range patterns {
		for _, key := range []uint32{0, 1, 0x80000000, 0xffffffff} {
			id, sum := rng.Uint32(), rng.Uint32()
			for mode := 0; mode < 3; mode++ {
				got := legacy.PortTestCreate(id, value, key, sum, mode)
				if got.Result != 1 || got.Count != 1 || got.ID != id^key || got.Value != value^key || got.Sum != sum^id^value || !got.LinksValid {
					t.Fatalf("mode=%d bits=%08x id=%08x key=%08x sum=%08x got=%+v", mode, value, id, key, sum, got)
				}
			}
		}
	}
}
