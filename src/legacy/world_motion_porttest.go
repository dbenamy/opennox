//go:build porttest

package legacy

/*
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_2488620;
extern uint32_t dword_5d4594_2386576, dword_5d4594_2386564;
static void motionSentryCandidate(nox_object_t* target,nox_object_t* sentry,float4* ray) {uintptr_t context[2]={(uintptr_t)sentry,(uintptr_t)ray};nox_xxx_sentry_511020((int)target,(int)context);}
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestWorldMotionList(op string, u *server.Object, delay int32) uint32 {
	switch op {
	case "decay-set":
		return uint32(C.nox_xxx_unitSetDecayTime_511660(asObjectC(u), C.int(delay)))
	case "decay-remove":
		return uint32(C.nox_xxx_decay_5116F0(asObjectC(u)))
	case "decay-tick":
		C.nox_xxx_decay_511750()
	case "decay-clear":
		return uint32(C.nox_xxx_decayDestroy_5117B0())
	case "sentry-update":
		return uint32(C.nox_xxx_updateSentryGlobe_510E60(C.int(uintptr(u.CObj()))))
	case "sentry-remove":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_sentryUpdateList_510FD0((*C.uint32_t)(u.CObj())))))
	case "sentry-clear":
		C.sub_510E50()
	default:
		panic(op)
	}
	return 0
}
func PortTestWorldMotionListGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"decay": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2386576)), "sentry": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2386564))}
	old := map[string]uint32{}
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

func PortTestWorldMotionPhysics(op string, u *server.Object) int8 {
	switch op {
	case "projectile":
		C.sub_51B810(asObjectC(u))
	case "fall":
		C.nox_xxx_updateFallLogic_51B870(asObjectC(u))
	case "activate-nonsimple":
		return int8(C.sub_5117F0(asObjectC(u)))
	case "activate":
		return int8(C.sub_51B860(C.int(uintptr(u.CObj()))))
	case "remove-nonsimple":
		C.nox_xxx_unit_511810(asObjectC(u))
	default:
		panic(op)
	}
	return 0
}

func PortTestWorldMotionVelocity(step float32) int32 {
	return int32(C.nox_xxx_updateObjectsVelocity_5118A0(C.float(step)))
}
func PortTestWorldMotionVelocityGlobals() (*[10]uint32, func()) {
	p := (*[10]uint32)(memmap.PtrOff(0x5D4594, 2386580))
	old := *p
	*p = [10]uint32{}
	return p, func() { *p = old }
}

func PortTestWorldMotionSentryCandidate(target, sentry *server.Object, ray *[4]float32) {
	C.motionSentryCandidate(asObjectC(target), asObjectC(sentry), (*C.float4)(unsafe.Pointer(ray)))
}
func PortTestWorldMotionSentryReport(player int32) { C.sub_511100(C.int(player)) }
func PortTestWorldMotionSentryPacket(player int32, u *server.Object) int32 {
	return int32(C.sub_511250(C.int(player), (*C.float)(u.CObj())))
}

func PortTestWorldMotionRoundScratch() (*[4]uint32, func()) {
	p := (*[4]uint32)(memmap.PtrOff(0x5D4594, 527668))
	old := *p
	*p = [4]uint32{}
	return p, func() { *p = old }
}

func PortTestWorldMotionTrace(u *server.Object, target *uint32, normal *types.Pointf) int8 {
	return int8(C.nox_xxx_projectileTraceHit_537850(C.int(uintptr(u.CObj())), (*C.int)(unsafe.Pointer(target)), (*C.float2)(unsafe.Pointer(normal))))
}
func PortTestWorldMotionDispatch(u *server.Object) { C.sub_537770(asObjectC(u)) }
func PortTestWorldMotionGlobals() (map[string]*uint32, func()) {
	out := map[string]*uint32{"trace": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2488620))}
	for name, off := range map[string]uintptr{"grid-x": 2488612, "grid-y": 2488616, "fist-small": 2488624, "fist-medium": 2488628, "fist-large": 2488632, "scorch-ready": 2488636, "trap-reachable": 2491764, "trap-arrow": 2491768, "trap-one": 2491772, "trap-two": 2491776, "trap-fx-one": 2491780, "trap-fx-two": 2491784} {
		out[name] = memmap.PtrUint32(0x5D4594, off)
	}
	old := map[string]uint32{}
	for name, p := range out {
		old[name] = *p
		*p = 0
	}
	return out, func() {
		for name, p := range out {
			*p = old[name]
		}
	}
}
func PortTestWorldMotionScorch(pos *types.Pointf, size int32) {
	C.nox_xxx_sMakeScorch_537AF0((*C.float)(unsafe.Pointer(pos)), C.int(size))
}
func PortTestWorldMotionScorchInit() int32 { return int32(C.nox_xxx_scorchInit_537BD0()) }
func PortTestWorldMotionScorchNames() (*[6]uint32, func()) {
	p := (*[6]uint32)(memmap.PtrOff(0x587000, 276824))
	old := *p
	var frees []func()
	for i, name := range []string{"ScorchMarkFloorSmallA", "ScorchMarkFloorMediumA", "ScorchMarkFloorLargeB"} {
		s, free := alloc.CString(name)
		frees = append(frees, free)
		p[i*2] = uint32(uintptr(unsafe.Pointer(s)))
		p[i*2+1] = 0
	}
	return p, func() {
		*p = old
		for _, f := range frees {
			f()
		}
	}
}
func PortTestWorldMotionTimed(op string, u, target *server.Object, arg int32) int32 {
	raw := C.int(uintptr(u.CObj()))
	switch op {
	case "mover":
		C.nox_xxx_unitUpdateMover_54F740(raw)
	case "trap-update":
		return int32(C.nox_xxx_updateShootingTrap_54F9A0(raw))
	case "trap-projectile":
		C.nox_xxx_createArrowTrapProjectile_54FA80(raw, C.int(arg))
	case "trap-state":
		C.sub_54FBB0(raw)
	case "trap-scan":
		return int32(C.sub_54FBF0(raw))
	case "trap-candidate":
		C.nox_xxx_unitIsAttackReachable_54FC50(C.int(uintptr(target.CObj())), raw)
	case "trigger":
		C.nox_xxx_collideTrigger_54FCD0(raw, C.int(uintptr(target.CObj())))
	default:
		panic(op)
	}
	return 0
}
