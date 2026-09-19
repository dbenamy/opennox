//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestClientAudioStreamsReads(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	var rows []map[string]any
	for _, mode := range []string{"bag", "missing", "bad-header", "mono8", "stereo16", "extra-chunk", "no-format", "zero-channels", "short-format", "truncated-format"} {
		t.Run(mode, func(t *testing.T) {
			cat, free := alloc.Make([]uint32{}, 72)
			defer free()
			entries, free := alloc.Make([]uint32{}, 9)
			defer free()
			dst, free := alloc.Make([]byte{}, 32)
			defer free()
			cp := audioStreamPointer(unsafe.Pointer(&cat[0]))
			ep := audioStreamPointer(unsafe.Pointer(&entries[0]))
			dp := audioStreamPointer(unsafe.Pointer(&dst[0]))
			cat[0] = ep
			cat[1] = 1
			copy(audioStreamMemory(ep, 16), "sample")
			entries[4], entries[5], entries[6], entries[7] = 3, 9, 11025, 0
			bag := []byte("---bag-bytes-tail")
			files, _ := prefabScriptsFiles(t, bag)
			cat[67] = audioStreamPointer(files[0])
			dir := t.TempDir()
			copy(audioStreamMemory(cp+8, 260), dir+string(os.PathSeparator))
			expected := []byte("bag-bytes")
			override := false
			if mode != "bag" {
				cat[69] = 1
			}
			if mode == "bad-header" {
				if err := os.WriteFile(filepath.Join(dir, "sample.wav"), []byte("invalid"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "mono8" || mode == "stereo16" || mode == "extra-chunk" {
				override = true
				channels, width := uint16(1), uint16(1)
				if mode == "stereo16" {
					channels, width = 2, 2
				}
				var wav bytes.Buffer
				wav.WriteString("RIFF")
				binary.Write(&wav, binary.LittleEndian, uint32(0))
				wav.WriteString("WAVE")
				wav.WriteString("fmt ")
				binary.Write(&wav, binary.LittleEndian, uint32(16))
				for _, v := range []uint16{1, channels} {
					binary.Write(&wav, binary.LittleEndian, v)
				}
				binary.Write(&wav, binary.LittleEndian, uint32(22050))
				binary.Write(&wav, binary.LittleEndian, uint32(22050)*uint32(channels*width))
				binary.Write(&wav, binary.LittleEndian, channels*width)
				binary.Write(&wav, binary.LittleEndian, width*8)
				if mode == "extra-chunk" {
					wav.WriteString("JUNK")
					binary.Write(&wav, binary.LittleEndian, uint32(4))
					wav.WriteString("abcd")
				}
				expected = []byte("wave-data-contents")
				wav.WriteString("data")
				binary.Write(&wav, binary.LittleEndian, uint32(len(expected)))
				wav.Write(expected)
				b := wav.Bytes()
				binary.LittleEndian.PutUint32(b[4:], uint32(len(b)-8))
				if err := os.WriteFile(filepath.Join(dir, "sample.wav"), b, 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "no-format" || mode == "zero-channels" || mode == "short-format" || mode == "truncated-format" {
				var wav bytes.Buffer
				wav.WriteString("RIFF")
				binary.Write(&wav, binary.LittleEndian, uint32(64))
				wav.WriteString("WAVE")
				if mode != "no-format" {
					wav.WriteString("fmt ")
					n := uint32(16)
					if mode == "short-format" {
						n = 4
					}
					binary.Write(&wav, binary.LittleEndian, n)
					if mode == "truncated-format" {
						wav.Write(make([]byte, 4))
					} else {
						wav.Write(make([]byte, n))
					}
				}
				if mode != "truncated-format" {
					wav.WriteString("data")
					binary.Write(&wav, binary.LittleEndian, uint32(4))
					wav.WriteString("abcd")
				}
				if err := os.WriteFile(filepath.Join(dir, "sample.wav"), wav.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
			}

			call := legacy.PortTestAudioStreamCall
			defer call("sub_486E00", cp)
			if call("sub_486B60", cp, 0) != 1 || cat[71] != uint32(len(expected)) {
				t.Fatal("open sample", cat[71])
			}
			if (cat[68] != 0) != override || (cat[70] != cat[67]) != override {
				t.Fatal("file ownership")
			}
			flags, rate := uint32(0), uint32(11025)
			if override {
				flags, rate = 2, 22050
				if mode == "stereo16" {
					flags = 7
				}
			}
			if entries[6] != rate || entries[7] != flags {
				t.Fatal("sample metadata")
			}
			var reads []uint32
			var data []byte
			for _, request := range []uint32{0, 0xffffffff, 1, 3, 2, 30, 1} {
				for i := range dst {
					dst[i] = 0xcc
				}
				n := call("sub_486DB0", cp, dp, request)
				want := 0
				if int32(request) > 0 {
					want = int(request)
					if want > len(expected)-len(data) {
						want = len(expected) - len(data)
					}
				}
				if n != uint32(want) || !bytes.Equal(dst[:n], expected[len(data):len(data)+want]) {
					t.Fatal("partial read", request, n, want)
				}
				for _, v := range dst[n:] {
					if v != 0xcc {
						t.Fatal("read exceeded returned length")
					}
				}
				data = append(data, dst[:n]...)
				reads = append(reads, n)
				if cat[71] != uint32(len(expected)-len(data)) {
					t.Fatal("remaining bytes")
				}
			}
			if call("sub_486E00", cp) != 0 || cat[68] != 0 || cat[70] != 0 || cat[67] != audioStreamPointer(files[0]) {
				t.Fatal("close ownership")
			}
			if call("sub_486DB0", cp, dp, 1) != 0 {
				t.Fatal("closed read")
			}
			rows = append(rows, map[string]any{"mode": mode, "reads": reads, "data": fmt.Sprintf("%x", data), "flags": flags, "rate": rate})
		})
	}
	spellbookCapture(t, "client-audio-streams-reads", rows, "a57303642d989dc3b084434f6bd921a2959bf47ee6f6c9176a04303c48523577")
}
