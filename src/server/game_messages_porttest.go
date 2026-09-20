//go:build porttest

package server

// PortTestGameMessageTeamTitles supplies nonempty titles to the real team lookup.
// The caller owns the actual team records independently.
func (s *Server) PortTestGameMessageTeamTitles() func() {
	oldDefs, oldStrings := s.Teams.defs, s.Teams.sm
	s.Teams.sm = s.Strings()
	s.Teams.defs = map[TeamColor]*TeamDef{1: {Title: "Red Ω"}, 2: {Title: "Blue é"}}
	return func() { s.Teams.defs = oldDefs; s.Teams.sm = oldStrings }
}
