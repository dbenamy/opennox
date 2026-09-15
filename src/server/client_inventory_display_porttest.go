//go:build porttest

package server

import "github.com/opennox/libs/balance"

// PortTestInventoryDisplayBalance supplies real XP and weapon calculation data.
func (s *Server) PortTestInventoryDisplayBalance() func() {
	old := s.Balance.file
	xp := make(balance.Array, 12)
	for i := range xp {
		xp[i] = float64(i * i * 100)
	}
	s.Balance.file = &balance.File{Global: balance.Config{"xptable": xp, "boltsolodamagemin": balance.Array{9}}, Tags: make(map[balance.Tag]balance.Config), Parent: old}
	return func() { s.Balance.file = old }
}
