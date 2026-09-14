//go:build porttest

package server

import "github.com/opennox/libs/balance"

// PortTestClientUpdateBalance overlays the real scalar balance value used by
// mana-bomb charge creation; restoration preserves the prior balance owner.
func (s *Server) PortTestClientUpdateBalance(radius float64) func() {
	old := s.Balance.file
	s.Balance.file = &balance.File{Global: balance.Config{"manabomboutradius": balance.Array{radius}}, Tags: make(map[balance.Tag]balance.Config), Parent: old}
	return func() { s.Balance.file = old }
}
