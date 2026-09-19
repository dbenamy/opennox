//go:build porttest

package opennox

import (
	"encoding/binary"
	"testing"

	"github.com/opennox/libs/things"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestClientInteractionEscapeBookAddition(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	quick, restoreQuick := legacy.PortTestQuickbarWords()
	defer restoreQuick()
	capture := serverConfigOwnBytes(t, 0x5D4594, 1047928, 4)
	cursor := serverConfigOwnBytes(t, 0x5D4594, 1096672, 4)
	pause, restorePause := legacy.PortTestBookPauseOwner()
	defer restorePause()
	oldTicks, oldFrame, oldLoading := nox_gameTicks_371764, nox_gameFrame_371772, dword_5d4594_1563080
	defer func() {
		nox_gameTicks_371764, nox_gameFrame_371772, dword_5d4594_1563080 = oldTicks, oldFrame, oldLoading
	}()
	dword_5d4594_1563080 = false
	oldFade := legacy.Get_nox_gameDisableMapDraw_5d4594_2650672()
	defer legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(oldFade)
	type row struct {
		Width, Kind   uint32
		Gate          int
		Pending, Slot uint32
		Particles     [][]uint32
		RNG           [2]int
	}
	var captured []row
	for _, width := range []uint32{640, 750, 1000} {
		for _, kind := range []uint32{2, 3, 4} {
			for gate := 0; gate < 5; gate++ {
				func() {
					prepareSpellbookAddition(t, o, bar, width)
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
					pause()
					noxflags.SetGame(noxflags.GameModeCoop | noxflags.GamePause)
					o.c.Inp.Tick()
					particles, freeParticles := legacy.PortTestEffectsScreenParticles(512)
					defer freeParticles()
					o.bookCall("nox_xxx_bookFillAll_45D570", kind, 1)
					if *o.words["dword_5d4594_1047520"] != 1 || bar[0] != 0 || len(particles()) != 0 {
						t.Fatal("pending addition setup")
					}
					chat := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 30, 20, nil)
					defer chat.Destroy()
					edit := o.c.GUI.NewWindowRaw(chat, 8, 0, 0, 30, 20, nil)
					chat.Hide()
					*words["dword_5d4594_1064856"] = uint32(uintptr(chat.C()))
					*words["dword_5d4594_1064860"] = uint32(uintptr(edit.C()))
					*words["dword_5d4594_1064868"] = 0
					*quick["dword_5d4594_1049532"], *quick["dword_5d4594_1047932"] = 0, 0
					clear(capture)
					clear(cursor)
					legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(0)
					switch gate {
					case 1:
						chat.Show()
						*words["dword_5d4594_1064868"] = 1
					case 2:
						*quick["dword_5d4594_1049532"] = uint32(uintptr(chat.C()))
						chat.Capture(true)
					case 3:
						binary.LittleEndian.PutUint32(cursor, 1)
					case 4:
						legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(1)
					}
					interactionCall("nox_xxx_consoleEsc_49B7A0")
					pending, slot, count := uint32(1), uint32(0), 0
					if gate == 0 {
						pending, slot, count = 0, 1, 50
					}
					if *o.words["dword_5d4594_1047520"] != pending || bar[0] != slot || len(particles()) != count {
						t.Fatal("Escape addition priority/effects", width, kind, gate, *o.words["dword_5d4594_1047520"], bar[0], len(particles()))
					}
					captured = append(captured, row{width, kind, gate, *o.words["dword_5d4594_1047520"], bar[0], particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}})
				}()
			}
		}
	}
	interactionCapture(t, "escape-book-addition", captured)
}
