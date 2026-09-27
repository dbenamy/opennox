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
