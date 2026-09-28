//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

var portTestEffectsParticleCalls int
var portTestEffectsParticleKey byte

func init() {
	effectParticleCallbacks[unsafe.Pointer(&portTestEffectsParticleKey)] = portTestEffectsParticleUpdate
}

func portTestEffectsParticleUpdate(p unsafe.Pointer) {
	words := unsafe.Slice((*uint32)(unsafe.Pointer(p)), 32)
	words[0]++
	words[20] += 0x10000
	words[21] ^= 0x10000
	portTestEffectsParticleCalls++
}

func PortTestEffectsParticleEnvironment() (unsafe.Pointer, *int, func()) {
	old := portTestEffectsParticleCalls
	portTestEffectsParticleCalls = 0
	return unsafe.Pointer(&portTestEffectsParticleKey), &portTestEffectsParticleCalls, func() { portTestEffectsParticleCalls = old }
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
	var pool alloc.ClassT[Nox_screenParticle]
	if capacity > 0 {
		pool = alloc.NewClassT("porttest screen particles", Nox_screenParticle{}, capacity)
		legacyGlobals.nox_alloc_screenParticles_806044 = pool.UPtr()
	} else {
		legacyGlobals.nox_alloc_screenParticles_806044 = nil
	}
	legacyGlobals.nox_screenParticles_head = nil
	legacyGlobals.dword_5d4594_806052 = nil
	snapshot := func() [][]uint32 {
		var nodes []*Nox_screenParticle
		ids := make(map[*Nox_screenParticle]uint32)
		var prev *Nox_screenParticle
		for p := (*Nox_screenParticle)(unsafe.Pointer(legacyGlobals.nox_screenParticles_head)); p != nil; p = p.Field_44 {
			if len(nodes) >= capacity || ids[p] != 0 || p.Field_48 != prev {
				panic("invalid screen particle list")
			}
			nodes = append(nodes, p)
			ids[p] = uint32(len(nodes))
			prev = p
		}
		if prev != (*Nox_screenParticle)(unsafe.Pointer(legacyGlobals.dword_5d4594_806052)) {
			panic("invalid screen particle tail")
		}
		var out [][]uint32
		for _, p := range nodes {
			words := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(p)), 13)...)
			if unsafe.Pointer(uintptr(words[0])) != screenParticleCallbackKey() {
				panic("unexpected screen draw callback")
			}
			words[0] = 0xe2000001
			words[11] = ids[p.Field_44]
			words[12] = ids[p.Field_48]
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
