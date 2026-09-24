package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

// Go owns the browser state. Fixed-width ABI types preserve the existing legacy
// callers and fixture word views while their C global symbols retire.
var browserUI = struct {
	listMode           C.uint32_t
	region             int32
	polygonsReady      C.uint32_t
	popup              uint32
	popupCount         C.uint32_t
	connectionState    C.uint
	selected           unsafe.Pointer
	mapWindow          unsafe.Pointer
	overview           C.uint32_t
	filter             C.uint32_t
	label              unsafe.Pointer
	detailPanel        unsafe.Pointer
	detailList         *C.nox_window
	playersColumn      C.uint32_t
	modeColumn         C.uint32_t
	mapColumn          C.uint32_t
	pingColumn         C.uint32_t
	statusColumn       C.uint32_t
	transition         C.uint32_t
	hosting            C.uint32_t
	hasSelection       C.uint32_t
	pendingKicked      C.uint32_t
	pendingTimeout     C.uint32_t
	retry              C.int
	connectionError    C.uint32_t
	creating           C.uint32_t
	animation          *C.nox_gui_animation
	resultCount        C.uint32_t
	sort               C.uint32_t
	gameList           C.uint32_t
	world              *C.nox_window
	connectionDeadline uint64
	refreshDeadline    uint64
	servers            *legacyListNode
}{listMode: 1, region: -1, sort: 6}
