//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_3.h"
#include "GAME4_3.h"
extern uint32_t dword_5d4594_527656;
extern unsigned int dword_5d4594_2650652;
extern uint32_t dword_5d4594_1567988;
extern uint32_t dword_5d4594_527660;
extern nox_server_netCodeCacheStruct nox_server_netCodeCache;
extern uint32_t nox_server_needInitNetCodeCache;
int sub_417F50(int a1);
void nox_xxx_pickupFlagCtf_4EA490(int a1, int a2);
int sub_4EB9B0(int a1, int a2);
void nox_xxx_collideBall_4EBA00(int a1, int a2);
int sub_4EBB50(int a1, int a2);
short nox_xxx_collideHomeBase_4EBB80(int a1, int a2);
int sub_4ECBD0(int a1);
int sub_4ECC00(char** a1);
signed int nox_xxx_updateObelisk_53C580(int a1);
int nox_xxx_updateFlag_53DDF0(int a1);
void nox_xxx_updateGameBall_53DF40(int a3);
void nox_xxx_updateCrown_53E1D0(int a1);
void sub_4EA400(int a1, int a2);
int sub_4EA7A0(int a1);
short sub_4EA800(int a1, int a2);
static void* objectiveFunction(int id){switch(id){
case 0:return (void*)sub_417F50;
case 1:return (void*)nox_xxx_pickupFlagCtf_4EA490;
case 2:return (void*)sub_4EB9B0;
case 3:return (void*)nox_xxx_collideBall_4EBA00;
case 4:return (void*)sub_4EBB50;
case 5:return (void*)nox_xxx_collideHomeBase_4EBB80;
case 6:return (void*)sub_4ECBD0;
case 7:return (void*)sub_4ECC00;
case 8:return (void*)nox_xxx_updateObelisk_53C580;
case 9:return (void*)nox_xxx_updateFlag_53DDF0;
case 10:return (void*)nox_xxx_updateGameBall_53DF40;
case 11:return (void*)nox_xxx_updateCrown_53E1D0;
case 12:return (void*)sub_4EA400;
case 13:return (void*)sub_4EA7A0;
case 14:return (void*)sub_4EA800;
default:return 0;}}
static uint32_t objectiveCall(int id,nox_object_t* u,nox_object_t* target,int value){switch(id){
case 0: return (uint32_t)sub_417F50(value?(int)u:0);
case 1: nox_xxx_pickupFlagCtf_4EA490((int)u,(int)target);return 0;
case 2: return (uint32_t)sub_4EB9B0((int)u,(int)target);
case 3: nox_xxx_collideBall_4EBA00((int)u,(int)target);return 0;
case 4: return (uint32_t)sub_4EBB50((int)u,(int)target);
case 5: return (uint32_t)nox_xxx_collideHomeBase_4EBB80((int)u,(int)target);
case 6: return (uint32_t)sub_4ECBD0((int)u);
case 7: return (uint32_t)sub_4ECC00(*(char***)(*(char**)((char*)u+692)+4));
case 8: return (uint32_t)nox_xxx_updateObelisk_53C580((int)u);
case 9: return (uint32_t)nox_xxx_updateFlag_53DDF0((int)u);
case 10: nox_xxx_updateGameBall_53DF40((int)u);return 0;
case 11: nox_xxx_updateCrown_53E1D0((int)u);return 0;
case 12: sub_4EA400((int)u,(int)target);return 0;
case 13: return (uint32_t)sub_4EA7A0((int)target);
case 14: return (uint32_t)sub_4EA800((int)u,(int)target);
default:return 0;}}
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

const (
	PortTestObjective417F50 = 800
	PortTestObjective4EA490 = 801
	PortTestObjective4EB9B0 = 802
	PortTestObjective4EBA00 = 803
	PortTestObjective4EBB50 = 804
	PortTestObjective4EBB80 = 805
	PortTestObjective4ECBD0 = 806
	PortTestObjective4ECC00 = 807
	PortTestObjective53C580 = 808
	PortTestObjective53DDF0 = 809
	PortTestObjective53DF40 = 810
	PortTestObjective53E1D0 = 811
	PortTestObjective4EA400 = 812
	PortTestObjective4EA7A0 = 813
	PortTestObjective4EA800 = 814
)

// References 100..102 address the existing three real fixture players.
type PortTestObjectivesSpec struct {
	SpellDefinitions                                []server.PortTestSpellClassDef
	Ticks                                           []uint64
	Players                                         int
	ObjectList                                      []int
	PlayerWords, PlayerUpdateWords, PlayerDataWords []map[int]uint32
	PlayerRefs, PlayerUpdateRefs, PlayerDataRefs    []map[int]int
	UseWords                                        []map[int]uint32
	ColorNames                                      []string
	MissingTypes                                    []string
	ManaMultipliers                                 [3]float32
	ScoreLimit                                      uint16
}
type portTestObjectives struct {
	packets   []uint32
	ticks     []uint64
	blocks    [][]byte
	frees     []func()
	initWords []uint32
}

func (p *portTestShopPools) objectiveRegion(size int) unsafe.Pointer {
	r := p.temporary.world.objectives
	size = (size + 3) &^ 3
	b, f := alloc.Make([]byte{}, size+16)
	for i := 0; i < 8; i++ {
		b[i] = 0xa5
		b[len(b)-8+i] = 0x5a
	}
	r.blocks = append(r.blocks, b)
	r.frees = append(r.frees, f)
	ptr := unsafe.Pointer(&b[8])
	p.identify(ptr, 84000+uint32(len(r.blocks)))
	return ptr
}
func (p *portTestShopPools) objectiveString(s string) unsafe.Pointer {
	ptr := p.objectiveRegion(len(s) + 1)
	copy(unsafe.Slice((*byte)(ptr), len(s)), s)
	return ptr
}
func (p *portTestShopPools) objectivesPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives
	if sp == nil {
		return func() {}
	}
	p.temporary.world.objectives = &portTestObjectives{}
	core := p.proxy.core
	restoreTypes := core.PortTestObjectiveTypes(C.objectiveFunction(10), C.objectiveFunction(3), sp.MissingTypes)
	setSpells, restoreSpells := core.PortTestAISpellDefs()
	if len(sp.SpellDefinitions) != 0 {
		setSpells(sp.SpellDefinitions)
	}
	oldSend := core.NetSendPacketXxx
	core.NetSendPacketXxx = func(a int, b []byte, c, d, e int) int {
		r := p.temporary.world.objectives
		r.packets = append(r.packets, uint32(a), uint32(c), uint32(d), uint32(e), uint32(len(b)))
		for _, v := range b {
			r.packets = append(r.packets, uint32(v))
		}
		return len(b)
	}
	oldTicks := PlatformTicks
	PlatformTicks = func() uint64 {
		r := p.temporary.world.objectives
		v := uint64(10000)
		if len(sp.Ticks) != 0 {
			i := len(r.ticks)
			if i >= len(sp.Ticks) {
				i = len(sp.Ticks) - 1
			}
			v = sp.Ticks[i]
		}
		r.ticks = append(r.ticks, v)
		return v
	}
	oldList := core.Objs.List
	oldNetCache, oldNetInit := C.nox_server_netCodeCache, C.nox_server_needInitNetCodeCache
	C.nox_server_netCodeCache_initArray_4ECE50()
	oldStart, oldTeamBall := C.dword_5d4594_1567988, C.dword_5d4594_527660
	C.dword_5d4594_527660 = 0
	C.dword_5d4594_1567988 = 0
	oldCache, oldQuest := C.dword_5d4594_527656, C.dword_5d4594_2650652
	C.dword_5d4594_527656, C.dword_5d4594_2650652 = 0, 0
	offsets := []uintptr{1567992, 1567996, 1568008, 1567720}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		v := memmap.PtrUint32(0x5d4594, off)
		old[i] = *v
		*v = 0
	}
	status := unsafe.Slice((*byte)(memmap.PtrOff(0x5d4594, 1567736)), 208)
	oldStatus := bytes.Clone(status)
	clear(status)
	limits := unsafe.Slice(memmap.PtrUint16(0x5d4594, 3488), 6)
	oldLimits := append([]uint16(nil), limits...)
	for i := range limits {
		limits[i] = sp.ScoreLimit
	}
	colors := unsafe.Slice(memmap.PtrUint32(0x587000, 205224), 8)
	oldColors := append([]uint32(nil), colors...)
	clear(colors)
	for i, name := range []string{"RedFlag", "BlueFlag", "GoldFlag"} {
		colors[2*i] = uint32(uintptr(p.objectiveString(name)))
		colors[2*i+1] = uint32(i + 1)
	}
	oldMult := core.Players.Mult
	core.Players.Mult.Warrior.Mana = sp.ManaMultipliers[0]
	core.Players.Mult.Wizard.Mana = sp.ManaMultipliers[1]
	core.Players.Mult.Conjurer.Mana = sp.ManaMultipliers[2]
	var savedUnits []server.Object
	var savedUD [][]byte
	var savedPlayers []server.Player
	for i := range p.proxy.life.players {
		u := &p.proxy.life.players[i]
		pl := u.UpdateDataPlayer().Player
		savedUnits = append(savedUnits, *u)
		savedUD = append(savedUD, bytes.Clone(unsafe.Slice((*byte)(u.UpdateData), int(unsafe.Sizeof(server.PlayerUpdateData{})))))
		savedPlayers = append(savedPlayers, *pl)
		pl.Active = 0
		pl.PlayerUnit = nil
		if i < sp.Players {
			pl.Active = 1
			pl.PlayerUnit = u
		}
		u.NetCode = uint32(3000 + i)
		pl.NetCodeVal = u.NetCode
	}
	return func() {
		core.Objs.List = oldList
		C.nox_server_netCodeCache, C.nox_server_needInitNetCodeCache = oldNetCache, oldNetInit
		C.dword_5d4594_1567988 = oldStart
		C.dword_5d4594_527660 = oldTeamBall
		for i := range savedUnits {
			u := &p.proxy.life.players[i]
			pl := u.UpdateDataPlayer().Player
			// Inventory teardown still owns newly created minimap nodes.
			minimap := pl.Field4580
			*pl = savedPlayers[i]
			pl.Field4580 = minimap
			*u = savedUnits[i]
			copy(unsafe.Slice((*byte)(u.UpdateData), len(savedUD[i])), savedUD[i])
		}
		for i, v := range p.temporary.world.objectives.initWords {
			*equipmentWord(p.items[i].u.InitData, 4) = v
		}
		PlatformTicks = oldTicks
		core.NetSendPacketXxx = oldSend
		restoreSpells()
		C.dword_5d4594_527656, C.dword_5d4594_2650652 = oldCache, oldQuest
		for i, off := range offsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		copy(status, oldStatus)
		copy(limits, oldLimits)
		copy(colors, oldColors)
		core.Players.Mult = oldMult
		restoreTypes()
		for _, f := range p.temporary.world.objectives.frees {
			f()
		}
	}
}
func (p *portTestShopPools) objectivesItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives
	if sp == nil {
		return
	}
	r := p.temporary.world.objectives
	for i := 0; i < 15; i++ {
		p.identify(C.objectiveFunction(C.int(i)), 81000+uint32(i))
	}
	for i := range p.proxy.core.Teams.Arr {
		p.identify(p.proxy.core.Teams.Arr[i].C(), 82000+uint32(i))
	}
	apply := func(ptr unsafe.Pointer, size int, words map[int]uint32, refs map[int]int) {
		for off, v := range words {
			if off < 0 || off+4 > size || off%4 != 0 {
				panic("objective word offset")
			}
			*(*uint32)(unsafe.Add(ptr, off)) = v
		}
		for off, id := range refs {
			if off < 0 || off+4 > size || off%4 != 0 {
				panic("objective reference offset")
			}
			*(*unsafe.Pointer)(unsafe.Add(ptr, off)) = p.temporaryRef(id).CObj()
		}
	}
	getW := func(a []map[int]uint32, i int) map[int]uint32 {
		if i < len(a) {
			return a[i]
		}
		return nil
	}
	getR := func(a []map[int]int, i int) map[int]int {
		if i < len(a) {
			return a[i]
		}
		return nil
	}
	for i := range p.proxy.life.players {
		u := &p.proxy.life.players[i]
		apply(u.CObj(), 772, getW(sp.PlayerWords, i), getR(sp.PlayerRefs, i))
		u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
		apply(u.UpdateData, int(unsafe.Sizeof(server.PlayerUpdateData{})), getW(sp.PlayerUpdateWords, i), getR(sp.PlayerUpdateRefs, i))
		apply(unsafe.Pointer(u.UpdateDataPlayer().Player), int(unsafe.Sizeof(server.Player{})), getW(sp.PlayerDataWords, i), getR(sp.PlayerDataRefs, i))
	}
	for i, it := range p.items {
		p.identify(unsafe.Pointer(&it.u.TeamVal), 83000+uint32(i))
		r.initWords = append(r.initWords, *equipmentWord(it.u.InitData, 4))
		*equipmentWord(it.u.InitData, 4) = 0
		if i < len(sp.ColorNames) && sp.ColorNames[i] != "" {
			name := p.objectiveString(sp.ColorNames[i])
			header := p.objectiveRegion(4)
			*(*unsafe.Pointer)(header) = name
			*(*unsafe.Pointer)(unsafe.Add(it.u.InitData, 4)) = header
		}
		apply(it.u.UseData.Ptr, 128, getW(sp.UseWords, i), nil)
	}
	p.proxy.core.Objs.List = nil
	for i := len(sp.ObjectList) - 1; i >= 0; i-- {
		u := p.temporaryRef(sp.ObjectList[i])
		u.ObjNext = p.proxy.core.Objs.List
		p.proxy.core.Objs.List = u
	}
}
func (p *portTestShopPools) objectivesAction(a PortTestShopAction) uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	p.temporary.result = uint32(C.objectiveCall(C.int(a.Op-800), asObjectC(p.items[a.Item].u), asObjectC(p.temporaryRef(sp.Target)), C.int(a.Value)))
	return p.temporary.result
}
func (p *portTestShopPools) objectivesSnapshot(out []uint32) []uint32 {
	r := p.temporary.world.objectives
	if r == nil {
		return out
	}
	out = append(out, uint32(len(r.packets)))
	out = append(out, r.packets...)
	out = append(out, uint32(len(r.ticks)))
	for _, v := range r.ticks {
		out = append(out, uint32(v), uint32(v>>32))
	}
	words := func(ptr unsafe.Pointer, n int) {
		for _, v := range unsafe.Slice((*uint32)(ptr), n) {
			out = append(out, p.normalize(v))
		}
	}
	out = append(out, uint32(C.dword_5d4594_527656), p.normalize(uint32(uintptr(p.proxy.core.Objs.List.CObj()))))
	words(memmap.PtrOff(0x5d4594, 1567736), 208/4)
	words(memmap.PtrOff(0x5d4594, 1567992), 2)
	out = append(out, uint32(C.dword_5d4594_1567988), *memmap.PtrUint32(0x5d4594, 1568008))
	for i := range p.proxy.life.players {
		u := &p.proxy.life.players[i]
		words(u.CObj(), 193)
		words(u.UpdateData, int(unsafe.Sizeof(server.PlayerUpdateData{}))/4)
		words(unsafe.Pointer(u.UpdateDataPlayer().Player), int(unsafe.Sizeof(server.Player{}))/4)
	}
	out = append(out, uint32(p.proxy.core.Teams.ActiveCnt), uint32(C.dword_5d4594_527660))
	for i := range p.proxy.core.Teams.Arr {
		words(p.proxy.core.Teams.Arr[i].C(), int(unsafe.Sizeof(server.Team{}))/4)
	}
	for _, b := range r.blocks {
		for j := 0; j < 8; j++ {
			if b[j] != 0xa5 || b[len(b)-8+j] != 0x5a {
				panic("objective buffer guard")
			}
		}
		words(unsafe.Pointer(&b[8]), (len(b)-16)/4)
	}
	out = append(out, math.Float32bits(p.proxy.core.Players.Mult.Warrior.Mana), math.Float32bits(p.proxy.core.Players.Mult.Wizard.Mana), math.Float32bits(p.proxy.core.Players.Mult.Conjurer.Mana))
	return out
}
