//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/music"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestPlayerFilesClientLoadFraming(t *testing.T) {
	q := newQuickbarOwner(t)
	defer flags.PortTestGameFlags(0)()
	t.Cleanup(legacy.PortTestPlayerFileClientSections())
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	pathField := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 10984), 1024)
	oldPath := bytes.Clone(pathField)
	t.Cleanup(func() { copy(pathField, oldPath) })
	oldMusic := legacy.MusicModule
	legacy.MusicModule = &music.Module{}
	t.Cleanup(func() { legacy.MusicModule = oldMusic })
	_, restore := legacy.PortTestAudioEventGlobals()
	t.Cleanup(restore)
	queue := unsafe.Slice(memmap.PtrUint8(0x5D4594, 815772), 320)
	oldQueue := bytes.Clone(queue)
	t.Cleanup(func() { copy(queue, oldQueue) })
	var rows []map[string]any
	for _, which := range []string{"empty", "unknown", "gui", "music", "both", "gui-fails", "music-fails", "missing"} {
		q.reset(t)
		clear(queue)
		clear(pathField)
		copy(pathField, "prior.plr")
		legacy.MusicModule = &music.Module{}
		for i := 0; i < 25; i++ {
			q.bar[2*i] = 99
			q.bar[2*i+1] = 0xaabbccff
		}
		path := filepath.Join(t.TempDir(), "load.plr")
		if which != "missing" {
			f, err := cryptfile.OpenFile(path, cryptfile.WriteOnly, 27)
			if err != nil {
				t.Fatal(err)
			}
			section := func(id uint32, payload []byte) {
				t.Helper()
				if err := f.WriteU32(id); err != nil {
					t.Fatal(err)
				}
				f.SectionStart()
				if _, err := f.Write(payload); err != nil {
					t.Fatal(err)
				}
				f.SectionEnd()
			}
			if which != "empty" {
				section(99, []byte{0xe7, 0x33, 0x82, 1, 2, 3, 4, 5, 6})
			}
			if which == "gui" || which == "both" {
				data := append([]byte{3, 0}, playerFileQuickbarBytes(1)...)
				data = append(data, 3, 2, 7)
				section(7, data)
			}
			if which == "gui-fails" {
				section(7, []byte{4, 0})
			}
			if which == "music" || which == "both" {
				data := []byte{11, 0, 1}
				data = playerFileMusicBytes(data, music.MusicState{D: 17, Position: 987, MusicIdx: 23, Volume: 75})
				data = binary.LittleEndian.AppendUint32(data, 0)
				section(12, data)
			}
			if which == "music-fails" {
				section(12, []byte{12, 0})
			}
			if err = f.WriteU32(0); err != nil {
				t.Fatal(err)
			}
			if err = f.Close(); err != nil {
				t.Fatal(err)
			}
		}
		var before []byte
		if which != "missing" {
			var err error
			before, err = os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
		}
		name, release := alloc.CString(path)
		ret := legacy.PortTestPlayerFileCall("nox_xxx_plrLoad_41A480", uint32(uintptr(unsafe.Pointer(name))))
		release()
		success := which != "missing" && which != "gui-fails" && which != "music-fails"
		wantRet := uint32(0)
		wantPath := "prior.plr"
		if success {
			wantRet = 1
			wantPath = path
		}
		if ret != wantRet || cryptfile.Global() != nil || alloc.GoStringS(pathField) != wantPath {
			t.Fatal("client load result", which, ret, wantRet, alloc.GoStringS(pathField), wantPath)
		}
		for i := 0; i < 25; i++ {
			wantID, wantFlag := uint32(0), byte(0)
			if which == "missing" {
				wantID = 99
				wantFlag = 255
			}
			if which == "gui" || which == "both" {
				wantID = uint32(i % 6)
				wantFlag = byte(i)
			}
			if q.bar[2*i] != wantID || q.bar[2*i+1] != 0xaabbcc00|uint32(wantFlag) {
				t.Fatal("client load slots", which, i, q.bar[2*i:2*i+2])
			}
		}
		if which != "missing" {
			after, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(after, before) {
				t.Fatal("client load modified source", which, err)
			}
		}
		state := legacy.MusicModule.GetCurrentBlock()
		if which == "music" || which == "both" {
			if state != (music.MusicState{D: 17, Position: 987, MusicIdx: 23, Volume: 75}) {
				t.Fatal("client music", state)
			}
		}
		rows = append(rows, map[string]any{"case": which, "return": ret, "path_updated": success, "slots": append([]uint32(nil), q.bar[:50]...), "music": state})
	}
	spellbookCapture(t, "player-files-client-load-framing", rows, "702d5394b25b7e24191777ec89dc509213e685511eba963212c7c7df222c5d21")
}
