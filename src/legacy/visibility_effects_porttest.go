//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME4_2.h"
extern uint32_t dword_5d4594_1565512, dword_5d4594_1565516, dword_5d4594_1565520, dword_5d4594_2649712;
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Retained entries exercise the C ABI; retired bridges call production Go directly.
// Operation 13 belonged to the proven orphaned throttle helper and is retired.
func PortTestVisibilityEffects(op int, a, b *server.Object, pos *types.Pointf, words *[4]int32, args [5]int32, data []byte, command string) uint32 {

	fp := (*C.float)(unsafe.Pointer(pos))
	switch op {
	case 0:
		return uint32(C.nox_xxx_netSendPointFx_522FF0(C.char(args[0]), (*C.float2)(unsafe.Pointer(pos))))
	case 1:
		return uint32(visibilityFXSend(*pos, data))
	case 2:
		return uint32(visibilityFXPointExtra(byte(args[0]), byte(args[1]), *pos))
	case 3:
		return uint32(visibilityFXSpark(*pos, byte(args[0])))
	case 4:
		visibilityFXGeneratorBreak(*pos, byte(args[0]))
	case 5:
		return uint32(visibilityFXVampire(byte(args[0]), *words, uint16(args[1])))
	case 6:
		return uint32(visibilityFXPrediction(a))
	case 7:
		return uint32(visibilityFXShield(a, pos))
	case 8:
		return uint32(visibilityFXSummonStart(uint16(args[0]), *pos, byte(args[1]), uint16(args[2]), uint16(args[3])))
	case 9:
		return uint32(visibilityFXSummonCancel(uint16(args[0])))
	case 10:
		visibilityFXGeneratorSpawn(*words, uint16(args[0]))
	case 11:
		C.nox_xxx_sendArrowTrapFX_5238A0(fp, C.char(args[0]))
	case 12:
		return uint32(visibilitySpecialUpdate(a, b))
	case 14:
		return uint32(visibilityKillable(a))
	case 15:
		return uint32(C.nox_xxx_frameCounterSetCopyToNextFrame_5281D0())
	case 16:
		return visibilityFrameCopy(false)
	case 17:
		visibilityUpdateSight(a)
	case 18:
		return uint32(visibilityLost(a, int(args[0])))
	case 19:
		visibilitySelectTarget(a)
	case 20:
		visibilityCandidate(a, b)
	case 21:
		visibilitySee(a, b)
	case 22:
		return uint32(visibilityRemove(a, b))
	case 23:
		return uint32(bool2int(visibilityContains(a, b)))
	case 24:
		return uint32(visibilityGlobalRemove(a))
	case 25:
		visibilityDestroyReport(a)
	case 26:
		return uint32(visibilityOutOfSight(int(args[0]), b))
	case 27:
		return uint32(visibilityInShadows(int(args[0]), b))
	case 28:
		return uint32(visibilityMonsterCommand(a, b, command, uint16(args[0])))
	case 29:
		return uint32(visibilityClearChats())
	default:
		panic("unknown visibility/effect operation")
	}
	return 0
}

// Own the real reliable-message pool, using the same layout as lifecycle/shop
// fixtures. Captures expose defined packet fields, never allocation addresses.
func PortTestVisibilityReliableOwner() (reset func(), snapshot func() []PortTestShopPacketResult, free func()) {
	pool := memmap.PtrPtr(0x5D4594, 1565508)
	oldPool := *pool
	head, tail, size, mask := C.dword_5d4594_1565512, C.dword_5d4594_1565516, C.dword_5d4594_1565520, C.dword_5d4594_2649712
	seq := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565524)), 64)
	oldSeq := bytes.Clone(seq)
	*pool = nil
	C.dword_5d4594_1565512 = 0
	C.dword_5d4594_1565516 = 0
	C.dword_5d4594_2649712 = 0x80000082
	reset = func() {
		if *pool != nil {
			alloc.AsClass(*pool).Free()
			*pool = nil
		}
		C.dword_5d4594_1565512 = 0
		C.dword_5d4594_1565516 = 0
		clear(seq)
	}
	snapshot = func() (out []PortTestShopPacketResult) {
		for q := uint32(C.dword_5d4594_1565512); q != 0; {
			b := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(q))), 416)
			p := PortTestShopPacketResult{Recipient: b[250], Ordered: b[184], A4: binary.LittleEndian.Uint32(b[404:]), A5: binary.LittleEndian.Uint32(b[180:]), Data: bytes.Clone(b[251 : 251+int(b[401])])}
			for i := range p.Sequence {
				p.Sequence[i] = binary.LittleEndian.Uint16(b[186+2*i:])
			}
			out = append(out, p)
			q = binary.LittleEndian.Uint32(b[408:])
		}
		return
	}
	free = func() {
		reset()
		*pool = oldPool
		C.dword_5d4594_1565512, C.dword_5d4594_1565516, C.dword_5d4594_1565520, C.dword_5d4594_2649712 = head, tail, size, mask
		copy(seq, oldSeq)
	}
	return
}
