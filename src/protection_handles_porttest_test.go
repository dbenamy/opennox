//go:build porttest

package opennox

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionHandlesABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56f250))
	sequences := []uint32{0, 1, 657757279, 0x7fffffff, 0x80000000, 0xfffffffc, 0xffffffff}
	for trial := 0; trial < 400; trial++ {
		seed := int(gen.Int31())
		if trial%9 == 0 {
			seed = -1
		}
		key, sum := gen.Uint32(), gen.Uint32()
		if trial%7 == 0 {
			key = 0
		}
		sequenceIndex := trial % (len(sequences) + 1)
		sequence := uint32(0)
		if sequenceIndex == len(sequences) {
			sequence = gen.Uint32()
		} else {
			sequence = sequences[sequenceIndex]
		}
		initial := make([][2]uint32, trial%8)
		for i := range initial {
			initial[i] = [2]uint32{gen.Uint32(), gen.Uint32()}
		}
		ops := make([]legacy.PortTestHandleOp, 1+trial%19)
		for i := range ops {
			ops[i] = legacy.PortTestHandleOp{Reserved: (trial+i)%4 == 0, Value: int32(gen.Uint32())}
			if i%5 == 0 {
				ops[i].Value = int32(i - trial)
			}
		}
		got := legacy.PortTestHandles(initial, key, sum, sequence, seed, ops)
		if len(got) != len(ops) {
			t.Fatal("missing handle snapshots")
		}
		rng, other := prand.New(seed), prand.New(seed+1)
		values := append([][2]uint32(nil), initial...)
		for step, op := range ops {
			result := uint32(1)
			creates := 1
			if op.Reserved {
				creates = 7
			}
			for i := 0; i < creates; i++ {
				value := uint32(op.Value)
				if op.Reserved {
					value = 0
				}
				id := sequence
				if !op.Reserved {
					result = id
				}
				if len(values) == 0 {
					values = append(values, [2]uint32{id, value})
				} else {
					index := rng.IntClamp(0, len(values)-1)
					values = append(values, [2]uint32{})
					copy(values[index+1:], values[index:])
					values[index] = [2]uint32{id, value}
				}
				sum ^= id ^ value
				sequence++
			}
			s := got[step]
			if s.Result != result || s.Sum != sum || s.Key != key || s.Sequence != sequence || s.Count != uint16(len(values)) || s.LogicIndex != rng.Index() || s.OtherIndex != other.Index() || !s.LinksValid || !reflect.DeepEqual(s.Values, values) {
				t.Fatalf("trial=%d step=%d op=%+v got=%+v want result=%08x sum=%08x sequence=%08x count=%d random=%d values=%v", trial, step, op, s, result, sum, sequence, len(values), rng.Index(), values)
			}
		}
	}
}
