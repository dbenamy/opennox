//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestItemRespawnGlobals() (*uint32, func()) {
	old := itemRespawnCrown
	return &itemRespawnEnabled, func() { itemRespawnCrown = old }
}
func PortTestItemRespawn(op string, u *server.Object) uintptr {
	switch op {
	case "add":
		return itemRespawnAdd(u)
	case "remove":
		itemRespawnRemove(u)
	case "reset":
		itemRespawnReset()
	case "tick":
		itemRespawnTick()
	}
	return 0
}
func PortTestItemRespawnRecords() []*[15]uint32 {
	var out []*[15]uint32
	for p := itemRespawnHead; p != nil; p = p.next {
		if len(out) > 1024 {
			panic("respawn list cycle")
		}
		out = append(out, (*[15]uint32)(unsafe.Pointer(p)))
	}
	return out
}
func PortTestItemRespawnSameTeam(a, b *server.Object) bool        { return itemOwnerSameTeam(a, b) }
func PortTestItemRespawnDropCrown(u *server.Object, stamp uint32) { itemDropCrowns(u, stamp) }
