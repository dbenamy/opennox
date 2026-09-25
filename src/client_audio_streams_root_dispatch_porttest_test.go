//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// These contracts enter through the root service functions, including their raw
// callback fields, rather than calling the legacy fixture exports directly.
func TestClientAudioStreamsRootService(t *testing.T) {
	o := newAudioStreamOwner(t, 2, 1)
	if o.device() == 0 {
		t.Fatal("device")
	}
	a, b := o.call("sub_487150", 0, 0), o.call("sub_487150", 1, 0)
	if a == 0 || b == 0 || a == b {
		t.Fatal("contexts", a, b)
	}
	aw, bw := audioStreamWords(a, 66), audioStreamWords(b, 66)
	nativeA, nativeB := aw[54], bw[54]
	defer func() { aw[54], bw[54], aw[3], bw[3], aw[53], o.root[6] = nativeA, nativeB, 0, 0, 0, 0 }()
	aw[56], aw[57], bw[56], bw[57] = 0, 0, 0, 0
	aw[54], bw[54] = audioStreamPointer(o.callbacks[10]), audioStreamPointer(o.callbacks[10])
	o.events = nil
	sub_486EF0()
	if len(o.events) != 2 || o.events[0] != (audioStreamCallback{10, o.refs[a]}) || o.events[1] != (audioStreamCallback{10, o.refs[b]}) {
		t.Fatalf("ordered callbacks: %+v", o.events)
	}
	for _, guard := range []string{"disabled", "mutating", "context"} {
		t.Run(guard, func(t *testing.T) {
			before := o.callbackCounts[10]
			switch guard {
			case "disabled":
				dword_5d4594_1193336 = 0
			case "mutating":
				o.root[6] = 1
			case "context":
				aw[3], bw[3] = 2, 2
			}
			sub_486EF0()
			if o.callbackCounts[10] != before {
				t.Fatal("guard dispatched callback")
			}
			dword_5d4594_1193336, o.root[6], aw[3], bw[3] = 1, 0, 0, 0
		})
	}
	aw[54] = nativeA
	aw[3], bw[3] = 1, 2 // Only bit 2 suppresses a context.
	o.ticks = 137
	sub_486EF0()
	if aw[62] != 137 || aw[63] != 0 || aw[58] != 137 || aw[59] != 0 {
		t.Fatal("native tick timing", aw[58:64])
	}
	aw[53] = 1
	o.ticks = 200
	sub_486EF0()
	if aw[62] != 137 {
		t.Fatal("mutating native context advanced")
	}
	aw[53], aw[56] = 0, 100
	sub_486EF0()
	if aw[62] != 137 {
		t.Fatal("period guard advanced")
	}
	aw[56] = 0
	sub_486EF0()
	if aw[62] != 200 || aw[58] != 63 {
		t.Fatal("native tick delta", aw[58:64])
	}
}

func TestClientAudioStreamsRootCompletion(t *testing.T) {
	for _, native := range []bool{false, true} {
		name := "foreign"
		if native {
			name = "native"
		}
		t.Run(name, func(t *testing.T) {
			o := newAudioStreamOwner(t, 1, 1)
			if o.device() == 0 {
				t.Fatal("device")
			}
			ctx := o.call("sub_487150", 0, 0)
			v := o.call("sub_487750", ctx)
			if v == 0 {
				t.Fatal("voice")
			}
			vw := audioStreamWords(v, 78)
			sample, free := alloc.New(legacy.AudioSample{})
			defer free()
			buffer, freeBuffer := alloc.Make([]uint32{}, 7)
			defer freeBuffer()
			defer func() { vw[68], vw[72], vw[31], vw[36] = 0, 0, 0, 0 }()
			sample.Field1 = unsafe.Pointer(uintptr(v))
			vw[68] = audioStreamPointer(unsafe.Pointer(sample))
			vw[72] = audioStreamPointer(unsafe.Pointer(&buffer[0]))
			vw[31], vw[32], vw[36] = 13, 7, audioStreamPointer(o.callbacks[12])
			op := 12
			if !native {
				op = 10
				vw[71] = audioStreamPointer(o.callbacks[op])
			}
			o.failOp = op // A nonzero callback result is deliberately discarded.
			before := o.callbackCounts[op]
			if got := sub_43EFD0(unsafe.Pointer(uintptr(v))); got != 0 {
				t.Fatal("completion result", got)
			}
			if sample.Flag7 != 1 || o.callbackCounts[op] != before+1 {
				t.Fatal("completion notification", sample.Flag7, o.callbackCounts)
			}
			if native {
				if vw[72] != 0 || vw[31] != 8 || vw[32] != 0 {
					t.Fatal("native completion state", vw[31:33], vw[72])
				}
			} else if vw[72] == 0 || vw[31] != 13 || vw[32] != 7 {
				t.Fatal("foreign callback state changed")
			}
			if sub_43EFD0(unsafe.Pointer(uintptr(v))) != 0 || o.callbackCounts[op] != before+1 {
				t.Fatal("completion was not idempotent")
			}
		})
	}
}
