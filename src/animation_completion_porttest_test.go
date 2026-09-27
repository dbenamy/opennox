//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestAnimationInCompletionOrdering(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, callback := range []bool{false, true} {
			for _, free := range []bool{false, true} {
				if free && !callback {
					continue
				}
				t.Run(fmt.Sprintf("reverse=%t/callback=%t/free=%t", reverse, callback, free), func(t *testing.T) {
					g := gui.New(nil)
					t.Cleanup(g.DestroyAll)
					t.Cleanup(gui.PortTestOwnAnimations())
					oldBg := gui.MainBg
					t.Cleanup(func() { gui.MainBg = oldBg })
					var events []string
					bg := g.NewWindowRaw(nil, 8, 0, 0, 10, 10, func(_ *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
						if e.EventCode() == 23 {
							events = append(events, "focus")
							return gui.RawEventResp(1)
						}
						return nil
					})
					gui.MainBg = bg
					w := g.NewWindowRaw(nil, 8, 0, 0, 12, 12, nil)
					start, end, step := image.Point{}, image.Pt(9, 13), image.Pt(4, 5)
					positions := []image.Point{image.Pt(4, 5), image.Pt(8, 10), image.Pt(9, 13)}
					if reverse {
						start, end, step = end, start, step.Mul(-1)
						positions = []image.Point{image.Pt(5, 8), image.Pt(1, 3), image.Point{}}
					}
					a := gui.NewAnim(w, end, start, step, step.Mul(-1))
					a.StateID = 0x6f005
					ready := false
					if callback {
						key, restore := legacy.PortTestAnimationCompleteIn(func() {
							events = append(events, "callback")
							ready = a.State() == gui.AnimInDone && gui.AnimGlobalState() == gui.AnimInDone && w.Offs() == end && g.Focused() != bg
							if free {
								a.Free()
							}
						})
						t.Cleanup(restore)
						a.FncDoneInPtr = key
					}
					for i, want := range positions {
						gui.AnimTick()
						if w.Offs() != want {
							t.Fatalf("tick %d position %v, want %v", i, w.Offs(), want)
						}
						if i < 2 && (len(events) != 0 || a.State() != gui.AnimIn || gui.AnimGlobalState() != gui.AnimIn) {
							t.Fatalf("early completion at tick %d: %v", i, events)
						}
					}
					want := []string{"focus"}
					if callback {
						want = []string{"callback", "focus"}
						if !ready {
							t.Fatal("callback preceded completion state/position or followed focus")
						}
					}
					if !reflect.DeepEqual(events, want) || g.Focused() != bg {
						t.Fatalf("completion order/focus: %v, want %v", events, want)
					}
					if (gui.FindAnimForStateID(0x6f005) == nil) != free {
						t.Fatal("callback-driven animation unlink")
					}
					gui.AnimTick()
					if !reflect.DeepEqual(events, want) {
						t.Fatal("completed animation invoked callback or focus again")
					}
				})
			}
		}
	}
}
