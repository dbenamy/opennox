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

type PortTestHandleOp struct {
	Reserved bool
	Value    int32
}

type PortTestHandleSnapshot struct {
	Result                 uint32
	Values                 [][2]uint32
	Sum, Key, Sequence     uint32
	Count                  uint16
	LogicIndex, OtherIndex int
	LinksValid             bool
}

func PortTestHandles(initial [][2]uint32, key, sum, sequence uint32, seed int, ops []PortTestHandleOp) []PortTestHandleSnapshot {
	oldGet := GetServer
	core := new(server.Server)
	core.Rand.Logic, core.Rand.Other = prand.New(seed), prand.New(seed+1)
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	head, tail, oldKey, oldSum, oldSequence := C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	defer func() {
		Sub_56F3B0()
		C.dword_5d4594_2516344, C.dword_5d4594_2516352 = head, tail
		C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = oldKey, oldSum, oldSequence
		*count = oldCount
		GetServer = oldGet
	}()
	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = 0, 0
	C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = C.uint(key), C.uint(sum), C.uint(sequence)
	*count = uint16(len(initial))
	var first, last *protection.Record
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
	}
	C.dword_5d4594_2516344 = C.uint(uintptr(unsafe.Pointer(first)))
	C.dword_5d4594_2516352 = C.uint(uintptr(unsafe.Pointer(last)))
	snapshot := func(result C.int) PortTestHandleSnapshot {
		out := PortTestHandleSnapshot{
			Result:     uint32(result),
			Sum:        uint32(C.dword_5d4594_2516328),
			Key:        uint32(C.dword_5d4594_2516348),
			Sequence:   uint32(C.dword_5d4594_2516356),
			Count:      *count,
			LogicIndex: core.Rand.Logic.Index(),
			OtherIndex: core.Rand.Other.Index(),
			LinksValid: true,
		}
		var prev *protection.Record
		for p := protectionHead(); p != nil; p = p.Next {
			if len(out.Values) >= len(initial)+7*len(ops) {
				out.LinksValid = false
				break
			}
			out.LinksValid = out.LinksValid && p.Prev == prev
			out.Values = append(out.Values, [2]uint32{p.ID ^ key, p.Value ^ key})
			prev = p
		}
		out.LinksValid = out.LinksValid && uint32(uintptr(unsafe.Pointer(prev))) == uint32(C.dword_5d4594_2516352)
		return out
	}
	var out []PortTestHandleSnapshot
	for _, op := range ops {
		var result C.int
		if op.Reserved {
			result = C.sub_56F250()
		} else {
			result = C.nox_xxx_protectionCreateInt_56F400(C.int(op.Value))
		}
		out = append(out, snapshot(result))
	}
	return out
}
