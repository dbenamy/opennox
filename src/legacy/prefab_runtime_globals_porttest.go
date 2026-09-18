//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestPrefabRuntimeGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"metadata":     prefabGlobal(prefabMetadata),
		"count":        prefabGlobal(prefabCount),
		"loaded":       prefabGlobal(prefabLoaded),
		"placed":       prefabGlobal(prefabPlaced),
		"selected":     prefabGlobal(prefabSelected),
		"objects":      prefabGlobal(prefabObjects),
		"walls":        prefabGlobal(prefabWalls),
		"tiles":        prefabGlobal(prefabTiles),
		"waypoints":    prefabGlobal(prefabWaypoints),
		"path":         prefabGlobal(prefabPath),
		"alternate":    prefabGlobal(prefabAlternate),
		"intro":        prefabGlobal(prefabIntro),
		"script":       prefabGlobal(prefabScript),
		"instance":     prefabGlobal(prefabInstance),
		"lastWaypoint": prefabGlobal(prefabLastWaypoint),
		"file":         prefabGlobal(prefabFile),
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
