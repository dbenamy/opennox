//go:build porttest

package legacy

/*
#include "GAME5.h"
*/
import "C"

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestCreationSpec struct {
	Type                              string
	TypeIndex                         uint16
	Class, Subclass, Flags            uint32
	Cache                             [13]uint32
	Disabled                          map[string]bool
	WeaponType, ArmorType, Durability uint32
	Balance                           map[string]float64
	Effects, NilHealth                bool
	InitFill, UseFill                 byte
	AutoFlag                          byte
}
type PortTestCreationResult struct {
	Globals, Init, Use, Health []uint32
	Intact                     bool
}
type portTestCreationState struct {
	ids       map[string]uint16
	configure func(map[string]bool, uint32, uint32, uint32, map[string]float64, bool)
	init, use []byte
}

func portTestCreationEnvironment(proxy *portTestRoamOwnerServer) func() {
	ids, configure, freeServer := proxy.core.PortTestCreationEnvironment()
	init, freeInit := alloc.Make([]byte{}, 256+16)
	use, freeUse := alloc.Make([]byte{}, 128+16)
	st := &portTestCreationState{ids: ids, configure: configure, init: init, use: use}
	proxy.callbacks.creation = st
	old := make([]uint32, 13)
	for i := range old {
		old[i] = *memmap.PtrUint32(0x5D4594, 2491624+uintptr(4*i))
	}
	return func() {
		for i, v := range old {
			*memmap.PtrUint32(0x5D4594, 2491624+uintptr(4*i)) = v
		}
		freeUse()
		freeInit()
		freeServer()
	}
}
func portTestCreationPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestCreationSpec) {
	st := proxy.callbacks.creation
	st.configure(sp.Disabled, sp.WeaponType, sp.ArmorType, sp.Durability, sp.Balance, sp.Effects)
	for i, v := range sp.Cache {
		*memmap.PtrUint32(0x5D4594, 2491624+uintptr(4*i)) = v
	}
	u.TypeInd = sp.TypeIndex
	if sp.Type != "" {
		u.TypeInd = st.ids[sp.Type]
	}
	u.ObjClass = object.Class(sp.Class)
	u.ObjSubClass = object.SubClass(sp.Subclass)
	u.ObjFlags = object.Flags(sp.Flags)
	for i, b := range [][]byte{st.init, st.use} {
		fill := sp.InitFill
		if i == 1 {
			fill = sp.UseFill
		}
		for j := range b {
			b[j] = fill
		}
		for j := 0; j < 8; j++ {
			b[j] = 0xa5
			b[len(b)-8+j] = 0x5a
		}
	}
	u.InitData = unsafe.Pointer(&st.init[8])
	u.UseData.Ptr = unsafe.Pointer(&st.use[8])
	if sp.NilHealth {
		u.HealthData = nil
	}
	*(*byte)(unsafe.Add(u.UpdateData, 2036)) = sp.AutoFlag
	proxy.life.ids[uint32(uintptr(u.InitData))] = 980
	proxy.life.ids[uint32(uintptr(u.UseData.Ptr))] = 981
	proxy.life.ids[uint32(uintptr(u.UpdateData))] = 982
	for i, name := range []string{"Lightning4", "Vampirism2", "Lightning3"} {
		mod := proxy.core.Modif.Nox_xxx_modifGetDescById413330(proxy.core.Modif.Nox_xxx_modifGetIdByName413290(name))
		if mod != nil {
			proxy.life.ids[uint32(uintptr(unsafe.Pointer(mod)))] = uint32(983 + i)
		}
	}
	for i, p := range []*server.Modifier{proxy.core.Modif.Nox_xxx_getProjectileClassById413250(int(sp.WeaponType)), proxy.core.Modif.Nox_xxx_equipClothFindDefByTT413270(int(sp.ArmorType))} {
		if p != nil {
			proxy.life.ids[uint32(uintptr(unsafe.Pointer(p)))] = uint32(986 + i)
		}
	}
}
func portTestCreationCall(u *server.Object, op int) uint32 {
	p := combatPtr(u)
	switch op {
	case 0:
		return uint32(int32(C.nox_xxx_monsterAutoSpells_54C0C0(asObjectC(u))))
	case 1:
		return uint32(C.nox_xxx_createWeapon_54C710(p))
	case 2:
		return uint32(C.sub_54C950(p))
	case 3:
		return uint32(C.nox_xxx_createFnObelisk_54CA10(p))
	case 4:
		C.nox_xxx_createFnAnim_54CA50(p)
	case 5:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_createTrigger_54CA60(p))))
	case 6:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_createMonsterGen_54CA90(p))))
	case 7:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_createRewardMarker_54CAC0(p))))
	case 8:
		return uint32(C.nox_xxx_dieImpEgg_54CAE0(p))
	case 9:
		C.nox_xxx_diePolyp_54CB10(p)
	case 10:
		C.nox_xxx_diePotion_54CBB0(p)
	default:
		panic("invalid creation operation")
	}
	return 0
}
func portTestCreationTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestCreationResult {
	st := proxy.callbacks.creation
	r := &PortTestCreationResult{Intact: true}
	for i := 0; i < 13; i++ {
		r.Globals = append(r.Globals, *memmap.PtrUint32(0x5D4594, 2491624+uintptr(4*i)))
	}
	for i, b := range [][]byte{st.init, st.use} {
		for j := 0; j < 8; j++ {
			if b[j] != 0xa5 || b[len(b)-8+j] != 0x5a {
				r.Intact = false
			}
		}
		var words []uint32
		for j := 8; j < len(b)-8; j += 4 {
			words = append(words, normalize(binary.LittleEndian.Uint32(b[j:])))
		}
		if i == 0 {
			r.Init = words
		} else {
			r.Use = words
		}
	}
	u := proxy.combat.actor
	if u.HealthData != nil {
		b := unsafe.Slice((*byte)(unsafe.Pointer(u.HealthData)), int(unsafe.Sizeof(server.HealthData{})))
		for i := 0; i < len(b); i += 4 {
			r.Health = append(r.Health, binary.LittleEndian.Uint32(b[i:]))
		}
	}
	// Auto-spells may change permissions; the entire actor update-data diff is
	// captured by the parent fixture. Keep its unrelated duration/health guards.
	proxy.spells.beforePermissions = append(proxy.spells.beforePermissions[:0], unsafe.Slice((*uint32)(unsafe.Add(u.UpdateData, 1492)), 136)...)
	return r
}
