//go:build porttest

package legacy

import "unsafe"

func PortTestMinimapWords() (map[string]*uint32, func()) {
	m := map[string]*uint32{
		"messageHead":   (*uint32)(unsafe.Pointer(&interactionMessageHead)),
		"debugIterator": &minimapDebugIterator,
		"zoom":          &minimapZoom,
		"polygons":      &mapPolygonNext,
		"gui":           (*uint32)(unsafe.Pointer(&nox_client_gui_flag_1556112)),
		"width":         (*uint32)(unsafe.Pointer(&nox_win_width)),
		"height":        (*uint32)(unsafe.Pointer(&nox_win_height)),
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
