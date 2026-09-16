//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
	"strings"
	"testing"
	"unsafe"
)

func bindingResource() string {
	var b strings.Builder
	b.WriteString("FONT = small;\nWINDOW\n 900 0 0 640 480 USER;\n STATUS = ENABLED;\n CHILD\n")
	for i := 0; i < 4; i++ {
		fmt.Fprintf(&b, " WINDOW\n %d %d 40 140 260 SCROLLLISTBOX;\n STATUS = ENABLED;\n DATA = 128 1 0 0 1 0 0;\n END\n", 910+i, 10+150*i)
	}
	for _, id := range []int{931, 932, 933, 971, 972, 973, 974} {
		fmt.Fprintf(&b, " WINDOW\n %d 0 320 40 24 PUSHBUTTON;\n STATUS = ENABLED;\n END\n", id)
	}
	b.WriteString(" WINDOW\n 980 0 100 300 80 USER;\n STATUS = ENABLED+HIDDEN;\n CHILD\n WINDOW\n 981 0 0 280 50 STATICTEXT;\n STATUS = ENABLED;\n DATA = 0 0 InputCfg.wnd:PressKey;\n END\n END\n END\n END\nEND\n")
	return b.String()
}

func TestBindingEditorConstruction(t *testing.T) {
	type record struct {
		Menu, Missing         bool
		Width, Height, Result int
		RootPos, ModalPos     [2]int
		Hidden, ModalHidden   bool
		IDs                   []uint32
		Counts                []uint16
		AnimState             int
	}
	var records []record
	for _, menu := range []bool{false, true} {
		for _, size := range [][2]int{{640, 480}, {1024, 768}, {480, 640}} {
			for _, missing := range []bool{false, true} {
				t.Run(fmt.Sprintf("menu=%v/size=%v/missing=%v", menu, size, missing), func(t *testing.T) {
					o := newBindingOwner(t, menu)
					oldStrings := strMan
					strMan = o.c.Strings()
					t.Cleanup(func() { strMan = oldStrings })
					oldBinding := keyBinding
					keyBinding = keybind.New(o.c.Strings())
					t.Cleanup(func() { keyBinding = oldBinding })
					_, restoreFont := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "small")
					t.Cleanup(restoreFont)
					oldParser := legacy.Nox_new_window_from_file
					t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
					legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
						if name != "InputCfg.wnd" {
							t.Fatalf("unexpected resource %q", name)
						}
						if missing {
							return newWindowFromString(o.c.GUI, "", fn)
						}
						return newWindowFromString(o.c.GUI, bindingResource(), fn)
					}
					for _, off := range []int{1321256, 1522636} {
						b := memmap.BlobByAddr(0x5D4594).Data[off : off+512]
						old := append([]byte(nil), b...)
						t.Cleanup(func() { copy(b, old) })
						clear(b)
					}
					dimensions := legacy.PortTestBindingDimensions()
					old := [2]int32{*dimensions[0], *dimensions[1]}
					t.Cleanup(func() { *dimensions[0], *dimensions[1] = old[0], old[1] })
					*dimensions[0], *dimensions[1] = int32(size[0]), int32(size[1])
					oldAnim := gui.AnimGlobalState()
					t.Cleanup(func() { gui.SetAnimGlobalState(oldAnim) })
					for _, p := range o.words {
						*p = 0
					}
					result := legacy.PortTestBindingConstruct(menu)
					r := record{Menu: menu, Missing: missing, Width: size[0], Height: size[1], Result: result, AnimState: -1}
					if missing {
						if result != 0 {
							t.Fatalf("missing resource returned %d", result)
						}
					} else {
						rootOff, modalOff, base := 1321228, 1321232, 1321236
						if menu {
							rootOff, modalOff, base = 1522604, 1522612, 1522616
						}
						root := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[rootOff])))
						modal := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[modalOff])))
						if root == nil || modal == nil {
							t.Fatalf("constructor returned %d root=%p modal=%p", result, root, modal)
						}
						r.RootPos = [2]int{root.Off.X, root.Off.Y}
						r.ModalPos = [2]int{modal.Off.X, modal.Off.Y}
						r.Hidden = root.Flags.IsHidden()
						r.ModalHidden = modal.Flags.IsHidden()
						for i := 0; i < 4; i++ {
							w := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[base+4*i])))
							r.IDs = append(r.IDs, uint32(w.ID()))
							r.Counts = append(r.Counts, (*gui.ScrollListBoxData)(w.WidgetData).Field_11_0)
						}
						if result != 1 || !r.ModalHidden {
							t.Fatalf("constructor result %+v", r)
						}
						if menu {
							a := gui.FindAnimForStateID(900)
							if a == nil {
								t.Fatal("missing animation")
							}
							r.AnimState = int(a.State())
							a.Free()
						} else {
							if !r.Hidden || r.RootPos[0] != bindingCenterX(size[0], 640) || r.ModalPos[0] != (size[0]-300)/2 {
								t.Fatalf("in-game geometry %+v", r)
							}
							legacy.PortTestBindingDestroy()
							for _, n := range []int{1321228, 1321232, 1321236, 1321240, 1321244, 1321248} {
								if *o.words[n] != 0 {
									t.Fatalf("destroy retained %d", n)
								}
							}
						}
					}
					records = append(records, r)
				})
			}
		}
	}
	spellbookCapture(t, "binding-construction", records, "ae635dff6390b90fd34d71f58863d32ee8f1f1e4745f1825f72c489bafa6f445")
}

// The C width subtraction is unsigned. Window.SetPos orders wrapped endpoints.
func bindingCenterX(screen, width int) int {
	x := int32((uint32(screen) - uint32(width)) / 2)
	end := x + int32(width)
	if end < x {
		return int(end)
	}
	return int(x)
}
