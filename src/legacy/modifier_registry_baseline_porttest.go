//go:build porttest

package legacy

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
		{"nox_xxx_effectDamageMultiplier_4E04C0", "DamageMultiplierEffect", "damage", "float", modifierKey(modifierIDDamageMultiplierEffect)},
		{"nox_xxx_stunEffect_4E04D0", "StunEffect", "damage", "int", modifierKey(modifierIDStunEffect)},
		{"nox_xxx_fireEffect_4E0550", "FireEffect", "damage", "float", modifierKey(modifierIDFireEffect)},
		{"nox_xxx_fireRingEffect_4E05B0", "FireRingEffect", "damage", "int", modifierKey(modifierIDFireRingEffect)},
		{"nox_xxx_blueFREffect_4E05F0", "BlueFireRingEffect", "damage", "int", modifierKey(modifierIDBlueFireRingEffect)},
		{"nullsub_38", "FrostEffect", "damage", "int", modifierKey(modifierIDFrostEffect)},
		{"nox_xxx_recoilEffect_4E0640", "RecoilEffect", "damage", "float", modifierKey(modifierIDRecoilEffect)},
		{"nox_xxx_confuseEffect_4E0670", "ConfuseEffect", "damage", "int", modifierKey(modifierIDConfuseEffect)},
		{"nox_xxx_lightngEffect_4E06F0", "LightningEffect", "damage", "float", modifierKey(modifierIDLightningEffect)},
		{"nox_xxx_drainMEffect_4E0740", "DrainManaEffect", "damage", "float", modifierKey(modifierIDDrainManaEffect)},
		{"nox_xxx_vampirismEffect_4E07C0", "VampirismEffect", "damage", "float", modifierKey(modifierIDVampirismEffect)},
		{"nox_xxx_poisonEffect_4E0850", "PoisonEffect", "damage", "int", modifierKey(modifierIDPoisonEffect)},
		{"nullsub_39", "PanicEffect", "damage", "int", modifierKey(modifierIDPanicEffect)},
		{"nox_xxx_sympathyEffect_4E08E0", "SympathyEffect", "damage", "float", modifierKey(modifierIDSympathyEffect)},
		{"nullsub_22", "ReadinessEffect", "damage", "int", modifierKey(modifierIDReadinessEffect)},
		{"nox_xxx_effectProjectileSpeed_4E09B0", "ProjectileSpeedEffect", "damage", "float", modifierKey(modifierIDProjectileSpeedEffect)},
		{"nullsub_36", "ReplenishmentEffect", "damage", "int", modifierKey(modifierIDReplenishmentEffect)},
		{"sub_4E0370", "ArmorMultiplierEffect", "defend", "float", modifierKey(modifierIDArmorMultiplierEffect)},
		{"sub_4E0380", "DurabilityMultiplierEffect", "defend", "float", modifierKey(modifierIDDurabilityMultiplierEffect)},
		{"nullsub_40", "ResilienceEffect", "defend", "float", modifierKey(modifierIDResilienceEffect)},
		{"nox_xxx_inversionEffect_4E03D0", "InversionEffect", "defend", "int", modifierKey(modifierIDInversionEffect)},
		{"nox_xxx_gripEffect_4E0480", "GripEffect", "defend", "int", modifierKey(modifierIDGripEffect)},
		{"nullsub_41", "BreakingEffect", "defend", "float", modifierKey(modifierIDBreakingEffect)},
		{"nullsub_42", "PunctureProneEffect", "defend", "float", modifierKey(modifierIDPunctureProneEffect)},
		{"nox_xxx_effectRegeneration_4E01D0", "RegenerationUpdate", "update", "int", modifierKey(modifierIDRegenerationUpdate)},
		{"nullsub_43", "ParasiteUpdate", "update", "int", modifierKey(modifierIDParasiteUpdate)},
		{"nullsub_44", "AttractionUpdate", "update", "int", modifierKey(modifierIDAttractionUpdate)},
		{"nox_xxx_attribContinualReplen_4E02C0", "ContinualReplenishmentUpdate", "update", "int", modifierKey(modifierIDContinualReplenishmentUpdate)},
		{"sub_4DFB50", "BrillianceEngage", "engage", "int", modifierKey(modifierIDBrillianceEngage)},
		{"sub_4DFB80", "BrillianceDisengage", "engage", "int", modifierKey(modifierIDBrillianceDisengage)},
		{"nox_xxx_effectSpeedEngage_4DFC30", "SpeedEngage", "engage", "float", modifierKey(modifierIDSpeedEngage)},
		{"nox_xxx_effectSpeedDisengage_4DFCA0", "SpeedDisengage", "engage", "float", modifierKey(modifierIDSpeedDisengage)},
		{"sub_4DFD10", "FireProtectEngage", "engage", "float", modifierKey(modifierIDFireProtectEngage)},
		{"nox_xxx_modifFireProtection_4DFD40", "FireProtectDisengage", "engage", "float", modifierKey(modifierIDFireProtectDisengage)},
		{"nox_xxx_buff_4DFD80", "LightningProtectEngage", "engage", "float", modifierKey(modifierIDLightningProtectEngage)},
		{"sub_4DFDB0", "LightningProtectDisengage", "engage", "float", modifierKey(modifierIDLightningProtectDisengage)},
		{"nox_xxx_checkPoisonProtectEnch_4DFDE0", "PoisonProtectEngage", "engage", "float", modifierKey(modifierIDPoisonProtectEngage)},
		{"sub_4DFE10", "PoisonProtectDisengage", "engage", "float", modifierKey(modifierIDPoisonProtectDisengage)},
		{"sub_4E0140", "RegenerationEngage", "engage", "none", modifierKey(modifierIDRegenerationEngage)},
		{"sub_4E0170", "RegenerationDisengage", "engage", "none", modifierKey(modifierIDRegenerationDisengage)},
	}
}
