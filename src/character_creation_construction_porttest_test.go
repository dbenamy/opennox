//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestCharacterCreationConstructionFailures(t *testing.T) {
	type row struct {
		Color, Missing     bool
		Resource           string
		Return, Animations int
		State              gui.StateID
		WindowAlive        bool
		Appearance         byte
	}
	var rows []row
	for _, color := range []bool{false, true} {
		for _, missing := range []bool{false, true} {
			t.Run(map[bool]string{false: "class", true: "color"}[color]+map[bool]string{false: "/animation", true: "/resource"}[missing], func(t *testing.T) {
				o := newEntryOwner(t)
				defer legacy.PortTestCharacterAppearanceOwner(nil, make([]unsafe.Pointer, 15))()
				oldBack := guiOptsBack
				guiOptsBack = nil
				defer func() { guiOptsBack = oldBack }()
				t.Cleanup(legacy.PortTestQuestProgressOwner())
				host := serverConfigOwnBytes(t, 0x5D4594, 807172, 128)
				host[67] = 7
				defer legacy.PortTestCharacterClassOwner(nil)()
				defer legacy.PortTestCharacterClassWindow(nil)()
				defer legacy.PortTestCharacterPaletteOwner(nil, nil)()
				defer legacy.PortTestCharacterAnimationOwner(color, nil)()
				var loaded *gui.Window
				resource := ""
				animations := 0
				oldLoad := legacy.Nox_new_window_from_file
				defer func() { legacy.Nox_new_window_from_file = oldLoad }()
				legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
					resource = name
					if missing {
						return nil
					}
					loaded = o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 640, 480, fn)
					return loaded
				}
				oldAnim := legacy.Nox_gui_makeAnimation_43C5B0
				defer func() { legacy.Nox_gui_makeAnimation_43C5B0 = oldAnim }()
				legacy.Nox_gui_makeAnimation_43C5B0 = func(_ *gui.Window, _, _, _, _, _, _, _, _ int) *gui.Anim { animations++; return nil }
				ret := legacy.PortTestCharacterConstruct(color)
				wantResource, wantState, wantAppearance := "SelClass.wnd", gui.StateID(600), byte(7)
				if color {
					wantResource = "SelColor.wnd"
					wantState = 700
					wantAppearance = 0
				}
				wantAnimations := 1
				if missing {
					wantAnimations = 0
				}
				alive := loaded != nil && loaded.Flags&gui.StatusDestroyed == 0
				if ret != 0 || resource != wantResource || o.c.GameGetStateCode() != wantState || animations != wantAnimations || alive == missing || host[67] != wantAppearance {
					t.Fatal("construction failure contract", color, missing, ret, resource, animations, alive, host[67])
				}
				rows = append(rows, row{color, missing, resource, ret, animations, o.c.GameGetStateCode(), alive, host[67]})
			})
		}
	}
	spellbookCapture(t, "character-creation-construction-failures", rows, "91244c69d14d1b7af0f34867dd4dd59969fc86432a014482f7a74caca61fc31a")
}
func TestCharacterCreationAdmissionWords(t *testing.T) {
	defer legacy.PortTestCharacterDefaultOwner(0)()
	serverConfigOwnBytes(t, 0x5D4594, 1308168, 4)
	var rows [][3]uint32
	for _, v := range []int32{0, 1, 2, -1, -2147483648, 2147483647} {
		a := legacy.PortTestCharacterDefaultSet(v)
		b := legacy.PortTestCharacterAdmissionSet(v)
		if a != v || b != v || legacy.PortTestCharacterDefaultGet() != uint32(v) || memmap.Uint32(0x5D4594, 1308168) != uint32(v) {
			t.Fatal("full admission/default words", v)
		}
		rows = append(rows, [3]uint32{uint32(v), uint32(a), uint32(b)})
	}
	spellbookCapture(t, "character-creation-admission-words", rows, "6556110b7ca7eaab61969ef66aa9784d9c29e1a70212f7c0c2642600fddaeab3")
}
