//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/ifs"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestSessionMOTDFileOwnership(t *testing.T) {
	restoreHandles := handles.PortTestInit()
	defer restoreHandles()
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	defer sessionCall("motdFree", nil, 0, 0, 0)
	sizes := serverConfigOwnBytes(t, 0x5D4594, 826040, 8)
	old, err := ifs.Workdir()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err = ifs.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ifs.Chdir(old); err != nil {
			t.Error(err)
		}
	}()
	contents := [][]byte{nil, {}, []byte("x"), []byte("first\r\n\rsecond\nlast"), {0, 1, 0x80, 0xff, 0}}
	for _, n := range []int{254, 255, 256, 4096} {
		contents = append(contents, bytes.Repeat([]byte{'x'}, n))
	}
	type row struct {
		Slot        uint32
		Case        int
		Length      int
		Missing     bool
		Size        uint32
		Free, Again uintptr
	}
	var rows []row
	for slot := uint32(0); slot < 2; slot++ {
		for i, data := range contents {
			path := filepath.Join(dir, "motd.txt")
			if data == nil {
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			for j := 0; j < 2; j++ {
				binary.LittleEndian.PutUint32(sizes[4*j:], 0x11223344)
			}
			if sessionCall("motdRead", nil, 0, slot, 0) != 0 {
				t.Fatal("MOTD file return")
			}
			count := binary.LittleEndian.Uint32(sizes[4*slot:])
			want := uint32(0)
			if data != nil {
				want = uint32(len(data) + 1)
			}
			if count != want || binary.LittleEndian.Uint32(sizes[4*(1-slot):]) != 0x11223344 {
				t.Fatal("MOTD size/neighbor", slot, i, count, want)
			}
			if data == nil {
				if *words["motdFile"] != 0 {
					t.Fatal("missing file allocated")
				}
			} else {
				if *words["motdFile"] == 0 {
					t.Fatal("missing MOTD buffer")
				}
				raw := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(*words["motdFile"]))), int(count))
				if !bytes.Equal(raw[:len(data)], data) || raw[len(data)] != 0 {
					t.Fatal("MOTD contents/terminator", slot, i)
				}
			}
			first := sessionCall("motdFree", nil, 0, slot, 0)
			second := sessionCall("motdFree", nil, 0, slot, 0)
			wantFree := uintptr(0)
			if data != nil {
				wantFree = uintptr(slot)
			}
			if first != wantFree || second != 0 || *words["motdFile"] != 0 || binary.LittleEndian.Uint32(sizes[4*slot:]) != 0 {
				t.Fatal("MOTD release", slot, i, first, second)
			}
			rows = append(rows, row{Slot: slot, Case: i, Length: len(data), Missing: data == nil, Size: count, Free: first, Again: second})
		}
	}
	sessionCapture(t, "motd-files", rows)
}
