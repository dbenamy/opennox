//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// The section writer only needs the real sparse player registry and C-owned
// player/unit records. Award paths use the fuller gameplay owners separately.
func playerFileBookOwner(t *testing.T) (*server.Object, []byte) {
	t.Helper()
	s := new(server.Server)
	units, configure, _, free := s.PortTestEscortPlayers()
	t.Cleanup(free)
	configure(1)
	u := &units[0]
	p := u.UpdateDataPlayer().Player
	u.NetCode = 1001
	p.NetCodeVal = 1001
	old := legacy.GetServer
	legacy.GetServer = func() legacy.Server { return &playerFileServerOwner{s: s} }
	t.Cleanup(func() { legacy.GetServer = old })
	return u, unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
}

func TestPlayerFilesSpellbookWrite(t *testing.T) {
	u, raw := playerFileBookOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, class := range []byte{0, 1, 2} {
		for _, gf := range []flags.GameFlag{0, 2048, 4096, 8192, 8192 | 4096, 8192 | 2048} {
			for _, pattern := range []string{"empty", "edges", "all"} {
				flags.ResetGame()
				flags.SetGame(gf)
				raw[2251] = class
				clear(raw[3700:4244])
				limit := 136
				if class == 0 {
					limit = 5
				}
				for id := 1; id <= 136; id++ {
					if pattern == "all" || pattern == "edges" && (id == 1 || id == limit || id == limit+1 || id == 136) {
						binary.LittleEndian.PutUint32(raw[3700+4*(id-1):], uint32(id*17+1))
					}
				}
				before := bytes.Clone(raw)
				present := byte(1)
				if gf&8192 != 0 && gf&4096 == 0 {
					present = 0
				}
				want := []byte{3, 0, present}
				retWant := uint32(1)
				if present != 0 {
					if gf&(2048|4096) == 0 {
						retWant = 0
					} else {
						var entries []byte
						count := byte(0)
						for id := 1; id <= limit; id++ {
							level := binary.LittleEndian.Uint32(raw[3700+4*(id-1):])
							if level == 0 {
								continue
							}
							name := spell.ID(id).String()
							if class == 0 {
								name = server.Ability(id).String()
							}
							if name == "" {
								t.Fatal("missing canonical name", class, id)
							}
							count++
							entries = append(entries, byte(len(name)))
							entries = append(entries, name...)
							entries = binary.LittleEndian.AppendUint32(entries, level)
						}
						want = append(want, count)
						want = append(want, entries...)
					}
				}
				ret, got, pos := playerFileSection(t, "nox_xxx_guiSpellbook_41B660", nil, uint32(uintptr(unsafe.Pointer(u))), 0)
				if ret != retWant || !bytes.Equal(got, want) || pos != int64(len(want)) || !bytes.Equal(raw, before) {
					t.Fatal("spellbook write", class, gf, pattern, ret, pos, len(want))
				}
				rows = append(rows, map[string]any{"class": class, "flags": uint32(gf), "pattern": pattern, "return": ret, "bytes": got, "position": pos})
			}
		}
	}
	spellbookCapture(t, "player-files-spellbook-write", rows, "28b567c24791640d4fe503a790bfc18b16eb5b896ebdb0506b297ba134462acd")
}

func TestPlayerFilesBookReadGates(t *testing.T) {
	u, raw := playerFileBookOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, section := range []struct {
		name     string
		version  uint16
		maxCount byte
	}{
		{"nox_xxx_guiFieldbook_41B420", 1, 41}, {"nox_xxx_guiSpellbook_41B660", 3, 137},
	} {
		for _, version := range []uint16{0, 1, 2, 3, 4, 0x8000, 0xffff} {
			for _, gf := range []flags.GameFlag{0, 2048, 4096, 8192, 8192 | 4096} {
				for _, present := range []byte{0, 1} {
					for _, count := range []byte{0, section.maxCount + 1, 255} {
						flags.ResetGame()
						flags.SetGame(gf)
						input := binary.LittleEndian.AppendUint16(nil, version)
						input = append(input, present, count, 0xde, 0xad, 0xbe, 0xef)
						wantRet, wantPos := uint32(1), int64(4)
						if int16(version) > int16(section.version) {
							wantRet = 0
							wantPos = 2
						} else if present == 0 {
							wantPos = 3
						} else if gf&(2048|4096) == 0 {
							wantRet = 0
							wantPos = 3
						} else if count > section.maxCount {
							wantRet = 0
						}
						before := bytes.Clone(raw)
						ret, got, pos := playerFileSection(t, section.name, input, uint32(uintptr(unsafe.Pointer(u))), 0)
						if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || !bytes.Equal(raw, before) {
							t.Fatal("book gate", section.name, version, gf, present, count, ret, pos, wantRet, wantPos)
						}
						rows = append(rows, map[string]any{"section": section.name, "version": version, "flags": uint32(gf), "present": present, "count": count, "return": ret, "position": pos})
					}
				}
			}
		}
		u.NetCode = 9999
		ret, got, pos := playerFileSection(t, section.name, nil, uint32(uintptr(unsafe.Pointer(u))), 0)
		if ret != 0 || len(got) != 0 || pos != 0 {
			t.Fatal("missing player", section.name)
		}
		u.NetCode = 1001
		rows = append(rows, map[string]any{"section": section.name, "missing_player": true, "return": ret, "position": pos})
	}
	spellbookCapture(t, "player-files-book-read-gates", rows, "700a93921052f26facf5e0c41c71e32e96b5b3519004bf4f8fcf5649710bf516")
}
