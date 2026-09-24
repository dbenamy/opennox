//go:build porttest

package legacy

/*
#include "GAME3_2.h"
#include "GAME3_3.h"
int nox_objectCollideDefault(int,int,float*);
static void* stateFunction(int id){switch(id){
case 40:return (void*)nox_xxx_collideMonsterEventProc_4E83B0;
case 41:return (void*)nox_xxx_collideMimic_4E83D0;
case 42:return (void*)nox_xxx_collidePlayer_4E8460;
case 43:return (void*)nox_objectCollideDefault;
default:return 0;}}
*/
import "C"

import (
	"bytes"
	"math"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
)

type PortTestObjectStateSpec struct {
	RegisteredCollision bool
	HealthRefs          []int
	ActorName           string
	Target              int
	X, Y, Z             int32
	FloatBits           uint32
	Globals             map[int]uint32
	MissileList         []int
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
	oldX, oldY := dword_5d4594_1565628, dword_5d4594_1565632
	dword_5d4594_1565628 = 0
	dword_5d4594_1565632 = 0
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
		dword_5d4594_1565628, dword_5d4594_1565632 = oldX, oldY
		copy(raw, saved)
	}
}
func (p *portTestShopPools) objectStateItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.State
	if sp == nil {
		return
	}
	p.reservedFunctionIDs += 39 // retain two previous reservations and reserve 37 removed callback slots.
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
		if fn := C.stateFunction(C.int(i)); fn != nil {
			p.identify(fn, 89000+uint32(i))
		}
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
func objectStateCall(id int, u, t *server.Object, x, y, z int32, bits uint32, record unsafe.Pointer) uint64 {
	switch id {
	case 0:
		u.NeedSync()
	case 1:
		u.Sub_4E4500(uint32(x), uint32(y), z != 0)
		return uint64(uint32(uintptr(stateSyncEnd(u))))
	case 2:
		return uint64(uint32(uintptr(stateOnOff(u, x != 0))))
	case 3:
		stateRaise(u, math.Float32frombits(bits))
	case 4:
		return uint64(uint32(uintptr(stateAnimation(u, uint32(x)))))
	case 5:
		return uint64(uint32(uintptr(stateBuffs(u, uint32(x)))))
	case 6:
		return uint64(uint32(uintptr(stateAttributes(u, record))))
	case 7:
		return math.Float64bits(float64(u.Mass))
	case 8:
		stateRemoveSpawned(u)
	case 9:
		return uint64(uint32(bool2int(stateIsUnit(u))))
	case 10:
		return uint64(uint32(bool2int(stateIsPixie(u))))
	case 11:
		stateCleanup(x)
	case 12:
		return uint64(uint32(bool2int(u.HealthData != nil && GetServer().S().Frame()-u.Frame134 <= 1)))
	case 13:
		return math.Float64bits(stateDistance(u, t))
	case 14:
		return uint64(uint32(stateDirection((*types.Pointf)(unsafe.Add(u.CObj(), 56)), (*types.Pointf)(record))))
	case 15:
		return uint64(uint32(stateFront((*types.Pointf)(unsafe.Add(u.CObj(), 56)), x, (*types.Pointf)(record))))
	case 16:
		stateTeleport(u, (*types.Pointf)(record))
	case 17:
		u.Nox_xxx_objectUnkUpdateCoords_4E7290()
		return uint64(uint32(uintptr(unsafe.Pointer(u))))
	case 18:
		stateLoot(u, (*types.Pointf)(record))
	case 19:
		stateRememberAttacker(u, t)
	case 20:
		return uint64(uint32(int32(stateOn(u))))
	case 21:
		return uint64(uint32(stateOff(u)))
	case 25:
		*memmap.PtrUint32(0x5d4594, 1567712) = uint32(x)
		return uint64(uint32(x))
	case 26:
		return uint64(uint32(int32(stateFreeze(u, x))))
	case 27:
		return uint64(uint32(int32(stateUnfreeze(u, x))))
	case 28:
		statePet(u, t)
	case 29:
		stateRemoveMonitors(u, t)
	case 30:
		if u == nil {
			return 0
		}
		return uint64(uint32(u.ObjClass) >> 2 & 1)
	case 31:
		return uint64(uint32(bool2int(stateOwns(u, 1567716, "Crown"))))
	case 32:
		return uint64(uint32(bool2int(stateOwns(u, 1567720, "GameBall"))))
	case 33:
		return uint64(uint32(stateCount(u, uint32(x), uint32(y))))
	case 34:
		return uint64(uint32(bool2int(stateEqual(u, t))))
	case 35:
		statePostCreate(u)
	case 36:
		statePlayerVisibility(x)
	case 37:
		return uint64(uint32(stateResetPixie(u)))
	case 38:
		stateCloseDoor(u, record)
	case 39:
		return uint64(uint32(stateDoorNotify(u)))
	case 40:
		return uint64(uint32(uintptr(stateMonsterCollision(u, t))))
	case 41:
		return uint64(uint32(uintptr(stateMimicCollision(u, t))))
	case 42:
		statePlayerCollision(u, t)
	case 43:
		return 0
	}
	return 0
}

func (p *portTestShopPools) objectStateAction(a PortTestShopAction) uint32 {
	attack := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	sp := attack.State
	state := p.temporary.world.objectives.attack
	u := p.temporaryRef(attack.Actor)
	pointerReturn := a.Op == 1226 && u != nil && u.ObjClass&2 != 0 && u.ObjFlags&0x8002 == 0
	if sp.RegisteredCollision && a.Op >= 1240 && a.Op <= 1242 {
		name := []string{"MonsterCollide", "MimicCollide", "PlayerCollide"}[a.Op-1240]
		PortTestRegisteredCollision(u, p.temporaryRef(sp.Target), nil, name)
		state.state.result = 0
	} else if a.Op == 1222 {
		state.state.result = uint64(stateChecksum(u))
	} else if a.Op == 1223 {
		state.state.result = uint64(uint32(nox_xxx_inventoryGetFirst_4E7980(C.int(uintptr(u.CObj())))))
	} else if a.Op == 1224 {
		state.state.result = uint64(uint32(nox_xxx_inventoryGetNext_4E7990(C.int(uintptr(u.CObj())))))
	} else {
		state.state.result = objectStateCall(int(a.Op-1200), p.temporaryRef(attack.Actor), p.temporaryRef(sp.Target), sp.X, sp.Y, sp.Z, sp.FloatBits, state.record)
	}
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
	out = append(out, p.normalize(uint32(s.result)), uint32(s.result>>32), uint32(dword_5d4594_1565628), uint32(dword_5d4594_1565632))
	for _, off := range objectStateOffsets {
		out = append(out, p.normalize(*memmap.PtrUint32(0x5d4594, off)))
	}
	return out
}
