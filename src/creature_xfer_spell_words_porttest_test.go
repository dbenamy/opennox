//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func creatureXferSpellWords(p *mapDrawableStream, named bool, sets ...[]uint32) {
	words := make([]uint32, 137)
	if len(sets) > 0 {
		copy(words, sets[0])
	}
	if !named {
		for _, w := range words {
			p.u32(w)
		}
		return
	}
	count := uint32(0)
	for _, w := range words {
		if w != 0 {
			count++
		}
	}
	p.u32(count)
	for id, w := range words {
		if w != 0 {
			itemXferRewardName(p, spell.ID(id).String())
			p.u32(w)
		}
	}
}
func TestCreatureXferSpellWords(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "spell-words.bin")
	type row struct {
		Case              string
		Wire              []byte
		ReadCRC, WriteCRC uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "creature-xfer-spell-words", rows, "cf46d2dbabd9249ac09dfe091e8e607ef94a841267269850a87bed89cc2c7826")
	}()
	for _, name := range []string{"Monster", "NPC"} {
		for _, version := range []int{32, 33, 34, 35, 60} {
			for _, mode := range []int{0, 1, 2, 3} {
				t.Run(fmt.Sprintf("%s-v%d-mode%d", name, version, mode), func(t *testing.T) {
					words := make([]uint32, 137)
					for i := range words {
						if mode == 3 || mode == 2 && i%2 == 0 || mode == 1 && (i == 1 || i == 136) {
							words[i] = 0x80000001 + uint32(i)
						}
					}
					u := newCreatureXferObject(t, s, name)
					// A named record must clear all absent entries.
					for i := range words {
						objectXferSetWord(u.UpdateData, 1488+4*i, 0x51515151)
					}
					p := new(mapDrawableStream)
					p.u16(uint16(version))
					objectXferEmptyBase(p, int16(version))
					if name == "Monster" {
						creatureXferMonsterHistorical(p, version, words)
					} else {
						creatureXferNPCHistorical(p, version, words)
					}
					if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					if err := u.CallXfer(nil); err != nil {
						t.Fatal(err)
					}
					readCRC := cryptfile.Global().PortTestChecksum()
					pos, err := cryptfile.Global().File.Seek(0, 1)
					cryptfile.Close()
					if err != nil || pos != int64(p.Len()) {
						t.Fatal("spell record position")
					}
					for i, w := range words {
						if got := objectXferGetWord(u.UpdateData, 1488+4*i); got != w {
							t.Fatalf("spell %d=%08x want=%08x", i, got, w)
						}
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
						t.Fatal(err)
					}
					if err := u.CallXfer(nil); err != nil {
						t.Fatal(err)
					}
					writeCRC := cryptfile.Global().PortTestChecksum()
					cryptfile.Close()
					got, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					section := new(mapDrawableStream)
					creatureXferSpellWords(section, true, words)
					if !bytes.Contains(got, section.Bytes()) {
						t.Fatal("current writer spell count/names/values")
					}
					rows = append(rows, row{t.Name(), got, readCRC, writeCRC})
				})
			}
		}
	}
}
