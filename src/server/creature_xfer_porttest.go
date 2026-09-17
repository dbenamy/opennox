//go:build porttest

package server

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Extend the existing transfer fixture with real monster/NPC callbacks and
// correctly sized, independently owned type buffers.
func (s *Server) PortTestCreatureXferTypes() func() {
	oldInd := s.Types.byInd
	oldID := make(map[string]*ObjectType, len(s.Types.byID))
	for k, v := range s.Types.byID {
		oldID[k] = v
	}
	var owned []unsafe.Pointer
	for _, name := range []string{"Monster", "NPC"} {
		id := uint16(len(s.Types.byInd))
		typ := &ObjectType{s: &s.Types, ind: id, ind2: id, id: strings.ToLower("PortCreature" + name), class: object.ClassMonster, allowed: true, Mass: 1}
		if name == "NPC" {
			typ.subclass = object.SubClass(0x20)
		}
		typ.Xfer = xferFuncs[name+"Xfer"]
		if typ.Xfer == nil {
			panic("missing real creature xfer: " + name)
		}
		typ.InitDataSize = 1724
		typ.InitData, _ = alloc.Malloc(typ.InitDataSize)
		owned = append(owned, typ.InitData)
		typ.UpdateDataSize = unsafe.Sizeof(MonsterUpdateData{})
		typ.UpdateData, _ = alloc.Malloc(typ.UpdateDataSize)
		owned = append(owned, typ.UpdateData)
		s.Types.byInd = append(s.Types.byInd, typ)
		s.Types.byID[typ.id] = typ
	}
	return func() {
		s.Types.byInd = oldInd
		s.Types.byID = oldID
		for _, p := range owned {
			alloc.FreePtr(p)
		}
	}
}
