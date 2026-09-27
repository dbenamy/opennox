package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var (
	_ = [1]struct{}{}[516-unsafe.Sizeof(server.Waypoint{})]
)
