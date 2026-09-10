//go:build porttest

package opennox

import (
	"math/bits"
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionBitsetABI(t *testing.T) {
	for i := int32(0); i <= 4096; i++ {
		for _, enabled := range []int32{0, 1, -1} {
			got := legacy.PortTestProtectionBit(i, enabled)
			if enabled == 0 {
				if got != 0 {
					t.Fatal("disabled bit")
				}
				continue
			}
			if bits.OnesCount32(got) != 1 || bits.TrailingZeros32(got) != int(i%32) {
				t.Fatalf("bit %d = %08x", i, got)
			}
		}
	}
	rng := rand.New(rand.NewSource(0x56fce0))
	ids := []int32{-2147483648, -1, 0, 657757278, 657757279, 657757280, 2147483647}
	counts := []int{-100, 0, 1, 2, 31, 32, 33, 64, 65, 137}
	for trial := 0; trial < 5000; trial++ {
		id, count, mode := ids[trial%len(ids)], counts[trial%len(counts)], trial%3
		key, sum := rng.Uint32(), rng.Uint32()
		if trial%7 == 0 {
			key = 0
		}
		var values []int32
		if count > 1 {
			values = make([]int32, count)
			for i := range values {
				if rng.Intn(4) == 0 {
					values[i] = int32(rng.Uint32() | 1)
				}
			}
		}
		var lanes [32]bool
		for i := 1; i < len(values); i++ {
			if values[i] != 0 {
				lanes[i%32] = true
			}
		}
		var flags uint32
		for i, set := range lanes {
			if set {
				flags |= 1 << uint(i)
			}
		}
		value := flags ^ key
		if trial%2 == 0 {
			value ^= 1 << uint(trial%32)
		}
		index, enabled := int32(rng.Intn(1025)), int32(rng.Uint32())
		if trial%5 == 0 {
			enabled = 0
		}
		wantValue, wantSum, wantAward := value, sum, int32(id)
		before, after := 0, 0
		if id >= 657757279 && mode == 1 {
			if value == flags^key {
				before = 1
			}
			decoded := value ^ key
			if enabled != 0 {
				decoded |= 1 << uint(index%32)
			}
			wantValue = decoded ^ key
			wantSum = sum ^ value ^ wantValue
			wantAward = int32(wantValue)
			if wantValue == flags^key {
				after = 1
			}
		}
		got := legacy.PortTestBitset(id, key, value, sum, index, enabled, values, count, mode)
		if got.Award != wantAward || got.Value != wantValue || got.Checksum != wantSum || got.BeforeValid != before || got.AfterValid != after || !got.Unchanged {
			t.Fatalf("trial=%d id=%d count=%d mode=%d got=%+v want award=%08x value=%08x sum=%08x valid=%d/%d", trial, id, count, mode, got, uint32(wantAward), wantValue, wantSum, before, after)
		}
	}
}
