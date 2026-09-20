package legacy

/*
#include "defs.h"
extern int nox_win_width, nox_win_height;
extern uint32_t dword_5d4594_1090276;
extern uint32_t dword_5d4594_1090280;
extern uint32_t dword_5d4594_1090284;
extern uint32_t dword_5d4594_1090292;
extern uint32_t dword_5d4594_1090828;
extern uint32_t dword_5d4594_1091364;
extern uint32_t dword_5d4594_1096252;
extern uint32_t dword_5d4594_1096256;
extern uint32_t dword_5d4594_1096260;
extern uint32_t dword_5d4594_1096264;
extern uint32_t dword_5d4594_1096272;
extern uint32_t dword_5d4594_1096276;
extern uint32_t dword_5d4594_1096280;
extern uint32_t dword_5d4594_1096284;
extern uint32_t dword_5d4594_1096288;
extern uint32_t dword_8531A0_2576;
extern uint32_t nox_client_renderBubbles_80844;
extern uint32_t nox_color_black_2650656;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_yellow_2589772;
extern unsigned int nox_gameDisableMapDraw_5d4594_2650672;
extern unsigned int nox_player_netCode_85319C;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

type uiMeterRecord struct {
	Window                             *gui.Window
	Current, Maximum, Color, Alternate uint32
}

var uiMeterRecords [7]uiMeterRecord

func uiMeters() *[7]uiMeterRecord { return &uiMeterRecords }
func uiMeterHide(w *gui.Window, hidden bool) int {
	if w == nil {
		return -2
	}
	w.SetHidden(hidden)
	return 0
}
func uiMeterPlayer() unsafe.Pointer   { return unsafe.Pointer(uintptr(C.dword_8531A0_2576)) }
func uiMeterColor(r, g, b int) uint32 { return uint32(nox_color_rgb_4344A0(r, g, b)) }
func uiMeterMode() uint32             { return uint32(C.dword_5d4594_1096252) }

func nox_xxx_cliShowHideTubes_470AA0(v int) {
	C.dword_5d4594_1096252 = C.uint32_t(v)
	if memmap.Uint32(0x5D4594, 1093176) != 0 {
		uiMeterHide(uiMeters()[2].Window, v == 0)
		uiMeterHide(uiMeters()[3].Window, v == 0)
	}
}
func uiMeterInitColors() unsafe.Pointer {
	m := uiMeters()
	C.dword_5d4594_1090284 = C.uint32_t(uiMeterColor(255, 0, 0))
	C.dword_5d4594_1090280 = C.uint32_t(uiMeterColor(100, 0, 0))
	*memmap.PtrUint32(0x5D4594, 1091964) = uiMeterColor(0, 255, 0)
	*memmap.PtrUint32(0x5D4594, 1092992) = uiMeterColor(0, 100, 0)
	m[0].Color, m[0].Alternate = uint32(C.dword_5d4594_1090284), uint32(C.dword_5d4594_1090280)
	m[1].Color, m[1].Alternate = uiMeterColor(0, 0, 255), uiMeterColor(0, 0, 100)
	m[4].Color, m[4].Alternate = uiMeterColor(240, 0, 240), uiMeterColor(50, 0, 50)
	m[5].Color, m[5].Alternate = uiMeterColor(255, 0, 255), uiMeterColor(50, 0, 50)
	m[6].Color, m[6].Alternate = uiMeterColor(255, 0, 255), uiMeterColor(50, 0, 50)
	for i := uintptr(0); i < 64; i++ {
		*memmap.PtrUint32(0x5D4594, 1093196+24*i) = 0
		*memmap.PtrUint32(0x5D4594, 1094732+24*i) = 0
	}
	return memmap.PtrOff(0x5D4594, 1096268)
}

func sub_470C40(v int) int {
	C.dword_5d4594_1096264 = C.uint32_t(v)
	m := &uiMeters()[0]
	if v != 0 {
		m.Color = memmap.Uint32(0x5D4594, 1091964)
		m.Alternate = memmap.Uint32(0x5D4594, 1092992)
		return int(m.Color)
	}
	m.Color, m.Alternate = uint32(C.dword_5d4594_1090284), uint32(C.dword_5d4594_1090280)
	return int(m.Alternate)
}
func uiMeterSetTotal(index int, playerOffset uintptr, current, maximum int) int {
	if p := uiMeterPlayer(); p != nil {
		*(*uint32)(unsafe.Add(p, playerOffset)) = uint32(maximum)
	}
	uiMeters()[index].Maximum = uint32(maximum)
	uiMeters()[index].Current = uint32(current)
	C.dword_5d4594_1096260 = 32
	return current
}

func nox_xxx_cliSetTotalHealth_470C80(current, maximum int) int {
	return uiMeterSetTotal(0, 2247, current, maximum)
}

func sub_470CB0(current int) int { uiMeters()[0].Current = uint32(current); return current }

//export sub_470CC0
func sub_470CC0() int { return int(uiMeters()[0].Current) }

//export sub_470CD0
func sub_470CD0() int { return int(uiMeters()[0].Maximum) }

func nox_xxx_cliSetManaAndMax_470CE0(current, maximum int) int {
	return uiMeterSetTotal(1, 2243, current, maximum)
}

func nox_xxx_cliSetMana_470D10(current int) int {
	uiMeters()[1].Current = uint32(current)
	return current
}

func sub_470D20(current, maximum int) int {
	uiMeters()[4].Current, uiMeters()[4].Maximum = uint32(current), uint32(maximum)
	if current != maximum {
		return nox_xxx_setKeybTimeout_4160D0(17)
	}
	return current
}

//export sub_470D70
func sub_470D70() { uiMeterHide(uiMeters()[5].Window, true); uiMeterHide(uiMeters()[6].Window, true) }

//export sub_470D90
func sub_470D90(current, maximum int) int {
	m := uiMeters()
	uiMeterHide(m[5].Window, false)
	uiMeterHide(m[6].Window, false)
	m[5].Current, m[5].Maximum = uint32(current), uint32(maximum)
	m[6].Current, m[6].Maximum = uint32(current), uint32(maximum)
	return current
}

//export nox_xxx_cliGetMana_470DD0
func nox_xxx_cliGetMana_470DD0() int { return int(uiMeters()[1].Current) }
func uiMeterHeartbeat() int {
	ret := uint32(C.nox_player_netCode_85319C)
	m := &uiMeters()[0]
	if ret == 0 || m.Current < 1 {
		return int(ret)
	}
	ret = uint32(3435973838) * m.Maximum
	threshold := 2 * m.Maximum / 5
	if m.Current < threshold {
		fps := GetServer().S().TickRate()
		duration := fps/3 + m.Current*((3*fps)>>2)/threshold
		*memmap.PtrUint32(0x5D4594, 1091960) = duration
		ret = 0
		if InputKeyCheckTimeoutLegacy(4, duration-1) {
			volume := int32(66)*(int32(threshold)-int32(m.Current))/int32(threshold) + 33
			Nox_xxx_clientPlaySoundSpecial_452D80(896, int(volume))
			ret = uint32(nox_xxx_setKeybTimeout_4160D0(4))
		}
	}
	return int(ret)
}
func uiMeterAdvanceCharge() int {
	m := &uiMeters()[4]
	if m.Current == m.Maximum {
		return int(m.Current)
	}
	if !Sub_416120(17) {
		return 0
	}
	n := uint32(100) / GetServer().S().TickRate()
	m.Current += n
	return int(n)
}

//export sub_4721A0
func sub_4721A0(show int) int {
	return uiMeterHide((*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1090276))), show == 0)
}
