//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsFormats(t *testing.T) {
	input, free := alloc.Make([]uint32{}, 12)
	defer free()
	output, freeOut := alloc.Make([]uint32{}, 8)
	defer freeOut()
	catalog, freeCat := alloc.Make([]uint32{}, 1)
	defer freeCat()
	entries, freeEnt := alloc.Make([]uint32{}, 9*3)
	defer freeEnt()
	in, out, cat := audioStreamPointer(unsafe.Pointer(&input[0])), audioStreamPointer(unsafe.Pointer(&output[0])), audioStreamPointer(unsafe.Pointer(&catalog[0]))
	catalog[0] = audioStreamPointer(unsafe.Pointer(&entries[0]))
	type record struct {
		Kind   string
		Args   []uint32
		Return uint32
		Words  []uint32
	}
	var rows []record
	for _, encoding := range []uint32{0, 1, 2, 0xffffffff} {
		for _, rate := range []uint32{0, 1, 22050, 44100, 0x80000000, 0xffffffff} {
			for _, channels := range []uint32{0, 1, 2, 0xffff} {
				for _, width := range []uint32{0, 1, 2, 4, 0xffffffff} {
					for i := range input {
						input[i] = uint32(i) + 0x12340000
					}
					input[1], input[2], input[3], input[4] = encoding, rate, channels, width
					want := append([]uint32(nil), input...)
					product := int32(rate * channels * width)
					if encoding == 1 {
						product >>= 2
					}
					want[5] = uint32(product)
					got := legacy.PortTestAudioStreamCall("sub_487D00", in)
					if got != uint32(product) || !reflect.DeepEqual(input, want) {
						t.Fatal("sample byte rate/32-bit signed quarter")
					}
					rows = append(rows, record{"byte-rate", []uint32{encoding, rate, channels, width}, got, append([]uint32(nil), input...)})
				}
			}
		}
	}
	for index := 0; index < 3; index++ {
		for flags := uint32(0); flags < 16; flags++ {
			for _, rate := range []uint32{0, 22050, 0xffffffff} {
				for i := range entries {
					entries[i] = 0xdead0000 + uint32(i)
				}
				entries[9*index+6], entries[9*index+7], entries[9*index+8] = rate, flags, 0x76543210
				for i := range output {
					output[i] = 0xaabb0000 + uint32(i)
				}
				want := append([]uint32(nil), output...)
				want[0], want[1], want[2], want[3], want[4], want[6] = 4, 0, rate, 1+(flags&1), 1+((flags>>2)&1), 0x76543210
				if flags&8 != 0 {
					want[1], want[4] = 2, 2
				}
				p := legacy.PortTestAudioStreamCall("sub_4866D0", cat, uint32(index))
				if p != catalog[0]+uint32(36*index) {
					t.Fatal("catalog entry stride")
				}
				got := legacy.PortTestAudioStreamCall("sub_486AA0", cat, uint32(index), out)
				if got != want[4] || !reflect.DeepEqual(output, want) {
					t.Fatal("catalog format mapping/untouched words")
				}
				rows = append(rows, record{"format", []uint32{uint32(index), flags, rate}, got, append([]uint32(nil), output...)})
			}
		}
	}
	for _, raw := range []uint32{0, 1, 0xffff, 0x10000, 0x80000000, 0xffffffff} {
		for _, percent := range []uint32{0, 1, 50, 100, 0x7fffffff, 0xffffffff} {
			input[9] = raw
			want := (percent * (raw >> 16)) / 100
			got := legacy.PortTestAudioStreamCall("sub_486640", in, percent)
			if got != want {
				t.Fatal("volume scale uses unsigned 32-bit product")
			}
			rows = append(rows, record{"volume", []uint32{raw, percent}, got, nil})
		}
	}
	spellbookCapture(t, "client-audio-streams-formats", rows, "55c5d1ef0b2f877c59e932257ffea3adb650d36db69ab395a33d989c5435731e")
}
