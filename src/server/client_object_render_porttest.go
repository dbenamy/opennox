//go:build porttest

package server

import "github.com/opennox/opennox/v1/legacy/common/alloc"

// Own the actual player list used by active/ID iteration without replacing the
// server that owns drawables, teams, NPCs and RNG in the rendering fixture.
func (s *Server) PortTestObjectRenderPlayers() ([]Player, func()) {
	old := s.Players.list
	players, free := alloc.Make([]Player{}, 32)
	s.Players.list = players
	return players, func() { s.Players.list = old; free() }
}
