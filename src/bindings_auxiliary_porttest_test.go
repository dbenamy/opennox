//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestBindingEditorAuxiliaryWindow(t *testing.T) {
	type record struct {
		Size            [2]int
		Missing         bool
		Present, Hidden bool
		Position        [2]int
	}
	var records []record
	for _, size := range [][2]int{{640, 480}, {1024, 768}, {0, 0}} {
		for _, missing := range []bool{false, true} {
			t.Run(fmt.Sprintf("%v/missing=%v", size, missing), func(t *testing.T) {
				o := newBindingOwner(t, false)
				oldParser := legacy.Nox_new_window_from_file
				t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
				legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
					if name != "yesno.wnd" {
						t.Fatalf("unexpected resource %q", name)
					}
					if missing {
						return newWindowFromString(o.c.GUI, "", fn)
					}
					return newWindowFromString(o.c.GUI, "WINDOW\n 950 0 0 320 240 USER;\n STATUS = ENABLED;\n END\n", fn)
				}
				dims := legacy.PortTestBindingDimensions()
				old := [2]int32{*dims[0], *dims[1]}
				t.Cleanup(func() { *dims[0], *dims[1] = old[0], old[1] })
				*dims[0], *dims[1] = int32(size[0]), int32(size[1])
				got := legacy.PortTestBindingInvoke("sub_4C3500", [4]uint32{})
				r := record{Size: size, Missing: missing, Present: got != 0}
				if got != *o.words[1321224] || r.Present == missing {
					t.Fatalf("auxiliary creation %+v", r)
				}
				if got != 0 {
					w := (*gui.Window)(unsafe.Pointer(uintptr(got)))
					r.Hidden = w.Flags.IsHidden()
					r.Position = [2]int{w.Off.X, w.Off.Y}
					if !r.Hidden || r.Position != [2]int{(size[0] - 320) / 2, (size[1] - 240) / 2} {
						t.Fatalf("auxiliary geometry %+v", r)
					}
				}
				records = append(records, r)
			})
		}
	}
	spellbookCapture(t, "binding-auxiliary", records, "a29f9eb5aa7ead40c99d9ff7d103f206aea0e2bedf73fcd3c110284fd7a98d72")
}

func TestBindingEditorCancel(t *testing.T) {
	type record struct {
		Initial                int
		Visible, Cancel, After uint32
		Hidden                 bool
		Writes                 int
	}
	var records []record
	for initial := 0; initial < 3; initial++ {
		t.Run(fmt.Sprint(initial), func(t *testing.T) {
			o := newBindingOwner(t, false)
			noxflags.ResetGame()
			writes := 0
			oldWriter := legacy.WriteConfigLegacy
			t.Cleanup(func() { legacy.WriteConfigLegacy = oldWriter })
			legacy.WriteConfigLegacy = func(string) { writes++ }
			w := o.parent
			if initial == 0 {
				*o.words[1321228] = 0
			} else {
				*o.words[1321228] = uint32(uintptr(w.C()))
				w.SetHidden(initial == 1)
			}
			r := record{Initial: initial, Visible: legacy.PortTestBindingInvoke("sub_4C4280", [4]uint32{})}
			r.Cancel = legacy.PortTestBindingInvoke("sub_4C35B0", [4]uint32{1})
			r.After = legacy.PortTestBindingInvoke("sub_4C4280", [4]uint32{})
			r.Hidden = w.Flags.IsHidden()
			r.Writes = writes
			want := uint32(0)
			if initial == 2 {
				want = 1
			}
			if r.Visible != want || r.Cancel != want || r.After != 0 || r.Writes != 0 {
				t.Fatalf("cancel contract %+v", r)
			}
			records = append(records, r)
		})
	}
	spellbookCapture(t, "binding-cancel", records, "167ba3826021173a744054ba1f28b1cb88bea697f5528f9f70f7d77c0fa6b3c9")
}
