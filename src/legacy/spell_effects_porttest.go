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
int nox_xxx_summonStart_500DA0(int a1);
int nox_xxx_summonFinish_5010D0(int a1);
void nox_xxx_summonCancel_5011C0(int a1);
int nox_xxx_charmCreature1_5011F0(int* a1);
int nox_xxx_charmCreatureFinish_5013E0(int* a1);
int nox_xxx_charmCreature2_501690(int a1);
static void* spellEffectsFunction(int op){switch(op){
case 1:return nox_xxx_summonStart_500DA0;
case 3:return nox_xxx_summonFinish_5010D0;
case 4:return nox_xxx_summonCancel_5011C0;
case 5:return nox_xxx_charmCreature1_5011F0;
case 6:return nox_xxx_charmCreatureFinish_5013E0;
case 7:return nox_xxx_charmCreature2_501690;
default:return 0;}}

*/
import "C"
import (
	"bytes"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestSpellEffectGuide struct {
	Index int
	Name  string
	Size  byte
}
type PortTestSpellEffectsSpec struct {
	Sustained                *PortTestSustainedSpellsSpec
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
	if sp.Sustained != nil {
		names = append(names, sustainedTypeNames...)
	}
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
	restoreSustained := p.sustainedPrepare()
	return func() {
		restoreSustained()
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
		if f := C.spellEffectsFunction(C.int(i)); f != nil {
			p.identify(f, 95000+uint32(i))
		}
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
	p.sustainedItems()
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
	out := spellEffectsInvoke(int32(action.Op-1600), ctrl.X, u, p.temporaryRef(sp.Args[0]), p.temporaryRef(sp.Args[1]), p.temporaryRef(sp.Args[2]), record, output, ctrl.SpellLifecycle.Z, math.Float32frombits(sp.Floats[0]), math.Float32frombits(sp.Floats[1]), math.Float32frombits(sp.Floats[2]), q, sp.Ints[1])
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
	return p.sustainedSnapshot(out)
}

func (p *portTestShopPools) spellEffectsActive() bool {
	if p == nil || p.temporary == nil || p.temporary.world == nil || p.temporary.world.objectives == nil || p.temporary.world.objectives.attack == nil || p.temporary.world.objectives.attack.controls == nil || p.temporary.world.objectives.attack.controls.spellLifecycle == nil {
		return false
	}
	return p.spellLifeSpec().Effects != nil
}

// PortTestSpellEffectsSummonLimit binds the actual root capacity owner in this fixture.
var PortTestSpellEffectsSummonLimit func(*server.Object, int) bool

func spellEffectsInvoke(op, id int32, u, a, b, c *server.Object, record, output unsafe.Pointer, level int32, x, y, z float32, q, r int32) uint32 {
	switch op {
	case 0:
		return uint32(spellEffectSummonCost(id, u))
	case 1:
		return uint32(C.nox_xxx_summonStart_500DA0(C.int(uintptr(record))))
	case 2:
		return uint32(spellEffectSummonPosition(record, output))
	case 3:
		return uint32(C.nox_xxx_summonFinish_5010D0(C.int(uintptr(record))))
	case 4:
		C.nox_xxx_summonCancel_5011C0(C.int(uintptr(record)))
		return 0
	case 5:
		return uint32(C.nox_xxx_charmCreature1_5011F0((*C.int)(record)))
	case 6:
		return uint32(C.nox_xxx_charmCreatureFinish_5013E0((*C.int)(record)))
	case 7:
		return uint32(C.nox_xxx_charmCreature2_501690(C.int(uintptr(record))))
	case 8:
		spellEffectBanish(u)
		return 0
	case 9:
		return uint32(spellEffectInversion(id, a, b, c, record, level))
	case 10:
		return uint32(spellEffectRestoreHealth(id, a, b, c, record, level))
	case 11:
		return uint32(spellEffectRestoreMana(id, a, b, c, record, level))
	case 12:
		return uint32(spellEffectPull(id, a, b, c, record, level))
	case 13:
		return uint32(spellEffectPush(id, a, b, c, record, level))
	case 14:
		return uint32(spellEffectFumble(id, a, b, c, record, level))
	case 15:
		return uint32(spellEffectConfuse(id, a, b, c, record, level))
	case 16:
		return uint32(spellEffectStun(id, a, b, c, record, level))
	case 17:
		return uint32(spellEffectBurn(id, a, b, c, record, level))
	case 18:
		return uint32(spellEffectShock(id, a, b, c, record, level))
	case 19:
		return uint32(spellEffectPoison(id, a, b, c, record, level))
	case 20:
		return uint32(spellEffectFireball(id, a, b, c, record, level))
	case 21:
		return uint32(spellEffectPortal(id, a, b, c, record, level))
	case 22:
		return uint32(spellEffectNamedPortal(id, a, b, c, record, level))
	case 23:
		return uint32(spellEffectDetonateGlyph(id, a, b, c, record, level))
	case 24:
		return uint32(spellEffectCurePoison(id, a, b, c, record, level))
	case 25:
		spellEffectDoorLink(u)
		return 0
	case 26:
		return uint32(spellEffectLock(id, a, b, c, record, level))
	case 27:
		spellEffectDoorCandidate(u, a)
		return 0
	case 28:
		spellEffectDoorPropagate(u, a)
		return 0
	case 29:
		return uint32(spellEffectTelekinesis(id, a, b, c, record, level))
	case 30:
		return uint32(spellEffectFist(id, a, b, c, record, level))
	case 31:
		return uint32(spellEffectCleansingFlame(id, a, b, c, record, level))
	case 32:
		return uint32(spellEffectMeteorShower(id, a, b, c, record, level))
	case 33:
		return uint32(spellEffectMeteor(id, a, b, c, record, level))
	case 34:
		return uint32(spellEffectToxicCloud(id, a, b, c, record, level))
	case 35:
		return uint32(spellEffectArachna(id, a, b, c, record, level))
	case 36:
		return uint32(spellEffectLesserHeal(id, a, b, c, record, level))
	case 37:
		return uint32(spellEffectQuake(id, a, b, c, record, level))
	case 38:
		return uint32(int32(spellEffectQuakeDamage(u, a)))
	case 39:
		return spellEffectMovable(u)
	case 40:
		spellEffectPushAround(spellEffectPos(record, 0), x, y, z, u, unsafe.Pointer(uintptr(uint32(q))), unsafe.Pointer(uintptr(uint32(r))))
		return 0
	case 41:
		spellEffectPushUnit(u, record)
		return 0
	default:
		panic("spell effects fixture operation")
	}
}
