//go:build porttest

package legacy

/*
#include "GAME1_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "client__draw__partscrn.h"
*/
import "C"
import (
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"image"
	"unsafe"
)

func PortTestScreenMaidenCallback() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_maiden_draw)
}
func PortTestScreenEffectCallback(op int) unsafe.Pointer {
	switch op {
	case 0:
		return drawableDrawKey(drawKey_nox_thing_harpoon_draw)
	case 1:
		return drawableDrawKey(drawKey_nox_thing_harpoon_rope_draw)
	case 2:
		return drawableDrawKey(drawKey_nox_thing_undead_killer_draw)
	case 3:
		return drawableDrawKey(drawKey_nox_thing_player_waypoint_draw)
	case 4:
		return drawableDrawKey(drawKey_nox_thing_maiden_draw)
	}
	panic("unknown screen effect callback")
}
func PortTestScreenParticleCreate(a [10]int32) *Nox_screenParticle {
	return (*Nox_screenParticle)(unsafe.Pointer(C.nox_client_newScreenParticle_431540(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5]), C.char(a[6]), C.char(a[7]), C.char(a[8]), C.char(a[9]))))
}
func PortTestScreenParticleDraw(p *Nox_screenParticle, vp *noxrender.Viewport) int {
	return int(C.nox_client_screenParticleDraw_489700(vp.C(), (*C.nox_screenParticle)(unsafe.Pointer(p))))
}
func PortTestScreenParticleDelete(p *Nox_screenParticle) {
	screenParticleDelete(p)
}
func PortTestScreenPrimitive(op int, a [4]int32) uint32 {
	switch op {
	case 0:
		return screenSqrt(uint32(a[0]))
	case 1:
		return uint32(C.sub_48C6B0(C.int(a[0]), C.int(a[1])))
	case 2:
		return screenDistanceBetween(a[0], a[1], a[2], a[3])
	case 3:
		return uint32(screenRopeLine(image.Pt(int(a[0]), int(a[1])), image.Pt(int(a[2]), int(a[3]))))
	case 4:
		return uint32(screenCircle(int(a[0]), int(a[1]), int(a[2]), uint32(a[3])))
	}
	panic("unknown screen primitive")
}

type PortTestScreenEnvironment struct{ restore []func() }

func PortTestNewScreenEnvironment() *PortTestScreenEnvironment {
	e := new(PortTestScreenEnvironment)
	old := [4]uint32{dword_5d4594_3807140, dword_5d4594_3807136, dword_5d4594_3807116, dword_5d4594_3807152}
	e.restore = append(e.restore, func() {
		dword_5d4594_3807140, dword_5d4594_3807136, dword_5d4594_3807116, dword_5d4594_3807152 = old[0], old[1], old[2], old[3]
	})
	for _, reg := range [][3]uintptr{{0x5D4594, 1312492, 8}, {0x5D4594, 1313728, 12}, {0x5D4594, 1197368, 4}, {0x5D4594, 1200916, 512}, {0x85B3FC, 940, 4}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(reg[0], reg[1])), reg[2])
		saved := append([]byte(nil), b...)
		clear(b)
		e.restore = append(e.restore, func() { copy(b, saved) })
	}
	e.Reset()
	return e
}
func (e *PortTestScreenEnvironment) Reset() {
	dword_5d4594_3807140, dword_5d4594_3807136, dword_5d4594_3807116, dword_5d4594_3807152 = 0, 0, 0, 0
	for _, reg := range [][2]uintptr{{1312492, 2}, {1313728, 3}} {
		clear(unsafe.Slice(memmap.PtrUint32(0x5D4594, reg[0]), reg[1]))
	}
	*memmap.PtrUint32(0x85B3FC, 940) = uint32(noxcolor.RGB5551Color(80, 220, 250).Color32())
}
func (e *PortTestScreenEnvironment) State() []uint32 {
	out := []uint32{uint32(dword_5d4594_3807140), uint32(dword_5d4594_3807136), uint32(dword_5d4594_3807116), uint32(dword_5d4594_3807152)}
	for _, reg := range [][2]uintptr{{1312492, 2}, {1313728, 3}} {
		out = append(out, unsafe.Slice(memmap.PtrUint32(0x5D4594, reg[0]), reg[1])...)
	}
	return out
}
func (e *PortTestScreenEnvironment) Restore() {
	for i := len(e.restore) - 1; i >= 0; i-- {
		e.restore[i]()
	}
}

// This adapter is available in every target; the game-loop wrapper is !server.
func PortTestScreenParticlesDraw(vp *noxrender.Viewport) {
	screenParticlesDraw(vp)
}
