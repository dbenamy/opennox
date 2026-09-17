//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestServerConfigRuleControlDrawing(t *testing.T) {
	o, w, _ := serverConfigRuleOwner(t)
	list, entry, draw := w.ChildByID(10170), w.ChildByID(10171), w.ChildByID(10176)
	list.Func94(gui.AsWindowEvent(16397, uintptr(unsafe.Pointer(alloc.InternCString16("easy"))), ^uintptr(0)))
	draw.DrawData().BgColorVal = 0x80000000
	type row struct {
		Flags                          uint32
		Selected, Focus, Text, Initial bool
		Enabled                        []bool
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 0x4000, 0x8000} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		for _, selected := range []bool{false, true} {
			for _, focus := range []bool{false, true} {
				for _, text := range []bool{false, true} {
					for _, initial := range []bool{false, true} {
						index := -1
						if selected {
							index = 0
						}
						list.Func94(gui.AsWindowEvent(16403, uintptr(index), 0))
						value := ""
						if text {
							value = "new"
						}
						serverPanelsSetText(entry, value)
						if focus {
							o.c.GUI.Focus(entry)
						} else {
							o.c.GUI.Focus(nil)
						}
						for _, id := range []uint{10172, 10173, 10174, 10175} {
							w.ChildByID(id).Flags &^= 8
							if initial {
								w.ChildByID(id).Flags |= 8
							}
						}
						draw.Draw()
						r := row{Flags: flags, Selected: selected, Focus: focus, Text: text, Initial: initial}
						for _, id := range []uint{10172, 10173, 10174, 10175} {
							want := selected
							if id == 10172 {
								want = initial
								if focus {
									want = text
								}
							} else if id == 10174 && selected {
								want = initial || flags&49152 == 0
							}
							got := w.ChildByID(id).Flags.IsEnabled()
							if got != want {
								t.Fatal("draw availability", flags, selected, focus, text, initial, id, got, want)
							}
							r.Enabled = append(r.Enabled, got)
						}
						rows = append(rows, r)
					}
				}
			}
		}
		restore()
	}
	spellbookCapture(t, "server-config-rule-control-drawing", rows, "ee0cc3d492e810e0f7361ab16a73d067db1c43f5c9abef0dc7f7ef2d31f6511b")
}
func TestServerConfigRulePixels(t *testing.T) {
	type row struct {
		Mode     int
		Position image.Point
		Pixels   string
	}
	var rows []row
	for mode := 0; mode < 3; mode++ {
		for _, pos := range []image.Point{image.Pt(5, 7), image.Pt(-3, -2), image.Pt(17, 13)} {
			t.Run(fmt.Sprintf("%d/%v", mode, pos), func(t *testing.T) {
				o, w, _ := serverConfigRuleOwner(t)
				draw := w.ChildByID(10176)
				o.options.SetPos(image.Pt(3, 2))
				draw.SetPos(pos)
				draw.SizeVal = image.Pt(21, 17)
				draw.DrawData().BgColorVal = 0x12345678
				if mode == 1 {
					draw.DrawData().BgColorVal = 0x80000000
				}
				if mode == 2 {
					draw.Flags |= gui.StatusImage
					draw.DrawData().SetBackgroundImage(o.images[3])
				}
				for i := range o.pix.Pix {
					o.pix.Pix[i] = 0xffff
				}
				before := effectsPixelHash(o.pix)
				draw.Draw()
				after := effectsPixelHash(o.pix)
				if (before != after) != (mode != 1) {
					t.Fatal("rule draw path", mode, before, after)
				}
				rows = append(rows, row{mode, pos, after})
			})
		}
	}
	spellbookCapture(t, "server-config-rule-pixels", rows, "660f0eca66c265c94e0615da2b6469110024ffd9389572be2b2482b234aa98f0")
}
func TestServerConfigRuleEventReturns(t *testing.T) {
	o, w, _ := serverConfigRuleOwner(t)
	entry := w.ChildByID(10171)
	type row struct {
		Map, Text               string
		Event, Code, ID, Result int
		Enabled                 bool
	}
	var rows []row
	for _, mapName := range []string{"arena", "a", "user"} {
		for _, text := range []string{"", "arena", "apple", "user", "new", "Ā"} {
			for _, code := range []int{0, 1, 2} {
				for _, id := range []int{10171, 99999} {
					clear(o.settings[:9])
					copy(o.settings, mapName)
					serverPanelsSetText(entry, text)
					w.ChildByID(10172).Flags |= 8
					got := gui.EventRespInt(w.Func94(gui.AsWindowEvent(16387, uintptr(code), uintptr(id))))
					want := 1
					if code == 1 || id == 99999 {
						want = 0
					}
					if got != want {
						t.Fatal("completion return", code, id, got, want)
					}
					rows = append(rows, row{mapName, text, 16387, code, id, got, w.ChildByID(10172).Flags.IsEnabled()})
				}
			}
		}
	}
	for _, event := range []int{23, 16384, 16400, 16415, 16416} {
		got := gui.EventRespInt(w.Func94(gui.AsWindowEvent(event, uintptr(entry.C()), 0)))
		if got != 1 {
			t.Fatal("ignored event result", event, got)
		}
		rows = append(rows, row{Event: event, Result: got})
	}
	spellbookCapture(t, "server-config-rule-event-returns", rows, "0bd3cef598594795b356b9631dffb8109dd02d972c4c0a5dc1cb5510ea52844e")
}
