package legacy

/*
#include "GAME1.h"
#include "GAME4_1.h"
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

var (
	Sub_526CA0                               func(a1 string) int
	Nox_xxx_mapSetWallInGlobalDir0pr1_5004D0 func()
	Nox_xxx_map_5004F0                       func()
	Sub_4FF990                               func(a1 uint32)
	Sub_5000B0                               func(a1 *server.Object) int
)

var _ = [1]struct{}{}[12332-unsafe.Sizeof(server.WallDef{})]

//export nox_xxx_wallFlags
func nox_xxx_wallFlags(ind int) uint32 {
	return GetServer().S().Walls.DefByInd(ind).Flags32
}

func nox_xxx_mapDamageToWalls_534FC0(a1 *C.int4, a2 unsafe.Pointer, a3 C.float, a4, a5 int, a6 unsafe.Pointer) C.bool {
	rect := image.Rect(int(a1.field_0), int(a1.field_4), int(a1.field_8), int(a1.field_C))
	return C.bool(GetServer().Nox_xxx_mapDamageToWalls_534FC0(rect, *(*types.Pointf)(a2), float32(a3), a4, object.DamageType(a5), AsObjectP(a6)))
}

func Sub_5071C0() bool {
	return voteHead != nil
}

func Nox_xxx_math_509ED0(pos types.Pointf) int {
	cpos, free := alloc.New(types.Pointf{})
	defer free()
	*cpos = pos
	return int(C.int(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(cpos))))))
}

func Nox_xxx_math_509EA0(a1 int) int {
	return int(C.int(geometryDirection4Index(int32(a1))))
}
