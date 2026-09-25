package legacy

/*
#include "GAME2.h"
#include "GAME3_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"math"
	"unsafe"
)

func drawableStatePredicate(dr *client.Drawable) int {
	return bool2int(dr.ObjClass&0x400000 != 0 && dr.ObjSubClass&8 != 0)
}

//export sub_45A010
func sub_45A010(dr *nox_drawable) *nox_drawable { return (*nox_drawable)(asDrawable(dr).Field_104.C()) }

//export nox_drawable_next_45A070
func nox_drawable_next_45A070(dr *nox_drawable) *nox_drawable {
	if dr == nil {
		return nil
	}
	return (*nox_drawable)(asDrawable(dr).NextPtr.C())
}

func nox_xxx_spriteSetActiveMB_45A990_drawable(p C.int) C.int {
	(*client.Drawable)(unsafe.Pointer(uintptr(uint32(p)))).SetActive()
	return p
}

func nox_xxx_spriteSetFrameMB_45AB80(p, value C.int) C.int {
	(*client.Drawable)(unsafe.Pointer(uintptr(uint32(p)))).SetFrameMB(int(value))
	return p
}
func drawableMotionDamped(vp *noxrender.Viewport, dr *client.Drawable) int {
	vx, vy := float64(math.Float32frombits(dr.Field_117)), float64(math.Float32frombits(dr.Field_118))
	drag := float64(math.Float32frombits(dr.Field_119))
	var dx, dy float64
	n := GetServer().S().Frame() - dr.AnimStart + 1
	for {
		vx -= vx * drag
		vy += float64(float32(-(vy * drag)))
		dx += vx
		dy += vy
		n--
		if n == 0 {
			break
		}
	}
	// Match C's float stores before integer conversion, including the extra Y store.
	x := int(floatToInt32(float32(float64(int32(dr.Field_81)) + dx)))
	y := int(floatToInt32(float32(float64(int32(dr.Field_82)) + float64(float32(dy)))))
	if x > 0 && y > 0 && x < 5888 && y < 5888 {
		GetClient().Cli().Nox_xxx_updateSpritePosition_49AA90(dr, x, y)
		if GetClient().Cli().Sight.Sub_4992B0(vp.Screen.Min.X+dr.PosVec.X-vp.World.Min.X, dr.PosVec.Y+vp.Screen.Min.Y-vp.World.Min.Y) != 0 {
			return 1
		}
	}
	GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	return 0
}
func drawableMotionTarget(dr *client.Drawable) int {
	tx := int32(*(*uint16)(unsafe.Add(dr.C(), 432)))
	ty := int32(*(*uint16)(unsafe.Add(dr.C(), 434)))
	x, y := int32(dr.PosVec.X), int32(dr.PosVec.Y)
	dx, dy := tx-x, ty-y
	speed := int32(*(*byte)(unsafe.Add(dr.C(), 443)))
	distance := int32(screenDistance(dx, dy)) + 1
	nx, ny := x+dx*speed/distance, y+dy*speed/distance
	if distance <= 10 || (x-tx)*(nx-tx)+(y-ty)*(ny-ty) < 0 {
		GetClient().Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	GetClient().Cli().Nox_xxx_updateSpritePosition_49AA90(dr, int(nx), int(ny))
	return 1
}
