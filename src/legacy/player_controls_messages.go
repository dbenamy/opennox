package legacy

/*
#include "GAME1_1.h"
#include "GAME1.h"
#include "GAME1_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
static void* controlInversionAddress(void) { return nox_xxx_inversionEffect_4E03D0; }
*/
import "C"
import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

func controlRewardNotify(u *server.Object, kind int32, target *server.Object, selector byte) int32 {
	if u == nil || u.ObjClass&4 == 0 {
		return int32(controlRaw(u))
	}
	if kind < 0 || kind > 2 {
		return kind - 2
	}
	msg := [5]byte{0xf0, byte(30 + kind), selector}
	binary.LittleEndian.PutUint16(msg[3:], uint16(target.NetCode))
	return int32(C.nox_xxx_netSendPacket0_4E5420(C.int(*controlByte(controlPlayer(u), 2064)), unsafe.Pointer(&msg[0]), 5, 0, 1))
}
func controlLockedDoor(u *server.Object, key *C.char, selector byte) {
	if u == nil || u.ObjClass&4 == 0 || key == nil {
		return
	}
	s := C.GoString(key)
	if len(s) == 0 || len(s) > 48 {
		return
	}
	// The fixed wire format includes unused bytes; keep their padding zero.
	msg := [52]byte{0xf0, 33}
	copy(msg[2:], s)
	msg[51] = selector
	C.nox_xxx_netSendPacket0_4E5420(C.int(*controlByte(controlPlayer(u), 2064)), unsafe.Pointer(&msg[0]), 52, 0, 1)
}
func controlRespawnNotify(u *server.Object, flag byte) int32 {
	msg := [9]byte{0xe9}
	binary.LittleEndian.PutUint16(msg[1:], uint16(u.NetCode))
	binary.LittleEndian.PutUint32(msg[3:], GetServer().S().Frame())
	msg[7] = byte(controlRespawnFlags())
	msg[8] = flag
	return int32(C.nox_xxx_netSendPacket1_4E5390(255, C.int(uintptr(unsafe.Pointer(&msg[0]))), 9, 0, 0))
}
func controlGuideLevel(u, target *server.Object) int32 {
	if u == nil || target == nil || u.ObjClass&4 == 0 || target.ObjClass&2 == 0 {
		return 0
	}
	name := C.nox_xxx_getUnitName_4E39D0((*C.nox_object_t)(target.CObj()))
	id := C.nox_xxx_guide_427010(name)
	if id == 0 {
		return 0
	}
	return int32(*equipmentWord(controlPlayer(u), 4244+4*int(id)))
}
func controlGuideDamage(u, target *server.Object, damage *int32) int32 {
	level := controlGuideLevel(u, target)
	if level == 0 {
		return 0
	}
	v := float32(float64(C.nox_xxx_gamedataGetFloat_419D40(internCStr("FieldGuideDamageBonus")))*float64(*damage) + 0.5)
	*damage = floatToInt32(v)
	return *damage
}
func controlInversion(u, target *server.Object) int32 {
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjFlags&0x100 == 0 || it.ObjClass&0x13001000 == 0 {
			continue
		}
		for i := 2; i < 4; i++ {
			mod := *controlPtr(it.InitData, i*4)
			if mod == nil {
				continue
			}
			fn := *controlPtr(mod, 88)
			if fn == nil || fn != C.controlInversionAddress() {
				continue
			}
			var result int32
			owner := target.FindOwnerChainPlayer()
			ccall.CallVoidPtr6(fn, mod, it.CObj(), u.CObj(), target.CObj(), owner.CObj(), unsafe.Pointer(&result))
			if result == 1 {
				return 1
			}
		}
	}
	return 0
}
func controlScheduledSpell(u, target *server.Object, back bool) int32 {
	d := u.UpdateData
	count := controlByte(d, 212)
	if *count == 0 {
		return 0
	}
	off := 192
	if back {
		off = 188 + 4*int(*count)
	}
	id := *equipmentWord(d, off)
	reason := C.nox_xxx_checkPlrCantCastSpell_4FD150((*C.nox_object_t)(u.CObj()), C.int(id), 0)
	args := [3]uint32{controlRaw(target), math.Float32bits(float32(int32(*equipmentWord(d, 220)))), math.Float32bits(float32(int32(*equipmentWord(d, 224))))}
	if reason != 0 {
		C.nox_xxx_netInformTextMsg_4DA0F0(C.int(*controlByte(controlPlayer(u), 2064)), 0, &reason)
		C.nox_xxx_aud_501960(231, (*C.nox_object_t)(u.CObj()), 0, 0)
	} else {
		C.nox_xxx_castSpellByUser_4FDD20(C.int(id), (*C.nox_object_t)(u.CObj()), unsafe.Pointer(&args[0]))
	}
	// Front removal shifts and clears; back removal only decrements the count.
	if !back {
		for i := 1; i < int(*count); i++ {
			*equipmentWord(d, 188+4*i) = *equipmentWord(d, 192+4*i)
		}
		*equipmentWord(d, 188+4*int(*count)) = 0
	}
	*count--
	return 1
}
