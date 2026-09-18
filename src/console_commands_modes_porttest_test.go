//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestConsoleCommandsModeAndReentry(t *testing.T) {
	o := newConsoleCommandOwner(t)
	table := serverConfigOwnBytes(t, 0x587000, 94400, 32)
	clear(table)
	for i, entry := range []struct {
		name string
		bits uint16
	}{{"arena", 0x100}, {"quest", 0x1000}, {"flagball", 0x40}} {
		name, free := alloc.CString16(entry.name)
		t.Cleanup(free)
		binary.LittleEndian.PutUint32(table[8*i:], uint32(uintptr(unsafe.Pointer(name))))
		binary.LittleEndian.PutUint16(table[8*i+4:], entry.bits)
	}
	slot := serverConfigOwnBytes(t, 0x5D4594, 371438, 58)
	type row struct {
		Name          string
		Before, After uint16
	}
	var rows []row
	for _, before := range []uint16{0, 0xffff, 0x5a5a} {
		for _, tc := range []struct {
			name  string
			bits  uint16
			found bool
		}{{"arena", 0x100, true}, {"QUEST", 0x1000, true}, {"flagball", 0x40, true}, {"missing", 0, false}} {
			binary.LittleEndian.PutUint16(slot[52:], before)
			if !o.call(t, "set mode", false, tc.name) {
				t.Fatal("mode result")
			}
			got := binary.LittleEndian.Uint16(slot[52:])
			want := before
			if tc.found {
				want = before&0xe80f | tc.bits
			}
			if got != want {
				t.Fatal(tc, got, want)
			}
			rows = append(rows, row{tc.name, before, got})
		}
	}
	reentry := serverConfigOwnBytes(t, 0x5D4594, 3508, 4)
	for _, flags := range []uint32{0, 8192} {
		binary.LittleEndian.PutUint32(reentry, 17)
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		if !o.call(t, "cheat re-enter", false) {
			t.Fatal("reentry result")
		}
		restore()
		want := uint32(1)
		if flags == 8192 {
			want = 17
		}
		if binary.LittleEndian.Uint32(reentry) != want {
			t.Fatal("reentry state")
		}
	}
	spellbookCapture(t, "console-commands-modes", rows, "f0ad09bb6a241c8d20368f84ef7827f32b6ab48db9c25d69f1b5f88e4fee8516")
}
func TestConsoleCommandsCreateTeams(t *testing.T) {
	o := newConsoleCommandOwner(t)
	t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
	o.c.srv.Teams.Reset()
	o.c.srv.Teams.ActiveCnt = 0
	if !o.call(t, "set team", false) || o.c.srv.Teams.Count() != 2 || !noxflags.HasGamePlay(4) {
		t.Fatal("team creation", fmt.Sprint(o.c.srv.Teams.Count()))
	}
}
