package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "client__gui__guimeter.h"
extern uint32_t dword_5d4594_1090276, dword_5d4594_1090292, dword_5d4594_1090828, dword_5d4594_1091364;
extern uint32_t dword_5d4594_1096252, dword_5d4594_1096256, dword_5d4594_1096260, dword_5d4594_1096264;
extern uint32_t dword_5d4594_1096272, dword_5d4594_1096276, dword_5d4594_1096280, dword_5d4594_1096284, dword_5d4594_1096288;
extern int nox_win_width, nox_win_height;
*/
import "C"

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"image"
	"unsafe"
)

func uiMeterWindow(parent *gui.Window, flags gui.StatusFlags, x, y, width, height int, index int, event, draw, tooltip unsafe.Pointer) *gui.Window {
	w := GetClient().Cli().GUI.NewWindowRaw(parent, flags, x, y, width, height, nil)
	w.SetAllFuncs(gui.WrapFuncC(event), gui.WrapDrawFuncC(draw), tooltip)
	if index >= 0 {
		w.WidgetData = unsafe.Pointer(uintptr(index))
	}
	return w
}
func uiMeterWindowTooltip(w *gui.Window, key string) {
	w.DrawData().SetTooltip(GetServer().S().Strings(), uiMeterString(key))
}
func uiMeterLoadImage(name string) uint32 { return uint32(uintptr(Nox_xxx_gLoadImg(name).C())) }

//export nox_win_init_cur_weapon
func nox_win_init_cur_weapon(p *C.nox_window, x, y, width, height int) {
	uiMeters()[4].Window = uiMeterWindow((*gui.Window)(unsafe.Pointer(p)), 1032, x, y, width, height, 4, C.sub_470E90, C.sub_470F40_draw, C.sub_4710B0)
}

//export sub_471160
func sub_471160(parent, x, y, width, height int) int {
	m := uiMeters()
	p := (*gui.Window)(unsafe.Pointer(uintptr(uint32(parent))))
	// Preserve creation order before installing callbacks and labels.
	m[5].Window = GetClient().Cli().GUI.NewWindowRaw(p, 1032, x, y, width, height, nil)
	m[6].Window = GetClient().Cli().GUI.NewWindowRaw(p, 1032, x-17, y-15, 15, 15, nil)
	m[5].Window.SetAllFuncs(nil, gui.WrapDrawFuncC(C.sub_471250), nil)
	uiMeterWindowTooltip(m[5].Window, "ToolTipCharges")
	m[6].Window.SetAllFuncs(nil, gui.WrapDrawFuncC(C.sub_471450), nil)
	uiMeterWindowTooltip(m[6].Window, "ToolTipCharges")
	m[5].Window.WidgetData = unsafe.Pointer(uintptr(5))
	m[6].Window.WidgetData = unsafe.Pointer(uintptr(6))
	return int(uintptr(m[5].Window.C()))
}
func uiMeterInit() int {
	player := uiMeterPlayer()
	if player == nil {
		return 0
	}
	C.dword_5d4594_1096288 = C.uint32_t(uintptr(GetClient().R2().GetFonts().FontPtrByName("small")))
	C.dword_5d4594_1096264 = 0
	C.dword_5d4594_1096256 = 0
	C.dword_5d4594_1096260 = 0
	width, height := int(C.nox_win_width), int(C.nox_win_height)
	vp := (*noxrender.Viewport)(memmap.PtrOff(0x5D4594, 1091908))
	vp.Screen = image.Rect(0, 0, width, height)
	vp.Size = image.Pt(width, height)
	vp.World.Min = image.Point{}
	for i := 0; i < 10; i++ {
		*memmap.PtrUint32(0x5D4594, 1092996+uintptr(i)*4) = uiMeterLoadImage(fmt.Sprintf("HealthMana%d", i+1))
	}
	if *memmap.PtrUint32(0x5D4594, 1096268) == 0 {
		types := &GetClient().Cli().Things
		*memmap.PtrUint32(0x5D4594, 1096268) = uint32(types.IndByID("RedPotion"))
		C.dword_5d4594_1096272 = C.uint32_t(types.IndByID("BluePotion"))
		C.dword_5d4594_1096276 = C.uint32_t(types.IndByID("CurePoisonPotion"))
		C.dword_5d4594_1096280 = C.uint32_t(types.IndByID("RedApple"))
		C.dword_5d4594_1096284 = C.uint32_t(types.IndByID("Meat"))
	}
	g := GetClient().Cli().GUI
	main := g.NewWindowRaw(nil, 136, width-91, height-201, 91, 201, nil)
	C.dword_5d4594_1090276 = C.uint32_t(uintptr(main.C()))
	uiMeterSetIcon(main, memmap.Uint32(0x5D4594, 1092996))
	cure := uiMeterWindow(main, 8, 6, 166, 28, 30, 2, C.nox_xxx_guiBottleSlotProc_471B90, C.nox_xxx_guiBottleSlotDrawFn_471A80, nil)
	C.dword_5d4594_1091364 = C.uint32_t(uintptr(cure.C()))
	uiMeterWindowTooltip(cure, "CurePoisonSlotTT")
	uiMeterSlot(2).Count = 0
	uiMeterLinkPotion(2, uint32(C.dword_5d4594_1096276))
	uiMeterSlot(2).Type = uint32(C.dword_5d4594_1096276)
	health := uiMeterWindow(main, 8, 34, 166, 28, 30, 0, C.nox_xxx_guiBottleSlotProc_471B90, C.nox_xxx_guiBottleSlotDrawFn_471A80, nil)
	C.dword_5d4594_1090292 = C.uint32_t(uintptr(health.C()))
	uiMeterWindowTooltip(health, "HealthSlotTT")
	uiMeterSlot(0).Count = 0
	uiMeterSlot(0).Image = nil
	uiMeterSlot(0).Type = 0
	m := uiMeters()
	if *(*byte)(unsafe.Add(player, 2251)) != 0 {
		mana := uiMeterWindow(main, 8, 62, 166, 28, 30, 1, C.nox_xxx_guiBottleSlotProc_471B90, C.nox_xxx_guiBottleSlotDrawFn_471A80, nil)
		C.dword_5d4594_1090828 = C.uint32_t(uintptr(mana.C()))
		uiMeterWindowTooltip(mana, "ManaSlotTT")
		uiMeterSlot(1).Count = 0
		uiMeterLinkPotion(1, uint32(C.dword_5d4594_1096272))
		uiMeterSlot(1).Type = uint32(C.dword_5d4594_1096272)
		*memmap.PtrUint32(0x5D4594, 1091900) = uiMeterLoadImage("PoisonTube")
		tubes := g.NewWindowRaw(main, 136, 0, 0, 91, 159, nil)
		uiMeterSetIcon(tubes, uiMeterLoadImage("HealthManaTubes"))
		m[1].Window = uiMeterWindow(tubes, 8, 60, 34, 25, 125, 1, C.nox_xxx_guiHealthManaTubeProc_472100, C.nox_xxx_guiHealthManaTubeDraw_471D10, nil)
		uiMeterWindowTooltip(m[1].Window, "ToolTipMana")
		m[0].Window = uiMeterWindow(tubes, 8, 34, 34, 25, 125, 0, C.nox_xxx_guiHealthManaTubeProc_472100, C.nox_xxx_guiHealthManaTubeDraw_471D10, nil)
		uiMeterWindowTooltip(m[0].Window, "ToolTipHealth")
		m[2].Window = uiMeterWindow(nil, 8, 0, 0, 0, 0, 0, nil, C.nox_xxx_drawHealthManaBar_471C00, nil)
		m[3].Window = uiMeterWindow(nil, 8, 0, 0, 0, 0, 1, nil, C.nox_xxx_drawHealthManaBar_471C00, nil)
		*memmap.PtrUint32(0x5D4594, 1093176) = 1
	} else {
		*memmap.PtrUint32(0x5D4594, 1091900) = uiMeterLoadImage("WarriorPoisonTube")
		tubes := g.NewWindowRaw(main, 136, 0, 0, 91, 159, nil)
		uiMeterSetIcon(tubes, uiMeterLoadImage("WarriorHealthTube"))
		m[0].Window = uiMeterWindow(tubes, 8, 34, 34, 25, 125, 0, C.nox_xxx_guiHealthManaTubeProc_472100, C.nox_xxx_guiHealthManaTubeDraw_471D10, nil)
		uiMeterWindowTooltip(m[0].Window, "ToolTipHealth")
		m[2].Window = uiMeterWindow(nil, 24, 0, 0, 0, 0, 0, nil, C.nox_xxx_drawHealthManaBar_471C00, nil)
	}
	uiMeterBindings()
	uiMeterInitColors()
	if noxflags.HasGame(4096) {
		C.dword_5d4594_1096252 = 1
	} else {
		C.dword_5d4594_1096252 = 0
		uiMeterHide(m[2].Window, true)
		uiMeterHide(m[3].Window, true)
	}
	return 1
}

func nox_xxx_cliPrepareGameplay2_4721D0() int {
	if w := uiMeterMainWindow(); w != nil {
		w.Destroy()
	}
	if w := uiMeters()[2].Window; w != nil {
		w.Destroy()
	}
	if w := uiMeters()[3].Window; w != nil {
		w.Destroy()
	}
	uiMeterInit()
	uiMeterRefreshPotions()
	return sub_4721A0(Nox_client_getRenderGUI())
}
