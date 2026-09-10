//go:build porttest

package opennox

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionObjectABI(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56fb60))
	check := func(spec legacy.PortTestProtectionObject, initial [][2]uint32, id uint32, repeat int) {
		t.Helper()
		key, sum := gen.Uint32(), gen.Uint32()
		if repeat == 2 {
			key = 0
		}
		var digest uint32
		if !spec.Nil && !spec.GuardObject {
			digest = spec.NetCode ^ uint32(spec.TypeInd)
			if spec.WithHealth {
				digest ^= uint32(spec.HP)
			}
			if spec.InitPresent && spec.TypePresent && spec.TypeInd != 0xffff && int32(spec.InitSize) > 0 {
				if spec.GuardInit { // Positive sizes 1..3 still contain no full word.
					if spec.InitSize >= 4 {
						t.Fatal("unsafe guard fixture")
					}
				} else {
					digest ^= expectedProtectionChecksum(spec.Init[:spec.InitSize])
				}
			}
			if spec.NamePresent {
				name := spec.Name
				if i := bytes.IndexByte(name, 0); i >= 0 {
					name = name[:i]
				}
				digest ^= expectedProtectionChecksum(name)
			}
		}
		for mode := 0; mode < 2; mode++ {
			got := legacy.PortTestObjectProtection(initial, key, sum, id, spec, mode, repeat)
			values := append([][2]uint32(nil), initial...)
			wantSum := sum
			wantResults := make([]uint32, repeat)
			for i := range wantResults {
				wantResults[i] = id // Missing and ineligible both retain the ID.
				if int32(id) >= 657757279 {
					for j := range values {
						if values[j][0] == id {
							if spec.GuardObject {
								t.Fatal("guarded object must never be evaluated")
							}
							values[j][1] ^= digest
							wantSum ^= digest
							wantResults[i] = wantSum
							break
						}
					}
				}
			}
			s := got.State
			if (!spec.GuardObject && got.Digest != digest) || !reflect.DeepEqual(got.Results, wantResults) || s.Result != wantResults[repeat-1] || s.Key != key || s.Sum != wantSum || s.Sequence != 0x87654321 || s.Count != uint16(len(initial)) || s.SwapCount != 0xffffffff || s.RekeyCount != 0xffffffff || s.LogicIndex != prand.New(123).Index() || s.OtherIndex != prand.New(124).Index() || s.FloatState != s.BeforeFloatState || s.FloatRange != s.BeforeFloatRange || !s.NodesAndLinksStable || !got.ObjectUnchanged || !reflect.DeepEqual(s.Values, values) {
				t.Fatalf("mode=%d id=%08x spec=%+v got=%+v want digest=%08x sum=%08x results=%v values=%v", mode, id, spec, got, digest, wantSum, wantResults, values)
			}
		}
	}
	ids := []uint32{0, 657757278, 657757279, 657757280, 0x7fffffff, 0x80000000, 0xffffffff}
	for trial := 0; trial < 500; trial++ {
		spec := legacy.PortTestProtectionObject{Nil: trial%17 == 0, NetCode: gen.Uint32(), TypeInd: uint16(trial % 5), WithHealth: trial%3 != 0, HP: uint16(gen.Uint32()), TypePresent: trial%7 != 0, InitPresent: trial%5 != 0, NamePresent: trial%11 != 0}
		spec.Init = make([]byte, trial%33)
		spec.Name = make([]byte, trial%31)
		gen.Read(spec.Init)
		gen.Read(spec.Name)
		spec.InitSize = uint32(len(spec.Init))
		if trial%4 == 0 && len(spec.Name) > 0 {
			spec.Name[trial%len(spec.Name)] = 0
		}
		id := ids[trial%len(ids)]
		initial := [][2]uint32{{id, gen.Uint32()}, {id ^ 0x12345678, gen.Uint32()}, {id, gen.Uint32()}}
		if trial%9 == 0 {
			initial = nil
		} else if trial%4 == 0 {
			initial = initial[1:2]
		}
		check(spec, initial, id, 1+trial%3)
	}
	check(legacy.PortTestProtectionObject{TypePresent: true, InitPresent: true, InitSize: 5, Init: []byte{1, 2, 3, 4, 255}, TypeInd: 0, NamePresent: true, Name: []byte{'a', 'b', 'c', 'd', 0, 'x', 'y', 'z'}, WithHealth: true, HP: 0x8000, NetCode: 0x80000000}, [][2]uint32{{657757279, 0}}, 657757279, 2)

	for _, size := range []uint32{0, 1, 2, 3, 0x80000000, 0xffffffff} {
		check(legacy.PortTestProtectionObject{TypePresent: true, InitPresent: true, GuardInit: true, InitSize: size, TypeInd: 1, NetCode: 0xffffffff, WithHealth: true, HP: 0xffff}, [][2]uint32{{657757279, 0}}, 657757279, 2)
	}
	check(legacy.PortTestProtectionObject{TypePresent: false, InitPresent: true, GuardInit: true, InitSize: 4, TypeInd: 2}, [][2]uint32{{657757279, 0}}, 657757279, 1)
	check(legacy.PortTestProtectionObject{TypePresent: true, InitPresent: true, GuardInit: true, InitSize: 4, TypeInd: 0xffff}, [][2]uint32{{657757279, 0}}, 657757279, 1)
	for _, id := range ids {
		check(legacy.PortTestProtectionObject{GuardObject: true}, nil, id, 2)
		if int32(id) < 657757279 {
			check(legacy.PortTestProtectionObject{GuardObject: true}, [][2]uint32{{id, 0}}, id, 1)
		}
	}
}
