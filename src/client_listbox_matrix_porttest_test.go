//go:build porttest

package opennox

import (
	"image"
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/client/seat"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
)

type listboxResult struct {
	Case, Step    int
	Return        uint32
	Windows, Data [][]uint32
	Rows          []gui.ScrollListBoxItem
	Selection     []int32
	Notices       []listboxNotice
	Input         []uint32
	Named         []string
	Focus         uint32
	Render        objectRenderResult
}

func (o *listboxOwner) snapshot(t *testing.T, id, step, ret int, draw bool) listboxResult {
	r := listboxResult{Case: id, Step: step, Return: o.norm(uint32(ret))}
	windows := []*gui.Window{o.parent, o.win}
	var add func(*gui.Window)
	add = func(w *gui.Window) {
		for c := w.Field100Ptr; c != nil; c = c.Prev() {
			windows = append(windows, c)
			add(c)
		}
	}
	add(o.win)
	for i, w := range windows {
		if draw && i != 0 {
			w.Draw()
		}
		words := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
		for _, j := range []int{8, 13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
			words[j] = o.norm(words[j])
		}
		r.Windows = append(r.Windows, words)
		var data []uint32
		if i == 1 {
			data = append(data, unsafe.Slice((*uint32)(w.WidgetData), 14)...)
			for _, j := range []int{6, 7, 8, 9} {
				data[j] = o.norm(data[j])
			}
			if o.data().Field_4 != 0 {
				data[12] = o.norm(data[12])
			}
		} else if i > 1 && w.WidgetData != nil && w.DrawData().Style&(gui.StyleVertSlider|gui.StyleHorizSlider) != 0 {
			data = append(data, unsafe.Slice((*uint32)(w.WidgetData), 4)...)
		}
		r.Data = append(r.Data, data)
	}
	d := o.data()
	r.Rows = append([]gui.ScrollListBoxItem(nil), unsafe.Slice(d.Items, int(d.Count))...)
	r.Selection = o.selection()
	r.Notices = append([]listboxNotice(nil), o.notices...)
	r.Input = o.c.Inp.PortTestEntryState()
	r.Named = append([]string(nil), o.named...)
	r.Focus = o.norm(uint32(uintptr(o.c.GUI.Focused().C())))
	r.Render = o.renderResult(t, id, step, int(r.Return))
	return r
}
func listboxUnits(s string) []uint16 {
	out := make([]uint16, 0, len(s))
	for _, v := range s {
		out = append(out, uint16(v))
	}
	return out
}
func TestClientListboxRowsAndCapacity(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newListboxOwner(t)
	var out []listboxResult
	id := 0
	for _, capacity := range []uint16{1, 2, 4, 8} {
		for multi := 0; multi < 2; multi++ {
			for scroll := 0; scroll < 2; scroll++ {
				for rolling := 0; rolling < 2; rolling++ {
					for automatic := 0; automatic < 2; automatic++ {
						id++
						o.resetRender(uint32(id), 120)
						o.create(t, 8, 70, 65, func(draw *gui.WindowData, d *gui.ScrollListBoxData) {
							d.Count = capacity
							d.Field_4 = uint32(multi)
							d.Field_3 = uint32(scroll)
							d.Field_2 = uint32(rolling)
							d.Field_1 = uint32(automatic)
							if id%2 != 0 {
								draw.SetText("Rows")
							}
						})
						for i := range unsafe.Slice(o.data().Items, int(capacity)) {
							row := &unsafe.Slice(o.data().Items, int(capacity))[i]
							row.Field_0 = 0x12345678
							row.Field_129 = 0x24682468
							row.Field_130 = 0xaabbcc00
							for j := 1; j < len(row.Text); j++ {
								row.Text[j] = 0x7000 + uint16(j)
							}
						}
						step := 0
						save := func(ret int) { out = append(out, o.snapshot(t, id, step, ret, false)); step++; o.notices = nil }
						save(0)
						for i := 0; i < int(capacity)+2; i++ {
							var text []uint16
							switch i % 4 {
							case 0:
								text = []uint16{'A', uint16('0' + i)}
							case 1:
								text = []uint16{}
							case 2:
								text = nil
							case 3:
								text = []uint16{'x', 0xd800, '\n'}
							}
							save(o.textEvent(16397, text, int32(i%19-1)))
						}
						save(o.event(16402, 1, 0))
						save(o.textEvent(16397, listboxUnits("middle"), -1))
						save(o.event(16403, 0, 0))
						save(o.event(16405, uint32(capacity-1), 0))
						save(o.textEvent(16407, listboxUnits("replaced"), 0))
						save(o.event(16406, 0, 0))
						save(o.event(16404, 0, 0))
						save(o.event(16411, 1, 0))
						save(o.event(16398, 0, 0))
						save(o.event(16399, 1, 0))
						save(o.textEvent(16397, listboxUnits("restart"), 16))
						save(o.event(16399, 0, 0))
					}
				}
			}
		}
	}
	effectsCapture(t, "listbox-rows-capacity", out, len(out), "13988d7aed540c47474782d208ef6378799a91706fa15de7919d6f7ae6779c52")
}
func TestClientListboxInputAndSelection(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newListboxOwner(t)
	var out []listboxResult
	id := 0
	for multi := 0; multi < 2; multi++ {
		for sticky := 0; sticky < 2; sticky++ {
			for scroll := 0; scroll < 2; scroll++ {
				for modifier := 0; modifier < 4; modifier++ {
					for _, flags := range []gui.StatusFlags{0, 8, 0x408} {
						id++
						o.resetRender(uint32(id), 120)
						o.create(t, flags, 70, 65, func(draw *gui.WindowData, d *gui.ScrollListBoxData) {
							d.Field_4 = uint32(multi)
							d.Field_5 = uint32(sticky)
							d.Field_3 = uint32(scroll)
							if id%2 != 0 {
								draw.SetText("Rows")
							}
						})
						for _, s := range []string{"one", "two", "three", "four", "five"} {
							o.textEvent(16397, listboxUnits(s), -1)
						}
						if modifier&1 != 0 {
							o.seat.send(&seat.KeyboardEvent{Key: keybind.Key(29), Pressed: true})
						}
						if modifier&2 != 0 {
							o.seat.send(&seat.KeyboardEvent{Key: keybind.Key(42), Pressed: true})
						}
						if modifier != 0 {
							o.c.Inp.Tick()
						}
						step := 0
						save := func(ret int) { out = append(out, o.snapshot(t, id, step, ret, false)); step++; o.notices = nil }
						save(0)
						for _, a := range []struct {
							code int
							y    uint32
						}{{5, 16}, {17, 16}, {6, 16}, {7, 30}, {6, 30}, {7, 43}, {6, 65}, {7, 0xffff}, {8, 16}, {10, 16}, {11, 16}, {18, 16}, {19, 16}, {20, 16}} {
							save(gui.EventRespInt(o.win.Func93(&gui.RawEvent{Event: a.code, Arg1: uintptr(20 | a.y<<16)})))
						}
						for _, key := range []uint32{1, 15, 28, 57, 200, 203, 205, 208, 199, 201, 207, 209, 0x100d0} {
							for state := uint32(0); state < 4; state++ {
								o.event(16403, 2, 0)
								o.notices = nil
								save(gui.EventRespInt(o.win.Func93(&gui.RawEvent{Event: 21, Arg1: uintptr(key), Arg2: uintptr(state)})))
							}
						}
						for _, index := range []uint32{0xffffffff, 0, 1, 4, 7, 8, 0x7fffffff} {
							save(o.event(16403, index, 0))
							save(o.event(16405, index, 0))
							save(o.event(16404, 0, 0))
						}
						o.c.GUI.Focus(o.win)
						save(0)
						o.c.GUI.Focus(nil)
						save(0)
					}
				}
			}
		}
	}
	effectsCapture(t, "listbox-input-selection", out, len(out), "7f7fb19d91ecc15a68e57cd7560b77480674a16220aa5380bea97dcc7bc71a28")
}
func TestClientListboxDrawing(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newListboxOwner(t)
	var out []listboxResult
	id := 0
	for mode := 0; mode < 2; mode++ {
		for multi := 0; multi < 2; multi++ {
			for enabled := 0; enabled < 2; enabled++ {
				for smoothing := 0; smoothing < 2; smoothing++ {
					for label := 0; label < 2; label++ {
						for mask := 0; mask < 64; mask++ {
							id++
							o.resetRender(uint32(id), 120)
							size := []image.Point{image.Pt(70, 65), image.Pt(28, 40), image.Pt(50, 20)}[id%3]
							flags := gui.StatusFlags(mode*128 + enabled*8 + smoothing*0x2000)
							if id%2 != 0 {
								flags |= 0x4000
							}
							o.create(t, flags, size.X, size.Y, func(draw *gui.WindowData, d *gui.ScrollListBoxData) {
								d.Field_4 = uint32(multi)
								d.Field_3 = uint32(id % 2)
								if label != 0 {
									draw.SetText("Rows")
								}
								if id%3 == 0 {
									draw.Field0 = 2
								}
								colors := []*uint32{&draw.BgColorVal, &draw.EnColorVal, &draw.HlColorVal, &draw.DisColorVal, &draw.SelColorVal, &draw.TextColorVal}
								images := []*noxrender.ImageHandle{&draw.BgImageHnd, &draw.EnImageHnd, &draw.HlImageHnd, &draw.DisImageHnd, &draw.SelImageHnd}
								for i, p := range colors {
									if mask&(1<<i) != 0 {
										*p = 0x80000000
										if i < len(images) {
											*images[i] = nil
										}
									}
								}
							})
							for i, s := range []string{"one", "wrapped row text", "three", "four"} {
								o.textEvent(16397, listboxUnits(s), int32((id+i)%19-1))
							}
							o.event(16403, 2, 0)
							if multi != 0 {
								o.event(16405, 0, 0)
							}
							out = append(out, o.snapshot(t, id, 0, 0, true))
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "listbox-drawing", out, len(out), "a4966f318753a7b35313e7f8e022b0b2217191cf8bd3940079218ccb7ab13724")
}

func TestClientListboxControlsAndBoundaries(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newListboxOwner(t)
	var out []listboxResult
	id := 0
	for _, lineHeight := range []uint16{0, 10, 13, 255, 65535} {
		for mode := 0; mode < 2; mode++ {
			for multi := 0; multi < 2; multi++ {
				for scroll := 0; scroll < 2; scroll++ {
					id++
					o.resetRender(uint32(id), 120)
					o.create(t, gui.StatusFlags(8|mode*128), 70, 65, func(draw *gui.WindowData, d *gui.ScrollListBoxData) {
						d.Line_height = lineHeight
						d.Field_4 = uint32(multi * 7)
						d.Field_3 = uint32(scroll * 3)
						if id%3 == 0 {
							draw.SetText("Label")
						}
					})
					step := 0
					save := func(ret int) { out = append(out, o.snapshot(t, id, step, ret, false)); step++; o.notices = nil }
					save(0)
					for _, text := range []string{"one", "two wrapped words", "three", "four", "five"} {
						save(o.textEvent(16397, listboxUnits(text), -1))
					}
					for _, index := range []uint32{0xffffffff, 0, 1, 4, 5, 8, 65535, 0x7fffffff} {
						save(o.event(16402, index, 0))
						save(o.event(16412, index, 0))
					}
					if scroll != 0 {
						for _, value := range []uint32{0, 1, 20, 64, 255, 65535, 0xffffffff} {
							save(o.event(16393, uint32(uintptr(o.data().Field_9)), value))
						}
						for _, code := range []int{16391, 16384} {
							for _, button := range []unsafe.Pointer{o.data().Field_7, o.data().Field_8, o.parent.C()} {
								save(o.event(code, uint32(uintptr(button)), 0))
							}
						}
					}
					for _, size := range [][2]uint32{{70, 65}, {28, 40}, {80, 90}} {
						save(o.event(16388, size[0], size[1]))
					}
					for _, length := range []int{0, 1, 63, 64, 255, 256, 400} {
						text := make([]uint16, length)
						for i := range text {
							text[i] = uint16('a' + i%26)
						}
						save(o.textEvent(16385, text, 0))
						save(o.textEvent(16407, text, 2))
					}
				}
			}
		}
	}
	effectsCapture(t, "listbox-controls-boundaries", out, len(out), "f5bd9dc52834c2edb335d1e8d46e8b26d115a07f8699e62cf0ae6c52a0fb69d7")
}
