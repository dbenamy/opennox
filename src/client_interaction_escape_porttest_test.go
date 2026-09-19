//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"slices"
	"testing"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientInteractionEscapeOrder(t *testing.T)   { interactionEscapeOrder(t, false) }
func TestClientInteractionChatEscapeKey(t *testing.T) { interactionEscapeOrder(t, true) }
func interactionEscapeOrder(t *testing.T, key bool) {
	o := newMeterOwner(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	inv, restoreInv := legacy.PortTestInventoryWindowWords()
	defer restoreInv()
	book, restoreBook := legacy.PortTestBookWords()
	defer restoreBook()
	quick, restoreQuick := legacy.PortTestQuickbarWords()
	defer restoreQuick()
	sessions, restoreSessions := legacy.PortTestSessionDialogWords()
	defer restoreSessions()
	vote, restoreVote := legacy.PortTestVoteGUIOwner()
	defer restoreVote()
	options, restoreOptions := legacy.PortTestOptionsWords()
	defer restoreOptions()
	bindings, restoreBindings := legacy.PortTestBindingWords()
	defer restoreBindings()
	serverOptions, restoreServerOptions := legacy.PortTestServerOptionsWords()
	defer restoreServerOptions()
	_ = options
	_ = bindings
	*serverOptions["root"] = 0
	for _, p := range quick {
		*p = 0
	}
	for _, off := range []uintptr{1047928, 1049476, 1049848, 1049868} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
	}
	cursor := serverConfigOwnBytes(t, 0x5D4594, 1096672, 4)
	choice := serverConfigOwnBytes(t, 0x5D4594, 1123516, 4)
	oldFlag := legacy.Get_dword_5d4594_251744()
	legacy.Set_dword_5d4594_251744(1)
	defer legacy.Set_dword_5d4594_251744(oldFlag)
	oldFade := legacy.Get_nox_gameDisableMapDraw_5d4594_2650672()
	defer legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(oldFade)
	restoreFlags := noxflags.PortTestGameFlags(2048)
	defer restoreFlags()
	oldDialog := nox_gui_curDialog_830224
	nox_gui_curDialog_830224 = nil
	defer func() { nox_gui_curDialog_830224 = oldDialog }()
	var roots []*gui.Window
	for i := 0; i < 11; i++ {
		roots = append(roots, o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 30, 30, nil))
	}
	defer func() {
		for _, w := range roots {
			w.Destroy()
		}
	}()
	// chat, edit, book, conversation, console, console input, console scroll,
	// MOTD, vote, save dialog, quickbar capture.
	ptr := func(i int) uint32 { return uint32(uintptr(roots[i].C())) }
	*words["dword_5d4594_1064856"], *words["dword_5d4594_1064860"] = ptr(0), ptr(1)
	*words["dword_5d4594_1123524"] = ptr(3)
	*book["nox_win_unk1"] = ptr(2)
	*sessions["motd"], *sessions["motdList"] = ptr(7), ptr(6)
	*vote["window"] = ptr(8)
	oldConsole := guiCon
	guiCon = &guiConsole{root: roots[4], input: roots[5], scrollbox: roots[6]}
	defer func() { guiCon = oldConsole }()
	oldSave, oldSaveArray := dword_5d4594_1082856, nox_savegame_arr_1064948
	dword_5d4594_1082856 = roots[9]
	defer func() { dword_5d4594_1082856 = oldSave; nox_savegame_arr_1064948 = oldSaveArray }()
	var trace []string
	oldConsoleCall, oldSaveCall, oldFadeCall := legacy.Nox_gui_console_Hide_4512B0, legacy.Sub_46D6F0, legacy.Nox_video_inFadeTransition_44E0D0
	defer func() {
		legacy.Nox_gui_console_Hide_4512B0 = oldConsoleCall
		legacy.Sub_46D6F0 = oldSaveCall
		legacy.Nox_video_inFadeTransition_44E0D0 = oldFadeCall
	}()
	legacy.Nox_gui_console_Hide_4512B0 = func() int {
		trace = append(trace, "console")
		if !roots[2].GetFlags().IsHidden() {
			t.Fatal("console reached before book close")
		}
		return oldConsoleCall()
	}
	legacy.Sub_46D6F0 = func() int {
		trace = append(trace, "save")
		for _, i := range []int{2, 4, 7, 8} {
			if !roots[i].GetFlags().IsHidden() {
				t.Fatal("save reached before earlier close", i)
			}
		}
		return oldSaveCall()
	}
	legacy.Nox_video_inFadeTransition_44E0D0 = func() int { trace = append(trace, "fade"); return oldFadeCall() }
	type row struct {
		Gate, Mask        int
		Hidden            []bool
		ShopMode, Capture uint32
		Message           []byte
		Sounds            [][2]int
		Trace             []string
	}
	var captured []row
	for gate := 0; gate < 11; gate++ {
		for mask := 0; mask < 32; mask++ {
			for _, w := range roots {
				w.Hide()
				w.Flags |= gui.StatusEnabled
			}
			for bit, i := range []int{2, 4, 7, 8, 9} {
				roots[i].SetHidden(mask&(1<<bit) == 0)
			}
			*words["dword_5d4594_1064868"], *words["dword_5d4594_1123520"] = 0, 0
			*book["dword_5d4594_1047520"], *book["dword_5d4594_1046872"] = 0, 0
			*quick["dword_5d4594_1049532"], *quick["dword_5d4594_1047932"] = 0, 0
			*inv["dword_5d4594_1098624"], *inv["dword_5d4594_1098628"], *inv["dword_5d4594_1049864"] = 0, 0, 0
			binary.LittleEndian.PutUint32(cursor, 0)
			legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(0)
			clear(choice)
			choice[0] = 0xaa
			switch gate {
			case 1:
				binary.LittleEndian.PutUint32(cursor, 1)
			case 2:
				legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(1)
			case 3, 5:
				*quick["dword_5d4594_1049532"] = ptr(10)
				roots[10].Capture(true)
			case 6, 7, 8:
				*inv["dword_5d4594_1098628"] = uint32(gate - 4)
			case 9:
				*inv["dword_5d4594_1098624"] = 1
			case 10:
				roots[3].Show()
			}
			if gate == 4 || gate == 5 {
				roots[0].Show()
				*words["dword_5d4594_1064868"] = 1
			}
			o.sounds = nil
			o.c.srv.NetList.ResetAll()
			trace = nil
			if key {
				for _, state := range []uintptr{0, 1, 3, 0xffffffff} {
					if interactionCall("sub_46A7E0", uintptr(roots[1].C()), 21, 1, state) != 1 || len(trace) != 0 {
						t.Fatal("chat Escape non-release")
					}
				}
				if interactionCall("sub_46A7E0", uintptr(roots[1].C()), 21, 1, 2) != 1 {
					t.Fatal("chat Escape release return")
				}
			} else {
				interactionCall("nox_xxx_consoleEsc_49B7A0")
			}
			var msg []byte
			o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { msg = append(msg, b...); return false })
			var wantMsg []byte
			if gate == 9 {
				wantMsg = []byte{0xc9, 18}
			}
			if gate == 10 {
				wantMsg = []byte{0xd0, 2, 0}
			}
			var wantSounds [][2]int
			if gate == 0 {
				if mask&1 != 0 {
					wantSounds = append(wantSounds, [2]int{787, 100})
				}
				if mask == 0 {
					wantSounds = append(wantSounds, [2]int{231, 100})
				}
			}
			wantTrace := []string{"fade"}
			if gate == 1 {
				wantTrace = nil
			}
			if gate == 0 {
				wantTrace = append(wantTrace, "console", "save")
			}
			var hidden []bool
			for bit, i := range []int{2, 4, 7, 8, 9} {
				h := roots[i].GetFlags().IsHidden()
				want := gate == 0 || mask&(1<<bit) == 0
				if h != want {
					t.Fatal("Escape window gate/order", gate, mask, i, h, want)
				}
				hidden = append(hidden, h)
			}
			mode := *inv["dword_5d4594_1098628"]
			wantMode := uint32(0)
			if gate >= 6 && gate <= 8 {
				wantMode = 1
			}
			if mode != wantMode || *quick["dword_5d4594_1049532"] != 0 || !roots[0].GetFlags().IsHidden() || *words["dword_5d4594_1064868"] != 0 || !bytes.Equal(msg, wantMsg) || !slices.Equal(o.sounds, wantSounds) || !slices.Equal(trace, wantTrace) {
				t.Fatalf("Escape gate%d mask%d: mode%d message%x want%x sounds%v want%v trace%v want%v", gate, mask, mode, msg, wantMsg, o.sounds, wantSounds, trace, wantTrace)
			}
			captured = append(captured, row{gate, mask, hidden, mode, *quick["dword_5d4594_1049532"], msg, append([][2]int(nil), o.sounds...), append([]string(nil), trace...)})
		}
	}
	name := "escape-order"
	if key {
		name = "chat-escape-key"
	}
	interactionCapture(t, name, captured)
}
