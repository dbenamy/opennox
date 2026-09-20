package legacy

/*
#include <stdlib.h>
#include <stdint.h>
extern uint32_t dword_5d4594_1197352, dword_5d4594_1197356;
*/
import "C"

import (
	"encoding/binary"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

func clientSequenceHead() *legacyListNode {
	return (*legacyListNode)(memmap.PtrOff(0x5D4594, 1197340))
}
func clientSequenceOld(current, sequence uint16) bool {
	if current < 10000 && sequence >= 65535-(10000-current) {
		return true
	}
	distance := int(current) - int(sequence)
	return distance > 0 && distance <= 10000
}
func clientSequenceContains(sequence uint16) bool {
	for p := listNext(clientSequenceHead()); p != nil; p = listNext(p) {
		if p.tag == uintptr(sequence) {
			return true
		}
	}
	return false
}
func clientSequenceEnqueue(data []byte) {
	seq := binary.LittleEndian.Uint16(data[1:])
	if clientSequenceOld(memmap.Uint16(0x5D4594, 1197360), seq) || clientSequenceContains(seq) {
		return
	}
	size := int(data[3])
	// Retain the foreign queue allocation/layout while shared globals remain ABI words.
	ptr := C.calloc(C.size_t(32+size), 1)
	if ptr == nil {
		return
	}
	node := listInit((*legacyListNode)(ptr))
	node.tag = uintptr(seq)
	*(*uint16)(unsafe.Add(ptr, 24)) = uint16(size)
	*(*uint64)(unsafe.Add(ptr, 16)) = uint64(uint32(PlatformTicks()))
	copy(unsafe.Slice((*byte)(unsafe.Add(ptr, 32)), size), data[4:4+size])
	if memmap.Uint16(0x5D4594, 1197360) == seq {
		C.dword_5d4594_1197352 = C.uint32_t(uintptr(ptr))
	}
	listAscending(clientSequenceHead(), node)
}
func clientSequencePoll() {
	ready := (*legacyListNode)(unsafe.Pointer(uintptr(C.dword_5d4594_1197352)))
	if ready == nil {
		pending := (*legacyListNode)(unsafe.Pointer(uintptr(C.dword_5d4594_1197356)))
		if pending != nil && uint64(uint32(PlatformTicks()))-*(*uint64)(unsafe.Add(unsafe.Pointer(pending), 16)) > 30000 {
			*memmap.PtrUint16(0x5D4594, 1197360) = uint16(pending.tag)
			ready = pending
			C.dword_5d4594_1197352 = C.uint32_t(uintptr(unsafe.Pointer(pending)))
			C.dword_5d4594_1197356 = C.uint32_t(uintptr(unsafe.Pointer(listNext(pending))))
		}
	}
	p := ready
	for p != nil && uint32(p.tag) == uint32(memmap.Uint16(0x5D4594, 1197360)) {
		next := listNext(p)
		if next == nil {
			next = listNext(clientSequenceHead())
			if next == p {
				next = nil
			}
		}
		size := int(*(*uint16)(unsafe.Add(unsafe.Pointer(p), 24)))
		Nox_xxx_netOnPacketRecvCli_48EA70(31, (*byte)(unsafe.Add(unsafe.Pointer(p), 32)), size)
		*memmap.PtrUint16(0x5D4594, 1197360)++
		listRemove(p)
		C.free(unsafe.Pointer(p))
		p = next
	}
	// Original polling drops the pending cursor when a gap is polled before its
	// timeout. Preserve this observable behavior separately from queue membership.
	C.dword_5d4594_1197356 = C.uint32_t(uintptr(unsafe.Pointer(p)))
	C.dword_5d4594_1197352 = 0
}
func clientSequenceInit() {
	listClear(clientSequenceHead())
	C.dword_5d4594_1197352 = 0
	C.dword_5d4594_1197356 = 0
	*memmap.PtrUint16(0x5D4594, 1197360) = 0
}
func clientSequenceFree() {
	for p := listNext(clientSequenceHead()); p != nil; {
		next := listNext(p)
		listRemove(p)
		C.free(unsafe.Pointer(p))
		p = next
	}
	listClear(clientSequenceHead())
	*memmap.PtrUint16(0x5D4594, 1197360) = 0
}
