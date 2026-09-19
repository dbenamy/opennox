//go:build porttest

package server

// PortTestPlayerFileAbilities connects the actual ability map to the fixture
// server clock/player registry and restores the complete previous owner.
func (s *Server) PortTestPlayerFileAbilities() func() {
	old := s.Abils
	s.Abils = serverAbilities{s: s}
	s.Abils.Reset()
	return func() { s.Abils = old }
}
