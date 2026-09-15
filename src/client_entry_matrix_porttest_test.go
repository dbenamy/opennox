//go:build porttest

package opennox

import (
	"crypto/sha256"
	"fmt"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/client/seat"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"runtime"
	"testing"
	"unsafe"
)

type entryResult struct {
	Case, Step, Return int
	Windows, Data      [][]uint32
	Items              string
	Notices            [][3]uint32
	Environment        [3]uint32
	Focus              uint32
	Input              []uint32
	Enabled            bool
	Changes            []bool
	Named              []string
	Render             *objectRenderResult
}

func (o *entryOwner) snapshot(t *testing.T, id, step, ret int, draw bool) entryResult {
	r := entryResult{Case: id, Step: step, Return: ret}
	if draw {
		o.win.Draw()
		v := o.renderResult(t, id, step, ret)
		r.Render = &v
	}
	windows := []*gui.Window{o.parent, o.win}
	data := (*gui.EntryFieldData)(o.win.WidgetData)
	var list *gui.ScrollListBoxData
	if data.Field_1048 != 0 {
		w := (*gui.Window)(unsafe.Pointer(uintptr(data.Field_1048)))
		windows = append(windows, w)
		var addChildren func(*gui.Window)
		addChildren = func(w *gui.Window) {
			for ch := w.Field100Ptr; ch != nil; ch = ch.Prev() {
				windows = append(windows, ch)
				addChildren(ch)
			}
		}
		addChildren(w)
		list = (*gui.ScrollListBoxData)(w.WidgetData)
	}
	refs := map[uint32]uint32{}
	for i, w := range windows {
		refs[uint32(uintptr(w.C()))] = 0xe1000001 + uint32(i)
		if w.WidgetData != nil {
			refs[uint32(uintptr(w.WidgetData))] = 0xe2000001 + uint32(i)
		}
	}
	if list != nil && list.Items != nil {
		refs[uint32(uintptr(unsafe.Pointer(list.Items)))] = 0xe3000001
	}
	norm := func(v uint32) uint32 {
		if v == 0 {
			return 0
		}
		if n, ok := refs[v]; ok {
			return n
		}
		if n, ok := o.c.imageRefs[v]; ok {
			return n
		}
		return v
	}
	for i, w := range windows {
		words := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
		for _, j := range []int{8, 13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
			words[j] = norm(words[j])
		}
		r.Windows = append(r.Windows, words)
		var d []uint32
		if i == 1 {
			d = append(d, unsafe.Slice((*uint32)(w.WidgetData), 264)...)
			d[262] = norm(d[262])
		} else if i == 2 && list != nil {
			d = append(d, unsafe.Slice((*uint32)(w.WidgetData), 14)...)
			for _, j := range []int{6, 7, 8, 9, 12} {
				d[j] = norm(d[j])
			}
			if list.Items != nil {
				r.Items = fmt.Sprintf("%x", sha256.Sum256(unsafe.Slice((*byte)(unsafe.Pointer(list.Items)), int(list.Count)*524)))
			}
		} else if i > 2 && w.WidgetData != nil {
			if w.DrawData().Style&(gui.StyleVertSlider|gui.StyleHorizSlider) != 0 {
				d = append(d, unsafe.Slice((*uint32)(w.WidgetData), 4)...)
			}
		}
		r.Data = append(r.Data, d)
	}
	r.Notices = append([][3]uint32(nil), o.notices...)
	r.Environment = o.env.State()
	r.Environment[1] = norm(r.Environment[1])
	r.Focus = norm(uint32(uintptr(o.c.GUI.Focused().C())))
	r.Input = o.c.Inp.PortTestEntryState()
	r.Enabled = o.seat.enabled
	r.Changes = append([]bool(nil), o.seat.changes...)
	r.Named = append([]string(nil), o.named...)
	return r
}
func TestClientEntryKeyboardMatrix(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newEntryOwner(t)
	var out []entryResult
	id := 0
	keys := []uint32{1, 14, 211, 15, 205, 208, 200, 203, 28, 156, 30, 2, 39, 57, 58, 59, 68, 87, 88, 199, 201, 207, 209, 0, 255, 256, 0xe9, 0x391, 0xff11, 0xd800, 0xffff, 0x1000e}
	for _, lang := range []int{0, 2, 6, 8} {
		for filter := 0; filter < 4; filter++ {
			for modifier := 0; modifier < 3; modifier++ {
				id++
				d := o.create(t, lang, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) {
					switch filter {
					case 1:
						d.Field_1028 = 1
					case 2:
						d.Field_1032 = 1
					case 3:
						d.Field_1036 = 1
					}
				})
				if modifier > 0 {
					key := keybind.Key(42)
					if modifier == 2 {
						key = 58
					}
					o.seat.send(&seat.KeyboardEvent{Key: key, Pressed: true})
					o.c.Inp.Tick()
				}
				o.c.GUI.Focus(o.win)
				step := 0
				for _, key := range keys {
					for state := uint32(0); state < 4; state++ {
						o.text([]uint16{'Q', 0xd800})
						d.Field_1044 = 0
						o.seat.send(&seat.TextEditEvent{Text: ""})
						ret := o.key(key, state)
						out = append(out, o.snapshot(t, id, step, ret, false))
						step++
					}
				}
			}
		}
	}
	effectsCapture(t, "entry-keyboard", out, len(out), "b280a48c1c476320fe5242caea7f56155778736f5ac430e02cb0818cd7da50c6")
}
func TestClientEntryLimitsAndComposition(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newEntryOwner(t)
	var out []entryResult
	id := 0
	for _, lang := range []int{0, 6, 8} {
		for _, limit := range []uint16{0, 1, 2, 3, 255, 256, 257, 32767, 32768, 65535} {
			for _, length := range []int{0, 1, 2, 254, 255, 256, 300} {
				id++
				o.create(t, lang, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) { d.Field_1040 = limit })
				units := make([]uint16, length)
				for i := range units {
					units[i] = uint16('A' + i%26)
				}
				if length > 1 {
					units[1] = 0xd800
				}
				o.text(units)
				out = append(out, o.snapshot(t, id, 0, 0, false))
				o.c.GUI.Focus(o.win)
				for step, key := range []uint32{30, 14, 30, 28} {
					out = append(out, o.snapshot(t, id, step+1, o.key(key, 2), false))
				}
			}
		}
	}
	for _, lang := range []int{6, 8} {
		for _, text := range []string{"", "ab", "水", "😀"} {
			for _, ch := range []uint16{0, 7, 8, 9, 10, 11, 12, 13, 'x', 0xd800, 0xffff} {
				id++
				o.create(t, lang, 8, nil)
				o.c.GUI.Focus(o.win)
				o.seat.send(&seat.TextEditEvent{Text: text})
				o.key(30, 0)
				out = append(out, o.snapshot(t, id, 0, 0, false))
				legacy.NoxInputOnChar(ch)
				out = append(out, o.snapshot(t, id, 1, 0, false))
				o.c.GUI.Focus(nil)
				out = append(out, o.snapshot(t, id, 2, 0, false))
			}
		}
	}
	effectsCapture(t, "entry-limits-composition", out, len(out), "4b8e826f3cc125a3df790b232423be6f004b0f745f8f03e89a57b2f3cea14e37")
}
func TestClientEntryDrawingMatrix(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newEntryOwner(t)
	var out []entryResult
	id := 0
	for mode := 0; mode < 2; mode++ {
		for _, lang := range []int{0, 6} {
			for _, width := range []int{5, 10, 11, 30, 70} {
				for password := 0; password < 2; password++ {
					for _, label := range []string{"", "Name"} {
						for mask := 0; mask < 8; mask++ {
							id++
							o.resetRender(uint32(id), 120)
							flags := gui.StatusFlags(8 | 0x4000 | mode*128)
							if mask&1 != 0 {
								flags &^= 8
							}
							if mask&2 != 0 {
								flags |= 0x2000
							}
							d := o.create(t, lang, flags, func(draw *gui.WindowData, d *gui.EntryFieldData) {
								draw.SetText(label)
								d.Field_1024 = uint32(password)
								if mask&4 != 0 {
									d.Field_1042 = 15
									draw.BgColorVal = 0x80000000
									draw.TextColorVal = 0x80000000
								}
							})
							o.win.SizeVal.X = width
							o.win.EndPos.X = o.win.Off.X + width
							o.text([]uint16{'a', 'b', 0xd83d, 0xde00, 'X', 0xd800, 'y', 'z'})
							if lang == 6 {
								copy(d.Text[256:], []uint16{'水', 'a', 'b', 0})
								d.Field_1052 |= 3 << 16
							}
							o.c.GUI.Focus(o.win)
							for step, blink := range []byte{0, 7, 8, 15, 255} {
								o.env.Blink(blink)
								out = append(out, o.snapshot(t, id, step, 0, true))
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "entry-drawing", out, len(out), "559d1969e894ef516eff5e7e28619ba84d1948b9d0132f651f5bef94d0968005")
}
