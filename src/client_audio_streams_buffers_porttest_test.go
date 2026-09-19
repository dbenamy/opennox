//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsBuffers(t *testing.T) {
	var rows []map[string]any
	for _, lengths := range [][]uint32{nil, {0}, {1}, {3, 7, 0, 12}, {0xffffffff, 2}} {
		t.Run(fmt.Sprint(lengths), func(t *testing.T) {
			o := newAudioStreamOwner(t, 1, 1)
			if o.device() == 0 {
				t.Fatal("device")
			}
			ctx := o.call("sub_487150", 0, 0)
			voice := o.call("sub_487750", ctx)
			vw := audioStreamWords(voice, 78)
			buffer, free := alloc.Make([]uint32{}, 7)
			defer free()
			bp := audioStreamPointer(unsafe.Pointer(&buffer[0]))
			chunks, free := alloc.Make([]uint32{}, 6*(len(lengths)+1))
			defer free()
			for i := range buffer {
				buffer[i] = 0xabcdef01
			}
			o.call("sub_487C30", bp)
			if buffer[0] != 0 || buffer[1] != 0 || buffer[5] != 0 || buffer[6] != 0 || o.call("sub_487C80", bp) != 0 {
				t.Fatal("buffer init")
			}
			var total uint32
			for i, n := range lengths {
				chunk := audioStreamPointer(unsafe.Pointer(&chunks[6*i]))
				data := uint32(0x12340000 + i*256)
				o.call("sub_487D30", chunk, data, n)
				if chunks[6*i+3] != data || chunks[6*i+4] != n || chunks[6*i+5] != 0 {
					t.Fatal("chunk init")
				}
				total += n
				if o.call("sub_487C50", bp, chunk) != total || buffer[1] != total || chunks[6*i+5] != bp {
					t.Fatal("chunk attachment")
				}
			}
			if len(lengths) == 0 {
				buffer[0] = 0x76543210
				buffer[1] = 19
			}
			o.call("sub_4BDB90", voice, bp)
			var observations [][]uint32
			if vw[72] != bp {
				t.Fatal("buffer binding")
			}
			if len(lengths) == 0 {
				if vw[73] != 0 || vw[74] != buffer[0] || vw[75] != 19 || vw[76] != 19 {
					t.Fatal("direct buffer")
				}
				observations = append(observations, append([]uint32(nil), vw[74:77]...))
			} else {
				for i, n := range lengths {
					if vw[73] != audioStreamPointer(unsafe.Pointer(&chunks[6*i])) || vw[74] != uint32(0x12340000+i*256) || vw[75] != n || vw[76] != n {
						t.Fatal("chunk traversal", i)
					}
					observations = append(observations, append([]uint32(nil), vw[74:77]...))
					if o.call("sub_4BD8C0", voice) != 0 {
						t.Fatal("chunk advance")
					}
				}
			}
			o.call("sub_4BD8C0", voice)
			if vw[75] != 0 {
				t.Fatal("buffer end")
			}
			before := append([]uint32(nil), vw[73:77]...)
			o.call("sub_4BDB90", voice, 0)
			if vw[72] != 0 || !reflect.DeepEqual(before, vw[73:77]) {
				t.Fatal("unbind preserves cursor metadata")
			}
			for i := range lengths {
				chunk := audioStreamPointer(unsafe.Pointer(&chunks[6*i]))
				before := append([]uint32(nil), chunks[6*i:6*i+6]...)
				before[5] = 0
				if o.call("sub_487D60", chunk) != chunk || !reflect.DeepEqual(before, chunks[6*i:6*i+6]) {
					t.Fatal("detach owner only")
				}
			}
			rows = append(rows, map[string]any{"lengths": lengths, "total": total, "segments": observations, "end": before})
		})
	}
	spellbookCapture(t, "client-audio-streams-buffers", rows, "c949de7087a94661e02ffaaab2a990098ea039b04dfe6bfabce103a7ef873107")
}
