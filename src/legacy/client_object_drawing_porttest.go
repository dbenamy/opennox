//go:build porttest

package legacy

/*
#include "GAME3_1.h"
#include "GAME1_1.h"
#include "client__draw__doordraw.h"
#include "client__draw__arrowdraw.h"
#include "client__draw__glyphdraw.h"
#include "client__draw__summondraw.h"
#include "client__draw__weapondraw.h"
#include "client__draw__armordraw.h"
#include "client__draw__mgendraw.h"
#include "client__draw__pressureplatedraw.h"
#include "client__draw__triggerdraw.h"
#include "client__draw__powderdraw.h"
#include "client__draw__basedraw.h"
#include "client__draw__flagdraw.h"
*/
import "C"

import (
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestObjectDrawHelper(op int, dr *client.Drawable, arg unsafe.Pointer, value int) uint32 {
	switch op {
	case 0:
		return uint32(C.sub_4B9470((**C.char)(arg)))
	case 1:
		return uint32(objectDrawableTeamColor(dr))
	case 3:
		return uint32(uintptr(unsafe.Pointer(C.sub_4B9650(C.int(value)))))
	case 5:
		return uint32(sub_4BC720(C.int(uintptr(dr.C()))))
	case 6:
		return uint32(C.nox_xxx_updDrawMonsterGen_4BC920())
	case 7:
		return uint32(uintptr(unsafe.Pointer(objectTeamByColor(value))))
	}
	panic("unknown object drawing helper")
}

func PortTestObjectDoorParse(obj *client.ObjectType, f *binfile.MemFile, attr unsafe.Pointer) bool {
	return bool(C.nox_things_door_draw_parse((*C.nox_thing)(obj.C()), (*C.nox_memfile)(f.C()), (*C.char)(attr)))
}

func PortTestObjectDrawCallback(op int) unsafe.Pointer {
	switch op {
	case 0:
		return C.nox_thing_door_draw
	case 1:
		return C.nox_thing_arrow_draw
	case 2:
		return C.nox_thing_weak_arrow_draw
	case 3:
		return C.nox_thing_arrow_tail_link_draw
	case 4:
		return C.nox_thing_weak_arrow_tail_link_draw
	case 5:
		return C.nox_thing_glyph_draw
	case 6:
		return C.nox_thing_summon_effect_draw
	case 7:
		return C.nox_thing_weapon_draw
	case 8:
		return C.nox_thing_weapon_animate_draw
	case 9:
		return C.nox_thing_armor_draw
	case 10:
		return C.nox_thing_armor_animate_draw
	case 11:
		return C.nox_thing_spherical_shield_draw
	case 12:
		return C.nox_thing_monster_gen_draw
	case 13:
		return C.nox_thing_pressure_plate_draw
	case 14:
		return C.nox_thing_trigger_draw
	case 15:
		return C.nox_thing_black_powder_draw
	case 16:
		return C.nox_thing_base_draw
	case 17:
		return C.nox_thing_flag_draw
	}
	panic("unknown object drawing callback")
}

type PortTestObjectDrawEnvironment struct{ restore []func() }

func PortTestNewObjectDrawEnvironment() *PortTestObjectDrawEnvironment {
	e := new(PortTestObjectDrawEnvironment)
	black, blue, text, summon := nox_color_black_2650656, nox_color_blue_2650684, dword_8531A0_2572, drawableSummonSpark
	e.restore = append(e.restore, func() {
		nox_color_black_2650656, nox_color_blue_2650684, dword_8531A0_2572, drawableSummonSpark = black, blue, text, summon
	})
	for _, region := range [][3]uintptr{{0x587000, 177488, 64}, {0x5D4594, 1313720, 8}, {0x852978, 8, 4}} {
		buf := unsafe.Slice((*byte)(memmap.PtrOff(region[0], region[1])), int(region[2]))
		old := append([]byte(nil), buf...)
		clear(buf)
		e.restore = append(e.restore, func() { copy(buf, old) })
	}
	e.Reset()
	return e
}
func (e *PortTestObjectDrawEnvironment) Reset() {
	nox_color_black_2650656 = C.uint32_t(noxcolor.RGB5551Color(0, 0, 0).Color32())
	nox_color_blue_2650684 = C.uint32_t(noxcolor.RGB5551Color(0, 0, 255).Color32())
	dword_8531A0_2572 = C.uint32_t(noxcolor.RGB5551Color(220, 220, 60).Color32())
	drawableSummonSpark = 0
	*memmap.PtrUint32(0x5D4594, 1313720) = 0
	*memmap.PtrUint32(0x5D4594, 1313724) = 0
	*memmap.PtrUint32(0x852978, 8) = 0
	*memmap.PtrUint32(0x5D4594, 1321512) = 0
	*memmap.PtrUint32(0x5D4594, 1321516) = 0
}
func (e *PortTestObjectDrawEnvironment) Restore() {
	for i := len(e.restore) - 1; i >= 0; i-- {
		e.restore[i]()
	}
}

// PortTestObjectDrawNamedState observes the actual C globals, including lazy cache.
func (e *PortTestObjectDrawEnvironment) NamedState() []uint32 {
	return []uint32{uint32(nox_color_black_2650656), uint32(nox_color_blue_2650684), uint32(dword_8531A0_2572), uint32(drawableSummonSpark)}
}
