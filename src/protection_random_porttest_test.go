//go:build porttest

package opennox

import (
	"math"
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionRandomABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56fe30))
	seeds := []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff}
	pairs := [][2]uint32{{0, 0}, {0, 99}, {1, 0xffffffff}, {0xffffffff, 0}, {1, 0}, {0, 0xffffffff}, {0x80000000, 0x7fffffff}, {0x7fffffff, 0x80000000}}
	for trial := 0; trial < 506; trial++ {
		seed := gen.Uint32()
		if trial < len(seeds) {
			seed = seeds[trial]
		}
		ops := make([]legacy.PortTestRandomOp, 90)
		for i := range ops {
			ops[i] = legacy.PortTestRandomOp{Mode: i % 3, A: gen.Uint32(), B: gen.Uint32()}
			if i < len(pairs)*3 {
				p := pairs[(i/3)%len(pairs)]
				ops[i].A, ops[i].B = p[0], p[1]
			}
			if i == 80 || i == 85 {
				ops[i].Mode = 3
			}
		}
		snapshots := legacy.PortTestProtectionRandom(seed, ops)
		var want protection.Random
		want.Seed(seed)
		for i, got := range snapshots {
			var result uint64
			if i > 0 {
				op := ops[i-1]
				switch op.Mode {
				case 0:
					result = math.Float64bits(want.Next())
				case 1:
					result = uint64(want.Range(op.A, op.B))
				case 2:
					result = uint64(want.Draw())
				case 3:
					want.Seed(op.A)
				}
			}
			if got.Result != result || got.Random.Span != want.Span || got.Random.Min != want.Min || got.Random.Max != want.Max {
				t.Fatalf("trial=%d step=%d got=%+v want=%+v result=%x", trial, i, got, want, result)
			}
			for j := range want.State {
				if math.Float64bits(got.Random.State[j]) != math.Float64bits(want.State[j]) {
					t.Fatalf("trial=%d seed=%08x step=%d index=%d got=%016x want=%016x", trial, seed, i, j, math.Float64bits(got.Random.State[j]), math.Float64bits(want.State[j]))
				}
			}
		}
	}
}
