package legacy

/*
#include "defs.h"
#include "client__draw__partscrn.h"
extern void* nox_alloc_screenParticles_806044;
extern nox_screenParticle* nox_screenParticles_head;
extern nox_screenParticle* dword_5d4594_806052;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"image"
	"unsafe"
)

func screenParticleHead() *Nox_screenParticle {
	return (*Nox_screenParticle)(unsafe.Pointer(C.nox_screenParticles_head))
}
func screenParticleTail() *Nox_screenParticle {
	return (*Nox_screenParticle)(unsafe.Pointer(C.dword_5d4594_806052))
}
func screenParticleAdd(p *Nox_screenParticle) {
	p.Field_44, p.Field_48 = screenParticleHead(), nil
	if p.Field_44 != nil {
		p.Field_44.Field_48 = p
	} else {
		C.dword_5d4594_806052 = (*C.nox_screenParticle)(unsafe.Pointer(p))
	}
	C.nox_screenParticles_head = (*C.nox_screenParticle)(unsafe.Pointer(p))
}
func screenParticleUnlink(p *Nox_screenParticle) {
	if p == screenParticleTail() {
		C.dword_5d4594_806052 = (*C.nox_screenParticle)(unsafe.Pointer(p.Field_48))
	}
	if p.Field_44 != nil {
		p.Field_44.Field_48 = p.Field_48
	}
	if p.Field_48 != nil {
		p.Field_48.Field_44 = p.Field_44
	} else {
		C.nox_screenParticles_head = (*C.nox_screenParticle)(unsafe.Pointer(p.Field_44))
	}
}
func screenParticleDelete(p *Nox_screenParticle) {
	screenParticleUnlink(p)
	alloc.AsClassT[Nox_screenParticle](C.nox_alloc_screenParticles_806044).FreeObjectFirst(p)
}
func screenParticleCreate(kind, x, y, vx, vy, gravity int, size, timer, phase, mode byte) *Nox_screenParticle {
	if C.nox_alloc_screenParticles_806044 == nil || kind < 0 || kind > 4 {
		return nil
	}
	colors := [5][2]uintptr{{806016, 806036}, {806028, 806004}, {806032, 806040}, {806020, 806012}, {806008, 806024}}
	glow, core := *effectMapped(colors[kind][0]), *effectMapped(colors[kind][1])
	p := alloc.AsClassT[Nox_screenParticle](C.nox_alloc_screenParticles_806044).NewObject()
	if p == nil {
		p = screenParticleTail()
		if p == nil {
			return nil
		}
		screenParticleUnlink(p)
	}
	// Tail reuse preserves unused bytes; only the low byte of Field_32 is assigned.
	p.Field_24, p.Field_28 = uint32(x)<<16, uint32(y)<<16
	p.Field_40 = [4]byte{size, timer, phase, timer}
	p.Draw_fnc = C.nox_client_screenParticleDraw_489700
	p.Field_16, p.Field_20, p.Field_36 = uint32(vx)<<16, uint32(vy)<<16, uint32(gravity)<<16
	p.Field_32 = p.Field_32&0xffffff00 | uint32(mode)
	p.Field_4, p.Field_8, p.Field_12 = uint32(kind), glow, core
	screenParticleAdd(p)
	if p.Field_36 == 0 && p.Field_40[1] == 0 {
		p.Field_40[1], p.Field_40[2], p.Field_40[3] = 3, 2, 3
	}
	return p
}
func screenParticleDraw(vp *noxrender.Viewport, p *Nox_screenParticle) int {
	pos := image.Pt(int(p.Field_24>>16), int(p.Field_28>>16))
	if pos.X <= 0 || pos.Y <= 0 || pos.X >= vp.Size.X || pos.Y >= vp.Size.Y {
		screenParticleDelete(p)
		return 0
	}
	effectGlow(pos, p.Field_8, int(p.Field_40[0]), int(p.Field_40[0]))
	effectColor(p.Field_12)
	effectPoint(pos, int(p.Field_40[0]>>1))
	p.Field_20 += p.Field_36
	if p.Field_40[1] != 0 {
		p.Field_40[1]--
		if p.Field_40[1] == 0 {
			if p.Field_40[2] == 1 {
				p.Field_40[0]++
				if p.Field_40[0] >= 4 {
					p.Field_40[2] = 2
				}
			} else {
				p.Field_40[0]--
				if p.Field_40[0] == 0 {
					screenParticleDelete(p)
					return 0
				}
			}
			p.Field_40[1] = p.Field_40[3]
		}
	}
	if byte(p.Field_32) == 1 && effectRand(0, 10) >= 8 {
		timer := byte(effectRand(3, 5))
		size := byte(effectRand(2, 3))
		dx := effectRand(-2, 2)
		screenParticleCreate(int(p.Field_4), pos.X+dx, pos.Y, 0, 0, 1, size, timer, 2, 2)
	}
	p.Field_24 += p.Field_16
	p.Field_28 += p.Field_20
	return 1
}
func screenParticlesDraw(vp *noxrender.Viewport) {
	if vp == nil {
		return
	}
	Sub_430B50(vp.Screen.Min.X, vp.Screen.Min.Y, vp.Screen.Max.X, vp.Screen.Max.Y)
	for p := screenParticleHead(); p != nil; {
		GetClient().Cli().GUI.ValYYY = 1
		// The callback can delete, recycle, or create nodes, including the current one.
		next := p.Field_44
		ccall.CallIntPtr2(p.Draw_fnc, vp.C(), unsafe.Pointer(p))
		p = next
	}
}
