//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabRuntimeLibrary(t *testing.T) {
	handles.Init()
	defer handles.Release()
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) }()
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	pathBuf, freePath := alloc.Make([]byte{}, 2048)
	defer freePath()
	alternate, freeAlt := alloc.Make([]byte{}, 2048)
	defer freeAlt()
	metadata, freeMeta := alloc.Make([]byte{}, 2048*76)
	defer freeMeta()
	*words["path"] = uint32(uintptr(unsafe.Pointer(&pathBuf[0])))
	*words["alternate"] = uint32(uintptr(unsafe.Pointer(&alternate[0])))
	*words["metadata"] = uint32(uintptr(unsafe.Pointer(&metadata[0])))
	path := filepath.Join(t.TempDir(), "library.lib")
	copy(pathBuf, path)
	type row struct {
		Name    string
		Return  uint64
		Count   uint32
		Records []byte
		Closed  bool
	}
	var rows []row
	for _, count := range []int{0, 1, 3, 2048, 2049} {
		for _, sentinel := range []bool{false, true} {
			var b bytes.Buffer
			binary.Write(&b, binary.LittleEndian, uint32(0xcafedead))
			var want []byte
			for i := 0; i < count; i++ {
				name := fmt.Sprintf("Room-%04d", i)
				if i%3 == 2 {
					name = strings.Repeat("x", 63)
				}
				offset := uint32(b.Len())
				size := 1 + len(name) + 10 + (i % 7)
				binary.Write(&b, binary.LittleEndian, uint32(size))
				b.WriteByte(byte(len(name)))
				b.WriteString(name)
				b.WriteByte(byte(i))
				b.WriteByte(1)
				width := math.Float32bits(float32(i) + 0.25)
				height := math.Float32bits(-float32(i) - 0.75)
				binary.Write(&b, binary.LittleEndian, width)
				binary.Write(&b, binary.LittleEndian, height)
				b.Write(bytes.Repeat([]byte{0x5a}, i%7))
				if i < 2048 {
					rec := make([]byte, 76)
					copy(rec, name)
					binary.LittleEndian.PutUint32(rec[64:], width)
					binary.LittleEndian.PutUint32(rec[68:], height)
					binary.LittleEndian.PutUint32(rec[72:], offset)
					want = append(want, rec...)
				}
			}
			if sentinel {
				binary.Write(&b, binary.LittleEndian, uint32(0))
			}
			clear(metadata)
			*words["count"] = 0xa5a5a5a5
			if err := os.WriteFile(path, b.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			ret := legacy.PortTestPrefabCall(7, [6]uint32{})
			wantReturn := uint64(1)
			if count > 2048 {
				wantReturn = 0
			}
			if ret != wantReturn || *words["count"] != uint32(min(count, 2048)) || !bytes.Equal(metadata[:len(want)], want) {
				t.Fatal("library metadata/capacity differs")
			}
			closed := *words["file"] == 0 && cryptfile.Global() == nil
			if !closed {
				t.Fatal("metadata file left open")
			}
			rows = append(rows, row{fmt.Sprintf("count%d/sentinel%t", count, sentinel), ret, *words["count"], append([]byte(nil), metadata[:len(want)]...), closed})
		}
	}
	for _, tc := range []struct {
		name   string
		data   []byte
		remove bool
	}{
		{"bad-magic", []byte{0, 0, 0, 0}, false}, {"missing", nil, true},
	} {
		if tc.remove {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		} else if err := os.WriteFile(path, tc.data, 0600); err != nil {
			t.Fatal(err)
		}
		ret := legacy.PortTestPrefabCall(7, [6]uint32{})
		if ret != 0 || *words["count"] != 0 || *words["file"] != 0 || cryptfile.Global() != nil {
			t.Fatal("failed library state")
		}
		rows = append(rows, row{tc.name, ret, *words["count"], nil, true})
	}
	spellbookCapture(t, "prefab-runtime-library", rows, "247b891633bef2c5cbb718145b39fc270a8649f9828b4322fa13a9c9ca2293d9")
}
