//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"time"
	"unsafe"
)

func TestPlayerFilesMetadataWrite(t *testing.T) {
	_, player := playerFileBookOwner(t)
	player[3684] = 0xfe
	meters := legacy.PortTestNewMeterEnvironment()
	t.Cleanup(meters.Restore)
	const base = 10980
	raw := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 10980), 1284)
	old := bytes.Clone(raw)
	t.Cleanup(func() { copy(raw, old) })
	current := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 36), 32)
	oldMap := bytes.Clone(current)
	t.Cleanup(func() { copy(current, oldMap) })
	clear(current)
	copy(current, "forest.map")
	quest := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1556160), 2)
	oldQuest := append([]uint32(nil), quest...)
	t.Cleanup(func() { copy(quest, oldQuest) })
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{0, 2048, 8192, 8192 | 4096} {
		for _, q := range [][2]uint32{{0, 0}, {1, 0}, {0, 1}} {
			for _, hasPlayer := range []bool{false, true} {
				flags.ResetGame()
				flags.SetGame(gf)
				copy(quest, q[:])
				*meters.NamedWord("dword_8531A0_2576") = 0
				mode := byte(1)
				if hasPlayer {
					*meters.NamedWord("dword_8531A0_2576") = uint32(uintptr(unsafe.Pointer(&player[0])))
					mode = 0xfe
				}
				for i := range raw {
					raw[i] = byte(i*13 + 7)
				}
				binary.LittleEndian.PutUint32(raw, 0xa5a5a5a7)
				putString := func(off int, s string) { copy(raw[off-base:], s); raw[off-base+len(s)] = 0 }
				path, title, name := "Save/slot.plr", "Café Ω", "A😀Ω"
				putString(10984, path)
				putString(12008, title)
				nb := playerFileName(name)
				copy(raw[12204-base:], nb)
				binary.LittleEndian.PutUint16(raw[12204-base+len(nb):], 0)
				expected := bytes.Clone(raw)
				fv := uint32(0xa5a5a5a7)
				if gf&8192 != 0 {
					fv &^= 1
					if gf&4096 != 0 || q[0] != 0 || q[1] != 0 {
						fv |= 4
					} else {
						fv |= 2
					}
				} else {
					fv = fv&^6 | 1
				}
				binary.LittleEndian.PutUint32(expected, fv)
				expected[12255-base] = 0
				expected[12256-base] = mode
				copy(expected[12136-base:], "forest.map\x00")
				clear(expected[12168-base : 12184-base])
				want := binary.LittleEndian.AppendUint16(nil, 12)
				want = binary.LittleEndian.AppendUint32(want, fv)
				want = binary.LittleEndian.AppendUint16(want, uint16(len(path)))
				want = append(want, path...)
				want = append(want, byte(len(title)))
				want = append(want, title...)
				clockOffset := len(want)
				want = append(want, make([]byte, 16)...)
				for _, off := range []int{12187, 12184, 12190, 12193, 12196} {
					want = append(want, raw[off-base:off-base+3]...)
				}
				want = append(want, raw[12199-base:12204-base]...)
				want = append(want, byte(len(nb)/2))
				want = append(want, nb...)
				want = append(want, raw[12254-base], 0, mode, 10)
				want = append(want, "forest.map"...)
				want = append(want, raw[12257-base])
				before := time.Now()
				ret, got, pos := playerFileSection(t, "nox_xxx_parseFileInfoData_41C3B0", nil, 0)
				after := time.Now()
				if len(got) != len(want) {
					t.Fatal("metadata length", len(got), len(want))
				}
				fields := [8]uint16{}
				for i := range fields {
					fields[i] = binary.LittleEndian.Uint16(got[clockOffset+2*i:])
				}
				stamp := time.Date(int(fields[0]), time.Month(fields[1]), int(fields[3]), int(fields[4]), int(fields[5]), int(fields[6]), int(fields[7])*1000000, time.Local)
				if stamp.Before(before.Truncate(time.Millisecond)) || stamp.After(after) || uint16(stamp.Weekday()) != fields[2] || !bytes.Equal(got[clockOffset:clockOffset+16], raw[12168-base:12184-base]) {
					t.Fatal("metadata clock", fields, before, after)
				}
				// Only the identified and bounded wall-clock field is normalized.
				clear(got[clockOffset : clockOffset+16])
				actual := bytes.Clone(raw)
				clear(actual[12168-base : 12184-base])
				if ret != 1 || pos != int64(len(want)) || !bytes.Equal(got, want) || !bytes.Equal(actual, expected) {
					t.Fatal("metadata write", gf, q, hasPlayer, ret, pos)
				}
				rows = append(rows, map[string]any{"flags": uint32(gf), "quest": q, "player": hasPlayer, "return": ret, "position": pos, "bytes": got, "info": actual})
			}
		}
	}
	spellbookCapture(t, "player-files-metadata-write", rows, "")
}
