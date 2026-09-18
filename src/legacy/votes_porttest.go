//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME4.h"
#include "GAME2_3.h"
extern void* nox_alloc_vote_1599652;
extern uint32_t dword_5d4594_1599656;
extern uint32_t nox_server_resetQuestMinVotes_229988;
extern uint32_t nox_server_kickQuestPlayerMinVotes_229992;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// Own the actual production allocation class and list, restoring prior state.
func PortTestVoteOwner() func() {
	pool, head := C.nox_alloc_vote_1599652, C.dword_5d4594_1599656
	C.nox_alloc_vote_1599652 = nil
	C.dword_5d4594_1599656 = 0
	a, b := memmap.PtrUint32(0x587000, 229980), memmap.PtrUint32(0x587000, 229984)
	av, bv := *a, *b
	*a, *b = 5, 9
	if C.nox_xxx_allocVoteArray_5066D0() != 1 {
		panic("vote allocation")
	}
	return func() {
		C.sub_506720()
		C.nox_alloc_vote_1599652, C.dword_5d4594_1599656 = pool, head
		*a, *b = av, bv
	}
}

// Direct entry dispatch, without reproducing the selected algorithms.
func PortTestVoteCreate(kind int, player unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.sub_506A20(C.int(kind), C.int(uintptr(player))))
}
func PortTestVoteDelete(p unsafe.Pointer)        { C.sub_5067B0(C.int(uintptr(p))) }
func PortTestVoteRemovePlayer(p unsafe.Pointer)  { C.sub_506740((*C.nox_object_t)(p)) }
func PortTestVoteThreshold(p unsafe.Pointer) int { return int(C.sub_507000(C.int(uintptr(p)))) }
func PortTestVoteNameMessage(name *uint16, withdraw bool) {
	if withdraw {
		C.nox_xxx_voteSend_48D260((*C.wchar2_t)(unsafe.Pointer(name)))
	} else {
		C.nox_xxx_netSendRenameMb_48D2D0((*C.wchar2_t)(unsafe.Pointer(name)))
	}
}
func PortTestVoteList() []unsafe.Pointer {
	var out []unsafe.Pointer
	var previous unsafe.Pointer
	for p := unsafe.Pointer(uintptr(C.dword_5d4594_1599656)); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 44)) {
		if len(out) >= 128 || *(*unsafe.Pointer)(unsafe.Add(p, 48)) != previous {
			panic("vote list links")
		}
		out = append(out, p)
		previous = p
	}
	return out
}

func PortTestVoteCast(kind int, player unsafe.Pointer, name *uint16, withdraw bool) {
	if withdraw {
		C.sub_506C90(C.int(kind), C.int(uintptr(player)), (*C.wchar2_t)(unsafe.Pointer(name)))
	} else {
		C.sub_506870(C.int(kind), C.int(uintptr(player)), (*C.wchar2_t)(unsafe.Pointer(name)))
	}
}

func PortTestVoteSettings() ([2]*uint32, func()) {
	p := [2]*uint32{(*uint32)(&C.nox_server_resetQuestMinVotes_229988), (*uint32)(&C.nox_server_kickQuestPlayerMinVotes_229992)}
	a, b := *p[0], *p[1]
	return p, func() { *p[0], *p[1] = a, b }
}
