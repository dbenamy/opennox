//go:build porttest

package server

import "github.com/opennox/opennox/v1/legacy/common/alloc"

func (s *Server) PortTestScreenNPCs(capacity int) func() {
	old := s.NPCs.arr
	arr, free := alloc.Make([]NPC{}, capacity)
	s.NPCs.arr = arr
	return func() { s.NPCs.arr = old; free() }
}
