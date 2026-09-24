//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestDurationIdentityLifetime(t *testing.T) {
	record, freeRecord := alloc.New(server.DurSpell{})
	defer freeRecord()
	keys := legacy.PortTestDurationCallbackKeys()
	if len(keys) != 53 {
		t.Fatalf("duration callback count: got %d want 53", len(keys))
	}
	seen := make(map[unsafe.Pointer]string, len(keys))
	for _, row := range keys {
		if row.Key == nil {
			t.Fatalf("nil duration callback key: %s", row.Name)
		}
		if old, ok := seen[row.Key]; ok {
			t.Fatalf("duration callback alias: %s and %s", old, row.Name)
		}
		seen[row.Key] = row.Name
		record.Create, record.Update, record.Destroy = row.Key, row.Key, row.Key
		collisionRegistryGrow(128)
		if record.Create != row.Key || record.Update != row.Key || record.Destroy != row.Key {
			t.Fatalf("foreign duration record key changed: %s", row.Name)
		}
	}
	runtime.GC()
	fresh := legacy.PortTestDurationCallbackKeys()
	if len(fresh) != len(keys) {
		t.Fatal("callback table changed length")
	}
	for i, row := range fresh {
		if row.Key == nil || seen[row.Key] != row.Name || row != keys[i] {
			t.Fatalf("duration key changed after GC: %s", row.Name)
		}
	}
}

func TestDurationRawCallbackForwarding(t *testing.T) {
	record, freeRecord := alloc.New(server.DurSpell{})
	defer freeRecord()
	*record = server.DurSpell{ID: 0x1234, Spell: 0x87654321, Level: 0x89abcdef, Field72: -1234567, Field80: 0xfedcba98, Flags88: 0x55aa33cc}
	before := *record
	values := []int32{0, 1, 255, 256, 0x7fffffff, -0x80000000, -1}
	records := []*server.DurSpell{nil, record}
	for _, value := range values {
		for _, p := range records {
			wantArg := uintptr(0)
			if p != nil {
				wantArg = uintptr(unsafe.Pointer(p))
			}
			resultKey, voidKey := legacy.PortTestDurationObserverReset(value)
			if resultKey == nil || voidKey == nil || resultKey == voidKey {
				t.Fatal("invalid duration observer keys")
			}
			got := legacy.PortTestDurationCallResult(resultKey, p)
			if got != value {
				t.Fatalf("raw duration result: got %#x want %#x", uint32(got), uint32(value))
			}
			state := legacy.PortTestDurationObserverSnapshot()
			if state.ResultArgument != wantArg || state.ResultCalls != 1 || state.VoidArgument != 0 || state.VoidCalls != 0 {
				t.Fatalf("result observer state: %+v want arg=%#x one result call", state, wantArg)
			}
			if p != nil && *p != before {
				t.Fatal("result callback modified duration record")
			}

			legacy.PortTestDurationObserverReset(value)
			legacy.PortTestDurationCallDiscard(resultKey, p)
			state = legacy.PortTestDurationObserverSnapshot()
			if state.ResultArgument != wantArg || state.ResultCalls != 1 || state.VoidArgument != 0 || state.VoidCalls != 0 {
				t.Fatalf("discarded-result observer state: %+v want arg=%#x one result call", state, wantArg)
			}
			if p != nil && *p != before {
				t.Fatal("discarded-result callback modified duration record")
			}

			legacy.PortTestDurationObserverReset(value)
			legacy.PortTestDurationCallDiscard(voidKey, p)
			state = legacy.PortTestDurationObserverSnapshot()
			if state.ResultArgument != 0 || state.ResultCalls != 0 || state.VoidArgument != wantArg || state.VoidCalls != 1 {
				t.Fatalf("void observer state: %+v want arg=%#x one void call", state, wantArg)
			}
			if p != nil && *p != before {
				t.Fatal("void callback modified duration record")
			}
		}
	}
}

func TestDurationWallCallbacksKeepLateHooks(t *testing.T) {
	record, freeRecord := alloc.New(server.DurSpell{})
	defer freeRecord()
	createKey := legacy.Get_nox_xxx_spellWallCreate_4FFA90()
	updateKey := legacy.Get_nox_xxx_spellWallUpdate_500070()
	destroyKey := legacy.Get_nox_xxx_spellWallDestroy_500080()
	if createKey == nil || updateKey == nil || destroyKey == nil {
		t.Fatal("nil spell-wall callback key")
	}
	oldCreate, oldUpdate, oldDestroy := legacy.Nox_xxx_spellWallCreate_4FFA90, legacy.Nox_xxx_spellWallUpdate_500070, legacy.Nox_xxx_spellWallDestroy_500080
	defer func() {
		legacy.Nox_xxx_spellWallCreate_4FFA90 = oldCreate
		legacy.Nox_xxx_spellWallUpdate_500070 = oldUpdate
		legacy.Nox_xxx_spellWallDestroy_500080 = oldDestroy
	}()
	for _, value := range []int32{0, 1, 255, 256, 0x7fffffff, -0x80000000, -1} {
		for _, p := range []*server.DurSpell{nil, record} {
			var createGot, updateGot, destroyGot *server.DurSpell
			createCalls, updateCalls, destroyCalls := 0, 0, 0
			legacy.Nox_xxx_spellWallCreate_4FFA90 = func(got *server.DurSpell) int { createGot = got; createCalls++; return int(value) }
			legacy.Nox_xxx_spellWallUpdate_500070 = func(got *server.DurSpell) int { updateGot = got; updateCalls++; return int(^value) }
			legacy.Nox_xxx_spellWallDestroy_500080 = func(got *server.DurSpell) { destroyGot = got; destroyCalls++ }
			collisionRegistryGrow(128)
			if got := legacy.PortTestDurationCallResult(createKey, p); got != value {
				t.Fatal("wall create result", got, value)
			}
			if got := legacy.PortTestDurationCallResult(updateKey, p); got != ^value {
				t.Fatal("wall update result", got, ^value)
			}
			legacy.PortTestDurationCallDiscard(destroyKey, p)
			if createGot != p || updateGot != p || destroyGot != p || createCalls != 1 || updateCalls != 1 || destroyCalls != 1 {
				t.Fatal("wall hook forwarding", createCalls, updateCalls, destroyCalls)
			}
		}
	}
}
