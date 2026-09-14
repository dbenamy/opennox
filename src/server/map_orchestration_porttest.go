//go:build porttest

package server

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func (p *PortTestPaintOwners) OrchestrationTypes() {
	p.GrowthTypes()
	for len(p.walls) < 4096 {
		b, _ := alloc.Make([]byte{}, int(unsafe.Sizeof(Wall{}))+16)
		p.blocks = append(p.blocks, b)
		p.walls = append(p.walls, (*Wall)(unsafe.Pointer(&b[0])))
	}
	typ := *p.S.Types.byInd[1]
	id := uint16(len(p.S.Types.byInd))
	typ.ind = id
	typ.ind2 = id
	typ.id = "playerstart"
	p.S.Types.byInd = append(p.S.Types.byInd, &typ)
	p.S.Types.byID[typ.id] = &typ
}
