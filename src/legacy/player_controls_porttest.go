//go:build porttest

package legacy

/*
#include <string.h>
#include <stdlib.h>
#include <stdint.h>
extern void* nox_monsterBin_head_2386924;
extern uint32_t dword_5d4594_1565616,dword_5d4594_1568868;
#include "GAME3_3.h"
#include "GAME4.h"
static void* controlsPlayerUpdatePtr(void) { return nox_xxx_updatePlayer_4F8100; }
void nox_xxx_playerLeaveObserver_0_4E6AA0(nox_playerInfo* pl);
void nox_xxx_unitRemoveChild_4EC470(nox_object_t* a1);
void nox_xxx_unitTransferSlaves_4EC4B0(nox_object_t* a1p);
int nox_xxx_plrReadVals_4EEDC0(nox_object_t* a1p, int a2);
int sub_4EF140(int a1);
double nox_xxx_calcBoltDamage_4EF1E0(int a1, int a2);
void sub_4EF410(int a1, unsigned char a2);
int sub_4EF6F0(int a1);
nox_object_t* nox_xxx_playerRespawnItem_4EF750(nox_object_t* a1p, char* a2, int* a3, int a4, int a5);
char nox_xxx_playerMakeDefItems_4EF7D0(int a1, int a2, int a3);
char nox_xxx_unitInitPlayer_4EFE80(nox_object_t* a1p);
int sub_4EFF10(int a1);
int nox_xxx_equipedItemByCode_4F7920(int a1, int a2);
void nox_xxx_playerSetCustomWP_4F79A0(int a1, int a2, int a3);
void nox_xxx_mapFindPlayerStart_4F7AB0(float2* a1, nox_object_t* a2p);
int nox_xxx_weaponGetStaminaByType_4F7E80(int a1);
short nox_xxx_playerRespawn_4F7EF0(nox_object_t* a1p);
int sub_4FA280(int a1);
int nox_common_mapPlrActionToStateId_4FA2B0(nox_object_t* a1p);
int nox_xxx_checkInversionEffect_4FA4F0(int a1, int a2);
char nox_xxx_mobMorphFromPlayer_4FAAC0(uint32_t* a1);
char nox_xxx_mobMorphToPlayer_4FAAF0(uint32_t* a1);
int nox_xxx_updatePlayerMonsterBot_4FAB20(uint32_t* a1);
int nox_xxx_netSendRewardNotify_4FAD50(int a1, int a2, int a3, char a4);
void sub_4FADD0(int a1, char* a2, char a3);
int sub_4FB050(int a1, int a2, int* a3);
static uint32_t controlsInitLog[256];static int controlsInitN;
static void controlsInit(nox_object_t* u,int x){if(controlsInitN+2>=256)abort();controlsInitLog[controlsInitN++]=(uint32_t)u;controlsInitLog[controlsInitN++]=x;}
static void controlsInitReset(){controlsInitN=0;}
static int controlsInitCount(){return controlsInitN;}
static uint32_t controlsInitValue(int i){return controlsInitLog[i];}
static void* controlsInitPtr(){return controlsInit;}
static void* controlsFunction(int id){switch(id){
case 5:return (void*)nox_xxx_playerLeaveObserver_0_4E6AA0;
case 8:return (void*)nox_xxx_unitRemoveChild_4EC470;
case 9:return (void*)nox_xxx_unitTransferSlaves_4EC4B0;
case 11:return (void*)nox_xxx_plrReadVals_4EEDC0;
case 12:return (void*)sub_4EF140;
case 13:return (void*)nox_xxx_calcBoltDamage_4EF1E0;
case 14:return (void*)sub_4EF410;
case 16:return (void*)sub_4EF6F0;
case 17:return (void*)nox_xxx_playerRespawnItem_4EF750;
case 18:return (void*)nox_xxx_playerMakeDefItems_4EF7D0;
case 20:return (void*)nox_xxx_unitInitPlayer_4EFE80;
case 21:return (void*)sub_4EFF10;
case 22:return (void*)nox_xxx_equipedItemByCode_4F7920;
case 24:return (void*)nox_xxx_playerSetCustomWP_4F79A0;
case 26:return (void*)nox_xxx_mapFindPlayerStart_4F7AB0;
case 31:return (void*)nox_xxx_weaponGetStaminaByType_4F7E80;
case 32:return (void*)nox_xxx_playerRespawn_4F7EF0;
case 41:return (void*)sub_4FA280;
case 42:return (void*)nox_common_mapPlrActionToStateId_4FA2B0;
case 43:return (void*)nox_xxx_checkInversionEffect_4FA4F0;
case 45:return (void*)nox_xxx_mobMorphFromPlayer_4FAAC0;
case 46:return (void*)nox_xxx_mobMorphToPlayer_4FAAF0;
case 47:return (void*)nox_xxx_updatePlayerMonsterBot_4FAB20;
case 50:return (void*)nox_xxx_netSendRewardNotify_4FAD50;
case 51:return (void*)sub_4FADD0;
case 53:return (void*)sub_4FB050;
default:return 0;}}
static uint64_t controlsCall(int id,nox_object_t* u,nox_object_t* t,int x,int y,void* record,char* name){
 switch(id){
case 5:{nox_xxx_playerLeaveObserver_0_4E6AA0((u?*(nox_playerInfo**)(*(char**)((char*)u+748)+276):0));return 0;}
case 8:{nox_xxx_unitRemoveChild_4EC470(u);return 0;}
case 9:{nox_xxx_unitTransferSlaves_4EC4B0(u);return 0;}
case 11:{return (uint32_t)nox_xxx_plrReadVals_4EEDC0(u,x);}
case 12:{return (uint32_t)sub_4EF140((int)u);}
case 13:{double d=nox_xxx_calcBoltDamage_4EF1E0(x,(int)record);uint64_t bits;memcpy(&bits,&d,8);return bits;}
case 14:{sub_4EF410((int)u,(unsigned char)x);return 0;}
case 16:{return (uint32_t)sub_4EF6F0((int)u);}
case 17:{return (uint32_t)nox_xxx_playerRespawnItem_4EF750(u,name,(int*)record,x,y);}
case 18:{return (uint32_t)nox_xxx_playerMakeDefItems_4EF7D0((int)u,x,y);}
case 20:{return (uint32_t)nox_xxx_unitInitPlayer_4EFE80(u);}
case 21:{return (uint32_t)sub_4EFF10((int)u);}
case 22:{return (uint32_t)nox_xxx_equipedItemByCode_4F7920((int)u,x);}
case 24:{nox_xxx_playerSetCustomWP_4F79A0((int)u,x,y);return 0;}
case 26:{nox_xxx_mapFindPlayerStart_4F7AB0((float2*)record,u);return 0;}
case 31:{return (uint32_t)nox_xxx_weaponGetStaminaByType_4F7E80(x);}
case 32:{return (uint32_t)nox_xxx_playerRespawn_4F7EF0(u);}
case 41:{return (uint32_t)sub_4FA280(x);}
case 42:{return (uint32_t)nox_common_mapPlrActionToStateId_4FA2B0(u);}
case 43:{return (uint32_t)nox_xxx_checkInversionEffect_4FA4F0((int)u,(int)t);}
case 45:{return (uint32_t)nox_xxx_mobMorphFromPlayer_4FAAC0((uint32_t*)u);}
case 46:{return (uint32_t)nox_xxx_mobMorphToPlayer_4FAAF0((uint32_t*)u);}
case 47:{return (uint32_t)nox_xxx_updatePlayerMonsterBot_4FAB20((uint32_t*)u);}
case 50:{return (uint32_t)nox_xxx_netSendRewardNotify_4FAD50((int)u,x,(int)t,(char)y);}
case 51:{sub_4FADD0((int)u,name,(char)x);return 0;}
case 53:{return (uint32_t)sub_4FB050((int)u,(int)t,(int*)record);}
default:return 0;}}
*/
import "C"

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestPlayerControlsSpec struct {
	SpellLifecycle        *PortTestSpellLifecycleSpec
	ByteReturn            int // 1: player record address; 2: last created item address, checked before normalization.
	Corpse                bool
	Guide                 bool
	Disallowed            uint8
	NullRecord, Modifiers bool
	Equipment             bool
	Stats                 []server.ClassStats
	MonsterRefs           []int
	UpdateByRef           map[int]map[int]uint32
	Target                int
	X, Y                  int32
	Name                  *string
	Bot                   bool
	BotWords              map[int]uint32
}
type portTestPlayerControls struct {
	spellLifecycle *portTestSpellLifecycle
	classes        map[*server.Object][2]uint32
	transitions    []uint32
	freshBots      []unsafe.Pointer
	result         uint64
	name           unsafe.Pointer
	canonical      map[*server.Object]unsafe.Pointer
}

var controlsOffsets = []uintptr{1564960, 1565600, 1568264, 1568268, 1568872}
var controlsTypeNames = []string{"Glyph", "PlayerStart", "PlayerWaypoint", "NPC", "ArcherBolt", "Bat", "PortControlsCorpse"}

func (p *portTestShopPools) controlsPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls
	if sp == nil {
		return func() {}
	}
	st := &portTestPlayerControls{canonical: make(map[*server.Object]unsafe.Pointer), classes: make(map[*server.Object][2]uint32)}
	p.temporary.world.objectives.attack.controls = st
	old := make([]uint32, len(controlsOffsets))
	for i, off := range controlsOffsets {
		v := memmap.PtrUint32(0x5d4594, off)
		old[i] = *v
		*v = 0
	}
	ball, start := C.dword_5d4594_1565616, C.dword_5d4594_1568868
	C.dword_5d4594_1565616 = 0
	C.dword_5d4594_1568868 = 0
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 215824)), 108)
	C.controlsInitReset()
	restoreTypes, restoreMods := func() {}, func() {}
	if sp.Equipment {
		restoreTypes = p.proxy.core.PortTestControlsTypes(C.controlsInitPtr(), sp.Disallowed)
		var mods []*server.ModifierEff
		var names []*byte
		for _, name := range []string{"UserColor1", "ArmorQuality1", "Material1", "Replenishment1"} {
			mods = append(mods, (*server.ModifierEff)(p.objectiveRegion(144)))
			names = append(names, (*byte)(p.objectiveString(name)))
		}
		restoreMods = p.proxy.core.PortTestControlsModifiers(mods, names)
	}
	oldPlace := Nox_xxx_inventoryServPlace_4F36F0
	Nox_xxx_inventoryServPlace_4F36F0 = func(u, t *server.Object, a, b int) bool { p.controlsAdopt(t); return oldPlace(u, t, a, b) }
	restoreCorpse := func() {}
	if sp.Corpse {
		cache := memmap.PtrUint32(0x5d4594, 2488736)
		oldCache := *cache
		*cache = 1
		table := unsafe.Slice(memmap.PtrUint32(0x5d4594, 2488740), 99)
		oldTable := append([]uint32(nil), table...)
		for i := range table {
			table[i] = uint32(p.proxy.core.Types.IndByID("PortControlsCorpse"))
		}
		points := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 280376)), 792)
		oldPoints := bytes.Clone(points)
		copy(points, blobdata.PortTestPlayerCorpsePoints())
		restoreCorpse = func() { *cache = oldCache; copy(table, oldTable); copy(points, oldPoints) }
	}
	restoreGuide := func() {}
	if sp.Guide {
		table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 70500)), 164)
		old := bytes.Clone(table)
		for i := 0; i < 41; i++ {
			name := fmt.Sprintf("controls-guide-%d", i)
			if i == 1 {
				name = "bat"
			}
			*memmap.PtrPtr(0x587000, 70500+uintptr(4*i)) = p.objectiveString(name)
		}
		restoreGuide = func() { copy(table, old) }
	}
	oldMonster := C.nox_monsterBin_head_2386924
	def := (*server.MonsterDef)(p.objectiveRegion(int(unsafe.Sizeof(server.MonsterDef{}))))
	*def = *p.proxy.combat.actor.UpdateDataMonster().MonsterDef
	def.Next244 = nil
	def.TypeInd240 = uint32(p.proxy.core.Types.IndByID("NPC"))
	clear(def.Name0[:])
	copy(def.Name0[:], "NPC")
	C.nox_monsterBin_head_2386924 = unsafe.Pointer(def)
	oldStats := p.proxy.core.Players.Stats
	if len(sp.Stats) == 4 {
		p.proxy.core.Players.Stats.Base = sp.Stats[0]
		p.proxy.core.Players.Stats.Warrior = sp.Stats[1]
		p.proxy.core.Players.Stats.Wizard = sp.Stats[2]
		p.proxy.core.Players.Stats.Conjurer = sp.Stats[3]
	}
	abilityTable := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 206108)), 40)
	oldAbility := bytes.Clone(abilityTable)
	copy(abilityTable, blobdata.PortTestPlayerAbilityTable())
	weight := unsafe.Slice((*byte)(memmap.PtrOff(0x581450, 10216)), 8)
	oldWeight := bytes.Clone(weight)
	copy(weight, blobdata.PortTestPlayerWeight())
	saved := bytes.Clone(table)
	copy(table, blobdata.PortTestPlayerActionTable())
	if sp.Name != nil {
		st.name = p.objectiveString(*sp.Name)
	}
	restoreSpellLife := p.spellLifePrepare()
	return func() {
		restoreSpellLife()
		restoreCorpse()
		restoreGuide()
		C.nox_monsterBin_head_2386924 = oldMonster
		Nox_xxx_inventoryServPlace_4F36F0 = oldPlace
		restoreMods()
		restoreTypes()
		for _, ptr := range st.freshBots {
			C.free(ptr)
		}
		p.proxy.core.Players.Stats = oldStats
		copy(abilityTable, oldAbility)
		copy(weight, oldWeight)
		for u, ptr := range st.canonical {
			u.UpdateData = ptr
			if values, ok := st.classes[u]; ok {
				*equipmentWord(u.CObj(), 8) = values[0]
				*equipmentWord(u.CObj(), 12) = values[1]
			}
		}
		for i, off := range controlsOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		C.dword_5d4594_1565616, C.dword_5d4594_1568868 = ball, start
		copy(table, saved)
	}
}
func (p *portTestShopPools) controlsItems() {
	a := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	if a.Controls == nil {
		return
	}
	st := p.temporary.world.objectives.attack.controls
	p.identify(C.controlsInitPtr(), 91600)
	p.identify(C.controlsPlayerUpdatePtr(), 91601)
	for i := 0; i < 56; i++ {
		if fn := C.controlsFunction(C.int(i)); fn != nil {
			p.identify(fn, 91000+uint32(i))
		}
	}
	for i := range p.proxy.life.players {
		u := &p.proxy.life.players[i]
		st.canonical[u] = u.UpdateData
		st.classes[u] = [2]uint32{uint32(u.ObjClass), uint32(u.ObjSubClass)}
	}
	for _, id := range a.Controls.MonsterRefs {
		u := p.temporaryRef(id)
		st.canonical[u] = u.UpdateData
		st.classes[u] = [2]uint32{uint32(u.ObjClass), uint32(u.ObjSubClass)}
		u.UpdateData = p.objectiveRegion(2200)
		for off, v := range a.Controls.UpdateByRef[id] {
			if off < 0 || off+4 > 2200 || off%4 != 0 {
				panic("controls monster word")
			}
			*(*uint32)(unsafe.Add(u.UpdateData, off)) = v
		}
	}

	if a.Controls.Modifiers {
		for i := 0; i < 4; i++ {
			*(*unsafe.Pointer)(unsafe.Add(p.temporary.world.objectives.attack.record, 4*i)) = p.proxy.core.Modif.Nox_xxx_modifGetDescById413330(i + 1).C()
		}
	}
	if a.Controls.Bot {
		u := p.temporaryRef(a.Actor)
		bot := p.objectiveRegion(2200)
		*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 292)) = bot
		*(*unsafe.Pointer)(unsafe.Add(bot, 2180)) = u.UpdateData
		for off, v := range a.Controls.BotWords {
			if off < 0 || off+4 > 2200 || off%4 != 0 {
				panic("controls bot word")
			}
			*(*uint32)(unsafe.Add(bot, off)) = v
		}
	}
	p.spellLifeItems()
}
func (p *portTestShopPools) controlsAction(a PortTestShopAction) uint32 {
	spec := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	st := p.temporary.world.objectives.attack.controls
	sp := spec.Controls
	record := p.temporary.world.objectives.attack.record
	if sp.NullRecord {
		record = nil
	}
	if a.Op == 1456 {
		u := p.temporaryRef(spec.Actor)
		st.transitions = append(st.transitions, uint32(C.controlsCall(45, asObjectC(u), nil, 0, 0, record, nil)))
		for _, v := range unsafe.Slice((*uint32)(u.CObj()), 193) {
			st.transitions = append(st.transitions, p.normalize(v))
		}
		st.result = uint64(C.controlsCall(46, asObjectC(u), nil, 0, 0, record, nil))
	} else {
		st.result = controlsInvoke(a.Op-1400, p.temporaryRef(spec.Actor), p.temporaryRef(sp.Target), sp.X, sp.Y, record, st.name)
	}

	if a.Op == 1444 {
		u := p.temporaryRef(spec.Actor)
		ptr := *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 292))
		if !sp.Bot && ptr != nil && len(st.freshBots) == 0 {
			p.identify(ptr, 91500)
			st.freshBots = append(st.freshBots, ptr)
		}
	}
	if sp.ByteReturn != 0 {
		var ptr unsafe.Pointer
		switch sp.ByteReturn {
		case 1:
			ptr = unsafe.Pointer(p.temporaryRef(spec.Actor).UpdateDataPlayer().Player)
		case 2:
			if len(p.proxy.life.created) == 0 {
				panic("controls missing pointer result")
			}
			ptr = p.proxy.life.created[len(p.proxy.life.created)-1].CObj()
		default:
			panic("controls return kind")
		}
		if uint32(st.result) != uint32(int32(int8(uintptr(ptr)))) {
			panic("controls byte pointer return")
		}
		st.result = uint64(p.normalize(uint32(uintptr(ptr))))
	}
	if a.Op == 1418 || a.Op == 1420 || a.Op == 1432 || a.Op == 1449 {
		for _, u := range p.proxy.life.created {
			if u.ObjClass&0x3001000 != 0 && *equipmentWord(u.InitData, 16) != 0 {
				panic("controls default modifier fifth word")
			}
		}
	}
	p.temporary.result = uint32(st.result)
	return p.temporary.result
}
func (p *portTestShopPools) controlsSnapshot(out []uint32) []uint32 {
	st := p.temporary.world.objectives.attack.controls
	if st == nil {
		return out
	}
	out = append(out, p.normalize(uint32(st.result)), uint32(st.result>>32), uint32(C.dword_5d4594_1565616), uint32(C.dword_5d4594_1568868))
	out = append(out, uint32(len(st.transitions)))
	out = append(out, st.transitions...)
	out = append(out, uint32(C.controlsInitCount()))
	for i := 0; i < int(C.controlsInitCount()); i++ {
		out = append(out, p.normalize(uint32(C.controlsInitValue(C.int(i)))))
	}
	out = append(out, uint32(len(st.freshBots)))
	for _, ptr := range st.freshBots {
		for _, v := range unsafe.Slice((*uint32)(ptr), 550) {
			out = append(out, p.normalize(v))
		}
	}
	for _, off := range controlsOffsets {
		out = append(out, *memmap.PtrUint32(0x5d4594, off))
	}
	return p.spellLifeSnapshot(out)
}

func (p *portTestShopPools) controlsAdopt(u *server.Object) {
	if u == nil {
		return
	}
	for _, old := range p.proxy.life.created {
		if old == u {
			return
		}
	}
	for _, item := range p.items {
		if item.u == u {
			return
		}
	}
	p.proxy.life.created = append(p.proxy.life.created, u)
	n := uint32(len(p.proxy.life.created))
	p.identify(u.CObj(), 1000+n)
	if u.InitData != nil {
		p.identify(u.InitData, 4000+n)
	}
	if u.UseData.Ptr != nil {
		p.identify(u.UseData.Ptr, 5000+n)
	}
	if u.UpdateData != nil {
		p.identify(u.UpdateData, 3000+n)
	}
}

// controlsInvoke calls retired internal entry points directly; the C dispatcher
// covers only ABIs still used by production C. Expected captures are unchanged.
func controlsInvoke(op int, u, t *server.Object, x, y int32, record, name unsafe.Pointer) uint64 {
	switch op {
	case 0:
		return uint64(uint32(controlRemoveGlyphs(u)))
	case 1:
		controlRemoveCreatures(u)
		return 0
	case 2:
		return uint64(uint32(controlRaw(controlNextObserver(controlPlayer(u)))))
	case 3:
		return uint64(uint32(controlRaw(controlFindBall())))
	case 4:
		return uint64(uint32(controlRaw(controlObserverSlave(controlPlayer(u)))))
	case 6:
		return uint64(uint32(controlRaw(controlSlave(u, false))))
	case 7:
		return uint64(uint32(controlRaw(controlSlave(u, true))))
	case 10:
		controlGiveAbilities(u, int8(x), y)
		return 0
	case 15:
		return uint64(uint32(int32(controlRespawnFlags())))
	case 19:
		return uint64(uint32(controlRespawnNotify(u, byte(x))))
	case 23:
		controlClearWaypoints(u)
		return 0
	case 25:
		return uint64(uint32(controlConfusedDirection(u)))
	case 27:
		return uint64(uint32(bool2int(controlStartEligible(u, x))))
	case 28:
		return uint64(uint32(controlSubStamina(u, x)))
	case 29:
		controlAdjustStamina(u, int8(x))
		return 0
	case 30:
		return uint64(uint32(controlDropBall(u)))
	case 33:
		return uint64(uint32(controlNearStart(u, (*types.Pointf)(record))))
	case 34:
		return uint64(uint32(bool2int(controlHasWaypoint(u))))
	case 35:
		return uint64(uint32(controlWalkWaypoint(u)))
	case 36:
		return uint64(uint32(bool2int(controlCanMove(u))))
	case 37:
		return uint64(uint32(bool2int(controlCanAttack(u))))
	case 38:
		controlInputAttack(u)
		return 0
	case 39:
		return uint64(uint32(bool2int(controlAimsAtEnemy(u))))
	case 40:
		return uint64(uint32(controlFollowEnemy(u)))
	case 44:
		return uint64(uint32(controlBotCreate(u)))
	case 48:
		return uint64(uint32(int32(controlBotState(u))))
	case 49:
		return uint64(uint32(controlRespawnBot(u)))
	case 52:
		return uint64(uint32(controlGuideLevel(u, t)))
	case 54:
		return uint64(uint32(controlScheduledSpell(u, t, false)))
	case 55:
		return uint64(uint32(controlScheduledSpell(u, t, true)))
	default:
		return uint64(C.controlsCall(C.int(op), asObjectC(u), asObjectC(t), C.int(x), C.int(y), record, (*C.char)(name)))
	}
}
