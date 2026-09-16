//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"
)

func listsCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if base := os.Getenv("OPENNOX_LISTS_CAPTURE"); base != "" {
		if err := os.WriteFile(base+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Log(label, got)
	if want != "" && got != want {
		t.Fatalf("%s changed: %s", label, got)
	}
}
func listsCheck(t *testing.T, f *legacy.PortTestLists, head int, want []int) {
	t.Helper()
	var got []int
	for p := f.Op("first", head, -1); p >= 0; p = f.Op("next", p, -1) {
		got = append(got, p)
		if len(got) > len(want) {
			t.Fatal("list cycle or extra member")
		}
	}
	if !reflect.DeepEqual(got, want) && !(len(got) == 0 && len(want) == 0) {
		t.Fatalf("list order %v != %v", got, want)
	}
	var back []int
	for p := f.Op("prev", head, -1); p >= 0; p = f.Op("prev", p, -1) {
		back = append(back, p)
		if len(back) > len(want) {
			t.Fatal("reverse list cycle")
		}
	}
	for i, p := range back {
		if p != want[len(want)-1-i] {
			t.Fatal("reverse order")
		}
	}
	if len(back) != len(want) {
		t.Fatal("reverse count")
	}
	for _, index := range []int{-2, -1, 0, 1, 2, len(want) - 1, len(want), len(want) + 1, 128} {
		expected := -1
		if index >= 0 && index < len(want) {
			expected = want[index]
		}
		if v := f.At(head, index); v != expected {
			t.Fatalf("index %d got %d want %d", index, v, expected)
		}
	}
}
func TestIntrusiveListsOrdering(t *testing.T) {
	type row struct {
		Size, Mode int
		Results    []int
		State      []legacy.PortTestListNodeState
	}
	var rows []row
	for _, size := range []int{0, 1, 2, 7, 32, 128} {
		for mode := 0; mode < 6; mode++ {
			f := legacy.PortTestListsOpen(size + 1)
			f.Init(0, true, 0)
			want := []int(nil)
			r := row{Size: size, Mode: mode}
			keys := []uint32{0, 0x7fffffff, 0x80000000, 0xffffffff, 1, 0xfffffffe, 0}
			for i := 1; i <= size; i++ {
				key := keys[(i*13)%len(keys)]
				f.Init(i, false, key)
				switch mode {
				case 0:
					r.Results = append(r.Results, f.Op("append", 0, i))
					want = append(want, i)
				case 1:
					r.Results = append(r.Results, f.Op("prepend", 0, i))
					want = append([]int{i}, want...)
				default:
					if mode%2 == 0 {
						index := sort.Search(len(want), func(j int) bool { return int32(keys[(want[j]*13)%len(keys)]) >= int32(key) })
						got := f.Op("ascending", 0, i)
						if got != index {
							t.Fatalf("insertion rank %d != %d", got, index)
						}
						r.Results = append(r.Results, got)
						want = append(want, 0)
						copy(want[index+1:], want[index:])
						want[index] = i
					} else {
						index := sort.Search(len(want), func(j int) bool { return int32(keys[(want[j]*13)%len(keys)]) <= int32(key) })
						r.Results = append(r.Results, f.Op("descending", 0, i))
						want = append(want, 0)
						copy(want[index+1:], want[index:])
						want[index] = i
					}
				}
				listsCheck(t, f, 0, want)
				if mode >= 4 && i%3 == 0 {
					at := len(want) / 2
					p := want[at]
					f.Op("remove", p, -1)
					want = append(want[:at], want[at+1:]...)
					listsCheck(t, f, 0, want)
				}
			}
			r.State = f.State(0)
			rows = append(rows, r)
			f.Close()
		}
	}
	listsCapture(t, "ordering", rows, "4e88b199e8043a9b3ac5c6fca53b40df154a1406b0ba7d3d95c7124c590a967d")
}
func TestIntrusiveListsMixed(t *testing.T) {
	type row struct {
		Seed, Step, Result int
		State              []legacy.PortTestListNodeState
	}
	var rows []row
	for seed := 0; seed < 12; seed++ {
		f := legacy.PortTestListsOpen(34)
		f.Init(0, true, 0)
		f.Init(1, true, 0)
		for i := 2; i < 34; i++ {
			f.Init(i, false, uint32(i*17))
		}
		owner := make([]int, 34)
		for i := range owner {
			owner[i] = -1
		}
		var want [2][]int
		rng := rand.New(rand.NewSource(int64(seed)))
		for step := 0; step < 120; step++ {
			id := 2 + rng.Intn(32)
			head := rng.Intn(2)
			result := 0
			if old := owner[id]; old >= 0 {
				f.Op("remove", id, -1)
				for j, v := range want[old] {
					if v == id {
						want[old] = append(want[old][:j], want[old][j+1:]...)
						break
					}
				}
				owner[id] = -1
			}
			switch rng.Intn(4) {
			case 0:
				f.Op("remove", id, -1) // removing an initialized detached item is defined
			case 1:
				result = f.Op("append", head, id)
				want[head] = append(want[head], id)
				owner[id] = head
			case 2:
				result = f.Op("prepend", head, id)
				want[head] = append([]int{id}, want[head]...)
				owner[id] = head
			case 3:
				if len(want[head]) == 0 {
					result = f.Op("append", head, id)
					want[head] = append(want[head], id)
				} else {
					at := rng.Intn(len(want[head]))
					result = f.Op("append", want[head][at], id)
					want[head] = append(want[head], 0)
					copy(want[head][at+1:], want[head][at:])
					want[head][at] = id
				}
				owner[id] = head
			}
			for h := 0; h < 2; h++ {
				listsCheck(t, f, h, want[h])
			}
			rows = append(rows, row{seed, step, result, f.State(0, 1)})
		}
		f.Close()
	}
	listsCapture(t, "mixed", rows, "b69b830142d62fbaa0a99130c31dd5edcfe73013aa963134454df44193d90415")
}
func TestIntrusiveListsNullable(t *testing.T) {
	f := legacy.PortTestListsOpen(3)
	defer f.Close()
	for _, op := range []string{"first", "next"} {
		if f.Op(op, -1, -1) != -1 {
			t.Fatal("nil traversal")
		}
	}
	for _, op := range []string{"first", "next", "raw-next"} {
		if f.Op(op, 0, -1) != -1 {
			t.Fatal("zeroed traversal")
		}
	}
	f.Init(1, false, 17)
	f.Op("append", 0, 1)
	// Preserve the existing partial append to a zeroed list: its forward pointer
	// remains nil, while the backward link and inserted item's forward link change.
	state := f.State()
	if state[0].Next != -1 || state[0].Prev != 1 || state[1].Next != 0 || state[1].Prev != -1 {
		t.Fatal("zeroed-list append")
	}
	listsCapture(t, "nullable", state, "108c8f0ee6ec5f56b57189d4acfd4315b6baa61207c1162f27a2ce133ed13244")
}
func TestPlayerGroupsLifecycle(t *testing.T) {
	type row struct {
		Seed, Step int
		Result     bool
		State      []legacy.PortTestGroupState
	}
	var rows []row
	names := [][]uint16{nil, {'A'}, {'N', 'o', 'x', 0x00e9}, {0xd83c, 0xdf0d}, {'1', '2', '3', '4', '5', '6', '7', '8', '9'}}
	for seed := 0; seed < 8; seed++ {
		f := legacy.PortTestGroupsOpen()
		rng := rand.New(rand.NewSource(int64(seed)))
		var want []legacy.PortTestGroupState
		check := func() {
			t.Helper()
			if f.Initialized() != 1 {
				t.Fatal("group initialization guard")
			}
			got := f.State()
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("seed %d state %v != %v", seed, got, want)
			}
			for id := uint32(0); id < 10; id++ {
				exists := false
				for _, v := range want {
					if v.ID == id {
						exists = true
					}
				}
				if f.Exists(id) != exists || (f.Find(id) != nil) != exists {
					t.Fatal("lookup")
				}
			}
		}
		for step := 0; step < 100; step++ {
			id := uint32(rng.Intn(10))
			at := -1
			for j, v := range want {
				if v.ID == id {
					at = j
					break
				}
			}
			result := true
			switch rng.Intn(6) {
			case 0, 1:
				name := names[rng.Intn(len(names))]
				result = f.Add(id, name)
				if result != (at < 0) {
					t.Fatal("duplicate insertion")
				}
				if at < 0 {
					buf := make([]uint16, 10)
					copy(buf, name)
					want = append(want, legacy.PortTestGroupState{ID: id, Name: buf})
				}
			case 2:
				if at >= 0 {
					index := int32(rng.Intn(7) - 2)
					f.Member(id, index)
					want[at].Members = append(want[at].Members, index)
				}
			case 3:
				if at >= 0 {
					index := int32(rng.Intn(7) - 2)
					result = f.Remove(id, index)
					if !result {
						t.Fatal("remove return identity")
					}
					for j, v := range want[at].Members {
						if v == index {
							want[at].Members = append(want[at].Members[:j], want[at].Members[j+1:]...)
							break
						}
					}
					if len(want[at].Members) == 0 {
						want = append(want[:at], want[at+1:]...)
						if len(want) == 0 {
							want = nil
						}
					}
				}
			case 4:
				f.Init()
			case 5:
				if step%13 == 0 {
					f.Free()
					want = nil
					f.Init()
				}
			}
			check()
			rows = append(rows, row{seed, step, result, f.State()})
		}
		f.Free()
		want = nil
		check()
		f.Close()
	}
	listsCapture(t, "groups", rows, "2928266a95fd7c7e09e9f102596996e941bc0be80da4327b995c07e1cf5c1fe4")
}
func TestPlayerGroupsBoundaries(t *testing.T) {
	f := legacy.PortTestGroupsOpen()
	defer f.Close()
	var rows [][]legacy.PortTestGroupState
	for _, id := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		if !f.Add(id, []uint16{'g', 0x00ff}) {
			t.Fatal("group ID boundary")
		}
		if f.Add(id, []uint16{'x'}) {
			t.Fatal("duplicate ID")
		}
		for _, member := range []int32{-2147483648, -1, 0, 1, 2147483647, 0} {
			f.Member(id, member)
		}
		rows = append(rows, f.State())
		f.Remove(id, 77)
		rows = append(rows, f.State())
		for _, member := range []int32{0, -2147483648, 2147483647, -1, 1, 0} {
			if !f.Remove(id, member) {
				t.Fatal("member return identity")
			}
			rows = append(rows, f.State())
		}
		if f.Exists(id) {
			t.Fatal("empty group retained")
		}
	}
	listsCapture(t, "group-boundaries", rows, "b8638e8868f2506bfb102667a2766b35bfcc52de9479a2d92e2c36963feb82ad")
}
