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

// Primitive dispatch only; the baseline exercises the actual production functions.
func PortTestVisibilityEffects(op int, a, b *server.Object, pos *types.Pointf, words *[4]int32, args [5]int32, data []byte, command string) uint32 {
	ai, bi := C.int(uintptr(unsafe.Pointer(a))), C.int(uintptr(unsafe.Pointer(b)))
	ap, bp := (*C.nox_object_t)(unsafe.Pointer(a)), (*C.nox_object_t)(unsafe.Pointer(b))
	fp := (*C.float)(unsafe.Pointer(pos))
	switch op {
	case 0:
		return uint32(C.nox_xxx_netSendPointFx_522FF0(C.char(args[0]), (*C.float2)(unsafe.Pointer(pos))))
	case 1:
		return uint32(C.nox_xxx_netSendFxAllCli_523030((*C.float2)(unsafe.Pointer(pos)), unsafe.Pointer(unsafe.SliceData(data)), C.int(len(data))))
	case 2:
		return uint32(C.sub_523150(C.char(args[0]), C.char(args[1]), fp))
	case 3:
		return uint32(C.nox_xxx_netSparkExplosionFx_5231B0(fp, C.char(args[0])))
	case 4:
		C.nox_xxx_sendGeneratorBreakFX_523200(fp, C.char(args[0]))
	case 5:
		return uint32(C.nox_xxx_netSendVampFx_523270(C.char(args[0]), (*C.short)(unsafe.Pointer(words)), C.short(args[1])))
	case 6:
		return uint32(C.nox_xxx_netClientPredictLinear_523530(ai))
	case 7:
		return uint32(C.nox_xxx_netSendShieldFx_523670(ai, fp))
	case 8:
		return uint32(C.nox_xxx_sendSummonStartFX_5236F0(C.short(args[0]), fp, C.char(args[1]), C.short(args[2]), C.short(args[3])))
	case 9:
		return uint32(C.nox_xxx_sendSummonCancelFX_523760(C.short(args[0])))
	case 10:
		C.nox_xxx_sendGeneratorSpawnFX_523830((*C.int4)(unsafe.Pointer(words)), C.short(args[0]))
	case 11:
		C.nox_xxx_sendArrowTrapFX_5238A0(fp, C.char(args[0]))
	case 12:
		return uint32(C.nox_xxx_netUpdateObjectSpecial_527E50(ap, bp))
	case 13:
		return uint32(int32(C.sub_528030(ai)))
	case 14:
		return uint32(C.nox_xxx_checkIsKillable_528190(ap))
	case 15:
		return uint32(C.nox_xxx_frameCounterSetCopyToNextFrame_5281D0())
	case 16:
		return uint32(C.nox_xxx_frameCounterSetCopy_5281E0())
	case 17:
		C.nox_xxx_unitUpdateSightMB_5281F0(ap)
	case 18:
		return uint32(C.nox_xxx_aiLostSight_528560(ai, C.int(args[0])))
	case 19:
		C.sub_528610(ai)
	case 20:
		C.nox_xxx_monsterUpdateSeenEnemies_5286D0(ai, bi)
	case 21:
		C.nox_xxx_monsterVisionSeeEnemy_5287B0(ai, bi)
	case 22:
		return uint32(C.sub_528910(ai, bi))
	case 23:
		return uint32(C.sub_528950(ai, bi))
	case 24:
		return uint32(C.sub_528990(ap))
	case 25:
		C.nox_xxx_netReportDestroyObject_5289D0(ap)
	case 26:
		return uint32(C.nox_xxx_netObjectOutOfSight_528A60(C.int(args[0]), (*C.uint32_t)(unsafe.Pointer(b))))
	case 27:
		return uint32(C.nox_xxx_netObjectInShadows_528A90(C.int(args[0]), (*C.uint32_t)(unsafe.Pointer(b))))
	case 28:
		text := append([]byte(command), 0)
		return uint32(C.nox_xxx_monsterCmdSend_528BD0(ai, bi, (*C.char)(unsafe.Pointer(&text[0])), C.short(args[0])))
	case 29:
		return uint32(C.nox_xxx_destroyEveryChatMB_528D60())
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
