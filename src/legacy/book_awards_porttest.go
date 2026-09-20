//go:build porttest

package legacy

/*
#include "GAME1_1.h"
#include "GAME4_3.h"
#include "server__ability__ability.h"
#include "server__magic__plyrspel.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

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

// PortTestBookAwardCall uses actual retained C exports or native private helpers.
func PortTestBookAwardCall(name string, args ...uint32) uint32 {
	var a [5]uint32
	copy(a[:], args)
	ptr := unsafe.Pointer(uintptr(a[0]))
	u := (*server.Object)(ptr)
	switch name {
	case "sub_424CB0":
		return uint32(bookEnchantCountActive(u))
	case "sub_424D00":
		return uint32(bookEnchantFirst())
	case "sub_424D20":
		return uint32(bookEnchantNext(int32(a[0])))
	case "nox_xxx_abilityNameToN_424D80":
		return uint32(bookAbilityID(alloc.GoString((*byte)(ptr))))
	case "nox_xxx_guide_427010":
		return uint32(C.nox_xxx_guide_427010((*C.char)(ptr)))
	case "nox_xxx_guideNameByN_427230":
		return uint32(uintptr(unsafe.Pointer(bookGuideName(int32(a[0])))))
	case "nox_xxx_guiCreatureGetName_427240":
		return uint32(C.nox_xxx_guiCreatureGetName_427240(C.int(a[0])))
	case "nox_xxx_creatureIsCharmableByTT_4272B0":
		return uint32(bookGuideCharmable(a[0]))
	case "nox_xxx_guideGetDescById_4272E0":
		return bookGuideDescription(int32(a[0]))
	case "nox_xxx_bookGetFirstCreMB_427300":
		return uint32(bookGuideFirst())
	case "nox_xxx_bookGetNextCre_427320":
		return uint32(bookGuideNext(int32(a[0])))
	case "nox_xxx_bookGetCreatureImg_427400":
		return bookGuideImage(int32(a[0]))
	case "sub_427430":
		return bookGuideCage(int32(a[0]))
	case "nox_xxx_guideGetUnitSize_427460":
		return uint32(bookGuideSize(int32(a[0])))
	case "nox_xxx_loadGuides_427070":
		return uint32(bookLoadGuides())
	case "nox_xxx_netAbilityReport_4D8060":
		return bookReportAbility(u, int32(a[1]), int32(a[2]))
	case "nox_xxx_abilGetSuccess_4FB960_ability":
		bookAwardClientMessage(int32(a[0]), true)
		return 0
	case "nox_xxx_abilityRewardServ_4FB9C0_ability":
		return uint32(bookAwardAbility(u, int32(a[1]), int32(a[2])))
	case "nox_xxx_netReportGuideAward_4D8000":
		return bookReportGuide(u, int32(a[1]), int32(a[2]), int32(a[3]))
	case "nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide":
		return uint32(bookAwardGuide(u, int32(a[1]), int32(a[2])))
	case "nox_xxx_netSendSpellAward_4D7F90":
		return bookReportSpell(u, int32(a[1]), int32(a[2]), int32(a[3]))
	case "nox_xxx_abilGetError_4FB0B0_magic_plyrspel":
		bookAwardClientMessage(int32(a[0]), false)
		return 0
	case "nox_xxx_spellGrantToPlayer_4FB550":
		return uint32(bookAwardSpell(u, int32(a[1]), int32(a[2]), int32(a[3]), int32(a[4])))
	case "sub_53F930":
		return uint32(C.sub_53F930(C.int(a[0]), C.int(a[1])))
	case "nox_xxx_useSpellReward_53F9E0":
		return uint32(C.nox_xxx_useSpellReward_53F9E0(C.int(a[0]), C.int(a[1])))
	case "nox_xxx_useAbilityReward_53FAE0":
		return uint32(C.nox_xxx_useAbilityReward_53FAE0(C.int(a[0]), C.int(a[1])))
	}
	panic("unknown book award operation: " + name)
}
func PortTestBookEnchantCount() *int32 { return &bookEnchantN }
