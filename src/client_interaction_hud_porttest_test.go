//go:build porttest

package opennox

import (
	"encoding/binary"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionHUDVisibility(t *testing.T) {
	o := newMeterOwner(t)
	interaction, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	inv, restoreInv := legacy.PortTestInventoryWindowWords()
	defer restoreInv()
	book, restoreBook := legacy.PortTestBookWords()
	defer restoreBook()
	quick, restoreQuick := legacy.PortTestQuickbarWords()
	defer restoreQuick()
	team, restoreTeam := legacy.PortTestTeamUIWords()
	defer func() { *team["ctf"], *team["ball"] = 0, 0; restoreTeam() }()
	for _, p := range quick {
		*p = 0
	}
	serverConfigOwnBytes(t, 0x5D4594, 1049848, 4) // own the empty dragged-item pointer.
	expanded := serverConfigOwnBytes(t, 0x5D4594, 1049476, 4)
	clear(expanded)
	ctfVisible := serverConfigOwnBytes(t, 0x5D4594, 1045608, 4)
	savedFlag := serverConfigOwnBytes(t, 0x5D4594, 811064, 4)
	record, free := alloc.New([64]uint32{})
	defer free()
	restoreMain := legacy.PortTestBookQuickbar(unsafe.Pointer(record))
	defer restoreMain()
	var wins []*gui.Window
	for i := 0; i < 7; i++ {
		w := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 20, 20, nil)
		wins = append(wins, w)
	}
	defer func() {
		for _, w := range wins {
			w.Destroy()
		}
	}()
	word := func(w *gui.Window) uint32 { return uint32(uintptr(w.C())) }
	*o.meters.NamedWord("dword_5d4594_1090276") = word(wins[0])
	*inv["nox_win_unk5"] = word(wins[1])
	*book["nox_win_unk1"] = word(wins[2])
	record[52] = word(wins[3]) // actual quickbar Window offset208.
	*team["ctf"], *team["ball"] = word(wins[4]), word(wins[5])
	oldFPS := o.c.guiFPS
	defer func() { o.c.guiFPS = oldFPS }()
	o.c.guiFPS.win = wins[6]
	o.c.guiFPS.enabled = true
	oldGUI := nox_client_renderGUI_80828
	defer func() { nox_client_renderGUI_80828 = oldGUI }()
	oldEngine := noxflags.GetEngine()
	defer func() { noxflags.UnsetEngine(^noxflags.EngineFlag(0)); noxflags.SetEngine(oldEngine) }()
	oldFPSCall, oldDrag := legacy.Sub_4706C0, legacy.Sub_478000
	defer func() { legacy.Sub_4706C0 = oldFPSCall; legacy.Sub_478000 = oldDrag }()
	type row struct {
		Show                                                bool
		Old                                                 uint32
		Disabled, Book, Conversation, Team, InitiallyHidden bool
		Stored                                              uint32
		Hidden                                              []bool
		Calls                                               []int
	}
	var captured []row
	var calls []int
	var beforeFPS []bool
	legacy.Sub_4706C0 = func(v int) {
		calls = append(calls, 1)
		beforeFPS = nil
		for _, w := range wins[:6] {
			beforeFPS = append(beforeFPS, w.GetFlags().IsHidden())
		}
		oldFPSCall(v)
	}
	// Preserve the real drag-cancel implementation; observe its ordering only.
	legacy.Sub_478000 = func() int { calls = append(calls, 2); return oldDrag() }
	for _, show := range []bool{false, true} {
		for _, old := range []uint32{0, 1, 2} {
			for _, disabled := range []bool{false, true} {
				for _, bookOpen := range []bool{false, true} {
					for _, conversation := range []bool{false, true} {
						for _, teamVisible := range []bool{false, true} {
							for _, hidden := range []bool{false, true} {
								nox_client_renderGUI_80828 = show
								noxflags.UnsetEngine(noxflags.EngineNoRendering)
								if disabled {
									noxflags.SetEngine(noxflags.EngineNoRendering)
								}
								binary.LittleEndian.PutUint32(savedFlag, old)
								*book["dword_5d4594_1046864"] = 0
								if bookOpen {
									*book["dword_5d4594_1046864"] = 1
								}
								*interaction["dword_5d4594_1123520"] = 0
								if conversation {
									*interaction["dword_5d4594_1123520"] = 1
								}
								binary.LittleEndian.PutUint32(ctfVisible, 0)
								*team["ball-visible"] = 0
								if teamVisible {
									binary.LittleEndian.PutUint32(ctfVisible, 1)
									*team["ball-visible"] = 1
								}
								for _, w := range wins {
									w.SetHidden(hidden)
								}
								calls = nil
								beforeFPS = nil
								interactionCall("sub_437100")
								value := uint32(0)
								if show {
									value = 1
								}
								changed := old != value && !disabled
								stored := old
								want := make([]bool, 7)
								for i := range want {
									want[i] = hidden
								}
								var wantCalls []int
								if changed {
									stored = value
									wantCalls = []int{1}
									if !show {
										wantCalls = append(wantCalls, 2)
									}
									want[0], want[1], want[3] = !show, !show, !show
									if bookOpen && (show || !conversation) {
										want[2] = !show
									}
									want[4], want[5] = !(show && teamVisible && hidden), !(show && teamVisible && hidden)
									want[6] = !(show && hidden)
								}
								var got []bool
								for _, w := range wins {
									got = append(got, w.GetFlags().IsHidden())
								}
								if binary.LittleEndian.Uint32(savedFlag) != stored || !slices.Equal(got, want) || !slices.Equal(calls, wantCalls) || changed && !slices.Equal(beforeFPS, want[:6]) {
									t.Fatalf("HUD show%v old%d disabled%v book%v conversation%v team%v hidden%v: states%v want%v calls%v want%v", show, old, disabled, bookOpen, conversation, teamVisible, hidden, got, want, calls, wantCalls)
								}
								captured = append(captured, row{show, old, disabled, bookOpen, conversation, teamVisible, hidden, stored, got, append([]int(nil), calls...)})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "hud-visibility", captured)
}
