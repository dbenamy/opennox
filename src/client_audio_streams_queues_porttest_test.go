//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsDeviceQueues(t *testing.T) {
	o := newAudioStreamOwner(t, 1, 2)
	first := o.device()
	if first == 0 {
		t.Fatal("first device")
	}
	secondDesc, free := alloc.Make([]uint32{}, 10)
	defer free()
	copy(secondDesc, o.descriptor)
	secondDesc[3] = 0
	o.slots = 2
	second := o.call("sub_486FA0", audioStreamPointer(unsafe.Pointer(&secondDesc[0])))
	if second == 0 {
		t.Fatal("second device")
	}
	// Explicit cleanup precedes freeing the second descriptor; the owner's cleanup is idempotent.
	defer func() { o.call("sub_4875F0"); o.call("sub_4870A0") }()
	out, free := alloc.Make([]uint32{}, 3)
	defer free()
	outp := audioStreamPointer(unsafe.Pointer(&out[0]))
	format, free := alloc.Make([]uint32{}, 7)
	defer free()
	fp := audioStreamPointer(unsafe.Pointer(&format[0]))
	for i := range format {
		format[i] = uint32(i)*17 + 3
	}
	expected := append([]uint32(nil), format...)
	var rows []map[string]any
	for i := 0; i < 5; i++ {
		out[0], out[1], out[2] = 0xabcdef00, 0xabcdef01, 0xabcdef02
		ret := o.call("sub_487360", uint32(i), outp, outp+4)
		device, slot := uint32(0), ^uint32(0)
		deviceID := -1
		if i == 0 {
			device, slot, deviceID = first, 0, 0
		} else if i < 3 {
			device, slot, deviceID = second, uint32(i-1), 1
		}
		wantRet := uint32(0)
		if device != 0 {
			wantRet = outp + 4
		}
		if ret != wantRet || out[0] != device || out[1] != slot || out[2] != 0xabcdef02 {
			t.Fatal("flattened device slots", i)
		}
		ctx := o.call("sub_487150", uint32(i), fp)
		if device == 0 {
			if ctx != 0 {
				t.Fatal("out of range context")
			}
		} else {
			w := audioStreamWords(ctx, 66)
			if w[5] != device || w[6] != slot || !reflect.DeepEqual(w[15:22], expected) {
				t.Fatal("context format copy")
			}
			for j := range format {
				format[j] = ^expected[j]
			}
			if o.call("sub_487150", uint32(i), fp) != ctx || !reflect.DeepEqual(w[15:22], expected) {
				t.Fatal("cached context keeps format")
			}
			if o.call("sub_487590", ctx, fp) != ctx || !reflect.DeepEqual(w[15:22], format) {
				t.Fatal("explicit format replacement")
			}
			copy(format, expected)
		}
		rows = append(rows, map[string]any{"index": i, "device": deviceID, "slot": slot})
	}
	o.call("sub_4875F0")
	o.call("sub_4870A0")
	if o.descriptor[3] != 0 || secondDesc[3] != 0 || o.callbackCounts[1] != 2 || o.callbackCounts[3] != 3 {
		t.Fatal("queue destruction")
	}
	rows = append(rows, map[string]any{"events": o.events})
	spellbookCapture(t, "client-audio-streams-device-queues", rows, "0c739b1f28c82f45e7e096fc6e8988677f72e933384353b9b40f4d5a3719d930")
}

func TestClientAudioStreamsBulkVoices(t *testing.T) {
	var rows []map[string]any
	for _, capacity := range []int{0, 1, 5} {
		for _, requested := range []uint32{0, 1, 2, 7, 0xffffffff} {
			for _, failAt := range []int{0, 2} {
				t.Run(fmt.Sprintf("%d/%d/%d", capacity, requested, failAt), func(t *testing.T) {
					o := newAudioStreamOwner(t, 1, capacity)
					if o.device() == 0 {
						t.Fatal("device")
					}
					ctx := o.call("sub_487150", 0, 0)
					if failAt != 0 {
						o.failOp = 4
						o.failAt = failAt
					}
					got := o.call("sub_487790", ctx, requested)
					want := capacity
					if requested != 0 && requested < uint32(want) {
						want = int(requested)
					}
					if failAt != 0 && want >= failAt {
						want = failAt - 1
					}
					if got != uint32(want) || audioStreamWords(ctx, 66)[48] != uint32(want) {
						t.Fatal("bulk capacity/failure", got, want)
					}
					o.call("sub_487910", ctx, 2)
					if audioStreamWords(ctx, 66)[48] != uint32(want) {
						t.Fatal("kind filtering")
					}
					o.call("sub_487910", ctx, 1)
					if audioStreamWords(ctx, 66)[48] != 0 {
						t.Fatal("bulk cleanup")
					}
					rows = append(rows, map[string]any{"capacity": capacity, "requested": requested, "fail_at": failAt, "allocated": got, "events": o.events})
				})
			}
		}
	}
	spellbookCapture(t, "client-audio-streams-bulk-voices", rows, "780b0b7ddd68bde5f2796d8732f52cb8e2f88053aba28c0af8791606bc30a72f")
}
