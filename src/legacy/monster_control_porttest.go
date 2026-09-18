//go:build porttest

package legacy

/*
#include "GAME4_1.h"
extern uint32_t dword_587000_237036;
extern unsigned int nox_gameDisableMapDraw_5d4594_2650672;
int sub_516570();
extern void* nox_alloc_pendingOwn_2386916;
extern uint32_t dword_5d4594_2386920;
void nox_xxx_monsterActionMelee_515A30(nox_object_t*,float2*);
void nox_xxx_monsterMissileAttack_515B80(nox_object_t*,float2*);
void nox_server_scriptFleeFrom_515F70(nox_object_t*,void*);
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// The adapter invokes original C control commands; target pointers belong to real
// engine objects and the argument block contains scalar or owned-object words.
func PortTestMonsterControl(op string, u, target *server.Object, args *[8]uint32, arg int32) uint32 {
	p := (*C.nox_object_t)(u.CObj())
	q := (*C.nox_object_t)(target.CObj())
	raw := C.int(uintptr(u.CObj()))
	switch op {
	case "death":
		return uint32(C.nox_xxx_monsterCallDieFn_50A3D0((*C.uint32_t)(u.CObj())))
	case "revive":
		return uint32(C.sub_516D00(p))
	case "chapter":
		return uint32(C.sub_516570())
	case "animation":
		C.nox_xxx_updateNPCAnimData_50A850(p) // Its sole production caller ignores the address-derived return byte.
		return 0
	case "head":
		return uint32(C.nox_xxx_mobActionGet_50A020(raw))
	case "ensure":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_monsterAction_50A360(raw, C.int(arg)))))
	case "refresh":
		return uint32(C.nox_xxx_mobAction_50A910(p))
	case "look":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_monsterLookAt_5125A0(p, C.int(arg)))))
	case "walk":
		C.nox_xxx_monsterWalkTo_514110(p, *(*C.float)(unsafe.Pointer(&args[0])), *(*C.float)(unsafe.Pointer(&args[1])))
	case "patrol":
		C.nox_xxx_monsterGoPatrol_515680(p, unsafe.Pointer(args))
	case "hunt":
		C.nox_xxx_unitHunt_5157A0(p)
	case "idle":
		C.nox_xxx_unitIdle_515820(p)
	case "follow":
		C.nox_xxx_unitSetFollow_5158C0(p, q)
	case "melee":
		C.nox_xxx_monsterActionMelee_515A30(p, (*C.float2)(unsafe.Pointer(args)))
	case "missile":
		C.nox_xxx_monsterMissileAttack_515B80(p, (*C.float2)(unsafe.Pointer(args)))
	case "control-byte":
		return uint32(C.sub_515C80(raw, (*C.uint8_t)(unsafe.Pointer(args))))
	case "fight":
		C.nox_xxx_mobSetFightTarg_515D30(p, q)
	case "flee":
		C.nox_server_scriptFleeFrom_515F70(p, unsafe.Pointer(args))
	case "wait":
		C.sub_516090(p, C.uint32_t(uint32(arg)))
	default:
		panic(op)
	}
	return 0
}

func PortTestMonsterCacheInit() (*uint32, func()) {
	p := (*uint32)(unsafe.Pointer(&C.dword_587000_237036))
	old := *p
	*p = 1
	return p, func() { *p = old }
}
func PortTestMonsterCache(op string, u *server.Object, id int32) uint32 {
	switch op {
	case "reset":
		return uint32(C.sub_511D20())
	case "prepare":
		return uint32(C.nox_xxx_scriptPrepareFoundUnit_511D70((*C.nox_object_t)(u.CObj())))
	case "find":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_script_511C50(C.int(id)))))
	case "remove":
		return uint32(C.sub_511DE0((*C.nox_object_t)(u.CObj())))
	case "clear":
		return uint32(C.sub_511E20())
	default:
		panic(op)
	}
}

func PortTestMonsterPendingOwner() (*uint32, func()) {
	oldPool, oldHead := C.nox_alloc_pendingOwn_2386916, C.dword_5d4594_2386920
	C.nox_alloc_pendingOwn_2386916 = nil
	C.dword_5d4594_2386920 = 0
	return (*uint32)(unsafe.Pointer(&C.dword_5d4594_2386920)), func() {
		if C.nox_alloc_pendingOwn_2386916 != nil {
			C.sub_516F10()
		}
		C.nox_alloc_pendingOwn_2386916 = oldPool
		C.dword_5d4594_2386920 = oldHead
	}
}
func PortTestMonsterPending(op string, a, b int32) uint32 {
	switch op {
	case "init":
		return uint32(C.nox_xxx_allocPendingOwnsArray_516EE0())
	case "free":
		return uint32(C.sub_516F10())
	case "clear":
		C.sub_516F30()
	case "add":
		return uint32(uintptr(unsafe.Pointer(C.sub_516F90(C.int(a), C.int(b)))))
	case "resolve":
		C.sub_516FC0()
	default:
		panic(op)
	}
	return 0
}

func PortTestMonsterChapterOwner() (*uint32, func()) {
	p := (*uint32)(unsafe.Pointer(&C.nox_gameDisableMapDraw_5d4594_2650672))
	old := *p
	return p, func() { *p = old }
}
