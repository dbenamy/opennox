package legacy

/*
#include <stdint.h>
*/
import "C"
import (
	"math"
	"unsafe"
)

//export nox_xxx_spellDrainMana_52E210
func nox_xxx_spellDrainMana_52E210(a C.float) C.int {
	return C.int(sustainedDrainMana(unsafe.Pointer(uintptr(math.Float32bits(float32(a))))))
}

//export nox_xxx_spellEnergyBoltStop_52E820
func nox_xxx_spellEnergyBoltStop_52E820(a C.int) C.int {
	return C.int(sustainedEnergyStart(unsafe.Pointer(uintptr(a))))
}

//export nox_xxx_spellEnergyBoltTick_52E850
func nox_xxx_spellEnergyBoltTick_52E850(a C.float) C.int {
	return C.int(sustainedEnergyTick(unsafe.Pointer(uintptr(math.Float32bits(float32(a))))))
}

//export nox_xxx_firewalkTick_52ED40
func nox_xxx_firewalkTick_52ED40(a *C.float) C.int {
	return C.int(sustainedFirewalk(unsafe.Pointer(a)))
}

//export sub_52EF30
func sub_52EF30(a C.int) C.int { return C.int(sustainedForceStart(unsafe.Pointer(uintptr(a)))) }

//export sub_52EFD0
func sub_52EFD0(a C.int) C.int { return C.int(sustainedForceTick(unsafe.Pointer(uintptr(a)))) }

//export sub_52F1D0
func sub_52F1D0(a C.int) C.int { return C.int(sustainedForceCancel(unsafe.Pointer(uintptr(a)))) }

//export sub_52F220
func sub_52F220(a *C.int) C.int { return C.int(sustainedGreaterHealStart(unsafe.Pointer(a))) }

//export sub_52F2E0
func sub_52F2E0(a C.float) C.int {
	return C.int(sustainedGreaterHealTick(unsafe.Pointer(uintptr(math.Float32bits(float32(a))))))
}

//export sub_52F460
func sub_52F460(a C.float) C.int {
	return C.int(sustainedChannelLife(unsafe.Pointer(uintptr(math.Float32bits(float32(a))))))
}

//export nox_xxx_castShield1_52F5A0
func nox_xxx_castShield1_52F5A0(a *C.uint32_t) C.int {
	return C.int(sustainedShieldStart(unsafe.Pointer(a)))
}

//export sub_52F650
func sub_52F650(a C.int) C.int { return C.int(sustainedShieldTick(unsafe.Pointer(uintptr(a)))) }

//export sub_52F670
func sub_52F670(a C.int) C.int { return C.int(sustainedShieldCancel(unsafe.Pointer(uintptr(a)))) }

//export nox_xxx_onStartLightning_52F820
func nox_xxx_onStartLightning_52F820(a C.int) C.int {
	return C.int(sustainedLightningStart(unsafe.Pointer(uintptr(a))))
}

//export nox_xxx_onFrameLightning_52F8A0
func nox_xxx_onFrameLightning_52F8A0(a C.float) C.int {
	return C.int(sustainedLightningTick(unsafe.Pointer(uintptr(math.Float32bits(float32(a))))))
}

//export sub_530100
func sub_530100(a *C.uint32_t) C.char { return C.char(sustainedLightningCancel(unsafe.Pointer(a))) }

//export nox_xxx_spellTagCreature_530160
func nox_xxx_spellTagCreature_530160(a *C.uint32_t) C.int {
	return C.int(sustainedTagStart(unsafe.Pointer(a)))
}

//export sub_530250
func sub_530250(a C.int) C.uint { return C.uint(sustainedTagTick(unsafe.Pointer(uintptr(a)))) }

//export sub_530270
func sub_530270(a C.int) C.int { return C.int(sustainedTagCancel(unsafe.Pointer(uintptr(a)))) }

//export nox_xxx_spellBlink2_530310
func nox_xxx_spellBlink2_530310(a *C.uint32_t) C.int {
	return C.int(sustainedBlinkStart(unsafe.Pointer(a)))
}

//export nox_xxx_spellBlink1_530380
func nox_xxx_spellBlink1_530380(a *C.int) C.int { return C.int(sustainedBlinkTick(unsafe.Pointer(a))) }

//export sub_5305D0
func sub_5305D0(a *C.uint32_t) C.int { return C.int(sustainedGlyphStart(unsafe.Pointer(a))) }

//export sub_530650
func sub_530650(a *C.int) C.int { return C.int(sustainedGlyphTick(unsafe.Pointer(a))) }

//export nox_xxx_castTele_530820
func nox_xxx_castTele_530820(a C.int) C.int {
	return C.int(sustainedTeleportStart(unsafe.Pointer(uintptr(a))))
}

//export sub_530880
func sub_530880(a *C.int) C.int { return C.int(sustainedRandomGlyphTick(unsafe.Pointer(a))) }

//export nox_xxx_castTTT_530B70
func nox_xxx_castTTT_530B70(a *C.int) C.int {
	return C.int(sustainedTeleportToPointTick(unsafe.Pointer(a)))
}

//export sub_530CA0
func sub_530CA0(a C.int) C.int { return C.int(sustainedSwapStart(unsafe.Pointer(uintptr(a)))) }

//export sub_530D30
func sub_530D30(a *C.int) C.int { return C.int(sustainedSwapTick(unsafe.Pointer(a))) }

//export nox_xxx_manaBomb_530F90
func nox_xxx_manaBomb_530F90(a *C.uint32_t) C.int {
	return C.int(sustainedManaBombStart(unsafe.Pointer(a)))
}

//export nox_xxx_manaBombBoom_5310C0
func nox_xxx_manaBombBoom_5310C0(a *C.int) C.int {
	return C.int(sustainedManaBombTick(unsafe.Pointer(a)))
}

//export sub_531290
func sub_531290(a C.int) C.int { return C.int(sustainedManaBombCancel(unsafe.Pointer(uintptr(a)))) }

//export nox_xxx_spellTurnUndeadCreate_531310
func nox_xxx_spellTurnUndeadCreate_531310(a *C.uint32_t) C.int {
	return C.int(sustainedTurnUndeadStart(unsafe.Pointer(a)))
}

//export nox_xxx_spellTurnUndeadUpdate_531410
func nox_xxx_spellTurnUndeadUpdate_531410() C.int { return C.int(sustainedTurnUndeadTick()) }

//export nox_xxx_spellTurnUndeadDelete_531420
func nox_xxx_spellTurnUndeadDelete_531420(a C.int) C.int {
	return C.int(sustainedTurnUndeadCancel(unsafe.Pointer(uintptr(a))))
}

//export sub_531490
func sub_531490(a *C.uint32_t) C.int { return C.int(sustainedOvalShieldStart(unsafe.Pointer(a))) }

//export sub_5314F0
func sub_5314F0(a C.int) C.int { return C.int(sustainedOvalShieldTick(unsafe.Pointer(uintptr(a)))) }

//export sub_531560
func sub_531560(a C.int) C.int { return C.int(sustainedOvalShieldCancel(unsafe.Pointer(uintptr(a)))) }

//export nox_xxx_plasmaSmth_531580
func nox_xxx_plasmaSmth_531580(a C.int) C.int {
	return C.int(sustainedPlasmaStart(unsafe.Pointer(uintptr(a))))
}

//export nox_xxx_plasmaShot_531600
func nox_xxx_plasmaShot_531600(a C.int) C.int {
	return C.int(sustainedPlasmaTick(unsafe.Pointer(uintptr(a))))
}

//export sub_5319E0
func sub_5319E0(a C.int) C.int { return C.int(sustainedPlasmaCancel(unsafe.Pointer(uintptr(a)))) }

//export nox_xxx_spellCreateMoonglow_531A00
func nox_xxx_spellCreateMoonglow_531A00(a *C.uint32_t) C.int {
	return C.int(sustainedMoonglowStart(unsafe.Pointer(a)))
}

//export sub_531AF0
func sub_531AF0(a C.int) C.int { return C.int(sustainedMoonglowCancel(unsafe.Pointer(uintptr(a)))) }
