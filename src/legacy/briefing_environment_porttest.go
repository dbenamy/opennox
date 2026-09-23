//go:build porttest

package legacy

/*
#include "defs.h"
extern uint32_t dword_5d4594_831236;
extern void* dword_5d4594_805984;
extern nox_window* nox_win_unk1;
extern nox_screenParticle* nox_screenParticles_head;
extern nox_screenParticle* dword_5d4594_806052;
extern void* nox_alloc_screenParticles_806044;
*/
import "C"
import "unsafe"

func PortTestBriefingWords() (map[string]*uint32, func()) {
	m := map[string]*uint32{
		"dword_5d4594_826028":              (*uint32)(unsafe.Pointer(&sessionMOTDRoot)),
		"dword_5d4594_826032":              (*uint32)(unsafe.Pointer(&sessionMOTDList)),
		"dword_5d4594_805984":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_805984)),
		"nox_win_unk1":                     (*uint32)(unsafe.Pointer(&C.nox_win_unk1)),
		"nox_screenParticles_head":         (*uint32)(unsafe.Pointer(&C.nox_screenParticles_head)),
		"dword_5d4594_806052":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_806052)),
		"nox_alloc_screenParticles_806044": (*uint32)(unsafe.Pointer(&C.nox_alloc_screenParticles_806044)),
		"dword_5d4594_1046864":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046864)),
		"dword_5d4594_1046868":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046868)),
		"dword_5d4594_1046872":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046872)),
		"nox_xxx_aNox_cfg_0_587000_132132": (*uint32)(unsafe.Pointer(&nox_xxx_aNox_cfg_0_587000_132132)),
		"dword_5d4594_1046936":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046936)),
		"dword_5d4594_1046952":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046952)),

		"dword_587000_122956":                   (*uint32)(&dword_587000_122956),
		"nox_xxx_aSpellphoneme_3_587000_123008": (*uint32)(&nox_xxx_aSpellphoneme_3_587000_123008),
		"dword_5d4594_832480":                   (*uint32)(&dword_5d4594_832480),
		"dword_5d4594_832520":                   (*uint32)(&dword_5d4594_832520),
		"dword_5d4594_832500":                   (*uint32)(&dword_5d4594_832500),
		"dword_5d4594_832528":                   (*uint32)(&dword_5d4594_832528),
		"dword_5d4594_832524":                   (*uint32)(&dword_5d4594_832524),
		"dword_5d4594_832512":                   (*uint32)(&dword_5d4594_832512),
		"dword_5d4594_832496":                   (*uint32)(&dword_5d4594_832496),
		"dword_5d4594_832516":                   (*uint32)(&dword_5d4594_832516),
		"dword_5d4594_832508":                   (*uint32)(&dword_5d4594_832508),
		"dword_5d4594_832504":                   (*uint32)(&dword_5d4594_832504),
		"dword_5d4594_832492":                   (*uint32)(&dword_5d4594_832492),
		"dword_5d4594_832532":                   (*uint32)(&dword_5d4594_832532),
		"dword_5d4594_832536":                   (*uint32)(&dword_5d4594_832536),
		"nox_wnd_briefing_831232":               (*uint32)(&nox_wnd_briefing_831232),
		"dword_5d4594_832476":                   (*uint32)(&dword_5d4594_832476),
		"dword_5d4594_832484":                   (*uint32)(&dword_5d4594_832484),
		"dword_5d4594_831236":                   (*uint32)(&C.dword_5d4594_831236),
		"dword_5d4594_831220":                   (*uint32)(&dword_5d4594_831220),
		"dword_5d4594_831224":                   (*uint32)(&dword_5d4594_831224),
		"dword_5d4594_831240":                   (*uint32)(&dword_5d4594_831240),
		"dword_5d4594_831244":                   (*uint32)(&dword_5d4594_831244),
		"dword_5d4594_831256":                   (*uint32)(&dword_5d4594_831256),
		"dword_5d4594_831260":                   (*uint32)(&dword_5d4594_831260),
		"dword_5d4594_831276":                   (*uint32)(&dword_5d4594_831276),
	}
	old := make(map[string]uint32, len(m))
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
