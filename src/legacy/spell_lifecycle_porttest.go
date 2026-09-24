//go:build porttest

package legacy

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
	oldPool, oldHead := legacyGlobals.nox_alloc_magicEnt_1569668, dword_5d4594_1569672
	legacyGlobals.nox_alloc_magicEnt_1569668 = st.pool.UPtr()
	dword_5d4594_1569672 = 0
	// Only type caches belong to this fixture. The gap before Hecubah and
	// Necromancer contains retired duration allocator/list slots; accessing
	// those slots through memmap is rejected by the safe runtime.
	cacheOffsets := [...]uintptr{
		1569676, 1569680, 1569684, 1569688, 1569692, 1569696, 1569700,
		1569704, 1569708, 1569712, 1569716, 1569720, 1569740, 1569744,
	}
	var oldCaches [len(cacheOffsets)]uint32
	for i, off := range cacheOffsets {
		oldCaches[i] = *memmap.PtrUint32(0x5d4594, off)
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	for i, name := range []string{"Pixie", "MagicMissile", "SmallFist", "MediumFist", "LargeFist", "DeathBall", "Meteor"} {
		*memmap.PtrUint32(0x5d4594, cacheOffsets[i]) = uint32(p.proxy.core.Types.IndByID(name))
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
		legacyGlobals.nox_alloc_magicEnt_1569668 = oldPool
		dword_5d4594_1569672 = oldHead
		for i, v := range oldCaches {
			*memmap.PtrUint32(0x5d4594, cacheOffsets[i]) = v
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
	// Keep the 13 normalized callback-ID slots formerly registered from C.
	p.reservedFunctionIDs += 14
	st := p.spellLifeState()
	u := p.temporaryRef(p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Actor)
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
		dword_5d4594_1569672 = uint32(uintptr(st.books[0]))
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
	for ptr := unsafe.Pointer(uintptr(dword_5d4594_1569672)); ptr != nil; ptr = *controlPtr(ptr, 52) {
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
	out = append(out, p.normalize(uint32(dword_5d4594_1569672)), uint32(len(st.books)))
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
		off := uintptr(1569676 + 4*i)
		switch off {
		case 1569724, 1569728, 1569732, 1569736:
			// Keep the historical zero columns for the retired duration-state
			// gap. The old fixture cleared these words; production no longer
			// uses them. Live duration list/record state is captured above.
			out = append(out, 0)
		default:
			out = append(out, *memmap.PtrUint32(0x5d4594, off))
		}
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
