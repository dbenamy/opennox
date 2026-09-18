//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Own the actual production allocation class and list, restoring prior state.
func PortTestVoteOwner() func() {
	pool, head := votePool, voteHead
	votePool.Class = nil
	voteHead = nil
	a, b := memmap.PtrUint32(0x587000, 229980), memmap.PtrUint32(0x587000, 229984)
	av, bv := *a, *b
	*a, *b = 5, 9
	if voteInit() != 1 {
		panic("vote allocation")
	}
	return func() {
		voteClose()
		votePool, voteHead = pool, head
		*a, *b = av, bv
	}
}

// Direct entry dispatch, without reproducing the selected algorithms.
func PortTestVoteCreate(kind int, player unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(voteCreate(uint32(kind), (*server.Object)(player)))
}
func PortTestVoteDelete(p unsafe.Pointer)        { voteDelete((*voteRecord)(p)) }
func PortTestVoteRemovePlayer(p unsafe.Pointer)  { voteRemovePlayer((*server.Object)(p)) }
func PortTestVoteThreshold(p unsafe.Pointer) int { return voteThreshold((*voteRecord)(p)) }
func PortTestVoteNameMessage(name *uint16, withdraw bool) {
	if withdraw {
		voteSendName(name, true)
	} else {
		voteSendName(name, false)
	}
}
func PortTestVoteList() []unsafe.Pointer {
	var out []unsafe.Pointer
	var previous unsafe.Pointer
	for p := unsafe.Pointer(voteHead); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 44)) {
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
		voteWithdraw(uint32(kind), (*server.Object)(player), name)
	} else {
		voteCast(uint32(kind), (*server.Object)(player), name)
	}
}

func PortTestVoteSettings() ([2]*uint32, func()) {
	p := [2]*uint32{&voteQuestResetSetting, &voteQuestKickSetting}
	a, b := *p[0], *p[1]
	return p, func() { *p[0], *p[1] = a, b }
}
