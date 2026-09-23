//go:build porttest

package legacy

/*
#include "defs.h"


*/
import "C"
import "unsafe"

func PortTestScoreboardWords() (map[string]*uint32, func()) {
	m := map[string]*uint32{
		"nox_player_netCode_85319C": (*uint32)(unsafe.Pointer(&nox_player_netCode_85319C)),
		"dword_587000_145664":       (*uint32)(unsafe.Pointer(&dword_587000_145664)),
		"dword_587000_145668":       (*uint32)(unsafe.Pointer(&dword_587000_145668)),
		"dword_587000_145672":       (*uint32)(unsafe.Pointer(&dword_587000_145672)),
		"dword_5d4594_1090040":      (*uint32)(unsafe.Pointer(&dword_5d4594_1090040)),
		"dword_5d4594_1090044":      (*uint32)(unsafe.Pointer(&dword_5d4594_1090044)),
		"dword_5d4594_1090048":      (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1090048)),
		"dword_5d4594_1090100":      (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1090100)),
		"dword_5d4594_1090108":      (*uint32)(unsafe.Pointer(&dword_5d4594_1090108)),
		"dword_5d4594_1090112":      (*uint32)(unsafe.Pointer(&dword_5d4594_1090112)),
		"dword_5d4594_1090120":      (*uint32)(unsafe.Pointer(&dword_5d4594_1090120)),
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
