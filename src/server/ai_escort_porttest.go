//go:build porttest

package server

import (
	"bytes"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestEscortPlayers installs real sparse player iteration with one active
// nil-unit slot. All player, object and update-data storage is C-owned.
func (s *Server) PortTestEscortPlayers() (units []Object, configure func(int), unchanged func() bool, free func()) {
	players, fp := alloc.Make([]Player{}, 32)
	units, fu := alloc.Make([]Object{}, 3)
	data, fd := alloc.Make([]PlayerUpdateData{}, 3)
	old := s.Players.list
	s.Players.list = players
	slots := []int{1, 7, 31}
	for i, slot := range slots {
		p := &players[slot]
		p.PlayerInd = byte(slot)
		units[i].ObjClass = object.ClassPlayer
		units[i].PosVec = types.Pointf{X: float32(100 + i*40), Y: 100}
		units[i].UpdateData = unsafe.Pointer(&data[i])
		data[i].Player = p
	}
	raw := [][]byte{unsafe.Slice((*byte)(unsafe.Pointer(&players[0])), len(players)*int(unsafe.Sizeof(Player{}))), unsafe.Slice((*byte)(unsafe.Pointer(&units[0])), len(units)*int(unsafe.Sizeof(Object{}))), unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*int(unsafe.Sizeof(PlayerUpdateData{})))}
	var before [][]byte
	configure = func(count int) {
		for i := range players {
			players[i].Active = 0
			players[i].PlayerUnit = nil
		}
		players[3].Active = 1
		players[3].PlayerInd = 3
		for i, slot := range slots {
			if i < count {
				players[slot].Active = 1
				players[slot].PlayerUnit = &units[i]
			}
		}
		before = before[:0]
		for _, b := range raw {
			before = append(before, bytes.Clone(b))
		}
	}
	unchanged = func() bool {
		for i, b := range raw {
			if !bytes.Equal(b, before[i]) {
				return false
			}
		}
		return true
	}
	free = func() { s.Players.list = old; fd(); fu(); fp() }
	return
}
