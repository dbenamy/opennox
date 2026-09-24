package legacy

import "unsafe"

func sub_4540E0(p int32) int32 {
	return int32(serverPanelsSpellApply((*uint32)(unsafe.Pointer(uintptr(uint32(p))))))
}
