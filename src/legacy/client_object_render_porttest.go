//go:build porttest

package legacy

/*
#include "defs.h"
extern uint32_t dword_5d4594_1321520;
extern uint32_t dword_5d4594_1321800;
extern uint32_t dword_5d4594_1305748;
extern uint32_t dword_8531A0_2576;
extern unsigned int nox_player_netCode_85319C;
extern nox_render_data_t* nox_draw_curDrawData_3799572;
extern int nox_win_height;
int sub_4C5020(int);
void sub_4C5050();
void nox_xxx_wndDraw_49F7F0();
int sub_49F860();
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"unsafe"
)

type PortTestObjectRenderEnvironment struct{ restore []func() }

func PortTestNewObjectRenderEnvironment(data *noxrender.RenderData) *PortTestObjectRenderEnvironment {
	e := new(PortTestObjectRenderEnvironment)
	ghost, count, saved, player, netcode := C.dword_5d4594_1321520, C.dword_5d4594_1321800, C.dword_5d4594_1305748, C.dword_8531A0_2576, C.nox_player_netCode_85319C
	height, render := C.nox_win_height, C.nox_draw_curDrawData_3799572
	e.restore = append(e.restore, func() {
		C.dword_5d4594_1321520, C.dword_5d4594_1321800, C.dword_5d4594_1305748, C.dword_8531A0_2576, C.nox_player_netCode_85319C = ghost, count, saved, player, netcode
		C.nox_win_height, C.nox_draw_curDrawData_3799572 = height, render
	})
	C.nox_win_height = 96
	C.nox_draw_curDrawData_3799572 = (*C.nox_render_data_t)(data.C())
	for _, reg := range [][3]uintptr{{0x587000, 80808, 4}, {0x5D4594, 1305732, 44}, {0x5D4594, 1321512, 16}, {0x5D4594, 1321532, 268}, {0x85B3FC, 956, 4}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(reg[0], reg[1])), reg[2])
		old := append([]byte(nil), b...)
		clear(b)
		e.restore = append(e.restore, func() { copy(b, old) })
	}
	for _, reg := range blobdata.PortTestObjectRenderTables() {
		b := unsafe.Slice((*byte)(memmap.PtrOff(reg.Base, reg.Offset)), len(reg.Data))
		old := append([]byte(nil), b...)
		copy(b, reg.Data)
		e.restore = append(e.restore, func() { copy(b, old) })
	}
	e.Reset()
	return e
}
func (e *PortTestObjectRenderEnvironment) Reset() {
	C.dword_5d4594_1321520 = 0x7fffffff
	C.dword_5d4594_1321800 = 0
	C.dword_5d4594_1305748 = 0
	C.dword_8531A0_2576 = 0
	C.nox_player_netCode_85319C = 7
	for _, reg := range [][3]uintptr{{0x587000, 80808, 4}, {0x5D4594, 1305732, 44}, {0x5D4594, 1321512, 16}, {0x5D4594, 1321532, 268}} {
		clear(unsafe.Slice((*byte)(memmap.PtrOff(reg[0], reg[1])), reg[2]))
	}
}
func (e *PortTestObjectRenderEnvironment) Restore() {
	for i := len(e.restore) - 1; i >= 0; i-- {
		e.restore[i]()
	}
}
func (e *PortTestObjectRenderEnvironment) GhostType(v uint32) { C.dword_5d4594_1321520 = C.uint32_t(v) }
func (e *PortTestObjectRenderEnvironment) State() []uint32 {
	out := []uint32{uint32(C.dword_5d4594_1321520), uint32(C.dword_5d4594_1321800), uint32(C.dword_5d4594_1305748), uint32(C.nox_player_netCode_85319C), *memmap.PtrUint32(0x5D4594, 1321512)}
	for _, reg := range [][2]uintptr{{1305732, 11}, {1321532, 67}} {
		out = append(out, unsafe.Slice(memmap.PtrUint32(0x5D4594, reg[0]), reg[1])...)
	}
	return out
}
func (e *PortTestObjectRenderEnvironment) LastDrawable() *client.Drawable {
	return (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1321516))
}
func PortTestObjectRenderGhost(vp *noxrender.Viewport, dr *client.Drawable) byte {
	return objectRenderGhost(vp, dr)
}
func PortTestObjectRenderShiny(vp *noxrender.Viewport, dr *client.Drawable) uint16 {
	return objectRenderShiny(vp, dr)
}
func PortTestObjectRenderClip(save bool) int32 {
	if save {
		C.nox_xxx_wndDraw_49F7F0()
		return 0
	}
	return int32(C.sub_49F860())
}

func PortTestObjectRenderBeam(op int, vp *noxrender.Viewport, a [4]int32) uint32 {
	switch op {
	case 0:
		return uint32(objectRenderBeamColors())
	case 1:
		// Packet layout is one leading opcode byte followed by two unaligned words.
		packet, free := alloc.New([9]byte{})
		defer free()
		*(*uint16)(unsafe.Pointer(&packet[1])) = uint16(a[0])
		*(*uint16)(unsafe.Pointer(&packet[3])) = uint16(a[1])
		*(*uint16)(unsafe.Pointer(&packet[5])) = uint16(a[2])
		*(*uint16)(unsafe.Pointer(&packet[7])) = uint16(a[3])
		return uint32(C.sub_4C5020(C.int(uintptr(unsafe.Pointer(&packet[0])))))
	case 2:
		C.sub_4C5050()
		return 0
	case 3:
		return uint32(objectRenderBeamDraw(vp))
	case 4:
		return uint32(objectRenderBeamLine(image.Pt(int(a[0]), int(a[1])), image.Pt(int(a[2]), int(a[3]))))
	}
	panic("unknown beam operation")
}
