//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2649712;
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"runtime"
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
	Rules         *PortTestGameplayReportRules
	RecipientMask *uint32
	Records       []PortTestGameplayReportRecord
	Args          [5]PortTestGameplayReportArg
	Calls         map[uint32][5]PortTestGameplayReportArg
	Caches        [3]uint32
}
type PortTestGameplayReportsResult struct {
	Rules    *PortTestGameplayReportRulesResult `json:",omitempty"`
	Caches   [3]uint32
	Messages [3][]byte
}

func (p *portTestShopPools) gameplayReportsPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.Reports
	if sp == nil {
		return func() {}
	}
	p.reports = sp
	restoreRules := p.gameplayReportRulesPrepare(sp.Rules)
	oldMask := C.dword_5d4594_2649712
	if sp.RecipientMask != nil {
		C.dword_5d4594_2649712 = C.uint32_t(*sp.RecipientMask)
	}
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
		restoreRules()
		C.dword_5d4594_2649712 = oldMask
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
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	control := PortTestProtectionFloatCW()
	if control&0x0f00 != 0x0200 {
		panic("report x87 precision/rounding")
	}
	p.gameplayReportRulesMembers()

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
	ret := gameplayReportsInvoke(a.Op-1800, args)
	if PortTestProtectionFloatCW() != control {
		panic("report changed x87 control")
	}
	p.temporary.result = ret
	p.temporary.world.objectives.attack.controls.result = uint64(ret)
	return ret
}
func (p *portTestShopPools) gameplayReportsSnapshot() *PortTestGameplayReportsResult {
	if p.reports == nil {
		return nil
	}
	out := new(PortTestGameplayReportsResult)
	if st := p.reportRules; st != nil {
		out.Rules = &PortTestGameplayReportRulesResult{Notified: *memmap.PtrUint32(0x5D4594, 3536), Countdown: append([]PortTestGameplayReportCountdown(nil), st.countdown...)}
		for i := range out.Rules.PlayerStatus {
			out.Rules.PlayerStatus[i] = *(*uint32)(unsafe.Add(unsafe.Pointer(p.proxy.life.players[i].UpdateDataPlayer().Player), 3680))
		}
	}
	for i := range out.Caches {
		out.Caches[i] = *memmap.PtrUint32(0x5D4594, 1556320+uintptr(4*i))
	}
	for i, ind := range []ntype.PlayerInd{1, 7, 31} {
		out.Messages[i] = p.proxy.core.NetList.CopyPacketsA(ind, 1)
	}
	return out
}

func gameplayReportsInvoke(op int, args [5]uint32) uint32 {
	switch op {
	case 0:
		return uint32(gameplayReportInterestingID((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 1:
		return uint32(gameplayReportReset((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 2:
		return uint32(gameplayReportResetAll())
	case 3:
		return uint32(gameplayReportCreature(int(args[0]), byte(args[1])))
	case 4:
		return uint32(gameplayReportRate(int(args[0])))
	case 5:
		return uint32(gameplayReportPoison((*server.Object)(unsafe.Pointer(uintptr(args[0]))), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2])))
	case 6:
		return uint32(gameplayReportExperience((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 7:
		return uint32(gameplayReportAnimation(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 8:
		return uint32(gameplayReportXStatus(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 9:
		return uint32(gameplayReportPlayerStatus((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 10:
		return uint32(gameplayReportCharges(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2]), byte(args[3])))
	case 11:
		return uint32(gameplayReportEquipment(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 12:
		return uint32(gameplayReportDequipFlags(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 13:
		return uint32(gameplayReportEquip(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), int(args[2])))
	case 14:
		return uint32(gameplayReportDequipItem(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 15:
		return uint32(gameplayReportTotalHealth(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 16:
		return uint32(gameplayReportCurrentHealth(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 17:
		return uint32(gameplayReportTeam(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 18:
		return uint32(gameplayReportPlayerHealthToTeam(int(args[0])))
	case 19:
		return uint32(gameplayReportHealthDelta(int(args[0]), uint16(args[1]), int16(args[2])))
	case 20:
		return uint32(gameplayReportItemHealth(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 21:
		return uint32(gameplayReportStamina(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 22:
		return uint32(gameplayReportObjectByte(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 23:
		return uint32(gameplayReportPlayerStat(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 24:
		return uint32(gameplayReportTotalMana(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 25:
		return uint32(gameplayReportMana(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 26:
		return uint32(gameplayReportStats(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2])))
	case 27:
		return uint32(gameplayReportArmor(int(args[0]), uint32(args[1])))
	case 28:
		return uint32(gameplayReportPickup(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 29:
		return uint32(gameplayReportModifiablePickup(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 30:
		return uint32(gameplayReportDrop(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 31:
		return uint32(gameplayReportDMWinner((*server.Object)(unsafe.Pointer(uintptr(args[0]))), byte(args[1])))
	case 32:
		return uint32(gameplayReportDMTeamWinner((*server.Team)(unsafe.Pointer(uintptr(args[0]))), byte(args[1])))
	case 33:
		return uint32(gameplayReportFlagballWinner((*server.Team)(unsafe.Pointer(uintptr(args[0])))))
	case 34:
		return uint32(gameplayReportFlagWinner((*server.Team)(unsafe.Pointer(uintptr(args[0]))), byte(args[1])))
	case 35:
		return uint32(gameplayReportScavenger((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 36:
		return uint32(gameplayReportEliminationDeath((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 37:
		return uint32(gameplayReportChangeScore((*server.Object)(unsafe.Pointer(uintptr(args[0]))), int(args[1])))
	case 38:
		return uint32(gameplayReportSubtractLessons((*server.Object)(unsafe.Pointer(uintptr(args[0]))), int(args[1])))
	case 39:
		return uint32(gameplayReportLesson((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 40:
		return uint32(gameplayReportTimer(int(args[0]), uint32(args[1])))
	case 41:
		return uint32(gameplayReportEnchant(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 42:
		return uint32(gameplayReportHidden(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 43:
		return uint32(gameplayReportHeight(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 44:
		return uint32(gameplayReportEarthquakeByte(int(args[0]), byte(args[1])))
	case 45:
		return uint32(gameplayReportEarthquake((*types.Pointf)(unsafe.Pointer(uintptr(args[0]))), int(args[1])))
	case 46:
		return uint32(gameplayReportAcquireCreature(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 47:
		return uint32(gameplayReportShield(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 48:
		return uint32(gameplayReportMonitor(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 49:
		return uint32(gameplayReportUnmonitor(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 50:
		return uint32(gameplayReportTeamBase(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 51:
		return uint32(gameplayReportSpeed(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2]), uint32(args[3])))
	case 52:
		return uint32(gameplayReportNPC(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 53:
		return uint32(gameplayReportJournal(int(args[0]), unsafe.Pointer(uintptr(args[1])), 1))
	case 54:
		return uint32(gameplayReportJournal(int(args[0]), unsafe.Pointer(uintptr(args[1])), 2))
	case 55:
		return uint32(gameplayReportJournal(int(args[0]), unsafe.Pointer(uintptr(args[1])), 3))
	case 56:
		return uint32(gameplayReportChapter(int(args[0]), byte(args[1]), int(args[2])))
	case 57:
		return uint32(gameplayReportFlag(int(args[0]), byte(args[1]), byte(args[2]), byte(args[3]), uint16(args[4])))
	case 58:
		return uint32(gameplayReportBall(int(args[0]), byte(args[1]), uint16(args[2])))
	case 59:
		return uint32(gameplayReportSpellStat(int(args[0]), uint32(args[1]), byte(args[2])))
	case 60:
		return uint32(gameplayReportSecondary(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2])))
	case 61:
		return uint32(gameplayReportQuiver(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 62:
		return uint32(gameplayReportInventoryLoaded(int(args[0])))
	case 63:
		return uint32(gameplayReportFriend(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), int(args[2])))
	case 64:
		return uint32(gameplayReportFriendReset(int(args[0])))
	case 65:
		return uint32(gameplayReportFade(int(args[0]), int(args[1]), int(args[2])))
	case 66:
		return uint32(gameplayReportAnything((*server.Object)(unsafe.Pointer(uintptr(args[0])))))
	case 67:
		return uint32(gameplayReportQuestStart(int(args[0])))
	case 68:
		return uint32(gameplayReportQuestObject(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 69:
		return uint32(gameplayReportQuestLevel(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1])))))
	case 70:
		return uint32(gameplayReportSilverKey(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2])))
	case 71:
		return uint32(gameplayReportGoldKey(int(args[0]), (*server.Object)(unsafe.Pointer(uintptr(args[1]))), byte(args[2])))
	case 72:
		return uint32(gameplayReportGauntlet(int(args[0])))
	}
	panic("unknown gameplay reporting operation")
}
