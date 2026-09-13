//go:build porttest

package server

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Population uses real object-pool and per-type data allocation.
func (p *PortTestPaintOwners) PopulationTypes(exitXfer unsafe.Pointer) {
	s := p.S
	s.Types.byID["spellbook"] = s.Types.byInd[3]
	s.Types.byID["playerstart"] = s.Types.byInd[1]
	for _, name := range []string{"PopulationExit", "ExitNorthMarker", "ExitSouthMarker", "ExitEastMarker", "ExitWestMarker"} {
		id := uint16(len(s.Types.byInd))
		typ := &ObjectType{s: &s.Types, ind: id, ind2: id, id: strings.ToLower(name), class: object.ClassSimple, allowed: true, Mass: 1}
		if name == "PopulationExit" {
			typ.Xfer = exitXfer
			typ.CollideDataSize = 80
			typ.CollideData, _ = alloc.Malloc(80)
			p.typeData = append(p.typeData, typ.CollideData)
		}
		s.Types.byInd = append(s.Types.byInd, typ)
		s.Types.byID[typ.id] = typ
	}
}
