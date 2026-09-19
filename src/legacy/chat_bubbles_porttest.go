//go:build porttest

package legacy

/*
#include "GAME2_3.h"
void nox_xxx_createTextBubble_48D880(void* a1, wchar2_t* a2);
extern uint32_t dword_5d4594_1197372;
extern int nox_win_width, nox_win_height;
extern uint32_t nox_color_white_2523948;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestChatBubbleGlobals() (head, tail *unsafe.Pointer, restore func()) {
	head = (*unsafe.Pointer)(memmap.PtrOff(0x5D4594, 1197368))
	tail = (*unsafe.Pointer)(unsafe.Pointer(&C.dword_5d4594_1197372))
	oldHead, oldTail := *head, *tail
	*head, *tail = nil, nil
	return head, tail, func() { *head, *tail = oldHead, oldTail }
}
func PortTestChatBubbleCreate(data unsafe.Pointer, text *uint16) {
	C.nox_xxx_createTextBubble_48D880(data, (*C.wchar2_t)(unsafe.Pointer(text)))
}
func PortTestChatBubbleLookup(code uint32) unsafe.Pointer {
	return unsafe.Pointer(uintptr(uint32(C.nox_xxx_netCode2ChatBubble_48D850(C.int(code)))))
}
func PortTestChatBubbleRemove(code uint32) { C.sub_48E8E0(C.int(code)) }
func PortTestChatBubbleClear()             { C.sub_48E940() }
func PortTestChatBubbleDestroy()           { C.sub_48D800() }

func PortTestChatBubbleLayout(v *noxrender.Viewport) { C.sub_48DCF0((*C.uint32_t)(v.C())) }
func PortTestChatBubbleDraw(v *noxrender.Viewport)   { C.sub_48D990((*C.nox_draw_viewport_t)(v.C())) }
func PortTestChatBubblePlace(rect *[4]int32, tail *uint32) {
	C.sub_48E000((*C.int4)(unsafe.Pointer(rect)), (*C.uint32_t)(tail))
}
func PortTestChatBubbleOverlap(a, b unsafe.Pointer) int32 {
	return int32(C.sub_48E480((*C.uint32_t)(a), (*C.uint32_t)(b)))
}
func PortTestChatBubbleArrange(v *noxrender.Viewport, p unsafe.Pointer) {
	C.sub_48E240(C.int(uintptr(v.C())), (*C.uint32_t)(p))
}

func PortTestChatBubbleRegion(x, y int32) int32 { return int32(C.sub_48E530(C.int(x), C.int(y))) }
func PortTestChatBubbleCandidate(p unsafe.Pointer, x, y int32) bool {
	return C.sub_48E5C0((*C.uint32_t)(p), C.int(x), C.int(y)) != 0
}
func PortTestChatBubbleShift(mask byte, a, b unsafe.Pointer) [2]int32 {
	var pos [2]int32
	C.sub_48E6A0(C.char(mask), (*C.uint32_t)(a), (*C.uint32_t)(b), (*C.int)(unsafe.Pointer(&pos[0])), (*C.int)(unsafe.Pointer(&pos[1])))
	return pos
}

func PortTestChatBubbleRenderGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"width": (*uint32)(unsafe.Pointer(&C.nox_win_width)), "height": (*uint32)(unsafe.Pointer(&C.nox_win_height)), "white": (*uint32)(unsafe.Pointer(&C.nox_color_white_2523948))}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
