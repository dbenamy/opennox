//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/spell"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func playerFileQuickbarBytes(class byte) []byte {
	out := []byte{class}
	appendSlot := func(i int, flag byte) {
		name := spell.ID(i % 6).String()
		if class == 0 {
			name = server.Ability(i % 6).String()
		}
		out = append(out, byte(len(name)))
		out = append(out, name...)
		out = append(out, flag)
	}
	for i := 0; i < 25; i++ {
		appendSlot(i, byte(i))
	}
	if class != 0 {
		for row := 0; row < 3; row++ {
			for slot := 0; slot < 3; slot++ {
				i := row*5 + slot
				appendSlot(i, byte(0x40+i))
			}
		}
	}
	return out
}
func TestPlayerFilesGUIWrite(t *testing.T) {
	q := newQuickbarOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, class := range []byte{0, 1, 2} {
		for _, row := range []uint32{0, 4} {
			for _, trapRow := range []uint32{0, 2} {
				q.reset(t)
				*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = class
				*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 3648)) = 7
				traps := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1047940), 50)
				for i := 0; i < 25; i++ {
					q.bar[2*i] = uint32(i % 6)
					q.bar[2*i+1] = 0x11223300 | uint32(i)
					traps[2*i] = uint32(i % 6)
					traps[2*i+1] = 0x22334440 | uint32(i)
				}
				q.call("nox_xxx_clientUpdateButtonRow_45E110", row)
				q.call("nox_client_trapSetSelect_4604B0", trapRow)
				want := append([]byte{3, 0}, playerFileQuickbarBytes(class)...)
				want = append(want, byte(row), byte(trapRow), 7)
				ret, got, pos := playerFileSection(t, "sub_41C280", nil, 0)
				if ret != 1 || pos != int64(len(want)) || !bytes.Equal(got, want) {
					t.Fatal("GUI write", class, row, trapRow, ret, pos, len(want))
				}
				rows = append(rows, map[string]any{"class": class, "row": row, "trap_row": trapRow, "return": ret, "bytes": got, "position": pos})
			}
		}
	}
	spellbookCapture(t, "player-files-gui-write", rows, "e3c0a77762c7f303349b5cfa1915a4d746599f17e179f6478584bd621cf13fca")
}
func TestPlayerFilesGUIRead(t *testing.T) {
	q := newQuickbarOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	// Spell-class reads avoid ability name-table ownership already covered by
	// TestQuickbarSaveOwners; this tests the enclosing section's version gates.
	for _, version := range []uint16{0, 1, 2, 3, 4, 0x8000, 0xffff} {
		q.reset(t)
		input := binary.LittleEndian.AppendUint16(nil, version)
		input = append(input, playerFileQuickbarBytes(1)...)
		if int16(version) >= 2 {
			input = append(input, 3, 2)
		}
		if int16(version) >= 3 {
			input = append(input, 9)
		}
		wantRet, wantPos := uint32(1), int64(len(input))
		if int16(version) > 3 {
			wantRet = 0
			wantPos = 2
		}
		input = append(input, 0xde, 0xad, 0xbe, 0xef)
		ret, got, pos := playerFileSection(t, "sub_41C280", input, 0)
		if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) {
			t.Fatal("GUI read", version, ret, pos, wantPos)
		}
		selected, trap := uint32(0), uint32(0)
		if int16(version) <= 3 && int16(version) >= 2 {
			selected = 3
			trap = 2
		}
		if q.call("nox_xxx_buttonsGetSelectedRow_45E180") != selected || q.call("sub_4604E0") != trap {
			t.Fatal("GUI selections", version)
		}
		if wantRet == 1 {
			for i := 0; i < 25; i++ {
				if q.bar[2*i] != uint32(i%6) || byte(q.bar[2*i+1]) != byte(i) {
					t.Fatal("GUI slots", version, i)
				}
			}
		}
		rows = append(rows, map[string]any{"version": version, "return": ret, "position": pos, "row": selected, "trap_row": trap})
	}
	spellbookCapture(t, "player-files-gui-read", rows, "ed2464d247bfd1dc8e7f9e3f3715b6fe061a8ae72b8b0647d9b352401d697416")
}
