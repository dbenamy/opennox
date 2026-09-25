//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestMonsterCallbackBomberForwarding(t *testing.T) {
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 287096)), 640)
	before := bytes.Clone(table)
	restore := legacy.PortTestInstallCallbackTables()
	defer restore()
	key := *memmap.PtrPtr(0x587000, 287236)
	all := legacy.PortTestCallbackTableFunctions()
	restore()
	if !bytes.Equal(table, before) {
		t.Fatal("callback-table fixture did not restore its region")
	}
	if key == nil {
		t.Fatal("nil Bomber dead callback key")
	}
	if len(all) != 26 {
		t.Fatalf("monster callback table keys: got %d want 26", len(all))
	}
	seen := make(map[unsafe.Pointer]bool, len(all))
	for _, other := range all {
		if other == nil || seen[other] {
			t.Fatalf("nil or duplicate monster callback identity %p", other)
		}
		seen[other] = true
	}
	if !seen[key] {
		t.Fatal("Bomber dead key is absent from callback tables")
	}

	u, free := alloc.New(server.Object{})
	defer free()
	old := legacy.Nox_bomberDead_54A150
	defer func() { legacy.Nox_bomberDead_54A150 = old }()

	values := []int32{-0x80000000, -1, 0, 1, 0x7fffffff}
	args := []*server.Object{nil, u}
	var calls int
	var gotArg *server.Object
	for _, value := range values {
		for _, arg := range args {
			// Resolve the callback identity before replacing the hook. This proves
			// the C/native wrapper reads the mutable function variable per call.
			wantArg, wantValue := arg, value
			legacy.Nox_bomberDead_54A150 = func(actual *server.Object) int {
				calls++
				gotArg = actual
				return int(wantValue)
			}
			calls, gotArg = 0, nil
			got := legacy.PortTestMonsterCallbackCallResult(key, arg)
			if got != wantValue || calls != 1 || gotArg != wantArg {
				t.Fatalf("result callback value=%#x arg=%p: got %#x calls=%d arg=%p", uint32(wantValue), wantArg, uint32(got), calls, gotArg)
			}
		}
	}

	for _, arg := range args {
		legacy.Nox_bomberDead_54A150 = func(actual *server.Object) int {
			calls++
			gotArg = actual
			return -0x1234567
		}
		calls, gotArg = 0, nil
		legacy.PortTestMonsterCallbackCallDiscard(key, arg)
		if calls != 1 || gotArg != arg {
			t.Fatalf("void lifecycle route: calls=%d got=%p want=%p", calls, gotArg, arg)
		}
	}
}
