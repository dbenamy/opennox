//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_1064856;
extern uint32_t dword_5d4594_1064860;
extern uint32_t dword_5d4594_1064864;
extern uint32_t dword_5d4594_1064868;
extern uint32_t dword_5d4594_1096636;
extern void* dword_5d4594_1096640;
extern uint32_t dword_5d4594_1123520;
extern void* dword_5d4594_1123524;
extern uint32_t dword_5d4594_1193712;
extern uint32_t dword_5d4594_1303452;
extern uint32_t dword_5d4594_1305680;
extern uint32_t dword_5d4594_1305684;
extern uint32_t dword_5d4594_1319056;
extern uint32_t dword_5d4594_1319060;
extern uint32_t dword_5d4594_1321216;
extern uint32_t dword_5d4594_805820;
extern uint32_t dword_5d4594_811904;
extern uint32_t dword_5d4594_825736;
extern uint32_t dword_5d4594_825744;
extern uint32_t dword_8531A0_2576;
extern void* nox_client_spriteUnderCursorXxx_1096644;
extern uint32_t nox_xxx_useAudio_587000_80772;
*/
import "C"
import "unsafe"

// PortTestClientInteractionWords owns the real C globals, separately from mapped blobs.
func PortTestClientInteractionWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1064856":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1064856)),
		"dword_5d4594_1064860":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1064860)),
		"dword_5d4594_1064864":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1064864)),
		"dword_5d4594_1064868":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1064868)),
		"dword_5d4594_1096636":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1096636)),
		"dword_5d4594_1096640":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1096640)),
		"dword_5d4594_1123520":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1123520)),
		"dword_5d4594_1123524":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1123524)),
		"dword_5d4594_1193712":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1193712)),
		"dword_5d4594_1303452":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1303452)),
		"dword_5d4594_1305680":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305680)),
		"dword_5d4594_1305684":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305684)),
		"dword_5d4594_1319056":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1319056)),
		"dword_5d4594_1319060":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1319060)),
		"dword_5d4594_1321216":                    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1321216)),
		"dword_5d4594_805820":                     (*uint32)(unsafe.Pointer(&C.dword_5d4594_805820)),
		"dword_5d4594_811904":                     (*uint32)(unsafe.Pointer(&C.dword_5d4594_811904)),
		"dword_5d4594_825736":                     (*uint32)(unsafe.Pointer(&C.dword_5d4594_825736)),
		"dword_5d4594_825744":                     (*uint32)(unsafe.Pointer(&C.dword_5d4594_825744)),
		"dword_8531A0_2576":                       (*uint32)(unsafe.Pointer(&C.dword_8531A0_2576)),
		"nox_client_spriteUnderCursorXxx_1096644": (*uint32)(unsafe.Pointer(&C.nox_client_spriteUnderCursorXxx_1096644)),
		"nox_xxx_useAudio_587000_80772":           (*uint32)(unsafe.Pointer(&C.nox_xxx_useAudio_587000_80772)),
	}
	old := make(map[string]uint32)
	for n, p := range words {
		old[n] = *p
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}
