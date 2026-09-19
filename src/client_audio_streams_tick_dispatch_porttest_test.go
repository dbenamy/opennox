//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/timer"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsTickDispatch(t *testing.T) {
	var rows []map[string]any
	for _, flags := range []uint32{0, 1, 4} {
		for _, bound := range []bool{false, true} {
			for change := 0; change < 5; change++ {
				t.Run(fmt.Sprintf("%d/%t/%d", flags, bound, change), func(t *testing.T) {
					o := newAudioStreamOwner(t, 1, 1)
					if o.device() == 0 {
						t.Fatal("device")
					}
					ctx := o.call("sub_487150", 0, 0)
					v := o.call("sub_487750", ctx)
					cw, vw := audioStreamWords(ctx, 66), audioStreamWords(v, 78)
					extras, free := alloc.Make([]timer.TimerGroup{}, 2)
					defer free()
					for i := range extras {
						extras[i].Init()
					}
					own := (*timer.TimerGroup)(unsafe.Pointer(uintptr(v + 16)))
					parent := (*timer.TimerGroup)(unsafe.Pointer(uintptr(ctx + 88)))
					groups := []*timer.TimerGroup{own, parent, &extras[0], &extras[1]}
					for _, g := range groups {
						g.ClearUpdated()
					}
					vw[28] = audioStreamPointer(unsafe.Pointer(&extras[0]))
					vw[29] = audioStreamPointer(unsafe.Pointer(&extras[1]))
					vw[31] = flags
					buffer, free := alloc.Make([]uint32{}, 7)
					defer free()
					bp := audioStreamPointer(unsafe.Pointer(&buffer[0]))
					o.call("sub_487C30", bp)
					if bound {
						o.call("sub_4BDB90", v, bp)
					}
					if change != 0 {
						g := groups[change-1]
						g.Timers[1].SetRaw(50)
						// External groups have their own service owner; the audio context observes their updated flags.
						if change >= 3 {
							g.Update()
						}
					}
					cw[56], cw[57] = 0, 0
					before := o.callbackCounts[8]
					if o.call("sub_4873C0", ctx) != 0 {
						t.Fatal("tick")
					}
					want := 0
					if flags&1 != 0 && bound && change != 0 {
						want = 1
					}
					if o.callbackCounts[8]-before != want || parent.IsUpdated() {
						t.Fatal("conditional driver update", o.callbackCounts[8]-before, want)
					}
					rows = append(rows, map[string]any{"flags": flags, "bound": bound, "change": change, "updates": want})
					vw[28], vw[29], vw[72], vw[31] = 0, 0, 0, 0
				})
			}
		}
	}
	spellbookCapture(t, "client-audio-streams-tick-dispatch", rows, "d783f5cc3d55de1770cf2c229589e1fa581b3e1c5eba60308f4ea1fde03e45ef")
}
