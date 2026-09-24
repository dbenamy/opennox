//go:build porttest

package legacy

import "unsafe"

type PortTestEmptyCallback struct {
	Symbol  string
	Effect  string
	Kind    string
	Arity   int
	Pointer unsafe.Pointer
}

// PortTestEmptyCallbacks reports the exact registered callback addresses in a
// stable order. It does not invoke or replace the callbacks.
func PortTestEmptyCallbacks() []PortTestEmptyCallback {
	return []PortTestEmptyCallback{
		{Symbol: "nullsub_22", Effect: "ReadinessEffect", Kind: "damage", Arity: 5, Pointer: modifierKey(modifierIDReadinessEffect)},
		{Symbol: "nullsub_36", Effect: "ReplenishmentEffect", Kind: "damage", Arity: 5, Pointer: modifierKey(modifierIDReplenishmentEffect)},
		{Symbol: "nullsub_38", Effect: "FrostEffect", Kind: "damage", Arity: 5, Pointer: modifierKey(modifierIDFrostEffect)},
		{Symbol: "nullsub_39", Effect: "PanicEffect", Kind: "damage", Arity: 5, Pointer: modifierKey(modifierIDPanicEffect)},
		{Symbol: "nullsub_40", Effect: "ResilienceEffect", Kind: "defend", Arity: 6, Pointer: modifierKey(modifierIDResilienceEffect)},
		{Symbol: "nullsub_41", Effect: "BreakingEffect", Kind: "defend", Arity: 6, Pointer: modifierKey(modifierIDBreakingEffect)},
		{Symbol: "nullsub_42", Effect: "PunctureProneEffect", Kind: "defend", Arity: 6, Pointer: modifierKey(modifierIDPunctureProneEffect)},
		{Symbol: "nullsub_43", Effect: "ParasiteUpdate", Kind: "update", Arity: 3, Pointer: modifierKey(modifierIDParasiteUpdate)},
		{Symbol: "nullsub_44", Effect: "AttractionUpdate", Kind: "update", Arity: 3, Pointer: modifierKey(modifierIDAttractionUpdate)},
		{Symbol: "nullsub_29", Effect: "EnergyBoltDestroy", Kind: "duration", Arity: 1, Pointer: Get_nullsub_29()},
	}
}
