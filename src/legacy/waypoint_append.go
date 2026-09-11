package legacy

/*
#include "GAME4_1.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func appendWaypointLink(source, target *server.Waypoint, kind int8) bool {
	count := source.PointsCnt
	// C deliberately leaves the last of the 32 physical slots unused.
	if count >= 31 || source == target {
		return false
	}
	for i := byte(0); i < count; i++ {
		p := &source.Points[i]
		// C promotes stored uint8 and incoming signed char to int. Kinds with
		// the high bit set therefore do not match an existing stored byte.
		if p.Waypoint == target && int(p.Ind) == int(kind) {
			return false
		}
	}
	source.Points[count].Waypoint = target
	source.Points[count].Ind = byte(kind)
	source.PointsCnt = count + 1
	return true
}

//export sub_51D2C0
func sub_51D2C0(source, target C.int) C.int {
	return C.int(bool2int(appendWaypointLink(waypointFromRaw(source), waypointFromRaw(target), int8(memmap.Uint8(0x973F18, 35972)))))
}

//export sub_51D300
func sub_51D300(source, target C.int, kind C.char) C.int {
	return C.int(bool2int(appendWaypointLink(waypointFromRaw(source), waypointFromRaw(target), int8(kind))))
}
