//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/timer"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsTimerMix(t *testing.T) {
	var rows []map[string]any
	for mask := 0; mask < 8; mask++ {
		t.Run(fmt.Sprint(mask), func(t *testing.T) {
			o := newAudioStreamOwner(t, 1, 1)
			if o.device() == 0 {
				t.Fatal("device")
			}
			ctx := o.call("sub_487150", 0, 0)
			v := o.call("sub_487750", ctx)
			cw, vw := audioStreamWords(ctx, 66), audioStreamWords(v, 78)
			extras, free := alloc.Make([]timer.TimerGroup{}, 3)
			defer free()
			group := func(p uint32) *timer.TimerGroup { return (*timer.TimerGroup)(unsafe.Pointer(uintptr(p))) }
			own, parent, effective := group(v+16), group(ctx+88), group(v+176)
			values := [][3]uint32{{8192, 80, 7000}, {12288, 60, 10000}, {4096, 50, 16000}, {14000, 90, 1000}, {10000, 70, 13000}}
			groups := []*timer.TimerGroup{own, &extras[0], &extras[1], parent, &extras[2]}
			for i, g := range groups {
				g.Init()
				for j, n := range values[i] {
					g.Timers[j].SetRaw(n)
					g.Timers[j].Update()
				}
			}
			vw[28], vw[29], cw[46] = 0, 0, 0
			if mask&1 != 0 {
				vw[28] = audioStreamPointer(unsafe.Pointer(&extras[0]))
			}
			if mask&2 != 0 {
				vw[29] = audioStreamPointer(unsafe.Pointer(&extras[1]))
			}
			if mask&4 != 0 {
				cw[46] = audioStreamPointer(unsafe.Pointer(&extras[2]))
			}
			want := [3]uint32{0x4000, 100, 0x2000}
			for i, value := range values {
				if i == 1 && mask&1 == 0 || i == 2 && mask&2 == 0 || i == 4 && mask&4 == 0 {
					continue
				}
				want[0] = ((want[0] * value[0]) / 0x4000) & 0xffff
				want[1] = ((want[1] * value[1]) / 100) & 0xffff
				if value[2] != 0x2000 {
					n := int32(want[2] + value[2] - 0x2000)
					if n < 0 {
						n = 0
					}
					if n > 0x4000 {
						n = 0x4000
					}
					want[2] = uint32(n)
				}
			}
			o.call("sub_4BD840", v)
			var got [3]uint32
			for i := range got {
				got[i] = effective.Timers[i].Current >> 16
			}
			if got != want || own.IsUpdated() || (mask&1 != 0 && extras[0].IsUpdated()) || !parent.IsUpdated() || !extras[1].IsUpdated() || !extras[2].IsUpdated() {
				t.Fatal("timer composition/updated flags", got, want)
			}
			rows = append(rows, map[string]any{"mask": mask, "values": got})
			vw[28], vw[29], cw[46] = 0, 0, 0
		})
	}
	spellbookCapture(t, "client-audio-streams-timer-mix", rows, "017bd63d5ed700098509c8d657443f6c52c093b0e0ca8260904a7adc24eff32e")
}
