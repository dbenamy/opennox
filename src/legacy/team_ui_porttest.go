//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME2.h"
#include "GAME1_2.h"
#include "client__gui__gui_ctf.h"
#include "client__gui__servopts__playrlst.h"
extern uint32_t dword_5d4594_1045604;
extern uint32_t dword_5d4594_527656;
extern uint32_t dword_5d4594_1045636;
extern uint32_t dword_5d4594_1045640;
extern uint32_t dword_5d4594_1045684;
extern uint32_t dword_5d4594_1045688;
extern uint32_t dword_5d4594_1045692;
int sub_456DF0(int a1);
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// The thin accessor calls original production C; it contains no UI algorithms.
func PortTestTeamUI(op string, w *gui.Window, tm *server.Team, code, value int, text string) int {
	switch op {
	case "ctf-construct":
		return int(C.sub_455C30())
	case "ctf-show":
		return int(C.sub_455A00(C.int(value)))
	case "ctf-open":
		return int(C.sub_455A50(C.char(value)))
	case "ctf-hide":
		return int(C.sub_455C10())
	case "ctf-select":
		return int(C.sub_455E70(C.uchar(code)))
	case "ctf-tooltip":
		C.sub_455D80(C.uchar(code), C.char(value))
	case "ctf-destroy":
		return int(C.sub_455EE0())
	case "ball-construct":
		return int(C.sub_456070())
	case "ball-show":
		return int(C.sub_455F10(C.int(value)))
	case "ball-open":
		return int(C.sub_455F60())
	case "ball-hide":
		return int(C.sub_456050())
	case "ball-destroy":
		return int(C.sub_456240())
	case "players-construct":
		return int(C.nox_xxx_guiServerPlayersLoad_456270(C.int(code)))
	case "players-refresh":
		return int(C.sub_456500())
	case "players-refresh-if-open":
		return int(C.sub_4573A0())
	case "players-destroy":
		return int(uintptr(unsafe.Pointer(C.sub_456D60(C.int(value)))))
	case "player-add":
		return int(C.sub_457140(C.int(code), internWStr(text)))
	case "player-remove":
		return int(C.sub_456DF0(C.int(code)))
	case "player-find":
		return int(C.sub_456E40(C.int(code), C.int(value)))
	case "player-team":
		return int(C.sub_4571A0(C.int(code), C.int(value)))
	case "team-add":
		return int(uintptr(unsafe.Pointer(C.sub_457230((*C.wchar2_t)(tm.C())))))
	case "team-color":
		return int(C.sub_457120(C.int(uintptr(tm.C()))))
	case "team-rename":
		return int(C.sub_457010(C.int(uintptr(tm.C())), internWStr(text)))
	case "team-remove":
		return int(C.sub_456EA0(internWStr(text)))
	case "team-find":
		return int(C.sub_456F10(internWStr(text), C.int(value)))
	case "team-clear":
		return int(C.sub_456FA0())
	case "team-join":
		C.sub_456BB0(C.int(uintptr(tm.C())))
	case "requests-reset":
		C.sub_4573B0()
	case "map-ctf":
		return int(C.nox_xxx_mapInfoSetCapflag_417EA0())
	case "map-ball":
		return int(C.nox_xxx_mapInfoSetFlagball_417F30())
	default:
		panic(op)
	}
	return 0
}
func PortTestTeamUIDraw(op string, w *gui.Window) int {
	switch op {
	case "ctf":
		return int(C.sub_455CD0((*C.uint8_t)(w.C()), (*C.uint32_t)(w.DrawData().C())))
	case "ball":
		return int(C.sub_4560D0(C.int(uintptr(w.C())), C.int(uintptr(w.DrawData().C()))))
	case "players":
		return int(C.sub_456640(C.int(uintptr(w.C())), C.int(uintptr(w.DrawData().C()))))
	default:
		panic(op)
	}
}
func PortTestTeamUIEvent(w, child *gui.Window, event, value int) int {
	return int(C.sub_4567C0(C.int(uintptr(w.C())), C.int(event), (*C.int)(child.C()), C.int(value)))
}
func PortTestTeamUISelectedName(index int) string {
	var b [128]uint16
	C.sub_456D00(C.int(index), (*C.wchar2_t)(unsafe.Pointer(&b[0])))
	return alloc.GoString16S(b[:])
}

// Extracted window pointers are independent of the old backing-blob slots.
func PortTestTeamUIWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"ball-start-type": (*uint32)(unsafe.Pointer(&C.dword_5d4594_527656)),
		"ctf":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045604)),
		"ball":            (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045636)),
		"ball-visible":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045640)),
		"players":         (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045684)),
		"join":            (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045688)),
		"rename":          (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045692)),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	region := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1045608)), 92)
	saved := append([]byte(nil), region...)
	clear(region)
	C.nox_common_list_clear_425760((*C.nox_list_item_t)(memmap.PtrOff(0x5D4594, 1045652)))
	C.nox_common_list_clear_425760((*C.nox_list_item_t)(memmap.PtrOff(0x5D4594, 1045668)))
	return words, func() {
		C.sub_456D60(1)
		C.sub_455EE0()
		C.sub_456240()
		copy(region, saved)
		for k, p := range words {
			*p = old[k]
		}
	}
}

// Snapshot the real intrusive row storage, without reproducing list UI behavior.
type PortTestTeamUIRow struct {
	Name            string
	Code, TeamColor uint32
	Palette         byte
}

func PortTestTeamUIRows(teams bool) []PortTestTeamUIRow {
	off := uintptr(1045652)
	if teams {
		off = 1045668
	}
	var rows []PortTestTeamUIRow
	for n := listNext((*legacyListNode)(memmap.PtrOff(0x5D4594, off))); n != nil; n = listNext(n) {
		if len(rows) >= 1024 {
			panic("team UI row cycle")
		}
		p := unsafe.Pointer(n)
		rows = append(rows, PortTestTeamUIRow{Name: alloc.GoString16((*uint16)(unsafe.Add(p, 12))), Code: *(*uint32)(unsafe.Add(p, 60)), Palette: *(*byte)(unsafe.Add(p, 64)), TeamColor: *(*uint32)(unsafe.Add(p, 68))})
	}
	return rows
}
