//go:build porttest

package legacy

import (
	"math"
)

func PortTestPrefabCall(op int, v [6]uint32) uint64 {
	switch op {
	case 1:
		return uint64(prefabScriptScan(v[0], v[1]))
	case 2:
		return uint64(prefabFindName(v[0]))
	case 3:
		return uint64(prefabMetadataAt(int32(v[0])))
	case 4:
		return uint64(*prefabGlobal(prefabCount))
	case 5:
		return uint64(prefabSetPath(v[0], false))
	case 6:
		return uint64(prefabSetPath(v[0], true))
	case 7:
		return uint64(prefabLibrary())
	case 8:
		return uint64(prefabSelect(int32(v[0])))
	case 9:
		return uint64(prefabOpen(v[0]))
	case 10:
		return uint64(prefabClose())
	case 11:
		return uint64(prefabSeek(int32(v[0])))
	case 12:
		return math.Float64bits(prefabDimension(int32(v[0]), 64))
	case 13:
		return math.Float64bits(prefabDimension(int32(v[0]), 68))
	case 14:
		return uint64(prefabLoad(int32(v[0])))
	case 15:
		return uint64(prefabInstantiate(prefabPoint(v[0])))
	case 16:
		return uint64(prefabRelativePosition(v[0], prefabPoint(v[1])))
	case 17:
		return uint64(prefabTileNew(int32(v[0]), int32(v[1]), v[2]))
	case 18:
		return uint64(prefabPlaceTiles(int32(v[0]), int32(v[1])))
	case 19:
		return uint64(prefabWallNew(byte(v[0]), byte(v[1])))
	case 20:
		return uint64(prefabWallFind(int32(v[0]), int32(v[1])))
	case 21:
		return uint64(prefabPlaceWalls(int32(v[0]), int32(v[1])))
	case 22:
		return uint64(prefabWaypointNew(v[0], math.Float32frombits(v[1]), math.Float32frombits(v[2])))
	case 23:
		return uint64(prefabPlaceWaypoints(int32(v[0]), int32(v[1])))
	case 24:
		return uint64(prefabObjectNew(v[0]))
	case 25:
		return uint64(prefabPlaceObjects(int32(v[0]), int32(v[1])))
	case 26:
		return uint64(prefabObjectHead())
	case 27:
		return uint64(prefabObjectNext(v[0]))
	case 28:
		return uint64(*prefabGlobal(prefabObjects))
	case 29:
		return uint64(prefabNodeNext(v[0]))
	case 30:
		return uint64(prefabObjectRemove(v[0]))
	case 31:
		return uint64(prefabIntroFree())
	case 32:
		return uint64(prefabIntroSection())
	case 33:
		return uint64(prefabGroupSection())
	case 34:
		return uint64(prefabWaypointSection((*[8]uint32)(mapRoomPointer(v[0]))))
	case 35:
		prefabResetWaypoint()
		return 0
	case 36:
		return uint64(prefabSetWaypointKind(byte(v[0])))
	case 37:
		return uint64(prefabCreateWaypoint(prefabPoint(v[0])))
	case 38:
		return uint64(prefabFindWaypoint(prefabPoint(v[0])))
	case 39:
		return uint64(prefabConnectWaypoint(prefabPoint(v[0]), prefabPoint(v[1])))
	}
	panic("prefab operation")
}
