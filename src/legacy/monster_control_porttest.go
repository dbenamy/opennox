//go:build porttest

package legacy

/*
#include <stdint.h>
extern unsigned int nox_gameDisableMapDraw_5d4594_2650672;
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

// The native adapter preserves the original operation/argument contract and
// identifies the same real object, stack and cache owners as the frozen C tests.
func PortTestMonsterControl(op string, u, target *server.Object, args *[8]uint32, arg int32) uint32 {
	switch op {
	case "death":
		return monsterControlDeath(u)
	case "revive":
		return monsterControlRevive(u)
	case "chapter":
		return uint32(monsterControlChapter())
	case "animation":
		monsterControlAnimation(u)
	case "head":
		return monsterControlHead(u)
	case "ensure":
		return uint32(uintptr(monsterControlEnsure(u, uint32(arg)).C()))
	case "refresh":
		return uint32(monsterControlRefresh(u))
	case "look":
		return monsterControlLook(u, arg)
	case "walk":
		monsterControlWalk(u, types.Pointf{X: math.Float32frombits(args[0]), Y: math.Float32frombits(args[1])})
	case "patrol":
		monsterControlPatrol(u, types.Pointf{X: math.Float32frombits(args[0]), Y: math.Float32frombits(args[1])}, types.Pointf{X: math.Float32frombits(args[2]), Y: math.Float32frombits(args[3])}, math.Float32frombits(args[4]))
	case "hunt":
		monsterControlIdle(u, true)
	case "idle":
		monsterControlIdle(u, false)
	case "follow":
		monsterControlFollow(u, target)
	case "melee":
		monsterControlMelee(u, (*types.Pointf)(unsafe.Pointer(args)))
	case "missile":
		monsterControlMissile(u, (*types.Pointf)(unsafe.Pointer(args)))
	case "control-byte":
		return monsterControlByte(u, (*byte)(unsafe.Pointer(args)))
	case "fight":
		monsterControlFight(u, target)
	case "flee":
		var target *server.Object
		var duration uint32
		if args != nil {
			target = motionObject(args[0])
			duration = args[1]
		}
		monsterControlFlee(u, target, duration)
	case "wait":
		monsterControlWait(u, uint32(arg))
	default:
		panic(op)
	}
	return 0
}
func PortTestMonsterCacheInit() (*uint32, func()) {
	old := monsterScriptCacheCold
	monsterScriptCacheCold = 1
	return &monsterScriptCacheCold, func() { monsterScriptCacheCold = old }
}
func PortTestMonsterCache(op string, u *server.Object, id int32) uint32 {
	switch op {
	case "reset":
		return monsterCacheReset()
	case "prepare":
		return monsterCachePrepare(u)
	case "find":
		return motionAddress(monsterCacheFind(id))
	case "remove":
		return monsterCacheRemove(u)
	case "clear":
		return monsterCacheClear()
	default:
		panic(op)
	}
}
func PortTestMonsterPendingOwner() (*uint32, func()) {
	oldPool, oldHead := monsterPendingClass, monsterPendingHead
	monsterPendingClass = nil
	monsterPendingHead = 0
	return &monsterPendingHead, func() {
		if monsterPendingClass != nil {
			monsterPendingFree()
		}
		monsterPendingClass = oldPool
		monsterPendingHead = oldHead
	}
}
func PortTestMonsterPending(op string, a, b int32) uint32 {
	switch op {
	case "init":
		return uint32(monsterPendingInit())
	case "free":
		return monsterPendingFree()
	case "clear":
		monsterPendingClear()
	case "add":
		return monsterPendingAdd(uint32(a), uint32(b))
	case "resolve":
		monsterPendingResolve()
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
