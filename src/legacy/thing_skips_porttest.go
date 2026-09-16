//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME2.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestThingSkip(op int, f *binfile.MemFile, scratch []byte) int {
	p := (*C.nox_memfile)(f.C())
	switch op {
	case 0:
		return int(C.nox_thing_skip_AUD_414D40(p))
	case 1:
		return int(C.nox_thing_skip_spells_415100(p))
	case 2:
		return int(C.nox_thing_read_ability_415320(p))
	case 3:
		return int(C.nox_thing_read_image_415240(p))
	case 4:
		return int(C.nox_thing_skip_AVNT_452B00(p))
	case 5:
		return int(C.nox_thing_skip_AVNT_inner_452B30(p))
	case 6:
		return int(C.nox_thing_read_WALL_414F60(p, unsafe.Pointer(&scratch[0])))
	default:
		panic(op)
	}
}
