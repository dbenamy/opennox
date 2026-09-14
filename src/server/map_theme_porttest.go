//go:build porttest

package server

import (
	"bytes"
	"github.com/opennox/libs/player"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// Install sparse real player slots and check every byte remains unchanged.
func (s *Server) PortTestThemePlayers(classes []byte) (unchanged func() bool, restore func()) {
	if len(classes) > 32 {
		panic("theme fixture player count")
	}
	list, free := alloc.Make([]Player{}, 32)
	old := s.Players.list
	s.Players.list = list
	for i, class := range classes {
		slot := (i*7 + 1) % 32
		p := &list[slot]
		p.Active = 1
		p.PlayerInd = byte(slot)
		p.Info().SetPlayerClass(player.Class(class))
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&list[0])), int(unsafe.Sizeof(Player{}))*len(list))
	before := bytes.Clone(raw)
	return func() bool { return bytes.Equal(raw, before) }, func() { s.Players.list = old; free() }
}
