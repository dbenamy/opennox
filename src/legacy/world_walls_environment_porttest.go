//go:build porttest

package legacy

/*
#include "defs.h"

*/
import "C"
import "unsafe"

func PortTestWorldWallWords() (map[string]*uint32, func()) {
	m := map[string]*uint32{
		"highFloors":  (*uint32)(unsafe.Pointer(&nox_client_highResFloors_154952)),
		"highFront":   (*uint32)(unsafe.Pointer(&nox_client_highResFrontWalls_80820)),
		"translucent": (*uint32)(unsafe.Pointer(&nox_client_translucentFrontWalls_805844)),
		"edgeMinX":    (*uint32)(unsafe.Pointer(&dword_5d4594_3807140)),
		"edgeMinY":    (*uint32)(unsafe.Pointer(&dword_5d4594_3807136)),
		"edgeMaxX":    (*uint32)(unsafe.Pointer(&dword_5d4594_3807116)),
		"edgeMaxY":    (*uint32)(unsafe.Pointer(&dword_5d4594_3807152)),
		"imageClip":   (*uint32)(unsafe.Pointer(&dword_5d4594_3799452)),
	}
	old := map[string]uint32{}
	for n, p := range m {
		old[n] = *p
	}
	rows := legacyGlobals.nox_pixbuffer_rows_3798784
	return m, func() {
		for n, p := range m {
			*p = old[n]
		}
		legacyGlobals.nox_pixbuffer_rows_3798784 = rows
	}
}
