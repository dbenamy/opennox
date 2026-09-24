//go:build porttest

package legacy

/*
#include "defs.h"
#include "client__draw__partscrn.h"



extern void portTestEffectsParticleUpdate(uint32_t*);
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

var portTestEffectsParticleCalls int

//export portTestEffectsParticleUpdate
func portTestEffectsParticleUpdate(p *C.uint32_t) {
	words := unsafe.Slice((*uint32)(unsafe.Pointer(p)), 32)
	words[0]++
	words[20] += 0x10000
	words[21] ^= 0x10000
	portTestEffectsParticleCalls++
}

func PortTestEffectsParticleEnvironment() (unsafe.Pointer, *int, func()) {
	old := portTestEffectsParticleCalls
	portTestEffectsParticleCalls = 0
	return unsafe.Pointer(C.portTestEffectsParticleUpdate), &portTestEffectsParticleCalls, func() { portTestEffectsParticleCalls = old }
}

// PortTestEffectsScreenParticles uses the production allocation/list machinery
// behind screen-particle creation. A small capacity also exercises tail reuse.
func PortTestEffectsScreenParticles(capacity int) (func() [][]uint32, func()) {
	oldPool, oldHead, oldTail := legacyGlobals.nox_alloc_screenParticles_806044, legacyGlobals.nox_screenParticles_head, legacyGlobals.dword_5d4594_806052
	colors := unsafe.Slice(memmap.PtrUint32(0x5D4594, 806004), 10)
	oldColors := append([]uint32(nil), colors...)
	for i := range colors {
		colors[i] = uint32(i+1) * 0x421
	}
	var pool alloc.ClassT[C.nox_screenParticle]
	if capacity > 0 {
		pool = alloc.NewClassT("porttest screen particles", C.nox_screenParticle{}, capacity)
		legacyGlobals.nox_alloc_screenParticles_806044 = pool.UPtr()
	} else {
		legacyGlobals.nox_alloc_screenParticles_806044 = nil
	}
	legacyGlobals.nox_screenParticles_head = nil
	legacyGlobals.dword_5d4594_806052 = nil
	snapshot := func() [][]uint32 {
		var nodes []*C.nox_screenParticle
		ids := make(map[*C.nox_screenParticle]uint32)
		var prev *C.nox_screenParticle
		for p := (*C.nox_screenParticle)(unsafe.Pointer(legacyGlobals.nox_screenParticles_head)); p != nil; p = p.field_44 {
			if len(nodes) >= capacity || ids[p] != 0 || p.field_48 != prev {
				panic("invalid screen particle list")
			}
			nodes = append(nodes, p)
			ids[p] = uint32(len(nodes))
			prev = p
		}
		if prev != (*C.nox_screenParticle)(unsafe.Pointer(legacyGlobals.dword_5d4594_806052)) {
			panic("invalid screen particle tail")
		}
		var out [][]uint32
		for _, p := range nodes {
			words := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(p)), 13)...)
			if unsafe.Pointer(uintptr(words[0])) != unsafe.Pointer(C.nox_client_screenParticleDraw_489700) {
				panic("unexpected screen draw callback")
			}
			words[0] = 0xe2000001
			words[11] = ids[p.field_44]
			words[12] = ids[p.field_48]
			out = append(out, words)
		}
		return out
	}
	return snapshot, func() {
		legacyGlobals.nox_alloc_screenParticles_806044 = oldPool
		legacyGlobals.nox_screenParticles_head = oldHead
		legacyGlobals.dword_5d4594_806052 = oldTail
		copy(colors, oldColors)
		if capacity > 0 {
			pool.Free()
		}
	}
}
