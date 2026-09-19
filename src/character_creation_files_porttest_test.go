//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/datapath"
	"github.com/opennox/libs/ifs"
	"github.com/opennox/opennox/v1/client/gui"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/music"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestCharacterCreationFiles(t *testing.T) {
	q := newQuickbarOwner(t)
	q.reset(t) // Reset the shared GUI before creating character controls.
	colors := characterPalette(t)
	q.entryOwner.create(t, 0, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) { d.Field_1040 = 25 })
	var controls []*gui.Window
	var ptrs []unsafe.Pointer
	for i := 0; i < 14; i++ {
		w := q.c.GUI.NewWindowRaw(q.parent, 8, 0, 0, 1, 1, nil)
		controls = append(controls, w)
		ptrs = append(ptrs, w.C())
		*characterWord(w, 32) = uint32(i%3) | uint32(3+i)<<16
	}
	ptrs = append(ptrs, q.entryOwner.win.C())
	player, free := alloc.Calloc(1, 128)
	defer free()
	defer legacy.PortTestCharacterAppearanceOwner(player, ptrs)()
	t.Cleanup(legacy.PortTestPlayerFileClientSections())
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) }()
	raw := serverConfigOwnBytes(t, 0x85B3FC, 10980, 1280)
	current := serverConfigOwnBytes(t, 0x85B3FC, 36, 32)
	mapPath := serverConfigOwnBytes(t, 0x5D4594, 2598188, 80)
	_, restoreConfig := legacy.PortTestServerOptionsWords()
	defer restoreConfig()
	// Restore the shipped directory fragments used by the constructor.
	fragment := serverConfigOwnBytes(t, 0x587000, 171764, 20)
	copy(fragment, []byte("\\Save\\\x00\x00WORKING\x00\\\x00\x00\x00"))
	words, restore := legacy.PortTestAudioEventGlobals()
	defer restore()
	*words["dword_5d4594_816368"] = 0
	*words["dword_5d4594_816372"] = 0
	oldMusic := legacy.MusicModule
	legacy.MusicModule = &music.Module{}
	defer func() { legacy.MusicModule = oldMusic }()
	oldData := datapath.Data()
	defer datapath.SetData(oldData)
	oldDir, err := ifs.Workdir()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := ifs.Chdir(oldDir); err != nil {
			t.Error(err)
		}
	}()
	defer flags.PortTestGameFlags(0)()
	type row struct {
		Class, Mask, Occupied, Mode, Return int
		Name                                string
		Info                                []byte
		Payloads                            [][]byte
	}
	var rows []row
	type fileCase struct{ class, mask, occupied, mode int }
	var cases []fileCase
	for class := 0; class < 3; class++ {
		for _, mask := range []int{0, 5, 15} {
			for _, occupied := range []int{0, 1, 99, 100} {
				cases = append(cases, fileCase{class, mask, occupied, 0})
			}
		}
		cases = append(cases, fileCase{class, 15, 0, 2048}, fileCase{class, 15, -1, 0})
	}
	for _, tc := range cases {
		class, mask, occupied, mode := tc.class, tc.mask, tc.occupied, tc.mode
		clear(q.bar[:50])
		flags.ResetGame()
		flags.SetGame(flags.GameFlag(mode))
		clear(mapPath)
		*(*byte)(unsafe.Add(player, 66)) = byte(class)
		*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = byte(class)
		clear(current)
		copy(current, "forest.map")
		clear(raw)
		q.entryOwner.text([]uint16{' ', 'H', 'e', 'r', 'o', ' '})
		for i := 0; i < 4; i++ {
			controls[10+i].Flags = 0
			if mask&(1<<i) != 0 {
				controls[10+i].Flags = 8
			}
		}
		dir := t.TempDir()
		datapath.SetData(dir)
		save := filepath.Join(dir, "Save")
		if occupied < 0 {
			if err := os.WriteFile(save, []byte("block directory"), 0600); err != nil {
				t.Fatal(err)
			}
		} else {
			if mode == 2048 {
				save = filepath.Join(save, "WORKING")
			}
			if err := os.MkdirAll(save, 0700); err != nil {
				t.Fatal(err)
			}
			if mode == 2048 {
				if err := os.WriteFile(filepath.Join(save, "stale"), []byte("old"), 0600); err != nil {
					t.Fatal(err)
				}
			}
		}
		if err := ifs.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < occupied; i++ {
			if err := os.WriteFile(filepath.Join(save, fmt.Sprintf("H%02d.plr", i)), []byte("existing"), 0600); err != nil {
				t.Fatal(err)
			}
		}
		ret := legacy.PortTestCharacterCreateFile()
		wantRet := 1
		if occupied == 100 || occupied < 0 {
			wantRet = 0
		}
		if ret != wantRet {
			t.Fatal("create result", class, mask, occupied, ret)
		}
		wd, err := ifs.Workdir()
		if err != nil || ifs.Normalize(wd) != ifs.Normalize(dir) {
			t.Fatal("working directory not restored", wd, dir, err)
		}
		if occupied == 100 {
			if !bytes.Equal(raw, make([]byte, 1280)) {
				t.Fatal("exhausted names reached writer")
			}
			rows = append(rows, row{Class: class, Mask: mask, Occupied: occupied, Return: ret})
			continue
		}
		slot := occupied
		if slot < 0 {
			slot = 0
		}
		name := fmt.Sprintf("H%02d.plr", slot)
		if mode == 2048 {
			name = "Player.plr"
			if _, err := os.Stat(filepath.Join(save, "stale")); !os.IsNotExist(err) {
				t.Fatal("working cleanup", err)
			}
			wantMap := []string{"war01a", "wiz01a", "con01a"}[class]
			if alloc.GoStringS(current) != wantMap {
				t.Fatal("campaign map", alloc.GoStringS(current), wantMap)
			}
		}
		path := filepath.Join(save, name)
		if ifs.Normalize(alloc.GoStringS(raw[4:1028])) != ifs.Normalize(path) {
			t.Fatal("filename", alloc.GoStringS(raw[4:1028]), path)
		}
		if !bytes.Equal(raw[1224:1234], append(playerFileName("Hero"), 0, 0)) || raw[1274] != byte(class) {
			t.Fatal("name/class record")
		}
		skin := colors[3*3:][:3]
		if !bytes.Equal(raw[1204:1207], skin) {
			t.Fatal("saved skin")
		}
		for i, off := range []int{1207, 1210, 1213, 1216} {
			want := skin
			if mask&(1<<i) != 0 {
				j := i + 1
				index := 32*(j%3) + 3 + j
				want = colors[3*index:][:3]
			}
			if !bytes.Equal(raw[off:off+3], want) {
				t.Fatal("saved override", mask, i)
			}
		}
		for i := 0; i < 5; i++ {
			if raw[1219+i] != byte(8+i) {
				t.Fatal("saved index", i)
			}
		}
		if occupied < 0 {
			if cryptfile.Global() != nil {
				t.Fatal("failed writer retained file")
			}
			info := append([]byte(nil), raw...)
			// Original C leaves these nonserialized tail padding bytes uninitialized.
			clear(info[1278:1280])
			clear(info[4:1028])
			rows = append(rows, row{class, mask, occupied, mode, ret, name, info, nil})
			continue
		}
		f, err := cryptfile.OpenFile(path, cryptfile.ReadOnly, 27)
		if err != nil {
			t.Fatal(err)
		}
		var payloads [][]byte
		for _, idWant := range []uint32{7, 1, 12} {
			id, err := f.ReadU32()
			if err != nil || id != idWant {
				t.Fatal("section", id, idWant, err)
			}
			n, err := f.ReadAlignedU32()
			if err != nil || n > 4096 {
				t.Fatal("length", n, err)
			}
			data := make([]byte, n)
			if _, err := io.ReadFull(f, data); err != nil {
				t.Fatal(err)
			}
			if id == 1 {
				pathN := int(binary.LittleEndian.Uint16(data[6:]))
				if string(data[8:8+pathN]) != alloc.GoStringS(raw[4:1028]) {
					t.Fatal("metadata pathname")
				}
				title := 8 + pathN
				clock := title + 1 + int(data[title])
				if !bytes.Equal(data[clock:clock+16], raw[1188:1204]) {
					t.Fatal("metadata clock")
				}
				clear(data[clock : clock+16])
				normalized := append([]byte(nil), data[:6]...)
				normalized = append(normalized, 0, 0)
				data = append(normalized, data[8+pathN:]...)
			}
			payloads = append(payloads, data)
		}
		end, err := f.ReadU32()
		if err != nil || end != 0 {
			t.Fatal("file terminator", end, err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		info := append([]byte(nil), raw...)
		// Original C leaves these nonserialized tail padding bytes uninitialized.
		clear(info[1278:1280])
		clear(info[4:1028])
		clear(info[1188:1204])
		rows = append(rows, row{class, mask, occupied, mode, ret, name, info, payloads})
	}
	spellbookCapture(t, "character-creation-files", rows, "18408533122efd720371bb4999e81fcd60d24456d4f248065c1edb7f78d06055")
}
