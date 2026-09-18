//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_3.h"
#include "client__gui__guivote.h"
extern uint32_t dword_5d4594_1197308;
extern uint32_t dword_5d4594_1197312;
extern uint32_t dword_5d4594_1197316;
extern uint32_t dword_5d4594_1197320;
extern uint32_t dword_5d4594_1197324;
extern uint32_t dword_5d4594_1197328;
extern uint32_t dword_5d4594_1197332;
extern uint32_t dword_5d4594_1197336;
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestVoteGUIOwner() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"topic": (*uint32)(&C.dword_5d4594_1197308), "window": (*uint32)(&C.dword_5d4594_1197312),
		"players": (*uint32)(&C.dword_5d4594_1197316), "topics": (*uint32)(&C.dword_5d4594_1197320),
		"count": (*uint32)(&C.dword_5d4594_1197324), "previousCount": (*uint32)(&C.dword_5d4594_1197328),
		"choice": (*uint32)(&C.dword_5d4594_1197332), "previousChoice": (*uint32)(&C.dword_5d4594_1197336),
	}
	old := map[string]uint32{}
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	var regions, saved [][]byte
	for _, off := range []int{1193720, 1195512} {
		b := memmap.BlobByAddr(0x5D4594).Data[off : off+1792]
		regions = append(regions, b)
		saved = append(saved, bytes.Clone(b))
		clear(b)
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
		for i, b := range regions {
			copy(b, saved[i])
		}
	}
}
func PortTestVoteGUI(op string, topic int) int {
	switch op {
	case "init":
		return int(C.sub_48D000_initGuiKick())
	case "show":
		C.sub_48CB10(C.int(topic))
	case "hide":
		return int(C.sub_48CAD0())
	case "selection":
		return int(C.sub_48D120())
	case "choice":
		return int(C.sub_48D340())
	case "topic":
		C.sub_48D410()
	default:
		panic(op)
	}
	return 0
}
func PortTestVoteGUIWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1197312)))
}
