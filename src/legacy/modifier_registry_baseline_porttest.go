//go:build porttest

package legacy

/*
#include "defs.h"

void nullsub_22();
void nullsub_36();
void nullsub_38();
void nullsub_39();
void nullsub_40();
void nullsub_41();
void nullsub_42();
void nullsub_43();
void nullsub_44();

void sub_4DFE10(int a1, int a2);
float* sub_4E0370(int a1, int a2, int a3, int a4, int a5, float* a6);
float* sub_4E0380(int a1, int a2, int a3, int a4, int a5, float* a6);
float* nox_xxx_effectDamageMultiplier_4E04C0(int a1, int a2, int a3, int a4, float* a5);
void nox_xxx_attribContinualReplen_4E02C0(int a1, uint32_t* a2);
void nox_xxx_confuseEffect_4E0670(int a1, int a2, int a3, int a4);
void nox_xxx_drainMEffect_4E0740(int a1, int a2, int a3, int a4, int* a5);
void nox_xxx_sympathyEffect_4E08E0(int a1, int a2, int a3, int a4, int* a5);
int nox_xxx_effectProjectileSpeed_4E09B0(int a1, int a2, int a3, int a4, int a5);
void nox_xxx_buff_4DFD80(int a1, int a2);
void nox_xxx_checkPoisonProtectEnch_4DFDE0(int a1, int a2);
int nox_xxx_gripEffect_4E0480(int a1, int a2, int a3, int a4, int a5, int* a6);
void nox_xxx_effectRegeneration_4E01D0(int a1, int a2);
void nox_xxx_stunEffect_4E04D0(int a1, int a2, int a3, int a4);
void nox_xxx_fireEffect_4E0550(void* a1, nox_object_t* a2, nox_object_t* a3, nox_object_t* a4);
void nox_xxx_fireRingEffect_4E05B0(void* a1, nox_object_t* a2, nox_object_t* a3, nox_object_t* a4);
void nox_xxx_blueFREffect_4E05F0(void* a1, nox_object_t* a2, nox_object_t* a3, nox_object_t* a4);
void nox_xxx_recoilEffect_4E0640(int a1, int a2, int a3, int a4);
void nox_xxx_lightngEffect_4E06F0(int a1, int a2, int a3, int a4);
void nox_xxx_vampirismEffect_4E07C0(int a1, int a2, int a3, int a4, int* a5);
void nox_xxx_poisonEffect_4E0850(int a1, int a2, int a3, int a4);
int nox_xxx_inversionEffect_4E03D0(int a1, int a2, int a3, int a4, int a5, int* a6);
void sub_4DFB50(int a1, int a2);
void sub_4DFB80(int a1, int a2);
void nox_xxx_effectSpeedEngage_4DFC30(int a1, int a2);
void nox_xxx_effectSpeedDisengage_4DFCA0(int a1, int a2);
void sub_4DFD10(int a1, int a2);
void nox_xxx_modifFireProtection_4DFD40(int a1, int a2, int a3);
void sub_4DFDB0(int a1, int a2);
void sub_4E0140(int a1, int a2);
void sub_4E0170(int a1, int a2);
*/
import "C"

import "unsafe"

type PortTestModifierRegistryEntry struct {
	Symbol string
	Name   string
	Kind   string
	Parser string
	Key    unsafe.Pointer
}

func PortTestModifierRegistryEntries() []PortTestModifierRegistryEntry {
	return []PortTestModifierRegistryEntry{
		{"nox_xxx_effectDamageMultiplier_4E04C0", "DamageMultiplierEffect", "damage", "float", C.nox_xxx_effectDamageMultiplier_4E04C0},
		{"nox_xxx_stunEffect_4E04D0", "StunEffect", "damage", "int", C.nox_xxx_stunEffect_4E04D0},
		{"nox_xxx_fireEffect_4E0550", "FireEffect", "damage", "float", C.nox_xxx_fireEffect_4E0550},
		{"nox_xxx_fireRingEffect_4E05B0", "FireRingEffect", "damage", "int", C.nox_xxx_fireRingEffect_4E05B0},
		{"nox_xxx_blueFREffect_4E05F0", "BlueFireRingEffect", "damage", "int", C.nox_xxx_blueFREffect_4E05F0},
		{"nullsub_38", "FrostEffect", "damage", "int", C.nullsub_38},
		{"nox_xxx_recoilEffect_4E0640", "RecoilEffect", "damage", "float", C.nox_xxx_recoilEffect_4E0640},
		{"nox_xxx_confuseEffect_4E0670", "ConfuseEffect", "damage", "int", C.nox_xxx_confuseEffect_4E0670},
		{"nox_xxx_lightngEffect_4E06F0", "LightningEffect", "damage", "float", C.nox_xxx_lightngEffect_4E06F0},
		{"nox_xxx_drainMEffect_4E0740", "DrainManaEffect", "damage", "float", C.nox_xxx_drainMEffect_4E0740},
		{"nox_xxx_vampirismEffect_4E07C0", "VampirismEffect", "damage", "float", C.nox_xxx_vampirismEffect_4E07C0},
		{"nox_xxx_poisonEffect_4E0850", "PoisonEffect", "damage", "int", C.nox_xxx_poisonEffect_4E0850},
		{"nullsub_39", "PanicEffect", "damage", "int", C.nullsub_39},
		{"nox_xxx_sympathyEffect_4E08E0", "SympathyEffect", "damage", "float", C.nox_xxx_sympathyEffect_4E08E0},
		{"nullsub_22", "ReadinessEffect", "damage", "int", C.nullsub_22},
		{"nox_xxx_effectProjectileSpeed_4E09B0", "ProjectileSpeedEffect", "damage", "float", C.nox_xxx_effectProjectileSpeed_4E09B0},
		{"nullsub_36", "ReplenishmentEffect", "damage", "int", C.nullsub_36},
		{"sub_4E0370", "ArmorMultiplierEffect", "defend", "float", C.sub_4E0370},
		{"sub_4E0380", "DurabilityMultiplierEffect", "defend", "float", C.sub_4E0380},
		{"nullsub_40", "ResilienceEffect", "defend", "float", C.nullsub_40},
		{"nox_xxx_inversionEffect_4E03D0", "InversionEffect", "defend", "int", C.nox_xxx_inversionEffect_4E03D0},
		{"nox_xxx_gripEffect_4E0480", "GripEffect", "defend", "int", C.nox_xxx_gripEffect_4E0480},
		{"nullsub_41", "BreakingEffect", "defend", "float", C.nullsub_41},
		{"nullsub_42", "PunctureProneEffect", "defend", "float", C.nullsub_42},
		{"nox_xxx_effectRegeneration_4E01D0", "RegenerationUpdate", "update", "int", C.nox_xxx_effectRegeneration_4E01D0},
		{"nullsub_43", "ParasiteUpdate", "update", "int", C.nullsub_43},
		{"nullsub_44", "AttractionUpdate", "update", "int", C.nullsub_44},
		{"nox_xxx_attribContinualReplen_4E02C0", "ContinualReplenishmentUpdate", "update", "int", C.nox_xxx_attribContinualReplen_4E02C0},
		{"sub_4DFB50", "BrillianceEngage", "engage", "int", C.sub_4DFB50},
		{"sub_4DFB80", "BrillianceDisengage", "engage", "int", C.sub_4DFB80},
		{"nox_xxx_effectSpeedEngage_4DFC30", "SpeedEngage", "engage", "float", C.nox_xxx_effectSpeedEngage_4DFC30},
		{"nox_xxx_effectSpeedDisengage_4DFCA0", "SpeedDisengage", "engage", "float", C.nox_xxx_effectSpeedDisengage_4DFCA0},
		{"sub_4DFD10", "FireProtectEngage", "engage", "float", C.sub_4DFD10},
		{"nox_xxx_modifFireProtection_4DFD40", "FireProtectDisengage", "engage", "float", C.nox_xxx_modifFireProtection_4DFD40},
		{"nox_xxx_buff_4DFD80", "LightningProtectEngage", "engage", "float", C.nox_xxx_buff_4DFD80},
		{"sub_4DFDB0", "LightningProtectDisengage", "engage", "float", C.sub_4DFDB0},
		{"nox_xxx_checkPoisonProtectEnch_4DFDE0", "PoisonProtectEngage", "engage", "float", C.nox_xxx_checkPoisonProtectEnch_4DFDE0},
		{"sub_4DFE10", "PoisonProtectDisengage", "engage", "float", C.sub_4DFE10},
		{"sub_4E0140", "RegenerationEngage", "engage", "none", C.sub_4E0140},
		{"sub_4E0170", "RegenerationDisengage", "engage", "none", C.sub_4E0170},
	}
}
