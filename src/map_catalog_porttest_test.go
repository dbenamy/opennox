//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/opennox/libs/datapath"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func mapCatalogCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if base := os.Getenv("OPENNOX_MAP_CATALOG_CAPTURE"); base != "" {
		if err := os.WriteFile(base+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Log(label, hash)
	if want != "" && hash != want {
		t.Fatalf("%s capture changed: %s", label, hash)
	}
}
func TestMapCatalogOrdering(t *testing.T) {
	var rows [][]legacy.PortTestMapCatalogEntry
	for _, size := range []int{0, 1, 2, 7, 32, 128, 129, 180} {
		for order := 0; order < 4; order++ {
			f := legacy.PortTestMapCatalogOpen(1)
			var input []legacy.PortTestMapCatalogEntry
			for i := 0; i < size; i++ {
				n := i
				if order == 1 {
					n = size - 1 - i
				}
				if order == 2 {
					n = i * 37 % 19
				}
				if order == 3 {
					n = i % 5
				}
				name := fmt.Sprintf("map%03d", n)
				if i%3 == 0 {
					name = strings.ToUpper(name)
				}
				if i%11 == 0 {
					name = ""
				}
				if i%13 == 0 {
					name = "\x80\xff" + name
				}
				v := legacy.PortTestMapCatalogEntry{Name: name, Enabled: int32(i % 3), Flags: uint32(i), Tag: uint16(i)}
				input = append(input, v)
				f.Add(v)
			}
			got := f.Entries()
			// C inserts duplicates before existing equal entries; model this by sorting
			// by bytewise name then decreasing insertion identity.
			sort.Slice(input, func(i, j int) bool {
				if input[i].Name != input[j].Name {
					return input[i].Name < input[j].Name
				}
				return input[i].Tag > input[j].Tag
			})
			if len(got) != len(input) {
				t.Fatal("catalog count")
			}
			for i := range got {
				if got[i] != input[i] {
					t.Fatalf("size %d order %d row %d: %#v != %#v", size, order, i, got[i], input[i])
				}
			}
			rows = append(rows, got)
			f.Close()
		}
	}
	mapCatalogCapture(t, "ordering", rows, "01ab41e61d7721e87c34da3cb50be72bdcc3e010807f822bd35e5d0966e9f552")
}
func TestMapCatalogQuestGroups(t *testing.T) {
	type row struct {
		Result  int
		Entries []legacy.PortTestMapCatalogEntry
		State   legacy.PortTestMapCatalogState
	}
	var rows []row
	for _, size := range []int{0, 1, 2, 7, 127, 128, 129, 180} {
		f := legacy.PortTestMapCatalogOpen(1)
		for i := 0; i < size; i++ {
			name := fmt.Sprintf("q%05d%02d", i/3, i%3)
			if i%4 == 0 {
				name = strings.ToUpper(name)
			}
			flags := uint32(2)
			if i%7 == 0 {
				flags = 0
			}
			enabled := int32(1)
			if i%11 == 0 {
				enabled = 0
			}
			f.Add(legacy.PortTestMapCatalogEntry{Name: name, Enabled: enabled, Flags: flags, Tag: uint16(i)})
		}
		n := f.BuildQuest()
		s := f.State()
		entries := f.Entries()
		var expected []string
		for _, v := range entries {
			if v.Enabled != 0 && v.Flags&2 != 0 && len(expected) < 128 {
				expected = append(expected, v.Name+".map")
			}
		}
		if n != len(expected) || int(s.Count) != n {
			t.Fatal("quest count")
		}
		groups := map[string]uint32{}
		next := uint32(1)
		for i, name := range expected {
			data := s.Quest[i*32 : (i+1)*32]
			j := bytes.IndexByte(data[4:], 0)
			if j < 0 || string(data[4:4+j]) != name {
				t.Fatalf("quest name %d", i)
			}
			key := strings.ToLower(name[:6])
			id, ok := groups[key]
			if !ok {
				id = next
				next++
				groups[key] = id
			}
			if binary.LittleEndian.Uint32(data) != id {
				t.Fatalf("quest group %d", i)
			}
		}
		f.ResetQuest()
		after := f.State()
		if after.Clock != 1000 {
			t.Fatal("quest clock reset")
		}
		rows = append(rows, row{n, entries, after})
		f.Close()
	}
	mapCatalogCapture(t, "groups", rows, "56e58d720d5804f504ffe730cf1a194d5dec4288f3ae03d77cddb6dba51de228")
}
func TestMapCatalogQuestSequences(t *testing.T) {
	type row struct {
		Name  string
		RNG   int
		State legacy.PortTestMapCatalogState
	}
	var rows []row
	for _, size := range []int{0, 1, 2, 3, 7, 32, 128} {
		for seed := 1; seed <= 5; seed++ {
			f := legacy.PortTestMapCatalogOpen(seed)
			for i := 0; i < size; i++ {
				f.Add(legacy.PortTestMapCatalogEntry{Name: fmt.Sprintf("q%05d%02d", i/2, i%2), Enabled: 1, Flags: 2})
			}
			f.BuildQuest()
			f.ResetQuest()
			if seed >= 4 {
				s := f.State()
				s.Clock = 0xfffffffd
				s.Last = 0
				for i := 0; i < size; i++ {
					binary.LittleEndian.PutUint32(s.Quest[i*32+24:], []uint32{0, 0x7fffffff, 0x80000000, 0xffffffff}[i%4])
					binary.LittleEndian.PutUint32(s.Quest[i*32+28:], s.Clock-uint32(i%8))
				}
				f.RestoreState(s)
			}
			for step := 0; step < 30; step++ {
				before := f.State()
				name, index := f.ChooseQuest()
				after := f.State()
				wantIndex := seed
				if size > 1 {
					wantIndex += step + 1
				}
				if index != wantIndex {
					t.Fatalf("size %d seed %d step %d RNG index %d want %d", size, seed, step, index, wantIndex)
				}
				if !bytes.Equal(before.Quest, after.Quest) || before.Clock != after.Clock || before.Last != after.Last {
					t.Fatal("choosing mutated history")
				}
				if size == 0 && name != "" {
					t.Fatal("empty quest catalog result")
				}
				if size > 0 && !strings.HasSuffix(name, ".map") {
					t.Fatal("quest result")
				}
				rows = append(rows, row{name, index, after})
				if step%5 == 0 {
					missing := "absent.map"
					f.Played(&missing)
					f.Played(nil)
				}
				if name != "" {
					name = strings.ToUpper(name)
					f.Played(&name)
				}
			}
			f.Close()
		}
	}
	mapCatalogCapture(t, "sequences", rows, "0d88809833169ae07f681ecaf7cfa2504aa2515d5acf5a6d8bf669afb668cc01")
}
func TestMapCycleGroups(t *testing.T) {
	f := legacy.PortTestMapCatalogOpen(1)
	defer f.Close()
	type row struct {
		Flags           uint32
		Group, Set, Get int
	}
	var rows []row
	masks := make([]uint32, 6)
	for i := range masks {
		masks[i] = legacy.PortTestMapCycleMask(i)
	}
	flags := []uint32{0, 0xffffffff, 0x80000000}
	for bit := 0; bit < 32; bit++ {
		flags = append(flags, 1<<bit)
	}
	for _, a := range masks {
		for _, b := range masks {
			flags = append(flags, a|b)
		}
	}
	for _, flag := range flags {
		expected := 0
		for i, m := range masks {
			if flag&m != 0 {
				expected = i
				break
			}
		}
		got := legacy.PortTestMapCycleFlagGroup(flag)
		if got != expected {
			t.Fatalf("flags %x group %d want %d", flag, got, expected)
		}
		for _, index := range []int{-1, 0, 1, 24, 25, 26, 2147483647, -2147483648} {
			set := legacy.PortTestMapCycleSetIndex(flag, index)
			get := legacy.PortTestMapCycleGetIndex(flag)
			if set != expected || get != index {
				t.Fatal("cycle index roundtrip")
			}
			rows = append(rows, row{flag, got, set, get})
		}
	}
	mapCatalogCapture(t, "cycle-groups", rows, "a5c5ff64353809bc4e773bc2a4c205910eefd314523666d94c7e71b031cdd0ee")
}
func TestMapCycleLineEndings(t *testing.T) {
	var rows [][]byte
	for _, text := range []string{"", "map", "a\rb", "a\nb", "a\r\nb", "a\n\rb", "\r\n", "Map  \t", "x.y.map"} {
		input := append([]byte(text), 0, 0x93, 0x82)
		want := bytes.Clone(input)
		if i := strings.IndexByte(text, '\r'); i >= 0 {
			want[i] = 0
		}
		prefix := want[:bytes.IndexByte(want, 0)]
		if i := bytes.IndexByte(prefix, '\n'); i >= 0 {
			want[i] = 0
		}
		got := legacy.PortTestMapCycleStrip(input)
		if !bytes.Equal(got, want) {
			t.Fatalf("line %q: %x != %x", text, got, want)
		}
		rows = append(rows, got)
	}
	if legacy.PortTestMapCycleStrip(nil) != nil {
		t.Fatal("nil line")
	}
	mapCatalogCapture(t, "line-endings", rows, "6b737614e1b4b4fefa4073871534ace7489c92bc0b1cdf1c073626820742860c")
}

func TestMapCycleControlsAndNext(t *testing.T) {
	oldGame, oldEngine := noxflags.GetGame(), noxflags.GetEngine()
	defer func() {
		noxflags.ResetGame()
		noxflags.SetGame(oldGame)
		noxflags.ResetEngine()
		noxflags.SetEngine(oldEngine)
	}()
	f := legacy.PortTestMapCatalogOpen(1)
	defer f.Close()
	var rows []any
	for _, render := range []bool{false, true} {
		noxflags.ResetEngine()
		if render {
			noxflags.SetEngine(noxflags.EngineNoRendering)
		}
		for _, v := range []int{0, 1, 2, -1, -2147483648, 2147483647} {
			set := legacy.PortTestMapCycleEnable(v)
			got := legacy.PortTestMapCycleEnabled()
			want := 0
			if render || v != 0 {
				want = 1
			}
			if set != v || got != want {
				t.Fatal("cycle enable contract")
			}
			rows = append(rows, []int{set, got})
		}
	}
	for group := 0; group < 6; group++ {
		noxflags.ResetGame()
		noxflags.SetGame(noxflags.GameFlag(legacy.PortTestMapCycleMask(group)))
		for _, count := range []int{0, 1, 2, 24, 25} {
			for _, delta := range []int{-1, 0, 1} {
				index := count + delta
				if index < 0 {
					index = 0
				}
				// The last group's one-past row aliases the count table. Test the strict
				// boundary on groups with another allocated row instead of reading its bytes
				// as a string beyond the name array.
				if group == 5 && count == 25 && index == count {
					continue
				}
				s := f.State()
				clear(s.Cycle)
				s.Counts = [6]uint32{}
				s.Indices = [6]uint32{}
				s.Counts[group] = uint32(count)
				s.Indices[group] = uint32(index)
				for i := 0; i < 25; i++ {
					copy(s.Cycle[(group*25+i)*128:], fmt.Sprintf("g%d-map%02d", group, i))
				}
				f.RestoreState(s)
				got := legacy.PortTestMapCycleNext()
				after := f.State()
				expectedIndex := index
				if index > count {
					expectedIndex = 0
				}
				want := ""
				next := uint32(index)
				if count > 0 {
					if expectedIndex < 25 {
						want = fmt.Sprintf("g%d-map%02d", group, expectedIndex)
					}
					next = uint32(expectedIndex+1) % uint32(count)
				}
				if got != want || after.Indices[group] != next {
					t.Fatalf("group %d count %d index %d: %q next %d; want %q next %d", group, count, index, got, after.Indices[group], want, next)
				}
				rows = append(rows, []any{group, count, index, got, after.Indices})
			}
		}
	}
	legacy.PortTestMapCycleReset()
	s := f.State()
	if s.Counts != [6]uint32{} || s.Indices != [6]uint32{} {
		t.Fatal("cycle reset")
	}
	mapCatalogCapture(t, "controls", rows, "989631b482535c3baa34027db4da1b06143a3adf67d3d7dd0353432438ce745e")
}
func TestMapCycleRealFiles(t *testing.T) {
	assets := os.Getenv("OPENNOX_MAP_CATALOG_ASSETS")
	if assets == "" {
		t.Skip("asset path not configured")
	}
	dir := t.TempDir()
	if err := os.Symlink(filepath.Join(assets, "maps"), filepath.Join(dir, "maps")); err != nil {
		t.Fatal(err)
	}
	old := datapath.Data()
	datapath.SetData(dir)
	defer datapath.SetData(old)
	closeHandles := handles.PortTestInit()
	defer closeHandles()
	metadata := unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464)
	saved := bytes.Clone(metadata)
	defer copy(metadata, saved)
	f := legacy.PortTestMapCatalogOpen(1)
	defer f.Close()
	files, err := filepath.Glob(filepath.Join(assets, "maps", "*", "*.nxz"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no real map fixtures")
	}
	type entry struct {
		name  string
		flags uint32
	}
	var maps []entry
	for _, p := range files {
		name := strings.TrimSuffix(filepath.Base(p), ".nxz")
		if err := nox_common_checkMapFile(name); err != nil {
			t.Fatalf("map %s: %v", name, err)
		}
		maps = append(maps, entry{name, uint32(nox_mapToGameFlags(int(memmap.Uint32(0x973F18, 3800))))})
	}
	headers := []string{"[ELIMINATION]", "[DEATHMATCH]", "[CAPTURE THE FLAG]", "[KING OF THE REALM]", "[FLAGBALL]", "[QUEST]"}
	var rows []legacy.PortTestMapCatalogState
	for mode := 0; mode < 4; mode++ {
		var text strings.Builder
		var expected [6][]string
		if mode == 1 {
			text.WriteString("not-a-real-map\n")
			expected[1] = append(expected[1], "not-a-real-map")
		}
		if mode == 2 {
			text.WriteString("\n")
			expected[1] = append(expected[1], "")
		}
		for group, header := range headers {
			if mode == 3 {
				header = strings.ToLower(header)
			}
			text.WriteString(header + "\r\n")
			text.WriteString("\nmissing-map-for-port-test.map\n")
			for repeat := 0; repeat < 3; repeat++ {
				for _, m := range maps {
					name := m.name + ".map"
					text.WriteString(name + "\n")
					if m.flags&legacy.PortTestMapCycleMask(group) != 0 && len(expected[group]) < 25 {
						expected[group] = append(expected[group], name)
					}
				}
			}
		}
		// A final non-newline-terminated line is consumed but rejected by the existing
		// fgets facade; preserve this behavior independently of scanner defaults.
		text.WriteString(maps[0].name + ".map")
		if err := os.WriteFile(filepath.Join(dir, "mapcycle.txt"), []byte(text.String()), 0600); err != nil {
			t.Fatal(err)
		}
		beforeHandles := legacy.PortTestMapCycleOpenHandles()
		legacy.PortTestMapCycleLoad()
		if legacy.PortTestMapCycleOpenHandles() != beforeHandles {
			t.Fatal("cycle file handle leaked")
		}
		if legacy.PortTestMapCyclePath() != dir+"\\mapcycle.txt" {
			t.Fatal("cycle path")
		}
		s := f.State()
		for group, names := range expected {
			if int(s.Counts[group]) != len(names) {
				t.Fatalf("mode %d group %d count %d want %d", mode, group, s.Counts[group], len(names))
			}
			for i, name := range names {
				raw := s.Cycle[(group*25+i)*128 : (group*25+i+1)*128]
				j := bytes.IndexByte(raw, 0)
				if j < 0 || string(raw[:j]) != name {
					t.Fatalf("mode %d group %d name %d", mode, group, i)
				}
			}
		}
		rows = append(rows, s)
	}
	if err := os.Remove(filepath.Join(dir, "mapcycle.txt")); err != nil {
		t.Fatal(err)
	}
	legacy.PortTestMapCycleEnable(1)
	legacy.PortTestMapCycleLoad()
	if legacy.PortTestMapCycleEnabled() != 0 {
		t.Fatal("missing file should disable map cycle")
	}
	mapCatalogCapture(t, "files", rows, "3d43cfa5d73332d3f5f537b0b4022169a490109cdc460252bfbb99aeda66b9a4")
}

func TestMapCycleFirstLineAndLength(t *testing.T) {
	dir := t.TempDir()
	old := datapath.Data()
	datapath.SetData(dir)
	defer datapath.SetData(old)
	closeHandles := handles.PortTestInit()
	defer closeHandles()
	type row struct {
		Counts [6]uint32
		First  []byte
	}
	var rows []row
	inputs := []string{"", "\n", "first\n", "first", "\nlast", "first\x00tail\n", "\r\n", "first\rbreak\n", "\x80\xff\n"}
	for _, size := range []int{124, 125, 126, 127, 128, 129, 255, 1024} {
		inputs = append(inputs, strings.Repeat("x", size)+"\n")
	}
	for _, input := range inputs {
		f := legacy.PortTestMapCatalogOpen(1)
		if err := os.WriteFile(filepath.Join(dir, "mapcycle.txt"), []byte(input), 0600); err != nil {
			f.Close()
			t.Fatal(err)
		}
		beforeHandles := legacy.PortTestMapCycleOpenHandles()
		legacy.PortTestMapCycleLoad()
		s := f.State()
		if legacy.PortTestMapCycleOpenHandles() != beforeHandles {
			f.Close()
			t.Fatal("cycle handle leaked")
		}
		var expected [6]uint32
		want := ""
		if end := strings.IndexByte(input, '\n'); end >= 0 {
			expected[1] = 1
			want = input[:end]
			if len(want) > 126 {
				want = want[:126]
			}
			if end := strings.IndexAny(want, "\x00\r"); end >= 0 {
				want = want[:end]
			}
		}
		raw := s.Cycle[25*128 : 26*128]
		end := bytes.IndexByte(raw, 0)
		if s.Counts != expected || end < 0 || string(raw[:end]) != want {
			f.Close()
			t.Fatalf("first-line contract: counts %v want %v, bytes %x want %x", s.Counts, expected, raw[:min(128, len(want)+2)], want)
		}
		rows = append(rows, row{s.Counts, raw})
		f.Close()
	}
	mapCatalogCapture(t, "first-lines", rows, "6e52285522fdac4cdffb7b4195468e3a721b4642f8d075df7f196da3d290d140")
}
