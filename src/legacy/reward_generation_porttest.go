//go:build porttest

package legacy

/*
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1568280,dword_5d4594_1568288;
extern uint64_t qword_581450_10256;
static void* rewardFunction(int id){switch(id){
case 0:return (void*)nox_xxx_unitSparkInit_4F0390;
case 1:return (void*)nox_xxx_initFrog_4F03B0;
case 2:return (void*)nox_xxx_initChest_4F0400;
case 3:return (void*)nox_xxx_unitBoulderInit_4F0420;
case 4:return (void*)sub_4F0450;
case 5:return (void*)sub_4F0490;
case 6:return (void*)nox_xxx_unitInitGold_4F04B0;
case 7:return (void*)nox_xxx_breakInit_4F0570;
case 8:return (void*)nox_xxx_unitInitGenerator_4F0590;
case 9:return (void*)nox_server_rewardgen_activateMarker_4F0720;
case 10:return (void*)nox_xxx_rewardSpellBook_4F09F0;
case 11:return (void*)nox_server_rewardGen_pickRandomSlots_4F0B60;
case 12:return (void*)nox_xxx_rewardAbilityBook_4F0C70;
case 13:return (void*)nox_xxx_rewardFieldGuide_4F0D20;
case 14:return (void*)nox_xxx_rewardMakeArmor_4F0E80;
case 15:return (void*)nox_xxx_rewardMakeWeapon_4F14E0;
case 16:return (void*)nox_xxx_rewardMakePotion_4F1C40;
case 17:return (void*)nox_xxx_createGem_4F1D30;
case 18:return (void*)nox_xxx_createGem2_4F1F00;
case 19:return (void*)sub_4F2110;
case 20:return (void*)sub_4F2210;
default:return 0;}}
static uint32_t rewardCall(int id,nox_object_t* u,uint32_t stage){switch(id){
case 0:return (uint32_t)nox_xxx_unitSparkInit_4F0390((int)u);
case 1:return (uint32_t)nox_xxx_initFrog_4F03B0((int)u);
case 2:return (uint32_t)nox_xxx_initChest_4F0400((int)u);
case 3:return (uint32_t)nox_xxx_unitBoulderInit_4F0420((uint32_t*)u);
case 4:return (uint32_t)sub_4F0450((int)u);
case 5:return (uint32_t)sub_4F0490((int)u);
case 6:return (uint32_t)nox_xxx_unitInitGold_4F04B0((int)u);
case 7:return (uint32_t)nox_xxx_breakInit_4F0570((int)u);
case 8:return (uint32_t)nox_xxx_unitInitGenerator_4F0590((int)u);
case 9:return (uint32_t)nox_server_rewardgen_activateMarker_4F0720((int)u,stage);
case 10:return (uint32_t)nox_xxx_rewardSpellBook_4F09F0((int)u,stage);
case 11:return (uint32_t)nox_server_rewardGen_pickRandomSlots_4F0B60(stage);
case 12:return (uint32_t)nox_xxx_rewardAbilityBook_4F0C70((int)u);
case 13:return (uint32_t)nox_xxx_rewardFieldGuide_4F0D20((int)u,stage);
case 14:return (uint32_t)nox_xxx_rewardMakeArmor_4F0E80((int)u,stage);
case 15:return (uint32_t)nox_xxx_rewardMakeWeapon_4F14E0((int)u,stage);
case 16:return (uint32_t)nox_xxx_rewardMakePotion_4F1C40((int)u,stage);
case 17:return (uint32_t)nox_xxx_createGem_4F1D30((int)u,stage);
case 18:return (uint32_t)nox_xxx_createGem2_4F1F00((int)u,stage);
case 19:sub_4F2110();return 0;
case 20:return (uint32_t)sub_4F2210();
default:return 0;}}
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"sort"
	"unsafe"
)

type PortTestRewardSpec struct {
	InitByRef              map[int]map[int]uint32
	UpdateName             string
	Stage                  uint32
	GeneratorStage         uint32
	InitWords, UpdateWords map[int]uint32
	ActorType              string
	Types                  map[int]string
	Missing                []string
	Allowed                bool
	ArmorBit, WeaponBit    uint32
	TableMode              int
}
type portTestReward struct {
	oldInit, oldUpdate map[*server.Object]unsafe.Pointer
	result             uint32
	configureSpells    func([]server.PortTestSpellClassDef)
	restoreModifier    func()
}

var rewardTypeNames = []string{"RewardMarker", "RewardMarkerPlus", "PortTestRewardArmor", "PortTestRewardWeapon", "PortTestRewardWand", "PortTestRewardPotion", "ConjurerSpellBook", "WizardSpellBook", "CommonSpellBook", "AbilityBook", "FieldGuide", "QuestGoldChest", "QuestGoldPile", "RubyGem", "EmeraldGem", "DiamondGem", "RedPotion", "Ankh"}
var rewardGlobalOffsets = []uintptr{1568276, 1568284, 1568292, 1568296, 2388660}

func (p *portTestShopPools) rewardPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Reward
	if sp == nil {
		return func() {}
	}
	state := &portTestReward{oldInit: make(map[*server.Object]unsafe.Pointer), oldUpdate: make(map[*server.Object]unsafe.Pointer)}
	p.temporary.world.objectives.attack.reward = state
	restoreTypes := p.proxy.core.PortTestRewardTypes(rewardTypeNames, sp.Missing, sp.Allowed, sp.ArmorBit, sp.WeaponBit)
	var old []uint32
	for _, off := range rewardGlobalOffsets {
		old = append(old, memmap.Uint32(0x5d4594, off))
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	configureSpells, restoreSpells := p.proxy.core.PortTestAISpellDefs()
	state.configureSpells = configureSpells
	goldA, goldB := memmap.Uint64(0x581450, 10248), memmap.Uint64(0x581450, 10264)
	goldC := C.qword_581450_10256
	constants := blobdata.PortTestRewardGoldConstants()
	*memmap.PtrUint64(0x581450, 10248) = binary.LittleEndian.Uint64(constants)
	C.qword_581450_10256 = C.uint64_t(binary.LittleEndian.Uint64(constants[8:]))
	*memmap.PtrUint64(0x581450, 10264) = binary.LittleEndian.Uint64(constants[16:])
	*memmap.PtrUint32(0x5d4594, 2388660) = sp.GeneratorStage
	oldA, oldB := C.dword_5d4594_1568280, C.dword_5d4594_1568288
	C.dword_5d4594_1568280 = 0
	C.dword_5d4594_1568288 = 0
	oldStage := memmap.Uint32(0x587000, 202028)
	*memmap.PtrUint32(0x587000, 202028) = sp.Stage
	// Mutable reward tables are isolated from other fixture families. Their
	// compact, explicit inputs are installed in rewardItems after object setup.
	raw := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 207044)), 4132)
	saved := bytes.Clone(raw)
	clear(raw)
	guide := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 70500)), 41*4)
	oldGuide := bytes.Clone(guide)
	return func() {
		for u, v := range state.oldInit {
			u.InitData = v
		}
		for u, v := range state.oldUpdate {
			u.UpdateData = v
		}
		copy(raw, saved)
		copy(guide, oldGuide)
		for i, off := range rewardGlobalOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		*memmap.PtrUint32(0x587000, 202028) = oldStage
		C.dword_5d4594_1568280, C.dword_5d4594_1568288 = oldA, oldB
		if state.restoreModifier != nil {
			state.restoreModifier()
		}
		*memmap.PtrUint64(0x581450, 10248) = goldA
		*memmap.PtrUint64(0x581450, 10264) = goldB
		C.qword_581450_10256 = goldC
		restoreSpells()
		restoreTypes()
	}
}
func (p *portTestShopPools) rewardAction(a PortTestShopAction) uint32 {
	attack := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	state := p.temporary.world.objectives.attack.reward
	state.result = uint32(C.rewardCall(C.int(a.Op-1300), asObjectC(p.temporaryRef(attack.Actor)), C.uint32_t(attack.Reward.Stage)))

	// These factories return newly allocated objects without placing them. Adopt
	// their returned allocation for capture and teardown; do not synthesize a
	// CreateObjectAt call or change any production object fields.
	if a.Op >= 1309 && a.Op <= 1318 && a.Op != 1311 && state.result != 0 {
		u := (*server.Object)(unsafe.Pointer(uintptr(state.result)))
		for _, old := range p.proxy.life.created {
			if old == u {
				panic("reward duplicate factory result")
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
	p.temporary.result = state.result
	return state.result
}
func (p *portTestShopPools) rewardItems() {
	attack := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	sp := attack.Reward
	if sp == nil {
		return
	}
	st := p.temporary.world.objectives.attack.reward
	u := p.temporaryRef(attack.Actor)
	setBuffers := func(o *server.Object) {
		if _, ok := st.oldInit[o]; ok {
			return
		}
		st.oldInit[o] = o.InitData
		st.oldUpdate[o] = o.UpdateData
		o.InitData = p.objectiveRegion(256)
		o.UpdateData = p.objectiveRegion(256)
	}
	setBuffers(u)
	for off, v := range sp.InitWords {
		if off < 0 || off+4 > 256 || off%4 != 0 {
			panic("reward init offset")
		}
		*equipmentWord(u.InitData, off) = v
	}
	for off, v := range sp.UpdateWords {
		if off < 0 || off+4 > 256 || off%4 != 0 {
			panic("reward update offset")
		}
		*equipmentWord(u.UpdateData, off) = v
	}
	if sp.ActorType != "" {
		u.TypeInd = uint16(p.proxy.core.Types.IndByID(sp.ActorType))
	}
	refs := make([]int, 0, len(sp.Types))
	for ref := range sp.Types {
		refs = append(refs, ref)
	}
	sort.Ints(refs)
	for _, ref := range refs {
		name := sp.Types[ref]
		o := p.temporaryRef(ref)
		setBuffers(o)
		o.TypeInd = uint16(p.proxy.core.Types.IndByID(name))
	}

	refs = refs[:0]
	for ref := range sp.InitByRef {
		refs = append(refs, ref)
	}
	sort.Ints(refs)
	for _, ref := range refs {
		o := p.temporaryRef(ref)
		setBuffers(o)
		for off, v := range sp.InitByRef[ref] {
			if off < 0 || off+4 > 256 || off%4 != 0 {
				panic("reward per-object init offset")
			}
			*equipmentWord(o.InitData, off) = v
		}
	}
	if sp.UpdateName != "" {
		if len(sp.UpdateName) >= 240 {
			panic("reward update name size")
		}
		copy(unsafe.Slice((*byte)(unsafe.Add(u.UpdateData, 16)), 240), sp.UpdateName)
	}
	st.configureSpells([]server.PortTestSpellClassDef{
		{Index: 1, Valid: true, Flags: uint32(things.SpellClassAny)},
		{Index: 2, Valid: true, Flags: uint32(things.SpellClassWizard)},
		{Index: 3, Valid: true, Flags: uint32(things.SpellClassConjurer)},
		{Index: 4, Valid: true},
		{Index: 5, Valid: false, Flags: uint32(things.SpellClassAny)},
	})
	repl := (*server.ModifierEff)(p.objectiveRegion(144))
	st.restoreModifier = p.proxy.core.PortTestRewardModifier(repl, (*byte)(p.objectiveString("Replenishment1")))
	for i := 0; i < 21; i++ {
		p.identify(C.rewardFunction(C.int(i)), 90000+uint32(i))
	}
	for i := 0; i < 41; i++ {
		*memmap.PtrPtr(0x587000, 70500+uintptr(i*4)) = p.objectiveString(fmt.Sprintf("PortTestGuide%02d", i))
	}
	// Explicit small tables use the production row layouts. Zero sentinels and
	// excluded rows participate in the tests rather than being silently omitted.
	for i := 0; i < 8; i++ {
		*memmap.PtrUint32(0x587000, 207044+uintptr(i*8)) = uint32(i + 1)
	}
	for _, base := range []uintptr{207104, 207792} {
		for i := 0; i < 5; i++ {
			off := base + uintptr(i*12)
			*memmap.PtrUint32(0x587000, off) = uint32(i + 1)
			*memmap.PtrUint32(0x587000, off+4) = uint32(i + 1)
			*memmap.PtrUint32(0x587000, off+8) = 31
		}
	}
	for i, name := range []string{"PortTestRewardWeapon", "PortTestRewardArmor", "PortTestRewardPotion", "PortTestRewardWand"} {
		off := uintptr(208176 + i*20)
		*memmap.PtrUint32(0x587000, off) = uint32(i + 1)
		text := name
		if i == 3 {
			text = "#PortTestRewardWand"
		}
		*memmap.PtrPtr(0x587000, off+4) = p.objectiveString(text)
		*memmap.PtrUint32(0x587000, off+8) = uint32(p.proxy.core.Types.IndByID(name))
		*memmap.PtrUint32(0x587000, off+12) = []uint32{1, 2, 4, 1}[i]
		*memmap.PtrUint32(0x587000, off+16) = 31
	}
	for _, base := range []uintptr{209336, 210704, 210848, 210992} {
		for i := 0; i < 3; i++ {
			off := base + uintptr(i*24)
			mod := p.objectiveRegion(144)
			*equipmentWord(mod, 28) = sp.WeaponBit
			*equipmentWord(mod, 32) = sp.ArmorBit
			*equipmentWord(mod, 36) = uint32(i + 1)
			*memmap.PtrUint32(0x587000, off) = uint32(i + 1)
			*memmap.PtrPtr(0x587000, off+4) = mod
			*memmap.PtrPtr(0x587000, off+8) = p.objectiveString(fmt.Sprintf("RewardMod%d-%d", base, i))
			*memmap.PtrUint32(0x587000, off+12) = 31
		}
	}
	for i := 0; i < 5; i++ {
		*memmap.PtrUint32(0x587000, 211136+uintptr(i*8)) = uint32(10 + 10*i)
		*memmap.PtrUint32(0x587000, 211140+uintptr(i*8)) = uint32(15 + 10*i)
	}

	if sp.TableMode == 6 || sp.TableMode == 7 {
		for _, base := range []uintptr{209336, 210704, 210848, 210992} {
			for i := 0; i < 3; i++ {
				if sp.TableMode == 6 {
					*memmap.PtrUint32(0x587000, base+uintptr(i*24)+12) = 1 << uint(i)
				}
				if sp.TableMode == 7 {
					*memmap.PtrUint32(0x587000, base+uintptr(i*24)) = 1
				}
			}
		}
	}
	switch sp.TableMode {
	case 1:
		// No table entries, while explicit per-marker choices remain usable.
		for _, off := range []uintptr{207108, 207796, 208180, 209344, 210712, 210856, 211000} {
			*memmap.PtrUint32(0x587000, off) = 0
		}
	case 2:
		for i := 0; i < 4; i++ {
			*memmap.PtrUint32(0x587000, 208192+uintptr(i*20)) = 0
		}
	case 4:
		*memmap.PtrUint32(0x587000, 208188) = 0
	case 5:
		*memmap.PtrUint32(0x587000, 208248) = 0
	case 3:
		for _, base := range []uintptr{209336, 210704, 210848, 210992} {
			for i := 0; i < 3; i++ {
				*memmap.PtrUint32(0x587000, base+uintptr(i*24)+16) = sp.ArmorBit
				*memmap.PtrUint32(0x587000, base+uintptr(i*24)+20) = sp.WeaponBit
			}
		}
	}
}
func (p *portTestShopPools) rewardSnapshot(out []uint32) []uint32 {
	st := p.temporary.world.objectives.attack.reward
	if st == nil {
		return out
	}
	out = append(out, p.normalize(st.result), uint32(C.dword_5d4594_1568280), uint32(C.dword_5d4594_1568288))
	for _, off := range rewardGlobalOffsets {
		out = append(out, p.normalize(memmap.Uint32(0x5d4594, off)))
	}
	for _, u := range p.proxy.life.created {
		for _, region := range []struct {
			ptr   unsafe.Pointer
			size  int
			guard byte
		}{{u.InitData, 256, 0x5a}, {u.UpdateData, 64, 0xa5}, {u.UseData.Ptr, 64, 0x5a}} {
			if region.ptr == nil {
				panic("reward created storage missing")
			}
			data := unsafe.Slice((*byte)(region.ptr), region.size+16)
			for _, b := range data[region.size:] {
				if b != region.guard {
					panic("reward created storage guard")
				}
			}
			for _, v := range unsafe.Slice((*uint32)(region.ptr), region.size/4) {
				out = append(out, p.normalize(v))
			}
		}
	}
	return out
}
