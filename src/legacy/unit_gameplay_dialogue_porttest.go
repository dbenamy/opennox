//go:build porttest

package legacy

/*
#include "GAME4.h"
#include "GAME5.h"
void nox_xxx_scriptDialog_548D30(nox_object_t* a1, char a2);
*/
import "C"

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"unsafe"
)

type PortTestUnitDialogueSpec struct {
	Finish             bool
	Gate               int
	Response, Kind     byte
	Frozen, FreezeLock bool
}

func (p *portTestShopPools) unitOrderContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls
	pl := p.proxy.core.Players.ByInd(ntype.PlayerInd(sp.X))
	word := (*uint32)(unsafe.Add(unsafe.Pointer(pl), 3648))
	*word = 0xabcdef01
	result := C.nox_xxx_orderUnitLocal_500C70(C.int(sp.X), C.int(sp.Y))
	if *word != uint32(sp.Y) {
		panic("local order full-width state")
	}
	return []uint32{*word, uint32(result)}
}
func (p *portTestShopPools) unitDialogueContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.UnitDialogue
	u, m := p.resources.unit, p.proxy.combat.actor
	oldPlayerClass, oldMonsterClass := u.ObjClass, m.ObjClass
	defer func() { u.ObjClass, m.ObjClass = oldPlayerClass, oldMonsterClass }()
	vm := &p.proxy.core.NoxScriptVM
	vm.Init(p.proxy.core)
	var observed []uint32
	callback := func(which uint32) func() {
		return func() {
			partner := *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 284))
			observed = append(observed, which, uint32(bool2int(vm.Caller() == u)), uint32(bool2int(vm.Trigger() == m)), uint32(bool2int(partner == nil)), uint32(*(*byte)(unsafe.Add(m.UpdateData, 2105))), uint32(u.ObjFlags&2))
		}
	}
	start, end := vm.AsFuncIndex("port-dialogue-start", callback(1)), vm.AsFuncIndex("port-dialogue-end", callback(2))
	a, b := (*int32)(unsafe.Add(m.UpdateData, 2096)), (*int32)(unsafe.Add(m.UpdateData, 2100))
	*a, *b = int32(start), int32(end)
	m.ObjClass = object.ClassMonster
	m.ObjFlags = 0
	m.Field5 = 0x10
	if sp.Frozen {
		u.ObjFlags |= 2
	} else {
		u.ObjFlags &^= 2
	}
	lock := memmap.PtrUint32(0x5D4594, 1567712)
	oldLock := *lock
	defer func() { *lock = oldLock }()
	*lock = uint32(bool2int(sp.FreezeLock))
	*(*byte)(unsafe.Add(m.UpdateData, 2104)) = sp.Kind
	*(*byte)(unsafe.Add(m.UpdateData, 2105)) = 0xa5
	partner := (*unsafe.Pointer)(unsafe.Add(u.UpdateData, 284))
	*partner = m.CObj()
	player, monster := u, m
	switch sp.Gate {
	case 1:
		*a = -1
	case 2:
		*b = -1
	case 3:
		u.ObjClass = object.ClassMonster
	case 4:
		m.ObjClass = object.ClassPlayer
	case 5:
		u.ObjFlags |= 0x20
	case 6:
		u.ObjFlags |= 0x8000
	case 7:
		m.Field5 = 0
	case 8:
		player = nil
	case 9:
		monster = nil
		*partner = nil
	}
	wanted := 0
	if sp.Finish {
		C.nox_xxx_scriptDialog_548D30(asObjectC(u), C.char(sp.Response))
		if sp.Gate != 1 && sp.Gate != 2 && sp.Gate != 9 {
			wanted = 2
		}
	} else {
		C.nox_xxx_script_forcedialog_548CD0(asObjectC(player), asObjectC(monster))
		if sp.Gate == 0 {
			wanted = 1
		}
	}
	if wanted == 0 && len(observed) != 0 || wanted != 0 && (len(observed) != 6 || observed[0] != uint32(wanted) || observed[1] != 1 || observed[2] != 1) {
		panic(fmt.Sprintf("dialogue callback %v want %d", observed, wanted))
	}
	response := *(*byte)(unsafe.Add(m.UpdateData, 2105))
	if wanted == 2 {
		wantResponse := byte(0)
		if sp.Kind == 1 {
			wantResponse = sp.Response
		}
		if *partner != nil || response != wantResponse || observed[3] != 1 || observed[4] != uint32(wantResponse) {
			panic("dialogue callback state/order")
		}
	} else if response != 0xa5 {
		panic("inactive dialogue response mutation")
	}
	if sp.Finish {
		wantFrozen := sp.Frozen && sp.FreezeLock
		if (u.ObjFlags&2 != 0) != wantFrozen {
			panic("dialogue unfreeze")
		}
	}
	out := []uint32{uint32(u.ObjFlags), uint32(bool2int(*partner == nil)), uint32(response), uint32(len(observed))}
	return append(out, observed...)
}
