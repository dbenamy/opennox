//go:build porttest

package server

// PortTestObjectInitSize installs a minimal type table for the protection ABI
// fixture without loading assets or changing the live server's type registry.
func (s *Server) PortTestObjectInitSize(ind uint16, size uint32) func() {
	old := s.Types.byInd
	s.Types.byInd = make([]*ObjectType, int(ind)+1)
	s.Types.byInd[ind] = &ObjectType{InitDataSize: uintptr(size)}
	return func() { s.Types.byInd = old }
}
