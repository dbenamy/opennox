//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func sessionQuitWindow(t *testing.T, o *listboxOwner, words map[string]*uint32) *gui.Window {
	t.Helper()
	w := legacy.Nox_new_window_from_file("QuitMenu.wnd", func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		return gui.RawEventResp(sessionCall("quitEvent", w, ev.EventCode(), uint32(a), uint32(b)))
	})
	if w == nil {
		t.Fatal("quit resource")
	}
	*words["quit"] = uint32(uintptr(w.C()))
	return w
}
func TestSessionQuitVisibilityAndModes(t *testing.T) {
	restoreAllFlags := noxflags.PortTestGameFlags(0)
	defer restoreAllFlags()
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	book, restoreBook := legacy.PortTestBookWords()
	defer restoreBook()
	colors, restoreColors := legacy.PortTestInventoryDisplayWords()
	defer restoreColors()
	*colors["nox_color_black_2650656"] = 0x12345678
	*colors["nox_color_orange_2614256"] = 0x76543210
	player := serverConfigOwnBytes(t, 0x852978, 8, 4)
	raw, free := alloc.Make([]byte{}, 320)
	defer free()
	binary.LittleEndian.PutUint32(player, uint32(uintptr(unsafe.Pointer(&raw[0]))))
	set, restoreStrings := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "guiquit.c:SoloSaveLabel", Vals: []strman.Variant{{Str: "Save game"}}}, strman.Entry{ID: "guiquit.c:MultiplayerSaveLabel", Vals: []strman.Variant{{Str: "Save player"}}})
	defer restoreStrings()
	set(0)
	var pause []int
	oldPause := legacy.Sub_413A00
	defer func() { legacy.Sub_413A00 = oldPause }()
	legacy.Sub_413A00 = func(v int) { pause = append(pause, v) }
	oldEngine := noxflags.GetEngine()
	defer func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) }()
	type row struct {
		Mode      uint32
		TeamCount int
		Headless  bool
		Gate      int
		Shown     bool
		Size      image.Point
		Children  map[uint]uint32
		Pause     []int
	}
	var rows []row
	oldTeams := o.c.srv.Teams.ActiveCnt
	defer func() { o.c.srv.Teams.ActiveCnt = oldTeams }()
	for _, teams := range []int{0, 2, 256} {
		o.c.srv.Teams.ActiveCnt = teams
		for _, mode := range []uint32{0, 2048, 4096, 49152} {
			for _, headless := range []bool{false, true} {
				for gate := 0; gate < 7; gate++ {
					flags := noxflags.GameFlag(mode)
					if gate == 2 {
						flags |= noxflags.GamePause
					}
					restoreFlags := noxflags.PortTestGameFlags(flags)
					noxflags.ResetEngine()
					if headless {
						noxflags.SetEngine(noxflags.EngineNoRendering)
					}
					*book["dword_5d4594_1047520"] = 0
					if gate == 1 {
						*book["dword_5d4594_1047520"] = 1
					}
					state := uint32(0)
					switch gate {
					case 3:
						state = 1
					case 4:
						state = 2
					case 5:
						state = 51
					case 6:
						state = 3
					}
					binary.LittleEndian.PutUint32(raw[276:], state)
					w := sessionQuitWindow(t, o, words)
					w.SetHidden(true)
					pause = nil
					if sessionCall("quitColors", nil, 0, 0, 0) != uintptr(w.ChildByID(9006).C()) || w.DrawData().BgColorVal != 0x12345678 {
						t.Fatal("quit color owner/return")
					}
					for id := uint(9001); id <= 9006; id++ {
						if w.ChildByID(id).DrawData().TextColorVal != 0x76543210 {
							t.Fatal("quit text color", id)
						}
					}
					if sessionCall("quitShown", nil, 0, 0, 0) != 0 {
						t.Fatal("hidden quit predicate")
					}
					sessionCall("quitToggle", nil, 0, 0, 0)
					blocked := gate == 1 || gate == 2 || (mode&2048 != 0 && (state == 1 || state == 2 || state == 51))
					shown := sessionCall("quitShown", nil, 0, 0, 0) != 0
					if shown == blocked {
						t.Fatal("quit opening gate", mode, headless, gate, shown)
					}
					r := row{TeamCount: teams, Mode: mode, Headless: headless, Gate: gate, Shown: shown, Size: w.SizeVal, Children: map[uint]uint32{}}
					if shown {
						solo := mode&2048 != 0
						size := image.Pt(220, 330)
						if solo {
							size.Y = 285
						}
						if w.SizeVal != size || o.c.GUI.Captured() != w {
							t.Fatal("quit size/capture", mode, w.SizeVal)
						}
						for _, id := range []uint{9001, 9002} {
							if w.ChildByID(id).Flags.Has(gui.StatusHidden) == solo {
								t.Fatal("solo action visibility", mode, id)
							}
						}
						for _, id := range []uint{9007, 9008, 9009} {
							if w.ChildByID(id).Flags.Has(gui.StatusHidden) != solo {
								t.Fatal("multiplayer action visibility", mode, id)
							}
						}
						if !solo {
							if w.ChildByID(9009).Flags.Has(gui.StatusEnabled) != (mode&49152 == 0 && uint8(teams) != 0) {
								t.Fatal("team action enabled", mode, teams)
							}
							for _, id := range []uint{9003, 9007} {
								want := mode&4096 == 0 && !headless
								if w.ChildByID(id).Flags.Has(gui.StatusEnabled) != want {
									t.Fatal("save/admin enabled", mode, headless, id)
								}
							}
							if w.ChildByID(9005).Flags.Has(gui.StatusEnabled) == headless {
								t.Fatal("headless options enable")
							}
						}
						for id := uint(9001); id <= 9009; id++ {
							r.Children[id] = uint32(w.ChildByID(id).Flags)
						}
						sessionCall("quitHide", nil, 0, 0, 0)
						if sessionCall("quitShown", nil, 0, 0, 0) != 0 || w.Flags.Has(gui.StatusEnabled) || o.c.GUI.Captured() != nil {
							t.Fatal("quit hiding")
						}
						if !solo && fmt.Sprint(pause) != "[0]" || solo && fmt.Sprint(pause) != "[1 0]" {
							t.Fatal("pause transitions", mode, pause)
						}
					} else if len(pause) != 0 {
						t.Fatal("blocked menu paused", pause)
					}
					r.Pause = append([]int{}, pause...)
					rows = append(rows, r)
					w.Destroy()
					o.c.GUI.FreeDestroyed()
					*words["quit"] = 0
					restoreFlags()
				}
			}
		}
	}
	sessionCapture(t, "quit-modes", rows)
}

func TestSessionQuitCallbacks(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	restoreFlags := noxflags.PortTestGameFlags(0)
	defer restoreFlags()
	w := sessionQuitWindow(t, o, words)
	defer w.Destroy()
	var calls []string
	pause, load, quit, dialog := legacy.Sub_413A00, legacy.Sub_445B40, legacy.Nox_client_quit_4460C0, legacy.Nox_xxx_dialogMsgBoxCreate_449A10
	defer func() {
		legacy.Sub_413A00 = pause
		legacy.Sub_445B40 = load
		legacy.Nox_client_quit_4460C0 = quit
		legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = dialog
	}()
	legacy.Sub_413A00 = func(v int) { calls = append(calls, fmt.Sprint("pause", v)) }
	legacy.Sub_445B40 = func() int { calls = append(calls, "load"); return 7 }
	legacy.Nox_client_quit_4460C0 = func() { calls = append(calls, "quit") }
	var yes, no func()
	legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(parent *gui.Window, title, text string, flags gui.DialogFlags, a, b func()) {
		if parent != w || flags != 56 || a == nil || b == nil {
			t.Fatal("quit confirmation arguments", parent, flags, a == nil, b == nil)
		}
		yes, no = a, b
		calls = append(calls, "confirm")
	}
	type row struct {
		Action      string
		Calls       []string
		Highlighted bool
		Hidden      bool
	}
	var rows []row
	invoke := func(id uint) {
		button := w.ChildByID(id)
		button.DrawData().Field0 = 0xa5a50006
		w.SetHidden(false)
		w.Flags |= gui.StatusEnabled
		w.Capture(true)
		if sessionCall("quitEvent", w, 16391, uint32(uintptr(button.C())), 0) != 0 {
			t.Fatal("quit click result")
		}
		if button.DrawData().Field0 != 0xa5a50004 {
			t.Fatal("highlight cleanup", id, button.DrawData().Field0)
		}
	}
	for _, code := range []int{0, 21, 22, 23, 16390} {
		if sessionCall("quitEvent", w, code, 0xffffffff, 0) != 0 {
			t.Fatal("non-click numeric notification")
		}
	}
	calls = nil
	invoke(9001)
	if fmt.Sprint(calls) != "[pause0 pause1 load]" {
		t.Fatal("load action", calls)
	}
	rows = append(rows, row{Action: "load", Calls: append([]string{}, calls...), Hidden: w.Flags.Has(gui.StatusHidden)})
	calls = nil
	invoke(9004)
	if fmt.Sprint(calls) != "[confirm]" || o.c.GUI.Captured() != nil {
		t.Fatal("confirmation capture", calls)
	}
	no()
	if o.c.GUI.Captured() != w {
		t.Fatal("cancel recapture")
	}
	yes()
	if fmt.Sprint(calls) != "[confirm quit pause0]" || !w.Flags.Has(gui.StatusHidden) {
		t.Fatal("confirmed quit", calls)
	}
	rows = append(rows, row{Action: "quit-confirmation", Calls: append([]string{}, calls...), Hidden: w.Flags.Has(gui.StatusHidden)})
	calls = nil
	invoke(9006)
	if fmt.Sprint(calls) != "[pause0]" || !w.Flags.Has(gui.StatusHidden) {
		t.Fatal("resume action", calls)
	}
	rows = append(rows, row{Action: "resume", Calls: append([]string{}, calls...), Hidden: w.Flags.Has(gui.StatusHidden)})
	sessionCapture(t, "quit-callbacks", rows)
}

func TestSessionQuitSaveAndLoadGates(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	restoreFlags := noxflags.PortTestGameFlags(noxflags.GameModeCoop)
	defer restoreFlags()
	player := serverConfigOwnBytes(t, 0x852978, 8, 4)
	raw, free := alloc.Make([]byte{}, 320)
	defer free()
	binary.LittleEndian.PutUint32(player, uint32(uintptr(unsafe.Pointer(&raw[0]))))
	w := sessionQuitWindow(t, o, words)
	defer w.Destroy()
	pause, load, save, name, dialog := legacy.Sub_413A00, legacy.Sub_445B40, legacy.Sub_4DB170, legacy.Nox_setSaveFileName_4DB130, legacy.Nox_xxx_dialogMsgBoxCreate_449A10
	defer func() {
		legacy.Sub_413A00 = pause
		legacy.Sub_445B40 = load
		legacy.Sub_4DB170 = save
		legacy.Nox_setSaveFileName_4DB130 = name
		legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = dialog
	}()
	var calls []string
	legacy.Sub_413A00 = func(v int) { calls = append(calls, fmt.Sprint("pause", v)) }
	legacy.Sub_445B40 = func() int { calls = append(calls, "load"); return 7 }
	legacy.Nox_setSaveFileName_4DB130 = func(s string) { calls = append(calls, s) }
	legacy.Sub_4DB170 = func(a bool, b unsafe.Pointer, c int) {
		if !a || b != nil || c != 0 {
			t.Fatal("autosave arguments")
		}
		calls = append(calls, "save")
	}
	var yes, no func()
	legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(parent *gui.Window, title, text string, flags gui.DialogFlags, a, b func()) {
		if parent != nil || flags != 56 || a == nil || b == nil {
			t.Fatal("load confirmation")
		}
		yes, no = a, b
		calls = append(calls, "confirm")
	}
	invoke := func(id uint) {
		w.SetHidden(false)
		w.ChildByID(id).Flags |= gui.StatusEnabled
		w.ChildByID(id).DrawData().Field0 = 0xa5a50006
		calls = nil
		if sessionCall("quitEvent", w, 16391, uint32(uintptr(w.ChildByID(id).C())), 0) != 0 || w.ChildByID(id).DrawData().Field0 != 0xa5a50004 {
			t.Fatal("button return/highlight", id)
		}
	}
	for _, state := range []uint32{0, 1, 2, 3, 51} {
		binary.LittleEndian.PutUint32(raw[276:], state)
		invoke(9001)
		if state == 1 || state == 2 || state == 51 {
			if fmt.Sprint(calls) != "[pause0 pause1 load]" {
				t.Fatal("direct load", state, calls)
			}
		} else {
			if fmt.Sprint(calls) != "[pause0 pause1 confirm]" {
				t.Fatal("load confirmation", state, calls)
			}
			no()
			yes()
			if fmt.Sprint(calls) != "[pause0 pause1 confirm pause0 load]" {
				t.Fatal("load callbacks", calls)
			}
		}
	}
	for _, status := range []uint32{0, 1, 0x8000, 0xffffffff} {
		binary.LittleEndian.PutUint32(raw[120:], status)
		invoke(9002)
		if status&0x8000 != 0 {
			if len(calls) != 0 || w.ChildByID(9002).Flags.Has(gui.StatusEnabled) || w.Flags.Has(gui.StatusHidden) {
				t.Fatal("autosave disabled", status, calls)
			}
		} else if fmt.Sprint(calls) != "[pause0 AUTOSAVE save]" || !w.Flags.Has(gui.StatusHidden) {
			t.Fatal("autosave", status, calls)
		}
	}
}
