//go:build porttest

package server

import (
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/strman"
)

// PortTestMapDrawableTeamMessages owns the real dependencies used when loading
// a previously unseen team from a client map record.
func (s *Server) PortTestMapDrawableTeamMessages() func() {
	oldStrings, oldPrinter := s.Teams.sm, s.Teams.pr
	s.Teams.sm = strman.New()
	s.Teams.pr = console.NewMultiPrinter()
	return func() { s.Teams.sm = oldStrings; s.Teams.pr = oldPrinter }
}
