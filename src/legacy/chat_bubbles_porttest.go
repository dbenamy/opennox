//go:build porttest

package legacy

/*
#include "GAME2_3.h"
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
	tail = (*unsafe.Pointer)(unsafe.Pointer(&chatBubbleTail))
	oldHead, oldTail := *head, *tail
	*head, *tail = nil, nil
	return head, tail, func() { *head, *tail = oldHead, oldTail }
}
func PortTestChatBubbleCreate(data unsafe.Pointer, text *uint16) {
	chatBubbleCreate(data, text)
}
func PortTestChatBubbleLookup(code uint32) unsafe.Pointer {
	return unsafe.Pointer(chatBubbleLookup(code))
}
func PortTestChatBubbleRemove(code uint32) { C.sub_48E8E0(C.int(code)) }
func PortTestChatBubbleClear()             { C.sub_48E940() }
func PortTestChatBubbleDestroy()           { chatBubbleDestroy() }

func PortTestChatBubbleLayout(v *noxrender.Viewport) { chatBubbleLayout(v) }
func PortTestChatBubbleDraw(v *noxrender.Viewport)   { chatBubbleDraw(v) }
func PortTestChatBubblePlace(rect *[4]int32, tail *uint32) {
	chatBubblePlace(rect, tail)
}
func PortTestChatBubbleOverlap(a, b unsafe.Pointer) int32 {
	return int32(bool2int(chatBubbleOverlap((*chatBubble)(a), (*chatBubble)(b))))
}
func PortTestChatBubbleArrange(v *noxrender.Viewport, p unsafe.Pointer) {
	chatBubbleArrange((*chatBubble)(p))
}

func PortTestChatBubbleRegion(x, y int32) int32 { return chatBubbleRegion(x, y) }
func PortTestChatBubbleCandidate(p unsafe.Pointer, x, y int32) bool {
	return chatBubbleCandidate((*chatBubble)(p), x, y)
}
func PortTestChatBubbleShift(mask byte, a, b unsafe.Pointer) [2]int32 {
	var pos [2]int32
	pos[0], pos[1] = chatBubbleShift(mask, (*chatBubble)(a), (*chatBubble)(b))
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
