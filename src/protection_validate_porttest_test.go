//go:build porttest

package opennox

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/memguard"
)

func TestProtectionValidateABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56fb00))
	check := func(initial [][2]uint32, id, key uint32, data []byte, size, checksum uint32, readable bool) {
		t.Helper()
		var before []byte
		if readable {
			before = bytes.Clone(data)
		}
		sum := gen.Uint32()
		want := uint32(0)
		if int32(id) >= 657757279 {
			for _, v := range initial {
				if v[0] == id {
					if v[1] == checksum {
						want = 1
					}
					break
				}
			}
		}
		got := legacy.PortTestProtectionValidate(initial, key, sum, id, data, size)
		values := append([][2]uint32(nil), initial...)
		if got.Result != want || got.Key != key || got.Sum != sum || got.Sequence != 0x87654321 || got.Count != uint16(len(initial)) || got.SwapCount != 0xffffffff || got.RekeyCount != 0xffffffff || got.LogicIndex != prand.New(123).Index() || got.OtherIndex != prand.New(124).Index() || got.FloatState != got.BeforeFloatState || got.FloatRange != got.BeforeFloatRange || !got.NodesAndLinksStable || !reflect.DeepEqual(got.Values, values) {
			t.Fatalf("id=%08x key=%08x size=%d got=%+v want result=%d records=%v", id, key, size, got, want, values)
		}
		if readable && !bytes.Equal(data, before) {
			t.Fatal("input buffer mutated")
		}
	}
	ids := []uint32{0, 657757278, 657757279, 657757280, 0x7fffffff, 0x80000000, 0xffffffff}
	lengths := []int{0, 1, 2, 3, 4, 5, 7, 8, 31, 32, 63, 64, 255, 1024}
	for trial := 0; trial < 500; trial++ {
		n, off := lengths[trial%len(lengths)], trial%8
		backing := make([]byte, n+off)
		gen.Read(backing)
		data := backing[off:]
		checksum := expectedProtectionChecksum(data)
		id, key := ids[trial%len(ids)], gen.Uint32()
		if trial%11 == 0 {
			key = 0
		}
		initial := [][2]uint32{{id, checksum}, {id ^ 0x12345678, gen.Uint32()}, {id, checksum}}
		switch trial % 5 {
		case 0:
			initial = nil
		case 1:
			initial = initial[1:2]
		case 2:
			initial[0][1] ^= 1 // Later matching duplicate must not win.
		case 3:
			initial[0], initial[1] = initial[1], initial[0]
		}
		check(initial, id, key, data, uint32(n), checksum, true)
	}
	for _, size := range []uint32{0, 1, 3, 4, 0x80000000, 0xffffffff} {
		check([][2]uint32{{657757279, 0}}, 657757279, 0xffffffff, nil, size, 0, true)
	}
	// Any attempted read faults. Rejected/missing IDs must short-circuit even
	// with a huge declared length; incomplete words must never be read either.
	guard, free := memguard.New(4)
	defer free()
	for _, id := range []uint32{0, 657757278, 0x80000000, 0xffffffff} {
		check([][2]uint32{{id, 0}}, id, 0x87654321, guard, 0xffffffff, 0, false)
	}
	check([][2]uint32{{657757280, 0}}, 657757279, 0x87654321, guard, 0xffffffff, 0, false)
	for size := uint32(0); size < 4; size++ {
		check([][2]uint32{{657757279, 0}}, 657757279, 0x87654321, guard, size, 0, false)
	}
}
