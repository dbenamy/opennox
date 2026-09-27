package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

var bookTooltipKeys [5]byte
var bookImageEndKeys [2]byte

func bookTooltipKey(i int) unsafe.Pointer  { return unsafe.Pointer(&bookTooltipKeys[i]) }
func bookImageEndKey(i int) unsafe.Pointer { return unsafe.Pointer(&bookImageEndKeys[i]) }

func init() {
	gui.RegisterTooltipCallbackGo(bookTooltipKey(0), func(_ *gui.Window, _ *gui.WindowData, _ uintptr) { quickbarBookTooltip() })
	gui.RegisterTooltipCallbackGo(bookTooltipKey(1), func(w *gui.Window, _ *gui.WindowData, _ uintptr) { quickbarDirectionTooltip(w) })
	gui.RegisterTooltipCallbackGo(bookTooltipKey(2), func(w *gui.Window, _ *gui.WindowData, arg uintptr) { summonSlotTooltipCallback(w, uint32(arg)) })
	gui.RegisterTooltipCallbackGo(bookTooltipKey(3), func(_ *gui.Window, _ *gui.WindowData, _ uintptr) { summonCommandTooltip() })
	gui.RegisterTooltipCallbackGo(bookTooltipKey(4), func(w *gui.Window, _ *gui.WindowData, _ uintptr) { bookIconTooltip(w) })
}

func summonSlotTooltipCallback(w *gui.Window, arg uint32) int {
	Nox_xxx_cursorSetTooltip_4776B0(GoWStringP(summonSlotTooltip(w, bookPoint(arg))))
	return 1
}

// CallImageAnimationEnd invokes the completion callback with its image reference.
// Book callbacks ignore the argument. Other addresses retain the foreign ABI.
func CallImageAnimationEnd(key unsafe.Pointer, ref *ImageRef) {
	switch key {
	case bookImageEndKey(0):
		bookPageComplete(false)
	case bookImageEndKey(1):
		bookPageComplete(true)
	default:
		ccall.CallVoidPtr(key, ref.C())
	}
}
