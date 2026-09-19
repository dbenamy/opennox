//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"testing"
	"unsafe"
)

func TestPlayerFilesJournalGates(t *testing.T) {
	o := newJournalOwner(t)
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 0x8000, 0xffff} {
		for _, gf := range []flags.GameFlag{0, 2048, 4096, 8192, 2048 | 8192} {
			for _, present := range []byte{0, 1} {
				o.resetJournal(t)
				o.call(t, 0, 31, "PriorEntry", 4)
				old := o.players[31].Journal
				flags.ResetGame()
				flags.SetGame(gf)
				input := binary.LittleEndian.AppendUint16(nil, version)
				input = append(input, present, 0, 0, 0xde, 0xad, 0xbe, 0xef)
				wantRet, wantPos := uint32(1), int64(5)
				clears := true
				if int16(version) > 1 {
					wantRet = 0
					wantPos = 2
					clears = false
				} else if present == 0 {
					wantPos = 3
					clears = false
				} else if gf&2048 == 0 {
					wantRet = 0
					wantPos = 3
					clears = false
				}
				ret, got, pos := playerFileSection(t, "sub_41BEC0", input, uint32(uintptr(unsafe.Pointer(&o.units[2]))), 0)
				if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) {
					t.Fatal("journal gate", version, gf, present, ret, pos)
				}
				now := o.players[31].Journal
				if clears && now != nil || !clears && now != old {
					t.Fatal("journal replace/skip", version, gf, present)
				}
				rows = append(rows, map[string]any{"version": version, "flags": uint32(gf), "present": present, "return": ret, "position": pos, "cleared": now == nil})
			}
		}
	}
	spellbookCapture(t, "player-files-journal-gates", rows, "")
}
