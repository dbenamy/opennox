//go:build porttest

package legacy

/*
#include <stdint.h>
#include <string.h>
#include "GAME4.h"
extern void* nox_alloc_magicEnt_1569668;
extern uint32_t dword_5d4594_1569672;
unsigned short sub_4FD030(int a1, short a2);
void nox_xxx_collide_4FDF90(int a1, int a2);
int nox_xxx_spellGetPhoneme_4FE1C0(int a1, char a2);
int nox_xxx_spellByBookInsert_4FE340(int a1, int* a2, int a3, int a4, int a5);
int sub_4FEA70(int a1, float2* a2);
int nox_xxx_playerCancelSpells_4FEAE0(nox_object_t* a1p);
char* nox_xxx_netStartDurationRaySpell_4FF130(int a1);
int sub_4FF2D0(int a1, int a2);
int nox_xxx_testUnitBuffs_4FF350(nox_object_t* unit, char buff);
void nox_xxx_buffApplyTo_4FF380(nox_object_t* unit, int buff, short dur, char power);
int nox_xxx_unitGetBuffTimer_4FF550(nox_object_t* unit, int buff);
char nox_xxx_buffGetPower_4FF570(nox_object_t* unit, int buff);
void nox_xxx_unitClearBuffs_4FF580(nox_object_t* unit);
int nox_xxx_spellBuffOff_4FF5B0(nox_object_t* a1p, int a2);
static void* spellLifeFunction(int id){switch(id){
case 6:return sub_4FD030;
case 11:return nox_xxx_collide_4FDF90;
case 12:return nox_xxx_spellGetPhoneme_4FE1C0;
case 13:return nox_xxx_spellByBookInsert_4FE340;
case 16:return sub_4FEA70;
case 17:return nox_xxx_playerCancelSpells_4FEAE0;
case 20:return nox_xxx_netStartDurationRaySpell_4FF130;
case 21:return sub_4FF2D0;
case 22:return nox_xxx_testUnitBuffs_4FF350;
case 23:return nox_xxx_buffApplyTo_4FF380;
case 24:return nox_xxx_unitGetBuffTimer_4FF550;
case 25:return nox_xxx_buffGetPower_4FF570;
case 26:return nox_xxx_unitClearBuffs_4FF580;
case 27:return nox_xxx_spellBuffOff_4FF5B0;
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestSpellLifeRecord struct {
	Words map[int]uint32
	Refs  map[int]int
}
type PortTestSpellLifecycleSpec struct {
	Effects           *PortTestSpellEffectsSpec
	ClientSprite      bool
	ActorType         string
	Definitions       []server.PortTestSpellLifecycleDef
	Record            PortTestSpellLifeRecord
	NullRecord        bool
	Z                 int32
	Books             []PortTestSpellLifeRecord
	BookTree          []int
	Durations         []PortTestSpellLifeRecord
	RecordChildren    map[int]int
	TreeIndices       []int32
	TreeEdges         []map[int]int
	HalfPointerReturn bool
}
type portTestSpellLifecycle struct {
	effects   *portTestSpellEffects
	magicType uint16
	client    *client.Client
	record    unsafe.Pointer
	tree      []*server.PhonemeLeaf
	pool      *alloc.Class
	books     []unsafe.Pointer
	durations []*server.DurSpell
}

type portTestSpellLifeClient struct {
	Client
	core *client.Client
}

func (c *portTestSpellLifeClient) Cli() *client.Client { return c.core }

var spellLifecycleTypeNames = []string{"ImaginaryCaster", "Magic", "Crown", "GameBall", "Hecubah", "Necromancer", "Pixie", "MagicMissile", "SmallFist", "MediumFist", "LargeFist", "DeathBall", "Meteor"}

func (p *portTestShopPools) spellLifeSpec() *PortTestSpellLifecycleSpec {
	return p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
}
func (p *portTestShopPools) spellLifeState() *portTestSpellLifecycle {
	return p.temporary.world.objectives.attack.controls.spellLifecycle
}
func (p *portTestShopPools) spellLifePrepare() func() {
	sp := p.spellLifeSpec()
	if sp == nil {
		return func() {}
	}
	st := &portTestSpellLifecycle{magicType: uint16(p.proxy.core.Types.IndByID("Magic")), record: p.objectiveRegion(160), pool: alloc.NewClass("portSpellBook", 60, 64)}
	p.temporary.world.objectives.attack.controls.spellLifecycle = st
	oldPool, oldHead := C.nox_alloc_magicEnt_1569668, C.dword_5d4594_1569672
	C.nox_alloc_magicEnt_1569668 = st.pool.UPtr()
	C.dword_5d4594_1569672 = 0
	oldCaches := make([]uint32, 18)
	for i := range oldCaches {
		off := uintptr(1569676 + 4*i)
		oldCaches[i] = *memmap.PtrUint32(0x5d4594, off)
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	for i, name := range []string{"Pixie", "MagicMissile", "SmallFist", "MediumFist", "LargeFist", "DeathBall", "Meteor"} {
		*memmap.PtrUint32(0x5d4594, uintptr(1569676+4*i)) = uint32(p.proxy.core.Types.IndByID(name))
	}
	treeIDs := sp.TreeIndices
	if len(treeIDs) == 0 {
		treeIDs = []int32{0, 1}
	}
	for _, id := range treeIDs {
		leaf := (*server.PhonemeLeaf)(p.objectiveRegion(int(unsafe.Sizeof(server.PhonemeLeaf{}))))
		leaf.Ind = id
		st.tree = append(st.tree, leaf)
	}
	if len(sp.TreeEdges) == 0 && len(st.tree) > 1 {
		st.tree[0].Pho[0] = st.tree[1]
	}
	for i, edges := range sp.TreeEdges {
		for phon, j := range edges {
			st.tree[i].Pho[phon] = st.tree[j]
		}
	}
	defs := sp.Definitions
	if defs == nil {
		for i := 1; i <= 136; i++ {
			defs = append(defs, server.PortTestSpellLifecycleDef{Index: i, Flags: 0x1000000, Valid: true, Enabled: true, ManaCost: 10, Phonemes: []int{0}, Sounds: [3]int{100 + i, 300 + i, 500 + i}})
		}
	}
	restore := p.proxy.core.PortTestSpellLifecycle(defs, st.tree[0])
	oldClient := GetClient
	if sp.ClientSprite {
		st.client = new(client.Client)
		proxy := &portTestSpellLifeClient{core: st.client}
		GetClient = func() Client { return proxy }
	}
	restoreEffects := p.spellEffectsPrepare()
	return func() {
		restoreEffects()
		GetClient = oldClient
		restore()
		st.pool.Free()
		C.nox_alloc_magicEnt_1569668 = oldPool
		C.dword_5d4594_1569672 = oldHead
		for i, v := range oldCaches {
			*memmap.PtrUint32(0x5d4594, uintptr(1569676+4*i)) = v
		}
	}
}
func (p *portTestShopPools) spellLifeFill(ptr unsafe.Pointer, size int, r PortTestSpellLifeRecord) {
	for off, v := range r.Words {
		if off < 0 || off+4 > size || off%4 != 0 {
			panic("spell life word")
		}
		*equipmentWord(ptr, off) = v
	}
	for off, ref := range r.Refs {
		if off < 0 || off+4 > size || off%4 != 0 {
			panic("spell life ref")
		}
		*controlPtr(ptr, off) = p.temporaryRef(ref).CObj()
	}
}
func (p *portTestShopPools) spellLifeItems() {
	sp := p.spellLifeSpec()
	if sp == nil {
		return
	}
	st := p.spellLifeState()
	u := p.temporaryRef(p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Actor)
	for i := 0; i < 29; i++ {
		if f := C.spellLifeFunction(C.int(i)); f != nil {
			p.identify(f, 92000+uint32(i))
		}
	}
	if sp.ActorType != "" {
		u.TypeInd = uint16(p.proxy.core.Types.IndByID(sp.ActorType))
	}
	if st.client != nil {
		dr := (*client.Drawable)(p.objectiveRegion(int(unsafe.Sizeof(client.Drawable{}))))
		dr.NetCode32 = u.NetCode
		dr.ObjClass = u.Class()
		st.client.Objs.List1 = dr
	}
	p.spellLifeFill(st.record, 160, sp.Record)
	for i, r := range sp.Books {
		node := st.pool.NewObject()
		clear(unsafe.Slice((*byte)(node), 60))
		p.identify(node, 93000+uint32(i))
		*controlPtr(node, 4) = u.CObj()
		*controlPtr(node, 32) = st.tree[0].C()
		if i < len(sp.BookTree) {
			*controlPtr(node, 32) = st.tree[sp.BookTree[i]].C()
		}
		p.spellLifeFill(node, 60, r)
		st.books = append(st.books, node)
	}
	for i, node := range st.books {
		if i+1 < len(st.books) {
			*controlPtr(node, 52) = st.books[i+1]
		}
		if i > 0 {
			*controlPtr(node, 56) = st.books[i-1]
		}
	}
	if len(st.books) > 0 {
		C.dword_5d4594_1569672 = C.uint32_t(uintptr(st.books[0]))
	}
	for i, r := range sp.Durations {
		d := p.proxy.core.Spells.Dur.NewRaw()
		p.identify(d.C(), 94000+uint32(i))
		p.spellLifeFill(d.C(), 120, r)
		st.durations = append(st.durations, d)
	}
	for i := len(st.durations) - 1; i >= 0; i-- {
		p.proxy.core.Spells.Dur.Add(st.durations[i])
	}
	p.spellEffectsItems()
	for off, i := range sp.RecordChildren {
		*controlPtr(st.record, off) = st.durations[i].C()
	}
}
func (p *portTestShopPools) spellLifeAction(a PortTestShopAction) uint32 {
	sp := p.spellLifeSpec()
	st := p.spellLifeState()
	ctrl := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls
	u := p.temporaryRef(p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Actor)
	record := st.record
	if sp.NullRecord {
		record = nil
	}
	var result uint32
	target := p.temporaryRef(ctrl.Target)
	x, y := ctrl.X, ctrl.Y
	switch a.Op - 1500 {
	case 0:
		result = uint32(spellLifeBroadcastPhoneme(u, int8(x)))
	case 1:
		result = uint32(spellLifeReset(x, y))
	case 2:
		spellLifeCastBooks()
	case 3:
		result = uint32(spellLifeCancelDurations(x))
	case 4:
		result = uint32(spellLifeCheckMana(u, record, x))
	case 5:
		result = uint32(spellLifeSpendMana(u, x, y))
	case 6:
		result = uint32(spellLifeRefundMana(u, int16(x)))
	case 7:
		result = uint32(spellLifeCheckClass(u, x))
	case 8:
		result = uint32(spellLifeCantCast(u, x, y))
	case 9:
		result = uint32(spellLifeCaptureAllowed(x, u))
	case 10:
		result = controlRaw(spellLifeCreateFly(u, target, x))
	case 11:
		spellLifeCollide(u, target)
	case 12:
		result = uint32(spellLifePhoneme(int32(u.NetCode), int8(x)))
	case 13:
		result = uint32(spellLifeInsertBook(u, record, x, y, sp.Z))
	case 14:
		spellLifeCounterBooks(u, math.Float32frombits(uint32(x)))
	case 15:
		result = uint32(spellLifePower(x, u))
	case 16:
		result = uint32(spellLifeMoved(u, (*types.Pointf)(record)))
	case 17:
		result = uint32(spellLifeCancelPlayer(u))
	case 18:
		spellLifeCancelWand(u, target)
	case 19:
		spellLifeCancelSelected(u)
	case 20:
		result = spellLifeRayMessage((*server.DurSpell)(record))
	case 21:
		result = uint32(uintptr(unsafe.Pointer(spellLifeFindDuration(x, u))))
	case 22:
		result = uint32(bool2int(spellLifeHasBuff(u, int32(int8(x)))))
	case 23:
		spellLifeApplyBuff(u, x, int16(y), int8(sp.Z))
	case 24:
		result = uint32(spellLifeBuffTimer(u, x))
	case 25:
		result = uint32(int32(spellLifeBuffPower(u, x)))
	case 26:
		spellLifeClearBuffs(u)
	case 27:
		result = spellLifeBuffOff(u, x)
	case 28:
		spellLifeUpdateBuffs(u)
	}
	if sp.HalfPointerReturn {
		if result != uint32(uint16(uintptr(u.CObj()))) {
			panic("spell life half pointer return")
		}
		result = p.normalize(controlRaw(u))
	}
	for ptr := unsafe.Pointer(uintptr(C.dword_5d4594_1569672)); ptr != nil; ptr = *controlPtr(ptr, 52) {
		found := false
		for _, old := range st.books {
			if old == ptr {
				found = true
				break
			}
		}
		if !found {
			p.identify(ptr, 93000+uint32(len(st.books)))
			st.books = append(st.books, ptr)
		}
	}
	p.temporary.world.objectives.attack.controls.result = uint64(result)
	return p.normalize(result)
}
func (p *portTestShopPools) spellLifeSnapshot(out []uint32) []uint32 {
	st := p.spellLifeState()
	if st == nil {
		return out
	}
	out = append(out, p.normalize(uint32(C.dword_5d4594_1569672)), uint32(len(st.books)))
	for _, node := range st.books {
		for _, v := range unsafe.Slice((*uint32)(node), 15) {
			out = append(out, p.normalize(v))
		}
	}
	out = append(out, p.normalize(uint32(uintptr(p.proxy.core.Spells.Dur.List.C()))), uint32(len(st.durations)))
	for _, d := range st.durations {
		for _, v := range unsafe.Slice((*uint32)(d.C()), 30) {
			out = append(out, p.normalize(v))
		}
	}
	for i := 0; i < 18; i++ {
		out = append(out, *memmap.PtrUint32(0x5d4594, uintptr(1569676+4*i)))
	}
	return p.spellEffectsSnapshot(out)
}

// PortTestSpellLifePlayerSpell binds the actual root player-casting owner to the
// isolated fixture server. It does not supply a synthetic cast result.
var PortTestSpellLifePlayerSpell func(*server.Object)

func (s *portTestRoamOwnerServer) PlayerSpell(u *server.Object) {
	if PortTestSpellLifePlayerSpell != nil {
		s.trace = append(s.trace, 900, u.NetCode)
		PortTestSpellLifePlayerSpell(u)
		return
	}
	s.portTestRandomServer.Server.PlayerSpell(u)
}
