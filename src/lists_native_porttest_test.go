//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestIntrusiveListsRootWrappers(t *testing.T) {
	f := legacy.PortTestListsOpen(3)
	defer f.Close()
	if nox_common_list_getFirstSafe_425890(nil) != nil || nox_common_list_getNextSafe_4258A0(nil) != nil {
		t.Fatal("nil list traversal")
	}
	if nox_common_list_getNext_425940(f.Pointer(0)) != nil {
		t.Fatal("zeroed successor")
	}
	nox_common_list_clear_425760(f.Pointer(0))
	f.Init(1, false, 17)
	nox_common_list_append_4258E0(f.Pointer(0), f.Pointer(1))
	if nox_common_list_getFirstSafe_425890(f.Pointer(0)) != f.Pointer(1) || nox_common_list_getNextSafe_4258A0(f.Pointer(1)) != nil {
		t.Fatal("root list traversal")
	}
	got := f.State(0)
	if got[0].Next != 1 || got[0].Prev != 1 || got[1].Next != 0 || got[1].Prev != 0 || got[1].Key != 17 {
		t.Fatal("root list links")
	}
}

func TestPlayerGroupsBoundedNames(t *testing.T) {
	f := legacy.PortTestGroupsOpen()
	defer f.Close()
	for _, n := range []int{0, 1, 9, 10, 11, 128} {
		name := make([]uint16, n)
		for i := range name {
			name[i] = uint16(0xd800 + i)
		}
		f.AddLongName(uint32(n), name)
		rows := f.State()
		got := rows[len(rows)-1].Name
		for i := 0; i < 10; i++ {
			want := uint16(0)
			if i < min(n, 9) {
				want = name[i]
			}
			if got[i] != want {
				t.Fatalf("length %d name[%d] %04x != %04x", n, i, got[i], want)
			}
		}
	}
}
