//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
	"testing"
	"unsafe"
)

func TestServerPanelsLocaleResources(t *testing.T) {
	type row struct {
		Kind             string
		Language, Height int
		Resource         string
		Children         []uint
	}
	var rows []row
	for _, kind := range []string{"access", "general", "advserv"} {
		for lang := 0; lang < 9; lang++ {
			for _, height := range []int{10, 11} {
				t.Run(fmt.Sprintf("%s-%d-%d", kind, lang, height), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					o.installSubpanels(t, true)
					o.configureLanguage(lang)
					face := *basicfont.Face7x13
					face.Height = height
					face.Ascent = height
					t.Cleanup(o.c.Render().GetFonts().PortTestDefaultFont(&face))
					_, restore := o.c.Render().GetFonts().PortTestWindowFont(&face, "small", "large")
					t.Cleanup(restore)
					off := uintptr(127824)
					key := "panel-1045516"
					if kind == "general" {
						off = 173556
						key = "panel-1309812"
					}
					if kind == "advserv" {
						off = 180048
						key = "panel-1316972"
					}
					table := unsafe.Slice((*unsafe.Pointer)(memmap.PtrOff(0x587000, off)), 9)
					for i := range table {
						table[i] = unsafe.Pointer(alloc.InternCString(fmt.Sprintf("panel-%s-%d.wnd", kind, i)))
					}
					old := legacy.Nox_new_window_from_file
					var loaded string
					legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
						loaded = name
						return newWindowFromString(o.c.GUI, serverOptionsSubpanelResource(kind+".wnd", true), fn)
					}
					t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
					legacy.PortTestServerPanelsConstruct(kind, o.options, unsafe.Pointer(&o.settings[0]))
					index := lang
					if height > 10 {
						index = 2
					}
					if loaded != fmt.Sprintf("panel-%s-%d.wnd", kind, index) {
						t.Fatal("locale resource", loaded)
					}
					root := *o.optionWords[key]
					if root == 0 {
						t.Fatal("locale root")
					}
					w := (*gui.Window)(unsafe.Pointer(uintptr(root)))
					rows = append(rows, row{kind, lang, height, loaded, serverOptionsChildren(t, w)})
					if kind == "advserv" {
						w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(2130).C()), 0))
					}
				})
			}
		}
	}
	spellbookCapture(t, "server-panels-locale-resources", rows, "286e19dd552e820262b60540fccbda063fd0181f06248b114581b44bf3187c1e")
}
