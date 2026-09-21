//go:build porttest

package legacy

/*
char playerDropATrap(int playerObj);
*/
import "C"

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
)

type PortTestExtensionSpec struct {
	Drop          bool
	Nil           bool
	Status, State byte
	Direction     int8
	Current       int // -1 means no active weapon.
	Mode          int
	Order         []int
	Want          int // First eligible item; -1 means none.
}

// Exercise the extension entrypoints over the existing real inventory owner.
// Keep actual equip/dequip and drop dispatch; observe only the class policy hook.
func (p *portTestShopPools) extensionContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.Extension
	u := p.resources.unit
	pl := u.UpdateDataPlayer().Player
	*(*byte)(unsafe.Add(pl.C(), 3680)) = sp.Status
	*(*byte)(unsafe.Add(u.UpdateData, 88)) = sp.State
	*(*float32)(unsafe.Add(pl.C(), 3632)) = -12.5
	*(*float32)(unsafe.Add(pl.C(), 3636)) = 123.25
	u.InvFirstItem = nil
	var prev *server.Object
	for _, i := range sp.Order {
		it := p.items[i].u
		it.InvNextItem, it.Field125, it.InvHolder = nil, prev, u
		if prev == nil {
			u.InvFirstItem = it
		} else {
			prev.InvNextItem = it
		}
		prev = it
	}
	oldClass := Nox_xxx_playerClassCanUseItem_57B3D0
	defer func() { Nox_xxx_playerClassCanUseItem_57B3D0 = oldClass }()
	var calls []uint32
	Nox_xxx_playerClassCanUseItem_57B3D0 = func(it *server.Object, cl player.Class) bool {
		if cl != pl.PlayerClass() {
			panic("extension class argument")
		}
		index := -1
		for i := range p.items {
			if p.items[i].u == it {
				index = i
			}
		}
		if index < 0 {
			panic("extension unknown item")
		}
		calls = append(calls, uint32(index))
		return sp.Mode != 2 || index != 1
	}
	before := portTestInventoryDropCalls()
	wantY := float32(123.25)
	if sp.Drop {
		u.PosVec = types.Pointf{X: -12.5, Y: 100}
		if sp.Mode == 1 {
			u.PosVec.Y = 23.25
			wantY = 98.25
		}
	}
	var got int
	if sp.Drop {
		target := u
		if sp.Nil {
			target = nil
		}
		got = int(C.playerDropATrap(C.int(uintptr(target.CObj()))))
		want := 0
		if sp.Want >= 0 {
			want = 1
		}
		if got != want {
			panic(fmt.Sprintf("extension trap return %d want %d", got, want))
		}
		after := portTestInventoryDropCalls()
		if sp.Want < 0 {
			if len(after) != len(before) {
				panic("extension unexpected drop")
			}
		} else {
			if len(after) != len(before)+6 {
				panic("extension missing drop")
			}
			call := after[len(before):]
			if call[0] != 2 || call[1] != uint32(uintptr(u.CObj())) || call[2] != uint32(uintptr(p.items[sp.Want].u.CObj())) || call[3] != math.Float32bits(-12.5) || call[4] != math.Float32bits(wantY) {
				panic("extension drop arguments")
			}
		}
	} else {
		var current *server.Object
		if sp.Current >= 0 {
			current = p.items[sp.Current].u
		}
		*equipmentWord(u.UpdateData, 104) = uint32(uintptr(current.CObj()))
		for i := range p.items {
			p.items[i].u.ObjFlags &^= object.Flags(0x100)
		}
		if current != nil {
			current.ObjFlags |= 0x100
		}
		// Modes 4/5 deliberately make the first eligible candidate fail equip or
		// the current weapon fail dequip. Neither may fall through to a later item.
		if sp.Mode == 4 && sp.Want >= 0 {
			p.items[sp.Want].u.ObjFlags |= 0x100
		}
		if sp.Mode == 5 && current != nil {
			current.ObjClass = object.ClassSimple
		}
		got = Mix_MouseKeyboardWeaponRoll(u, sp.Direction)
		wantResult := 0
		next := current
		if sp.Want >= 0 {
			if sp.Mode == 5 && current != nil {
				next = current
			} else if sp.Mode == 4 {
				next = nil
			} else {
				wantResult = 1
				next = p.items[sp.Want].u
			}
		}
		if got != wantResult || *equipmentWord(u.UpdateData, 104) != uint32(uintptr(next.CObj())) {
			panic(fmt.Sprintf("extension roll current=%d mode=%d direction=%d candidate=%d return=%d/%d", sp.Current, sp.Mode, sp.Direction, sp.Want, got, wantResult))
		}
		// Every selected candidate must have passed the class hook; blocked states
		// and empty searches must not call it.
		if sp.Status&3 != 0 || sp.State == 1 {
			if len(calls) != 0 {
				panic("extension blocked class query")
			}
		}
		if wantResult != 0 && (len(calls) == 0 || calls[len(calls)-1] != uint32(sp.Want)) {
			panic("extension wrong equipped candidate")
		}
	}
	return append([]uint32{uint32(got)}, calls...)
}
