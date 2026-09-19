//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"io"
	"path/filepath"
	"testing"
	"unsafe"

	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/music"
	"github.com/opennox/opennox/v1/server"
)

func TestPlayerFilesClientWriteFraming(t *testing.T) {
	q := newQuickbarOwner(t)
	defer flags.PortTestGameFlags(0)()
	t.Cleanup(legacy.PortTestPlayerFileClientSections())
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	raw := serverConfigOwnBytes(t, 0x85B3FC, 10980, 1284)
	currentMap := serverConfigOwnBytes(t, 0x85B3FC, 36, 32)
	clear(currentMap)
	copy(currentMap, "forest.map")
	meters := legacy.PortTestNewMeterEnvironment()
	t.Cleanup(meters.Restore)
	*meters.NamedWord("dword_8531A0_2576") = 0
	words, restore := legacy.PortTestAudioEventGlobals()
	t.Cleanup(restore)
	*words["dword_5d4594_816368"] = 0
	*words["dword_5d4594_816372"] = 0
	oldMusic := legacy.MusicModule
	t.Cleanup(func() { legacy.MusicModule = oldMusic })
	state := music.MusicState{D: 17, Position: 987, MusicIdx: 23, Volume: 75}
	legacy.MusicModule = &music.Module{}
	legacy.MusicModule.SetNextMusic(state)
	var rows []map[string]any
	for _, mode := range []uint32{0, 1, 2, 0xffffffff} {
		for _, missing := range []bool{false, true} {
			q.reset(t)
			*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
			*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 3648)) = 7
			path := filepath.Join(t.TempDir(), "save.plr")
			if missing {
				path = filepath.Join(path, "missing.plr")
			}
			info, release := alloc.New(server.SaveGameInfo{})
			copy(info.PathBuf[:], path)
			copy(info.Field1028[:], "Port framing")
			info.Flags = 0xa5a5a5a7
			info.Stage = 19
			before := bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(info)), 1280))
			clear(raw)
			raw[1280] = 0xd1
			raw[1281] = 0xd2
			raw[1282] = 0xd3
			raw[1283] = 0xd4
			ret := legacy.PortTestPlayerFileCall("sub_41CEE0", uint32(uintptr(unsafe.Pointer(info))), mode)
			if !bytes.Equal(before, unsafe.Slice((*byte)(unsafe.Pointer(info)), 1280)) {
				t.Fatal("input mutation")
			}
			release()
			wantRet := uint32(1)
			if missing {
				wantRet = 0
			}
			if ret != wantRet || cryptfile.Global() != nil || !bytes.Equal(raw[1280:], []byte{0xd1, 0xd2, 0xd3, 0xd4}) {
				t.Fatal("writer result", mode, missing, ret)
			}
			if missing {
				if !bytes.Equal(raw[:1280], before) {
					t.Fatal("failed open must retain copied info")
				}
				rows = append(rows, map[string]any{"mode": mode, "missing": true, "return": ret})
				continue
			}
			f, err := cryptfile.OpenFile(path, cryptfile.ReadOnly, 27)
			if err != nil {
				t.Fatal(err)
			}
			ids := []uint32{1}
			if mode != 0 {
				ids = []uint32{7, 1, 12}
			}
			var payloads [][]byte
			for _, wantID := range ids {
				id, err := f.ReadU32()
				if err != nil || id != wantID {
					t.Fatal("section ID", id, wantID, err)
				}
				n, err := f.ReadAlignedU32()
				if err != nil || n > 4096 {
					t.Fatal("section length", n, err)
				}
				data := make([]byte, n)
				if _, err = io.ReadFull(f, data); err != nil {
					t.Fatal(err)
				}
				switch id {
				case 1:
					// Metadata's detailed independent contract covers its payload. Here check
					// the enclosing writer copied the input, selected it, and framed it whole.
					pathN := int(binary.LittleEndian.Uint16(data[6:]))
					if binary.LittleEndian.Uint16(data) != 12 || string(data[8:8+pathN]) != path {
						t.Fatal("metadata path")
					}
					titleAt := 8 + pathN
					titleN := int(data[titleAt])
					clockAt := titleAt + 1 + titleN
					if string(data[titleAt+1:clockAt]) != "Port framing" || data[len(data)-1] != 19 {
						t.Fatal("metadata fields")
					}
					if !bytes.Equal(data[clockAt:clockAt+16], raw[1188:1204]) {
						t.Fatal("metadata clock copy")
					}
					clear(data[clockAt : clockAt+16])
					// The temporary pathname varies between processes; assert it above and
					// remove only that identified field from the repeated capture.
					normalized := bytes.Clone(data[:6])
					normalized = append(normalized, 0, 0)
					data = append(normalized, data[8+pathN:]...)
				case 7:
					// Reset bars contain canonical invalid-spell names and zero flags.
					want := []byte{3, 0, 1}
					for i := 0; i < 34; i++ {
						want = append(want, 13)
						want = append(want, "SPELL_INVALID"...)
						want = append(want, 0)
					}
					want = append(want, 0, 0, 7)
					if !bytes.Equal(data, want) {
						t.Fatalf("GUI framing: %x != %x", data, want)
					}
				case 12:
					want := playerFileMusicBytes([]byte{11, 0, 1}, state)
					want = binary.LittleEndian.AppendUint32(want, 0)
					if !bytes.Equal(data, want) {
						t.Fatal("music framing", data, want)
					}
				}
				payloads = append(payloads, data)
			}
			end, err := f.ReadU32()
			if err != nil || end != 0 {
				t.Fatal("section terminator", end, err)
			}
			if err = f.Close(); err != nil {
				t.Fatal(err)
			}
			rows = append(rows, map[string]any{"mode": mode, "missing": false, "return": ret, "ids": ids, "payloads": payloads})
		}
	}
	spellbookCapture(t, "player-files-client-write-framing", rows, "")
}
