//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestThingSkipsConcatenatedSections(t *testing.T) {
	var rows []thingSkipResult
	for offset := 0; offset < 8; offset++ {
		for _, reverse := range []bool{false, true} {
			type section struct {
				op         int
				b, scratch []byte
			}
			var sections []section
			input := bytes.Repeat([]byte{0x83}, offset)
			for j := 0; j < 7; j++ {
				op := j
				if reverse {
					op = 6 - j
				}
				var b, scratch []byte
				switch op {
				case 0:
					b = thingSkipWord(nil, 1)
					b = thingSkipName(b, 31, 0x41)
					b = thingSkipFill(b, 9, 0x81)
					b = thingSkipName(b, 255, 0x61)
					b = append(b, 0)
				case 1, 2:
					b = thingSkipWord(nil, 2)
					for i := 0; i < 2; i++ {
						b = thingSkipSpellRecord(b, op == 2, 3+17*i, 255, 2)
					}
				case 3:
					b = thingSkipWord(nil, 1)
					b = thingSkipName(b, 3, 0x41)
					b = append(b, 2, 3, 0x81)
					b = thingSkipName(b, 7, 0x51)
					for i := 0; i < 3; i++ {
						b = thingSkipRef(b, i%2 == 0, 31)
					}
				case 4, 5:
					if op == 4 {
						b = thingSkipName(b, 7, 0x51)
					}
					for tag := byte(1); tag <= 10; tag++ {
						b = thingSkipEventPayload(b, tag, 31)
					}
					b = append(b, 0)
				case 6:
					b, scratch, _, _ = thingSkipWallRecord(len(input), 8, 31, 2, 2, 0x454e4420)
				}
				sections = append(sections, section{op, b, scratch})
				input = append(input, b...)
			}
			input = append(input, bytes.Repeat([]byte{0xf3}, 16)...)
			func() {
				raw, _ := alloc.CloneSlice(input)
				f := binfile.NewMemFile(unsafe.Pointer(&raw[0]), len(raw))
				defer f.Free()
				f.Skip(offset)
				scratch, free := alloc.Make([]byte{}, 256*1024)
				defer free()
				for i := range scratch {
					scratch[i] = 0xa5
				}
				want := append([]byte(nil), scratch...)
				cursor := offset
				for _, s := range sections {
					ret := legacy.PortTestThingSkip(s.op, f, scratch)
					next := len(raw) - len(f.Data())
					copy(want, s.scratch)
					if ret != 1 || next != cursor+len(s.b) || !bytes.Equal(raw, input) || !bytes.Equal(scratch, want) {
						t.Fatalf("section%d offset%d reverse%v cursor%d want%d", s.op, offset, reverse, next, cursor+len(s.b))
					}
					rows = append(rows, thingSkipResult{s.op, cursor, ret, next - cursor, sha256.Sum256(scratch)})
					cursor = next
				}
			}()
		}
	}
	thingSkipCapture(t, "concatenated", rows, "198057bc2e0c28bab1385a6a52f2a437c73c7db526b09bf112e0a5d6d44bd127")
}
