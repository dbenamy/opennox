//go:build porttest

package server

import (
	"github.com/opennox/libs/balance"
	"strings"
	"unsafe"
)

// Overlay real balance lookups without replacing their implementation.
func (s *Server) PortTestWorldCollisionBalance() (func(map[string]float64), func()) {
	old := s.Balance.file
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: old}
	s.Balance.file = overlay
	return func(values map[string]float64) {
		clear(overlay.Global)
		for k, v := range values {
			overlay.Global[strings.ToLower(k)] = balance.Array{v}
		}
	}, func() { s.Balance.file = old }
}

// Expose the actual registered callback and its declared per-object data size.
func PortTestWorldCollisionRegistry(name string) (unsafe.Pointer, uintptr) {
	entry, ok := collideFuncs[name]
	if !ok {
		panic("unknown collision callback")
	}
	return entry.Func, entry.DataSize
}

func PortTestWorldPickupRegistry(name string) PickupFuncPtr {
	p, ok := pickupFuncs[name]
	if !ok {
		panic("unknown pickup callback")
	}
	return PickupFuncPtr{Ptr: p}
}
func PortTestWorldUseRegistry(name string) (UseFuncPtr, uintptr) {
	p, ok := useFuncs[name]
	if !ok {
		panic("unknown use callback")
	}
	return UseFuncPtr{Ptr: p.Func}, p.DataSize
}

// Give the real object factory a tradable ankh using the real pickup callback.
func (s *Server) PortTestWorldAnkhType() func() {
	old := s.Types
	s.Types.byInd = append([]*ObjectType(nil), old.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(old.byID)+1)
	for k, v := range old.byID {
		s.Types.byID[k] = v
	}
	typ := *s.Types.ByInd(1)
	typ.s = &s.Types
	typ.id = "AnkhTradable"
	typ.ind = uint16(len(s.Types.byInd))
	typ.ind2 = typ.ind
	typ.Pickup = PortTestWorldPickupRegistry("AnkhTradablePickup")
	s.Types.byInd = append(s.Types.byInd, &typ)
	s.Types.byID["ankhtradable"] = &typ
	return func() { s.Types = old }
}
