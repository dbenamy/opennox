//go:build porttest

package server

import "unsafe"

// PortTestModifierCallback returns an existing registered callback address.
// kind is one of "damage", "defend", or "update"; this function never
// mutates registration maps.
func PortTestModifierCallback(kind, name string) unsafe.Pointer {
	var table map[string]modFuncs
	switch kind {
	case "damage":
		table = modDamageEffects
	case "defend":
		table = modDefendEffects
	case "update":
		table = modUpdateEffects
	default:
		panic("unknown modifier callback kind: " + kind)
	}
	entry, ok := table[name]
	if !ok || entry.Func == nil {
		panic("missing registered modifier callback: " + kind + "/" + name)
	}
	return entry.Func
}
