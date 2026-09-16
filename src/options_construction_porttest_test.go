//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
)

func optionsResource() string {
	var b strings.Builder
	b.WriteString("FONT = small;\nWINDOW\n 300 0 0 640 480 USER;\n STATUS = ENABLED;\n CHILD\n")
	for _, id := range []int{310, 311, 312, 313, 314, 320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334} {
		group := 2
		if id >= 320 {
			group = 1
		}
		if id >= 331 {
			group = 3
		}
		if id >= 333 {
			group = 4
		}
		fmt.Fprintf(&b, " WINDOW\n %d 10 10 50 20 RADIOBUTTON;\n STATUS = ENABLED;\n GROUP = %d;\n DATA = 1;\n END\n", id, group)
	}
	for _, id := range []int{351, 352, 353} {
		fmt.Fprintf(&b, " WINDOW\n %d %d 100 24 145 VERTSLIDER;\n STATUS = ENABLED;\n DATA = 0 100;\n END\n", id, 490+(id-351)*32)
	}
	for _, id := range []int{361, 362, 363} {
		fmt.Fprintf(&b, " WINDOW\n %d 10 260 24 20 CHECKBOX;\n STATUS = ENABLED;\n END\n", id)
	}
	for _, id := range []int{341, 371} {
		fmt.Fprintf(&b, " WINDOW\n %d 10 320 40 20 PUSHBUTTON;\n STATUS = ENABLED;\n END\n", id)
	}
	b.WriteString(" END\nEND\n")
	return b.String()
}
func prepareOptionsConstruction(t *testing.T, o *optionsAudioOwner, missing bool) {
	oldStrings := strMan
	strMan = o.c.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	_, restoreFont := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "small", "large")
	t.Cleanup(restoreFont)
	o.c.guiAdv.Init(o.c.Client)
	oldBack := guiOptsBack
	guiOptsBack = o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 640, 480, nil)
	t.Cleanup(func() { guiOptsBack = oldBack })
	for _, id := range []int{151, 152} {
		w := o.c.GUI.NewWindowRaw(guiOptsBack, 8, 0, 0, 10, 10, nil)
		w.SetID(uint(id))
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if name == "OptionsVolumeSlider" {
			return o.images[0]
		}
		if name == "OptionsVolumeSliderLit" {
			return o.images[1]
		}
		return oldLoad(name)
	}
	oldParser := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "Options.wnd" {
			t.Fatalf("unexpected resource %q", name)
		}
		if missing {
			return newWindowFromString(o.c.GUI, "", fn)
		}
		w := newWindowFromString(o.c.GUI, optionsResource(), fn)
		if w != nil {
			guiParseHook(name, w)
		}
		return w
	}
	rects := memmap.BlobByAddr(0x587000).Data[174072:174108]
	oldRects := append([]byte(nil), rects...)
	t.Cleanup(func() { copy(rects, oldRects) })
	words := unsafe.Slice((*int32)(unsafe.Pointer(&rects[0])), 9)
	copy(words, []int32{10, 20, 30, 40, 50, 60, 20, 10, -1})
	animWord := legacy.PortTestOptionsAnimWord()
	oldAnim := *animWord
	*animWord = 0
	t.Cleanup(func() { *animWord = oldAnim })
	animState := gui.AnimGlobalState()
	t.Cleanup(func() { gui.SetAnimGlobalState(animState) })
	t.Cleanup(func() {
		if a := gui.FindAnimForStateID(300); a != nil {
			a.Free()
		}
	})
}

func TestOptionsConstruction(t *testing.T) {
	type record struct {
		Menu, Missing                  bool
		Width, Height, Enabled, Return int
		Pos                            [2]int
		Hidden                         bool
		Sliders                        [3]gui.SliderData
		ThumbSize                      [3][2]int
		Checked                        [3]bool
		ClosePos                       [2]int
		CloseHidden                    bool
		Anim                           int
		Advanced                       bool
	}
	var rows []record
	for _, menu := range []bool{false, true} {
		for _, size := range [][2]int{{640, 480}, {1024, 768}, {480, 640}} {
			for _, enabled := range []int{0, 1, 2} {
				for _, missing := range []bool{false, true} {
					t.Run(fmt.Sprintf("menu%v/size%v/enabled%d/missing%v", menu, size, enabled, missing), func(t *testing.T) {
						o := newOptionsAudioOwner(t, menu)
						prepareOptionsConstruction(t, o, missing)
						dims := legacy.PortTestBindingDimensions()
						old := [2]int32{*dims[0], *dims[1]}
						t.Cleanup(func() { *dims[0], *dims[1] = old[0], old[1] })
						*dims[0], *dims[1] = int32(size[0]), int32(size[1])
						for i, off := range []int{126996, 122848, 93156} {
							*o.words[off] = uint32(enabled)
							o.timers[i].Current = uint32(1000+i*2000) << 16
						}
						r := record{Menu: menu, Missing: missing, Width: size[0], Height: size[1], Enabled: enabled, Anim: -1}
						r.Return = legacy.PortTestOptionsConstruct(menu)
						rootOff := 1309820
						if menu {
							rootOff = 1309720
						}
						root := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[rootOff])))
						if missing {
							if r.Return != 0 || root != nil {
								t.Fatalf("missing resource result %+v", r)
							}
						} else {
							if r.Return != 1 || root == nil {
								t.Fatalf("constructor failed %+v", r)
							}
							r.Pos = [2]int{root.Off.X, root.Off.Y}
							r.Hidden = root.Flags.IsHidden()
							r.Advanced = root.ChildByID(2000) != nil
							for ch := 0; ch < 3; ch++ {
								w := root.ChildByID(uint(351 + ch))
								r.Sliders[ch] = *(*gui.SliderData)(w.WidgetData)
								r.ThumbSize[ch] = [2]int{w.Field100Ptr.SizeVal.X, w.Field100Ptr.SizeVal.Y}
								r.Checked[ch] = root.ChildByID(uint(361+ch)).DrawData().Field0&4 != 0
								if r.Checked[ch] != (enabled == 1) || r.ThumbSize[ch] != [2]int{24, 20} {
									t.Fatalf("channel%d %+v", ch, r)
								}
							}
							if !r.Advanced {
								t.Fatal("missing actual advanced-video owner")
							}
							if menu {
								a := gui.FindAnimForStateID(300)
								if a == nil {
									t.Fatal("missing menu animation")
								}
								r.Anim = int(a.State())
							} else {
								close := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[1309824])))
								if close == nil {
									t.Fatal("missing close overlay")
								}
								r.ClosePos = [2]int{close.Off.X, close.Off.Y}
								r.CloseHidden = close.Flags.IsHidden()
								want := [2]int{(size[0] - 640) / 2, 0}
								if r.Pos != want || r.ClosePos != want || !r.Hidden || !r.CloseHidden {
									t.Fatalf("position/visibility %+v want=%v", r, want)
								}
							}
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "options-construction", rows, "8c91640ffe180f958fefbd8cc699c8366e4707d28e73d35b17001b2baa67733d")
}
