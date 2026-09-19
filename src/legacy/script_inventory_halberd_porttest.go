//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestScriptHalberdSpec struct {
	Index      uint32
	Order      []int
	WantDelete int
	Equipped   bool
}
type portTestInventoryScriptVM struct {
	portTestCombatScriptBase
	vm *server.NoxScriptVM
}

func (s *portTestInventoryScriptVM) PopI32() int32    { return s.vm.PopI32() }
func (s *portTestInventoryScriptVM) PopU32() uint32   { return s.vm.PopU32() }
func (s *portTestInventoryScriptVM) PushI32(v int32)  { s.vm.PushI32(v) }
func (s *portTestInventoryScriptVM) PushU32(v uint32) { s.vm.PushU32(v) }
func (p *portTestShopPools) scriptHalberdContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptHalberd
	host := p.proxy.core.Players.ByInd(31).PlayerUnit
	host.InvFirstItem = nil
	var prev *server.Object
	for _, i := range sp.Order {
		u := p.items[i].u
		u.InvHolder = host
		u.Field125 = prev
		u.InvNextItem = nil
		if prev == nil {
			host.InvFirstItem = u
		} else {
			prev.InvNextItem = u
		}
		prev = u
	}
	names := []string{"OblivionHalberd", "OblivionHeart", "OblivionWierdling", "OblivionOrb"}
	for i, name := range names {
		if alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 247336+4*uintptr(i)))) != name {
			panic("shipped halberd table mismatch")
		}
	}
	vm := &p.proxy.core.NoxScriptVM
	vm.Init(p.proxy.core)
	old := p.proxy.combat.portTestCombatScriptBase
	p.proxy.combat.portTestCombatScriptBase = &portTestInventoryScriptVM{portTestCombatScriptBase: old, vm: vm}
	defer func() { p.proxy.combat.portTestCombatScriptBase = old }()
	vm.PushU32(0x13579bdf)
	vm.PushU32(sp.Index)
	before := len(p.proxy.life.created)
	st := p.temporary.world.objectives.attack.controls
	st.scriptDeletes = nil
	r, ok := CallScriptBuiltin(asm.BuiltinSetHalberd)
	if !ok || r != 0 || vm.PopU32() != 0x13579bdf {
		panic("halberd VM result/stack")
	}
	if len(p.proxy.life.created) != before+1 {
		panic("halberd creation count")
	}
	it := p.proxy.life.created[before]
	calls := st.scriptInit()
	if len(calls) != 1 || calls[0][0] != uintptr(it.CObj()) || calls[0][1] != 0 {
		panic("halberd initializer arguments/count")
	}
	if it.TypeInd != uint16(p.proxy.core.Types.IndByID(names[sp.Index])) || it.InvHolder != host {
		panic("halberd selected item/holder")
	}
	if (it.ObjFlags&0x100 != 0) != sp.Equipped {
		panic("halberd replacement equipment")
	}
	if it.ObjFlags&0x80000 != 0 || *equipmentWord(it.UpdateData, 4)&1 == 0 {
		panic("halberd respawn markers")
	}
	wantCount := 0
	if sp.WantDelete >= 0 {
		wantCount = 1
	}
	if len(st.scriptDeletes) != wantCount {
		panic("halberd deletion count")
	}
	if wantCount != 0 && st.scriptDeletes[0] != p.items[sp.WantDelete].u {
		panic("halberd first-match deletion")
	}

	return []uint32{sp.Index, uint32(it.TypeInd), uint32(it.ObjFlags)}
}

func (p *portTestShopPools) scriptInventoryPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls
	if sp.ScriptHalberd == nil {
		return func() {}
	}
	init, readInit := PortTestSessionEntryInitCallback()
	p.temporary.world.objectives.attack.controls.scriptInit = readInit
	restoreTypes := p.proxy.core.PortTestScriptInventoryTypes(init)
	oldPlace := Nox_xxx_inventoryServPlace_4F36F0
	Nox_xxx_inventoryServPlace_4F36F0 = func(u, it *server.Object, a, b int) bool {
		if a != 1 || b != 1 {
			panic("halberd placement arguments")
		}
		if !oldPlace(u, it, a, b) {
			panic("halberd fixture placement rejected")
		}
		// Reuse the real insertion implementation at the fixture's pickup boundary.
		inventoryInsert(u, it, 1)
		return true
	}
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 247336)), 80)
	old := bytes.Clone(table)
	copy(table[16:], blobdata.PortTestScriptInventoryNames())
	for i, off := range []uintptr{247352, 247368, 247384, 247404} {
		*memmap.PtrPtr(0x587000, 247336+4*uintptr(i)) = memmap.PtrOff(0x587000, off)
	}
	return func() { Nox_xxx_inventoryServPlace_4F36F0 = oldPlace; copy(table, old); restoreTypes() }
}
func (p *portTestShopPools) scriptInventoryObserveDelete(u *server.Object) {
	if p.temporary == nil || p.temporary.world == nil || p.temporary.world.objectives == nil || p.temporary.world.objectives.attack == nil || p.temporary.world.objectives.attack.controls == nil {
		return
	}
	st := p.temporary.world.objectives.attack.controls
	if p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptHalberd != nil {
		st.scriptDeletes = append(st.scriptDeletes, u)
	}
}
