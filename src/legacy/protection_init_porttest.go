//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
extern uint32_t dword_5d4594_2516356;
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/server"
)

type PortTestProtectionInitCandidate struct {
	Seed        int64
	Random, Key uint32
	FloatState  [40]byte
	FloatRange  [3]uint32
}

type PortTestProtectionInitSnapshot struct {
	Retry                  bool
	Result                 uint32
	Values                 [][2]uint32
	Key, Sum, Sequence     uint32
	Count                  uint16
	FirstSlot, LastSlot    uint32
	SwapCount, RekeyCount  uint32
	LogicIndex, OtherIndex int
	FloatState             [40]byte
	FloatRange             [3]uint32
	LinksValid             bool
	Candidates             []PortTestProtectionInitCandidate
}

// PortTestProtectionInit calls the initializer or its public Go wrapper with
// an isolated manager. Candidate seeds bracket the wall-clock timestamp; this
// fixture established the same behavior when the original C called time(NULL).
func PortTestProtectionInit(frame uint32, seed int, swaps, rekeys uint32, wrapper bool) PortTestProtectionInitSnapshot {
	oldGet := GetServer
	core := new(server.Server)
	core.SetFrame(frame)
	core.Rand.Logic, core.Rand.Other = prand.New(seed), prand.New(seed+1)
	GetServer = func() Server { return &portTestRandomServer{core: core} }

	head, tail := C.dword_5d4594_2516344, C.dword_5d4594_2516352
	key, sum, sequence := C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	firstSlot, lastSlot := memmap.PtrUint32(0x5D4594, 2516340), memmap.PtrUint32(0x5D4594, 2516332)
	oldFirstSlot, oldLastSlot := *firstSlot, *lastSlot
	swapCounter, rekeyCounter := memmap.PtrUint32(0x5D4594, 2516360), memmap.PtrUint32(0x5D4594, 2516364)
	oldSwaps, oldRekeys := *swapCounter, *rekeyCounter
	oldRandom := protectionRandom
	defer func() {
		Sub_56F3B0()
		C.dword_5d4594_2516344, C.dword_5d4594_2516352 = head, tail
		C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = key, sum, sequence
		*count, *firstSlot, *lastSlot = oldCount, oldFirstSlot, oldLastSlot
		*swapCounter, *rekeyCounter = oldSwaps, oldRekeys
		protectionRandom = oldRandom
		GetServer = oldGet
	}()

	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = 0, 0
	C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = 0xdeadbeef, 0x13579bdf, 0x2468ace0
	*count, *firstSlot, *lastSlot = 0x4567, 0x89abcdef, 0x76543210
	*swapCounter, *rekeyCounter = swaps, rekeys
	protectionRandom = protection.Random{State: [5]float64{1, 2, 3, 4, 5}, Span: 0x10203040, Max: 0x50607080, Min: 0x90a0b0c0}

	before := time.Now().Unix()
	var result uint32
	if wrapper {
		Sub_56F1C0()
	} else {
		result = initializeProtection()
	}
	after := time.Now().Unix()

	out := PortTestProtectionInitSnapshot{
		Result:     result,
		Key:        uint32(C.dword_5d4594_2516348),
		Sum:        uint32(C.dword_5d4594_2516328),
		Sequence:   uint32(C.dword_5d4594_2516356),
		Count:      *count,
		FirstSlot:  *firstSlot,
		LastSlot:   *lastSlot,
		SwapCount:  *swapCounter,
		RekeyCount: *rekeyCounter,
		LogicIndex: core.Rand.Logic.Index(),
		OtherIndex: core.Rand.Other.Index(),
		LinksValid: true,
	}
	out.FloatState, out.FloatRange = portTestRandomState(protectionRandom)
	var prev *protection.Record
	for p, n := protectionHead(), 0; p != nil; p, n = p.Next, n+1 {
		if n >= int(out.Count) {
			out.LinksValid = false
			break
		}
		out.LinksValid = out.LinksValid && p.Prev == prev
		out.Values = append(out.Values, [2]uint32{p.ID ^ out.Key, p.Value ^ out.Key})
		prev = p
	}
	out.LinksValid = out.LinksValid && len(out.Values) == int(out.Count) && uint32(uintptr(unsafe.Pointer(prev))) == uint32(C.dword_5d4594_2516352)

	if after < before || after-before > 2 {
		out.Retry = true
		return out
	}
	actualRandom := protectionRandom
	for candidate := before; candidate <= after; candidate++ {
		protectionRandom.Seed(uint32(candidate))
		random := protectionRandom.Draw()
		want := PortTestProtectionInitCandidate{Seed: candidate, Random: random, Key: frame ^ random}
		want.FloatState, want.FloatRange = portTestRandomState(protectionRandom)
		out.Candidates = append(out.Candidates, want)
	}
	protectionRandom = actualRandom
	return out
}
