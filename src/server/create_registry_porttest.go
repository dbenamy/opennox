//go:build porttest

package server

import (
	"unsafe"

	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestCreateWithRegisteredCreate clones the synthetic type at index 1,
// parses a real create procedure name, and constructs an object through the
// ordinary allocator/copy/Create callback path. An empty name uses raw as a
// temporary callback override (nil selects the no-callback case).
func (s *Server) PortTestCreateWithRegisteredCreate(name string, raw unsafe.Pointer, class object.Class, subclass object.SubClass) *Object {
	base := s.Types.ByInd(1)
	if base == nil {
		panic("PortTestCreateWithRegisteredCreate requires synthetic type 1")
	}
	oldType := s.Types.byInd[1]
	oldByID := s.Types.byID[base.id]
	oldModif := s.Modif
	oldBalance := s.Balance.file
	const gbase uintptr = 0x5D4594
	offs := make([]uintptr, 13)
	for i := range offs {
		offs[i] = 2491624 + uintptr(4*i)
	}
	oldGlobals := make([]uint32, len(offs))
	for i, off := range offs {
		oldGlobals[i] = memmap.Uint32(gbase, off)
		*memmap.PtrUint32(gbase, off) = 0
	}

	weapon, freeWeapon := alloc.New(Modifier{})
	armor, freeArmor := alloc.New(Modifier{})
	initData, freeInit := alloc.Make([]byte{}, 1724)
	updateData, freeUpdate := alloc.Make([]byte{}, 2200)
	useData, freeUse := alloc.Make([]byte{}, 128)
	health, freeHealth := alloc.New(HealthData{})
	*weapon = Modifier{TypeInd: 1, Durability52: 100}
	*armor = Modifier{TypeInd: 1, Durability52: 100}
	*health = HealthData{Cur: 40, Max: 50}
	for i := range initData {
		initData[i] = 0x11
	}
	for i := range useData {
		useData[i] = 0x22
	}

	defer func() {
		s.Types.byInd[1] = oldType
		s.Types.byID[base.id] = oldByID
		s.Modif = oldModif
		s.Balance.file = oldBalance
		for i, off := range offs {
			*memmap.PtrUint32(gbase, off) = oldGlobals[i]
		}
		freeHealth()
		freeUse()
		freeUpdate()
		freeInit()
		freeArmor()
		freeWeapon()
	}()

	overlay := &balance.File{
		Global: balance.Config{"questdurabilitymultiplier": balance.Array{1.5}, "defaultammoamountquest": balance.Array{30}, "defaultammoamount": balance.Array{20}},
		Tags:   make(map[balance.Tag]balance.Config),
		Parent: oldBalance,
	}
	s.Balance.file = overlay
	s.Modif.Dword_5d4594_251600 = weapon
	s.Modif.Dword_5d4594_251608 = armor

	typ := *base
	typ.ind, typ.ind2 = 1, 1
	typ.class, typ.subclass = class, subclass
	typ.Field9 = 0
	typ.health = health
	typ.InitData, typ.InitDataSize = unsafe.Pointer(&initData[0]), uintptr(len(initData))
	typ.UpdateData, typ.UpdateDataSize = unsafe.Pointer(&updateData[0]), uintptr(len(updateData))
	typ.UseData = UseDataPtr{Ptr: unsafe.Pointer(&useData[0])}
	typ.UseDataSize = uintptr(len(useData))
	if name != "" {
		if err := typ.parseCreate(&things.ProcFunc{Name: name}); err != nil {
			panic(err)
		}
	} else {
		typ.Create = raw
	}
	s.Types.byInd[1] = &typ
	s.Types.byID[base.id] = &typ

	return s.NewObjectByTypeInd(1)
}
