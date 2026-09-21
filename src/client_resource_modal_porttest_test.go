//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientResourceModalWindow(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestClientResourceWords()
	t.Cleanup(restore)
	dimensions, restore := legacy.PortTestChatBubbleRenderGlobals()
	t.Cleanup(restore)
	colors, restore := legacy.PortTestInventoryDisplayWords()
	t.Cleanup(restore)
	*words["modalWindow"] = 0
	t.Cleanup(func() { legacy.Nox_xxx_gui_43E1A0(0); o.c.GUI.FreeDestroyed() })
	type row struct {
		Size                  image.Point
		Color                 uint32
		Step, Command, Window int
		Flags                 uint32
		Position, ActualSize  image.Point
		Background            uint32
	}
	var rows []row
	for _, size := range []image.Point{{1, 1}, {96, 64}, {640, 480}, {1024, 768}} {
		for _, color := range []uint32{0, 0x80000000, 0x12345678} {
			*dimensions["width"], *dimensions["height"] = uint32(size.X), uint32(size.Y)
			*colors["nox_color_black_2650656"] = color
			created := 0
			var previous *gui.Window
			for step, command := range []int{0, 1, 1, 0, 0} {
				legacy.Nox_xxx_gui_43E1A0(command)
				r := row{Size: size, Color: color, Step: step, Command: command}
				p := *words["modalWindow"]
				if command == 0 {
					if p != 0 {
						t.Fatal("modal disposal kept pointer")
					}
					if previous != nil && !previous.Flags.Has(gui.StatusDestroyed) {
						t.Fatal("modal disposal did not queue destruction")
					}
				} else {
					if p == 0 {
						t.Fatal("modal allocation")
					}
					win := (*gui.Window)(unsafe.Pointer(uintptr(p)))
					if win == previous {
						t.Fatal("modal creation reused the live window")
					}
					created++
					r.Window = created
					r.Flags = uint32(win.Flags)
					r.Position = win.Off
					r.ActualSize = win.SizeVal
					r.Background = win.DrawData().BgColorVal
					if r.Flags != 552 || r.Position != (image.Point{}) || r.ActualSize != size || r.Background != color {
						t.Fatal("modal window fields", r)
					}
					previous = win
				}
				rows = append(rows, r)
			}
			o.c.GUI.FreeDestroyed()
		}
	}
	interactionCapture(t, "client-resource-modal-window", rows)
}
