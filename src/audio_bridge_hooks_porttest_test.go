//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientAudioBridgeMutableHooks(t *testing.T) {
	keys := legacy.PortTestAudioBridgeKeys()
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 93952)), 84)
	beforeTable := bytes.Clone(table)
	mapped, restoreTable := legacy.PortTestAudioBridgeTableBindings()
	defer restoreTable()
	if mapped != keys {
		t.Fatal("audio table field mapping differs from original callback keys")
	}
	seen := make(map[unsafe.Pointer]bool)
	for i, key := range keys {
		if key == nil || seen[key] {
			t.Fatalf("missing or duplicate audio bridge key %d", i)
		}
		seen[key] = true
	}
	p, free := alloc.New(byte(0x5a))
	defer free()
	args := []unsafe.Pointer{nil, unsafe.Pointer(p)}
	hooks := []struct {
		id       int
		original func(unsafe.Pointer) int
		set      func(func(unsafe.Pointer) int)
	}{
		{0, legacy.Sub_43EC30, func(f func(unsafe.Pointer) int) { legacy.Sub_43EC30 = f }},
		{1, legacy.Sub_43ECB0, func(f func(unsafe.Pointer) int) { legacy.Sub_43ECB0 = f }},
		{2, legacy.Sub_43ED00, func(f func(unsafe.Pointer) int) { legacy.Sub_43ED00 = f }},
		{3, legacy.Sub_43EFD0, func(f func(unsafe.Pointer) int) { legacy.Sub_43EFD0 = f }},
		{7, legacy.Sub_43F060, func(f func(unsafe.Pointer) int) { legacy.Sub_43F060 = f }},
		{9, legacy.Sub_43E940, func(f func(unsafe.Pointer) int) { legacy.Sub_43E940 = f }},
		{11, legacy.Sub_43EA20, func(f func(unsafe.Pointer) int) { legacy.Sub_43EA20 = f }},
	}
	for _, hook := range hooks {
		defer hook.set(hook.original)
		for _, arg := range args {
			for _, result := range []int32{-2147483648, -1, 0, 1, 2147483647} {
				var gotArg unsafe.Pointer
				calls := 0
				hook.set(func(p unsafe.Pointer) int { calls++; gotArg = p; return int(result) })
				got := legacy.AudioStreamCallbackInt(keys[hook.id], arg)
				if calls != 1 || gotArg != arg || got != result {
					t.Fatalf("hook %d: calls=%d argument=%p want=%p result=%d want=%d", hook.id, calls, gotArg, arg, got, result)
				}
				legacy.AudioStreamCallbackVoid(keys[hook.id], arg)
				if calls != 2 || gotArg != arg {
					t.Fatalf("hook %d: void dispatch calls=%d argument=%p", hook.id, calls, gotArg)
				}
			}
		}
	}
	oldFree, oldContextFree := legacy.Sub_43E9F0, legacy.Sub_43EC10
	defer func() { legacy.Sub_43E9F0, legacy.Sub_43EC10 = oldFree, oldContextFree }()
	for _, arg := range args {
		freeCalls, contextCalls := 0, 0
		legacy.Sub_43E9F0 = func() { freeCalls++ }
		legacy.Sub_43EC10 = func() int { contextCalls++; return -2147483648 }
		legacy.AudioStreamCallbackVoid(keys[10], arg)
		legacy.AudioStreamCallbackVoid(keys[12], arg)
		if freeCalls != 1 || contextCalls != 1 {
			t.Fatalf("zero-argument hook forwarding: free=%d context=%d", freeCalls, contextCalls)
		}
		if got := legacy.AudioStreamCallbackInt(keys[8], arg); got != 0 {
			t.Fatalf("restart no-op result %d, want zero", got)
		}
	}
	restoreTable()
	if !bytes.Equal(table, beforeTable) {
		t.Fatal("audio table not restored exactly")
	}

}
