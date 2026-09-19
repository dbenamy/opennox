//go:build porttest

package opennox

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestUnitGameplayActionNames(t *testing.T) {
	newCreatureXferOwner(t)
	names := blobdata.PortTestUnitActionNames()
	if len(names) != 72 || names[0] != "ACTION_IDLE" || names[39] != "" || names[71] != "DEPENDENCY_NOT_MOVED" {
		t.Fatal("shipped table anchors")
	}
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 261768), 72)
	for i, name := range names {
		p, free := alloc.CString(name)
		t.Cleanup(free)
		table[i] = uint32(uintptr(unsafe.Pointer(p)))
	}
	var rows []struct {
		ID         int32
		Restricted bool
		Name       string
		Valid      bool
		Index      int32
	}
	ids := []int32{-2147483648, -1, 72, 73, 2147483647}
	for i := int32(0); i < 72; i++ {
		ids = append(ids, i)
	}
	for _, id := range ids {
		for _, restricted := range []bool{false, true} {
			want, valid := "", id >= 0 && id < 72
			if valid {
				want = names[id]
			}
			if restricted {
				valid = true
				want = names[38]
				if id == 0 || id == 4 || id == 10 || id == 3 {
					want = names[id]
				}
			}
			got, ok := legacy.PortTestUnitActionName(id, restricted)
			if got != want || ok != valid {
				t.Fatalf("id=%d restricted=%v: %q/%v want %q/%v", id, restricted, got, ok, want, valid)
			}
			rows = append(rows, struct {
				ID         int32
				Restricted bool
				Name       string
				Valid      bool
				Index      int32
			}{id, restricted, got, ok, 0})
		}
	}
	checkLookups := func() {
		inputs := append([]string{"", "UNKNOWN_ACTION", "action_idle", "ACTION_IDLE_extra", "ACTION_IDLE\x00ignored", "\x00ACTION_IDLE"}, names...)
		for _, name := range names {
			inputs = append(inputs, strings.ToLower(name))
		}
		for _, name := range inputs {
			key := strings.SplitN(name, "\x00", 2)[0]
			for _, restricted := range []bool{false, true} {
				limit, want := 72, int32(0)
				if restricted {
					limit, want = 39, 38
				}
				for i, candidate := range names[:limit] {
					if candidate == key {
						want = int32(i)
						break
					}
				}
				got := legacy.PortTestUnitActionIndex(name, restricted)
				if got != want {
					t.Fatalf("name=%q restricted=%v: %d want %d", name, restricted, got, want)
				}
				rows = append(rows, struct {
					ID         int32
					Restricted bool
					Name       string
					Valid      bool
					Index      int32
				}{0, restricted, name, true, got})
			}
		}
	}
	checkLookups()
	// A controlled duplicate verifies first-match ordering in both search ranges.
	table[7], names[7] = table[3], names[3]
	table[51], names[51] = table[3], names[3]
	checkLookups()
	spellbookCapture(t, "unit-gameplay-action-names", rows, "85d1149befc354744bd00ac5ae551cb0b22c6913d6b6d87c75d4e97a3ea4fc13")
	t.Logf("%d action metadata cases", len(rows))
}
