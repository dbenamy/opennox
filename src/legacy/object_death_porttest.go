//go:build porttest

package legacy

/*
#include "GAME5.h"
#include "server__object__die__die.h"
extern uint32_t dword_5d4594_2491704, dword_5d4594_527656;
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestDeathSpec struct {
	Language                                int
	Description                             string
	Name                                    string
	Sound, Flags, Cache, Rotation, Material uint32
	Barrel                                  bool
	DropCount, DropThreshold                uint32
	Owner                                   bool
	Markers                                 uint32
}
type PortTestDeathResult struct {
	Globals, DeathData, OwnerData, Init []uint32
	Intact                              bool
}
type portTestDeathState struct {
	init             []byte
	configureStrings func(int, string, bool)
	data             []byte
	beforeOwner      []byte
	drop             unsafe.Pointer
}

func portTestDeathEnvironment(proxy *portTestRoamOwnerServer) func() {
	freeTypes := proxy.core.PortTestDeathTypes()
	configureStrings, freeStrings := proxy.core.PortTestObjectDeathStrings()
	init, freeInit := alloc.Make([]byte{}, 36)
	data, freeData := alloc.Make([]byte{}, 132+16)
	a, freeA := alloc.CString("PortTestDebrisA")
	b, freeB := alloc.CString("PortTestDebrisB")
	st := &portTestDeathState{data: data, drop: unsafe.Pointer(a), init: init, configureStrings: configureStrings}
	proxy.callbacks.death = st
	regions := []struct {
		off  uintptr
		size int
	}{{203080, 24}, {203240, 24}, {290328, 16}, {291512, 8}}
	nameBuf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565660)), 2076)
	oldName := bytes.Clone(nameBuf)
	clear(nameBuf)
	before := make([][]byte, len(regions))
	for i, r := range regions {
		dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, r.off)), r.size)
		before[i] = bytes.Clone(dst)
		clear(dst)
	}
	*memmap.PtrPtr(0x587000, 291512) = unsafe.Pointer(a)
	*memmap.PtrPtr(0x587000, 291516) = unsafe.Pointer(b)
	*memmap.PtrUint32(0x587000, 290328) = 0x0000ff01
	*memmap.PtrUint32(0x587000, 290340) = 2
	oldCache := *memmap.PtrUint32(0x5D4594, 2491696)
	oldRot, oldBall := C.dword_5d4594_2491704, C.dword_5d4594_527656
	return func() {
		for i, r := range regions {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, r.off)), r.size), before[i])
		}
		*memmap.PtrUint32(0x5D4594, 2491696) = oldCache
		C.dword_5d4594_2491704, C.dword_5d4594_527656 = oldRot, oldBall
		copy(nameBuf, oldName)
		freeStrings()
		freeInit()
		freeB()
		freeA()
		freeData()
		freeTypes()
	}
}
func portTestDeathPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestDeathSpec) {
	st := proxy.callbacks.death
	clear(st.data)
	for i := 0; i < 8; i++ {
		st.data[i] = 0xa5
		st.data[len(st.data)-8+i] = 0x5a
	}
	copy(st.data[8:8+127], sp.Name)
	binary.LittleEndian.PutUint32(st.data[8+128:], sp.Sound)
	u.DeathData = unsafe.Pointer(&st.data[8])
	clear(st.init)
	for i := 0; i < 8; i++ {
		st.init[i] = 0xa5
		st.init[len(st.init)-8+i] = 0x5a
	}
	u.InitData = unsafe.Pointer(&st.init[8])
	proxy.life.ids[uint32(uintptr(u.InitData))] = 993
	if proxy.callbacks.spec.Op >= 57 {
		armor := proxy.callbacks.spec.Op == 57
		st.configureStrings(sp.Language, sp.Description, armor)
		if armor {
			u.ObjClass = object.ClassArmor
		} else {
			u.ObjClass = object.ClassWeapon
		}
	}
	proxy.life.ids[uint32(uintptr(unsafe.Pointer(&u.PosVec)))] = 994
	u.ObjFlags = object.Flags(sp.Flags)
	u.Material = uint16(sp.Material)
	u.TypeInd = 15
	if sp.Barrel {
		u.TypeInd = uint16(proxy.core.Types.IndByID("BarrelPortTest"))
	}
	proxy.life.ids[uint32(uintptr(u.DeathData))] = 992
	*memmap.PtrUint32(0x5D4594, 2491696) = sp.Cache
	C.dword_5d4594_2491704 = C.uint32_t(sp.Rotation)
	C.dword_5d4594_527656 = 0
	for _, off := range []uintptr{203080, 203240} {
		*memmap.PtrPtr(0x587000, off) = st.drop
		*memmap.PtrUint32(0x587000, off+4) = sp.DropCount
		*memmap.PtrUint32(0x587000, off+8) = sp.DropThreshold
	}
	owner := &proxy.life.players[0]
	d := owner.UpdateDataPlayer()
	proxy.life.ids[uint32(uintptr(unsafe.Pointer(d.Player)))] = 991
	proxy.life.ids[uint32(uintptr(owner.UpdateData))] = 990
	proxy.life.ids[uint32(uintptr(unsafe.Pointer(&owner.PosVec)))] = 995
	st.beforeOwner = bytes.Clone(unsafe.Slice((*byte)(owner.UpdateData), int(unsafe.Sizeof(*d))))
	if sp.Owner {
		u.ObjOwner = owner
		u.InvHolder = owner
	}
	for i := 0; i < 4; i++ {
		p := (*unsafe.Pointer)(unsafe.Add(owner.UpdateData, 116+4*i))
		*p = nil
		if sp.Markers&(1<<i) != 0 {
			*p = u.CObj()
		}
	}
}
func portTestDeathCall(u *server.Object, op int) uint32 {
	p := combatPtr(u)
	switch op {
	case 0:
		C.nox_xxx_dieBarrel_54DFA0(p)
	case 1:
		C.nox_xxx_dieCreateObject_54E010(p)
	case 2:
		return uint32(int32(C.nox_xxx_dieSpawnObject_54E070(p)))
	case 3:
		C.nox_xxx_dieMarker_54E460(p)
	case 4:
		C.nox_xxx_dieBoulder_54E4B0(p)
	case 5:
		return uint32(C.nox_xxx_dieGameBall_54E620(p))
	case 6:
		C.nox_xxx_dieArmor_54E170_obj_die(p)
	case 7:
		C.nox_xxx_dieWeapon_54E370_obj_die(p)
	default:
		panic("death operation")
	}
	return 0
}
func portTestDeathTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestDeathResult {
	st := proxy.callbacks.death
	r := &PortTestDeathResult{Intact: true}
	r.Globals = []uint32{*memmap.PtrUint32(0x5D4594, 2491696), uint32(C.dword_5d4594_2491704), uint32(C.dword_5d4594_527656)}
	for i := 0; i < 8; i++ {
		r.Intact = r.Intact && st.data[i] == 0xa5 && st.data[len(st.data)-8+i] == 0x5a && st.init[i] == 0xa5 && st.init[len(st.init)-8+i] == 0x5a
	}
	for i := 8; i < len(st.data)-8; i += 4 {
		r.DeathData = append(r.DeathData, binary.LittleEndian.Uint32(st.data[i:]))
	}
	for i := 8; i < len(st.init)-8; i += 4 {
		r.Init = append(r.Init, normalize(binary.LittleEndian.Uint32(st.init[i:])))
	}
	owner := &proxy.life.players[0]
	dst := unsafe.Slice((*byte)(owner.UpdateData), len(st.beforeOwner))
	for i := 0; i < len(dst); i += 4 {
		r.OwnerData = append(r.OwnerData, normalize(binary.LittleEndian.Uint32(dst[i:])))
	}
	copy(dst, st.beforeOwner)
	return r
}
