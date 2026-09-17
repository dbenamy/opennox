//go:build porttest

package server

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Real registry callbacks, type-owned templates and the production object pool.
// Sizes cover each callback's actual init/collide/use/update record layout.
func (p *PortTestPaintOwners) ObjectXferTypes() {
	for _, sp := range []struct {
		name                       string
		class                      object.Class
		init, collide, use, update uintptr
	}{
		{"SpellPagePedestal", object.ClassImmobile, 0, 4, 0, 0},
		{"Readable", object.ClassReadable, 0, 0, 260, 0},
		{"Exit", object.ClassExit, 0, 88, 0, 0},
		{"Door", object.ClassDoor, 0, 0, 0, 64},
		{"Trigger", object.ClassTrigger, 0, 0, 0, 60},
		{"Hole", object.ClassHole, 0, 28, 0, 0},
		{"Transporter", object.ClassTransporter, 0, 0, 0, 20},
		{"Elevator", object.ClassElevator, 0, 0, 0, 20},
		{"ElevatorShaft", object.ClassElevatorShaft, 0, 0, 0, 12},
		{"Mover", object.ClassImmobile, 0, 0, 0, 36},
		{"Glyph", object.ClassImmobile, 36, 0, 0, 0},
		{"InvisibleLight", object.ClassLight, 0, 0, 0, 0},
		{"Sentry", object.ClassImmobile, 0, 0, 0, 12},
	} {
		s := p.S
		id := uint16(len(s.Types.byInd))
		typ := &ObjectType{s: &s.Types, ind: id, ind2: id, id: strings.ToLower("Port" + sp.name), class: sp.class, allowed: true, Mass: 1}
		typ.Xfer = xferFuncs[sp.name+"Xfer"]
		if typ.Xfer == nil {
			panic("missing real xfer: " + sp.name)
		}
		allocate := func(n uintptr) unsafe.Pointer {
			if n == 0 {
				return nil
			}
			ptr, _ := alloc.Malloc(n)
			p.typeData = append(p.typeData, ptr)
			return ptr
		}
		typ.InitDataSize = sp.init
		typ.InitData = allocate(sp.init)
		typ.CollideDataSize = sp.collide
		typ.CollideData = allocate(sp.collide)
		typ.UseDataSize = sp.use
		if sp.use != 0 {
			typ.UseData.SetPtr(allocate(sp.use))
		}
		typ.UpdateDataSize = sp.update
		typ.UpdateData = allocate(sp.update)
		if sp.name == "Trigger" {
			typ.Shape.Kind = ShapeKindBox
			typ.Shape.Box.W = 20
			typ.Shape.Box.H = 30
			typ.Shape.Box.Calc()
		}
		s.Types.byInd = append(s.Types.byInd, typ)
		s.Types.byID[typ.id] = typ
	}
}

// Admission contracts configure real type fields read by the production policy.
func (s *Server) PortTestObjectXferAdmission(class object.Class, flags object.Flags, allowed bool) {
	typ := s.Types.ByInd(1)
	typ.class = class
	typ.flags = flags
	typ.allowed = allowed
	typ.Weight = 255
}
