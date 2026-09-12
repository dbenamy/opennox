//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
static uint32_t spellEffectsForceLog[384]; static int spellEffectsForceN;
static void spellEffectsForceRecord(int u,uint32_t distance,int* arg) {if(spellEffectsForceN+3>384)abort();spellEffectsForceLog[spellEffectsForceN++]=u;spellEffectsForceLog[spellEffectsForceN++]=distance;spellEffectsForceLog[spellEffectsForceN++]=(uint32_t)arg;}
static void* spellEffectsForcePtr(void){return spellEffectsForceRecord;}
static void spellEffectsForceReset(void){spellEffectsForceN=0;}
static int spellEffectsForceCount(void){return spellEffectsForceN;}
static uint32_t spellEffectsForceValue(int i){return spellEffectsForceLog[i];}
#include <string.h>
#include "GAME4.h"
#include "GAME1_1.h"
extern int nox_cheat_charmall;
extern uint32_t dword_5d4594_2487708,dword_5d4594_2487712,dword_5d4594_2487804;
#include "GAME4_2.h"
#include "GAME4_3.h"
int sub_500CA0(int a1, int a2);
int nox_xxx_summonStart_500DA0(int a1);
int sub_500F40(int a1, float a2);
int nox_xxx_summonFinish_5010D0(int a1);
void nox_xxx_summonCancel_5011C0(int a1);
int nox_xxx_charmCreature1_5011F0(int* a1);
int nox_xxx_charmCreatureFinish_5013E0(int* a1);
int nox_xxx_charmCreature2_501690(int a1);
void nox_xxx_banishUnit_5017F0(int unit);
int sub_52BEB0(int a1, int a2, int a3, int a4);
int nox_xxx_castSpellWinkORrestoreHealth_52BF20(int a1, int a2, int a3, int a4, int* a5);
int sub_52BF50(int a1, int a2, int a3, int a4, int* a5);
int nox_xxx_castPull_52BFA0(int a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_castPush_52C000(int a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_castFumble_52C060(int a1, int a2, int a3, int a4, int* a5);
int nox_xxx_castConfuse_52C1E0(int a1, int a2, int a3, int a4, int* a5, char a6);
int nox_xxx_castStun_52C2C0(int a1, int a2, int a3, int a4, int* a5, char a6);
int nox_xxx_castBurn_52C3E0(int a1, int a2, int a3, int a4, int a5);
int nox_xxx_useShock_52C5A0(int a1, int a2, int a3, int a4, int* a5, int a6);
int nox_xxx_castPoison_52C720(int a1, int a2, int a3, int a4, int* a5, int a6);
int nox_xxx_castFireball_52C790(int a1, int a2, int a3, int a4, int a5, int a6);
int sub_52CA80(int a1, int a2, int a3, int a4);
int sub_52CBD0(int a1, int a2, int a3, int a4);
int sub_52CCD0(int a1, int a2, int a3);
int nox_xxx_castCurePoison_52CDB0(int a1, int a2, int a3, int a4, int* a5, int a6);
void sub_52CE60(int a1);
int nox_xxx_castLock_52CE90(int a1, int a2, int a3, int a4);
void sub_52CF90(int a1, int a2);
void sub_52D060(int a1, int a2);
int nox_xxx_castTelekinesis_52D330(int a1, int a2, int a3, int a4, int* a5, char a6);
int nox_xxx_castFist_52D3C0(int a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_spellCastCleansingFlame_52D5C0(int a1, nox_object_t* a2p, nox_object_t* a3p, nox_object_t* a4p, void* a5p, int a6);
int nox_xxx_castMeteorShower_52D8A0(int a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_castMeteor_52D9D0(int a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_castToxicCloud_52DB60(int a1, int a2, int a3, int a4, int a5);
int nox_xxx_spellArachna_52DC80(int a1, int a2, int a3, int a4, int a5);
int sub_52DD50(int a1, int a2, int a3, int a4, void* a5);
int nox_xxx_castEquake_52DE40(int a1, int a2, int a3, int a4, int a5, int a6);
short nox_xxx_equakeDamage_52DEC0(int a1, int a2);
unsigned int nox_xxx_isObjectMovable_52E020(int a1);
void nox_xxx_mapPushUnitsAround_52E040(void* a1p, float a2, float a3p, float a4, nox_object_t* a5p, int a6, int a7);
void nox_xxx_unitPushAroundFn_52E0E0(int a1, int** a2);
static void* spellEffectsFunction(int op) {switch(op){
case 0: return sub_500CA0;
case 1: return nox_xxx_summonStart_500DA0;
case 2: return sub_500F40;
case 3: return nox_xxx_summonFinish_5010D0;
case 4: return nox_xxx_summonCancel_5011C0;
case 5: return nox_xxx_charmCreature1_5011F0;
case 6: return nox_xxx_charmCreatureFinish_5013E0;
case 7: return nox_xxx_charmCreature2_501690;
case 8: return nox_xxx_banishUnit_5017F0;
case 9: return sub_52BEB0;
case 10: return nox_xxx_castSpellWinkORrestoreHealth_52BF20;
case 11: return sub_52BF50;
case 12: return nox_xxx_castPull_52BFA0;
case 13: return nox_xxx_castPush_52C000;
case 14: return nox_xxx_castFumble_52C060;
case 15: return nox_xxx_castConfuse_52C1E0;
case 16: return nox_xxx_castStun_52C2C0;
case 17: return nox_xxx_castBurn_52C3E0;
case 18: return nox_xxx_useShock_52C5A0;
case 19: return nox_xxx_castPoison_52C720;
case 20: return nox_xxx_castFireball_52C790;
case 21: return sub_52CA80;
case 22: return sub_52CBD0;
case 23: return sub_52CCD0;
case 24: return nox_xxx_castCurePoison_52CDB0;
case 25: return sub_52CE60;
case 26: return nox_xxx_castLock_52CE90;
case 27: return sub_52CF90;
case 28: return sub_52D060;
case 29: return nox_xxx_castTelekinesis_52D330;
case 30: return nox_xxx_castFist_52D3C0;
case 31: return nox_xxx_spellCastCleansingFlame_52D5C0;
case 32: return nox_xxx_castMeteorShower_52D8A0;
case 33: return nox_xxx_castMeteor_52D9D0;
case 34: return nox_xxx_castToxicCloud_52DB60;
case 35: return nox_xxx_spellArachna_52DC80;
case 36: return sub_52DD50;
case 37: return nox_xxx_castEquake_52DE40;
case 38: return nox_xxx_equakeDamage_52DEC0;
case 39: return nox_xxx_isObjectMovable_52E020;
case 40: return nox_xxx_mapPushUnitsAround_52E040;
case 41: return nox_xxx_unitPushAroundFn_52E0E0;
default:return 0;}}
static uint32_t spellEffectsCall(int op,int id,nox_object_t* u,nox_object_t* a,nox_object_t* b,nox_object_t* c,void* record,void* output,int level,uint32_t fx,uint32_t fy,uint32_t fz,int q,int r) {
float x,y,z,pout;memcpy(&x,&fx,4);memcpy(&y,&fy,4);memcpy(&z,&fz,4);memcpy(&pout,&output,4);
switch(op){
case 0: return (uint32_t)sub_500CA0(id,(int)u);
case 1: return (uint32_t)nox_xxx_summonStart_500DA0((int)record);
case 2: return (uint32_t)sub_500F40((int)record,pout);
case 3: return (uint32_t)nox_xxx_summonFinish_5010D0((int)record);
case 4: nox_xxx_summonCancel_5011C0((int)record);return 0;
case 5: return (uint32_t)nox_xxx_charmCreature1_5011F0((int*)record);
case 6: return (uint32_t)nox_xxx_charmCreatureFinish_5013E0((int*)record);
case 7: return (uint32_t)nox_xxx_charmCreature2_501690((int)record);
case 8: nox_xxx_banishUnit_5017F0((int)u);return 0;
case 9: return (uint32_t)sub_52BEB0(id,(int)a,(int)b,(int)c);
case 10: return (uint32_t)nox_xxx_castSpellWinkORrestoreHealth_52BF20(id,(int)a,(int)b,(int)c,(int*)record);
case 11: return (uint32_t)sub_52BF50(id,(int)a,(int)b,(int)c,(int*)record);
case 12: return (uint32_t)nox_xxx_castPull_52BFA0(id,(int)a,(int)b,(int)c,(int)record,level);
case 13: return (uint32_t)nox_xxx_castPush_52C000(id,(int)a,(int)b,(int)c,(int)record,level);
case 14: return (uint32_t)nox_xxx_castFumble_52C060(id,(int)a,(int)b,(int)c,(int*)record);
case 15: return (uint32_t)nox_xxx_castConfuse_52C1E0(id,(int)a,(int)b,(int)c,(int*)record,(char)level);
case 16: return (uint32_t)nox_xxx_castStun_52C2C0(id,(int)a,(int)b,(int)c,(int*)record,(char)level);
case 17: return (uint32_t)nox_xxx_castBurn_52C3E0(id,(int)a,(int)b,(int)c,(int)record);
case 18: return (uint32_t)nox_xxx_useShock_52C5A0(id,(int)a,(int)b,(int)c,(int*)record,level);
case 19: return (uint32_t)nox_xxx_castPoison_52C720(id,(int)a,(int)b,(int)c,(int*)record,level);
case 20: return (uint32_t)nox_xxx_castFireball_52C790(id,(int)a,(int)b,(int)c,(int)record,level);
case 21: return (uint32_t)sub_52CA80(id,(int)a,(int)b,(int)c);
case 22: return (uint32_t)sub_52CBD0(id,(int)a,(int)b,(int)c);
case 23: return (uint32_t)sub_52CCD0(id,(int)a,(int)b);
case 24: return (uint32_t)nox_xxx_castCurePoison_52CDB0(id,(int)a,(int)b,(int)c,(int*)record,level);
case 25: sub_52CE60((int)u);return 0;
case 26: return (uint32_t)nox_xxx_castLock_52CE90(id,(int)a,(int)b,(int)c);
case 27: sub_52CF90((int)u,(int)a);return 0;
case 28: sub_52D060((int)u,(int)a);return 0;
case 29: return (uint32_t)nox_xxx_castTelekinesis_52D330(id,(int)a,(int)b,(int)c,(int*)record,(char)level);
case 30: return (uint32_t)nox_xxx_castFist_52D3C0(id,(int)a,(int)b,(int)c,(int)record,level);
case 31: return (uint32_t)nox_xxx_spellCastCleansingFlame_52D5C0(id,a,b,c,record,level);
case 32: return (uint32_t)nox_xxx_castMeteorShower_52D8A0(id,(int)a,(int)b,(int)c,(int)record,level);
case 33: return (uint32_t)nox_xxx_castMeteor_52D9D0(id,(int)a,(int)b,(int)c,(int)record,level);
case 34: return (uint32_t)nox_xxx_castToxicCloud_52DB60(id,(int)a,(int)b,(int)c,(int)record);
case 35: return (uint32_t)nox_xxx_spellArachna_52DC80(id,(int)a,(int)b,(int)c,(int)record);
case 36: return (uint32_t)sub_52DD50(id,(int)a,(int)b,(int)c,record);
case 37: return (uint32_t)nox_xxx_castEquake_52DE40(id,(int)a,(int)b,(int)c,(int)record,level);
case 38: return (uint32_t)nox_xxx_equakeDamage_52DEC0((int)u,(int)a);
case 39: return (uint32_t)nox_xxx_isObjectMovable_52E020((int)u);
case 40: nox_xxx_mapPushUnitsAround_52E040(record,x,y,z,u,q,r);return 0;
case 41: nox_xxx_unitPushAroundFn_52E0E0((int)u,(int**)record);return 0;
default:return 0;}}

*/
import "C"
import (
	"bytes"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestSpellEffectGuide struct {
	Index int
	Name  string
	Size  byte
}
type PortTestSpellEffectsSpec struct {
	RecordForce              bool
	Tile                     *int
	RecordOutput, NullOutput bool
	Args                     [3]int
	Floats                   [3]uint32
	Ints                     [2]int32
	Output                   PortTestSpellLifeRecord
	ObjectTypes              map[int]string
	Guides                   []PortTestSpellEffectGuide
	Caches                   map[uintptr]uint32
	CharmAll                 bool
}
type portTestSpellEffects struct {
	output     unsafe.Pointer
	guideNames []byte
	glyphCalls []uint32
}

var spellEffectsTypeNames = []string{"ArachnaphobiaFocus", "BlueFlameCleanse", "FlameCleanse", "LargeBlueFlameCleanse", "LargeFlameCleanse", "MediumBlueFlameCleanse", "MediumFlame", "MediumFlameCleanse", "MeteorShower", "SmallBlueFlameCleanse", "SmallFlameCleanse", "TelekinesisHand", "ToxicCloud", "Fireball", "BlackFireball", "StrongFireball", "TitanFireball", "TeleportGlyph1", "TeleportGlyph2", "TeleportGlyph3", "TeleportGlyph4"}
var spellEffectsCacheOffsets = []uintptr{1570276, 1570280, 2487700, 2487704, 2487716, 2487728, 2487732, 2487736, 2487740, 2487744, 2487748, 2487752, 2487756, 2487760, 2487764, 2487768, 2487772, 2487776, 2487780, 2487784, 2487788, 2487792, 2487796, 2487800, 2487808, 2487812}

func (p *portTestShopPools) spellEffectsPrepare() func() {
	sp := p.spellLifeSpec().Effects
	if sp == nil {
		return func() {}
	}
	p.spellLifeState().effects = &portTestSpellEffects{output: p.objectiveRegion(32), guideNames: bytes.Clone(memmap.Slice(0x587000, 70500)[:164])}
	names := append(append([]string{}, spellLifecycleTypeNames...), spellEffectsTypeNames...)
	names = append(names, "NPC", "Bat", "Glyph")
	restoreTypes := p.proxy.core.PortTestSpellEffectTypes(names, []string{"Bat"}, p.proxy.combat.actor)
	restoreTiles := func() {}
	if sp.Tile != nil {
		configure, intact, restore := portTestGeneratorTileEnvironment()
		configure(*sp.Tile, false)
		restoreTiles = func() {
			if !intact() {
				panic("spell effect tile guard")
			}
			restore()
		}
	}
	fireballTable := memmap.Slice(0x587000, 258868)[:20]
	oldFireballTable := bytes.Clone(fireballTable)
	for i, name := range []string{"Fireball", "StrongFireball", "TitanFireball", "TitanFireball", "TitanFireball"} {
		*memmap.PtrPtr(0x587000, 258868+uintptr(i*4)) = p.objectiveString(name)
	}
	oldLimit := Nox_xxx_checkSummonedCreaturesLimit_500D70
	if PortTestSpellEffectsSummonLimit == nil {
		panic("spell effects requires real summon limit owner")
	}
	Nox_xxx_checkSummonedCreaturesLimit_500D70 = PortTestSpellEffectsSummonLimit
	oldGlyphDeath := Nox_xxx_dieGlyph_54DF30
	Nox_xxx_dieGlyph_54DF30 = func(u *server.Object) {
		st := p.spellLifeState().effects
		st.glyphCalls = append(st.glyphCalls, p.normalize(uint32(uintptr(u.CObj()))))
	}
	C.spellEffectsForceReset()
	oldSummon := Nox_xxx_unitDoSummonAt_5016C0
	oldPending := p.proxy.core.Objs.Pending
	Nox_xxx_unitDoSummonAt_5016C0 = func(id int, pos types.Pointf, owner *server.Object, dir server.Dir16) *server.Object {
		u := oldSummon(id, pos, owner, dir)
		p.controlsAdopt(u)
		return u
	}
	oldCharm, oldDoor, oldGlyph, oldOther := C.nox_cheat_charmall, C.dword_5d4594_2487708, C.dword_5d4594_2487712, C.dword_5d4594_2487804
	C.nox_cheat_charmall = C.int(bool2int(sp.CharmAll))
	C.dword_5d4594_2487708 = 0
	C.dword_5d4594_2487712 = 0
	C.dword_5d4594_2487804 = 0
	old := make([]uint32, len(spellEffectsCacheOffsets))
	for i, off := range spellEffectsCacheOffsets {
		old[i] = *memmap.PtrUint32(0x5d4594, off)
		*memmap.PtrUint32(0x5d4594, off) = sp.Caches[off]
	}
	guide := memmap.Slice(0x5d4594, 740076)[:41*28]
	oldGuide := bytes.Clone(guide)
	for i := 0; i <= 40; i++ {
		*memmap.PtrUint32(0x5d4594, uintptr(740080+28*i)) = 0
		*(*byte)(memmap.PtrOff(0x5d4594, uintptr(740100+28*i))) = 0
	}
	entries := sp.Guides
	if entries == nil {
		entries = []PortTestSpellEffectGuide{{Index: 1, Size: 1}}
	}
	for _, g := range entries {
		if g.Index < 1 || g.Index > 40 {
			panic("spell effect guide index")
		}
		name := g.Name
		if name == "" {
			if g.Index != 1 {
				panic("spell effect guide requires a type name")
			}
			name = "Bat"
		}
		*memmap.PtrUint32(0x5d4594, uintptr(740080+28*g.Index)) = uint32(p.proxy.core.Types.IndByID(name))
		*(*byte)(memmap.PtrOff(0x5d4594, uintptr(740100+28*g.Index))) = g.Size
	}
	return func() {
		restoreTiles()
		copy(fireballTable, oldFireballTable)
		Nox_xxx_checkSummonedCreaturesLimit_500D70 = oldLimit
		Nox_xxx_dieGlyph_54DF30 = oldGlyphDeath
		Nox_xxx_unitDoSummonAt_5016C0 = oldSummon
		p.proxy.core.Objs.Pending = oldPending
		restoreTypes()
		copy(memmap.Slice(0x587000, 70500)[:164], p.spellLifeState().effects.guideNames)
		copy(guide, oldGuide)
		for i, off := range spellEffectsCacheOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		C.nox_cheat_charmall = oldCharm
		C.dword_5d4594_2487708 = oldDoor
		C.dword_5d4594_2487712 = oldGlyph
		C.dword_5d4594_2487804 = oldOther
	}
}
func (p *portTestShopPools) spellEffectsItems() {
	sp := p.spellLifeSpec().Effects
	if sp == nil {
		return
	}
	entries := sp.Guides
	if entries == nil {
		entries = []PortTestSpellEffectGuide{{Index: 1, Size: 1}}
	}
	for _, g := range entries {
		name := g.Name
		if name == "" {
			name = "Bat"
		}
		*memmap.PtrPtr(0x587000, 70500+uintptr(4*g.Index)) = p.objectiveString(name)
	}
	for i := 0; i < 42; i++ {
		p.identify(C.spellEffectsFunction(C.int(i)), 95000+uint32(i))
	}
	for ref, name := range sp.ObjectTypes {
		u := p.temporaryRef(ref)
		u.TypeInd = uint16(p.proxy.core.Types.IndByID(name))
	}
	for i, off := range []int{484, 488} {
		ptr := *controlPtr(p.proxy.combat.actor.UpdateData, off)
		if ptr != nil {
			p.identify(ptr, 96050+uint32(i))
		}
	}
	p.identify(C.spellEffectsForcePtr(), 95042)
	if sp.RecordForce {
		*controlPtr(p.spellLifeState().record, 20) = C.spellEffectsForcePtr()
	}
	p.spellLifeFill(p.spellLifeState().effects.output, 32, sp.Output)
	if sp.RecordOutput {
		*controlPtr(p.spellLifeState().record, 0) = p.spellLifeState().effects.output
	}
}
func (p *portTestShopPools) spellEffectsAction(action PortTestShopAction) uint32 {
	ctrl := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls
	sp := ctrl.SpellLifecycle.Effects
	st := p.spellLifeState()
	u := p.temporaryRef(p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Actor)
	record, output := st.record, st.effects.output
	if ctrl.SpellLifecycle.NullRecord {
		record = nil
	}
	if sp.NullOutput {
		output = nil
	}
	q := sp.Ints[0]
	if sp.RecordForce {
		q = int32(uintptr(C.spellEffectsForcePtr()))
	}
	out := uint32(C.spellEffectsCall(C.int(action.Op-1600), C.int(ctrl.X), asObjectC(u), asObjectC(p.temporaryRef(sp.Args[0])), asObjectC(p.temporaryRef(sp.Args[1])), asObjectC(p.temporaryRef(sp.Args[2])), record, output, C.int(ctrl.SpellLifecycle.Z), C.uint32_t(sp.Floats[0]), C.uint32_t(sp.Floats[1]), C.uint32_t(sp.Floats[2]), C.int(q), C.int(sp.Ints[1])))
	p.temporary.world.objectives.attack.controls.result = uint64(out)
	return p.normalize(out)
}
func (p *portTestShopPools) spellEffectsSnapshot(out []uint32) []uint32 {
	if p.spellLifeSpec().Effects == nil {
		return out
	}
	calls := p.spellLifeState().effects.glyphCalls
	out = append(out, uint32(len(calls)))
	out = append(out, calls...)
	out = append(out, uint32(C.spellEffectsForceCount()))
	for i := 0; i < int(C.spellEffectsForceCount()); i++ {
		out = append(out, p.normalize(uint32(C.spellEffectsForceValue(C.int(i)))))
	}
	out = append(out, p.normalize(uint32(C.dword_5d4594_2487708)), p.normalize(uint32(C.dword_5d4594_2487712)), p.normalize(uint32(C.dword_5d4594_2487804)), uint32(C.nox_cheat_charmall))
	for _, off := range spellEffectsCacheOffsets {
		out = append(out, p.normalize(*memmap.PtrUint32(0x5d4594, off)))
	}
	for _, v := range unsafe.Slice((*uint32)(memmap.PtrOff(0x5d4594, 740076)), 41*7) {
		out = append(out, p.normalize(v))
	}
	return out
}

func (p *portTestShopPools) spellEffectsActive() bool {
	if p == nil || p.temporary == nil || p.temporary.world == nil || p.temporary.world.objectives == nil || p.temporary.world.objectives.attack == nil || p.temporary.world.objectives.attack.controls == nil || p.temporary.world.objectives.attack.controls.spellLifecycle == nil {
		return false
	}
	return p.spellLifeSpec().Effects != nil
}

// PortTestSpellEffectsSummonLimit binds the actual root capacity owner in this fixture.
var PortTestSpellEffectsSummonLimit func(*server.Object, int) bool
