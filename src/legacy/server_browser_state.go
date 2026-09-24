package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

// Go owns the browser state. Fixed-width ABI types preserve the existing legacy
// callers and fixture word views while their C global symbols retire.
var browserUI = struct {
	listMode           uint32
	region             int32
	polygonsReady      uint32
	popup              uint32
	popupCount         uint32
	connectionState    uint32
	selected           unsafe.Pointer
	mapWindow          unsafe.Pointer
	overview           uint32
	filter             uint32
	label              unsafe.Pointer
	detailPanel        unsafe.Pointer
	detailList         *C.nox_window
	playersColumn      uint32
	modeColumn         uint32
	mapColumn          uint32
	pingColumn         uint32
	statusColumn       uint32
	transition         uint32
	hosting            uint32
	hasSelection       uint32
	pendingKicked      uint32
	pendingTimeout     uint32
	retry              int32
	connectionError    uint32
	creating           uint32
	animation          *C.nox_gui_animation
	resultCount        uint32
	sort               uint32
	gameList           uint32
	world              *C.nox_window
	connectionDeadline uint64
	refreshDeadline    uint64
	servers            *legacyListNode
}{listMode: 1, region: -1, sort: 6}
