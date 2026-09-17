//go:build porttest

package legacy

/*
#include "GAME1_1.h"
extern uint32_t nox_client_gui_flag_1556112;
extern uint32_t dword_5d4594_825736;
extern int nox_win_width, nox_win_height;
*/
import "C"
import "unsafe"

func PortTestMinimapWords() (map[string]*uint32, func()) {
	m := map[string]*uint32{
		"messageHead":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_825736)),
		"debugIterator": &minimapDebugIterator,
		"zoom":          &minimapZoom,
		"polygons":      &mapPolygonNext,
		"gui":           (*uint32)(unsafe.Pointer(&C.nox_client_gui_flag_1556112)),
		"width":         (*uint32)(unsafe.Pointer(&C.nox_win_width)),
		"height":        (*uint32)(unsafe.Pointer(&C.nox_win_height)),
	}
	old := map[string]uint32{}
	for k, p := range m {
		old[k] = *p
	}
	return m, func() {
		for k, p := range m {
			*p = old[k]
		}
	}
}
func PortTestMinimapNewPolygon() unsafe.Pointer { return unsafe.Pointer(mapPolygonNew()) }
