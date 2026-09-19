//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestPlayerFilesSectionExtraction(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	output, free := alloc.Make([]byte{}, 4096)
	t.Cleanup(free)
	var rows []map[string]any
	for _, size := range []int{0, 1, 7, 8, 9, 31} {
		for _, mixed := range []bool{false, true} {
			ids := []uint32{2, 3, 4, 8, 5, 6, 9, 10, 11}
			if mixed {
				ids = []uint32{99, 2, 1, 3, 4, 7, 8, 5, 6, 0xffff, 9, 10, 11, 12}
			}
			path := filepath.Join(t.TempDir(), "sections.plr")
			f, err := cryptfile.OpenFile(path, cryptfile.WriteOnly, 27)
			if err != nil {
				t.Fatal(err)
			}
			expected := bytes.Repeat([]byte{0xa5}, len(output))
			outPos, inputPos := 0, 0
			for i, id := range ids {
				payload := bytes.Repeat([]byte{byte(i*11 + 3)}, size)
				if err = f.WriteU32(id); err != nil {
					t.Fatal(err)
				}
				f.SectionStart()
				if _, err = f.Write(payload); err != nil {
					t.Fatal(err)
				}
				f.SectionEnd()
				afterID := inputPos + 4
				bodyAt := (afterID + 7) &^ 7
				bodyAt += 8 // Block-aligned length occupies a full block.
				known := id == 2 || id == 3 || id == 4 || id == 5 || id == 6 || id == 8 || id == 9 || id == 10 || id == 11
				if known {
					headerSize := bodyAt - inputPos
					binary.LittleEndian.PutUint32(expected[outPos:], id)
					// The extractor retains spacing for the aligned length block, leaving
					// non-field bytes in the destination untouched.
					binary.LittleEndian.PutUint32(expected[outPos+headerSize-8:], uint32(size))
					outPos += headerSize
					copy(expected[outPos:], payload)
					outPos += size
				}
				inputPos = bodyAt + size
			}
			if err = f.WriteU32(0); err != nil {
				t.Fatal(err)
			}
			if err = f.Close(); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint32(expected[outPos:], 0)
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			copy(output, bytes.Repeat([]byte{0xa5}, len(output)))
			name, release := alloc.CString(path)
			global := cryptfile.Global()
			legacy.PortTestPlayerFileCall("sub_41CAC0", uint32(uintptr(unsafe.Pointer(name))), uint32(uintptr(unsafe.Pointer(&output[0]))))
			release()
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(output, expected) || !bytes.Equal(before, after) || cryptfile.Global() != global {
				t.Fatal("section extraction", size, mixed)
			}
			rows = append(rows, map[string]any{"size": size, "mixed": mixed, "bytes": bytes.Clone(output[:outPos+4])})
		}
	}
	path, release := alloc.CString(filepath.Join(t.TempDir(), "missing.plr"))
	defer release()
	before := bytes.Clone(output)
	legacy.PortTestPlayerFileCall("sub_41CAC0", uint32(uintptr(unsafe.Pointer(path))), uint32(uintptr(unsafe.Pointer(&output[0]))))
	if !bytes.Equal(output, before) {
		t.Fatal("missing extraction changed output")
	}
	rows = append(rows, map[string]any{"missing": true})
	spellbookCapture(t, "player-files-section-extraction", rows, "51d2556dee434201136375c841b2a6a1b39bdfdd05d12ea671e61c34566355cb")
}
