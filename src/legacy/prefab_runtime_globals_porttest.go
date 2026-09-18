//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdio.h>
extern uint32_t dword_5d4594_1599576, dword_5d4594_1599596;
extern uint32_t dword_5d4594_1599480, dword_5d4594_1599476, dword_5d4594_3835396;
extern void* dword_5d4594_1599540;
extern void* dword_5d4594_1599532;
extern void* dword_5d4594_1599556;
extern void* dword_5d4594_1599548;
extern void* dword_5d4594_1599588;
extern void* dword_5d4594_1599592;
extern uint32_t dword_5d4594_1599616, dword_5d4594_1599644;
extern uint32_t dword_5d4594_3835312, dword_5d4594_2487244;
extern FILE* nox_file_8;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestPrefabRuntimeGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"metadata":     (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599576)),
		"count":        (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599596)),
		"loaded":       (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599480)),
		"placed":       (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599476)),
		"selected":     (*uint32)(unsafe.Pointer(&C.dword_5d4594_3835396)),
		"objects":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599540)),
		"walls":        (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599532)),
		"tiles":        (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599556)),
		"waypoints":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599548)),
		"path":         (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599588)),
		"alternate":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599592)),
		"intro":        (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599616)),
		"script":       (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599644)),
		"instance":     (*uint32)(unsafe.Pointer(&C.dword_5d4594_3835312)),
		"lastWaypoint": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2487244)),
		"file":         (*uint32)(unsafe.Pointer(&C.nox_file_8)),
		"waypointKind": memmap.PtrUint32(0x973F18, 35972),
		"autoConnect":  memmap.PtrUint32(0x973F18, 35976),
	}
	// Snapshot backing blob bytes independently from extracted C globals.
	block := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1599484)), 165)
	saved := append([]byte(nil), block...)
	clear(block)
	old := map[*uint32]uint32{}
	for _, p := range words {
		old[p] = *p
		*p = 0
	}
	return words, func() {
		copy(block, saved)
		for p, v := range old {
			*p = v
		}
	}
}
