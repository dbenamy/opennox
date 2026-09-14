//go:build porttest

package legacy

import "math"

func populationInvokeNative(op int, v [6]uint32, high *uint32) uint32 {
	*high = 0
	switch op {
	case 0:
		return mapPopulationStart()
	case 1:
		return mapPopulationSpellID(v[0])
	case 2:
		return mapPopulationFillCachedObjects(v[0], v[1])
	case 3:
		return mapPopulationPlacePrefabInRoom(v[0], v[1], v[2], v[3])
	case 4:
		mapPopulationInventory(v[0], v[1], v[2])
		return 0
	case 5:
		return mapPopulationSpellbook(v[0], v[1])
	case 6:
		return mapPopulationItem(v[0], v[1], v[2])
	case 7:
		return mapPopulationAttach(v[0], v[1])
	case 8:
		mapPopulationRoom(v[0], v[1])
		return 0
	case 9:
		mapPopulationGroup(v[0], v[1], v[2])
		return 0
	case 10:
		return mapPopulationSpawn(v[0], v[1], v[2])
	case 11:
		return mapPopulationMonster(v[0], v[1])
	case 12:
		mapPopulationFinish(v[0])
		return 0
	case 13:
		return mapPopulationExit(v[0])
	case 14:
		return mapPopulationRoomExit(v[0], v[1])
	case 15:
		return mapPopulationWaypoint(v[0])
	case 16:
		return mapPopulationHallwayWaypoints(v[0])
	case 17:
		return mapPopulationRoomWaypoints(v[0])
	case 18:
		mapPopulationProgress(byte(v[0]))
		return 0
	case 19:
		r := math.Float64bits(mapPopulationDistanceMax())
		*high = uint32(r >> 32)
		return uint32(r)
	case 20:
		return mapPopulationFarthest()
	case 21:
		mapPopulationDistance(v[0], v[1], math.Float32frombits(v[2]))
		return 0
	case 22:
		return mapPopulationThemes(v[0])
	case 23:
		mapPopulationFindFarthest(v[0])
		return 0
	case 24:
		return mapPopulationSort()
	case 25:
		return mapPopulationSelectPrefabs(v[0])
	case 26:
		mapPopulationPrefabPositions(math.Float32frombits(v[0]))
		return 0
	case 27:
		return mapPopulationPrefabRoom(v[0], v[1])
	case 28:
		r := uint64(mapPopulationNextPosition(v[0], v[1]))
		*high = uint32(r >> 32)
		return uint32(r)
	case 29:
		return mapPopulationPrefabInfo(v[0], v[1])
	case 30:
		return mapPopulationCandidates(v[0], v[1])
	case 31:
		return mapPopulationConnectPrefabs(v[0])
	case 32:
		return mapPopulationApplyPrefabs(v[0])
	case 33:
		return mapPopulationMetadataIndex(v[0])
	case 34:
		return mapPopulationMetadataInit()
	case 35:
		mapPopulationMetadataFree()
		return 0
	case 36:
		return mapPopulationMetadataAt(v[0])
	case 37:
		return mapPopulationExitName(v[0], v[1])
	}
	panic("population operation")
}
