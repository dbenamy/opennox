//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestCharacterCreationAnimationCleanup(t *testing.T) {
	type row struct {
		Color                                 bool
		Return, Calls                         int
		RootDestroyed, MenuDestroyed, Removed bool
	}
	var rows []row
	for _, color := range []bool{false, true} {
		t.Run(map[bool]string{false: "class", true: "appearance"}[color], func(t *testing.T) {
			o := newEntryOwner(t)
			root := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 90, 90, nil)
			menu := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 20, 20, nil)
			defer legacy.PortTestCharacterClassWindow(root.C())()
			defer legacy.PortTestCharacterPaletteOwner(root.C(), menu.C())()
			state := gui.AnimGlobalState()
			defer gui.SetAnimGlobalState(state)
			a := gui.NewAnim(root, image.Pt(0, 0), image.Pt(0, -460), image.Pt(0, 20), image.Pt(0, -40))
			id := gui.StateID(600)
			if color {
				id = 700
			}
			a.StateID = id
			released := false
			defer func() {
				if !released {
					a.Free()
				}
			}()
			defer legacy.PortTestCharacterAnimationOwner(color, unsafe.Pointer(a))()
			fn, count, restore := legacy.PortTestOptionsDone()
			defer restore()
			a.Func13Ptr = fn
			if !color {
				if legacy.PortTestCharacterClassStart() != 1 || a.State() != gui.AnimOut || gui.AnimGlobalState() != gui.AnimOut {
					t.Fatal("class animation start")
				}
			}
			ret := legacy.PortTestCharacterAnimationDone(color)
			released = true
			rootGone := root.Flags&gui.StatusDestroyed != 0
			menuGone := menu.Flags&gui.StatusDestroyed != 0
			removed := gui.FindAnimForStateID(id) == nil
			if ret != 1 || count() != 1 || !rootGone || !menuGone || !removed {
				t.Fatal("animation cleanup", ret, count(), rootGone, menuGone, removed)
			}
			rows = append(rows, row{color, ret, count(), rootGone, menuGone, removed})
		})
	}
	spellbookCapture(t, "character-creation-animation-cleanup", rows, "f6b84f79806e4922936852eb4718363032b8b05cda42c1b41cc2112bdabd1e92")
}
