//go:build porttest

package legacy

/*
#include "GAME1_1.h"
extern uint32_t dword_587000_66116;
#include "GAME1_2.h"
#include "GAME4_3.h"
#include "server__ability__ability.h"
#include "server__magic__plyrgide.h"
#include "server__magic__plyrspel.h"
int nox_xxx_loadGuides_427070(void);
int nox_xxx_netAbilityReport_4D8060(int,int,int);
int nox_xxx_netReportGuideAward_4D8000(int,char,char,int);
int nox_xxx_netSendSpellAward_4D7F90(int,int,char,int);
static uint32_t nox_porttest_book_award_call(int op,uint32_t a0,uint32_t a1,uint32_t a2,uint32_t a3,uint32_t a4){switch(op){
case 0: return (uint32_t)sub_424CB0(a0);
case 1: return (uint32_t)sub_424D00();
case 2: return (uint32_t)sub_424D20(a0);
case 3: return (uint32_t)nox_xxx_abilityNameToN_424D80((char*)a0);
case 4: return (uint32_t)nox_xxx_guide_427010((char*)a0);
case 5: return (uint32_t)nox_xxx_guideNameByN_427230(a0);
case 6: return (uint32_t)nox_xxx_guiCreatureGetName_427240(a0);
case 7: return (uint32_t)nox_xxx_creatureIsCharmableByTT_4272B0(a0);
case 8: return (uint32_t)nox_xxx_guideGetDescById_4272E0(a0);
case 9: return (uint32_t)nox_xxx_bookGetFirstCreMB_427300();
case 10: return (uint32_t)nox_xxx_bookGetNextCre_427320(a0);
case 11: return (uint32_t)nox_xxx_bookGetCreatureImg_427400(a0);
case 12: return (uint32_t)sub_427430(a0);
case 13: return (uint32_t)nox_xxx_guideGetUnitSize_427460(a0);
case 14: return (uint32_t)nox_xxx_loadGuides_427070();
case 15: return (uint32_t)nox_xxx_netAbilityReport_4D8060(a0,a1,a2);
case 16: nox_xxx_abilGetSuccess_4FB960_ability(a0); return 0;
case 17: return (uint32_t)nox_xxx_abilityRewardServ_4FB9C0_ability(a0,a1,a2);
case 18: return (uint32_t)nox_xxx_netReportGuideAward_4D8000(a0,a1,a2,a3);
case 19: return (uint32_t)nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide(a0,a1,a2);
case 20: return (uint32_t)nox_xxx_netSendSpellAward_4D7F90(a0,a1,a2,a3);
case 21: nox_xxx_abilGetError_4FB0B0_magic_plyrspel(a0); return 0;
case 22: return (uint32_t)nox_xxx_spellGrantToPlayer_4FB550((nox_object_t*)a0,a1,a2,a3,a4);
case 23: return (uint32_t)sub_53F930(a0,a1);
case 24: return (uint32_t)nox_xxx_useSpellReward_53F9E0(a0,a1);
case 25: return (uint32_t)nox_xxx_useAbilityReward_53FAE0(a0,a1);
} return 0; }
*/
import "C"
import "unsafe"

var portTestBookAwardNames = []string{
	"sub_424CB0",
	"sub_424D00",
	"sub_424D20",
	"nox_xxx_abilityNameToN_424D80",
	"nox_xxx_guide_427010",
	"nox_xxx_guideNameByN_427230",
	"nox_xxx_guiCreatureGetName_427240",
	"nox_xxx_creatureIsCharmableByTT_4272B0",
	"nox_xxx_guideGetDescById_4272E0",
	"nox_xxx_bookGetFirstCreMB_427300",
	"nox_xxx_bookGetNextCre_427320",
	"nox_xxx_bookGetCreatureImg_427400",
	"sub_427430",
	"nox_xxx_guideGetUnitSize_427460",
	"nox_xxx_loadGuides_427070",
	"nox_xxx_netAbilityReport_4D8060",
	"nox_xxx_abilGetSuccess_4FB960_ability",
	"nox_xxx_abilityRewardServ_4FB9C0_ability",
	"nox_xxx_netReportGuideAward_4D8000",
	"nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide",
	"nox_xxx_netSendSpellAward_4D7F90",
	"nox_xxx_abilGetError_4FB0B0_magic_plyrspel",
	"nox_xxx_spellGrantToPlayer_4FB550",
	"sub_53F930",
	"nox_xxx_useSpellReward_53F9E0",
	"nox_xxx_useAbilityReward_53FAE0",
}

// PortTestBookAwardCall dispatches the selected original implementations. It owns
// no gameplay state and substitutes no dependencies.
func PortTestBookAwardCall(name string, args ...uint32) uint32 {
	var a [5]uint32
	copy(a[:], args)
	for i, n := range portTestBookAwardNames {
		if name == n {
			return uint32(C.nox_porttest_book_award_call(C.int(i), C.uint32_t(a[0]), C.uint32_t(a[1]), C.uint32_t(a[2]), C.uint32_t(a[3]), C.uint32_t(a[4])))
		}
	}
	panic("unknown book award operation: " + name)
}

// PortTestBookEnchantCount exposes the owned live count for signed-boundary tests.
func PortTestBookEnchantCount() *int32 { return (*int32)(unsafe.Pointer(&C.dword_587000_66116)) }
