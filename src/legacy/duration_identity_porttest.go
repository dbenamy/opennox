//go:build porttest

package legacy

/*
#include <stdint.h>
static int32_t portDurationResultValue;
static uintptr_t portDurationResultArgument;
static int portDurationResultCalls;
static uintptr_t portDurationVoidArgument;
static int portDurationVoidCalls;

static int portDurationResultObserver(void* arg) {
    portDurationResultArgument = (uintptr_t)arg;
    portDurationResultCalls++;
    return (int)portDurationResultValue;
}
static void portDurationVoidObserver(void* arg) {
    portDurationVoidArgument = (uintptr_t)arg;
    portDurationVoidCalls++;
}
static void portDurationObserverReset(int32_t result) {
    portDurationResultValue = result;
    portDurationResultArgument = 0;
    portDurationResultCalls = 0;
    portDurationVoidArgument = 0;
    portDurationVoidCalls = 0;
}
static void* portDurationResultObserverKey(void) { return (void*)portDurationResultObserver; }
static void* portDurationVoidObserverKey(void) { return (void*)portDurationVoidObserver; }
static uintptr_t portDurationResultObservedArgument(void) { return portDurationResultArgument; }
static int portDurationResultObservedCalls(void) { return portDurationResultCalls; }
static uintptr_t portDurationVoidObservedArgument(void) { return portDurationVoidArgument; }
static int portDurationVoidObservedCalls(void) { return portDurationVoidCalls; }
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

type PortTestDurationCallbackKey struct {
	Name string
	Key  unsafe.Pointer
}

// PortTestDurationCallbackKeys reports the original C function addresses used
// by duration callback slots, in the audited stable order.
func PortTestDurationCallbackKeys() []PortTestDurationCallbackKey {
	return []PortTestDurationCallbackKey{
		{Name: "nox_xxx_spellBlink2_530310", Key: Get_nox_xxx_spellBlink2_530310()},
		{Name: "nox_xxx_spellBlink1_530380", Key: Get_nox_xxx_spellBlink1_530380()},
		{Name: "sub_52F460", Key: Get_sub_52F460()},
		{Name: "nox_xxx_charmCreature1_5011F0", Key: Get_nox_xxx_charmCreature1_5011F0()},
		{Name: "nox_xxx_charmCreatureFinish_5013E0", Key: Get_nox_xxx_charmCreatureFinish_5013E0()},
		{Name: "nox_xxx_charmCreature2_501690", Key: Get_nox_xxx_charmCreature2_501690()},
		{Name: "nox_xxx_spellTurnUndeadCreate_531310", Key: Get_nox_xxx_spellTurnUndeadCreate_531310()},
		{Name: "nox_xxx_spellTurnUndeadUpdate_531410", Key: Get_nox_xxx_spellTurnUndeadUpdate_531410()},
		{Name: "nox_xxx_spellTurnUndeadDelete_531420", Key: Get_nox_xxx_spellTurnUndeadDelete_531420()},
		{Name: "nox_xxx_spellDrainMana_52E210", Key: Get_nox_xxx_spellDrainMana_52E210()},
		{Name: "nox_xxx_spellEnergyBoltStop_52E820", Key: Get_nox_xxx_spellEnergyBoltStop_52E820()},
		{Name: "nox_xxx_spellEnergyBoltTick_52E850", Key: Get_nox_xxx_spellEnergyBoltTick_52E850()},
		{Name: "nullsub_29", Key: Get_nullsub_29()},
		{Name: "nox_xxx_firewalkTick_52ED40", Key: Get_nox_xxx_firewalkTick_52ED40()},
		{Name: "sub_52EF30", Key: Get_sub_52EF30()},
		{Name: "sub_52EFD0", Key: Get_sub_52EFD0()},
		{Name: "sub_52F1D0", Key: Get_sub_52F1D0()},
		{Name: "sub_52F220", Key: Get_sub_52F220()},
		{Name: "sub_52F2E0", Key: Get_sub_52F2E0()},
		{Name: "nox_xxx_onStartLightning_52F820", Key: Get_nox_xxx_onStartLightning_52F820()},
		{Name: "nox_xxx_onFrameLightning_52F8A0", Key: Get_nox_xxx_onFrameLightning_52F8A0()},
		{Name: "sub_530100", Key: Get_sub_530100()},
		{Name: "nox_xxx_castShield1_52F5A0", Key: Get_nox_xxx_castShield1_52F5A0()},
		{Name: "sub_52F650", Key: Get_sub_52F650()},
		{Name: "sub_52F670", Key: Get_sub_52F670()},
		{Name: "nox_xxx_spellCreateMoonglow_531A00", Key: Get_nox_xxx_spellCreateMoonglow_531A00()},
		{Name: "sub_531AF0", Key: Get_sub_531AF0()},
		{Name: "nox_xxx_manaBomb_530F90", Key: Get_nox_xxx_manaBomb_530F90()},
		{Name: "nox_xxx_manaBombBoom_5310C0", Key: Get_nox_xxx_manaBombBoom_5310C0()},
		{Name: "sub_531290", Key: Get_sub_531290()},
		{Name: "nox_xxx_plasmaSmth_531580", Key: Get_nox_xxx_plasmaSmth_531580()},
		{Name: "nox_xxx_plasmaShot_531600", Key: Get_nox_xxx_plasmaShot_531600()},
		{Name: "sub_5319E0", Key: Get_sub_5319E0()},
		{Name: "sub_531490", Key: Get_sub_531490()},
		{Name: "sub_5314F0", Key: Get_sub_5314F0()},
		{Name: "sub_531560", Key: Get_sub_531560()},
		{Name: "nox_xxx_summonStart_500DA0", Key: Get_nox_xxx_summonStart_500DA0()},
		{Name: "nox_xxx_summonFinish_5010D0", Key: Get_nox_xxx_summonFinish_5010D0()},
		{Name: "nox_xxx_summonCancel_5011C0", Key: Get_nox_xxx_summonCancel_5011C0()},
		{Name: "sub_530CA0", Key: Get_sub_530CA0()},
		{Name: "sub_530D30", Key: Get_sub_530D30()},
		{Name: "nox_xxx_spellTagCreature_530160", Key: Get_nox_xxx_spellTagCreature_530160()},
		{Name: "sub_530250", Key: Get_sub_530250()},
		{Name: "sub_530270", Key: Get_sub_530270()},
		{Name: "sub_5305D0", Key: Get_sub_5305D0()},
		{Name: "sub_530650", Key: Get_sub_530650()},
		{Name: "nox_xxx_castTele_530820", Key: Get_nox_xxx_castTele_530820()},
		{Name: "sub_530880", Key: Get_sub_530880()},
		{Name: "sub_530A30_spell_execdur", Key: Get_sub_530A30_spell_execdur()},
		{Name: "nox_xxx_castTTT_530B70", Key: Get_nox_xxx_castTTT_530B70()},
		{Name: "nox_xxx_spellWallCreate_4FFA90", Key: Get_nox_xxx_spellWallCreate_4FFA90()},
		{Name: "nox_xxx_spellWallUpdate_500070", Key: Get_nox_xxx_spellWallUpdate_500070()},
		{Name: "nox_xxx_spellWallDestroy_500080", Key: Get_nox_xxx_spellWallDestroy_500080()},
	}
}

type PortTestDurationObserverState struct {
	ResultArgument uintptr
	ResultCalls    int
	VoidArgument   uintptr
	VoidCalls      int
}

// PortTestDurationObserverReset resets two independent raw C callback observers
// and returns their callable C addresses for ccall forwarding contracts.
func PortTestDurationObserverReset(result int32) (unsafe.Pointer, unsafe.Pointer) {
	C.portDurationObserverReset(C.int32_t(result))
	return C.portDurationResultObserverKey(), C.portDurationVoidObserverKey()
}

func PortTestDurationObserverSnapshot() PortTestDurationObserverState {
	return PortTestDurationObserverState{
		ResultArgument: uintptr(C.portDurationResultObservedArgument()),
		ResultCalls:    int(C.portDurationResultObservedCalls()),
		VoidArgument:   uintptr(C.portDurationVoidObservedArgument()),
		VoidCalls:      int(C.portDurationVoidObservedCalls()),
	}
}

// PortTestDurationCallResult enters the original raw C-call path and preserves
// the C int return as its exact 32-bit word.
func PortTestDurationCallResult(key unsafe.Pointer, p *server.DurSpell) int32 {
	return int32(ccall.CallIntPtr(key, p.C()))
}

// PortTestDurationCallDiscard enters the original raw void C-call path.
func PortTestDurationCallDiscard(key unsafe.Pointer, p *server.DurSpell) {
	ccall.CallVoidPtr(key, p.C())
}
