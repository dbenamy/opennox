//go:build porttest

package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_2491552;
extern uint32_t dword_5d4594_2488604;
extern uint32_t dword_5d4594_2488608;
char nox_xxx_updateDoor_53AC50(int a1);
void nox_xxx_updatePush_53B030(int a1);
char nox_xxx_updateToggle_53B060(uint32_t* a1);
char nox_xxx_updateTrigger_53B1B0(int a1);
char sub_53B300(int a1);
char nox_xxx_updateSwitch_53B320(uint32_t* a1);
char nox_xxx_updateElevatorShaft_53B380(int a1);
void nox_xxx_fnElevatorShaft_53B410(int a1, int a2);
void nox_xxx_elevatorAud_53B490(int a1, int a2);
void nox_xxx_updateElevator_53B5D0(uint32_t* a1);
void nox_xxx_elevatorFn_53B750(int a1, int a2);
void nox_xxx_updatePhantomPlayer_53B860(int a1);
int nox_xxx_updateTeleportPentagram_53BEF0(int a1);
void nox_xxx_fnPentagramTeleport_53C060(float* a1, int a2);
int nox_xxx_updateInvisiblePentagram_53C0C0(int a1);
void sub_53C140(float* a1, int a2);
void nox_xxx_updateBlow_53C160(int a3);
void sub_53C240(float* a1, int arg4);
int* nox_xxx_updateTrapDoor_53DE80(uint32_t* a1);
void sub_548830(int a1);
void sub_548860(int a1, short a2);
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
static void* worldFunction(int id) {switch(id){
case 0: return (void*)nox_xxx_updateDoor_53AC50;
case 1: return (void*)nox_xxx_updatePush_53B030;
case 2: return (void*)nox_xxx_updateToggle_53B060;
case 3: return (void*)nox_xxx_updateTrigger_53B1B0;
case 4: return (void*)sub_53B300;
case 5: return (void*)nox_xxx_updateSwitch_53B320;
case 6: return (void*)nox_xxx_updateElevatorShaft_53B380;
case 7: return (void*)nox_xxx_fnElevatorShaft_53B410;
case 8: return (void*)nox_xxx_elevatorAud_53B490;
case 9: return (void*)nox_xxx_updateElevator_53B5D0;
case 10: return (void*)nox_xxx_elevatorFn_53B750;
case 11: return (void*)nox_xxx_updatePhantomPlayer_53B860;
case 12: return (void*)nox_xxx_updateTeleportPentagram_53BEF0;
case 13: return (void*)nox_xxx_fnPentagramTeleport_53C060;
case 14: return (void*)nox_xxx_updateInvisiblePentagram_53C0C0;
case 15: return (void*)sub_53C140;
case 16: return (void*)nox_xxx_updateBlow_53C160;
case 17: return (void*)sub_53C240;
case 18: return (void*)nox_xxx_updateTrapDoor_53DE80;
case 19: return (void*)sub_548830;
case 20: return (void*)sub_548860;
default:return 0;}}
static uint32_t worldCall(int id,nox_object_t* u,nox_object_t* target,int value){switch(id){
case 0: return (uint32_t)nox_xxx_updateDoor_53AC50((int)u);
case 1: nox_xxx_updatePush_53B030((int)u);return 0;
case 2: return (uint32_t)nox_xxx_updateToggle_53B060((uint32_t*)u);
case 3: return (uint32_t)nox_xxx_updateTrigger_53B1B0((int)u);
case 4: return (uint32_t)sub_53B300((int)u);
case 5: return (uint32_t)nox_xxx_updateSwitch_53B320((uint32_t*)u);
case 6: return (uint32_t)nox_xxx_updateElevatorShaft_53B380((int)u);
case 7: nox_xxx_fnElevatorShaft_53B410((int)target,(int)u);return 0;
case 8: nox_xxx_elevatorAud_53B490((int)u,value);return 0;
case 9: nox_xxx_updateElevator_53B5D0((uint32_t*)u);return 0;
case 10: nox_xxx_elevatorFn_53B750((int)target,(int)u);return 0;
case 11: nox_xxx_updatePhantomPlayer_53B860((int)u);return 0;
case 12: return (uint32_t)nox_xxx_updateTeleportPentagram_53BEF0((int)u);
case 13: nox_xxx_fnPentagramTeleport_53C060((float*)target,(int)((char*)u+56));return 0;
case 14: return (uint32_t)nox_xxx_updateInvisiblePentagram_53C0C0((int)u);
case 15: sub_53C140((float*)target,(int)((char*)u+56));return 0;
case 16: nox_xxx_updateBlow_53C160((int)u);return 0;
case 17: sub_53C240((float*)target,(int)u);return 0;
case 18: return (uint32_t)nox_xxx_updateTrapDoor_53DE80((uint32_t*)u);
case 19: sub_548830((int)*(uint32_t*)((char*)u+748));return 0;
case 20: sub_548860((int)u,(short)value);return 0;
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
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
	oldQueue := C.dword_5d4594_2491552
	C.dword_5d4594_2491552 = 0
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
		C.dword_5d4594_2491552 = oldQueue
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
		p.identify(C.worldFunction(C.int(i)), 71000+uint32(i))
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
func (p *portTestShopPools) worldAction(a PortTestShopAction) uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	p.temporary.result = uint32(C.worldCall(C.int(a.Op-700), asObjectC(p.items[a.Item].u), asObjectC(p.temporaryRef(sp.Target)), C.int(a.Value)))
	return p.temporary.result
}
func (p *portTestShopPools) worldSnapshot(out []uint32) []uint32 {
	w := p.temporary.world
	if w == nil {
		return out
	}
	for _, v := range []uint32{*memmap.PtrUint32(0x5d4594, 527672), uint32(C.dword_5d4594_2491552), uint32(C.dword_5d4594_2488604), uint32(C.dword_5d4594_2488608), *memmap.PtrUint32(0x5d4594, 2488680), *memmap.PtrUint32(0x5d4594, 2488676)} {
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
