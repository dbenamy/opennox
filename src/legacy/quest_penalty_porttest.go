//go:build porttest

package legacy

/*
#include "GAME5.h"
extern uint32_t dword_5d4594_2491676;
extern unsigned int gameex_flags;
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestPenaltyItem struct {
	Type                          uint16
	Class, Subclass, Flags, Worth uint32
	Mods                          [4]bool
	Eligible                      bool
}
type PortTestPenaltySpec struct {
	ProtectedGold          bool
	Class                  byte
	Gold                   uint32
	Items                  []PortTestPenaltyItem
	Armor                  map[uint16]uint32
	Spells                 [137]uint32
	Beasts                 [41]uint32
	SpellTable, BeastTable [][2]uint32
	Cache                  [6]uint32
}
type PortTestPenaltyResult struct {
	Protection                                        []uint32 `json:",omitempty"`
	Owner, Data, Player, Globals, Items, Initializers []uint32
	Intact                                            bool
}
type portTestPenaltyState struct {
	prepareProtection                     func(uint32)
	snapshotProtection                    func() ([]uint32, bool)
	configure                             func(map[uint16]uint32)
	items, init                           []byte
	beforeOwner, beforeData, beforePlayer []byte
	deleted                               *server.Object
	spec                                  *PortTestPenaltySpec
}

const penaltyItemStride = 772 + 16

func penaltyItem(st *portTestPenaltyState, i int) *server.Object {
	return (*server.Object)(unsafe.Pointer(&st.items[i*penaltyItemStride+8]))
}
func portTestPenaltyEnvironment(proxy *portTestRoamOwnerServer) func() {
	_, configure, freeServer := proxy.core.PortTestPenaltyEnvironment()
	items, freeItems := alloc.Make([]byte{}, 64*penaltyItemStride)
	init, freeInit := alloc.Make([]byte{}, 64*36)
	st := &portTestPenaltyState{configure: configure, items: items, init: init}
	proxy.callbacks.penalty = st
	prepareProtection, snapshotProtection, freeProtection := portTestPenaltyProtectionEnvironment()
	st.prepareProtection, st.snapshotProtection = prepareProtection, snapshotProtection
	off := []uintptr{2491680, 2491684, 2386504, 2386508, 2386512}
	saved := make([]uint32, len(off))
	for i, o := range off {
		saved[i] = *memmap.PtrUint32(0x5D4594, o)
	}
	diamond, gameex := C.dword_5d4594_2491676, C.gameex_flags
	C.gameex_flags = 0
	// Controlled eligibility records fit well inside the shipped table regions.
	regions := []struct {
		off uintptr
		n   int
	}{{207108, 36 * 12}, {207796, 16 * 12}}
	before := make([][]byte, len(regions))
	for i, r := range regions {
		before[i] = bytes.Clone(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, r.off)), r.n))
	}
	oldEligible := Nox_xxx_playerClassCanUseItem_57B3D0
	Nox_xxx_playerClassCanUseItem_57B3D0 = func(u *server.Object, cl player.Class) bool {
		id := proxy.life.ids[uint32(uintptr(u.CObj()))]
		proxy.trace = append(proxy.trace, 48, id, uint32(cl))
		return id >= 6000 && int(id-6000) < len(st.spec.Items) && st.spec.Items[id-6000].Eligible
	}
	return func() {
		freeProtection()
		Nox_xxx_playerClassCanUseItem_57B3D0 = oldEligible
		for i, r := range regions {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, r.off)), r.n), before[i])
		}
		for i, o := range off {
			*memmap.PtrUint32(0x5D4594, o) = saved[i]
		}
		C.dword_5d4594_2491676 = diamond
		C.gameex_flags = gameex
		freeInit()
		freeItems()
		freeServer()
	}
}
func portTestPenaltyPrepare(proxy *portTestRoamOwnerServer, sp *PortTestPenaltySpec) {
	st := proxy.callbacks.penalty
	st.spec = sp
	st.configure(sp.Armor)
	if len(sp.Items) > 64 || len(sp.SpellTable) >= 36 || len(sp.BeastTable) >= 16 {
		panic("penalty fixture capacity")
	}
	owner := &proxy.life.players[0]
	data := owner.UpdateDataPlayer()
	pl := data.Player
	st.beforeOwner = bytes.Clone(unsafe.Slice((*byte)(owner.CObj()), int(unsafe.Sizeof(*owner))))
	st.beforeData = bytes.Clone(unsafe.Slice((*byte)(owner.UpdateData), int(unsafe.Sizeof(*data))))
	st.beforePlayer = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl))))
	st.deleted = proxy.core.Objs.DeletedList
	proxy.core.Objs.DeletedList = nil
	pl.GoldVal = sp.Gold
	pl.SpellLvl = sp.Spells
	pl.BeastScrollLvl = sp.Beasts
	*(*byte)(unsafe.Add(unsafe.Pointer(pl), 2251)) = sp.Class
	*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 4588)) = 0
	st.prepareProtection(sp.Gold)
	if sp.ProtectedGold {
		*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 4588)) = 0x40000001
	}
	*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 4632)) = 0
	owner.InvFirstItem = nil
	proxy.life.ids[uint32(uintptr(owner.UpdateData))] = 990
	proxy.life.ids[uint32(uintptr(unsafe.Pointer(pl)))] = 991
	clear(st.items)
	clear(st.init)
	for i, item := range sp.Items {
		u := penaltyItem(st, i)
		u.TypeInd = item.Type
		u.ObjClass = object.Class(item.Class)
		u.ObjSubClass = object.SubClass(item.Subclass)
		u.ObjFlags = object.Flags(item.Flags)
		u.Worth = item.Worth
		u.NetCode = uint32(6000 + i)
		u.InvHolder = owner
		if i > 0 {
			u.Field125 = penaltyItem(st, i-1)
		}
		if i+1 < len(sp.Items) {
			u.InvNextItem = penaltyItem(st, i+1)
		}
		u.InitData = unsafe.Pointer(&st.init[i*36+8])
		for j, set := range item.Mods {
			if set {
				*(*unsafe.Pointer)(unsafe.Add(u.InitData, 4*j)) = unsafe.Pointer(proxy.callbacks.modifiers.WeaponPower1)
			}
		}
		proxy.life.ids[uint32(uintptr(u.CObj()))] = uint32(6000 + i)
		proxy.life.ids[uint32(uintptr(u.InitData))] = uint32(7000 + i)
		for _, b := range [][]byte{st.items[i*penaltyItemStride : (i+1)*penaltyItemStride], st.init[i*36 : (i+1)*36]} {
			for j := 0; j < 8; j++ {
				b[j] = 0xa5
				b[len(b)-8+j] = 0x5a
			}
		}
	}
	if len(sp.Items) > 0 {
		owner.InvFirstItem = penaltyItem(st, 0)
	}
	for i, table := range [][][2]uint32{sp.SpellTable, sp.BeastTable} {
		off, n := uintptr(207108), 36
		if i == 1 {
			off, n = 207796, 16
		}
		words := unsafe.Slice(memmap.PtrUint32(0x587000, off), n*3)
		clear(words)
		for j, entry := range table {
			words[j*3] = entry[0]
			words[j*3+1] = entry[1]
		}
	}
	C.dword_5d4594_2491676 = C.uint32_t(sp.Cache[0])
	for i, o := range []uintptr{2491680, 2491684, 2386504, 2386508, 2386512} {
		*memmap.PtrUint32(0x5D4594, o) = sp.Cache[i+1]
	}
}
func portTestPenaltyCall(proxy *portTestRoamOwnerServer, op int) uint32 {
	p := combatPtr(&proxy.life.players[0])
	switch op {
	case 0:
		C.sub_54CBD0(p)
	case 1:
		C.sub_54CC40(p)
	case 2:
		C.sub_54CD30(p)
	case 3:
		C.sub_54CE00(p)
	case 4:
		C.sub_54CEE0(p)
	case 5:
		return uint32(int32(C.sub_54CFB0(p)))
	case 6:
		C.sub_54D080(p)
	default:
		panic("penalty operation")
	}
	return 0
}
func portTestPenaltyTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestPenaltyResult {
	st := proxy.callbacks.penalty
	r := &PortTestPenaltyResult{Intact: true}
	owner := &proxy.life.players[0]
	data := owner.UpdateDataPlayer()
	pl := data.Player
	words := func(p unsafe.Pointer, n int) []uint32 {
		var out []uint32
		b := unsafe.Slice((*byte)(p), n)
		for i := 0; i+4 <= n; i += 4 {
			out = append(out, normalize(binary.LittleEndian.Uint32(b[i:])))
		}
		return out
	}
	r.Owner = words(owner.CObj(), int(unsafe.Sizeof(*owner)))
	r.Data = words(owner.UpdateData, int(unsafe.Sizeof(*data)))
	r.Player = words(unsafe.Pointer(pl), int(unsafe.Sizeof(*pl)))
	r.Globals = append(r.Globals, uint32(C.dword_5d4594_2491676))
	for _, o := range []uintptr{2491680, 2491684, 2386504, 2386508, 2386512} {
		r.Globals = append(r.Globals, *memmap.PtrUint32(0x5D4594, o))
	}
	for i := range st.spec.Items {
		u := penaltyItem(st, i)
		r.Items = append(r.Items, words(u.CObj(), 772)...)
		r.Initializers = append(r.Initializers, words(u.InitData, 20)...)
		for _, b := range [][]byte{st.items[i*penaltyItemStride : (i+1)*penaltyItemStride], st.init[i*36 : (i+1)*36]} {
			for j := 0; j < 8; j++ {
				r.Intact = r.Intact && b[j] == 0xa5 && b[len(b)-8+j] == 0x5a
			}
		}
	}
	if st.spec.ProtectedGold {
		var intact bool
		r.Protection, intact = st.snapshotProtection()
		r.Intact = r.Intact && intact && r.Protection[1] == pl.GoldVal
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(pl)), len(st.beforePlayer)), st.beforePlayer)
	copy(unsafe.Slice((*byte)(owner.UpdateData), len(st.beforeData)), st.beforeData)
	copy(unsafe.Slice((*byte)(owner.CObj()), len(st.beforeOwner)), st.beforeOwner)
	proxy.core.Objs.DeletedList = st.deleted
	return r
}
