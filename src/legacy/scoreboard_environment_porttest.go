//go:build porttest

package legacy

/*
#include "defs.h"
extern unsigned int nox_player_netCode_85319C;
extern uint32_t dword_587000_145664;
extern uint32_t dword_587000_145668;
extern uint32_t dword_587000_145672;
extern uint32_t dword_5d4594_1090040;
extern uint32_t dword_5d4594_1090044;
extern nox_window* dword_5d4594_1090048;
extern nox_window* dword_5d4594_1090100;
extern uint32_t dword_5d4594_1090108;
extern uint32_t dword_5d4594_1090112;
extern uint32_t dword_5d4594_1090120;
*/
import "C"
import "unsafe"

func PortTestScoreboardWords() (map[string]*uint32, func()) {
	m := map[string]*uint32{
		"nox_player_netCode_85319C": (*uint32)(unsafe.Pointer(&C.nox_player_netCode_85319C)),
		"dword_587000_145664":       (*uint32)(unsafe.Pointer(&C.dword_587000_145664)),
		"dword_587000_145668":       (*uint32)(unsafe.Pointer(&C.dword_587000_145668)),
		"dword_587000_145672":       (*uint32)(unsafe.Pointer(&C.dword_587000_145672)),
		"dword_5d4594_1090040":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090040)),
		"dword_5d4594_1090044":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090044)),
		"dword_5d4594_1090048":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090048)),
		"dword_5d4594_1090100":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090100)),
		"dword_5d4594_1090108":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090108)),
		"dword_5d4594_1090112":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090112)),
		"dword_5d4594_1090120":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1090120)),
	}
	old := map[string]uint32{}
	for n, p := range m {
		old[n] = *p
		*p = 0
	}
	return m, func() {
		for n, p := range m {
			*p = old[n]
		}
	}
}
