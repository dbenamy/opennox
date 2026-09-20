//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME1_1.h"
#include "GAME2_3.h"
extern uint32_t dword_5d4594_1200768, dword_5d4594_1200832;
extern uint32_t dword_5d4594_1197352, dword_5d4594_1197356;
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestClientSessionWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{"settingsChanged": (*uint32)(&C.dword_5d4594_1200768), "settingsNotice": (*uint32)(&C.dword_5d4594_1200832)}
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

// Own the existing C queue without substituting its insertion or lookup rules.
func PortTestClientSequenceOwner() (func(), func() PortTestClientSequenceState, func()) {
	raw := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1197340), 24)
	saved := bytes.Clone(raw)
	ready, pending := C.dword_5d4594_1197352, C.dword_5d4594_1197356
	C.sub_48D740()
	reset := func() { C.sub_48D760(); C.sub_48D740() }
	snapshot := func() PortTestClientSequenceState {
		r := PortTestClientSequenceState{Current: memmap.Uint16(0x5D4594, 1197360)}
		head := (*C.nox_list_item_t)(memmap.PtrOff(0x5D4594, 1197340))
		for p := C.nox_common_list_getFirstSafe_425890(head); p != nil; p = C.nox_common_list_getNextSafe_4258A0(p) {
			if len(r.Nodes) > 65536 {
				panic("sequence list cycle")
			}
			b := unsafe.Slice((*byte)(unsafe.Pointer(p)), 32)
			seq := binary.LittleEndian.Uint32(b[8:])
			size := binary.LittleEndian.Uint16(b[24:])
			r.Nodes = append(r.Nodes, PortTestClientSequenceNode{seq, binary.LittleEndian.Uint64(b[16:]), bytes.Clone(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(p), 32)), int(size)))})
			if uint32(uintptr(unsafe.Pointer(p))) == uint32(C.dword_5d4594_1197352) {
				r.Ready = seq + 1
			}
			if uint32(uintptr(unsafe.Pointer(p))) == uint32(C.dword_5d4594_1197356) {
				r.Pending = seq + 1
			}
		}
		if C.dword_5d4594_1197352 != 0 && r.Ready == 0 || C.dword_5d4594_1197356 != 0 && r.Pending == 0 {
			panic("sequence cursor outside owned list")
		}
		return r
	}
	return reset, snapshot, func() { reset(); copy(raw, saved); C.dword_5d4594_1197352, C.dword_5d4594_1197356 = ready, pending }
}
