//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
void nox_xxx_mapSwitchLevel_4D12E0_tileFree(void);
extern uint32_t dword_5d4594_1556856;
extern uint32_t dword_5d4594_1568024;
extern uint32_t dword_5d4594_1548700;
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestSessionEntryScalar(op string, value int32) int32 {
	switch op {
	case "character-count":
		return int32(C.nox_client_countPlayerFiles02_4DC630())
	case "tile-clear":
		C.nox_xxx_mapSwitchLevel_4D12E0_tileFree()
		return 0
	case "crown-clear":
		C.nox_xxx_mapFindCrown_4CFC30()
		return 0
	case "map-state-set":
		return int32(C.sub_4CFDF0(C.int(value)))
	case "map-state-get":
		return int32(C.sub_4CFE00())
	case "load-state-set":
		return int32(C.nox_xxx_mapLoadOrSaveMB_4DCC70(C.int(value)))
	case "scavenger-max":
		return int32(C.nox_xxx_scavengerTreasureMax_4D1600())
	case "scavenger-reset":
		C.sub_4D15C0()
		return 0
	case "scavenger-max-reset":
		C.sub_4D1610()
		return 0
	case "latency":
		C.nox_xxx_netReportAllLatency_4D3050()
		return 0
	case "save-count":
		return int32(C.nox_client_countSaveFiles_4DC550())
	case "mode-index":
		return int32(C.sub_409A70(C.short(value)))
	default:
		panic(op)
	}
}

func PortTestSessionEntryRoster(op string, index int) int {
	switch op {
	case "init":
		C.sub_4D11A0()
	case "clear":
		C.sub_4D11D0()
	case "add":
		C.sub_4D1210(C.int(index))
	case "remove":
		C.sub_4D1250(C.int(index))
	case "contains":
		return int(C.sub_4D12A0(C.int(index)))
	default:
		panic(op)
	}
	return 0
}
func PortTestSessionEntryRosterOwner() func() {
	raw := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1548492), 16)
	old := bytes.Clone(raw)
	clear(raw)
	C.sub_4D11A0()
	return func() { C.sub_4D11D0(); copy(raw, old) }
}
func PortTestSessionEntryRosterMembers() []int {
	type node struct {
		Next, Prev *node
		Key        uint32
		Player     *server.Player
	}
	head := (*node)(memmap.PtrOff(0x5D4594, 1548492))
	out := []int{}
	for p := head.Next; p != head; p = p.Next {
		if p == nil || p.Next == nil || p.Prev == nil || p.Next.Prev != p || p.Prev.Next != p || len(out) >= 32 {
			panic("invalid game roster links")
		}
		out = append(out, int(p.Player.PlayerInd))
	}
	return out
}

func PortTestSessionEntryShadow(op string, u *server.Object) *server.Object {
	switch op {
	case "add":
		return (*server.Object)(unsafe.Pointer(C.nox_xxx_unitNewAddShadow_4DA9A0((*C.nox_object_t)(unsafe.Pointer(u)))))
	case "remove":
		return (*server.Object)(unsafe.Pointer(C.nox_xxx_action_4DA9F0((*C.nox_object_t)(unsafe.Pointer(u)))))
	case "head":
		return (*server.Object)(unsafe.Pointer(uintptr(C.dword_5d4594_1556856)))
	default:
		panic(op)
	}
}
func PortTestSessionEntryShadowOwner() func() {
	old := C.dword_5d4594_1556856
	C.dword_5d4594_1556856 = 0
	return func() { C.dword_5d4594_1556856 = old }
}

func PortTestSessionEntryLatencyOwner() func() {
	old := C.dword_5d4594_1548700
	C.dword_5d4594_1548700 = 0
	return func() { C.dword_5d4594_1548700 = old }
}
func PortTestSessionEntryLatencyCursor() int {
	if C.dword_5d4594_1548700 == 0 {
		return -1
	}
	return int((*server.Player)(unsafe.Pointer(uintptr(C.dword_5d4594_1548700))).PlayerInd)
}

func PortTestSessionEntryMapText(op, value string) (string, uint32) {
	switch op {
	case "set-path":
		p := C.nox_xxx_gameSetMapPath_409D70(internCStr(value))
		var changed uint32
		if p != nil {
			changed = 1
		}
		return "", changed
	case "filename":
		return GoString(C.nox_server_currentMapGetFilename_409B30()), 0
	case "basename":
		return GoString(C.nox_xxx_mapGetMapName_409B40()), 0
	case "set-selected":
		return "", uint32(C.sub_409B50(internCStr(value)))
	case "selected":
		return GoString(C.sub_409B80()), 0
	default:
		panic(op)
	}
}

func PortTestSessionEntryValidate(name string, nilName bool, crc uint32) int {
	var p *C.char
	if !nilName {
		p = internCStr(name)
	}
	return int(C.nox_xxx_mapValidateMB_4CF470(p, C.int(crc)))
}
func PortTestSessionEntryMetadataTable() func() {
	raw := unsafe.Slice(memmap.PtrUint8(0x587000, 55936), 24)
	old := bytes.Clone(raw)
	*memmap.PtrUint32(0x587000, 55936) = 1
	*memmap.PtrUint32(0x587000, 55948) = 1
	*memmap.PtrPtr(0x587000, 55956) = unsafe.Pointer(C.nox_xxx_parseFileInfoData_41C3B0)
	return func() { copy(raw, old) }
}

func PortTestSessionEntryTileFreeHead() *uint32 { return mapPaintGlobal(paintFreeHead) }
func PortTestSessionEntryRespawns() []*server.Object {
	var out []*server.Object
	for p := uint32(C.dword_5d4594_1568024); p != 0; p = *(*uint32)(unsafe.Pointer(uintptr(p) + 52)) {
		if len(out) > 32 {
			panic("respawn list cycle")
		}
		out = append(out, *(**server.Object)(unsafe.Pointer(uintptr(p) + 4)))
	}
	return out
}

func PortTestSessionEntryMapFlags(record unsafe.Pointer) int32 {
	return int32(C.sub_4CFFC0(C.int(uintptr(record))))
}
