//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"sort"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientPaletteTables(t *testing.T) {
	palette := serverConfigOwnBytes(t, 0x973F18, 3880, 1024)
	sorted := serverConfigOwnBytes(t, 0x5D4594, 809604, 1024)
	indices := serverConfigOwnBytes(t, 0x5D4594, 808304, 256)
	type row struct {
		Mode                          int
		Palette, Sorted, Indices, RGB []byte
	}
	var rows []row
	for mode := 0; mode < 20; mode++ {
		rgb := make([]byte, 768)
		for i := range rgb {
			switch mode {
			case 0:
				rgb[i] = 0
			case 1:
				rgb[i] = 255
			case 2:
				rgb[i] = byte(i / 3)
			case 3:
				rgb[i] = byte(255 - i/3)
			default:
				rgb[i] = byte(i*i*17 + i*(mode+1) + mode*29)
			}
		}
		input := bytes.Clone(rgb)
		legacy.Sub_435120(unsafe.Pointer(&palette[0]), unsafe.Pointer(&rgb[0]))
		if !bytes.Equal(rgb, input) {
			t.Fatal("source changed")
		}
		for i := 0; i < 256; i++ {
			if !bytes.Equal(palette[4*i:4*i+3], rgb[3*i:3*i+3]) || palette[4*i+3] != 4 {
				t.Fatal("RGB expansion", mode, i)
			}
		}
		before := bytes.Clone(palette)
		legacy.Sub_435040()
		if !bytes.Equal(palette, before) {
			t.Fatal("sort changed source palette")
		}
		want := make([]int, 256)
		for i := range want {
			want[i] = i
		}
		// Packed C ordering is blue, green, red, then original index.
		sort.Slice(want, func(a, b int) bool {
			i, j := want[a], want[b]
			for k := 2; k >= 0; k-- {
				if rgb[3*i+k] != rgb[3*j+k] {
					return rgb[3*i+k] < rgb[3*j+k]
				}
			}
			return i < j
		})
		for i, index := range want {
			if indices[i] != byte(index) || !bytes.Equal(sorted[4*i:4*i+3], rgb[3*index:3*index+3]) || sorted[4*i+3] != 0 {
				t.Fatal("sorted RGB/index", mode, i, index)
			}
		}
		output := make([]byte, 768)
		legacy.Sub_435150(unsafe.Pointer(&output[0]), unsafe.Pointer(&palette[0]))
		if !bytes.Equal(output, rgb) {
			t.Fatal("RGB round trip", mode)
		}
		rows = append(rows, row{mode, bytes.Clone(palette), bytes.Clone(sorted), bytes.Clone(indices), output})
	}
	// Fourth-byte tags are ignored by packing and by palette sorting.
	for i := 0; i < 256; i++ {
		binary.LittleEndian.PutUint32(palette[4*i:], uint32(i)*0x01020304)
	}
	output := make([]byte, 768)
	legacy.Sub_435150(unsafe.Pointer(&output[0]), unsafe.Pointer(&palette[0]))
	for i := 0; i < 256; i++ {
		if !bytes.Equal(output[3*i:3*i+3], palette[4*i:4*i+3]) {
			t.Fatal("tag omitted", i)
		}
	}
	interactionCapture(t, "client-palette-tables", rows)
}
