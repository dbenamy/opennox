//go:build porttest

package legacy

/*
#include "GAME5.h"
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"unsafe"
)

// Install only the production callback table/name region, with real ABI addresses.
func portTestCallbackTablesEnvironment() func() {
	dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 287096)), 640)
	old := bytes.Clone(dst)
	copy(dst, blobdata.PortTestCallbackTables())
	*memmap.PtrPtr(0x587000, 287096) = memmap.PtrOff(0x587000, 287352)
	*memmap.PtrPtr(0x587000, 287100) = unsafe.Pointer(C.nox_xxx_strikeOgre_549220)
	*memmap.PtrPtr(0x587000, 287104) = memmap.PtrOff(0x587000, 287364)
	*memmap.PtrPtr(0x587000, 287108) = unsafe.Pointer(C.nox_xxx_strikeScorpion_5495B0)
	*memmap.PtrPtr(0x587000, 287112) = memmap.PtrOff(0x587000, 287380)
	*memmap.PtrPtr(0x587000, 287116) = unsafe.Pointer(C.nox_xxx_strikeVileZombie_549700)
	*memmap.PtrPtr(0x587000, 287120) = memmap.PtrOff(0x587000, 287400)
	*memmap.PtrPtr(0x587000, 287124) = unsafe.Pointer(C.nox_xxx_strikeStoneGolem_5497E0)
	*memmap.PtrPtr(0x587000, 287128) = memmap.PtrOff(0x587000, 287420)
	*memmap.PtrPtr(0x587000, 287132) = unsafe.Pointer(C.nox_xxx_strikeMechGolem_549960)
	*memmap.PtrPtr(0x587000, 287136) = memmap.PtrOff(0x587000, 287436)
	*memmap.PtrPtr(0x587000, 287140) = unsafe.Pointer(C.nox_xxx_strikeWasp_549980)
	*memmap.PtrPtr(0x587000, 287144) = memmap.PtrOff(0x587000, 287448)
	*memmap.PtrPtr(0x587000, 287148) = unsafe.Pointer(C.nox_xxx_strikeSpider_549BC0)
	*memmap.PtrPtr(0x587000, 287152) = memmap.PtrOff(0x587000, 287464)
	*memmap.PtrPtr(0x587000, 287156) = unsafe.Pointer(C.nox_xxx_strikeSpittingSpider_549CA0)
	*memmap.PtrPtr(0x587000, 287160) = memmap.PtrOff(0x587000, 287488)
	*memmap.PtrPtr(0x587000, 287164) = unsafe.Pointer(C.nox_xxx_strikeGhost_549A60)
	*memmap.PtrPtr(0x587000, 287168) = memmap.PtrOff(0x587000, 287500)
	*memmap.PtrPtr(0x587000, 287172) = unsafe.Pointer(C.nox_xxx_strikeBomber_549BB0)
	*memmap.PtrPtr(0x587000, 287176) = memmap.PtrOff(0x587000, 287516)
	*memmap.PtrPtr(0x587000, 287180) = unsafe.Pointer(C.nox_xxx_strikeMonsterDefault_549380)
	*memmap.PtrPtr(0x587000, 287192) = memmap.PtrOff(0x587000, 287532)
	*memmap.PtrPtr(0x587000, 287196) = unsafe.Pointer(C.sub_549D80)
	*memmap.PtrPtr(0x587000, 287200) = memmap.PtrOff(0x587000, 287548)
	*memmap.PtrPtr(0x587000, 287204) = unsafe.Pointer(C.sub_549E00)
	*memmap.PtrPtr(0x587000, 287208) = memmap.PtrOff(0x587000, 287560)
	*memmap.PtrPtr(0x587000, 287212) = unsafe.Pointer(C.sub_549E70)
	*memmap.PtrPtr(0x587000, 287216) = memmap.PtrOff(0x587000, 287568)
	*memmap.PtrPtr(0x587000, 287220) = unsafe.Pointer(C.sub_549E90)
	*memmap.PtrPtr(0x587000, 287224) = memmap.PtrOff(0x587000, 287584)
	*memmap.PtrPtr(0x587000, 287228) = unsafe.Pointer(C.sub_549FA0)
	*memmap.PtrPtr(0x587000, 287232) = memmap.PtrOff(0x587000, 287596)
	*memmap.PtrPtr(0x587000, 287236) = unsafe.Pointer(C.nox_bomberDead_54A150)
	*memmap.PtrPtr(0x587000, 287240) = memmap.PtrOff(0x587000, 287608)
	*memmap.PtrPtr(0x587000, 287244) = unsafe.Pointer(C.sub_54A250)
	*memmap.PtrPtr(0x587000, 287248) = memmap.PtrOff(0x587000, 287620)
	*memmap.PtrPtr(0x587000, 287252) = unsafe.Pointer(C.sub_54A310)
	*memmap.PtrPtr(0x587000, 287256) = memmap.PtrOff(0x587000, 287636)
	*memmap.PtrPtr(0x587000, 287260) = unsafe.Pointer(C.sub_54A750)
	*memmap.PtrPtr(0x587000, 287264) = memmap.PtrOff(0x587000, 287656)
	*memmap.PtrPtr(0x587000, 287268) = unsafe.Pointer(C.nox_xxx_monsterDeadTroll_54A270)
	*memmap.PtrPtr(0x587000, 287280) = memmap.PtrOff(0x587000, 287668)
	*memmap.PtrPtr(0x587000, 287284) = unsafe.Pointer(C.sub_54A890)
	*memmap.PtrPtr(0x587000, 287288) = memmap.PtrOff(0x587000, 287680)
	*memmap.PtrPtr(0x587000, 287292) = unsafe.Pointer(C.sub_54A900)
	*memmap.PtrPtr(0x587000, 287296) = memmap.PtrOff(0x587000, 287688)
	*memmap.PtrPtr(0x587000, 287300) = unsafe.Pointer(C.sub_54A7D0)
	*memmap.PtrPtr(0x587000, 287304) = memmap.PtrOff(0x587000, 287704)
	*memmap.PtrPtr(0x587000, 287308) = unsafe.Pointer(C.sub_54A850)
	*memmap.PtrPtr(0x587000, 287312) = memmap.PtrOff(0x587000, 287720)
	*memmap.PtrPtr(0x587000, 287316) = unsafe.Pointer(C.sub_54A950)
	return func() { copy(dst, old) }
}
