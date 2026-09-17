//go:build porttest

package legacy

/*
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
int sub_4DE4D0(char a1);
extern uint32_t dword_5d4594_1563276;
extern uint32_t dword_5d4594_1599688;
extern uint32_t dword_5d4594_526276;
extern uint32_t dword_5d4594_3484;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestMatchRosterSendPlayers(to int) { C.nox_xxx_newPlayerSendAllPlayers_4DE300(C.int(to)) }
func PortTestMatchRosterSendSettings()      { C.nox_xxx_netGameSettings_4DEF00() }
func PortTestMatchRosterOwnCache() func() {
	old := C.dword_5d4594_1563276
	return func() { C.dword_5d4594_1563276 = old }
}
func PortTestMatchRosterPlayerPacket(buf []byte, pl *server.Player) {
	C.nox_xxx_netNewPlayerMakePacket_4DDA90((*C.uchar)(unsafe.Pointer(&buf[0])), (*nox_playerInfo)(pl.C()))
}
func PortTestMatchRoster(op string, u *server.Object, pl *server.Player, data unsafe.Pointer, a [4]uint32) uint32 {
	x := C.int(a[0])
	switch op {
	case "gui-settings":
		return uint32(C.nox_xxx_netGuiGameSettings_4DD9B0(C.char(a[0]), data, C.int(a[1])))
	case "inventory":
		C.sub_4DDE10(x, (*nox_playerInfo)(pl.C()))
	case "player-ids":
		return uint32(C.nox_xxx_sendAllPlayerIDs_4DE270(C.int(uintptr(pl.C()))))
	case "objective-minimap":
		return uint32(C.nox_xxx_servMinimapRevealFlag_4DE380(x))
	case "object-report-mask":
		C.sub_4DE410(x)
	case "object-resync-mask":
		return uint32(C.sub_4DE4D0(C.char(a[0])))
	case "object-clear-mask":
		return uint32(C.sub_4E80C0(C.char(a[0])))
	case "wall-open":
		return uint32(C.sub_4DF120(data))
	case "wall-close":
		return uint32(C.sub_4DF180(data))
	case "team-roster":
		C.sub_4DF2E0(x)
	case "simple-object":
		return uint32(C.nox_xxx_netSendSimpleObject2_4DF360(x, asObjectC(u)))
	case "assign-team":
		C.sub_4DF3C0((*nox_playerInfo)(pl.C()))
	case "flag-state":
		return uint32(C.sub_4E82C0(C.uchar(a[0]), C.char(a[1]), C.char(a[2]), C.short(a[3])))
	case "winner-min":
		return uint32(C.sub_5095E0())
	case "winner-max":
		return uint32(C.sub_5098A0())
	case "winner-flag":
		return uint32(C.sub_5099B0())
	case "check-victory":
		C.nox_server_checkVictory_509A60()
	case "check-limit":
		return uint32(C.sub_5096F0())
	case "remember":
		C.sub_509C30((*nox_playerInfo)(pl.C()))
	case "remembered":
		return uint32(C.sub_509D80(C.int(uintptr(pl.C()))))
	case "forget":
		return uint32(bool2int(C.sub_509CB0() != nil))
	default:
		panic(op)
	}
	return 0
}
func PortTestMatchRosterQuery(name string, class int, group uint32) uint32 {
	return uint32(C.sub_509CF0(internCStr(name), C.char(class), C.int(group)))
}
func PortTestMatchRosterFlagPointers(index byte) (unsafe.Pointer, unsafe.Pointer) {
	return unsafe.Pointer(C.sub_4E8310()), unsafe.Pointer(C.sub_4E8320(C.uchar(index)))
}
func PortTestMatchRosterGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"flag-type":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1563276)),
		"remembered-initialized": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599688)),
		"team-cap":               (*uint32)(unsafe.Pointer(&C.dword_5d4594_526276)),
		"server-subflags":        (*uint32)(unsafe.Pointer(&C.dword_5d4594_3484)),
	}
	old := map[string]uint32{}
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	header := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1599676), 12)
	saved := append([]byte(nil), header...)
	clear(header)
	return words, func() {
		C.sub_509CB0()
		copy(header, saved)
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestMatchRosterRemembered() [][]byte {
	var out [][]byte
	if C.dword_5d4594_1599688 == 0 {
		return out
	}
	for p := C.nox_common_list_getFirstSafe_425890((*C.nox_list_item_t)(memmap.PtrOff(0x5D4594, 1599676))); p != nil; p = C.nox_common_list_getNextSafe_4258A0(p) {
		out = append(out, append([]byte(nil), unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(p), 12)), 20)...))
	}
	return out
}

func PortTestMatchRosterFlagWinner() uint32 { return uint32(C.sub_5099B0()) }
