//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_2488620;
static uint32_t motionRadialWords[65];
static void motionRadialObserve(uint32_t* unit,uint32_t code) {uint32_t n=motionRadialWords[0];if(n>=32)abort();motionRadialWords[1+2*n]=(uintptr_t)unit;motionRadialWords[2+2*n]=code;motionRadialWords[0]=n+1;}
static void* motionRadialCallback(void) { return motionRadialObserve; }

static uint32_t* motionRadialSnapshot(void) {return motionRadialWords;}
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
		return motionDecaySet(u, delay)
	case "decay-remove":
		return motionDecayRemove(u)
	case "decay-tick":
		motionDecayTick()
	case "decay-clear":
		return motionDecayClear()
	case "sentry-update":
		return motionSentryUpdate(u)
	case "sentry-remove":
		return motionSentryRemove(u)
	case "sentry-clear":
		motionSentryHead = 0
	default:
		panic(op)
	}
	return 0
}
func PortTestWorldMotionListGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"decay": &motionDecayHead, "sentry": &motionSentryHead}
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
		motionProjectileStep(u)
	case "fall":
		motionFall(u)
	case "activate-nonsimple":
		return motionActivate(u)
	case "activate":
		return collisionActivate(u)
	case "remove-nonsimple":
		motionDeactivate(u)
	default:
		panic(op)
	}
	return 0
}

func PortTestWorldMotionVelocity(step float32) int32 {
	return motionVelocity(step)
}
func PortTestWorldMotionVelocityGlobals() (*[10]uint32, func()) {
	p := &motionVelocityTypes
	old := *p
	*p = [10]uint32{}
	return p, func() { *p = old }
}

func PortTestWorldMotionSentryCandidate(target, sentry *server.Object, ray *[4]float32) {
	motionSentryContact(target, sentry, ray)
}
func PortTestWorldMotionSentryReport(player int32) { motionSentryReport(player) }
func PortTestWorldMotionSentryPacket(player int32, u *server.Object) int32 {
	return motionSentryPacket(player, u)
}

func PortTestWorldMotionRoundScratch() (*[4]uint32, func()) {
	p := (*[4]uint32)(memmap.PtrOff(0x5D4594, 527668))
	old := *p
	*p = [4]uint32{}
	return p, func() { *p = old }
}

func PortTestWorldMotionTrace(u *server.Object, target *uint32, normal *types.Pointf) int8 {
	return motionTrace(u, target, normal)
}
func PortTestWorldMotionDispatch(u *server.Object) { motionProjectileDispatch(u) }
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
	motionScorch(pos, size)
}
func PortTestWorldMotionScorchInit() int32 { return motionScorchInit() }
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
	switch op {
	case "mover":
		motionMover(u)
	case "trap-update":
		return motionTrapUpdate(u)
	case "trap-projectile":
		motionTrapProjectile(u, arg)
	case "trap-state":
		motionTrapState(u)
	case "trap-scan":
		return motionTrapScan(u)
	case "trap-candidate":
		motionTrapCandidate(target, u)
	case "trigger":
		motionTrigger(u, target)
	default:
		panic(op)
	}
	return 0
}

func PortTestWorldMotionRadial(p *types.Pointf, radius float32, code uint32) (uint32, [][2]uint32) {
	rv := motionRadial(p, radius, C.motionRadialCallback(), code)
	words := (*[65]uint32)(unsafe.Pointer(C.motionRadialSnapshot()))
	var rows [][2]uint32
	for i := uint32(0); i < words[0]; i++ {
		rows = append(rows, [2]uint32{words[1+2*i], words[2+2*i]})
	}
	*words = [65]uint32{}
	return rv, rows
}
