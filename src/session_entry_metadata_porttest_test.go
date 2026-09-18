//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestSessionEntrySaveMetadata(t *testing.T) {
	defer noxflags.PortTestGameFlags(0)()
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(original) }()
	current := serverConfigOwnBytes(t, 0x85B3FC, 10980, 1278)
	base := serverConfigOwnBytes(t, 0x85B3FC, 36, 80)
	clear(base)
	copy(base, "fixture")
	marker := serverConfigOwnBytes(t, 0x5D4594, 527728, 1)
	marker[0] = 0x6d
	restore := legacy.PortTestSessionEntryMetadataTable()
	defer restore()
	type row struct {
		Version  int
		Unknown  bool
		Result   int
		Output   []byte
		Restored bool
	}
	var rows []row
	for _, version := range []int{10, 11, 12, 13} {
		for _, unknown := range []bool{false, true} {
			path := filepath.Join(t.TempDir(), "player.plr")
			f, err := cryptfile.OpenFile(path, cryptfile.WriteOnly, 27)
			if err != nil {
				t.Fatal(err)
			}
			write := func(p []byte) {
				t.Helper()
				if _, err := f.Write(p); err != nil {
					t.Fatal(err)
				}
			}
			u16 := func(v uint16) { var b [2]byte; binary.LittleEndian.PutUint16(b[:], v); write(b[:]) }
			u32 := func(v uint32) { var b [4]byte; binary.LittleEndian.PutUint32(b[:], v); write(b[:]) }
			if unknown {
				u32(77)
				f.SectionStart()
				write([]byte{9, 8, 7, 6, 5, 4, 3, 2, 1})
				f.SectionEnd()
			}
			u32(1)
			f.SectionStart()
			u16(uint16(version))
			u32(0x11223302)
			u16(3)
			write([]byte("old"))
			write([]byte{3, 'a', 'b', 'c'})
			stamp := []byte{0xe8, 7, 9, 0, 3, 0, 18, 0, 12, 0, 34, 0, 56, 0, 78, 0}
			write(stamp)
			colors := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
			write(colors)
			write([]byte{2, 'A', 0, 'Z', 0, 2, 0x45, 0x67})
			if version >= 11 {
				write([]byte{5, 'a', 'r', 'e', 'n', 'a'})
			}
			if version >= 12 {
				write([]byte{9})
			}
			f.SectionEnd()
			u32(0)
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			for i := range current {
				current[i] = 0xa5
			}
			before := bytes.Clone(current)
			out, free := alloc.Make([]byte{}, 1280)
			for i := range out {
				out[i] = 0xcc
			}
			got := legacy.Sub_41A000(path, (*server.SaveGameInfo)(unsafe.Pointer(&out[0])))
			wantRet := 1
			if version > 12 {
				wantRet = 0
			}
			if got != wantRet {
				t.Fatalf("version%d unknown%v result%d", version, unknown, got)
			}
			if !bytes.Equal(current, before) {
				t.Errorf("version%d metadata global was not restored", version)
			}
			if cryptfile.Global() != nil {
				t.Fatal("metadata file left open")
			}
			want := bytes.Repeat([]byte{0xcc}, 1280)
			want[1028] = 0x6d
			if version <= 12 {
				copy(want, before)
				binary.LittleEndian.PutUint32(want, 0x11223302)
				copy(want[4:], path)
				want[4+len(path)] = 0
				copy(want[1028:], "abc\x00")
				copy(want[1188:], stamp)
				copy(want[1207:], colors[:3])
				copy(want[1204:], colors[3:6])
				copy(want[1210:], colors[6:])
				copy(want[1224:], []byte{'A', 0, 'Z', 0, 0, 0})
				want[1274] = 2
				want[1275] = 0x45
				want[1276] = 0x67
				want[1277] = 0
				if version >= 11 {
					copy(want[1156:], "fixture\x00")
					copy(want[1156:], "arena\x00")
				}
				if version >= 12 {
					want[1277] = 9
				}
			}
			if !bytes.Equal(out, want) {
				for i := range out {
					if out[i] != want[i] {
						t.Errorf("version%d unknown%v byte%d=%x want%x", version, unknown, i, out[i], want[i])
						break
					}
				}
			}
			// The temporary path is the only process-specific payload.
			normalized := bytes.Clone(out)
			if got != 0 {
				clear(normalized[4:1028])
				copy(normalized[4:], "<fixture>/player.plr")
			}
			rows = append(rows, row{version, unknown, got, normalized, bytes.Equal(current, before)})
			free()
		}
	}
	t.Run("missing", func(t *testing.T) {
		out, free := alloc.Make([]byte{}, 1280)
		defer free()
		for i := range out {
			out[i] = 0xcc
		}
		got := legacy.Sub_41A000(filepath.Join(t.TempDir(), "missing.plr"), (*server.SaveGameInfo)(unsafe.Pointer(&out[0])))
		want := bytes.Repeat([]byte{0xcc}, 1280)
		want[1028] = 0x6d
		if got != 0 || !bytes.Equal(out, want) {
			t.Fatal(fmt.Sprintf("missing metadata result%d", got))
		}
	})
	spellbookCapture(t, "session-entry-save-metadata", rows, "e01db7340caea46a3d87b232e80d79f4faa5477385d97f0bd9d92ca0a9ddb750")
}
