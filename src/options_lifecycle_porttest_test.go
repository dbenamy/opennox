//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestOptionsCloseAndVisibility(t *testing.T) {
	type record struct {
		Initial, Cancel, Visible, Return, After, Writes int
		Hidden, OverlayHidden, QuitHidden               bool
		Dialog                                          string
	}
	var rows []record
	for initial := 0; initial < 3; initial++ {
		for _, cancel := range []int{0, 1, -1} {
			t.Run(fmt.Sprintf("initial%d/cancel%d", initial, cancel), func(t *testing.T) {
				o := newOptionsAudioOwner(t, false)
				noxflags.ResetGame()
				inventory, restore := legacy.PortTestInventoryWindowWords()
				t.Cleanup(restore)
				binding, restore := legacy.PortTestBindingWords()
				t.Cleanup(restore)
				*binding[1321228] = 0
				quit := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
				*inventory["nox_wnd_quitMenu_825760"] = uint32(uintptr(quit.C()))
				overlay := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
				*o.words[1309824] = uint32(uintptr(overlay.C()))
				if initial == 0 {
					*o.words[1309820] = 0
				} else {
					*o.words[1309820] = uint32(uintptr(o.root.C()))
					o.root.SetHidden(initial == 1)
				}
				*o.words[122848] = 1
				legacy.Dialogs.PlayFile("before.wav", 100)
				writes := 0
				oldWriter := legacy.WriteConfigLegacy
				t.Cleanup(func() { legacy.WriteConfigLegacy = oldWriter })
				legacy.WriteConfigLegacy = func(name string) {
					if name != "nox.cfg" {
						t.Fatalf("unexpected path %q", name)
					}
					writes++
				}
				r := record{Initial: initial, Cancel: cancel, Visible: legacy.PortTestOptionsAction(3, 0)}
				r.Return = legacy.PortTestOptionsAction(2, cancel)
				r.After = legacy.PortTestOptionsAction(3, 0)
				r.Writes = writes
				r.Hidden = o.root.Flags.IsHidden()
				r.OverlayHidden = overlay.Flags.IsHidden()
				r.QuitHidden = quit.Flags.IsHidden()
				r.Dialog = legacy.Dialogs.FileToRead()
				want := 0
				if initial == 2 {
					want = 1
				}
				wantWrites := 0
				if initial == 2 && cancel == 0 {
					wantWrites = 1
				}
				if r.Visible != want || r.Return != want || r.After != 0 || writes != wantWrites {
					t.Fatalf("close %+v", r)
				}
				if initial == 2 && (!r.Hidden || !r.OverlayHidden || r.Dialog != "") {
					t.Fatalf("visible cleanup %+v", r)
				}
				if initial != 2 && r.Dialog != "before.wav" {
					t.Fatal("hidden close changed dialog")
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "options-close", rows, "e30e09f9dffc4ad3aa88805a4e2e772fa7c266c0bda5c0c4e12326852f4af215")
}

func TestOptionsMenuCompletion(t *testing.T) {
	o := newOptionsAudioOwner(t, true)
	prepareOptionsConstruction(t, o, false)
	if legacy.PortTestOptionsConstruct(true) != 1 {
		t.Fatal("construction failed")
	}
	anim := gui.FindAnimForStateID(300)
	if anim == nil {
		t.Fatal("missing animation")
	}
	fn, count, restore := legacy.PortTestOptionsDone()
	t.Cleanup(restore)
	anim.Func13Ptr = fn
	root := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[1309720])))
	ret := legacy.PortTestOptionsAction(5, 0)
	if ret != 1 || count() != 1 || gui.FindAnimForStateID(300) != nil {
		t.Fatalf("completion return=%d calls=%d", ret, count())
	}
	// Destruction is deferred by GUI ownership, but the window is already marked.
	if root.Flags&gui.StatusDestroyed == 0 {
		t.Fatal("menu root was not destroyed")
	}
	spellbookCapture(t, "options-menu-completion", [][2]int{{ret, count()}}, "5cc6a8a697bb6f2fcb468ef02cf2a76e9ad3f5da223682f1ca29288836ee66b3")
}
