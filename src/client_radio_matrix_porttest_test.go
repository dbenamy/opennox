//go:build porttest

package opennox

import (
	"image"
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientRadioInputAndGroups(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newRadioOwner(t)
	var out []radioResult
	type action struct {
		Input bool
		Code  int
		A, B  uint32
	}
	actions := []action{{true, 5, 0, 0}, {true, 8, 0x12345678, 0}, {true, 17, 0x12345678, 0}, {false, 23, 1, 0}, {false, 23, 0, 0}, {true, 18, 0x89abcdef, 0}, {true, 6, 0x12345678, 0}, {true, 7, 0x89abcdef, 0}, {false, 16392, 1, 0}, {true, 17, 0, 0}, {true, 7, 0, 0}, {true, 18, 0, 0}, {true, 6, 0, 0}}
	for _, key := range []uint32{15, 205, 208, 200, 203, 28, 57, 1} {
		for state := uint32(0); state < 4; state++ {
			actions = append(actions, action{true, 21, key, state})
		}
	}
	actions = append(actions, action{false, 16392, 0, 0}, action{false, 16392, 2, 0}, action{false, -1, 0xffffffff, 0xffffffff})
	id := 0
	for parent := 0; parent < 2; parent++ {
		for owner := 0; owner < 2; owner++ {
			for track := 0; track < 2; track++ {
				for mode := 0; mode < 2; mode++ {
					for _, flags := range []gui.StatusFlags{0, 8, 0x408} {
						for _, state := range []uint32{0, 2, 4, 6, 0xabcdef00} {
							id++
							o.create(t, id, mode, parent, owner, track, 0, flags, func(d *gui.WindowData) { d.Field0 = state })
							out = append(out, o.snapshot(t, id, -1, 0))
							for step, a := range actions {
								// Reset selection before each key/select action so accepted/rejected input
								// states remain discriminating after earlier mouse actions selected it.
								if a.Code == 21 || a.Code == 16392 {
									o.windows[0].DrawData().Field0 &^= 4
									for _, w := range o.windows[1:] {
										w.DrawData().Field0 |= 4
									}
								}
								ev := &gui.RawEvent{Event: a.Code, Arg1: uintptr(a.A), Arg2: uintptr(a.B)}
								var ret gui.WindowEventResp
								if a.Input {
									ret = o.windows[0].Func93(ev)
								} else {
									ret = o.windows[0].Func94(ev)
								}
								out = append(out, o.snapshot(t, id, step, gui.EventRespInt(ret)))
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "radio-input-groups", out, len(out), "a60b16d88850abb413528383633bcbfaaefe20166676cb0c930695baeebc298c")
}

func TestClientRadioDrawingAndText(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newRadioOwner(t)
	var out []radioResult
	id := 0
	for mode := 0; mode < 2; mode++ {
		for enabled := 0; enabled < 2; enabled++ {
			for center := 0; center < 2; center++ {
				for smoothing := 0; smoothing < 2; smoothing++ {
					for _, state := range []uint32{0, 2, 4, 6} {
						for mask := 0; mask < 64; mask++ {
							id++
							flags := gui.StatusFlags(enabled*8 + smoothing*0x2000)
							o.create(t, id, mode, 1, 1, 1, center, flags, func(d *gui.WindowData) {
								d.Field0 = state
								d.SetImagePoint(image.Pt(mask%3-1, mask%5-2))
								colors := []*uint32{&d.BgColorVal, &d.EnColorVal, &d.DisColorVal, &d.SelColorVal, &d.HlColorVal, &d.TextColorVal}
								images := []*noxrender.ImageHandle{&d.BgImageHnd, &d.EnImageHnd, &d.DisImageHnd, &d.SelImageHnd, &d.HlImageHnd}
								for i := range colors {
									if mask&(1<<i) != 0 {
										*colors[i] = 0x80000000
										if i < len(images) {
											*images[i] = nil
										}
									}
								}
							})
							out = append(out, o.snapshot(t, id, -1, 0))
						}
					}
				}
			}
		}
	}
	// Actual 16-bit text input: empty, embedded terminator, exact limit, truncation,
	// non-ASCII, surrogate pair and an unpaired surrogate retain raw window words.
	texts := [][]uint16{{0}, {'A', 0, 'B', 0}, {0x03a9, 0x4e2d, 0}, {0xd83d, 0xde00, 0}, {0xd800, 0}}
	for _, n := range []int{62, 63, 64, 90} {
		s := make([]uint16, n+1)
		for i := 0; i < n; i++ {
			s[i] = uint16('A' + i%26)
		}
		texts = append(texts, s)
	}
	for center := 0; center < 2; center++ {
		id++
		o.create(t, id, 0, 1, 1, 1, center, 8)
		for step, txt := range texts {
			data, free := alloc.Make(txt, len(txt))
			ret := o.windows[0].Func94(&gui.RawEvent{Event: 16385, Arg1: uintptr(unsafe.Pointer(&data[0]))})
			free()
			out = append(out, o.snapshot(t, id, step, gui.EventRespInt(ret)))
		}
	}
	effectsCapture(t, "radio-drawing-text", out, len(out), "07722865ef1a1e51f9f40370bcf7b186ae1651c0eb95d4e710339ac7c6f72b82")
}
