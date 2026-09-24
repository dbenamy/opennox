package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var (
	_ = [1]struct{}{}[516-unsafe.Sizeof(server.Waypoint{})]
	_ = [1]struct{}{}[516-unsafe.Sizeof(nox_waypoint_t{})]
)

type nox_waypoint_t = C.nox_waypoint_t

func asWaypointC(p *server.Waypoint) *nox_waypoint_t {
	return (*nox_waypoint_t)(p.C())
}

//export nox_xxx_waypointGetList_579860
func nox_xxx_waypointGetList_579860() *nox_waypoint_t {
	return asWaypointC(GetServer().S().WPs.First())
}

//export sub_579C80
func sub_579C80(a1 uint32) *nox_waypoint_t {
	return asWaypointC(GetServer().S().WPs.PendingByInd(int(a1)))
}
