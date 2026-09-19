//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"testing"
	"unsafe"
)

func TestClientAudioEventsRefill(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	cases := [][]int{nil, {0}, {1}, {16383}, {16384}, {16385}, {32768}, {1, 2, 3}, {8192, 8192}, {8191, 8194}, {16383, 2, 17}, {0, 3, 0, 7}, {16000, 16000, 8000}}
	var rows []map[string]any
	for index, lengths := range cases {
		for _, direct := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/%t", index, direct), func(t *testing.T) {
				o := newAudioStreamOwner(t, 1, 1)
				if o.device() == 0 {
					t.Fatal("device")
				}
				ctx := o.call("sub_487150", 0, 0)
				v := o.call("sub_487750", ctx)
				vw := audioStreamWords(v, 78)
				total := 0
				for _, n := range lengths {
					total += n
				}
				data, free := alloc.Make([]byte{}, total+1)
				defer free()
				for i := 0; i < total; i++ {
					data[i] = byte((i*37 + 11) % 251)
				}
				buffer, free := alloc.Make([]uint32{}, 7)
				defer free()
				bp := audioStreamPointer(unsafe.Pointer(&buffer[0]))
				o.call("sub_487C30", bp)
				chunks, free := alloc.Make([]uint32{}, 6*(len(lengths)+1))
				defer free()
				if direct {
					buffer[0], buffer[1] = audioStreamPointer(unsafe.Pointer(&data[0])), uint32(total)
				} else {
					off := 0
					for i, n := range lengths {
						chunk := audioStreamPointer(unsafe.Pointer(&chunks[6*i]))
						o.call("sub_487D30", chunk, audioStreamPointer(unsafe.Pointer(&data[off])), uint32(n))
						o.call("sub_487C50", bp, chunk)
						off += n
					}
				}
				o.call("sub_4BDB90", v, bp)
				vw[31] = 1
				vw[36] = audioStreamPointer(o.callbacks[12])
				sample, free := alloc.New(legacy.AudioSample{})
				defer free()
				sample.Smp = ail.Sample(handles.New())
				sample.Field1 = unsafe.Pointer(uintptr(v))
				scratch, free := alloc.Make([]byte{}, 2*(16384+32))
				defer free()
				for i := range scratch {
					scratch[i] = 0xcc
				}
				sample.Data1 = &scratch[16]
				sample.Data2 = &scratch[16384+48]
				ready := []int{0, 1, 0, -1}
				calls := 0
				type load struct {
					Buffer  uint32
					Bytes   int
					Storage string
					SHA     string
				}
				var loads []load
				var output []byte
				restore := legacy.PortTestAudioEventDevice(sample.Smp, sample, func() int {
					calls++
					if calls > len(ready) {
						return -1
					}
					return ready[calls-1]
				}, func(n uint32, b []byte) {
					where := "empty"
					if len(b) > 0 {
						where = "source"
						if unsafe.Pointer(&b[0]) == unsafe.Pointer(sample.Data1) {
							where = "scratch0"
						}
						if unsafe.Pointer(&b[0]) == unsafe.Pointer(sample.Data2) {
							where = "scratch1"
						}
					}
					loads = append(loads, load{n, len(b), where, fmt.Sprintf("%x", sha256.Sum256(b))})
					output = append(output, b...)
				})
				defer restore()
				ret := legacy.PortTestAudioEventCall("sub_43EE00", audioEventPointer(unsafe.Pointer(sample)))
				expectedTotal := total
				if !direct {
					expectedTotal = 0
					for i, n := range lengths {
						if i > 0 && n == 0 {
							break
						}
						expectedTotal += n
					}
				}
				if ret != audioEventSigned(-1) || calls != len(ready) || len(loads) != 3 || !bytes.Equal(output, data[:expectedTotal]) {
					t.Fatal("refill output/readiness", ret, calls, len(loads), len(output), expectedTotal)
				}
				for _, span := range [][2]int{{0, 16}, {16400, 16432}, {32816, 32832}} {
					for _, b := range scratch[span[0]:span[1]] {
						if b != 0xcc {
							t.Fatal("scratch boundary")
						}
					}
				}
				if sample.Field3 != 1 || vw[75] != 0 {
					t.Fatal("stream end state")
				}
				before := o.callbackCounts[12]
				legacy.PortTestAudioEventCall("sub_43EDB0", uint64(sample.Smp))
				legacy.PortTestAudioEventCall("sub_43EDB0", uint64(sample.Smp))
				if sample.Flag7 != 1 || o.callbackCounts[12] != before+1 || vw[31]&5 != 0 || vw[72] != 0 {
					t.Fatal("single end callback")
				}
				rows = append(rows, map[string]any{"lengths": lengths, "direct": direct, "loads": loads, "ready_calls": calls, "ended": sample.Field3})
			})
		}
	}
	spellbookCapture(t, "client-audio-events-refill", rows, "d86f646a4552e05e0b79022da44102c82a9f8b5d1ec5b52d4e6b27ecf774e44b")
}
