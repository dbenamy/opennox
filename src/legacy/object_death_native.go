package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "server__object__objutil.h"
extern uint32_t dword_5d4594_2491704;
static void deathLineMessage(int u, wchar2_t* format, wchar2_t* name) {
 nox_xxx_netSendLineMessage_4D9EB0(u,format,name);
}
*/
import "C"
import (
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

func objectDeathBarrel(u *server.Object) {
	core := GetServer().S()
	cache := memmap.PtrUint32(0x5D4594, 2491696)
	if *cache == 0 {
		*cache = uint32(core.Types.IndByID("BarrelBreaking"))
	}
	if t := core.NewObjectByTypeInd(int(*cache)); t != nil {
		GetServer().CreateObjectAt(t, nil, u.PosVec)
	}
	core.Audio.EventObj(286, u, 0, 0)
	C.nox_xxx_spawnSomeBarrel_4E7470(C.int(uintptr(u.CObj())), C.int(uintptr(unsafe.Pointer(&u.PosVec))))
	GetServer().DelayedDelete(u)
}
func objectDeathCreate(u *server.Object, spawn bool) int16 {
	core := GetServer().S()
	data := u.DeathData
	if t := core.NewObjectByTypeID(alloc.GoString((*byte)(data))); t != nil {
		GetServer().CreateObjectAt(t, nil, u.PosVec)
	}
	if id := *(*uint32)(unsafe.Add(data, 128)); id != 0 {
		core.Audio.EventObj(sound.ID(id), u, 0, 0)
	}
	if spawn {
		u.ObjFlags |= 0x8000
		return int16(u.ObjFlags)
	}
	GetServer().DelayedDelete(u)
	return 0
}
func objectDeathMarker(u *server.Object) {
	if owner := u.FindOwnerChainPlayer(); owner != nil {
		// The retained owner-chain service can return a non-player terminal owner;
		// preserve the C callback's raw four-slot layout access.
		for i := 0; i < 4; i++ {
			p := (**server.Object)(unsafe.Add(owner.UpdateData, 116+4*i))
			if *p == u {
				*p = nil
				break
			}
		}
	}
	monsterPointFX(u, 138)
	GetServer().DelayedDelete(u)
}
func objectDeathBoulder(u *server.Object) {
	core := GetServer().S()
	core.Audio.EventObj(757, u, 0, 0)
	monsterPointFX(u, 138)
	n := core.Rand.Logic.IntClamp(20, 30)
	index := uint32(C.dword_5d4594_2491704)
	for i := 0; i < n; i++ {
		name := *memmap.PtrPtr(0x587000, 291512+4*uintptr(index))
		t := core.NewObjectByTypeID(alloc.GoString((*byte)(name)))
		if t == nil {
			return
		} // Original early failure leaves the source alive.
		monsterDebrisPlace(u, t, 30)
		Nox_xxx_unitRaise_4E46F0(t, float32(core.Rand.Logic.FloatClamp(10, 70)))
		t.Field27 = float32(core.Rand.Logic.FloatClamp(-2, 0))
		t.ObjFlags |= 0x800000
		t.Field29 = math.Float32bits(float32(*(*byte)(memmap.PtrOff(0x587000, 290328+uintptr(C.dword_5d4594_2491704)))))
		GetServer().ApplyForce(t, u.PosVec, float64(float32(core.Rand.Logic.FloatClamp(5, 20))))
		monsterDebrisDecay(t, 45, 75)
		index = (uint32(C.dword_5d4594_2491704) + 1) % *memmap.PtrUint32(0x587000, 290340)
		C.dword_5d4594_2491704 = C.uint32_t(index)
	}
	GetServer().DelayedDelete(u)
}
func objectDeathArmor(u *server.Object) {
	core := GetServer().S()
	plural := false
	if lang := core.Strings().Lang(); lang == 0 || lang == 1 {
		if def := core.Modif.Nox_xxx_equipClothFindDefByTT413270(int(u.TypeInd)); def != nil {
			name := alloc.GoString16(def.Desc8)
			plural = len(name) != 0 && (name[len(name)-1] == 's' || name[len(name)-1] == 'S')
		}
	}
	holder := u.InvHolder
	pos := &u.PosVec
	if holder != nil {
		pos = &holder.PosVec
	}
	key := "ArmorDieGeneric"
	id := sound.ID(uint32(uintptr(unsafe.Pointer(pos))))
	switch {
	case u.Material&16 != 0:
		key = "ArmorDieMetal"
		id = 806
	case u.Material&8 != 0:
		key = "ArmorDieWood"
		id = 812
	case u.Material&4 != 0:
		key = "ArmorDieHide"
		id = 809
	case u.Material&2 != 0:
		key = "ArmorDieCloth"
		id = 815
	}
	if plural && key != "ArmorDieGeneric" {
		key += "Plural"
	}
	format := internWStr(core.Strings().GetStringInFile(strman.ID(key), "Die.c"))
	name := C.nox_xxx_itemGetName_4E77E0_obj_util(C.int(uintptr(u.CObj())))
	C.deathLineMessage(C.int(uintptr(unsafe.Pointer(holder))), format, name)
	core.Audio.EventPos(id, *pos, 0, 0)
	GetServer().DelayedDelete(u)
}
func objectDeathWeapon(u *server.Object) {
	core := GetServer().S()
	holder := u.InvHolder
	pos := &u.PosVec
	if holder != nil {
		pos = &holder.PosVec
	}
	key := "WeaponDieGeneric"
	var id sound.ID
	var name *C.wchar2_t
	switch {
	case u.Material&16 != 0:
		name = C.nox_xxx_itemGetName_4E77E0_obj_util(C.int(uintptr(u.CObj())))
		key = "WeaponDieMetal"
		id = 818
	case u.Material&8 != 0:
		name = C.nox_xxx_itemGetName_4E77E0_obj_util(C.int(uintptr(u.CObj())))
		key = "WeaponDieWood"
		id = 819
	default:
		name = internWStr(core.Armor.Sub_415B60(u))
	}
	format := internWStr(core.Strings().GetStringInFile(strman.ID(key), "Die.c"))
	C.deathLineMessage(C.int(uintptr(unsafe.Pointer(holder))), format, name)
	if id != 0 {
		core.Audio.EventPos(id, *pos, 0, 0)
	}
	GetServer().DelayedDelete(u)
}

//export nox_xxx_dieBarrel_54DFA0
func nox_xxx_dieBarrel_54DFA0(a C.int) { objectDeathBarrel(objectFromInt(a)) }

//export nox_xxx_dieCreateObject_54E010
func nox_xxx_dieCreateObject_54E010(a C.int) { objectDeathCreate(objectFromInt(a), false) }

//export nox_xxx_dieSpawnObject_54E070
func nox_xxx_dieSpawnObject_54E070(a C.int) C.short {
	return C.short(objectDeathCreate(objectFromInt(a), true))
}

//export nox_xxx_dieMarker_54E460
func nox_xxx_dieMarker_54E460(a C.int) { objectDeathMarker(objectFromInt(a)) }

//export nox_xxx_dieBoulder_54E4B0
func nox_xxx_dieBoulder_54E4B0(a C.int) { objectDeathBoulder(objectFromInt(a)) }

//export nox_xxx_dieGameBall_54E620
func nox_xxx_dieGameBall_54E620(a C.int) C.int { return C.int(objectiveBallReset(objectFromInt(a))) }

//export nox_xxx_dieArmor_54E170_obj_die
func nox_xxx_dieArmor_54E170_obj_die(a C.int) { objectDeathArmor(objectFromInt(a)) }

//export nox_xxx_dieWeapon_54E370_obj_die
func nox_xxx_dieWeapon_54E370_obj_die(a C.int) { objectDeathWeapon(objectFromInt(a)) }
