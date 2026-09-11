//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME3_2.h"
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1565628,dword_5d4594_1565632;
int nox_objectCollideDefault(int,int,float*);
static void* stateFunction(int id){switch(id){
case 0:return (void*)nox_xxx_unitNeedSync_4E44F0;
case 1:return (void*)sub_4E4500;
case 2:return (void*)nox_xxx_unitSetOnOff_4E4670;
case 3:return (void*)nox_xxx_unitRaise_4E46F0;
case 4:return (void*)nox_xxx_servMarkObjAnimFrame_4E4880;
case 5:return (void*)nox_xxx_setUnitBuffFlags_4E48F0;
case 6:return (void*)nox_xxx_modifSetItemAttrs_4E4990;
case 7:return (void*)nox_xxx_objectGetMass_4E4A70;
case 8:return (void*)nox_xxx_playerRemoveSpawnedStuff_4E5AD0;
case 9:return (void*)nox_xxx_isUnit_4E5B50;
case 10:return (void*)sub_4E5B80;
case 11:return (void*)sub_4E5BF0;
case 12:return (void*)sub_4E6BD0;
case 13:return (void*)nox_xxx_calcDistance_4E6C00;
case 14:return (void*)sub_4E6CE0;
case 15:return (void*)nox_server_testTwoPointsAndDirection_4E6E50;
case 16:return (void*)nox_xxx_teleportToMB_4E7190;
case 17:return (void*)nox_xxx_objectUnkUpdateCoords_4E7290;
case 18:return (void*)nox_xxx_spawnSomeBarrel_4E7470;
case 19:return (void*)sub_4E7540;
case 20:return (void*)nox_xxx_objectSetOn_4E75B0;
case 21:return (void*)nox_xxx_objectSetOff_4E7600;
case 22:return (void*)sub_4E7700;
case 23:return (void*)nox_xxx_inventoryGetFirst_4E7980;
case 24:return (void*)nox_xxx_inventoryGetNext_4E7990;
case 25:return (void*)sub_4E79B0;
case 26:return (void*)nox_xxx_unitFreeze_4E79C0;
case 27:return (void*)nox_xxx_unitUnFreeze_4E7A60;
case 28:return (void*)nox_xxx_unitBecomePet_4E7B00;
case 29:return (void*)nox_xxx_monsterRemoveMonitors_4E7B60;
case 30:return (void*)sub_4E7BC0;
case 31:return (void*)nox_xxx_unitIsCrown_4E7BE0;
case 32:return (void*)nox_xxx_unitIsGameball_4E7C30;
case 33:return (void*)nox_xxx_unitCountSlaves_4E7CF0;
case 34:return (void*)sub_4E7DE0;
case 35:return (void*)nox_xxx_unitPostCreateNotify_4E7F10;
case 36:return (void*)sub_4E8110;
case 37:return (void*)sub_4E81D0;
case 38:return (void*)nox_xxx_fnFindCloseDoors_4E8340;
case 39:return (void*)sub_4E8390;
case 40:return (void*)nox_xxx_collideMonsterEventProc_4E83B0;
case 41:return (void*)nox_xxx_collideMimic_4E83D0;
case 42:return (void*)nox_xxx_collidePlayer_4E8460;
case 43:return (void*)nox_objectCollideDefault;
default:return 0;}}
static uint64_t stateCall(int id,nox_object_t* u,nox_object_t* t,int x,int y,int z,uint32_t bits,void* record){float f;memcpy(&f,&bits,4);switch(id){
case 0:{nox_xxx_unitNeedSync_4E44F0(u);return 0;}
case 1:{return (uint32_t)sub_4E4500(u,x,y,z);}
case 2:{return (uint32_t)nox_xxx_unitSetOnOff_4E4670((int)u,x);}
case 3:{nox_xxx_unitRaise_4E46F0(u,f);return 0;}
case 4:{return (uint32_t)nox_xxx_servMarkObjAnimFrame_4E4880((int)u,x);}
case 5:{return (uint32_t)nox_xxx_setUnitBuffFlags_4E48F0((int)u,x);}
case 6:{return (uint32_t)nox_xxx_modifSetItemAttrs_4E4990(u,(int*)record);}
case 7:{double d=nox_xxx_objectGetMass_4E4A70((int)u);uint64_t raw;memcpy(&raw,&d,8);return raw;}
case 8:{nox_xxx_playerRemoveSpawnedStuff_4E5AD0(u);return 0;}
case 9:{return (uint32_t)nox_xxx_isUnit_4E5B50(u);}
case 10:{return (uint32_t)sub_4E5B80(u);}
case 11:{sub_4E5BF0(x);return 0;}
case 12:{return (uint32_t)sub_4E6BD0((int)u);}
case 13:{double d=nox_xxx_calcDistance_4E6C00(u,t);uint64_t raw;memcpy(&raw,&d,8);return raw;}
case 14:{return (uint32_t)sub_4E6CE0((float2*)((char*)u+56),(float2*)record);}
case 15:{return (uint32_t)nox_server_testTwoPointsAndDirection_4E6E50((float2*)((char*)u+56),x,(float2*)record);}
case 16:{nox_xxx_teleportToMB_4E7190((uint8_t*)u,(float*)record);return 0;}
case 17:{return (uint32_t)nox_xxx_objectUnkUpdateCoords_4E7290(u);}
case 18:{nox_xxx_spawnSomeBarrel_4E7470((int)u,(int)record);return 0;}
case 19:{sub_4E7540(u,t);return 0;}
case 20:{return (uint32_t)nox_xxx_objectSetOn_4E75B0(u);}
case 21:{return (uint32_t)nox_xxx_objectSetOff_4E7600(u);}
case 22:{return (uint32_t)sub_4E7700((int)u);}
case 23:{return (uint32_t)nox_xxx_inventoryGetFirst_4E7980((int)u);}
case 24:{return (uint32_t)nox_xxx_inventoryGetNext_4E7990((int)u);}
case 25:{return (uint32_t)sub_4E79B0(x);}
case 26:{return (uint32_t)nox_xxx_unitFreeze_4E79C0(u,x);}
case 27:{return (uint32_t)nox_xxx_unitUnFreeze_4E7A60(u,x);}
case 28:{nox_xxx_unitBecomePet_4E7B00((int)u,(int)t);return 0;}
case 29:{nox_xxx_monsterRemoveMonitors_4E7B60(u,t);return 0;}
case 30:{return (uint32_t)sub_4E7BC0((int)u);}
case 31:{return (uint32_t)nox_xxx_unitIsCrown_4E7BE0((int)u);}
case 32:{return (uint32_t)nox_xxx_unitIsGameball_4E7C30((int)u);}
case 33:{return (uint32_t)nox_xxx_unitCountSlaves_4E7CF0((int)u,x,y);}
case 34:{return (uint32_t)sub_4E7DE0((int)u,t);}
case 35:{return (uint32_t)nox_xxx_unitPostCreateNotify_4E7F10(u);}
case 36:{return (uint32_t)sub_4E8110(x);}
case 37:{return (uint32_t)sub_4E81D0(u);}
case 38:{nox_xxx_fnFindCloseDoors_4E8340((float*)u,(int)record);return 0;}
case 39:{return (uint32_t)sub_4E8390((int)u);}
case 40:{return (uint32_t)nox_xxx_collideMonsterEventProc_4E83B0((int)u,(int)t);}
case 41:{return (uint32_t)nox_xxx_collideMimic_4E83D0((int)u,(int)t);}
case 42:{nox_xxx_collidePlayer_4E8460((int)u,(int)t);return 0;}
case 43:{return (uint32_t)nox_objectCollideDefault((int)u,(int)t,(float*)record);}
default:return 0;}}
*/
import "C"

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestObjectStateSpec struct {
	HealthRefs  []int
	ActorName   string
	Target      int
	X, Y, Z     int32
	FloatBits   uint32
	Globals     map[int]uint32
	MissileList []int
}
type portTestObjectState struct {
	result    uint64
	oldHealth map[*server.Object]*server.HealthData
	oldDamage map[*server.Object]unsafe.Pointer
	disabled  []uint32
}

var objectStateOffsets = []uintptr{1564960, 1565592, 1565596, 1565636, 1565640, 1565652, 1565656, 1567708, 1567712, 1567716, 1567720, 1567728}
var objectStateLootPointers = [][2]uintptr{
	{203080, 203420}, {203092, 203432}, {203104, 203444}, {203116, 203452}, {203128, 203460}, {203140, 203468}, {203152, 203476}, {203164, 203484}, {203176, 203492}, {203188, 203504}, {203200, 203516}, {203212, 203524},
	{203240, 203536}, {203252, 203548}, {203264, 203556}, {203276, 203564}, {203288, 203572}, {203300, 203580}, {203312, 203588}, {203324, 203608}, {203336, 203620}, {203348, 203632}, {203360, 203648}, {203372, 203660}, {203384, 203684}, {203396, 203700},
}

func objectStateTypeNames() []string {
	out := []string{"TeamBase", "Pixie", "Moonglow", "Crown", "BarrelPortTest", "CratePortTest"}
	data := blobdata.PortTestObjectStateLoot()
	for _, pair := range objectStateLootPointers {
		raw := data[pair[1]-203080:]
		out = append(out, string(raw[:bytes.IndexByte(raw, 0)]))
	}
	return out
}
func (p *portTestShopPools) objectStatePrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.State
	if sp == nil {
		return func() {}
	}
	p.temporary.world.objectives.attack.state = &portTestObjectState{oldHealth: make(map[*server.Object]*server.HealthData), oldDamage: make(map[*server.Object]unsafe.Pointer)}
	var old []uint32
	for _, off := range objectStateOffsets {
		old = append(old, *memmap.PtrUint32(0x5d4594, off))
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	for off, v := range sp.Globals {
		found := false
		for _, valid := range objectStateOffsets {
			if uintptr(off) == valid {
				found = true
				break
			}
		}
		if !found {
			panic("object-state global offset")
		}
		*memmap.PtrUint32(0x5d4594, uintptr(off)) = v
	}
	oldX, oldY := C.dword_5d4594_1565628, C.dword_5d4594_1565632
	C.dword_5d4594_1565628 = 0
	C.dword_5d4594_1565632 = 0
	raw := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 203080)), 644)
	saved := bytes.Clone(raw)
	copy(raw, blobdata.PortTestObjectStateLoot())
	for _, pair := range objectStateLootPointers {
		*memmap.PtrPtr(0x587000, pair[0]) = memmap.PtrOff(0x587000, pair[1])
	}
	oldDisable := Sub_4FC300
	Sub_4FC300 = func(u *server.Object, abil int) {
		if abil != 1 {
			panic("object-state unexpected ability disable")
		}
		st := p.temporary.world.objectives.attack.state
		st.disabled = append(st.disabled, p.normalize(uint32(uintptr(u.CObj()))), uint32(abil))
		p.proxy.core.Abils.DisableAbilityAaa(u, server.Ability(abil))
	}
	oldList := p.proxy.core.Objs.MissileList
	return func() {
		Sub_4FC300 = oldDisable
		for u, fn := range p.temporary.world.objectives.attack.state.oldDamage {
			u.Damage = fn
		}
		for u, hp := range p.temporary.world.objectives.attack.state.oldHealth {
			u.HealthData = hp
		}
		p.proxy.core.Objs.MissileList = oldList
		for i, off := range objectStateOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		C.dword_5d4594_1565628, C.dword_5d4594_1565632 = oldX, oldY
		copy(raw, saved)
	}
}
func (p *portTestShopPools) objectStateItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.State
	if sp == nil {
		return
	}
	if sp.ActorName != "" {
		u := p.temporaryRef(p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Actor)
		u.IDPtr = p.objectiveString(sp.ActorName)
	}
	for _, id := range []int{1, 2, 100, 101, 102} {
		if u := p.temporaryRef(id); u != nil {
			st := p.temporary.world.objectives.attack.state
			if _, ok := st.oldDamage[u]; !ok {
				st.oldDamage[u] = u.Damage
			}
			u.Damage = p.proxy.combat.target.Damage
		}
	}
	for _, id := range sp.HealthRefs {
		u := p.temporaryRef(id)
		p.temporary.world.objectives.attack.state.oldHealth[u] = u.HealthData
		hp := (*server.HealthData)(p.objectiveRegion(int(unsafe.Sizeof(server.HealthData{}))))
		hp.Cur = 100
		hp.Max = 100
		u.HealthData = hp
	}
	for i := 0; i < 44; i++ {
		p.identify(C.stateFunction(C.int(i)), 89000+uint32(i))
	}
	for _, id := range []int{1, 2, 3, 4, 5, 100, 101, 102} {
		if u := p.temporaryRef(id); u != nil {
			p.identify(unsafe.Add(u.CObj(), 688), 89100+uint32(id))
		}
	}
	p.identify(unsafe.Add(p.temporary.world.objectives.attack.record, 16), 89300)
	p.proxy.core.Objs.MissileList = nil
	for i := len(sp.MissileList) - 1; i >= 0; i-- {
		u := p.temporaryRef(sp.MissileList[i])
		u.ObjNext = p.proxy.core.Objs.MissileList
		p.proxy.core.Objs.MissileList = u
	}
}
func (p *portTestShopPools) objectStateAction(a PortTestShopAction) uint32 {
	attack := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	sp := attack.State
	state := p.temporary.world.objectives.attack
	u := p.temporaryRef(attack.Actor)
	pointerReturn := a.Op == 1226 && u != nil && u.ObjClass&2 != 0 && u.ObjFlags&0x8002 == 0
	state.state.result = uint64(C.stateCall(C.int(a.Op-1200), asObjectC(p.temporaryRef(attack.Actor)), asObjectC(p.temporaryRef(sp.Target)), C.int(sp.X), C.int(sp.Y), C.int(sp.Z), C.uint32_t(sp.FloatBits), state.record))
	// Freeze's char return truncates the C-owned action-stack pointer. Check the
	// actual byte before replacing the address-dependent value with its identity.
	if pointerReturn {
		head := u.UpdateDataMonster().AIStackHead()
		want := uint32(int32(int8(uintptr(unsafe.Pointer(head)))))
		if uint32(state.state.result) == want {
			if head != nil {
				state.state.result = 89400
			}
		} else if state.state.result != 0 {
			panic("object-state freeze pointer return")
		}
	}
	p.temporary.result = uint32(state.state.result)
	return p.temporary.result
}
func (p *portTestShopPools) objectStateSnapshot(out []uint32) []uint32 {
	s := p.temporary.world.objectives.attack.state
	if s == nil {
		return out
	}
	out = append(out, uint32(len(s.disabled)))
	out = append(out, s.disabled...)
	out = append(out, p.normalize(uint32(s.result)), uint32(s.result>>32), uint32(C.dword_5d4594_1565628), uint32(C.dword_5d4594_1565632))
	for _, off := range objectStateOffsets {
		out = append(out, p.normalize(*memmap.PtrUint32(0x5d4594, off)))
	}
	return out
}
