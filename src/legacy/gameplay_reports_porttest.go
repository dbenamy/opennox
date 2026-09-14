//go:build porttest

package legacy

/*
#include "GAME3_2.h"
#include <stdint.h>
static uint32_t gameplayReportsInvoke(int op, uint32_t* args) {
 switch(op) {
 case 0: return (uint32_t)(uintptr_t)nox_xxx_netSendInterestingId_4D7BE0((int)(uintptr_t)args[0]);
 case 1: return (uint32_t)(uintptr_t)sub_4D7E50((nox_object_t*)(uintptr_t)args[0]);
 case 2: return (uint32_t)(uintptr_t)sub_4D7EA0();
 case 3: return (uint32_t)(uintptr_t)nox_xxx_netCreatureCmd_4D7EE0((int)(uintptr_t)args[0], (char)(uintptr_t)args[1]);
 case 4: return (uint32_t)(uintptr_t)nox_xxx_netNotifyRate_4D7F10((int)(uintptr_t)args[0]);
 case 5: return (uint32_t)(uintptr_t)nox_xxx_netReportObjectPoison_4D7F40((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1], (char)(uintptr_t)args[2]);
 case 6: sub_4D81A0((int)(uintptr_t)args[0]); return 0;
 case 7: return (uint32_t)(uintptr_t)nox_xxx_netReportAnimFrame_4D81F0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 8: return (uint32_t)(uintptr_t)nox_xxx_netReportXStatus_4D8230((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 9: return (uint32_t)(uintptr_t)nox_xxx_netReportPlrStatus_4D8270((int)(uintptr_t)args[0]);
 case 10: return (uint32_t)(uintptr_t)nox_xxx_netReportCharges_4D82B0((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1], (char)(uintptr_t)args[2], (char)(uintptr_t)args[3]);
 case 11: return (uint32_t)(uintptr_t)sub_4D82F0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 12: return (uint32_t)(uintptr_t)nox_xxx_netReportDequip_4D84C0((int)(uintptr_t)args[0], (const nox_object_t*)(uintptr_t)args[1]);
 case 13: return (uint32_t)(uintptr_t)nox_xxx_netReportEquip_4D8540((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1], (int)(uintptr_t)args[2]);
 case 14: return (uint32_t)(uintptr_t)nox_xxx_netReportDequip_4D8590((int)(uintptr_t)args[0], (const nox_object_t*)(uintptr_t)args[1]);
 case 15: return (uint32_t)(uintptr_t)nox_xxx_netReportTotalHealth_4D85C0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 16: return (uint32_t)(uintptr_t)nox_xxx_netReportUnitCurrentHP_4D8620((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 17: return (uint32_t)(uintptr_t)nox_xxx_netSendTeam_4D8670((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 18: return (uint32_t)(uintptr_t)nox_xxx_netSendPlrHealthToTeam_4D86E0((int)(uintptr_t)args[0]);
 case 19: return (uint32_t)(uintptr_t)nox_xxx_netReportHealthDelta_4D8760((int)(uintptr_t)args[0], (short)(uintptr_t)args[1], (short)(uintptr_t)args[2]);
 case 20: return (uint32_t)(uintptr_t)nox_xxx_itemReportHealth_4D87A0((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1]);
 case 21: return (uint32_t)(uintptr_t)nox_xxx_netReportStamina_4D8800((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 22: return (uint32_t)(uintptr_t)sub_4D8840((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 23: return (uint32_t)(uintptr_t)sub_4D8870((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 24: return (uint32_t)(uintptr_t)nox_xxx_netReportTotalMana_4D88C0((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 25: return (uint32_t)(uintptr_t)nox_xxx_netReportMana_4D8930((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 26: return (uint32_t)(uintptr_t)nox_xxx_servSendStats_4D8990((int)(uintptr_t)args[0], (int)(uintptr_t)args[1], (char)(uintptr_t)args[2]);
 case 27: return (uint32_t)(uintptr_t)nox_xxx_netReportArmorVal_4D8A30((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 28: return (uint32_t)(uintptr_t)nox_xxx_netReportPickup_4D8A60((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1]);
 case 29: return (uint32_t)(uintptr_t)nox_xxx_netReportModifiablePickup_4D8AD0((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1]);
 case 30: return (uint32_t)(uintptr_t)nox_xxx_netReportDrop_4D8B50((int)(uintptr_t)args[0], (const nox_object_t*)(uintptr_t)args[1]);
 case 31: return (uint32_t)(uintptr_t)nox_xxx_netSendDMWinner_4D8B90((int)(uintptr_t)args[0], (char)(uintptr_t)args[1]);
 case 32: return (uint32_t)(uintptr_t)nox_xxx_netSendDMTeamWinner_4D8BF0((int)(uintptr_t)args[0], (char)(uintptr_t)args[1]);
 case 33: return (uint32_t)(uintptr_t)nox_xxx_netFlagballWinner_4D8C40((int)(uintptr_t)args[0]);
 case 34: return (uint32_t)(uintptr_t)nox_xxx_netFlagWinner_4D8C40_4D8C80((int)(uintptr_t)args[0], (char)(uintptr_t)args[1]);
 case 35: return (uint32_t)(uintptr_t)nox_xxx_scavengerHuntReport_4D8CD0((int)(uintptr_t)args[0]);
 case 36: nox_xxx_playerIncrementElimDeath_4D8D40((int)(uintptr_t)args[0]); return 0;
 case 37: return (uint32_t)(uintptr_t)nox_xxx_changeScore_4D8E90((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 38: return (uint32_t)(uintptr_t)nox_xxx_playerSubLessons_4D8EC0((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 39: return (uint32_t)(uintptr_t)nox_xxx_netReportLesson_4D8EF0((nox_object_t*)(uintptr_t)args[0]);
 case 40: return (uint32_t)(uintptr_t)nox_xxx_netTimerStatus_4D8F50((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 41: return (uint32_t)(uintptr_t)nox_xxx_netReportEnchant_4D8F90((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 42: nox_xxx_netReportObjHidden_4D8FD0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]); return 0;
 case 43: return (uint32_t)(uintptr_t)nox_xxx_netReportUnitHeight_4D9020((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1]);
 case 44: return (uint32_t)(uintptr_t)sub_4D90E0((int)(uintptr_t)args[0], (char)(uintptr_t)args[1]);
 case 45: return (uint32_t)(uintptr_t)nox_xxx_earthquakeSend_4D9110((float*)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 46: return (uint32_t)(uintptr_t)nox_xxx_netReportAcquireCreature_4D91A0((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1]);
 case 47: return (uint32_t)(uintptr_t)nox_xxx_netFxShield_0_4D9200((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 48: return (uint32_t)(uintptr_t)nox_xxx_netMonitorCreature_4D9250((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 49: return (uint32_t)(uintptr_t)nox_xxx_netSendUnMonitorCrea_4D92A0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 50: return (uint32_t)(uintptr_t)nox_xxx_netReportTeamBase_4D92D0((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 51: return (uint32_t)(uintptr_t)nox_xxx_netReportStatsSpeed_4D9360((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1], (char)(uintptr_t)args[2], (int)(uintptr_t)args[3]);
 case 52: return (uint32_t)(uintptr_t)nox_xxx_netSendReportNPC_4D93A0((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 53: return (uint32_t)(uintptr_t)nox_xxx_netSendJournalAdd_4D9440((int)(uintptr_t)args[0], (nox_playerInfo_journal*)(uintptr_t)args[1]);
 case 54: return (uint32_t)(uintptr_t)nox_xxx_netSendJournalRemove_4D94A0((int)(uintptr_t)args[0], (const char*)(uintptr_t)args[1]);
 case 55: return (uint32_t)(uintptr_t)nox_xxx_netSendJournalUpdate_4D9500((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 56: return (uint32_t)(uintptr_t)nox_xxx_netSendChapterEnd_4D9560((int)(uintptr_t)args[0], (char)(uintptr_t)args[1], (int)(uintptr_t)args[2]);
 case 57: return (uint32_t)(uintptr_t)nox_xxx_netSendFlagStatus_4D95A0((int)(uintptr_t)args[0], (char)(uintptr_t)args[1], (char)(uintptr_t)args[2], (char)(uintptr_t)args[3], (short)(uintptr_t)args[4]);
 case 58: return (uint32_t)(uintptr_t)nox_xxx_netSendBallStatus_4D95F0((int)(uintptr_t)args[0], (char)(uintptr_t)args[1], (short)(uintptr_t)args[2]);
 case 59: return (uint32_t)(uintptr_t)nox_xxx_netReportSpellStat_4D9630((int)(uintptr_t)args[0], (int)(uintptr_t)args[1], (char)(uintptr_t)args[2]);
 case 60: return (uint32_t)(uintptr_t)nox_xxx_netSendSecondaryWeapon_4D9670((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1], (char)(uintptr_t)args[2]);
 case 61: return (uint32_t)(uintptr_t)nox_xxx_netMsgLastQuiver_4D96B0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1]);
 case 62: return (uint32_t)(uintptr_t)nox_xxx_netMsgInventoryLoaded_4D96E0((int)(uintptr_t)args[0]);
 case 63: return (uint32_t)(uintptr_t)nox_xxx_netFriendAddRemove_4D97A0((int)(uintptr_t)args[0], (uint32_t*)(uintptr_t)args[1], (int)(uintptr_t)args[2]);
 case 64: return (uint32_t)(uintptr_t)sub_4D97E0((int)(uintptr_t)args[0]);
 case 65: return (uint32_t)(uintptr_t)nox_xxx_netMsgFadeBeginPlayer((int)(uintptr_t)args[0], (int)(uintptr_t)args[1], (int)(uintptr_t)args[2]);
 case 66: nox_xxx_playerReportAnything_4D9900((int)(uintptr_t)args[0]); return 0;
 case 67: return (uint32_t)(uintptr_t)sub_4D9CF0((int)(uintptr_t)args[0]);
 case 68: return (uint32_t)(uintptr_t)sub_4D9D20((int)(uintptr_t)args[0], (nox_object_t*)(uintptr_t)args[1]);
 case 69: return (uint32_t)(uintptr_t)sub_4D9D60((int)(uintptr_t)args[0], (int)(uintptr_t)args[1]);
 case 70: return (uint32_t)(uintptr_t)sub_4D9DF0((int)(uintptr_t)args[0], (int)(uintptr_t)args[1], (char)(uintptr_t)args[2]);
 case 71: return (uint32_t)(uintptr_t)sub_4D9E30((int)(uintptr_t)args[0], (int)(uintptr_t)args[1], (char)(uintptr_t)args[2]);
 case 72: return (uint32_t)(uintptr_t)nox_xxx_netGauntlet_4D9E70((int)(uintptr_t)args[0]);
 }
 return 0;
}
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
)

type PortTestGameplayReportArg struct {
	Kind        string
	Ref, Offset int
	Value       uint32
}
type PortTestGameplayReportRecord struct {
	Text  string
	Flags uint16
}
type PortTestGameplayReportsSpec struct {
	Records []PortTestGameplayReportRecord
	Args    [5]PortTestGameplayReportArg
	Calls   map[uint32][5]PortTestGameplayReportArg
	Caches  [3]uint32
}
type PortTestGameplayReportsResult struct {
	Caches   [3]uint32
	Messages [3][]byte
}

func (p *portTestShopPools) gameplayReportsPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.Reports
	if sp == nil {
		return func() {}
	}
	p.reports = sp
	p.reportRecords = nil
	for _, v := range sp.Records {
		if len(v.Text) > 63 {
			panic("journal fixture text bounds")
		}
		ptr := p.objectiveRegion(80)
		clear(unsafe.Slice((*byte)(ptr), 80))
		copy(unsafe.Slice((*byte)(ptr), 64), v.Text)
		*(*uint16)(unsafe.Add(ptr, 72)) = v.Flags
		p.reportRecords = append(p.reportRecords, ptr)
	}
	var old [3]uint32
	for i := range old {
		v := memmap.PtrUint32(0x5D4594, 1556320+uintptr(4*i))
		old[i] = *v
		*v = sp.Caches[i]
	}
	return func() {
		for i, v := range old {
			*memmap.PtrUint32(0x5D4594, 1556320+uintptr(4*i)) = v
		}
		p.reports = nil
		p.reportRecords = nil
	}
}
func (p *portTestShopPools) gameplayReportsArg(a PortTestGameplayReportArg) uint32 {
	if a.Kind == "" || a.Kind == "value" {
		return a.Value
	}
	var ptr unsafe.Pointer
	var size int
	switch a.Kind {
	case "record":
		ptr = p.temporary.world.objectives.attack.record
		if a.Ref != 0 {
			if a.Ref < 1 || a.Ref > len(p.reportRecords) {
				panic("report record reference")
			}
			ptr = p.reportRecords[a.Ref-1]
		}
		size = 80
	case "name":
		ptr = p.temporary.world.objectives.attack.controls.name
		size = 1
	default:
		u := p.temporaryRef(a.Ref)
		if u != nil {
			switch a.Kind {
			case "object":
				ptr = u.CObj()
				size = 772
			case "update":
				ptr = u.UpdateData
				size = int(unsafe.Sizeof(*u.UpdateDataPlayer()))
			case "player":
				ptr = unsafe.Pointer(u.UpdateDataPlayer().Player)
				size = int(unsafe.Sizeof(*u.UpdateDataPlayer().Player))
			case "health":
				ptr = unsafe.Pointer(u.HealthData)
				size = int(unsafe.Sizeof(*u.HealthData))
			default:
				panic("gameplay report argument kind")
			}
		}
	}
	if ptr == nil {
		if a.Offset != 0 {
			panic("offset from null report argument")
		}
		return 0
	}
	if a.Offset < 0 || a.Offset >= size {
		panic("report argument bounds")
	}
	return uint32(uintptr(unsafe.Add(ptr, a.Offset)))
}
func (p *portTestShopPools) gameplayReportsAction(a PortTestShopAction) uint32 {
	if p.reports == nil {
		panic("gameplay report fixture not prepared")
	}
	spec := p.reports.Args
	if call, ok := p.reports.Calls[a.Value]; ok {
		spec = call
	}
	var args [5]uint32
	for i, v := range spec {
		args[i] = p.gameplayReportsArg(v)
	}
	ret := uint32(C.gameplayReportsInvoke(C.int(a.Op-1800), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
	p.temporary.result = ret
	p.temporary.world.objectives.attack.controls.result = uint64(ret)
	return ret
}
func (p *portTestShopPools) gameplayReportsSnapshot() *PortTestGameplayReportsResult {
	if p.reports == nil {
		return nil
	}
	out := new(PortTestGameplayReportsResult)
	for i := range out.Caches {
		out.Caches[i] = *memmap.PtrUint32(0x5D4594, 1556320+uintptr(4*i))
	}
	for i, ind := range []ntype.PlayerInd{1, 7, 31} {
		out.Messages[i] = p.proxy.core.NetList.CopyPacketsA(ind, 1)
	}
	return out
}
