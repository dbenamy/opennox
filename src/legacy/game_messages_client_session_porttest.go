//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_3.h"
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestClientSessionWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{"settingsChanged": (*uint32)(&dword_5d4594_1200768), "settingsNotice": (*uint32)(&dword_5d4594_1200832)}
	old := map[string]uint32{}
	for n, p := range words {
		old[n] = *p
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

type PortTestClientSequenceNode struct {
	Sequence uint32
	Ticks    uint64
	Data     []byte
}
type PortTestClientSequenceState struct {
	Current        uint16
	Ready, Pending uint32
	Nodes          []PortTestClientSequenceNode
}

// Own the live queue without substituting its insertion or lookup rules.
func PortTestClientSequenceOwner() (func(), func() PortTestClientSequenceState, func()) {
	raw := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1197340), 24)
	saved := bytes.Clone(raw)
	ready, pending := dword_5d4594_1197352, dword_5d4594_1197356
	clientSequenceInit()
	reset := func() { clientSequenceFree(); clientSequenceInit() }
	snapshot := func() PortTestClientSequenceState {
		r := PortTestClientSequenceState{Current: memmap.Uint16(0x5D4594, 1197360)}
		head := (*legacyListNode)(memmap.PtrOff(0x5D4594, 1197340))
		for p := listNext(head); p != nil; p = listNext(p) {
			if len(r.Nodes) > 65536 {
				panic("sequence list cycle")
			}
			b := unsafe.Slice((*byte)(unsafe.Pointer(p)), 32)
			seq := binary.LittleEndian.Uint32(b[8:])
			size := binary.LittleEndian.Uint16(b[24:])
			r.Nodes = append(r.Nodes, PortTestClientSequenceNode{seq, binary.LittleEndian.Uint64(b[16:]), bytes.Clone(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(p), 32)), int(size)))})
			if uint32(uintptr(unsafe.Pointer(p))) == uint32(dword_5d4594_1197352) {
				r.Ready = seq + 1
			}
			if uint32(uintptr(unsafe.Pointer(p))) == uint32(dword_5d4594_1197356) {
				r.Pending = seq + 1
			}
		}
		if dword_5d4594_1197352 != 0 && r.Ready == 0 || dword_5d4594_1197356 != 0 && r.Pending == 0 {
			panic("sequence cursor outside owned list")
		}
		return r
	}
	return reset, snapshot, func() { reset(); copy(raw, saved); dword_5d4594_1197352, dword_5d4594_1197356 = ready, pending }
}
