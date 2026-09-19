//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_1305680,dword_5d4594_1305684;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

func PortTestSessionDialogWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"otherDialogA":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305680)),
		"otherDialogB":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305684)),
		"filter":         (*uint32)(unsafe.Pointer(&sessionFilterRoot)),
		"filterControls": (*uint32)(unsafe.Pointer(&sessionFilterControls)),
		"disconnect":     (*uint32)(unsafe.Pointer(&sessionDisconnectRoot)),
		"disconnectIcon": (*uint32)(unsafe.Pointer(&sessionDisconnectIcon)),
		"motd":           (*uint32)(unsafe.Pointer(&sessionMOTDRoot)),
		"motdList":       (*uint32)(unsafe.Pointer(&sessionMOTDList)),
		"motdFile":       (*uint32)(unsafe.Pointer(&sessionMOTDFile)),
		"quit":           (*uint32)(unsafe.Pointer(&sessionQuitRoot)),
	}
	old := make(map[string]uint32)
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

func PortTestSessionDialogCall(op string, p unsafe.Pointer, code int, a, b uint32) uintptr {
	w := (*gui.Window)(p)
	event := func(fn gui.WindowFunc) uintptr {
		return uintptr(gui.EventRespInt(fn(w, gui.AsWindowEvent(code, uintptr(a), uintptr(b)))))
	}
	switch op {
	case "filterOpen":
		return uintptr(sessionFilterOpen(w).C())
	case "filterSave":
		return uintptr(sessionFilterSave())
	case "filterControls":
		sessionFilterUpdateControls()
		return 0
	case "filterClose":
		return uintptr(sessionFilterClose())
	case "filterConfig":
		return uintptr(sessionFilterConfig(int(a), int(b), p))
	case "filterEvent":
		return event(sessionFilterEvent)
	case "motdOpen":
		return uintptr(sessionMOTDOpen().C())
	case "motdClose":
		return uintptr(sessionMOTDClose())
	case "motdShown":
		return uintptr(sessionMOTDShown())
	case "motdShow":
		sessionMOTDShow()
		return 0
	case "motdEvent":
		return event(sessionMOTDEvent)
	case "motdAdd":
		return sessionMOTDAdd(p)
	case "motdRead":
		return uintptr(sessionMOTDRead(int(a)))
	case "motdFree":
		return sessionMOTDFree(int(a))
	case "disconnectOpen":
		return uintptr(sessionDisconnectOpen())
	case "disconnectClose":
		return uintptr(sessionDisconnectClose())
	case "disconnectIcon":
		return uintptr(sessionDisconnectIconShow(int(a)))
	case "disconnectShow":
		return uintptr(sessionDisconnectShow(int(a)))
	case "disconnectInput":
		return event(sessionDisconnectInput)
	case "disconnectEvent":
		return event(sessionDisconnectEvent)
	case "disconnectDraw":
		return uintptr(sessionDisconnectDraw(w, w.DrawData()))
	case "quitShown":
		return uintptr(sessionQuitShown())
	case "quitColors":
		return uintptr(sessionQuitColors().C())
	case "quitCapture":
		return uintptr(sessionQuitCapture())
	case "quitHide":
		sessionQuitHide()
		return 0
	case "quitToggle":
		sessionQuitToggle()
		return 0
	case "quitEvent":
		return event(SessionQuitEvent)
	}
	panic(op)
}
