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
extern uint32_t dword_5d4594_2516372;
extern uint32_t dword_5d4594_2516380;
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

// PortTestProtectionInit calls the original C initializer (or its public Go
// wrapper) with an otherwise isolated protection manager. Its time-derived
// floating RNG output is reported as candidates for the enclosing test to
// match, because production C calls time(NULL) directly.
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
	floatState := memmap.PtrOff(0x5D4594, 2516388)
	oldFloat := append([]byte(nil), unsafe.Slice((*byte)(floatState), 40)...)
	floatMax := memmap.PtrUint32(0x5D4594, 2516376)
	oldRange := [3]uint32{uint32(C.dword_5d4594_2516372), *floatMax, uint32(C.dword_5d4594_2516380)}
	floatConstants := [5]*float64{
		memmap.PtrFloat64(0x581450, 11344),
		memmap.PtrFloat64(0x581450, 11352),
		memmap.PtrFloat64(0x581450, 11360),
		memmap.PtrFloat64(0x581450, 11368),
		memmap.PtrFloat64(0x581450, 11376),
	}
	var oldFloatConstants [5]float64
	for i, p := range floatConstants {
		oldFloatConstants[i] = *p
	}
	defer func() {
		Sub_56F3B0()
		C.dword_5d4594_2516344, C.dword_5d4594_2516352 = head, tail
		C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = key, sum, sequence
		*count, *firstSlot, *lastSlot = oldCount, oldFirstSlot, oldLastSlot
		*swapCounter, *rekeyCounter = oldSwaps, oldRekeys
		copy(unsafe.Slice((*byte)(floatState), 40), oldFloat)
		C.dword_5d4594_2516372, *floatMax, C.dword_5d4594_2516380 = C.uint32_t(oldRange[0]), oldRange[1], C.uint32_t(oldRange[2])
		for i, p := range floatConstants {
			*p = oldFloatConstants[i]
		}
		GetServer = oldGet
	}()

	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = 0, 0
	C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = 0xdeadbeef, 0x13579bdf, 0x2468ace0
	*count, *firstSlot, *lastSlot = 0x4567, 0x89abcdef, 0x76543210
	*swapCounter, *rekeyCounter = swaps, rekeys
	for i := range unsafe.Slice((*byte)(floatState), 40) {
		unsafe.Slice((*byte)(floatState), 40)[i] = 0xa5
	}
	C.dword_5d4594_2516372, *floatMax, C.dword_5d4594_2516380 = 0x10203040, 0x50607080, 0x90a0b0c0
	*floatConstants[0], *floatConstants[1], *floatConstants[2], *floatConstants[3], *floatConstants[4] = 0x1p-32, 2111111111, 1492, 1776, 5115

	before := time.Now().Unix()
	var result uint32
	if wrapper {
		Sub_56F1C0()
	} else {
		result = uint32(C.sub_56F1C0())
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
	copy(out.FloatState[:], unsafe.Slice((*byte)(floatState), 40))
	out.FloatRange = [3]uint32{uint32(C.dword_5d4594_2516372), *floatMax, uint32(C.dword_5d4594_2516380)}
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
	actualFloat := out.FloatState
	actualRange := out.FloatRange
	for candidate := before; candidate <= after; candidate++ {
		C.sub_56FF00(C.int(candidate))
		random := uint32(C.nox_xxx_protect_56F240())
		want := PortTestProtectionInitCandidate{Seed: candidate, Random: random, Key: frame ^ random}
		copy(want.FloatState[:], unsafe.Slice((*byte)(floatState), 40))
		want.FloatRange = [3]uint32{uint32(C.dword_5d4594_2516372), *floatMax, uint32(C.dword_5d4594_2516380)}
		out.Candidates = append(out.Candidates, want)
	}
	copy(unsafe.Slice((*byte)(floatState), 40), actualFloat[:])
	C.dword_5d4594_2516372, *floatMax, C.dword_5d4594_2516380 = C.uint32_t(actualRange[0]), actualRange[1], C.uint32_t(actualRange[2])
	return out
}
