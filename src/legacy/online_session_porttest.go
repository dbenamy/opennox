//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
extern uint32_t dword_5d4594_10984;
extern unsigned int dword_5d4594_527988;
extern uint32_t dword_5d4594_528252, dword_5d4594_528256;
extern uint32_t dword_5d4594_528260, dword_5d4594_528264;
*/
import "C"
import "unsafe"

func PortTestOnlineSessionWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"failed":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_10984)),
		"status":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_527988)),
		"pending":  (*uint32)(unsafe.Pointer(&C.dword_5d4594_528252)),
		"active":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_528256)),
		"deadline": (*uint32)(unsafe.Pointer(&C.dword_5d4594_528260)),
		"origin":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_528264)),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestOnlineSessionCall(op string, a, b uint32) uint32 {
	switch op {
	case "reset":
		return uint32(C.sub_41E370())
	case "start":
		C.nox_xxx_reconStart_41E400()
		return 0
	case "attempt":
		return uint32(C.nox_xxx_reconAttempt_41E390())
	case "retry":
		return uint32(C.sub_41E470())
	case "marker":
		return uint32(C.sub_41E4B0(C.int(a)))
	case "briefing":
		return uint32(C.sub_41D1A0(C.int(a)))
	case "map":
		return uint32(C.sub_41D650())
	case "channels":
		return uint32(C.sub_41EC30())
	case "users":
		return uint32(C.sub_41F4B0())
	case "append":
		return uint32(C.sub_41DA70(C.int(a), C.short(b)))
	case "clear":
		return uint32(C.sub_41DA10(C.int(a)))
	case "status":
		return uint32(C.sub_41E2F0())
	case "login":
		return uint32(C.nox_xxx_officialStringCmp_41FDE0())
	case "player":
		return uint32(C.sub_41D670(nil))
	case "players":
		return uint32(C.sub_41D6C0())
	}
	panic(op)
}
