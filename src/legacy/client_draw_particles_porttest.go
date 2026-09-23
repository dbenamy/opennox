//go:build porttest

package legacy

/*
#include "GAME2_2.h"
#include "GAME3.h"
#include "client__draw__magicdrw.h"
#include "client__draw__bubbledraw.h"
#include "client__draw__partrain.h"
#include "client__draw__lvupdraw.h"
#include "client__draw__spiderspitdraw.h"
#include "client__draw__vortexdraw.h"
*/
import "C"

import (
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"math"
	"unsafe"
)

func PortTestClientDrawParticle(op int, vp *noxrender.Viewport, dr *client.Drawable, a [4]int32) int64 {
	view := (*C.int)(vp.C())
	words := (*C.uint32_t)(vp.C())
	drawable := (*C.nox_drawable)(dr.C())
	vi := C.int(uintptr(vp.C()))
	light := unsafe.Add(dr.C(), 136)
	li := C.int(uintptr(light))
	switch op {
	case 0:
		return int64(C.nox_thing_magic_draw(view, drawable))
	case 1:
		return int64(C.nox_thing_magic_missle_draw(view, drawable))
	case 2:
		return int64(C.nox_thing_magic_missle_tail_link_draw(words, drawable))
	case 3:
		return int64(C.nox_thing_magic_tail_link_draw(words, drawable))
	case 4:
		return int64(C.nox_thing_drain_mana_draw())
	case 5:
		return int64(C.nox_thing_bubble_draw(words, drawable))
	case 6:
		return int64(C.nox_thing_blue_rain_draw(vi, drawable))
	case 7:
		return int64(C.nox_thing_levelup_draw(vi, drawable))
	case 8:
		return int64(C.nox_thing_oblivion_up_draw(vi, drawable))
	case 9:
		return int64(uint32(uintptr(unsafe.Pointer(particleFallingSparks(int(a[0]), vp, dr)))))
	case 10:
		return int64(C.nox_thing_spider_spit_draw(words, drawable))
	case 11:
		return int64(C.nox_thing_vortex_draw(view, drawable))
	case 12:
		got := particleLightColor(light, int(a[0]), int(a[1]), int(a[2]))
		if got == light {
			return 1
		}
		return 0
	case 13:
		return int64(C.sub_484C00(li, C.int(a[0])))
	case 14:
		return int64(C.nox_xxx_spriteChangeLightSize_484C30(li, C.int(a[0])))
	case 15:
		return int64(C.sub_484CE0(li, C.float(math.Float32frombits(uint32(a[0])))))
	case 16:
		return int64(nox_xxx_spriteChangeIntensity_484D70_light_intensity(li, C.float(math.Float32frombits(uint32(a[0])))))
	case 17:
		return int64(initParticlePalettes())
	case 18:
		return int64(initParticleColors())
	}
	panic("unknown particle draw operation")
}

type PortTestClientParticleEnvironment struct {
	mapped    []portTestEffectsRegion
	named     [2]uint32
	constants [2]uint64
}

func PortTestNewClientParticleEnvironment() *PortTestClientParticleEnvironment {
	e := &PortTestClientParticleEnvironment{named: [2]uint32{uint32(dword_5d4594_1313804), uint32(dword_5d4594_1313816)}, constants: [2]uint64{uint64(qword_581450_9544), uint64(qword_581450_9552)}}
	for _, r := range []struct {
		base, off uintptr
		size      int
	}{
		{0x5D4594, 1312500, 1028}, {0x5D4594, 1313708, 16}, {0x5D4594, 1313804, 24},
		{0x587000, 154972, 12}, {0x85B3FC, 956, 4},
	} {
		data := unsafe.Slice((*byte)(memmap.PtrOff(r.base, r.off)), r.size)
		initial := make([]byte, r.size)
		if r.base == 0x587000 {
			copy(initial, blobdata.PortTestClientParticleLightTable())
		}
		e.mapped = append(e.mapped, portTestEffectsRegion{data: data, old: append([]byte(nil), data...), initial: initial})
	}
	e.Reset()
	return e
}
func (e *PortTestClientParticleEnvironment) Reset() {
	for _, r := range e.mapped {
		copy(r.data, r.initial)
	}
	dword_5d4594_1313804 = 0
	dword_5d4594_1313816 = 0
	qword_581450_9544 = C.uint64_t(math.Float64bits(0.5))
	qword_581450_9552 = C.uint64_t(math.Float64bits(65536))
	*memmap.PtrUint32(0x85B3FC, 956) = noxcolor.RGB5551Color(0, 0, 0).Color32()
	initParticlePalettes()
}
func (e *PortTestClientParticleEnvironment) Constants(half, scale float64) {
	qword_581450_9544 = C.uint64_t(math.Float64bits(half))
	qword_581450_9552 = C.uint64_t(math.Float64bits(scale))
}
func (e *PortTestClientParticleEnvironment) Restore() {
	for _, r := range e.mapped {
		copy(r.data, r.old)
	}
	dword_5d4594_1313804 = C.uint32_t(e.named[0])
	dword_5d4594_1313816 = C.uint32_t(e.named[1])
	qword_581450_9544 = C.uint64_t(e.constants[0])
	qword_581450_9552 = C.uint64_t(e.constants[1])
}
func (e *PortTestClientParticleEnvironment) Snapshot() []uint32 {
	out := []uint32{uint32(dword_5d4594_1313804), uint32(dword_5d4594_1313816), uint32(qword_581450_9544), uint32(qword_581450_9544 >> 32), uint32(qword_581450_9552), uint32(qword_581450_9552 >> 32)}
	for _, r := range e.mapped {
		out = append(out, unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(r.data))), len(r.data)/4)...)
	}
	return out
}

// Addresses installed in the drawable's actual draw callback slot by the owner.
func PortTestClientParticleDrawCallback(op int) unsafe.Pointer {
	switch op {
	case 0:
		return unsafe.Pointer(C.nox_thing_magic_draw)
	case 1:
		return unsafe.Pointer(C.nox_thing_magic_missle_draw)
	case 2:
		return unsafe.Pointer(C.nox_thing_magic_missle_tail_link_draw)
	case 3:
		return unsafe.Pointer(C.nox_thing_magic_tail_link_draw)
	case 5:
		return unsafe.Pointer(C.nox_thing_bubble_draw)
	case 6:
		return unsafe.Pointer(C.nox_thing_blue_rain_draw)
	case 7:
		return unsafe.Pointer(C.nox_thing_levelup_draw)
	case 8:
		return unsafe.Pointer(C.nox_thing_oblivion_up_draw)
	case 10:
		return unsafe.Pointer(C.nox_thing_spider_spit_draw)
	case 11:
		return unsafe.Pointer(C.nox_thing_vortex_draw)
	}
	return nil
}
