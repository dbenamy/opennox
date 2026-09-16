//go:build porttest

package server

// PortTestBookSpellOwner connects the real spell metadata methods to the
// lightweight GUI fixture server, which does not run full server startup.
func (s *Server) PortTestBookSpellOwner() func() {
	old := s.Spells.s
	s.Spells.s = s
	return func() { s.Spells.s = old }
}
