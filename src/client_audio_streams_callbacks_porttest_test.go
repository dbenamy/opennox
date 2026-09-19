//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsVoiceCallbacks(t *testing.T) {
	var rows []map[string]any
	for _, flags := range []uint32{0, 1, 4, 5, 0x10, 0x123456f0, 0xffffffff} {
		for _, fail := range []int{-1, 6, 10, 11, 12, 13} {
			t.Run(fmt.Sprintf("%x/%d", flags, fail), func(t *testing.T) {
				o := newAudioStreamOwner(t, 1, 1)
				if o.device() == 0 {
					t.Fatal("device")
				}
				ctx := o.call("sub_487150", 0, 0)
				v := o.call("sub_487750", ctx)
				w := audioStreamWords(v, 78)
				buffer, free := alloc.Make([]uint32{}, 7)
				defer free()
				bp := audioStreamPointer(unsafe.Pointer(&buffer[0]))
				o.call("sub_487C30", bp)
				buffer[0], buffer[1] = 0x12340000, 15
				for i := 0; i < 4; i++ {
					w[34+i] = audioStreamPointer(o.callbacks[10+i])
				}
				o.failOp = fail
				w[31] = flags
				if o.call("sub_4BDB20", v) != v || w[31] != flags|0x10 {
					t.Fatal("reserve voice")
				}
				if o.call("sub_4BDB30", v) != v || w[31] != (flags|0x10)&^0x10 {
					t.Fatal("unreserve voice")
				}
				w[31] = flags
				got := o.call("sub_4BDB40", v)
				want := uint32(0x80070000)
				if flags&5 != 0 {
					want = 0x800f0000
				}
				if got != want || o.callbackCounts[6] != 0 {
					t.Fatalf("start without buffer: %x want %x", got, want)
				}
				o.call("sub_4BDB90", v, bp)
				got = o.call("sub_4BDB40", v)
				want = 0
				if flags&5 != 0 {
					want = 0x800f0000
				} else if fail == 6 {
					want = 0xffffcfc7
				}
				if got != want {
					t.Fatalf("start result: %x want %x", got, want)
				}
				expectedFlags := flags
				if got == 0 {
					expectedFlags |= 1
				}
				if w[31] != expectedFlags {
					t.Fatal("start flags")
				}
				dataRet := o.call("sub_4BD8C0", v)
				if fail == 10 {
					if dataRet != 0xffffcfc7 || w[74] != 0 || w[75] != 0 || w[76] != 0 {
						t.Fatal("data callback failure")
					}
				} else if dataRet != 0 || w[74] != buffer[0] || w[75] != 15 || w[76] != 15 {
					t.Fatal("data callback success")
				}
				o.call("sub_4BDB90", v, bp)
				beforeStop := o.callbackCounts[7]
				stopRet := o.call("sub_4BDA80", v)
				wantStop := uint32(0)
				if fail == 13 {
					wantStop = 0xffffcfc7
				}
				driverStops := 0
				if expectedFlags&5 != 0 {
					driverStops = 1
				}
				if stopRet != wantStop || w[72] != 0 || w[31] != expectedFlags || o.callbackCounts[7]-beforeStop != driverStops {
					t.Fatal("stop contract")
				}
				var loops []uint32
				for _, n := range []uint32{0, 1, 2, 0xffffffff} {
					o.call("sub_4BDB90", v, bp)
					w[32] = n
					before := o.callbackCounts[9]
					if o.call("sub_4BD940", v) != 0 {
						t.Fatal("loop return")
					}
					want := n
					if n != 0 && n != 0xffffffff {
						want--
					}
					if w[32] != want || (w[72] != 0) != (n != 0) || o.callbackCounts[9]-before != boolAudioStream(n != 0) {
						t.Fatal("loop state", n)
					}
					loops = append(loops, w[32])
				}
				w[31] = flags
				w[32] = 2
				o.call("sub_4BDB90", v, bp)
				endRet := o.call("sub_4BD9B0", v)
				wantEnd := uint32(0)
				if fail == 12 {
					wantEnd = 0xffffcfc7
				}
				if endRet != wantEnd || w[31] != flags&^5 || w[32] != 0 || w[72] != 0 {
					t.Fatal("end contract")
				}
				rows = append(rows, map[string]any{"flags": flags, "failure": fail, "start": got, "data": dataRet, "stop": stopRet, "end": endRet, "loops": loops, "events": append([]audioStreamCallback(nil), o.events...)})
				w[31] = 0
				w[37] = 0
			})
		}
	}
	spellbookCapture(t, "client-audio-streams-voice-callbacks", rows, "c943b91267474994193c301ed0591041d011000c61eb4119ef82fcc586c9bca9")
}
func boolAudioStream(v bool) int {
	if v {
		return 1
	}
	return 0
}
