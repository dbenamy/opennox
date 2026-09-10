//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
extern uint32_t dword_5d4594_2516356;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/server"
)

type PortTestRekeySnapshot struct {
	Result, ExpectedRandom                           uint32
	Values                                           [][2]uint32
	Sum, Key, Sequence                               uint32
	SwapCount, RekeyCount                            uint32
	Count                                            uint16
	LogicIndex, OtherIndex                           int
	FloatState, ExpectedFloatState, BeforeFloatState [40]byte
	FloatRange, ExpectedFloatRange, BeforeFloatRange [3]uint32
	LinksValid, NodesAndLinksStable                  bool
}

func PortTestRekey(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, wrapper bool) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		if wrapper {
			Nox_xxx_protectData_56F5C0()
			return 0
		}
		return uint32(C.nox_xxx_protectData_56F5C0())
	})
}

func portTestRekeyOperation(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, run func() uint32) PortTestRekeySnapshot {
	oldGet := GetServer
	core := new(server.Server)
	core.SetFrame(frame)
	core.Rand.Logic, core.Rand.Other = prand.New(seed), prand.New(seed+1)
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	head, tail, oldKey, oldSum, oldSequence := C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	swaps, rekeys := memmap.PtrUint32(0x5D4594, 2516360), memmap.PtrUint32(0x5D4594, 2516364)
	oldSwaps, oldRekeys := *swaps, *rekeys
	oldRandom := protectionRandom
	defer func() {
		Sub_56F3B0()
		C.dword_5d4594_2516344, C.dword_5d4594_2516352 = head, tail
		C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = oldKey, oldSum, oldSequence
		*count, *swaps, *rekeys = oldCount, oldSwaps, oldRekeys
		protectionRandom = oldRandom
		GetServer = oldGet
	}()
	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = 0, 0
	C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = C.uint(key), C.uint(sum), C.uint(sequence)
	*count, *swaps, *rekeys = uint16(len(initial)), swapCount, rekeyCount
	var first, last *protection.Record
	nodes := make([]*protection.Record, 0, len(initial))
	for _, value := range initial {
		r := (*protection.Record)(C.calloc(1, C.size_t(unsafe.Sizeof(protection.Record{}))))
		if r == nil {
			panic("fixture allocation failed")
		}
		*r = protection.Record{ID: value[0] ^ key, Value: value[1] ^ key, Prev: last}
		if last == nil {
			first = r
		} else {
			last.Next = r
		}
		last = r
		nodes = append(nodes, r)
	}
	C.dword_5d4594_2516344 = C.uint(uintptr(unsafe.Pointer(first)))
	C.dword_5d4594_2516352 = C.uint(uintptr(unsafe.Pointer(last)))
	protectionRandom.Seed(floatSeed)
	beforeRandom := protectionRandom
	beforeFloat, beforeRange := portTestRandomState(protectionRandom)
	expectedRandom := protectionRandom.Draw()
	expectedFloat, expectedRange := portTestRandomState(protectionRandom)
	protectionRandom = beforeRandom
	result := run()
	out := PortTestRekeySnapshot{
		Result:             result,
		ExpectedRandom:     expectedRandom,
		Sum:                uint32(C.dword_5d4594_2516328),
		Key:                uint32(C.dword_5d4594_2516348),
		Sequence:           uint32(C.dword_5d4594_2516356),
		SwapCount:          *swaps,
		RekeyCount:         *rekeys,
		Count:              *count,
		LogicIndex:         core.Rand.Logic.Index(),
		OtherIndex:         core.Rand.Other.Index(),
		ExpectedFloatState: expectedFloat,
		ExpectedFloatRange: expectedRange,
		BeforeFloatRange:   beforeRange,
		LinksValid:         true,
	}
	out.BeforeFloatState = beforeFloat
	out.FloatState, out.FloatRange = portTestRandomState(protectionRandom)
	p := protectionHead()
	for i, want := range nodes {
		var prev, next *protection.Record
		if i != 0 {
			prev = nodes[i-1]
		}
		if i+1 < len(nodes) {
			next = nodes[i+1]
		}
		if p != want || p.Prev != prev || p.Next != next {
			out.LinksValid = false
			break
		}
		out.Values = append(out.Values, [2]uint32{p.ID ^ uint32(C.dword_5d4594_2516348), p.Value ^ uint32(C.dword_5d4594_2516348)})
		p = p.Next
	}
	out.LinksValid = out.LinksValid && p == nil && len(out.Values) == len(initial) && protectionHead() == first && uint32(uintptr(unsafe.Pointer(last))) == uint32(C.dword_5d4594_2516352)
	out.NodesAndLinksStable = out.LinksValid
	return out
}
