//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionSettingsNotices(t *testing.T) {
	catalog := legacy.PortTestMapCatalogOpen(9)
	t.Cleanup(catalog.Close)
	o := newServerOptionsOwner(t)
	root := *o.optionWords["root"]
	words, restore := legacy.PortTestClientSessionWords()
	t.Cleanup(restore)
	ui, restoreUI := legacy.PortTestClientInteractionWords()
	t.Cleanup(restoreUI)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	acquired := serverConfigOwnBytes(t, 0x5D4594, 371704, 4)
	records := serverConfigOwnBytes(t, 0x5D4594, 371380, 116)
	clear(serverConfigOwnBytes(t, 0x5D4594, 1556160, 4))
	serverConfigOwnBytes(t, 0x5D4594, 1304400, 1280)
	table := serverConfigOwnBytes(t, 0x587000, 164512, 32)
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint32(table[4*i:], uint32(uintptr(unsafe.Pointer(alloc.InternCString("fixture-help.wnd")))))
	}
	oldHelp := legacy.Get_nox_server_sanctuaryHelp_54276()
	t.Cleanup(func() { legacy.Set_nox_server_sanctuaryHelp_54276(oldHelp) })
	oldPlayer := legacy.Get_dword_8531A0_2576()
	legacy.Set_dword_8531A0_2576(&o.players[0])
	t.Cleanup(func() { legacy.Set_dword_8531A0_2576(oldPlayer) })
	o.players[0].SetName("Player Ω")
	set, restoreStrings := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "%s"}}}, strman.Entry{ID: "cdecode.c:NameChange", Vals: []strman.Variant{{Str: "Renamed %s"}}}, strman.Entry{ID: "Sanchlp.wnd:ClientHelp", Vals: []strman.Variant{{Str: "Help %s"}}}, strman.Entry{ID: "cdecode.c:KeyToChat", Vals: []strman.Variant{{Str: "Chat %s"}}})
	t.Cleanup(restoreStrings)
	set(0)
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	t.Cleanup(o.c.Render().GetFonts().PortTestDefaultFont(basicfont.Face7x13))
	o.c.ctrl = new(CtrlEventHandler)
	o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{keybind.KeyQ}, events: []keybind.Event{45}})
	o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{keybind.KeyF12}, events: []keybind.Event{8}})
	oldLoad := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldLoad })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "fixture-help.wnd" {
			return oldLoad(name, fn)
		}
		return newWindowFromString(o.c.GUI, "FONT = small; WINDOW 4100 0 0 300 120 USER; STATUS = ENABLED; CHILD WINDOW 4102 0 0 280 50 STATICTEXT; STATUS = ENABLED; DATA = 0 0 WindowDir:Blank; END WINDOW 4103 0 60 50 20 PUSHBUTTON; STATUS = ENABLED; END WINDOW 4104 80 60 50 20 CHECKBOX; STATUS = ENABLED; END END END", fn)
	}
	printer := new(teamRuntimePrinter)
	con := console.New(printer)
	oldConsole := legacy.GetConsole
	legacy.GetConsole = func() *console.Console { return con }
	t.Cleanup(func() { legacy.GetConsole = oldConsole })
	messageRegion := serverConfigOwnBytes(t, 0x5D4594, 823804, 1932)
	messageIndex := legacy.PortTestTeamRuntimeTextIndex()
	oldIndex := *messageIndex
	t.Cleanup(func() { *messageIndex = oldIndex })
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(func() { interactionCall("sub_49C7A0") })
	type row struct {
		On, Open, Acquired, Mode, Notice, Help int
		AfterNotice                            uint32
		Messages                               []string
		Text                                   string
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for open := 0; open < 2; open++ {
			for prior := 0; prior < 2; prior++ {
				for mode := 0; mode < 2; mode++ {
					for notice := 0; notice < 2; notice++ {
						for help := 0; help < 2; help++ {
							interactionCall("sub_49C7A0")
							o.c.GUI.FreeDestroyed()
							noxflags.ResetGame()
							if mode != 0 {
								noxflags.SetGame(128)
							}
							*o.optionWords["root"] = root
							if open == 0 {
								*o.optionWords["root"] = 0
							}
							*words["settingsNotice"] = uint32(notice)
							legacy.Set_nox_server_sanctuaryHelp_54276(help)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							binary.LittleEndian.PutUint32(acquired, uint32(prior))
							clear(records)
							binary.LittleEndian.PutUint16(records[52:], 0x100)
							data := make([]byte, 60)
							data[0] = 177
							copy(data[2:], records[:58])
							input := bytes.Clone(data)
							printer.lines = nil
							clear(messageRegion)
							*messageIndex = 0
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(177), data)
							gate := open == 0 && prior == 0 && mode != 0
							wantNotice := uint32(notice)
							var wantMessages []string
							if gate && notice != 0 {
								wantNotice = 0
								wantMessages = []string{"Renamed Player Ω"}
							}
							w := (*gui.Window)(unsafe.Pointer(uintptr(*ui["dword_5d4594_1305680"])))
							text := ""
							if w != nil {
								text = interactionStaticText(w.ChildByID(4102))
							}
							if n != 60 || !bytes.Equal(input, data) || *words["settingsNotice"] != wantNotice || !reflect.DeepEqual(printer.lines, wantMessages) || (w != nil) != (gate && help != 0) || (w != nil && text != "Help Q Chat F12") || binary.LittleEndian.Uint32(acquired) != 1 {
								t.Fatal("settings notice/help", on, open, prior, mode, notice, help, n, printer.lines, text)
							}
							rows = append(rows, row{on, open, prior, mode, notice, help, wantNotice, append([]string(nil), printer.lines...), text})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-settings-notices", rows)
}
