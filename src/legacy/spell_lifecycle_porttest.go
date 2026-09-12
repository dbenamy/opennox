//go:build porttest

package legacy

/*
#include <stdint.h>
#include <string.h>
#include "GAME4.h"
extern void* nox_alloc_magicEnt_1569668;
extern uint32_t dword_5d4594_1569672;
int sub_4FC960(int a1, char a2);
int nox_xxx_Fn_4FCAC0(int a1, int a2);
void nox_xxx_spellCastByBook_4FCB80();
int sub_4FCEB0(int a1);
int nox_xxx_spellCheckSmth_4FCEF0(int a1, int* a2, int a3);
int sub_4FCF90(nox_object_t* a1p, int a2, int a3);
unsigned short sub_4FD030(int a1, short a2);
int sub_4FD0E0(nox_object_t* a1p, int a2);
int nox_xxx_checkPlrCantCastSpell_4FD150(nox_object_t* a1p, int a2, int a3);
int nox_xxx_gameCaptureMagic_4FDC10(int a1, nox_object_t* a2p);
uint32_t* nox_xxx_createSpellFly_4FDDA0(nox_object_t* a1p, nox_object_t* a2p, int a3);
void nox_xxx_collide_4FDF90(int a1, int a2);
int nox_xxx_spellGetPhoneme_4FE1C0(int a1, char a2);
int nox_xxx_spellByBookInsert_4FE340(int a1, int* a2, int a3, int a4, int a5);
void nox_xxx_spell_4FE680(nox_object_t* a1p, float a2);
int nox_xxx_spellGetPower_4FE7B0(int a1, nox_object_t* a2p);
int sub_4FEA70(int a1, float2* a2);
int nox_xxx_playerCancelSpells_4FEAE0(nox_object_t* a1p);
void sub_4FEB60(int a1, int a2);
void nox_xxx_cancelAllSpells_4FEE90(nox_object_t* a1p);
char* nox_xxx_netStartDurationRaySpell_4FF130(int a1);
int sub_4FF2D0(int a1, int a2);
int nox_xxx_testUnitBuffs_4FF350(nox_object_t* unit, char buff);
void nox_xxx_buffApplyTo_4FF380(nox_object_t* unit, int buff, short dur, char power);
int nox_xxx_unitGetBuffTimer_4FF550(nox_object_t* unit, int buff);
char nox_xxx_buffGetPower_4FF570(nox_object_t* unit, int buff);
void nox_xxx_unitClearBuffs_4FF580(nox_object_t* unit);
int nox_xxx_spellBuffOff_4FF5B0(nox_object_t* a1p, int a2);
void nox_xxx_updateUnitBuffs_4FF620(nox_object_t* a1p);
static void* spellLifeFunction(int id){switch(id){
case 0:return sub_4FC960;
case 1:return nox_xxx_Fn_4FCAC0;
case 2:return nox_xxx_spellCastByBook_4FCB80;
case 3:return sub_4FCEB0;
case 4:return nox_xxx_spellCheckSmth_4FCEF0;
case 5:return sub_4FCF90;
case 6:return sub_4FD030;
case 7:return sub_4FD0E0;
case 8:return nox_xxx_checkPlrCantCastSpell_4FD150;
case 9:return nox_xxx_gameCaptureMagic_4FDC10;
case 10:return nox_xxx_createSpellFly_4FDDA0;
case 11:return nox_xxx_collide_4FDF90;
case 12:return nox_xxx_spellGetPhoneme_4FE1C0;
case 13:return nox_xxx_spellByBookInsert_4FE340;
case 14:return nox_xxx_spell_4FE680;
case 15:return nox_xxx_spellGetPower_4FE7B0;
case 16:return sub_4FEA70;
case 17:return nox_xxx_playerCancelSpells_4FEAE0;
case 18:return sub_4FEB60;
case 19:return nox_xxx_cancelAllSpells_4FEE90;
case 20:return nox_xxx_netStartDurationRaySpell_4FF130;
case 21:return sub_4FF2D0;
case 22:return nox_xxx_testUnitBuffs_4FF350;
case 23:return nox_xxx_buffApplyTo_4FF380;
case 24:return nox_xxx_unitGetBuffTimer_4FF550;
case 25:return nox_xxx_buffGetPower_4FF570;
case 26:return nox_xxx_unitClearBuffs_4FF580;
case 27:return nox_xxx_spellBuffOff_4FF5B0;
case 28:return nox_xxx_updateUnitBuffs_4FF620;
default:return 0;}}
static uint32_t spellLifeCall(int id,nox_object_t* u,nox_object_t* t,int x,int y,int z,void* record){float value;memcpy(&value,&x,4);switch(id){
case 0: return (uint32_t)sub_4FC960((int)u,(char)x);
case 1: return (uint32_t)nox_xxx_Fn_4FCAC0(x,y);
case 2: nox_xxx_spellCastByBook_4FCB80();return 0;
case 3: return (uint32_t)sub_4FCEB0(x);
case 4: return (uint32_t)nox_xxx_spellCheckSmth_4FCEF0((int)u,(int*)record,x);
case 5: return (uint32_t)sub_4FCF90(u,x,y);
case 6: return (uint32_t)sub_4FD030((int)u,(short)x);
case 7: return (uint32_t)sub_4FD0E0(u,x);
case 8: return (uint32_t)nox_xxx_checkPlrCantCastSpell_4FD150(u,x,y);
case 9: return (uint32_t)nox_xxx_gameCaptureMagic_4FDC10(x,u);
case 10: return (uint32_t)nox_xxx_createSpellFly_4FDDA0(u,t,x);
case 11: nox_xxx_collide_4FDF90((int)u,(int)t);return 0;
case 12: return (uint32_t)nox_xxx_spellGetPhoneme_4FE1C0(u->net_code,(char)x);
case 13: return (uint32_t)nox_xxx_spellByBookInsert_4FE340((int)u,(int*)record,x,y,z);
case 14: nox_xxx_spell_4FE680(u,value);return 0;
case 15: return (uint32_t)nox_xxx_spellGetPower_4FE7B0(x,u);
case 16: return (uint32_t)sub_4FEA70((int)u,(float2*)record);
case 17: return (uint32_t)nox_xxx_playerCancelSpells_4FEAE0(u);
case 18: sub_4FEB60((int)u,(int)t);return 0;
case 19: nox_xxx_cancelAllSpells_4FEE90(u);return 0;
case 20: return (uint32_t)nox_xxx_netStartDurationRaySpell_4FF130((int)record);
case 21: return (uint32_t)sub_4FF2D0(x,(int)u);
case 22: return (uint32_t)nox_xxx_testUnitBuffs_4FF350(u,(char)x);
case 23: nox_xxx_buffApplyTo_4FF380(u,x,(short)y,(char)z);return 0;
case 24: return (uint32_t)nox_xxx_unitGetBuffTimer_4FF550(u,x);
case 25: return (uint32_t)nox_xxx_buffGetPower_4FF570(u,x);
case 26: nox_xxx_unitClearBuffs_4FF580(u);return 0;
case 27: return (uint32_t)nox_xxx_spellBuffOff_4FF5B0(u,x);
case 28: nox_xxx_updateUnitBuffs_4FF620(u);return 0;
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestSpellLifeRecord struct {
	Words map[int]uint32
	Refs  map[int]int
}
type PortTestSpellLifecycleSpec struct {
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
	return func() {
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
		p.identify(C.spellLifeFunction(C.int(i)), 92000+uint32(i))
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
	result := uint32(C.spellLifeCall(C.int(a.Op-1500), asObjectC(u), asObjectC(p.temporaryRef(ctrl.Target)), C.int(ctrl.X), C.int(ctrl.Y), C.int(sp.Z), record))
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
	return out
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
