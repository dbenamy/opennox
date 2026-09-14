//go:build porttest

package server

// Use distinct types with the real door transfer/update layout already supplied
// by the painting owner. Shared immutable template bytes retain their one owner.
func (p *PortTestPaintOwners) GrowthTypes() {
	for _, name := range []string{"archeddoor", "archedhalfdoor"} {
		typ := *p.S.Types.byInd[2]
		id := uint16(len(p.S.Types.byInd))
		typ.ind = id
		typ.ind2 = id
		typ.id = name
		p.S.Types.byInd = append(p.S.Types.byInd, &typ)
		p.S.Types.byID[name] = &typ
	}
}
