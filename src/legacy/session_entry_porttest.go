//go:build porttest

package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestSessionEntryScalar(op string, value int32) int32 {
	switch op {
	case "character-count":
		return sessionCharacterCount()
	case "tile-clear":
		sessionClearTiles()
		return 0
	case "crown-clear":
		sessionClearCrowns()
		return 0
	case "map-state-set":
		return sessionMapStateSet(value)
	case "map-state-get":
		return sessionMapState()
	case "load-state-set":
		return sessionLoadStateSet(value)
	case "scavenger-max":
		return sessionScavengerMaximum()
	case "scavenger-reset":
		sessionScavengerReset()
		return 0
	case "scavenger-max-reset":
		sessionScavengerMaximumReset()
		return 0
	case "latency":
		sessionReportLatency()
		return 0
	case "save-count":
		return sessionSaveSlots()
	case "mode-index":
		return sessionModeIndex(int16(value))
	default:
		panic(op)
	}
}
func PortTestSessionEntryRoster(op string, index int) int {
	switch op {
	case "init":
		sessionRosterInit()
	case "clear":
		sessionRosterClear()
	case "add":
		sessionRosterAdd(int32(index))
	case "remove":
		sessionRosterRemove(int32(index))
	case "contains":
		return sessionRosterContains(int32(index))
	default:
		panic(op)
	}
	return 0
}
func PortTestSessionEntryRosterOwner() func() {
	old, initialized := sessionOperators, sessionOperatorsInitialized
	sessionOperators = nil
	sessionOperatorsInitialized = false
	sessionRosterInit()
	return func() { sessionRosterClear(); sessionOperators = old; sessionOperatorsInitialized = initialized }
}
func PortTestSessionEntryRosterMembers() []int {
	out := []int{}
	for _, pl := range sessionOperators {
		out = append(out, int(pl.PlayerInd))
	}
	return out
}
func PortTestSessionEntryShadow(op string, u *server.Object) *server.Object {
	switch op {
	case "add":
		return sessionShadowAdd(u)
	case "remove":
		return sessionShadowRemove(u)
	case "head":
		return sessionShadowHead
	default:
		panic(op)
	}
}
func PortTestSessionEntryShadowOwner() func() {
	old := sessionShadowHead
	sessionShadowHead = nil
	return func() { sessionShadowHead = old }
}
func PortTestSessionEntryLatencyOwner() func() {
	old := sessionLatencyPlayer
	sessionLatencyPlayer = nil
	return func() { sessionLatencyPlayer = old }
}
func PortTestSessionEntryLatencyCursor() int {
	if sessionLatencyPlayer == nil {
		return -1
	}
	return int(sessionLatencyPlayer.PlayerInd)
}
func PortTestSessionEntryMapText(op, value string) (string, uint32) {
	switch op {
	case "set-path":
		return "", uint32(bool2int(sessionSetMapPath(alloc.InternCString(value)) != nil))
	case "filename":
		return alloc.GoString(sessionMapFilename()), 0
	case "basename":
		return alloc.GoString(sessionMapName()), 0
	case "set-selected":
		return "", sessionSetSelectedMap(alloc.InternCString(value))
	case "selected":
		return alloc.GoString(sessionSelectedMap()), 0
	default:
		panic(op)
	}
}
func PortTestSessionEntryValidate(name string, nilName bool, crc uint32) int {
	var p *byte
	if !nilName {
		p = alloc.InternCString(name)
	}
	return int(sessionMapValidate(p, crc))
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
	for p := itemRespawnHead; p != nil; p = p.next {
		if len(out) > 32 {
			panic("respawn list cycle")
		}
		out = append(out, p.obj)
	}
	return out
}
func PortTestSessionEntryMapFlags(record unsafe.Pointer) int32 { return sessionMapFlags(record) }
