package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_815748, dword_5d4594_816412;
extern uint32_t nox_color_white_2523948, nox_color_black_2650656;
extern int nox_win_width, nox_win_height;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func clientFrameSample() {
	// Preserve the old platform C.uint boundary before widening the timestamp.
	now := uint64(uint32(PlatformTicks()))
	last := memmap.PtrUint64(0x5D4594, 815740)
	cursor := memmap.PtrUint64(0x5D4594, 815756)
	*memmap.PtrUint64(0x5D4594, 815220+8*uintptr(uint32(*cursor))) = now - *last
	*cursor = (*cursor + 1) % 60
	*last = now
	C.dword_5d4594_815748++
}
func clientFrameAverage() {
	n := min(int32(C.dword_5d4594_815748), 60)
	average := uint64(33)
	if n > 10 {
		var sum uint64
		for i := int32(0); i < n; i++ {
			sum += memmap.Uint64(0x5D4594, 815220+8*uintptr(i))
		}
		average = sum / uint64(n)
	}
	*memmap.PtrUint64(0x587000, 91880) = average
}
func clientModalWindow(command int) {
	if command != 0 {
		win := GetClient().Cli().GUI.NewWindowRaw(nil, 552, 0, 0, int(C.nox_win_width), int(C.nox_win_height), nil)
		C.dword_5d4594_816412 = C.uint32_t(uintptr(unsafe.Pointer(win)))
		win.DrawData().BgColorVal = uint32(C.nox_color_black_2650656)
	} else if C.dword_5d4594_816412 != 0 {
		win := (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_816412)))
		win.Destroy()
		C.dword_5d4594_816412 = 0
	}
}
func clientPlayerColors(pl *server.Player) {
	pl.InitColors()
	pl.Colors.UnkColor = uint32(C.nox_color_white_2523948)
}
func clientAllPlayerColors() {
	players := &GetServer().S().Players
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		clientPlayerColors(pl)
	}
}
func clientAnimationCachesLoad() {
	*memmap.PtrUint32(0x5D4594, 1096456) = uint32(uintptr(Nox_xxx_gLoadAnim("ConfusedBirdies").C()))
	*memmap.PtrUint32(0x5D4594, 1096460) = uint32(uintptr(Nox_xxx_gLoadAnim("SphericalShieldAnim").C()))
}
func clientAnimationCachesClear() {
	*memmap.PtrUint32(0x5D4594, 1096456) = 0
	*memmap.PtrUint32(0x5D4594, 1096460) = 0
}
