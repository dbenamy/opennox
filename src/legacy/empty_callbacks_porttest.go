//go:build porttest

package legacy

/*
void nullsub_22(void);
void nullsub_29(void);
void nullsub_36(void);
void nullsub_38(void);
void nullsub_39(void);
void nullsub_40(void);
void nullsub_41(void);
void nullsub_42(void);
void nullsub_43(void);
void nullsub_44(void);
*/
import "C"
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
		{Symbol: "nullsub_22", Effect: "ReadinessEffect", Kind: "damage", Arity: 5, Pointer: C.nullsub_22},
		{Symbol: "nullsub_36", Effect: "ReplenishmentEffect", Kind: "damage", Arity: 5, Pointer: C.nullsub_36},
		{Symbol: "nullsub_38", Effect: "FrostEffect", Kind: "damage", Arity: 5, Pointer: C.nullsub_38},
		{Symbol: "nullsub_39", Effect: "PanicEffect", Kind: "damage", Arity: 5, Pointer: C.nullsub_39},
		{Symbol: "nullsub_40", Effect: "ResilienceEffect", Kind: "defend", Arity: 6, Pointer: C.nullsub_40},
		{Symbol: "nullsub_41", Effect: "BreakingEffect", Kind: "defend", Arity: 6, Pointer: C.nullsub_41},
		{Symbol: "nullsub_42", Effect: "PunctureProneEffect", Kind: "defend", Arity: 6, Pointer: C.nullsub_42},
		{Symbol: "nullsub_43", Effect: "ParasiteUpdate", Kind: "update", Arity: 3, Pointer: C.nullsub_43},
		{Symbol: "nullsub_44", Effect: "AttractionUpdate", Kind: "update", Arity: 3, Pointer: C.nullsub_44},
		{Symbol: "nullsub_29", Effect: "EnergyBoltDestroy", Kind: "duration", Arity: 1, Pointer: C.nullsub_29},
	}
}
