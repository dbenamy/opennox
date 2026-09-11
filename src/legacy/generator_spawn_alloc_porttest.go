//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_1.h"
#include "common/alloc/classes/alloc_class.h"
extern void* nox_alloc_spawn_2386216;
extern void* nox_alloc_monsterList_2386220;
extern uint32_t dword_5d4594_2386212;
extern uint32_t dword_5d4594_2386224;
extern uint32_t dword_5d4594_2386228;
void sub_50D7E0(void);
static void portTestGeneratorSpawnClassesFree(void) {
	if (nox_alloc_spawn_2386216 && nox_alloc_monsterList_2386220) {
		sub_50D820();
		return;
	}
	if (nox_alloc_spawn_2386216) nox_free_alloc_class(nox_alloc_spawn_2386216);
	if (nox_alloc_monsterList_2386220) nox_free_alloc_class(nox_alloc_monsterList_2386220);
	nox_alloc_spawn_2386216 = 0;
	nox_alloc_monsterList_2386220 = 0;
	dword_5d4594_2386212 = 0;
	dword_5d4594_2386224 = 0;
	dword_5d4594_2386228 = 0;
}
*/
import "C"

import "unsafe"

// portTestGeneratorSpawnNode is a normalized SpawnClass record. Object is the
// real child object address; Prev/Next are list ordinals (-1 for nil), so the
// fixture can compare topology without comparing allocator addresses.
type portTestGeneratorSpawnNode struct {
	Object     uintptr
	Prev, Next int
}

type portTestGeneratorSpawnAllocatorSnapshot struct {
	SpawnClass, MonsterListClass uintptr
	SpawnHead, MonsterListHead   uintptr
	MonsterListCount             uint32
	Spawn                        []portTestGeneratorSpawnNode
}

// portTestGeneratorSpawnAllocator detaches any pre-existing allocator state,
// creates the real 50D780 SpawnClass/MonsterListClass pair, and later restores
// the detached state exactly. reset frees only this fixture's allocated class
// objects (50D7E0), retaining the classes for the next case. The caller must
// first release spawned child objects and clear their UpdateData links; C
// 50D7E0 intentionally invalidates those SpawnClass records.
func portTestGeneratorSpawnAllocator() (reset func(), snapshot func() portTestGeneratorSpawnAllocatorSnapshot, restore func()) {
	oldSpawnClass, oldMonsterClass := C.nox_alloc_spawn_2386216, C.nox_alloc_monsterList_2386220
	oldSpawnHead, oldMonsterHead, oldMonsterCount := C.dword_5d4594_2386212, C.dword_5d4594_2386224, C.dword_5d4594_2386228

	// Do not call C cleanup while the old classes are installed: they belong to
	// the surrounding server. 50D780 initializes these exact five globals.
	C.nox_alloc_spawn_2386216 = nil
	C.nox_alloc_monsterList_2386220 = nil
	C.dword_5d4594_2386212 = 0
	C.dword_5d4594_2386224 = 0
	C.dword_5d4594_2386228 = 0
	if C.nox_xxx_allocMonsterRelatedArrays_50D780() == 0 {
		// A partial 50D780 allocation, if any, is owned by this fixture.
		C.portTestGeneratorSpawnClassesFree()
		C.nox_alloc_spawn_2386216, C.nox_alloc_monsterList_2386220 = oldSpawnClass, oldMonsterClass
		C.dword_5d4594_2386212, C.dword_5d4594_2386224, C.dword_5d4594_2386228 = oldSpawnHead, oldMonsterHead, oldMonsterCount
		panic("generator SpawnClass allocation failed")
	}
	ownSpawnClass, ownMonsterClass := C.nox_alloc_spawn_2386216, C.nox_alloc_monsterList_2386220

	checkOwn := func() {
		if C.nox_alloc_spawn_2386216 != ownSpawnClass || C.nox_alloc_monsterList_2386220 != ownMonsterClass {
			panic("generator spawn allocator ownership changed")
		}
	}
	reset = func() {
		checkOwn()
		C.sub_50D7E0()
	}
	snapshot = func() (out portTestGeneratorSpawnAllocatorSnapshot) {
		checkOwn()
		out.SpawnClass = uintptr(C.nox_alloc_spawn_2386216)
		out.MonsterListClass = uintptr(C.nox_alloc_monsterList_2386220)
		out.SpawnHead = uintptr(C.dword_5d4594_2386212)
		out.MonsterListHead = uintptr(C.dword_5d4594_2386224)
		out.MonsterListCount = uint32(C.dword_5d4594_2386228)

		byAddr := make(map[uintptr]int)
		for p := uintptr(C.dword_5d4594_2386212); p != 0; {
			if len(out.Spawn) == 96 { // SpawnClass capacity from 50D780.
				panic("generator SpawnClass list exceeds capacity")
			}
			if _, exists := byAddr[p]; exists {
				panic("cycle in generator SpawnClass list")
			}
			byAddr[p] = len(out.Spawn)
			words := unsafe.Slice((*uint32)(unsafe.Pointer(p)), 3)
			out.Spawn = append(out.Spawn, portTestGeneratorSpawnNode{Object: uintptr(words[0]), Prev: -2, Next: -2})
			p = uintptr(words[1])
		}
		for p, i := range byAddr {
			words := unsafe.Slice((*uint32)(unsafe.Pointer(p)), 3)
			if words[2] != 0 {
				prev, ok := byAddr[uintptr(words[2])]
				if !ok {
					panic("generator SpawnClass previous link outside list")
				}
				out.Spawn[i].Prev = prev
			} else {
				out.Spawn[i].Prev = -1
			}
			if words[1] != 0 {
				next, ok := byAddr[uintptr(words[1])]
				if !ok {
					panic("generator SpawnClass next link outside list")
				}
				out.Spawn[i].Next = next
			} else {
				out.Spawn[i].Next = -1
			}
		}
		return out
	}
	restore = func() {
		checkOwn()
		C.sub_50D7E0()
		C.portTestGeneratorSpawnClassesFree()
		C.nox_alloc_spawn_2386216, C.nox_alloc_monsterList_2386220 = oldSpawnClass, oldMonsterClass
		C.dword_5d4594_2386212, C.dword_5d4594_2386224, C.dword_5d4594_2386228 = oldSpawnHead, oldMonsterHead, oldMonsterCount
	}
	return reset, snapshot, restore
}
