//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"strings"
	"testing"
	"unsafe"
)

func TestPlayerFilesMetadataRead(t *testing.T) {
	const base = 10980
	raw := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 10980), 1284)
	old := bytes.Clone(raw)
	t.Cleanup(func() { copy(raw, old) })
	current := unsafe.Slice(memmap.PtrUint8(0x85B3FC, 36), 32)
	oldMap := bytes.Clone(current)
	t.Cleanup(func() { copy(current, oldMap) })
	clear(current)
	copy(current, "current.map")
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 10, 11, 12, 13, 0x8000, 0xffff} {
		for _, kind := range []string{"empty", "unicode", "limits"} {
			path, title, name, mapName := "", "", "", ""
			if kind == "unicode" {
				path = "Save/slot.plr"
				title = "Café Ω"
				name = "A😀Ω"
				mapName = "forest.map"
			}
			if kind == "limits" {
				path = strings.Repeat("p", 1023)
				title = strings.Repeat("t", 127)
				name = strings.Repeat("n", 24)
				mapName = strings.Repeat("m", 31)
			}
			for i := range raw {
				raw[i] = 0xa5
			}
			expected := bytes.Clone(raw)
			input := binary.LittleEndian.AppendUint16(nil, version)
			input = binary.LittleEndian.AppendUint32(input, 0x87654321)
			input = binary.LittleEndian.AppendUint16(input, uint16(len(path)))
			input = append(input, path...)
			input = append(input, byte(len(title)))
			input = append(input, title...)
			binary.LittleEndian.PutUint32(expected, 0x87654321)
			putString := func(off int, value string) { copy(expected[off-base:], value); expected[off-base+len(value)] = 0 }
			putString(10984, path)
			putString(12008, title)
			timestamp := []uint16{2024, 2, 4, 29, 23, 58, 57, 999}
			for i, v := range timestamp {
				input = binary.LittleEndian.AppendUint16(input, v)
				binary.LittleEndian.PutUint16(expected[12168-base+2*i:], v)
			}
			for i, off := range []int{12187, 12184, 12190, 12193, 12196} {
				value := []byte{byte(11 + i), byte(101 + i), byte(201 + i)}
				input = append(input, value...)
				copy(expected[off-base:], value)
			}
			for off := 12199; off <= 12203; off++ {
				value := byte(off)
				input = append(input, value)
				expected[off-base] = value
			}
			nameBytes := playerFileName(name)
			input = append(input, byte(len(nameBytes)/2))
			input = append(input, nameBytes...)
			copy(expected[12204-base:], nameBytes)
			binary.LittleEndian.PutUint16(expected[12204-base+len(nameBytes):], 0)
			input = append(input, 0x91, 0x82, 0xf3)
			copy(expected[12254-base:], []byte{0x91, 0x82, 0xf3})
			if int16(version) >= 11 {
				input = append(input, byte(len(mapName)))
				input = append(input, mapName...)
				putString(12136, "current.map")
				putString(12136, mapName)
			}
			expected[12257-base] = 0
			if int16(version) >= 12 {
				input = append(input, 0xe4)
				expected[12257-base] = 0xe4
			}
			wantPos := int64(len(input))
			wantRet := uint32(1)
			if int16(version) > 12 {
				wantPos = 2
				wantRet = 0
				expected = bytes.Clone(raw)
			}
			input = append(input, 0xde, 0xad, 0xbe, 0xef)
			ret, got, pos := playerFileSection(t, "nox_xxx_parseFileInfoData_41C3B0", input, 0)
			if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) {
				t.Fatal("metadata read", version, kind, ret, wantRet, pos, wantPos)
			}
			rows = append(rows, map[string]any{"version": version, "case": kind, "return": ret, "position": pos, "info": bytes.Clone(raw)})
		}
	}
	spellbookCapture(t, "player-files-metadata-read", rows, "")
}
