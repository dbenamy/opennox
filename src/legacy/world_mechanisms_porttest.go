//go:build porttest

package legacy

/*
#include <stdint.h>
// A retained char return can contain the low byte of the collision callback
// address. Align this recorder so that byte is stable across builds and ASLR.
static uint32_t worldCalls[256]; static int worldCount;
static int __attribute__((aligned(256))) worldCollide(int u,int a,int b) {
 worldCalls[worldCount++]=u;worldCalls[worldCount++]=a;worldCalls[worldCount++]=b;return 0;
}
static void* worldCollidePtr(void){return worldCollide;}
static void worldReset(void){worldCount=0;}
static int worldN(void){return worldCount;}
static uint32_t worldValue(int i){return worldCalls[i];}
*/
import "C"
import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

const (
	PortTestWorld53AC50 = 700
	PortTestWorld53B030 = 701
	PortTestWorld53B060 = 702
	PortTestWorld53B1B0 = 703
	PortTestWorld53B300 = 704
	PortTestWorld53B320 = 705
	PortTestWorld53B380 = 706
	PortTestWorld53B410 = 707
	PortTestWorld53B490 = 708
	PortTestWorld53B5D0 = 709
	PortTestWorld53B750 = 710
	PortTestWorld53B860 = 711
	PortTestWorld53BEF0 = 712
	PortTestWorld53C060 = 713
	PortTestWorld53C0C0 = 714
	PortTestWorld53C140 = 715
	PortTestWorld53C160 = 716
	PortTestWorld53C240 = 717
	PortTestWorld53DE80 = 718
	PortTestWorld548830 = 719
	PortTestWorld548860 = 720
)

type PortTestWorldSpec struct {
	Objectives   *PortTestObjectivesSpec
	CollideWords []map[int]uint32
	ItemNames    []string
}
type portTestWorld struct {
	objectives *portTestObjectives
	blocks     [][]byte
	frees      []func()
	oldCD      []unsafe.Pointer
}

func (p *portTestShopPools) worldPrepare() func() {
	if p.proxy.callbacks.shop.spec.TemporaryUpdates.World == nil {
		return func() {}
	}
	p.temporary.world = &portTestWorld{}
	restoreTypes := p.proxy.core.PortTestWorldTypes()
	// The retained absolute-value helper uses this original relocated scratch pointer.
	scratch := memmap.PtrUint32(0x5d4594, 527672)
	slot := memmap.PtrPtr(0x587000, 55744)
	oldScratch, oldSlot := *scratch, *slot
	*scratch, *slot = 0, unsafe.Pointer(scratch)
	oldQueue := collisionAngleHead
	collisionAngleHead = 0
	a, b := memmap.PtrUint32(0x5d4594, 2488680), memmap.PtrUint32(0x5d4594, 2488676)
	oldA, oldB := *a, *b
	*a, *b = 0, 0
	restoreObjectives := p.objectivesPrepare()
	return func() {
		restoreObjectives()
		w := p.temporary.world
		for i, cd := range w.oldCD {
			// Filter-only class probes own no monster-generator children.
			// Snapshots are already captured; generic freeing must not follow UD as a generator.
			p.items[i].u.ObjClass &^= object.ClassMonsterGenerator
			p.items[i].u.CollideData = cd
		}
		for _, f := range w.frees {
			f()
		}
		*scratch, *slot = oldScratch, oldSlot
		collisionAngleHead = oldQueue
		*a, *b = oldA, oldB
		restoreTypes()
	}
}
func (p *portTestShopPools) worldItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World
	if sp == nil {
		return
	}
	w := p.temporary.world
	C.worldReset()
	p.identify(C.worldCollidePtr(), 71300)
	for i := 0; i < 21; i++ {
		p.identify(portTestWorldFunction(i), 71000+uint32(i))
	}
	for i, it := range p.items {
		it.u.Collide = C.worldCollidePtr()
		b, free := alloc.Make([]byte{}, 80)
		for j := 0; j < 8; j++ {
			b[j] = 0xa5
			b[72+j] = 0x5a
		}
		w.blocks = append(w.blocks, b)
		w.frees = append(w.frees, free)
		w.oldCD = append(w.oldCD, it.u.CollideData)
		it.u.CollideData = unsafe.Pointer(&b[8])
		p.identify(it.u.CollideData, 72000+uint32(i))
		if i < len(sp.CollideWords) {
			for off, v := range sp.CollideWords[i] {
				if off < 0 || off+4 > 64 || off%4 != 0 {
					panic("world collide offset")
				}
				*(*uint32)(unsafe.Add(it.u.CollideData, off)) = v
			}
		}
		if i < len(sp.ItemNames) && sp.ItemNames[i] != "" {
			id := p.proxy.core.Types.IndByID(sp.ItemNames[i])
			if id == 0 {
				panic("world item type")
			}
			it.u.TypeInd = uint16(id)
		}
	}
	p.objectivesItems()
}

func portTestWorldFunction(id int) unsafe.Pointer {
	switch id {
	case 0:
		return updateIdentityKey(updateIDDoor)
	case 1:
		return updateIdentityKey(updateIDPush)
	case 2:
		return updateIdentityKey(updateIDToggle)
	case 3:
		return updateIdentityKey(updateIDTrigger)
	case 4:
		return updateIdentityKey(updateIDLoopAndDamage)
	case 5:
		return updateIdentityKey(updateIDSwitch)
	case 6:
		return updateIdentityKey(updateIDElevatorShaft)
	case 9:
		return updateIdentityKey(updateIDElevator)
	case 11:
		return updateIdentityKey(updateIDPhantomPlayer)
	case 12:
		return updateIdentityKey(updateIDPentagram)
	case 14:
		return updateIdentityKey(updateIDInvisiblePentagram)
	case 16:
		return updateIdentityKey(updateIDBlow)
	case 18:
		return updateIdentityKey(updateIDTrapDoor)
	default:
		switch id {
		case 7:
			return portTestFixtureKey("nox_xxx_fnElevatorShaft_53B410")
		case 8:
			return portTestFixtureKey("nox_xxx_elevatorAud_53B490")
		case 10:
			return portTestFixtureKey("nox_xxx_elevatorFn_53B750")
		case 13:
			return portTestFixtureKey("nox_xxx_fnPentagramTeleport_53C060")
		case 15:
			return portTestFixtureKey("sub_53C140")
		case 17:
			return portTestFixtureKey("sub_53C240")
		case 19:
			return portTestFixtureKey("sub_548830")
		case 20:
			return portTestFixtureKey("sub_548860")
		default:
			return nil
		}
	}
}

func portTestWorldCall(id int, u, target *server.Object, value int) uint32 {
	switch id {
	case 0:
		return uint32(int32(int8(worldDoor(u))))
	case 1:
		worldPush(u)
		return 0
	case 2:
		return uint32(int32(int8(worldToggle(u))))
	case 3:
		return uint32(int32(int8(worldTrigger(u))))
	case 4:
		return uint32(int32(int8(worldEnabledCollision(u))))
	case 5:
		return uint32(int32(int8(worldSwitch(u))))
	case 6:
		return uint32(int32(int8(worldShaft(u))))
	case 9:
		worldElevator(u)
		return 0
	case 11:
		worldPhantom(u)
		return 0
	case 12:
		return uint32(int32(worldTeleport(u)))
	case 14:
		return uint32(int32(worldInvisibleTeleport(u)))
	case 16:
		worldBlow(u)
		return 0
	case 18:
		return worldTrapDoor(u)
	default:
		switch id {
		case 7:
			portTestInvoke_nox_xxx_fnElevatorShaft_53B410(C.int(uintptr(target.CObj())), C.int(uintptr(u.CObj())))
		case 8:
			portTestInvoke_nox_xxx_elevatorAud_53B490(C.int(uintptr(u.CObj())), C.int(value))
		case 10:
			portTestInvoke_nox_xxx_elevatorFn_53B750(C.int(uintptr(target.CObj())), C.int(uintptr(u.CObj())))
		case 13:
			portTestInvoke_nox_xxx_fnPentagramTeleport_53C060((*C.float)(unsafe.Pointer(target.CObj())), C.int(uintptr(unsafe.Add(u.CObj(), 56))))
		case 15:
			portTestInvoke_sub_53C140((*C.float)(unsafe.Pointer(target.CObj())), C.int(uintptr(unsafe.Add(u.CObj(), 56))))
		case 17:
			portTestInvoke_sub_53C240((*C.float)(unsafe.Pointer(target.CObj())), C.int(uintptr(u.CObj())))
		case 19:
			portTestInvoke_sub_548830(C.int(*(*uint32)(unsafe.Add(u.CObj(), 748))))
		case 20:
			portTestInvoke_sub_548860(C.int(uintptr(u.CObj())), C.short(value))
		}
		return 0
	}
}

func (p *portTestShopPools) worldAction(a PortTestShopAction) uint32 {
	if p.registeredUpdateAction(a) {
		return p.temporary.result
	}
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	p.temporary.result = portTestWorldCall(a.Op-700, p.items[a.Item].u, p.temporaryRef(sp.Target), int(a.Value))
	return p.temporary.result
}
func portTestWorldCollisionCheck(actor, target *server.Object, called bool) {
	if !called {
		if C.worldN() != 0 {
			panic("unexpected world collision callback")
		}
		return
	}
	if C.worldN() != 3 || uint32(C.worldValue(0)) != uint32(uintptr(target.CObj())) || uint32(C.worldValue(1)) != uint32(uintptr(actor.CObj())) || C.worldValue(2) != 0 {
		panic("world collision callback count or arguments")
	}
}

func (p *portTestShopPools) worldSnapshot(out []uint32) []uint32 {
	w := p.temporary.world
	if w == nil {
		return out
	}
	for _, v := range []uint32{*memmap.PtrUint32(0x5d4594, 527672), uint32(collisionAngleHead), uint32(collisionActiveHead), uint32(collisionActiveTail), *memmap.PtrUint32(0x5d4594, 2488680), *memmap.PtrUint32(0x5d4594, 2488676)} {
		out = append(out, p.normalize(v))
	}
	out = append(out, uint32(C.worldN()))
	for i := 0; i < int(C.worldN()); i++ {
		out = append(out, p.normalize(uint32(C.worldValue(C.int(i)))))
	}
	for _, b := range w.blocks {
		for j := 0; j < 8; j++ {
			if b[j] != 0xa5 || b[72+j] != 0x5a {
				panic("world collision guard")
			}
		}
		for _, v := range unsafe.Slice((*uint32)(unsafe.Pointer(&b[8])), 16) {
			out = append(out, p.normalize(v))
		}
	}
	return p.objectivesSnapshot(out)
}

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_nox_xxx_elevatorAud_53B490(a1 C.int, a2 C.int) {
	worldElevatorSound(objectFromInt(a1), a2 != 0)
}

func portTestInvoke_nox_xxx_elevatorFn_53B750(a1 C.int, a2 C.int) {
	worldElevatorCandidate(objectFromInt(a1), objectFromInt(a2))
}

func portTestInvoke_nox_xxx_fnElevatorShaft_53B410(a1 C.int, a2 C.int) {
	worldShaftCandidate(objectFromInt(a1), objectFromInt(a2))
}

func portTestInvoke_nox_xxx_fnPentagramTeleport_53C060(a1 *C.float, a2 C.int) {
	worldTeleportCandidate((*server.Object)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))), true)
}

func portTestInvoke_sub_53C140(a1 *C.float, a2 C.int) {
	worldTeleportCandidate((*server.Object)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))), false)
}

func portTestInvoke_sub_53C240(a1 *C.float, arg4 C.int) {
	worldBlowCandidate((*server.Object)(unsafe.Pointer(a1)), objectFromInt(arg4))
}

func portTestInvoke_sub_548830(a1 C.int) { worldAngleQueue(unsafe.Pointer(uintptr(uint32(a1)))) }

func portTestInvoke_sub_548860(a1 C.int, a2 C.short) { worldAngle(objectFromInt(a1), int16(a2)) }
