//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"runtime"
	"testing"
)

func TestClientSliderConstructionAndValues(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if legacy.PortTestProtectionFloatCW()&0x0f00 != 0x0200 {
		t.Fatal("slider x87 control")
	}
	o := newSliderOwner(t)
	var out []sliderResult
	id := 0
	for horizontal := 0; horizontal < 2; horizontal++ {
		for mode := 0; mode < 2; mode++ {
			for parent := 0; parent < 2; parent++ {
				for owner := 0; owner < 2; owner++ {
					for _, flags := range []gui.StatusFlags{0, 8, 0x408} {
						for _, bounds := range [][2]int32{{0, 100}, {10, 30}, {-20, 20}, {7, 7}, {100, 0}, {-1, 0}, {0, 3}} {
							for _, size := range []image.Point{image.Pt(61, 21), image.Pt(10, 10), image.Pt(17, 17)} {
								if horizontal == 0 {
									size.X, size.Y = size.Y, size.X
								}
								id++
								o.create(t, id, horizontal, mode, parent, owner, flags, uint32(bounds[0]), uint32(bounds[1]), size)
								out = append(out, o.snapshot(t, id, -1, 0))
								for step, value := range []int32{-21, -1, 0, 1, 7, 50, 100, 101, -2147483648, 2147483647} {
									ret := gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: 16394, Arg1: uintptr(uint32(value))}))
									out = append(out, o.snapshot(t, id, step, ret))
								}
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "slider-construction-values", out, len(out), "35627ac012d127b7a7b7093418b0977824b9a38dd72cd6865d210c2064ada195")
}

func TestClientSliderInputAndNotifications(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if legacy.PortTestProtectionFloatCW()&0x0f00 != 0x0200 {
		t.Fatal("slider x87 control")
	}
	o := newSliderOwner(t)
	var out []sliderResult
	id := 0
	type action struct {
		Input bool
		Code  int
		A, B  uint32
	}
	actions := []action{{true, 5, 0, 0}, {true, 8, 0, 0}, {true, 17, 0, 0}, {true, 18, 0, 0}, {false, 23, 1, 0}, {false, 23, 0, 0}, {false, 16391, 0, 0}, {false, 16394, 12, 0}}
	for _, coord := range []uint32{0, 9, 10, 15, 30, 60, 70, 95, 65535} {
		actions = append(actions, action{true, 6, coord | (coord << 16), 0}, action{true, 7, coord | (coord << 16), 0}, action{false, 16384, 0, coord | (coord << 16)})
	}
	for _, key := range []uint32{15, 200, 203, 205, 208, 1} {
		for state := uint32(0); state < 4; state++ {
			actions = append(actions, action{true, 21, key, state})
		}
	}
	actions = append(actions, action{false, 16395, 7, 7}, action{false, 16394, 7, 0}, action{false, 16384, 0, 30 | (30 << 16)}, action{false, 16395, 0, 100}, action{false, 16394, 50, 0}, action{false, -1, 0, 0})
	for horizontal := 0; horizontal < 2; horizontal++ {
		for mode := 0; mode < 2; mode++ {
			for parent := 0; parent < 2; parent++ {
				for owner := 0; owner < 2; owner++ {
					for interactive := 0; interactive < 2; interactive++ {
						for _, flags := range []gui.StatusFlags{8, 0x408} {
							for _, min := range []uint32{0, 10} {
								id++
								size := image.Pt(61, 21)
								if horizontal == 0 {
									size.X, size.Y = size.Y, size.X
								}
								o.create(t, id, horizontal, mode, parent, owner, flags, min, min+20, size)
								if interactive != 0 {
									o.win.DrawData().Style |= 0x100
								}
								for step, a := range actions {
									ev := &gui.RawEvent{Event: a.Code, Arg1: uintptr(a.A), Arg2: uintptr(a.B)}
									ret := 0
									if a.Input {
										ret = gui.EventRespInt(o.win.Func93(ev))
									} else {
										ret = gui.EventRespInt(o.win.Func94(ev))
									}
									out = append(out, o.snapshot(t, id, step, ret))
								}
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "slider-input-notifications", out, len(out), "19cdc4370e37b1a0b43d5c347011e6b51eb3c5dec3c2600b8fa2af97b40b8592")
}

// Arithmetic-only cases keep wrapped/extreme thumb coordinates out of raster APIs.
func TestClientSliderNumericState(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if legacy.PortTestProtectionFloatCW()&0x0f00 != 0x0200 {
		t.Fatal("slider x87 control")
	}
	o := newSliderOwner(t)
	var out []sliderResult
	id := 0
	bounds := [][2]int32{{-2147483648, 2147483647}, {2147483647, -2147483648}, {-1, -1}, {0, 0}, {-2147483648, -2147483648}, {2147483647, 2147483647}, {-100000, 100000}, {0, 3}, {10, 100}}
	for horizontal := 0; horizontal < 2; horizontal++ {
		for _, span := range []int{0, 1, 9, 10, 11, 17, 61, 255, 1024} {
			for _, b := range bounds {
				id++
				size := image.Pt(span, 20)
				if horizontal == 0 {
					size.X, size.Y = size.Y, size.X
				}
				o.create(t, id, horizontal, 0, 1, 1, 8, uint32(b[0]), uint32(b[1]), size)
				out = append(out, o.state(t, id, -1, 0))
				step := 0
				for _, v := range []int32{b[0], b[1], 0, 1, -1, 2147483647, -2147483648} {
					ret := gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: 16394, Arg1: uintptr(uint32(v))}))
					out = append(out, o.state(t, id, step, ret))
					step++
				}
				for _, r := range bounds {
					ret := gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: 16395, Arg1: uintptr(uint32(r[0])), Arg2: uintptr(uint32(r[1]))}))
					out = append(out, o.state(t, id, step, ret))
					step++
					ret = gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: 16394, Arg1: uintptr(uint32(r[0]))}))
					out = append(out, o.state(t, id, step, ret))
					step++
				}
			}
		}
	}
	effectsCapture(t, "slider-numeric-state", out, len(out), "c38f5e017ed06cfb30968dc7406f27dbef707e9ace552b4124713fe86dfc0032")
}
