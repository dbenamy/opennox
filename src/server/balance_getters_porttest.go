//go:build porttest

// Isolates mutable balance data without reading shipped assets.
package server

import "github.com/opennox/libs/balance"

func (s *Server) PortTestBalanceOverlay(global balance.Config, arena, solo balance.Config) func() {
	old := s.Balance.file
	s.Balance.file = &balance.File{
		Global: global,
		Tags: map[balance.Tag]balance.Config{
			balance.TagArena: arena,
			balance.TagSolo:  solo,
		},
	}
	return func() { s.Balance.file = old }
}
