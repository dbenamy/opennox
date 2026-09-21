//go:build porttest

package legacy

/*
#include "common__gamemech__pausefx.h"
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2523804, dword_5d4594_2523780, dword_5d4594_2523776;
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

type PortTestRuntimePauseSpec struct {
	Active                                               uint32
	Paused, Provided, Cached, Missing, StateOK, BookBusy bool
	Kind                                                 int32
	Ticks                                                uint64
}

func (p *portTestShopPools) runtimePauseContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.RuntimePause
	u := p.resources.unit
	a, b, c := C.dword_5d4594_2523804, C.dword_5d4594_2523780, C.dword_5d4594_2523776
	defer func() { C.dword_5d4594_2523804, C.dword_5d4594_2523780, C.dword_5d4594_2523776 = a, b, c }()
	C.dword_5d4594_2523804, C.dword_5d4594_2523780, C.dword_5d4594_2523776 = C.uint32_t(sp.Active), 0, 0
	if sp.Cached {
		C.dword_5d4594_2523780 = C.uint32_t(uintptr(u.CObj()))
	}
	offsets := []uintptr{2523772, 2523796, 2523800}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = memmap.Uint32(0x5D4594, off)
		*memmap.PtrUint32(0x5D4594, off) = 0x12345678
	}
	book := bookWord(1047520)
	oldBook, oldHold := *book, Get_dword_5d4594_251744()
	*book = uint32(bool2int(sp.BookBusy))
	Set_dword_5d4594_251744(0)
	defer func() { *book = oldBook; Set_dword_5d4594_251744(oldHold) }()
	defer func() {
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
	}()
	ticks := memmap.PtrUint64(0x5D4594, 2523788)
	oldTicks := *ticks
	*ticks = 0xabcdef0187654321
	defer func() { *ticks = oldTicks }()
	game := noxflags.GetGame()
	noxflags.UnsetGame(noxflags.GamePause)
	if sp.Paused {
		noxflags.SetGame(noxflags.GamePause)
	}
	defer func() { noxflags.ResetGame(); noxflags.SetGame(game) }()
	oldClock, oldState, oldPause := PlatformTicks, Nox_xxx_playerSetState_4FA020, Sub_413A00
	defer func() { PlatformTicks, Nox_xxx_playerSetState_4FA020, Sub_413A00 = oldClock, oldState, oldPause }()
	PlatformTicks = func() uint64 { return sp.Ticks }
	var events []uint32
	Nox_xxx_playerSetState_4FA020 = func(actor *server.Object, state server.PlayerState) bool {
		if actor != u {
			panic("pause animation actor")
		}
		events = append(events, 1, uint32(state))
		return sp.StateOK
	}
	Sub_413A00 = func(v int) { events = append(events, 2, uint32(v)); oldPause(v) }
	var missing []string
	if sp.Missing {
		missing = []string{"LevelUp", "OblivionUp"}
	}
	restoreTypes := p.proxy.core.PortTestRewardTypes([]string{"LevelUp", "OblivionUp"}, missing, true, 0, 0)
	defer restoreTypes()
	stateByte := (*byte)(unsafe.Add(u.UpdateData, 236))
	oldByte := *stateByte
	*stateByte = 73
	defer func() { *stateByte = oldByte }()
	snapshot := func() []uint32 {
		return []uint32{uint32(C.dword_5d4594_2523804), p.normalize(uint32(C.dword_5d4594_2523780)), p.normalize(uint32(C.dword_5d4594_2523776)), memmap.Uint32(0x5D4594, 2523772), memmap.Uint32(0x5D4594, 2523796), memmap.Uint32(0x5D4594, 2523800), uint32(*ticks), uint32(*ticks >> 32), uint32(bool2int(noxflags.HasGame(noxflags.GamePause))), uint32(*stateByte), uint32(len(events)), uint32(len(p.proxy.life.created))}
	}
	before := snapshot()
	var arg *server.Object
	if sp.Provided {
		arg = u
	}
	runtimePauseStart(arg, sp.Kind)
	started := snapshot()
	blocked := sp.Active == 1 || sp.Paused
	if blocked {
		if !reflect.DeepEqual(before, started) {
			panic("blocked pause start mutated state")
		}
	} else {
		delay := uint32(0)
		if sp.Kind == 0 || sp.Kind == 1 {
			delay = 5000
		}
		expectedByte := uint32(73)
		if (sp.Provided || sp.Cached) && sp.StateOK {
			expectedByte = 4
		}
		if started[0] != 1 || started[3] != uint32(sp.Kind) || started[4] != delay || started[5] != 0 || *ticks != uint64(uint32(sp.Ticks)) || started[8] != 1 || started[9] != expectedByte {
			panic(fmt.Sprintf("pause start state: %+v / %+v", sp, started))
		}
		actor := uint32(0)
		if sp.Provided || sp.Cached {
			actor = p.normalize(uint32(uintptr(u.CObj())))
		}
		if started[1] != actor {
			panic("pause cached actor")
		}
		wantCreated := 0
		if actor != 0 && !sp.Missing && (sp.Kind == 0 || sp.Kind == 1) {
			wantCreated = 1
		}
		if int(started[11]-before[11]) != wantCreated || (started[2] != 0) != (wantCreated != 0) {
			panic("pause effect ownership")
		}
	}
	// A second start while active/paused must preserve all state and callback counts.
	runtimePauseStart(arg, sp.Kind)
	if !reflect.DeepEqual(started, snapshot()) {
		panic("repeated pause start mutated state")
	}
	runtimePauseStop()
	stopped := snapshot()
	if started[0] != 0 {
		if stopped[0] != 0 || stopped[1] != 0 || stopped[2] != 0 {
			panic("pause stop left owned state")
		}
		wantPause := started[8]
		if !sp.BookBusy {
			wantPause = 0
		}
		if stopped[8] != wantPause {
			panic(fmt.Sprintf("pause stop book gate: %+v start=%v stop=%v", sp, started, stopped))
		}
	} else if !reflect.DeepEqual(started, stopped) {
		panic("inactive pause stop mutated state")
	}
	runtimePauseStop()
	if !reflect.DeepEqual(stopped, snapshot()) {
		panic("repeated pause stop mutated state")
	}
	out := append(before, started...)
	out = append(out, stopped...)
	out = append(out, uint32(len(events)))
	return append(out, events...)
}
