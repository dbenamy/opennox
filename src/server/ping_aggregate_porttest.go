//go:build porttest

package server

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestPingServer creates only the C-backed Player slice required by the
// legacy first/next player-info exports. Slots are both list positions and
// PlayerInd values, matching production iteration.
func PortTestPingServer(slots []int) (s *Server, raw []byte, free func()) {
	players, free := alloc.Make([]Player{}, 32)
	s = new(Server)
	s.Players.list = players
	for _, slot := range slots {
		if slot < 0 || slot >= len(players) {
			continue
		}
		players[slot].PlayerInd = byte(slot)
		players[slot].Active = 1
	}
	raw = unsafe.Slice((*byte)(unsafe.Pointer(&players[0])), len(players)*int(unsafe.Sizeof(Player{})))
	return s, raw, free
}
